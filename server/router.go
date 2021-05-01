package server

import (
	"cc-server/calculoid"
	ccPrometheus "cc-server/prometheus"
	"cc-server/proxmox"
	"cc-server/terraform"
	prometheusmiddleware "github.com/albertogviana/prometheus-middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

func CreateRouter(middleware *prometheusmiddleware.PrometheusMiddleware, tracer opentracing.Tracer) http.Handler {

	// Initialize router
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.InstrumentHandlerDuration)

	// handlers
	calc := &calculoid.Handler{}

	pxm := &proxmox.Proxmox{}

	prom := &ccPrometheus.Prometheus{}

	tf := &terraform.Terraform{}

	router.Path("/metrics").Handler(promhttp.Handler())

	router.HandleFunc("/", homeLinkHandler)

	// v1-arch handlers
	router.Handle("/calculoid/webhook", calc.CalculoidWebhookHandler())

	// v2-arch handlers
	//

	// proxmox/lxc
	router.HandleFunc("/proxmox-provisioning-server/container/all", pxm.ProvisioningServerGetContainerHandler)
	router.HandleFunc("/proxmox-provisioning-server/container/create", pxm.ProvisioningServerContainerCreateHandler)

	// monitoring handlers
	router.HandleFunc("/prometheus/remote/target/add", prom.RemoteTargetAddHandler)
	router.HandleFunc("/prometheus/remote/target/remove", prom.RemoteTargetRemoveHandler)

	// terraform handlers
	router.HandleFunc("/terraform/owncloudstack/create", tf.TerraformOwncloudstackCreateHandler)
	router.HandleFunc("/terraform/owncloudstackdocker/create", tf.TerraformOwncloudstackdockerCreateHandler)

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

	return loggedRouter
}
