package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var receivedData []byte

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
}

func CalculoidWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		code := http.StatusBadRequest
		if r.Method == "GET" {
			code = http.StatusOK
			result := fmt.Sprintln("CalculoidWebhook")
			_, err = w.Write([]byte(result))
			if err != nil {
				log.Fatal(err)
			}
		} else if r.Method == "POST" {
			code = http.StatusOK
			result := fmt.Sprintln("CalculoidWebhook")
			_, err = w.Write([]byte(result))
			if err != nil {
				log.Fatal(err)
			}
			queryParams(w, r)
			calculoidWebhookParser()
		} else {
			w.WriteHeader(code)
		}
	}
}

func queryParams(w http.ResponseWriter, r *http.Request) {
	var err error
	receivedData, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("received data: %s \n", receivedData)
}

func calculoidWebhookParser() {
	log.Printf("received data for parsing: %s \n", receivedData)
}

func main() {

	filename := os.Args[0]
	log.Printf("starting %s \n", filename)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.Handle("/calculoid/webhook", CalculoidWebhook())
	log.Fatal(http.ListenAndServe(":8080", router))
}
