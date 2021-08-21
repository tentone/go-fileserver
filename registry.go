package main

import "github.com/gorilla/mux"

// Create database tables
func RegistryDatabaseMigrate() {
	LogMigrate(db)
	LibraryMigrate(db)
	ResourceMigrate(db)
}

// Declaration of the routes available in the API
var routes []Route = []Route{
	{"GET", "/v1/resource/get/{library}/{uuid}", ResourceGet},
	{"POST", "/v1/resource/upload", ResourceUpload},
}

// Create a mux router object to api the API
func RegistryRouter() *mux.Router {
	var router *mux.Router = mux.NewRouter().StrictSlash(true)

	for _, value := range routes {
		router.Methods(value.Type).Path(value.Path).Handler(value.Handler)
	}

	router.Use(HandleCORS)
	router.Use(HandleOptions)

	return router
}
