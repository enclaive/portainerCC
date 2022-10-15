package ra

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/internal/endpointutils"
	"github.com/rs/zerolog/log"
)

func (handler *Handler) raCoordinatorBuild(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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
	client, err := handler.dockerClientFactory.CreateClient(&localEndpoint, "", nil)
	if err != nil {
		log.Err(err)
		// panic(err)
	}

	files, err := ioutil.ReadDir("/")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		log.Info().Msg(file.Name())
	}
	tar, err := archive.Tar("/coordinator", archive.Gzip)
	if err != nil {
		panic(err)
	}
	opts := types.ImageBuildOptions{
		Dockerfile: "./dockerfile/Dockerfile.coordinator",
		Tags:       []string{"coordinator"},
	}
	res, err := client.ImageBuild(r.Context(), tar, opts)
	if err != nil {
		return httperror.InternalServerError("Unable to build Coordinator image", err)
	}
	defer res.Body.Close()
	err = print(res.Body)

	return response.JSON(w, res)
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
		log.Info().Str("Docker", scanner.Text()).Msg(scanner.Text())
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
