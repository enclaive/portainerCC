package ra

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

func (handler *Handler) raCoordinatorList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	coordinators, err := handler.DataStore.Coordinator().Coordinators()
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve coordinators", err)
	}
	return response.JSON(w, coordinators)
}
