package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Structure to declare a route of the application
type Route struct {
	// Type of the API route.
	Type string

	// Path of the API.
	Path string

	// Handler to process the API call
	Handler http.HandlerFunc
}

// Declaration of the routes available in the API
var routes []Route = []Route{
	{"GET", "/v1/resource/get/{library}/{uuid}", ResourceGet},
	{"POST", "/v1/resource/upload", ResourceUpload},
}

// Create a mux router object to api the API
func RouterCreate() *mux.Router {
	var router *mux.Router = mux.NewRouter().StrictSlash(true)

	for _, value := range routes {
		router.Methods(value.Type).Path(value.Path).Handler(value.Handler)
	}

	return router
}
