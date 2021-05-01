package server

import (
	prometheusmiddleware "github.com/albertogviana/prometheus-middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics/prometheus"
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

	router := CreateRouter(middleware, tracer)

	// start server
	err = http.ListenAndServe(
		httpAddress,
		router)
	if err != nil {
		panic(err)
	}
}
