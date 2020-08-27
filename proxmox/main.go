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
	var key string

	const param = "proxmoxUrl"

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatal("cannot parse Form: ", err)
		}

		// TODO
		// fix empty map
		log.Println(r.PostForm)

		key = r.Form.Get(param)
		if key == "" {
			log.Printf("Url param '%s' (%s) is missing", param, r.Method)
			return defaultProxmoxServer
		}
	} else {
		keys, ok := r.URL.Query()[param]

		if !ok || len(keys[0]) < 1 {
			log.Printf("Url param '%s' (%s) is missing", param, r.Method)
			return defaultProxmoxServer
		}
		key = keys[0]
	}

	log.Printf("Url param \"%s\": %s", param, key)

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

	if port == "" {
		port = "4567"
	}
	if host == "" {
		host = defaultProxmoxServer
	}
	return host, port, nil
}
