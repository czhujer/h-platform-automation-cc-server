package main

import (
	"cc-server/calculoid"
	"fmt"
	prometheusmiddleware "github.com/albertogviana/prometheus-middleware"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"time"
	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
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

	// http listen port
	const httpAddress = ":8080"
	const jaegerHostPort = ":6831"

	filename := os.Args[0]
	var opts prometheusmiddleware.Opts

	log.Printf("starting %s \n", filename)

	middleware := prometheusmiddleware.NewPrometheusMiddleware(opts)

	calculoidHandler := &calculoid.Handler{}

	tracingcfg := jaegercfg.Configuration{
		ServiceName: "c-and-c-server",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  jaegerHostPort,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := prometheus.New()

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := tracingcfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Fatal("cannot initialize Jaeger Tracer", err)
	}

	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.InstrumentHandlerDuration)

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
	//loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, router)

	// start server
	err = http.ListenAndServe(
		httpAddress,
		router)
	if err != nil {
		panic(err)
	}
}
