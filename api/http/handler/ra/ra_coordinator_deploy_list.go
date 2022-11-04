package ra

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

func (handler *Handler) raCoordinatorDeployList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	coordinators, err := handler.DataStore.CoordinatorDeployment().CoordinatorDeployments()
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve coordinator deployments", err)
	}
	return response.JSON(w, coordinators)
}
