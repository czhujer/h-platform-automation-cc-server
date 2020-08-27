package proxmox

import (
	"log"
	"net/http"
	"net/url"
)

type Proxmox struct {
}

var proxmoxServer string
var proxmoxPort string

const defaultProxmoxServer = "192.168.121.10"

func getProxmoxUrl(r *http.Request) string {
	param := "proxmoxUrl"
	keys, ok := r.URL.Query()[param]

	if !ok || len(keys[0]) < 1 {
		log.Printf("Url Param '%s' is missing", param)
		return defaultProxmoxServer
	}

	key := keys[0]

	log.Printf("Url Param 's%' is: "+string(key), param)

	return string(key)
}

func setProxmoxUrl(r *http.Request) (string, string, error) {
	var (
		proxmoxUrl string
		host       string
		port       string
	)
	proxmoxUrl = getProxmoxUrl(r)

	// host, port, err := net.SplitHostPort(url)
	// https://golang.org/pkg/net/url/#Parse
	parsedUrl, err := url.Parse(proxmoxUrl)
	if err != nil {
		log.Fatal("cannot parse proxmox hostname: ", err)
	}

	host = parsedUrl.Hostname()
	port = parsedUrl.Port()

	if err != nil {
		log.Fatal("cannot set proxmox URL: ", err)
		return "", "", err
	}
	if port == "" {
		port = "4567"
	}
	if host == "" {
		host = defaultProxmoxServer
	}
	return host, port, nil
}
