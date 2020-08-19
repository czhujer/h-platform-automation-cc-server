package main

import (
	"cc-server/calculoid"
	"fmt"
	prometheusmiddleware "github.com/albertogviana/prometheus-middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

func homeLinkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "")
}

func main() {

	var opts prometheusmiddleware.Opts

	middleware := prometheusmiddleware.NewPrometheusMiddleware(opts)

	calculoidHandler := &calculoid.Handler{}

	filename := os.Args[0]
	log.Printf("starting %s \n", filename)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.InstrumentHandlerDuration)

	router.Path("/metrics").Handler(promhttp.Handler())

	router.HandleFunc("/", homeLinkHandler)

	router.Handle("/calculoid/webhook", calculoidHandler.CalculoidWebhook())

	notFoundHandlerFunc := http.HandlerFunc(notFoundHandler)
	router.NotFoundHandler = http.Handler(middleware.InstrumentHandlerDuration(notFoundHandlerFunc))

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
