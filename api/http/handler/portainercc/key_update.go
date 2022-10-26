package portainercc

import (
	"encoding/json"
	"net/http"
	"reflect"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

type UpdateKeyParams struct {
	TeamAccessPolicies portainer.TeamAccessPolicies
}

func (handler *Handler) updateKey(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	// read query id
	keyID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("invalid query parameter", err)
	}

	// create JSON object
	var params UpdateKeyParams
	err = json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return httperror.BadRequest("request body maleformed", err)
	}

	// get key object from database
	key, err := handler.DataStore.Key().Key(portainer.KeyID(keyID))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find a key with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("error retrieving key from database", err)
	}

	// update the key teams
	if params.TeamAccessPolicies != nil && !reflect.DeepEqual(params.TeamAccessPolicies, key.TeamAccessPolicies) {
		key.TeamAccessPolicies = params.TeamAccessPolicies
	}

	// update the key
	err = handler.DataStore.Key().Update(key.ID, key)
	if err != nil {
		return httperror.InternalServerError("error updating key in database", err)
	}

	result := KeyResponse{
		Id:                 key.ID,
		KeyType:            key.KeyType,
		Description:        key.Description,
		TeamAccessPolicies: key.TeamAccessPolicies,
	}

	return response.JSON(w, result)
}
