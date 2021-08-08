package source

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Structure to declare a route of the application
type Route struct {
	Type    string
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

// Declaration of the routes available in the API
var routes = Routes{

	Route{"GET", "/v1/resource/get/{library}/{uuid}", ResourceGet},
	Route{"POST", "/v1/resource/upload", ResourceUpload},
}

// Create a mux router object to api the API
func RouterCreate() *mux.Router {
	var router *mux.Router = mux.NewRouter().StrictSlash(true)

	for _, value := range routes {
		router.Methods(value.Type).Path(value.Path).Handler(value.Handler)
	}

	return router
}
