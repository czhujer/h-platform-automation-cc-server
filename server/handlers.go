package server

import (
	"fmt"
	"net/http"
)

func homeLinkHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
	// add html template
	fmt.Fprintf(w, "Welcome in C&C server API\n")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
