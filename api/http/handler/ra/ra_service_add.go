package ra

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"time"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/http/client"
	"github.com/portainer/portainer/api/internal/endpointutils"
	"github.com/portainer/portainer/api/internal/url"
	"github.com/rs/zerolog/log"
)

type ServiceAddParams struct {
	EnvironmentID int
	SignerID      string
	UniqueID      string
	ImageID       string
	Env           map[string]string
	Files         map[string]portainer.File
	Secrets       map[string]string
}

func (handler *Handler) raServiceAdd(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params ServiceAddParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return httperror.BadRequest("request body malformed", err)
	}

	for key, value := range params.Secrets {
		log.Info().Msg(key + ": " + value)

	}

	// get target endpoint
	endpoint, err := handler.DataStore.Endpoint().Endpoint(portainer.EndpointID(params.EnvironmentID))
	if err != nil {
		return httperror.InternalServerError("could not fetch endpoint from db", err)
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

	// create docker API client
	localClient, err := handler.dockerClientFactory.CreateClient(&localEndpoint, "", nil)
	if err != nil {
		log.Err(err)
		// panic(err)
	}

	_, err = localClient.Ping(r.Context())
	if err != nil {
		return httperror.InternalServerError("Could not ping docker env", err)
	}

	imageInspect, _, err := localClient.ImageInspectWithRaw(r.Context(), params.ImageID)
	if err != nil {
		return httperror.InternalServerError("Could not fetch image information", err)
	}

	for index, tag := range imageInspect.RepoTags {
		log.Info().Msg(strconv.FormatInt(int64(index), 10) + " " + tag)
	}
	log.Info().Msg("length of tags: " + strconv.FormatInt(int64(len(imageInspect.RepoTags)), 10))

	// check if coordinator already has a manifest
	coordinatorDeployments, err := handler.DataStore.CoordinatorDeployment().CoordinatorDeployments()
	if err != nil {
		return httperror.InternalServerError("Could not fetch coordinator Deployments from DB", err)
	}
	var coordinatorDeployment portainer.CoordinatorDeployment
	for _, deployment := range coordinatorDeployments {
		if deployment.EndpointID == params.EnvironmentID {
			coordinatorDeployment = deployment
		}
	}
	var manifest portainer.CoordinatorManifest = portainer.CoordinatorManifest{}

	// name of the package in manifest
	var packageName string
	if len(imageInspect.RepoTags) == 1 {
		packageName = imageInspect.RepoTags[0]
	} else {
		packageName = imageInspect.RepoTags[len(imageInspect.RepoTags)-1]
	}

	// parse Secrets
	var manifestSecrets = map[string]portainer.Secret{}
	var secrets = map[string]map[string]string{}
	for key, value := range params.Secrets {
		manifestSecrets[key] = portainer.Secret{
			Type:        "plain",
			UserDefined: true,
		}
		secrets[key] = make(map[string]string)
		secrets[key]["Key"] = value
	}

	params.Files["/dev/attestation/keys/default"] = portainer.File{
		Data:        "{{ raw .Secrets.app_defaultkey.Private }}",
		Encoding:    "string",
		NoTemplates: false,
	}
	manifestSecrets["app_defaultkey"] = portainer.Secret{
		Type: "symmetric-key",
		Size: 128,
	}

	if reflect.DeepEqual(coordinatorDeployment.Manifest, manifest) {
		log.Info().Msg("No manifest")

		//create User certificate
		userCertPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return httperror.InternalServerError("unable to create user certificate private key", err)
		}

		userCertBytes, err := CreateUserCert(userCertPrivKey)
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

		// add user cert and private key to db
		block, _ := pem.Decode(userCertPrivKeyPEM.Bytes())
		coordinatorDeployment.UserPrivateKey = *block
		block, _ = pem.Decode(userCertPEM.Bytes())
		coordinatorDeployment.UserCert = *block

		log.Info().Msg(userCertPEM.String())
		log.Info().Msg(userCertPrivKeyPEM.String())

		// create coordinator manifest
		manifest := portainer.CoordinatorManifest{
			Users: map[string]portainer.CoordinatorUser{
				"portainer": {
					Certificate: userCertPEM.String(),
					Roles: []string{
						"updatePackage",
						"secretManager",
					},
				},
			},
			Packages: map[string]portainer.PackageProperties{
				packageName: {
					UniqueID: params.UniqueID,
					// ProductID:       1,
					// SecurityVersion: 1,
				},
			},
			Marbles: map[string]portainer.Marble{
				"app_marble": {
					Package: packageName,
					Parameters: portainer.Parameters{
						Files: params.Files,
						//  map[string]portainer.File{
						// "/dev/attestation/keys/default": {
						// 	Data:        "{{ raw .Secrets.app_defaultkey.Private }}",
						// 	Encoding:    "string",
						// 	NoTemplates: false,
						// },
						// "/app/init.sql": {
						// 	Data:        "{{ raw .Secrets.init.Private }}",
						// 	Encoding:    "string",
						// 	NoTemplates: false,
						// },
						// },
						Env: params.Env,
						Argv: []string{
							"/app/mariadbd",
							"--init-file=/app/init.sql",
						},
					},
				},
			},
			Secrets: manifestSecrets,
			// map[string]portainer.Secret{
			// 	"app_defaultkey": {
			// 		Type: "symmetric-key",
			// 		Size: 128,
			// 	},
			// "password": {
			// 	Type:        "plain",
			// 	UserDefined: true,
			// },
			// 	"init": {
			// 		Type:        "plain",
			// 		UserDefined: true,
			// 	},
			// },
			Roles: map[string]portainer.CoordinatorRole{
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
			},
		}

		jsonManifest, err := json.Marshal(manifest)
		log.Info().Msg(string(jsonManifest))

		// add coordinator manifest to db
		coordinatorDeployment.Manifest = manifest
		err = handler.DataStore.CoordinatorDeployment().Update(coordinatorDeployment.ID, &coordinatorDeployment)
		if err != nil {
			return httperror.InternalServerError("Could not update coordinator manifest in DB", err)
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

		tlsConfig := &tls.Config{
			RootCAs: rootCAs,
		}

		// create custom http client
		client := CreateCustomClient(rootCAs, endpointUrl.Host, tlsConfig)

		resp, err := client.Post("https://coordinator:9001/manifest", "application/json", bytes.NewReader(jsonManifest))
		if err != nil {
			log.Err(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Err(err)
		}
		log.Info().Msg(string(body))

		// create secrets
		secretsBody := secrets
		// map[string]map[string]string{
		// 	"init": {
		// 		"Key": "Q1JFQVRFIE9SIFJFUExBQ0UgVVNFUiByb290IElERU5USUZJRUQgQlkgJ3Jvb3QnOwoJCQkJR1JBTlQgQUxMIE9OICouKiBUTyByb290IFdJVEggR1JBTlQgT1BUSU9OOw==",
		// 	},
		// }

		secretsBodyJson, err := json.Marshal(secretsBody)

		cert, _ := tls.X509KeyPair(userCertPEM.Bytes(), userCertPrivKeyPEM.Bytes())

		tlsConfig = &tls.Config{
			RootCAs:      rootCAs,
			Certificates: []tls.Certificate{cert},
		}

		client = CreateCustomClient(rootCAs, endpointUrl.Host, tlsConfig)

		// send secrets to coordinator
		secretsResp, err := client.Post("https://coordinator:9001/secrets", "application/json", bytes.NewReader(secretsBodyJson))
		if err != nil {
			return httperror.InternalServerError("Could not set secrets", err)
		}
		secretsResponseBody, err := ioutil.ReadAll(secretsResp.Body)
		log.Info().Msg(string(secretsResponseBody))

		defer secretsResp.Body.Close()

		return response.JSON(w, string(secretsResponseBody))

	} else {

		// create update manifest
		manifest := portainer.CoordinatorManifest{
			Packages: map[string]portainer.PackageProperties{
				packageName: {
					UniqueID: params.UniqueID,
					// ProductID:       1,
					// SecurityVersion: 1,
				},
			},
			Marbles: map[string]portainer.Marble{
				"app_marble": {
					Package: packageName,
					Parameters: portainer.Parameters{
						Files: map[string]portainer.File{
							"/dev/attestation/keys/default": {
								Data:        "{{ raw .Secrets.app_defaultkey.Private }}",
								Encoding:    "string",
								NoTemplates: false,
							},
						},
						Env: params.Env,
					},
				},
			},
		}

		jsonManifest, err := json.Marshal(manifest)

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

		// create custom http client
		client := CreateCustomClient(rootCAs, endpointUrl.Host, tlsConfig)

		resp, err := client.Post("https://coordinator:9001/update", "application/json", bytes.NewReader(jsonManifest))
		if err != nil {
			log.Err(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Err(err)
		}

		// TODO add update to manifest in DB

		return response.JSON(w, string(body))
	}
	// TODO if no manifest present, create x.509 certificate, create manifest from request params and user, post it and save cert and key to db

	// TODO if manifest present, build update manifest from request params, get user cert and key from db and post update manifest
}

func CreateUserCert(key *rsa.PrivateKey) ([]byte, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	userCertBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	return userCertBytes, nil
}

func CreateCustomClient(rootCAs *x509.CertPool, endpointUrl string, tlsConfig *tls.Config) client.HTTPClient {
	tr := &http.Transport{TLSClientConfig: tlsConfig}

	client := client.NewHTTPClient()
	client.Transport = tr
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		log.Info().Msg(endpointUrl)
		if addr == "coordinator:9001" {
			log.Info().Msg("Ich bin eine andere Adresse")
			addr = endpointUrl
		}
		return dialer.DialContext(ctx, network, addr)
	}
	return *client
}
