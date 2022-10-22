package portainercc

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	httperror "github.com/portainer/libhttp/error"
)

func (handler *Handler) coordinator(rw http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	// Loop over header names
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}

	fmt.Println("Hallo ich bin hier")

	//get the url from header
	// coordinatorURL := r.Header.Get("X-Coordinator-URL")
	coordinatorURL := "172.17.0.4:4444"

	//create a proxy
	coordProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   coordinatorURL,
	})
	coordProxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//alter request
	r.Host = coordinatorURL

	coordProxy.ServeHTTP(rw, r)
	return nil
}

// func (handler *Handler) fillHostInfo(hi *agent.HostInfo) error {
// 	devices, devicesError := handler.systemService.GetPciDevices()
// 	if devicesError != nil {
// 		return devicesError
// 	}
// 	hi.PCIDevices = devices

// 	disks, disksError := handler.systemService.GetDiskInfo()
// 	if disksError != nil {
// 		return disksError
// 	}
// 	hi.PhysicalDisks = disks
// 	return nil
// }
