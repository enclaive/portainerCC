package portainercc

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/http/handler/ra"
	"github.com/portainer/portainer/api/internal/endpointutils"
	"github.com/portainer/portainer/api/internal/url"
	"github.com/rs/zerolog/log"
)

type ConfTempDeployParams struct {
	Id     portainer.ConfidentialTemplateId
	EnvId  int
	Name   string
	Inputs map[string]string
}

func (handler *Handler) deployConfidentialTemplate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params ConfTempDeployParams
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return httperror.BadRequest("request body malefomred", err)
	}

	//check if all values are set
	template, err := handler.DataStore.ConfidentialTemplate().ConfidentialTemplate(portainer.ConfidentialTemplateId(params.Id))
	if err != nil {
		return httperror.BadRequest("invalid template id", err)
	}

	for _, val := range template.Inputs {
		if _, ok := params.Inputs[val.Label]; !ok {
			return httperror.BadRequest("request body malefomred", fmt.Errorf("values missing."))
			// return httperror.BadRequest("request body malefomred", fmt.Errorf("values missing. Expected: %s ", strings.Join(template.Inputs[:], ",")))
		}
	}

	//get endpoint
	endpoint, err := handler.DataStore.Endpoint().Endpoint(portainer.EndpointID(params.EnvId))
	if err != nil {
		return httperror.InternalServerError("unable to find requested endpoint", err)
	}

	//TODO its hardcoded for docker..?!?..
	// create docker API client
	client, err := handler.DockerClientFactory.CreateClient(endpoint, "", nil)
	if err != nil {
		return httperror.InternalServerError("could not create docker client", err)
	}

	//pull image and get mr enclave mr signer
	res, err := client.ImagePull(r.Context(), template.ImageName, types.ImagePullOptions{})
	if err != nil {
		return httperror.InternalServerError("Unable to pull image", err)
	}
	defer res.Close()

	//if we dont read the res, the image would not be tagged ..
	buf := new(strings.Builder)
	_, _ = io.Copy(buf, res)
	fmt.Println(buf.String())

	//read labels
	//read pcc.mrenclave, pcc.mrsigner

	_, inspectRaw, err := client.ImageInspectWithRaw(r.Context(), template.ImageName)
	if err != nil {
		return httperror.InternalServerError("Unable to inspect image", err)
	}

	var JSON map[string]interface{}
	json.Unmarshal(inspectRaw, &JSON)

	cfg := JSON["Config"].(map[string]interface{})
	labels := cfg["Labels"].(map[string]interface{})

	mrenclave := labels["pcc.mrenclave"].(string)
	mrsigner := labels["pcc.mrsigner"].(string)

	//add to marblemanifest
	fmt.Printf("I WILL PUT THIS INTO MARBLERUN MANIFEST / Container to deploy:\n")
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("Dockerimage to use: %s:\n", template.ImageName)
	fmt.Printf("MRENCLAVE: %s (extracted from Image)\n", mrenclave)
	fmt.Printf("MRSIGNER: %s (extracted from Image)\n", mrsigner)
	fmt.Printf("Packagename %s\n", params.Name)
	fmt.Printf("Secrets:\n")
	for _, val := range template.Inputs {
		fmt.Printf("\t%s: %s\n", val.Label, params.Inputs[val.Label])
	}

	//create updateManifest
	manifest := createUpdateManifest(*template, params, mrenclave, mrsigner)

	jsonManifest, err := json.Marshal(manifest)
	if err != nil {
		return httperror.InternalServerError("Could not marshal manifest", err)
	}

	// get local docker environment
	endpoints, err := handler.DataStore.Endpoint().Endpoints()
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve environments", err)
	}
	var localEndpoint portainer.Endpoint = portainer.Endpoint{}
	for _, endpoint := range endpoints {
		if endpointutils.IsLocalEndpoint(&endpoint) {
			localEndpoint = endpoint
			log.Info().Msg(localEndpoint.URL)
		}
	}

	// check if coordinator already has a manifest
	coordinatorDeployments, err := handler.DataStore.CoordinatorDeployment().CoordinatorDeployments()
	if err != nil {
		return httperror.InternalServerError("Could not fetch coordinator Deployments from DB", err)
	}
	var coordinatorDeployment portainer.CoordinatorDeployment
	for _, deployment := range coordinatorDeployments {
		if deployment.EndpointID == params.EnvId {
			coordinatorDeployment = deployment
		}
	}

	// https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()

	}

	// encode rootCert
	rootCert := new(bytes.Buffer)
	pem.Encode(rootCert, &coordinatorDeployment.RootCert)

	if err != nil {
		return httperror.InternalServerError("failed to apply coordinator root certificate", err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(rootCert.Bytes()); !ok {
		fmt.Println("No certs appended, using system certs only")
	}

	endpointUrl, err := url.ParseURL(endpoint.URL)
	if err != nil {
		return httperror.InternalServerError("Could not parse endpoint URL", err)
	}

	userCert := new(bytes.Buffer)
	pem.Encode(userCert, &coordinatorDeployment.UserCert)

	userCertPrivKey := new(bytes.Buffer)
	pem.Encode(userCertPrivKey, &coordinatorDeployment.UserPrivateKey)

	log.Info().Msg(userCert.String())
	log.Info().Msg(userCertPrivKey.String())

	cert, _ := tls.X509KeyPair(userCert.Bytes(), userCertPrivKey.Bytes())

	tlsConfig := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cert},
	}

	cl := ra.CreateCustomClient(rootCAs, endpointUrl.Host, tlsConfig)
	resp, err := cl.Post("https://coordinator:9001/update", "application/json", bytes.NewReader(jsonManifest))
	if err != nil {
		log.Err(err)
		return httperror.InternalServerError("error request", err)
	}
	defer resp.Body.Close()

	//build secrets
	replacedStrings := map[string]map[string]string{}

	//replace secrets in files/strings secrets
	for k := range template.Secrets {
		replacedStrings[k] = make(map[string]string)
		replacedStrings[k]["Key"] = template.Secrets[k]
		for _, val := range template.Inputs {
			if val.Type == "SECRET" && val.SecretName == k {
				replacedStrings[k]["Key"] = strings.Replace(replacedStrings[k]["Key"], val.ReplacePattern, params.Inputs[val.Label], -1)
			}
		}
		//encode secret as base64
		replacedStrings[k]["Key"] = base64.StdEncoding.EncodeToString([]byte(replacedStrings[k]["Key"]))
	}

	//create Volumemapping and get file encryption keys
	mountedVolumes := []mount.Mount{}
	for _, val := range template.Inputs {
		if val.Type == "VOLUME" {
			//get pfkeyid from volume label
			_, inspectRaw, err := client.VolumeInspectWithRaw(r.Context(), params.Inputs[val.Label])
			if err != nil {
				return httperror.InternalServerError("Unable to inspect image", err)
			}

			var JSON map[string]interface{}
			json.Unmarshal(inspectRaw, &JSON)

			labels := JSON["Labels"].(map[string]interface{})
			pfKeyId, _ := strconv.Atoi(labels["pfEncryptionKeyId"].(string))

			// get key from DB
			key, err := handler.DataStore.Key().Key(portainer.KeyID(pfKeyId))
			if handler.DataStore.IsErrObjectNotFound(err) {
				//TODO error handling
				fmt.Println("key not found TODO ERROR HANDLING")
				fmt.Println(err)
			} else if err != nil {
				//TODO error handling
				fmt.Println("error TODO ERROR HANDLING")
				fmt.Println(err)
			}

			//add to secrets payload - key is already saved as base64
			replacedStrings[val.SecretName] = make(map[string]string)
			replacedStrings[val.SecretName]["Key"] = key.PFKey

			//add to docker container config
			mountedVolumes = append(mountedVolumes, mount.Mount{
				Type:   mount.TypeVolume,
				Source: params.Inputs[val.Label],
				Target: val.Label,
			})

			// replacedStrings[k]["Key"] = strings.Replace(replacedStrings[k]["Key"], val.ReplacePattern, params.Inputs[val.Label], -1)
		}
	}

	secretsJson, err := json.Marshal(replacedStrings)

	fmt.Println("SECRETS POST COORDINATOR:")
	fmt.Printf(string(secretsJson))

	// send secrets to coordinator
	secretsResp, err := cl.Post("https://coordinator:9001/secrets", "application/json", bytes.NewReader(secretsJson))
	if err != nil {
		return httperror.InternalServerError("Could not set secrets", err)
	}
	secretsResponseBody, err := ioutil.ReadAll(secretsResp.Body)
	log.Info().Msg(string(secretsResponseBody))

	defer secretsResp.Body.Close()

	coordinatorDeployment.Manifest.Packages[params.Name] = manifest.Packages[params.Name]
	coordinatorDeployment.Manifest.Marbles[params.Name+"_marble"] = manifest.Marbles[params.Name+"_marble"]

	for key, value := range manifest.Secrets {
		coordinatorDeployment.Manifest.Secrets[key] = value
	}

	err = handler.DataStore.CoordinatorDeployment().Update(coordinatorDeployment.ID, &coordinatorDeployment)
	if err != nil {
		return httperror.InternalServerError("Could not update manifest in DB", err)
	}

	//
	//
	//deploy container
	//
	//

	//port mapping
	exposedPorts := make(nat.PortSet)
	portBinding := make(nat.PortMap)

	for _, val := range template.Inputs {
		if val.Type == "PORT" {
			p, err := nat.NewPort(val.PortType, val.PortContainer)
			if err != nil {
				return httperror.InternalServerError("unable to create port", err)
			}
			exposedPorts[p] = struct{}{}

			portBinding[p] = []nat.PortBinding{
				{
					HostIP:   "",
					HostPort: params.Inputs[val.Label],
				},
			}
		}
	}

	createContainer, err := client.ContainerCreate(r.Context(),
		&container.Config{
			Image:        template.ImageName,
			ExposedPorts: exposedPorts,
			Env: []string{
				"EDG_MARBLE_TYPE=" + params.Name + "_marble",
				"EDG_MARBLE_COORDINATOR_ADDR=coordinator:2001",
				"EDG_MARBLE_DNS_NAMES=localhost,app",
			},
			Domainname: "coordinator",
		},
		&container.HostConfig{
			PortBindings:    portBinding,
			PublishAllPorts: true,
			Mounts:          mountedVolumes,
			Resources: container.Resources{
				Devices: []container.DeviceMapping{
					{
						PathOnHost:        "/dev/sgx/enclave",
						PathInContainer:   "/dev/sgx/enclave",
						CgroupPermissions: "rw",
					},
					{
						PathOnHost:        "/dev/sgx/enclave",
						PathInContainer:   "/dev/sgx_enclave",
						CgroupPermissions: "rw",
					},
					{
						PathOnHost:        "/dev/sgx_provision",
						PathInContainer:   "/dev/sgx_provision",
						CgroupPermissions: "rw",
					},
				},
			},
		},
		&network.NetworkingConfig{},
		nil,
		params.Name)

	if err != nil {
		return httperror.InternalServerError("unable to create container", err)
	}

	// connect container to coordinator network
	err = client.NetworkConnect(r.Context(), "coordinator", createContainer.ID, &network.EndpointSettings{})
	if err != nil {
		return httperror.InternalServerError("could not connect container to coordinator network", err)
	}

	// start container
	err = client.ContainerStart(r.Context(), createContainer.ID, types.ContainerStartOptions{})
	if err != nil {
		return httperror.InternalServerError("Could not start container", err)
	}

	// remove coordinator from bridge network to fix SSL_ERROR_SYSCALL error
	err = client.NetworkDisconnect(r.Context(), "bridge", createContainer.ID, false)
	if err != nil {
		return httperror.InternalServerError("could not remove container from bridge network", err)
	}

	return response.JSON(w, http.StatusOK)
}
