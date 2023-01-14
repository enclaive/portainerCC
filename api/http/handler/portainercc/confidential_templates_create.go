package portainercc

import (
	"encoding/json"
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

type ConfTempParams struct {
	ImageName    string
	LogoURL      string
	TemplateName string
	Values       []string
}

func (handler *Handler) createConfidentialTemplate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params ConfTempParams
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return httperror.BadRequest("request body malefomred", err)
	}

	templateObject := &portainer.ConfidentialTemplate{
		ImageName:    params.ImageName,
		LogoURL:      params.LogoURL,
		TemplateName: params.TemplateName,
		Values:       params.Values,
	}

	err = handler.DataStore.ConfidentialTemplate().Create(templateObject)

	if err != nil {
		return httperror.InternalServerError("could not save template in db", err)
	}

	return response.JSON(w, templateObject)
}
