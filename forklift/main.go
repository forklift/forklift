package main

import (
	"flag"
	"log"
	"net/http"
)

var port = flag.String("port", "7070", "Define what TCP port to bind to")

func main() {

	flag.Parse()
	endpoint := ":" + *port

	mux := http.NewServeMux()
	mux.HandleFunc("/_ping", pong)

	log.Printf("Listening at %s", endpoint)

	log.Fatal(http.ListenAndServe(endpoint, mux))
}

func pong(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}
