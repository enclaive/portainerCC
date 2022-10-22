package portainercc

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/mux"

	"github.com/portainer/agent/http/proxy"
	"github.com/portainer/agent/http/security"
)

// Handler represents an HTTP API Handler for host specific actions
type Handler struct {
	*mux.Router
}

// NewHandler returns a new instance of Handler
func NewHandler(agentProxy *proxy.AgentProxy, notaryService *security.NotaryService) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	h.HandleFunc("/{useless}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HALLO")

		if r.Method == "CONNECT" {
			fmt.Println("CONNEEEEEEEEEEEEEECT")
		}

		var proxy = createProxy()
		proxy.ServeHTTP(w, r)
	})

	// h.coordinator)

	return h
}

func createProxy() *httputil.ReverseProxy {
	target, _ := url.Parse("https://172.17.0.4:4444/status")
	director := func(req *http.Request) {
		req.URL = target
	}

	return &httputil.ReverseProxy{Director: director, Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}}
}
