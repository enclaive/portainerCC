package portainercc

import (
	"encoding/json"
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	portainer "github.com/portainer/portainer/api"
)

type KeyParams struct {
	KeyType            string
	Description        string
	TeamAccessPolicies portainer.TeamAccessPolicies
	Data               string
}

type UpdateKeyParams struct {
	TeamAccessPolicies portainer.TeamAccessPolicies
}

type ExportKey struct {
	Id  int
	PEM string
}

func (handler *Handler) generateOrImport(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params KeyParams
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return httperror.BadRequest("request body malefomred", err)
	}

	return httperror.BadRequest("invalid body content", nil)
}
