package proxmox

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

func (proxmox *Proxmox) ProvisioningServerGetContainerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// add json headers

	// TODO
	// add get/post check

	fmt.Fprintf(w, "getting containers from proxmox..\n")
	fmt.Fprintf(w, "selected proxmox..\n")

	tracer := opentracing.GlobalTracer()

	// TODO
	// add err
	_, rs := proxmox.proxmoxProvisioningServerClient(tracer, "getall")
	fmt.Fprintf(w, "returned: %s\n", rs)

}

func (proxmox *Proxmox) PovisioningServerContainerCreateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// add json headers

	// TODO
	// add get/post check

	fmt.Fprintf(w, "create container on proxmox..\n")

	tracer := opentracing.GlobalTracer()

	// TODO
	// check err
	_, rs := proxmox.proxmoxProvisioningServerClient(tracer, "create")
	fmt.Fprintf(w, "returned: %s\n", rs)
}
