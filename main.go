package main

import (
	"cc-server/calculoid"
	"fmt"
	prometheusmiddleware "github.com/albertogviana/prometheus-middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func homeLinkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome in C&C server API\n")
}

func proxmoxProvisioningServerGetContainerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getting containers from proxmox..\n")
	fmt.Fprintf(w, "selected proxmox..\n")

	tracer := opentracing.GlobalTracer()
	_, rs := proxmoxProvisioningServerClient(tracer, "getall")
	fmt.Fprintf(w, "returned: %s\n", rs)

}

func proxmoxProvisioningServerContainerCreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create container on proxmox..\n")

	tracer := opentracing.GlobalTracer()
	_, rs := proxmoxProvisioningServerClient(tracer, "create")
	fmt.Fprintf(w, "returned: %s\n", rs)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func proxmoxProvisioningServerClient(tracer opentracing.Tracer, action string) (bool, string) {
	var (
		client            = "proxmoxProvisioningServerClient"
		requestBodyCreate = "{ \"disk\": 20}"
		req               *http.Request
		err               error
	)

	// nethttp.Transport from go-stdlib will do the tracing
	c := &http.Client{Transport: &nethttp.Transport{}}

	// create a top-level span to represent full work of the client
	span := tracer.StartSpan(client)
	span.SetTag(string(ext.Component), client)
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	if action == "getall" {
		req, err = http.NewRequest(
			"GET",
			fmt.Sprintf("http://192.168.121.10:%s", "4567"),
			nil,
		)
		if err != nil {
			onError(span, err)
			return false, ""
		}
	} else if action == "create" {

		req, err = http.NewRequest(
			"POST",
			fmt.Sprintf("http://192.168.121.10:%s%s", "4567", "/api/containers/create"),
			strings.NewReader(requestBodyCreate),
		)
		if err != nil {
			onError(span, err)
			return false, ""
		}
	} else {
		// no action selected
		return false, ""
	}

	req = req.WithContext(ctx)
	// wrap the request in nethttp.TraceRequest
	req, ht := nethttp.TraceRequest(tracer, req)
	defer ht.Finish()

	res, err := c.Do(req)
	if err != nil {
		onError(span, err)
		return false, ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		onError(span, err)
		return false, ""
	}
	fmt.Printf("Received result: %s\n", string(body))
	return true, string(body)
}

func onError(span opentracing.Span, err error) {
	// handle errors by recording them in the span
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
	log.Print(err)
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
	calculoidHandler := &calculoid.Handler{}

	router.Path("/metrics").Handler(promhttp.Handler())

	router.HandleFunc("/", homeLinkHandler)

	router.HandleFunc("/proxmox-provisioning-server/container/all", proxmoxProvisioningServerGetContainerHandler)

	router.HandleFunc("/proxmox-provisioning-server/container/create", proxmoxProvisioningServerContainerCreateHandler)

	router.Handle("/calculoid/webhook", calculoidHandler.CalculoidWebhook())

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
