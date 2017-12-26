package minit

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"github.com/amlun/minit/http/controller"
	"github.com/amlun/minit/services"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
	router.HandleFunc("/", controller.DefaultHandler)
	router.HandleFunc("/users", controller.UsersGetHandler).Methods("GET")
	router.HandleFunc("/users", controller.UsersAddHandler).Methods("POST")
	router.HandleFunc("/users/{user_id}/relationships", controller.UsersRelationshipsGetHandler).Methods("GET")
	router.HandleFunc("/users/{owner_id}/relationships/{user_id}", controller.UsersRelationshipsAddHandler).Methods("PUT")
}

// Run minit server, serve http service
func Run() {
	// use env or flag?
	services.Init("postgres://lunweiwei:123456@127.0.0.1:5432/minit?sslmode=disable")
	router.Headers("Content-Type", "application/json")
	log.Fatal(http.ListenAndServe(":8080", router))
}
