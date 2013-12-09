package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	// handle belkin related commands
	ch := make(chan BelkinRequest, 20)
	go BelkinProcessor(ch)
	router.Methods("POST").Subrouter().HandleFunc("/api/belkin/{name}/{action}", BelkinHandler(ch))

	// default to file server
	router.Methods("GET").Handler(http.FileServer(http.Dir("public")))

	http.ListenAndServe(":8080", router)
}
