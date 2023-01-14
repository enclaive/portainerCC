package portainercc

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

func (handler *Handler) listConfidentialTemplates(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	templates, err := handler.DataStore.ConfidentialTemplate().ConfidentialTemplates()

	if err != nil {
		return httperror.InternalServerError("couldn retrive templates from db", err)
	}

	return response.JSON(w, templates)
}
