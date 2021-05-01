package main_test

import (
	"cc-server/server"
	prometheusmiddleware "github.com/albertogviana/prometheus-middleware"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var s *httptest.Server
var serverURL string

func TestCCServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "c&c server")
}

var _ = BeforeSuite(func() {

	var opts prometheusmiddleware.Opts
	middleware := prometheusmiddleware.NewPrometheusMiddleware(opts)

	tracer, _ := server.CreateTracer()

	r := server.CreateRouter(middleware, tracer)
	s = httptest.NewServer(r)

	Expect(len(s.URL)).To(BeNumerically(">", 0))
	serverURL = s.URL
	log.Print("### " + serverURL + " ###\n\n")

})

var _ = Describe("c&c server", func() {

	Describe("prometheus metrics", func() {
		var url string

		BeforeEach(func() {
			url = serverURL + "/metrics"
		})

		Context("when endpoint exists", func() {

			var rdr *strings.Reader
			var req *http.Request
			var res *http.Response
			var err error

			log.Print(url)

			It("Makes a GET request", func() {
				rdr = strings.NewReader("")
				req, err = http.NewRequest("GET", url, rdr)
				Expect(err).NotTo(HaveOccurred())
			})

			It("retrieves a response", func() {
				res, err = http.DefaultClient.Do(req)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Returns HTTP 200 OK", func() {
				Expect(res.StatusCode).To(BeNumerically("==", http.StatusOK))
			})

		})
	})
})
