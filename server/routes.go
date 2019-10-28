package server

import (
	"github.com/gorilla/mux"
	"github.com/tentone/godonkey/server/api"
	"net/http"
)

type Route struct {
	Type string
	Path string
	Handler http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"GET", "/v1/resource/get/{library}/{uuid}", api.ResourceGet},
	Route{"POST", "/v1/resource/upload", api.ResourceUpload},
}

func RouterCreate() *mux.Router {
	var router *mux.Router = mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Type).Path(route.Path).Handler(route.Handler)
	}

	return router
}


