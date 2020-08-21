package proxmox

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"log"
	"net/http"
)

func (proxmox *Proxmox) ProvisioningServerGetContainerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// add json headers

	// TODO
	// add input params
	//     proxmoxServer
	proxmoxServer = "192.168.121.10"
	tracer := opentracing.GlobalTracer()

	log.Printf("getting all containers from proxmox")

	if r.Method == "GET" {
		// TODO
		// add err
		_, rs := proxmox.proxmoxProvisioningServerClient(tracer, "getall", proxmoxServer)
		fmt.Fprintf(w, "returned: %s\n", rs)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (proxmox *Proxmox) PovisioningServerContainerCreateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// add json headers

	// TODO
	// add input params
	//     proxmoxServer
	proxmoxServer = "192.168.121.10"
	tracer := opentracing.GlobalTracer()

	log.Printf("create container on proxmox")

	if r.Method == "GET" {
		// TODO
		// check err
		_, rs := proxmox.proxmoxProvisioningServerClient(tracer, "create", proxmoxServer)
		fmt.Fprintf(w, "returned: %s\n", rs)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
