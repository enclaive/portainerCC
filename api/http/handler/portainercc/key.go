package portainercc

import (
	"encoding/json"
	"log"
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
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

	keyObject := &portainer.Key{
		KeyType:            params.KeyType,
		Description:        params.Description,
		TeamAccessPolicies: params.TeamAccessPolicies,
	}

	handler.DataStore.Key().Create(keyObject, params.Data)

	log.Print("AHA?")
	log.Print(keyObject)
	log.Printf("BODY: %s", params.Data)

	log.Print(params.Data != "")

	return httperror.BadRequest("invalid body content", nil)
}

func (handler *Handler) getKeys(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	keys, err := handler.DataStore.Key().Keys()

	if err != nil {
		return httperror.InternalServerError("couldn retrive keys from db", err)
	}

	return response.JSON(w, keys)
}
