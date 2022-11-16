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
	"github.com/rs/zerolog/log"
)

type ServiceDeployParams struct {
	EnvironmentID int
	ImageID       string
	Name          string // HAS to be the same name as the package in the coordinator manifest
}

// Deploys Gramine Application Container to an environment with a running coordinator
func (handler *Handler) raServiceDeploy(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params ServiceDeployParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return httperror.BadRequest("request body malformed", err)
	}

	// get target endpoint
	endpoint, err := handler.DataStore.Endpoint().Endpoint(portainer.EndpointID(params.EnvironmentID))
	if err != nil {
		return httperror.InternalServerError("could not fetch endpoint from db", err)
	}

	// create target docker API client
	targetClient, err := handler.dockerClientFactory.CreateClient(endpoint, "", nil)
	if err != nil {
		return httperror.InternalServerError("unable to create docker client", err)
	}

	ping, err := targetClient.Ping(r.Context())
	if err != nil {
		return httperror.InternalServerError("Could not ping docker env", err)
	}
	log.Info().Msg(ping.APIVersion)

	pullRes, err := targetClient.ImagePull(r.Context(), params.ImageID, types.ImagePullOptions{})
	if err != nil {
		return httperror.InternalServerError("could not pull coordinator image to registry", err)
	}
	defer pullRes.Close()
	print(pullRes)

	// run gramine image on target environment
	port1, err := nat.NewPort("tcp", "3306")
	createdBody, err := targetClient.ContainerCreate(r.Context(),
		&container.Config{
			Image: params.ImageID,
			ExposedPorts: nat.PortSet{
				port1: struct{}{},
			},
			Env: []string{
				"EDG_MARBLE_TYPE=" + params.Name + "_marble",
				"EDG_MARBLE_COORDINATOR_ADDR=coordinator:2001",
				"EDG_MARBLE_DNS_NAMES=localhost,app",
			},
			Domainname: "coordinator",
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				port1: []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "3306",
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
		params.Name)

	if err != nil {
		return httperror.InternalServerError("unable to create container", err)
	}

	// connect container to coordinator network
	err = targetClient.NetworkConnect(r.Context(), "coordinator", createdBody.ID, &network.EndpointSettings{})
	if err != nil {
		return httperror.InternalServerError("could not connect container to coordinator network", err)
	}

	// start container
	err = targetClient.ContainerStart(r.Context(), createdBody.ID, types.ContainerStartOptions{})
	if err != nil {
		return httperror.InternalServerError("Could not start container", err)
	}

	// remove coordinator from bridge network to fix SSL_ERROR_SYSCALL error
	err = targetClient.NetworkDisconnect(r.Context(), "bridge", createdBody.ID, false)
	if err != nil {
		return httperror.InternalServerError("could not remove container from bridge network", err)
	}
	return response.JSON(w, http.StatusOK)
}
