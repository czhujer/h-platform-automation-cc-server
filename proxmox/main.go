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

	log.Printf("Url Param \"%s\" is: %s", param, key)

	return key
}

func setProxmoxUrl(r *http.Request) (string, string, error) {
	var (
		proxmoxUrl string
		host       string
		port       string
	)
	proxmoxUrl = getProxmoxUrl(r)

	parsedUrl, err := url.Parse("http://" + proxmoxUrl)
	if err != nil {
		log.Fatal("cannot parse proxmox hostname and port: ", err)
	}

	host = parsedUrl.Hostname()
	port = parsedUrl.Port()

	//log.Printf("parsedHost: ", parsedHost)

	//host, port, _ := net.SplitHostPort(parsedHost)

	//if err != nil {
	//	log.Fatal("cannot set proxmox URL: ", err)
	//	return "", "", err
	//}
	if port == "" {
		port = "4567"
	}
	if host == "" {
		host = defaultProxmoxServer
	}
	return host, port, nil
}
