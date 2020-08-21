package proxmox

import (
	"fmt"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (proxmox *Proxmox) proxmoxProvisioningServerClient(tracer opentracing.Tracer, action string, proxmoxServer string) (bool, string) {
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
			fmt.Sprintf("http://%s:%s%s", proxmoxServer, "4567", "/api/containers"),
			nil,
		)
		if err != nil {
			proxmox.onError(span, err)
			return false, ""
		}
	} else if action == "create" {

		req, err = http.NewRequest(
			"POST",
			fmt.Sprintf("http://%s:%s%s", proxmoxServer, "4567", "/api/containers/create"),
			strings.NewReader(requestBodyCreate),
		)
		if err != nil {
			proxmox.onError(span, err)
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
		proxmox.onError(span, err)
		return false, ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		proxmox.onError(span, err)
		return false, ""
	}
	log.Printf("Received result: %s\n", string(body))
	return true, string(body)
}

func (proxmox *Proxmox) onError(span opentracing.Span, err error) {
	// handle errors by recording them in the span
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
	log.Print(err)
}
