package ra

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	"github.com/portainer/portainer/api/demo"
	"github.com/portainer/portainer/api/docker"
)

type requestBouncer interface {
	AuthenticatedAccess(h http.Handler) http.Handler
	AdminAccess(h http.Handler) http.Handler
	RestrictedAccess(h http.Handler) http.Handler
	PublicAccess(h http.Handler) http.Handler
	AuthorizedEndpointOperation(r *http.Request, endpoint *portainer.Endpoint) error
	AuthorizedEdgeEndpointOperation(r *http.Request, endpoint *portainer.Endpoint) error
}

type Handler struct {
	*mux.Router
	requestBouncer      requestBouncer
	demoService         *demo.Service
	DataStore           dataservices.DataStore
	dockerClientFactory *docker.ClientFactory
}

func NewHandler(bouncer requestBouncer, dockerClientFactory *docker.ClientFactory) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
	}

	h.Handle("/ra/coordinator/build",
		bouncer.PublicAccess(httperror.LoggerHandler(h.raCoordinatorBuild))).Methods(http.MethodPost)
	h.Handle("/ra/coordinator/list",
		bouncer.PublicAccess(httperror.LoggerHandler(h.raCoordinatorList))).Methods(http.MethodGet)
	return h
}
