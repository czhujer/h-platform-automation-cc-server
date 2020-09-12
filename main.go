package main

import (
	"cc-server/calculoid"
	prometheusRemote "cc-server/prometheus/remote"
	"cc-server/proxmox"
	tfOwncloudstack "cc-server/terraform/owncloudstack"
	tfOwncloudstackDocker "cc-server/terraform/owncloudstackdocker"
	"fmt"
	prometheusmiddleware "github.com/albertogviana/prometheus-middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"log"
	"net/http"
	"os"
)

func homeLinkHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
	// add html template
	fmt.Fprintf(w, "Welcome in C&C server API\n")
}

func prometheusRemoteTargetAddHandler(w http.ResponseWriter, r *http.Request) {
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

func prometheusRemoteTargetRemoveHandler(w http.ResponseWriter, r *http.Request) {
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

func terraformOwncloudstackCreateHandler(w http.ResponseWriter, r *http.Request) {
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
}

func terraformOwncloudstackdockerCreateHandler(w http.ResponseWriter, r *http.Request) {
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
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func main() {

	const httpAddress = ":8080"
	var filename = os.Args[0]
	var opts prometheusmiddleware.Opts

	log.Printf("starting %s \n", filename)

	middleware := prometheusmiddleware.NewPrometheusMiddleware(opts)

	// Initialize tracer with a logger and a metrics factory
	tracingcfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Printf("Could not parse Jaeger env vars: %s", err.Error())
		return
	}

	if tracingcfg.ServiceName == "" {
		tracingcfg.ServiceName = "c-and-c-server"
	}
	log.Printf("Jaeger ServiceName: %s", tracingcfg.ServiceName)
	log.Printf("Jaeger LocalAgentHostPort: %s", tracingcfg.Reporter.LocalAgentHostPort)

	tracingcfg.Sampler.Type = jaeger.SamplerTypeConst
	tracingcfg.Sampler.Param = 1
	tracingcfg.Reporter.LogSpans = true

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := prometheus.New()

	tracer, closer, err := tracingcfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Fatal("cannot initialize Jaeger Tracer: ", err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// Initialize router
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.InstrumentHandlerDuration)

	// handlers
	calculoid := &calculoid.Handler{}

	proxmox := &proxmox.Proxmox{}

	router.Path("/metrics").Handler(promhttp.Handler())

	router.HandleFunc("/", homeLinkHandler)

	// v1-arch handlers
	router.Handle("/calculoid/webhook", calculoid.CalculoidWebhookHandler())

	// v2-arch handlers
	router.HandleFunc("/proxmox-provisioning-server/container/all", proxmox.ProvisioningServerGetContainerHandler)

	router.HandleFunc("/proxmox-provisioning-server/container/create", proxmox.PovisioningServerContainerCreateHandler)

	// monitoring handlers
	router.HandleFunc("/prometheus/remote/target/add", prometheusRemoteTargetAddHandler)
	router.HandleFunc("/prometheus/remote/target/remove", prometheusRemoteTargetRemoveHandler)

	// terraform handlers
	router.HandleFunc("/terraform/owncloudstack/create", terraformOwncloudstackCreateHandler)

	// default handler
	notFoundHandlermw := gorilla.Middleware(
		tracer,
		http.HandlerFunc(notFoundHandler),
	)
	notFoundHandlerFunc := http.Handler(notFoundHandlermw)
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
	err = http.ListenAndServe(
		httpAddress,
		loggedRouter)
	if err != nil {
		panic(err)
	}
}
