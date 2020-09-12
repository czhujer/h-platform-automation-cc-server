package proxmox

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"log"
	"net/http"
)

func (proxmox *Proxmox) ProvisioningServerGetContainerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" && r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rContentType := r.Header.Get("Content-type")
	log.Printf("request Content-type: %s", rContentType)

	if rContentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	//TODO
	// add input params
	//     proxmoxURL

	proxmoxServer, proxmoxPort, _ = setProxmoxUrl(r)

	tracer := opentracing.GlobalTracer()

	log.Printf("getting all containers from (server: %s, port: %s)", proxmoxServer, proxmoxPort)

	status, rs := proxmox.proxmoxProvisioningServerClient(tracer, "getall", proxmoxServer, proxmoxPort)
	if status != true {
		w.WriteHeader(http.StatusInternalServerError)

		log.Printf("get containers failed: %s", rs)
		if isJSON(rs) {
			fmt.Fprintf(w, "%s\n", rs)
		} else {
			fmt.Fprintf(w, "{\"returned_body\": \"%s\"}\n", rs)
		}

		return
	}

	if isJSON(rs) {
		fmt.Fprintf(w, "%s\n", rs)
	} else {
		fmt.Fprintf(w, "{\"returned_body\": \"%s\", \"status:\": \"%s\"}\n", rs, status)
	}
}

func (proxmox *Proxmox) ProvisioningServerContainerCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rContentType := r.Header.Get("Content-type")
	log.Printf("request Content-type: %s", rContentType)

	if rContentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	//TODO
	// add input params
	//     disk, ram - one of this is required

	proxmoxServer, proxmoxPort, _ = setProxmoxUrl(r)

	tracer := opentracing.GlobalTracer()

	log.Printf("create container on proxmox.. (server: %s, port: %s, disk: , ram: )", proxmoxServer, proxmoxPort)

	status, rs := proxmox.proxmoxProvisioningServerClient(tracer, "create", proxmoxServer, proxmoxPort)
	if status != true {
		w.WriteHeader(http.StatusInternalServerError)

		log.Printf("create container failed: %s", rs)
		if isJSON(rs) {
			fmt.Fprintf(w, "%s\n", rs)
		} else {
			fmt.Fprintf(w, "{\"returned\": \"%s\"}\n", rs)
		}

		return
	}

	if isJSON(rs) {
		fmt.Fprintf(w, "%s\n", rs)
	} else {
		fmt.Fprintf(w, "{\"returned\": \"%s\"}\n", rs)
	}

	//TODO
	// add monitoring targets
	// prometheusRemote.AddTarget()
}
