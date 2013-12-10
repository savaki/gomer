package main

import (
	"flag"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	port     = flag.String("port", "8080", "the port number to use")
	username = flag.String("hue", "", "the hue username")
)

func main() {
	flag.Parse()

	router := mux.NewRouter()

	// handle belkin related commands
	belkinCh := make(chan BelkinRequest, 20)
	go BelkinProcessor(belkinCh)
	router.HandleFunc("/api/belkin/{name}/{action}", BelkinHandler(belkinCh)).Methods("POST")

	// handle hue related commands
	hueCh := make(chan HueRequest, 20)
	go HueProcessor(*username, hueCh)
	router.HandleFunc("/api/hue/{name}/{action}", HueHandler(hueCh)).Methods("POST")

	// default to file server
	router.Methods("GET").Handler(http.FileServer(http.Dir("public")))

	http.ListenAndServe(":"+*port, router)
}
