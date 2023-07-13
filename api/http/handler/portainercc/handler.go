package portainercc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	httperror "github.com/portainer/libhttp/error"

	"github.com/gorilla/mux"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	"github.com/portainer/portainer/api/docker"
	"github.com/portainer/portainer/api/http/security"
)

// Handler is the HTTP handler used to handle MOTD operations.
type Handler struct {
	*mux.Router
	DataStore           dataservices.DataStore
	DockerClientFactory *docker.ClientFactory
}

// NewHandler returns a new Handler
func NewHandler(bouncer *security.RequestBouncer, dataStore dataservices.DataStore, dockerClientFactory *docker.ClientFactory) *Handler {
	h := &Handler{
		Router:              mux.NewRouter(),
		DataStore:           dataStore,
		DockerClientFactory: dockerClientFactory,
	}

	//fill with default templates
	jsonTemplates, err := ioutil.ReadFile("/confidential-templates.json")

	if err != nil {
		fmt.Println(err)
	}

	var templates []ConfTempParams

	err = json.Unmarshal(jsonTemplates, &templates)
	if err != nil {
		fmt.Println(err)
	}

	for _, t := range templates {
		templateObject := &portainer.ConfidentialTemplate{
			ImageName:    t.ImageName,
			LogoURL:      t.LogoURL,
			TemplateName: t.TemplateName,
			Inputs:       t.Inputs,
			Secrets:      t.Secrets,
			ManifestBoilerplate: struct {
				ManifestParameters portainer.Parameters        "json:\"ManifestParameters\""
				ManifestSecrets    map[string]portainer.Secret "json:\"ManifestSecrets\""
			}(t.ManifestBoilerplate),
		}

		err = h.DataStore.ConfidentialTemplate().Create(templateObject)

		if err != nil {
			fmt.Println(err)
		}
	}

	restrictedRouter := h.NewRoute().Subrouter()
	// restrictedRouter.Use(bouncer.RestrictedAccess)

	//keys
	restrictedRouter.Handle("/portainercc/keys", httperror.LoggerHandler(h.listKeysByType)).Methods(http.MethodGet)
	restrictedRouter.Handle("/portainercc/keys", httperror.LoggerHandler(h.createKey)).Methods(http.MethodPost)
	restrictedRouter.Handle("/portainercc/keys/{id}", httperror.LoggerHandler(h.exportKey)).Methods(http.MethodGet)
	restrictedRouter.Handle("/portainercc/keys/{id}", httperror.LoggerHandler(h.updateKey)).Methods(http.MethodPost)
	restrictedRouter.Handle("/portainercc/keys/{id}", httperror.LoggerHandler(h.deleteKey)).Methods(http.MethodDelete)

	//secure images
	restrictedRouter.Handle("/portainercc/secimages", httperror.LoggerHandler(h.listSecImages)).Methods(http.MethodGet)

	//confidential templates
	restrictedRouter.Handle("/portainercc/confidential-templates", httperror.LoggerHandler(h.listConfidentialTemplates)).Methods(http.MethodGet)
	restrictedRouter.Handle("/portainercc/confidential-templates", httperror.LoggerHandler(h.deployConfidentialTemplate)).Methods(http.MethodPost)
	restrictedRouter.Handle("/portainercc/confidential-templates/add", httperror.LoggerHandler(h.createConfidentialTemplate)).Methods(http.MethodPost)

	//build and run img
	restrictedRouter.Handle("/portainercc/confidential-templates/python", httperror.LoggerHandler(h.runConfidentialPython)).Methods(http.MethodPost)
	restrictedRouter.Handle("/portainercc/confidential-templates/node", httperror.LoggerHandler(h.runConfidentialNode)).Methods(http.MethodPost)

	return h
}
