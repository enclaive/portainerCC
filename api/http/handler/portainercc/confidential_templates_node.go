package portainercc

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/http/handler/ra"
	"github.com/portainer/portainer/api/internal/endpointutils"
	"github.com/portainer/portainer/api/internal/url"
)

func (handler *Handler) runConfidentialNode(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params ConfImgDeployParams
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return httperror.BadRequest("request body malefomred", err)
	}

	fmt.Println("Body:")
	fmt.Println(params)

	const BASE_IMG = "soelangen/pcc-nodejs-demo "

	// clone repo in tmpdir on host
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return httperror.InternalServerError("error while cloning git repo", err)
	}

	fmt.Println(tmpDir)

	_, err = git.PlainClone(tmpDir+"/app", false, &git.CloneOptions{
		URL:      params.Repository,
		Progress: os.Stdout,
	})

	//build docker
	// read signing key from db and convert it to pem format
	key, err := handler.DataStore.Key().Key(portainer.KeyID(params.SigningKeyId))
	if err != nil {
		return httperror.InternalServerError("could not fetch signing key from db", err)
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(key.SigningKey)
	keyPEM :=
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		}

	var PrivateKeyRow bytes.Buffer
	err = pem.Encode(&PrivateKeyRow, keyPEM)

	signingKey := PrivateKeyRow.String()

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

	//get target endpoint
	endpoint, err := handler.DataStore.Endpoint().Endpoint(portainer.EndpointID(params.EnvId))
	if err != nil {
		return httperror.InternalServerError("unable to find requested endpoint", err)
	}

	// create docker API client
	client, err := handler.DockerClientFactory.CreateClient(&localEndpoint, "", nil)
	if err != nil {
		return httperror.InternalServerError("could not create docker client", err)
	}

	//build dockerfile
	dfile, err := os.Create(tmpDir + "/Dockerfile")
	defer dfile.Close()

	/////////////////////////////////////////DOCKERFILE

	dfile.WriteString("FROM " + BASE_IMG + "\n")
	dfile.WriteString("\n")
	//TODO use buildkit secrets
	dfile.WriteString("ARG signingkey\n")
	dfile.WriteString("\n")
	dfile.WriteString("RUN echo \"$signingkey\" > /signing.pem\n")
	dfile.WriteString("\n")
	dfile.WriteString("COPY ./app /app/\n")
	dfile.WriteString("\n")

	//handle build args

	dfile.WriteString("RUN gramine-manifest -Darch_libdir=/lib/x86_64-linux-gnu node.manifest.template node.manifest \\\n")
	dfile.WriteString("&& gramine-sgx-sign --key \"/signing.pem\" --manifest node.manifest --output node.manifest.sgx \\\n")
	dfile.WriteString("&& gramine-sgx-get-token -s ./node.sig -o attributes \\\n")
	dfile.WriteString("&& cat ./attributes \\\n")
	dfile.WriteString("&& sed -i 's,https://localhost:8081/sgx/certification/v3/,https://172.17.0.1:8081/sgx/certification/v3/,g' /etc/sgx_default_qcnl.conf \\\n")
	dfile.WriteString("&& sed -i 's,\"use_secure_cert\": true,\"use_secure_cert\": false,g' /etc/sgx_default_qcnl.conf\\\n")

	///////////////////////////////////////////////////
	dfile.Sync()

	// archive repo and dockerfile
	tar, err := archive.Tar(tmpDir, archive.Gzip)
	if err != nil {
		panic(err)
	}

	imgName := "sgxdcaprastuff/pcc-node-demo:" + params.Name

	// set build options for image build
	opts := types.ImageBuildOptions{
		Dockerfile: "./Dockerfile",
		Tags:       []string{imgName},
		BuildArgs:  map[string]*string{"signingkey": &signingKey},
		Outputs: []types.ImageBuildOutput{
			{Type: "local"},
		},
		NoCache: true,
	}

	// send image build request
	res, err := client.ImageBuild(r.Context(), tar, opts)
	if err != nil {
		return httperror.InternalServerError("Unable to build Coordinator image", err)
	}
	defer res.Body.Close()

	var lastLine string

	var mrenclave string
	var mrsigner string

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		lastLine = scanner.Text()
		log.Info().Str("Docker", "").Msg(scanner.Text())
		if strings.Contains(lastLine, "mr_enclave") {
			split := strings.Split(lastLine, ",")
			for _, line := range split {
				fmt.Println(line)
				if strings.Contains(line, "mr_enclave") {
					uniqueID := strings.Split(line, "mr_enclave:")[1]
					uniqueID = strings.Split(uniqueID, "\\n")[0]
					mrenclave = strings.TrimSpace(uniqueID)
				}
				if strings.Contains(line, "mr_signer") {
					signerID := strings.Split(line, "mr_signer:")[1]
					signerID = strings.Split(signerID, "\\n")[0]
					mrsigner = strings.TrimSpace(signerID)
				}
			}
		}
	}

	fmt.Println(mrenclave)
	fmt.Println(mrsigner)

	// var params ConfTempDeployParams
	marbleparams := &portainer.Parameters{
		Argv: []string{"/usr/bin/node", params.RunArgs},
	}

	//create updateManifest
	manifest := createUpdateManifestForImage(params.Name, *marbleparams, mrenclave, mrsigner)

	fmt.Println("")
	fmt.Println("Manifest:")
	b, _ := json.MarshalIndent(manifest, "", "  ")
	fmt.Println(string(b))
	fmt.Println("")

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

	coordinatorURLEndpoint := "update"
	// if manifest in db is empty, create initial manifest + the deployment params
	if reflect.DeepEqual(coordinatorDeployment.Manifest, (portainer.CoordinatorManifest{})) {

		coordinatorURLEndpoint = "manifest"

		//create user/portainer cert to be able to update the coordinator later
		userCertPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return httperror.InternalServerError("unable to create user certificate private key", err)
		}

		userCertBytes, err := ra.CreateUserCert(userCertPrivKey)
		if err != nil {
			return httperror.InternalServerError("Could not create user certificate", err)
		}

		userCertPEM := new(bytes.Buffer)
		pem.Encode(userCertPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: userCertBytes,
		})

		userCertPrivKeyPEM := new(bytes.Buffer)
		pem.Encode(userCertPrivKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(userCertPrivKey),
		})

		//add to coordinator db object - saved later
		block, _ := pem.Decode(userCertPrivKeyPEM.Bytes())
		coordinatorDeployment.UserPrivateKey = *block
		block, _ = pem.Decode(userCertPEM.Bytes())
		coordinatorDeployment.UserCert = *block

		//add initial to manifest
		manifest.Users = map[string]portainer.CoordinatorUser{
			"portainer": {
				Certificate: userCertPEM.String(),
				Roles: []string{
					"updatePackage",
					"secretManager",
				},
			}}

		manifest.Roles = map[string]portainer.CoordinatorRole{
			"updatePackage": {
				ResourceType: "Packages",
				Actions:      []string{"UpdateSecurityVersion"},
			},
			"secretManager": {
				ResourceType:  "Secrets",
				ResourceNames: []string{},
				Actions: []string{
					"ReadSecret",
					"WriteSecret",
				},
			},
		}

		coordinatorDeployment.Manifest = manifest
	}

	//parse manifest
	jsonManifest, err := json.Marshal(manifest)
	if err != nil {
		return httperror.InternalServerError("Could not marshal manifest", err)
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

	resp, err := cl.Post("https://coordinator:9001/"+coordinatorURLEndpoint, "application/json", bytes.NewReader(jsonManifest))
	if err != nil {
		log.Err(err)
		return httperror.InternalServerError("error request", err)
	}
	defer resp.Body.Close()

	//update db manifest

	coordinatorDeployment.Manifest.Packages[params.Name] = manifest.Packages[params.Name]
	coordinatorDeployment.Manifest.Marbles[params.Name+"_marble"] = manifest.Marbles[params.Name+"_marble"]

	err = handler.DataStore.CoordinatorDeployment().Update(coordinatorDeployment.ID, &coordinatorDeployment)
	if err != nil {
		return httperror.InternalServerError("Could not update manifest in DB", err)
	}

	//push image
	authConfig := types.AuthConfig{
		Username:      "sgxdcaprastuff",
		Password:      "dckr_pat_ASB6_d6hVfhgHsNXByxEWYjfXtc",
		ServerAddress: "https://index.docker.io/v2/",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
	pushOptions := types.ImagePushOptions{RegistryAuth: authConfigEncoded}

	pushRes, err := client.ImagePush(r.Context(), imgName, pushOptions)
	if err != nil {
		return httperror.InternalServerError("could not push coordinator image to registry", err)
	}
	defer pushRes.Close()
	print(pushRes)

	//execute img
	//get target endpoint docker client
	targetClient, err := handler.DockerClientFactory.CreateClient(endpoint, "", nil)
	if err != nil {
		return httperror.InternalServerError("could not create docker client", err)
	}
	// pull coordinator image
	pullOptions := types.ImagePullOptions{RegistryAuth: authConfigEncoded}

	pullRes, err := targetClient.ImagePull(r.Context(), imgName, pullOptions)
	if err != nil {
		return httperror.InternalServerError("could not pull coordinator image to registry", err)
	}
	defer pullRes.Close()
	print(pullRes)

	//port mapping
	exposedPorts := make(nat.PortSet)
	portBinding := make(nat.PortMap)

	for _, val := range params.Ports {
		p, err := nat.NewPort(val.Type, val.Container)
		if err != nil {
			return httperror.InternalServerError("unable to create port", err)
		}
		exposedPorts[p] = struct{}{}

		portBinding[p] = []nat.PortBinding{
			{
				HostIP:   "",
				HostPort: val.Host,
			},
		}
	}

	createContainer, err := client.ContainerCreate(r.Context(),
		&container.Config{
			Image:        imgName,
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
	err = targetClient.NetworkConnect(r.Context(), "coordinator", createContainer.ID, &network.EndpointSettings{})
	if err != nil {
		return httperror.InternalServerError("could not connect container to coordinator network", err)
	}

	// start container
	err = targetClient.ContainerStart(r.Context(), createContainer.ID, types.ContainerStartOptions{})
	if err != nil {
		return httperror.InternalServerError("Could not start container", err)
	}

	// remove coordinator from bridge network to fix SSL_ERROR_SYSCALL error
	err = targetClient.NetworkDisconnect(r.Context(), "bridge", createContainer.ID, false)
	if err != nil {
		return httperror.InternalServerError("could not remove container from bridge network", err)
	}

	return response.JSON(w, createContainer.ID)
}
