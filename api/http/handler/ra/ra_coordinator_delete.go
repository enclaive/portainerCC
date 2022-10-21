package ra

import (
	"net/http"

	"github.com/docker/docker/api/types"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/internal/endpointutils"
	"github.com/rs/zerolog/log"
)

func (handler *Handler) raCoordinatorDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	coordinatorID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid coordinator identifier route variable", err)
	}

	// get coordinator from DB
	coordinator, err := handler.DataStore.Coordinator().Coordinator(portainer.CoordinatorID(coordinatorID))
	if err != nil {
		return httperror.InternalServerError("Could not retrieve coordinator from database", err)
	}

	// delete coordinator image from local docker env
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

	// remove coordinator docker image
	_, err = client.ImageRemove(r.Context(), coordinator.ImageID, types.ImageRemoveOptions{})
	if err != nil {
		return httperror.InternalServerError("unable to delete coordinator image", err)
	}

	// delete coordinator from db
	err = handler.DataStore.Coordinator().Delete(portainer.CoordinatorID(coordinatorID))
	if err != nil {
		return httperror.InternalServerError("unable to delete coordinator from database", err)
	}

	return response.JSON(w, http.StatusOK)
}
