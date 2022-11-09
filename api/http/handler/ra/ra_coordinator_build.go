package ra

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/internal/endpointutils"
	"github.com/rs/zerolog/log"
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type PostCoordinatorParams struct {
	Name         string
	SigningKeyId int
}

func (handler *Handler) raCoordinatorBuild(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params PostCoordinatorParams
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "request body maleformed", err}
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
	client, err := handler.dockerClientFactory.CreateClient(&localEndpoint, "", nil)
	if err != nil {
		log.Err(err)
		// panic(err)
	}

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
	// signingKey := string(fileContent)

	// archive coordinator source code
	tar, err := archive.Tar("/coordinator", archive.Gzip)
	if err != nil {
		panic(err)
	}

	// set build options for image build
	opts := types.ImageBuildOptions{
		Dockerfile: "./dockerfile/Dockerfile.coordinator",
		Tags:       []string{"sgxdcaprastuff/coordinatortest"},
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

	coordinatorObject := &portainer.Coordinator{
		Name:         params.Name,
		SigningKeyID: params.SigningKeyId,
	}

	// print(res.Body)
	// extract UniqueID and SignerID from Build Logs
	// scanner := bufio.NewScanner(res.Body)

	var lastLine string

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		lastLine = scanner.Text()
		log.Info().Str("Docker", "").Msg(scanner.Text())
		if strings.Contains(lastLine, "UniqueID") {
			split := strings.Split(lastLine, ",")
			for _, line := range split {
				fmt.Println(line)
				if strings.Contains(line, "UniqueID") {
					uniqueID := strings.Split(line, ":")[1]
					uniqueID = strings.ReplaceAll(uniqueID, `\"`, "")
					uniqueID = strings.ReplaceAll(uniqueID, ` `, "")
					coordinatorObject.UniqueID = uniqueID
				}
				if strings.Contains(line, "SignerID") {
					signerID := strings.Split(line, ":")[1]
					signerID = strings.ReplaceAll(signerID, `\"`, "")
					signerID = strings.ReplaceAll(signerID, ` `, "")
					coordinatorObject.SignerID = signerID
				}
			}
		}
	}

	// return nil

	// var lastLine string
	// for scanner.Scan() {
	// 	lastLine = scanner.Text()
	// 	if strings.Contains(lastLine, "UniqueID") {
	// 		split := strings.Split(lastLine, ",")
	// 		for _, line := range split {
	// 			fmt.Println(line)
	// 			if strings.Contains(line, "UniqueID") {
	// 				uniqueID := strings.Split(line, ":")[1]
	// 				uniqueID = strings.ReplaceAll(uniqueID, `\"`, "")
	// 				uniqueID = strings.ReplaceAll(uniqueID, ` `, "")
	// 				coordinatorObject.UniqueID = uniqueID
	// 			}
	// 			if strings.Contains(line, "SignerID") {
	// 				signerID := strings.Split(line, ":")[1]
	// 				signerID = strings.ReplaceAll(signerID, `\"`, "")
	// 				signerID = strings.ReplaceAll(signerID, ` `, "")
	// 				coordinatorObject.SignerID = signerID
	// 			}
	// 		}
	// 	}
	// }

	//push image
	authConfig := types.AuthConfig{
		Username:      "sgxdcaprastuff",
		Password:      "dckr_pat_ASB6_d6hVfhgHsNXByxEWYjfXtc",
		ServerAddress: "https://index.docker.io/v2/",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
	pushOptions := types.ImagePushOptions{RegistryAuth: authConfigEncoded}

	pushRes, err := client.ImagePush(r.Context(), "sgxdcaprastuff/coordinatortest", pushOptions)
	if err != nil {
		return httperror.InternalServerError("could not push coordinator image to registry", err)
	}
	defer pushRes.Close()
	print(pushRes)

	// get image id of built image
	imgMeta, _, err := client.ImageInspectWithRaw(r.Context(), "sgxdcaprastuff/coordinatortest")
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve new coordinators image id", err)
	}
	coordinatorObject.ImageID = strings.Split(imgMeta.ID, ":")[1]

	// create new coordinator in database
	err = handler.DataStore.Coordinator().Create(coordinatorObject)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to generate new coordinator", err}
	}
	return response.JSON(w, coordinatorObject)
}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		log.Info().Str("Docker", "").Msg(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
