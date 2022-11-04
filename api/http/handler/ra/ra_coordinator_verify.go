package ra

import (
	"bytes"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/edgelesssys/ego/eclient"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/http/client"
	"github.com/portainer/portainer/api/internal/url"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

type certQuoteResp struct {
	Cert  string
	Quote []byte
}

func (handler *Handler) raCoordinatorVerify(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	environmentID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid environment identifier route variable", err)
	}

	coordinatorDeployments, err := handler.DataStore.CoordinatorDeployment().CoordinatorDeployments()
	if err != nil {
		return httperror.InternalServerError("could not fetch coordinatorDeployments from db", err)
	}

	var coordinatorDeployment *portainer.CoordinatorDeployment
	for _, deployment := range coordinatorDeployments {
		if deployment.EndpointID == environmentID {
			coordinatorDeployment = &deployment
		}
	}
	if coordinatorDeployment == nil {
		return httperror.InternalServerError("no coordinator deployment found for requested environment", errors.New(""))
	}

	endpoint, err := handler.DataStore.Endpoint().Endpoint(portainer.EndpointID(environmentID))
	endpointUrl, err := url.ParseURL(endpoint.URL)
	// config, err := crypto.CreateTLSConfigurationFromDisk(endpoint.TLSConfig.TLSCACertPath, endpoint.TLSConfig.TLSCertPath, endpoint.TLSConfig.TLSKeyPath, endpoint.TLSConfig.TLSSkipVerify)
	// if err != nil {
	// 	return httperror.InternalServerError("", err)
	// }

	endpointUrl.Scheme = "https"

	// https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
	// Get the SystemCertPool, continue with an empty pool on error
	// localCertFile := "/coordinator/root.cert"
	// rootCAs, _ := x509.SystemCertPool()
	// if rootCAs == nil {
	// 	rootCAs = x509.NewCertPool()

	// }

	// Read in the cert file
	// certs, err := ioutil.ReadFile(localCertFile)
	// if err != nil {
	// 	return httperror.InternalServerError("failed to apply coordinator root certificate", err)
	// }

	// Append our cert to the system pool
	// if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
	// 	fmt.Println("No certs appended, using system certs only")
	// }

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: config}
	// client := &http.Client{Transport: tr}

	log.Info().Msg("hello coordinator")
	client := client.NewHTTPClient()
	client.Transport = tr
	// dialer := &net.Dialer{
	// 	Timeout:   30 * time.Second,
	// 	KeepAlive: 30 * time.Second,
	// 	// DualStack: true, // this is deprecated as of go 1.16
	// }
	// or create your own transport, there's an example on godoc.
	// tr.DialTLSContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
	// 	log.Info().Msg(endpointUrl.Host)
	// 	if addr == "coordinator:9001" {
	// 		log.Info().Msg("Hey im another address!!!")
	// 		addr = "endpointUrl.Host"
	// 	}
	// 	return dialer.DialContext(ctx, network, addr)
	// }
	resp, err := client.Get("https://20.169.251.207:4433/quote")
	if err != nil {
		log.Err(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Err(err)
	}

	var certQuoteData certQuoteResp

	quoteData := gjson.GetBytes(body, "data")
	err = json.Unmarshal([]byte(quoteData.String()), &certQuoteData)

	// https://github.com/edgelesssys/era/blob/master/era/era.go
	var certs []*pem.Block
	block, rest := pem.Decode([]byte(certQuoteData.Cert))
	if block == nil {
		return httperror.InternalServerError("could not parse certificate", err)
	}
	certs = append(certs, block)

	// If we get more than one certificate, append it to the slice
	for len(rest) > 0 {
		block, rest = pem.Decode([]byte(rest))
		if block == nil {
			return httperror.InternalServerError("could not parse certificate chain", err)
		}
		certs = append(certs, block)
	}

	rootCert := certs[len(certs)-1]

	coordinatorDeployment.RootCert = *rootCert

	report, err := eclient.VerifyRemoteReport(certQuoteData.Quote)
	if err != nil {
		return httperror.InternalServerError("could not verify remote report", err)
	}
	log.Info().Msg("Data: " + string(certQuoteData.Quote))
	log.Info().Msg("Report: " + string(report.Data))

	coordinator, err := handler.DataStore.Coordinator().Coordinator(portainer.CoordinatorID(1))
	if err != nil {
		return httperror.InternalServerError("could not fetch coordinator from db", err)
	}

	uniqueIdBytes, err := hex.DecodeString(coordinator.UniqueID)
	log.Info().Msg("uniqueID: " + hex.EncodeToString(report.UniqueID))
	log.Info().Msg("uniqueID db: " + coordinator.UniqueID)
	if !bytes.Equal(uniqueIdBytes, report.UniqueID) {
		return httperror.InternalServerError("coordinators unique id is not matching", errors.New(""))
	}

	signerIdBytes, err := hex.DecodeString(coordinator.SignerID)
	log.Info().Msg("signerID: " + hex.EncodeToString(report.SignerID))
	log.Info().Msg("signerID db: " + coordinator.SignerID)

	if !bytes.Equal(signerIdBytes, report.SignerID) {
		return httperror.InternalServerError("coordinators signer id is not matching", errors.New(""))
	}

	coordinatorDeployment.Verified = true
	err = handler.DataStore.CoordinatorDeployment().Update(coordinatorDeployment.ID, coordinatorDeployment)
	if err != nil {
		return httperror.InternalServerError("could not update coordinator deployment in db", err)
	}

	fmt.Println(string(body), err)
	return response.JSON(w, coordinatorDeployment)
}
