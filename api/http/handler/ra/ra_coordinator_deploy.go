package ra

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/internal/endpointutils"
	"github.com/rs/zerolog/log"
)

type DeployCoordinatorParams struct {
	CoordinatorID int `json:"coordinatorId"`
	EnvironmentID int `json:"environmentId"`
}

func (handler *Handler) raCoordinatorDeploy(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	var params DeployCoordinatorParams
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "request body maleformed", err}
	}

	// get coordinator
	coordinator, err := handler.DataStore.Coordinator().Coordinator(portainer.CoordinatorID(params.CoordinatorID))
	if err != nil {
		return httperror.InternalServerError("Could not retrieve coordinator from database", err)
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

	// create local docker API client
	localClient, err := handler.dockerClientFactory.CreateClient(&localEndpoint, "", nil)
	if err != nil {
		log.Err(err)
		// panic(err)
	}

	// get requested coordinator image
	image, err := localClient.ImageSave(r.Context(), []string{coordinator.ImageID})
	defer image.Close()
	localClient.Close()

	// get target environment
	targetEndpoint, err := handler.DataStore.Endpoint().Endpoint(portainer.EndpointID(params.EnvironmentID))
	log.Info().Msg(targetEndpoint.Name)
	if err != nil {
		return httperror.InternalServerError("unable to find requested endpoint", err)
	}

	// create target docker API client
	targetClient, err := handler.dockerClientFactory.CreateClient(targetEndpoint, "", nil)
	if err != nil {
		return httperror.InternalServerError("unable to create docker client", err)
	}

	// create coordinator image on target environment
	_, err = targetClient.ImageLoad(r.Context(), image, false)
	if err != nil {
		return httperror.InternalServerError("Unable to build Coordinator image", err)
	}

	port1, err := nat.NewPort("tcp", "9944")
	port2, err := nat.NewPort("tcp", "4433")

	// run coordinator image on target environment
	createdBody, err := targetClient.ContainerCreate(r.Context(),
		&container.Config{
			Image: "coordinator/" + coordinator.Name,
			ExposedPorts: nat.PortSet{
				port1: struct{}{},
				port2: struct{}{},
			},
			Env: []string{
				"OE_SIMULATION=0",
				"OE_LOG_LEVEL=INFO",
				"EDG_COORDINATOR_MESH_ADDR=coordinator:2001",
				"EDG_COORDINATOR_CLIENT_ADDR=coordinator:4433",
				"EDG_COORDINATOR_DNS_NAMES=coordinator",
				"EDG_COORDINATOR_PROMETHEUS_ADDR=0.0.0.0:9944",
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				port1: []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "9944",
					},
				},
				port2: []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "4433",
					},
				},
			},
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
		coordinator.Name)

	if err != nil {
		return httperror.InternalServerError("unable to create coordinator container", err)
	}

	res := targetClient.ContainerStart(r.Context(), createdBody.ID, types.ContainerStartOptions{})

	return response.JSON(w, res)
}
