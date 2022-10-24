package portainercc

import (
	"crypto/x509"
	"encoding/pem"
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

type ExportKey struct {
	Id          int
	KeyType     string
	Description string
	Export      string
}

func (handler *Handler) export(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	id, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("invalid query parameter", err)
	}

	key, err := handler.DataStore.Key().Key(portainer.KeyID(id))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find a key with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("error retrieving key from database", err)
	}

	result := ExportKey{
		Id:          int(key.ID),
		KeyType:     key.KeyType,
		Description: key.Description,
	}

	if key.KeyType == "SIGNING" {
		//encode rsa key to PEM
		pem := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(key.SigningKey),
			},
		)

		result.Export = string(pem[:])
	} else {
		//file encryptionkey
		// TODO hex
		result.Export = key.PFKey
	}

	return response.JSON(w, result)
}
