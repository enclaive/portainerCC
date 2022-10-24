package portainercc

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"

	"github.com/gorilla/mux"
	"github.com/portainer/portainer/api/dataservices"
	"github.com/portainer/portainer/api/http/security"
)

// Handler is the HTTP handler used to handle MOTD operations.
type Handler struct {
	*mux.Router
	DataStore dataservices.DataStore
}

// NewHandler returns a new Handler
func NewHandler(bouncer *security.RequestBouncer, dataStore dataservices.DataStore) *Handler {
	h := &Handler{
		Router:    mux.NewRouter(),
		DataStore: dataStore,
	}

	restrictedRouter := h.NewRoute().Subrouter()
	restrictedRouter.Use(bouncer.RestrictedAccess)

	restrictedRouter.Handle("/portainercc/keys", httperror.LoggerHandler(h.listByType)).Methods(http.MethodGet)
	restrictedRouter.Handle("/portainercc/keys", httperror.LoggerHandler(h.create)).Methods(http.MethodPost)
	restrictedRouter.Handle("/portainercc/keys/{id}", httperror.LoggerHandler(h.export)).Methods(http.MethodGet)

	return h
}
