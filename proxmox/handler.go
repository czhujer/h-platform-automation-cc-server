package proxmox

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

func (proxmox *Proxmox) ProvisioningServerGetContainerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getting containers from proxmox..\n")
	fmt.Fprintf(w, "selected proxmox..\n")

	tracer := opentracing.GlobalTracer()
	_, rs := proxmox.proxmoxProvisioningServerClient(tracer, "getall")
	fmt.Fprintf(w, "returned: %s\n", rs)

}

func (proxmox *Proxmox) PovisioningServerContainerCreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create container on proxmox..\n")

	tracer := opentracing.GlobalTracer()
	_, rs := proxmox.proxmoxProvisioningServerClient(tracer, "create")
	fmt.Fprintf(w, "returned: %s\n", rs)
}
