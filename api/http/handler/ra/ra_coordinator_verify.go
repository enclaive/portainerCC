package ra

import (
	"context"
	"net"
	"net/http"
	"time"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	"github.com/portainer/portainer/api/http/client"
	"github.com/rs/zerolog/log"
)

func (handler *Handler) raCoordinatorVerify(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Info().Msg("hello coordinator")
	client := client.NewHTTPClient()
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		// DualStack: true, // this is deprecated as of go 1.16
	}
	// or create your own transport, there's an example on godoc.
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if addr == "coordinator:443" {
			addr = "172.17.0.5:443"
		}
		return dialer.DialContext(ctx, network, addr)
	}
	resp, err := client.Get("https://coordinator")
	if err != nil {
		log.Err(err)
	}
	// fmt.Println(resp.Header, err)
	return response.JSON(w, resp.Header)
}
