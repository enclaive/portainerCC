package portainercc

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

type ExportListImage struct {
	Id                 int
	Description        string
	TeamAccessPolicies portainer.TeamAccessPolicies
}

func (handler *Handler) listSecImages(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	images, err := handler.DataStore.SecureImage().SecureImages()

	if err != nil {
		return httperror.InternalServerError("couldn retrive keys from db", err)
	}

	return response.JSON(w, images)
}
