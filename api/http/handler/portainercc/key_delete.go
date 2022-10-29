package portainercc

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

func (handler *Handler) deleteKey(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	// read query id
	keyID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("invalid query parameter", err)
	}

	_, err = handler.DataStore.Key().Key(portainer.KeyID(keyID))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find a key with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("Unable to find a key with the specified identifier inside the database", err)
	}

	err = handler.DataStore.Key().Delete(portainer.KeyID(keyID))
	if err != nil {
		return httperror.InternalServerError("Unable to delete the key from the database", err)
	}

	return response.Empty(w)
}
