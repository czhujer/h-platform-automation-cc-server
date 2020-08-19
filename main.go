package main

import (
	"cc-server/calculoid"
	"fmt"
	prometheusmiddleware "github.com/albertogviana/prometheus-middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
	"net/http"
	"os"
	//jaegerlog "github.com/uber/jaeger-client-go/log"
)

func homeLinkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "")
}

func main() {
	filename := os.Args[0]
	var opts prometheusmiddleware.Opts

	log.Printf("starting %s \n", filename)

	middleware := prometheusmiddleware.NewPrometheusMiddleware(opts)

	calculoidHandler := &calculoid.Handler{}

	jaegercfg := jaegercfg.Configuration{
		ServiceName: "c-and-c-server",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	//jLogger := jaegerlog.StdLogger
	//jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := jaegercfg.NewTracer(
	//jaegercfg.Logger(jLogger),
	//jaegercfg.Metrics(jMetricsFactory),
	)
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	//tracer := opentracing.Tracer()
	//if err != nil {
	//	log.Fatal("cannot initialize Jaeger Tracer", err)
	//}

	//okHandler := func(w http.ResponseWriter, r *http.Request) {
	//	// do something
	//	data := "Hello"
	//	fmt.Fprintf(w, data)
	//}

	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.InstrumentHandlerDuration)

	//tracingmiddleware := gorilla.Middleware(
	//	tracer,
	//	http.HandlerFunc(homeLinkHandler),
	//)

	router.Path("/metrics").Handler(promhttp.Handler())

	router.HandleFunc("/", homeLinkHandler)

	router.Handle("/calculoid/webhook", calculoidHandler.CalculoidWebhook())

	notFoundHandlerFunc := http.HandlerFunc(notFoundHandler)
	router.NotFoundHandler = http.Handler(middleware.InstrumentHandlerDuration(notFoundHandlerFunc))

	// Add tracing to all routes
	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		route.Handler(
			gorilla.Middleware(tracer, route.GetHandler()))
		return nil
	})

	// Add apache-like logging to all routes
	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, router)

	// start server
	err = http.ListenAndServe(":8080", loggedRouter)
	if err != nil {
		panic(err)
	}
}
