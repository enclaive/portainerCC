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
	Name          string
	UniqueID      string
	Username      string
	Password      string
}

// Adds a new service to the manifest of a running coordinator
func (handler *Handler) raServiceAdd(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params ServiceAddParams
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
	var manifest portainer.CoordinatorManifest = portainer.CoordinatorManifest{}

	// if coordinator has no manifest, create an initial manifest with the requested Marbles and secrets, if coordinator already has a manifest, create an update manifest
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
		manifest, secrets := createManifestMariadb(params.UniqueID, params.Username, params.Password, params.Name, true)
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

		// create secrets
		secretsBodyJson, err := json.Marshal(secrets)

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
		return response.JSON(w, http.StatusOK)

	} else {

		// create update manifest
		manifest, secrets := createManifestMariadb(params.UniqueID, params.Username, params.Password, params.Name, false)

		jsonManifest, err := json.Marshal(manifest)
		if err != nil {
			return httperror.InternalServerError("Could not marshal manifest", err)
		}
		log.Info().Msg(string(jsonManifest))

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
		log.Info().Msg(string(body))

		// create secrets
		secretsBodyJson, err := json.Marshal(secrets)

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

		// add update to manifest in DB
		coordinatorDeployment.Manifest.Packages[params.Name] = manifest.Packages[params.Name]
		coordinatorDeployment.Manifest.Marbles[params.Name+"_marble"] = manifest.Marbles[params.Name+"_marble"]

		for key, value := range manifest.Secrets {
			coordinatorDeployment.Manifest.Secrets[key] = value
		}

		err = handler.DataStore.CoordinatorDeployment().Update(coordinatorDeployment.ID, &coordinatorDeployment)
		if err != nil {
			return httperror.InternalServerError("Could not update manifest in DB", err)
		}
		return response.JSON(w, http.StatusOK)
	}

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
