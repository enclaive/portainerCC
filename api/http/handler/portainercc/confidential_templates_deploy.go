package portainercc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

type ConfTempDeployParams struct {
	Id     portainer.ConfidentialTemplateId
	EnvId  int
	Name   string
	Values map[string]string
}

func (handler *Handler) deployConfidentialTemplate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params ConfTempDeployParams
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return httperror.BadRequest("request body malefomred", err)
	}

	fmt.Println(params)

	//check if all values are set
	template, err := handler.DataStore.ConfidentialTemplate().ConfidentialTemplate(portainer.ConfidentialTemplateId(params.Id))
	if err != nil {
		return httperror.BadRequest("invalid template id", err)
	}

	for _, val := range template.Values {
		if _, ok := params.Values[val]; !ok {
			return httperror.BadRequest("request body malefomred", fmt.Errorf("values missing. Expected: %s ", strings.Join(template.Values[:], ",")))
		}
	}

	//pull image and get mr enclave mr signer

	return response.JSON(w, params)
}
