package main

import (
	"collector/internal/controller"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	// TODO: create conf, app and run app
	r := mux.NewRouter()
	r.HandleFunc("/", controller.HandleConnection)
	http.ListenAndServe("127.0.0.1:8382", r)
}
