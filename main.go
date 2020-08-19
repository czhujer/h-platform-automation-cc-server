package main

import (
	"cc-server/calculoid"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	//promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
}

func main() {

	//var (
	//	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	//		Name: "myapp_http_duration_seconds",
	//		Help: "Duration of HTTP requests.",
	//	}, []string{"path"})
	//)

	// prometheusMiddleware implements mux.MiddlewareFunc.
	//func prometheusMiddleware(next http.Handler) http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	route := mux.CurrentRoute(r)
	//	path, _ := route.GetPathTemplate()
	//	timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
	//	next.ServeHTTP(w, r)
	//	timer.ObserveDuration()
	//})
	//}

	calculoidHandler := &calculoid.Handler{}

	filename := os.Args[0]
	log.Printf("starting %s \n", filename)

	router := mux.NewRouter().StrictSlash(true)
	//router.Use(prometheusMiddleware)

	router.Path("/metrics").Handler(promhttp.Handler())

	router.HandleFunc("/", homeLink)

	router.Handle("/calculoid/webhook", calculoidHandler.CalculoidWebhook())

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
