package ra

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
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

	fileContent, err := ioutil.ReadFile("/coordinator/build/private.pem")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	//:TODO read signing key from db and convert it to pem format
	key, err := rsa.GenerateKey(rand.Reader, 3072)
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPEM :=
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		}

	var PrivateKeyRow bytes.Buffer
	err = pem.Encode(&PrivateKeyRow, keyPEM)

	// signingKey := PrivateKeyRow.String()
	signingKey := string(fileContent)

	// archive coordinator source code
	tar, err := archive.Tar("/coordinator", archive.Gzip)
	if err != nil {
		panic(err)
	}

	// set build options for image build
	opts := types.ImageBuildOptions{
		Dockerfile: "./dockerfile/Dockerfile.coordinator",
		Tags:       []string{"coordinator/" + params.Name},
		BuildArgs:  map[string]*string{"signingkey": &signingKey},
		Outputs:    []types.ImageBuildOutput{},
	}

	// send image build request
	res, err := client.ImageBuild(r.Context(), tar, opts)
	if err != nil {
		return httperror.InternalServerError("Unable to build Coordinator image", err)
	}
	defer res.Body.Close()
	err = print(res.Body)

	// get image id of built image
	imgMeta, _, err := client.ImageInspectWithRaw(r.Context(), "coordinator/"+params.Name)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve new coordinators image id", err)
	}

	// TODO extract MRENCLAVE and MRSIGNER

	// create new coordinator in database
	coordinatorObject := &portainer.Coordinator{
		Name:         params.Name,
		SigningKeyID: params.SigningKeyId,
		ImageID:      strings.Split(imgMeta.ID, ":")[1],
	}
	err = handler.DataStore.Coordinator().Create(coordinatorObject)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to generate new coordinator", err}
	}
	return response.JSON(w, coordinatorObject)
}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
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
