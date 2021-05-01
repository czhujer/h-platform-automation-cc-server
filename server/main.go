package server

import (
	prometheusmiddleware "github.com/albertogviana/prometheus-middleware"
	"log"
	"net/http"
	"os"
)

const httpAddress = ":8080"

func RunServer() {

	var filename = os.Args[0]
	var opts prometheusmiddleware.Opts

	log.Printf("starting %s \n", filename)

	middleware := prometheusmiddleware.NewPrometheusMiddleware(opts)

	tracer, _ := CreateTracer()

	router := CreateRouter(middleware, tracer)

	// start server
	err := http.ListenAndServe(
		httpAddress,
		router)
	if err != nil {
		panic(err)
	}
}
