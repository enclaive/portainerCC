package ra

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/internal/endpointutils"
	"github.com/portainer/portainer/api/internal/url"
	"github.com/rs/zerolog/log"
)

type param struct {
	EnvironmentID int
}

func (handler *Handler) raInitManifest(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params param
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return httperror.BadRequest("request body malformed", err)
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
		Secrets: map[string]portainer.Secret{
			"app_defaultkey": {
				Type: "symmetric-key",
				Size: 128,
			},
		},
	}
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

	return response.JSON(w, http.StatusOK)

}
