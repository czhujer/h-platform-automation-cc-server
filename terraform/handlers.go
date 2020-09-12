package terraform

import (
	tfOwncloudstack "cc-server/terraform/owncloudstack"
	tfOwncloudstackDocker "cc-server/terraform/owncloudstackdocker"
	"fmt"
	"net/http"
)

func (terraform *Terraform) TerraformOwncloudstackCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//TODO
	// check if request is GET/POST

	//TODO
	// add loading/generating vmNameFull variable

	err := tfOwncloudstack.Create()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"result\": \"%s\"}\n", err)
	} else {
		fmt.Fprintf(w, "{\"result\": \"terraform create executed\"}\n")
	}

	//TODO
	// add monitoring targets
	// prometheusRemote.AddTarget()
}

func (terraform *Terraform) TerraformOwncloudstackdockerCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//TODO
	// check if request is GET/POST

	//TODO
	// add loading/generating vmNameFull variable

	//TODO
	// add logic
	err := tfOwncloudstackDocker.Create()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"result\": \"%s\"}\n", err)
	} else {
		fmt.Fprintf(w, "{\"result\": \"terraform create executed\"}\n")
	}

	//TODO
	// add monitoring targets
	// prometheusRemote.AddTarget()
}
