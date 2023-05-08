package portainercc

import (
	"encoding/json"
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

type ConfTempParams struct {
	ImageName           string
	LogoURL             string
	TemplateName        string
	Inputs              []portainer.Input
	Secrets             map[string]string
	ManifestBoilerplate struct {
		ManifestParameters portainer.Parameters
		ManifestSecrets    map[string]portainer.Secret
	}
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
		Inputs:       params.Inputs,
		Secrets:      params.Secrets,
		ManifestBoilerplate: struct {
			ManifestParameters portainer.Parameters        "json:\"ManifestParameters\""
			ManifestSecrets    map[string]portainer.Secret "json:\"ManifestSecrets\""
		}(params.ManifestBoilerplate),
	}

	err = handler.DataStore.ConfidentialTemplate().Create(templateObject)

	if err != nil {
		return httperror.InternalServerError("could not save template in db", err)
	}

	return response.JSON(w, templateObject)
}
