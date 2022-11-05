package ra

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

func (handler *Handler) raCoordinatorGet(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	coordinatorID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid coordinator identifier route variable", err)
	}

	coordinator, err := handler.DataStore.Coordinator().Coordinator(portainer.CoordinatorID(coordinatorID))
	if err != nil {
		return httperror.InternalServerError("Could not retrieve coordinator from database", err)
	}

	return response.JSON(w, coordinator)
}
