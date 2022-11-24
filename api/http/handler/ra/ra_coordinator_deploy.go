package ra

import (
	"encoding/base64"
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

type DeployCoordinatorParams struct {
	CoordinatorID int `json:"coordinatorId"`
	EnvironmentID int `json:"environmentId"`
}

// FIXME image repo still hardcoded
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

	// get target environment
	targetEndpoint, err := handler.DataStore.Endpoint().Endpoint(portainer.EndpointID(params.EnvironmentID))
	if err != nil {
		return httperror.InternalServerError("unable to find requested endpoint", err)
	}

	// create target docker API client
	targetClient, err := handler.dockerClientFactory.CreateClient(targetEndpoint, "", nil)
	if err != nil {
		return httperror.InternalServerError("unable to create docker client", err)
	}

	ping, err := targetClient.Ping(r.Context())
	if err != nil {
		return httperror.InternalServerError("Could not ping docker env", err)
	}
	log.Info().Msg(ping.APIVersion)

	// pull coordinator image
	authConfig := types.AuthConfig{
		Username:      "sgxdcaprastuff",
		Password:      "dckr_pat_ASB6_d6hVfhgHsNXByxEWYjfXtc",
		ServerAddress: "https://index.docker.io/v2/",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
	pullOptions := types.ImagePullOptions{RegistryAuth: authConfigEncoded}

	pullRes, err := targetClient.ImagePull(r.Context(), "sgxdcaprastuff/coordinatortest", pullOptions)
	if err != nil {
		return httperror.InternalServerError("could not pull coordinator image to registry", err)
	}
	defer pullRes.Close()
	print(pullRes)

	port1, err := nat.NewPort("tcp", "9944")
	port2, err := nat.NewPort("tcp", "4433")

	// run coordinator image on target environment
	createdBody, err := targetClient.ContainerCreate(r.Context(),
		&container.Config{
			Image: "sgxdcaprastuff/coordinatortest",
			ExposedPorts: nat.PortSet{
				port1: struct{}{},
				port2: struct{}{},
			},
			Env: []string{
				"OE_SIMULATION=0",
				"OE_LOG_LEVEL=INFO",
				"EDG_COORDINATOR_MESH_ADDR=" + coordinator.Name + ":2001",
				"EDG_COORDINATOR_CLIENT_ADDR=" + coordinator.Name + ":4433",
				"EDG_COORDINATOR_DNS_NAMES=" + coordinator.Name,
				"EDG_COORDINATOR_PROMETHEUS_ADDR=0.0.0.0:9944",
			},
			Domainname: "coordinator",
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

	var networkID string

	// check if coordinator network exists
	networkResource, err := targetClient.NetworkInspect(r.Context(), "coordinator", types.NetworkInspectOptions{})
	if err != nil {
		// create coordinator network
		createNetworkResponse, err := targetClient.NetworkCreate(r.Context(), "coordinator", types.NetworkCreate{
			IPAM: &network.IPAM{
				Config: []network.IPAMConfig{
					{
						Subnet: "172.20.0.0/16",
					},
				},
			},
		})
		if err != nil {
			return httperror.InternalServerError("could not create coordinator network", err)
		}
		networkID = createNetworkResponse.ID
	} else {
		networkID = networkResource.ID
	}

	// connect container to coordinator network
	err = targetClient.NetworkConnect(r.Context(), networkID, createdBody.ID, &network.EndpointSettings{
		IPAddress:   "172.20.0.20",
		Gateway:     "172.20.0.1",
		IPPrefixLen: 16,
		IPAMConfig: &network.EndpointIPAMConfig{
			IPv4Address: "172.20.0.20",
		},
	})
	if err != nil {
		return httperror.InternalServerError("could not connect container to coordinator network", err)
	}

	// start coordinator container
	_ = targetClient.ContainerStart(r.Context(), createdBody.ID, types.ContainerStartOptions{})

	// remove coordinator from bridge network to fix SSL_ERROR_SYSCALL error
	err = targetClient.NetworkDisconnect(r.Context(), "bridge", createdBody.ID, false)
	if err != nil {
		return httperror.InternalServerError("could not remove coordinator container from bridge network", err)
	}

	targetClient.Close()

	// new coordinatorDeployment Object
	coordinatorDeployment := &portainer.CoordinatorDeployment{
		CoordinatorID: params.CoordinatorID,
		EndpointID:    params.EnvironmentID,
		Verified:      false,
	}

	// check if a coordinatorDeployment already exists for endpoint and if so, delete it
	deployments, err := handler.DataStore.CoordinatorDeployment().CoordinatorDeployments()
	for _, deployment := range deployments {
		if deployment.EndpointID == params.EnvironmentID {
			err := handler.DataStore.CoordinatorDeployment().Delete(deployment.ID)
			if err != nil {
				return httperror.InternalServerError("could not delete old deployment", err)
			}
		}
	}

	// create new coordinatorDeployment in DB
	err = handler.DataStore.CoordinatorDeployment().Create(coordinatorDeployment)
	if err != nil {
		return httperror.InternalServerError("could not create coordinatorDeployment in db", err)
	}

	return response.JSON(w, portainer.StatusOk)
}
