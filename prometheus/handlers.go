package prometheus

import (
	prometheusRemote "cc-server/prometheus/remote"
	"fmt"
	"net/http"
)

func (prometheus *Prometheus) RemoteTargetAddHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//TODO
	// check if request is GET/POST

	//TODO
	// add loading/generating vmNameFull variable

	err := prometheusRemote.AddTarget()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"result\": \"%s\"}\n", err)
	} else {
		fmt.Fprintf(w, "{\"result\": \"prometheus target added\"}\n")
	}
}

func (prometheus *Prometheus) RemoteTargetRemoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//TODO
	// check if request is POST

	//TODO
	// add loading vmNameFull variable from request

	err := prometheusRemote.RemoveTarget()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"result\": \"%s\"}\n", err)
	} else {
		fmt.Fprintf(w, "{\"result\": \"prometheus target removed\"}\n")
	}
}
