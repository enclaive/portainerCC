package portainercc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
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
	//get endpoint
	endpoint, err := handler.DataStore.Endpoint().Endpoint(portainer.EndpointID(params.EnvId))
	if err != nil {
		return httperror.InternalServerError("unable to find requested endpoint", err)
	}

	//TODO its hardcoded for docker..?!?..
	// create docker API client
	client, err := handler.DockerClientFactory.CreateClient(endpoint, "", nil)
	if err != nil {
		return httperror.InternalServerError("could not create docker client", err)
	}

	res, err := client.ImagePull(r.Context(), template.ImageName, types.ImagePullOptions{})
	if err != nil {
		return httperror.InternalServerError("Unable to pull image", err)
	}
	defer res.Close()

	//if we dont read the res, the image would not be tagged ..
	buf := new(strings.Builder)
	_, _ = io.Copy(buf, res)
	fmt.Println(buf.String())

	//read labels
	//read pcc.mrenclave, pcc.mrsigner

	_, inspectRaw, err := client.ImageInspectWithRaw(r.Context(), template.ImageName)
	if err != nil {
		return httperror.InternalServerError("Unable to inspect image", err)
	}

	var JSON map[string]interface{}
	json.Unmarshal(inspectRaw, &JSON)

	cfg := JSON["Config"].(map[string]interface{})
	labels := cfg["Labels"].(map[string]interface{})

	mrenclave := labels["pcc.mrenclave"].(string)
	mrsigner := labels["pcc.mrsigner"].(string)

	//add to marblemanifest
	fmt.Printf("I WILL PUT THIS INTO MARBLERUN MANIFEST / Container to deploy:\n")
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("Dockerimage to use: %s:\n", template.ImageName)
	fmt.Printf("MRENCLAVE: %s (extracted from Image)\n", mrenclave)
	fmt.Printf("MRSIGNER: %s (extracted from Image)\n", mrsigner)
	fmt.Printf("Packagename %s\n", params.Name)
	fmt.Printf("Secrets:\n")
	for _, val := range template.Values {
		fmt.Printf("\t%s: %s\n", val, params.Values[val])
	}

	return response.JSON(w, params)
}
