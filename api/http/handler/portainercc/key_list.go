package portainercc

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

type ExportListKey struct {
	Id                 int
	Description        string
	TeamAccessPolicies portainer.TeamAccessPolicies
}

func (handler *Handler) listKeysByType(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	keyType, err := request.RetrieveQueryParameter(r, "type", false)
	if err != nil {
		return httperror.BadRequest("invalid query parameter", err)
	}

	keys, err := handler.DataStore.Key().Keys()

	if err != nil {
		return httperror.InternalServerError("couldn retrive keys from db", err)
	}

	// TODO filter for teamaccess/admin

	// filter by type
	result := make([]ExportListKey, 0)

	for _, key := range keys {
		if key.KeyType == keyType {
			entry := ExportListKey{
				Id:                 int(key.ID),
				Description:        key.Description,
				TeamAccessPolicies: key.TeamAccessPolicies,
			}
			result = append(result, entry)
		}
	}

	return response.JSON(w, result)
}
