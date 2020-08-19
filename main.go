package main

import (
	"cc-server/calculoid"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
}

func main() {

	//calculoidHandler := &calculoid.CalculoidWebhook()

	filename := os.Args[0]
	log.Printf("starting %s \n", filename)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.Handle("/calculoid/webhook", calculoid.CalculoidWebhook())
	log.Fatal(http.ListenAndServe(":8080", router))
}
