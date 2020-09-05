package proxmox

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type Proxmox struct {
}

type UrlInputData struct {
	ProxmoxUrl string `yaml:"proxmoxUrl"`
}

var proxmoxServer string
var proxmoxPort string

const defaultProxmoxServer = "192.168.121.10"

func getProxmoxUrl(r *http.Request) string {
	var key string

	const param = "proxmoxUrl"

	if r.Method == http.MethodPost {
		var message UrlInputData

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		//TODO
		// fix rest of the unhandled cases
		// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

		err := dec.Decode(&message)
		if err != nil {
			log.Fatal("cannot decode body: ", err)
			return defaultProxmoxServer
		}

		key = message.ProxmoxUrl

		if key == "" {
			log.Printf("input param '%s' (%s) is missing", param, r.Method)
			return defaultProxmoxServer
		}
	} else {
		// GET method
		keys, ok := r.URL.Query()[param]

		if !ok || len(keys[0]) < 1 {
			log.Printf("input param '%s' (%s) is missing", param, r.Method)
			return defaultProxmoxServer
		}
		key = keys[0]
	}

	log.Printf("input param \"%s\": %s", param, key)

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

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
