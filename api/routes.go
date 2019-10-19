package api

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/tentone/godonkey/api/api"
	"github.com/tentone/godonkey/api/resource"
	"github.com/tentone/godonkey/global"
)

// Create all the routes for service access using fast http router.
//
// Attention some routes may be prefixed by a version number and other are not.
//
// Routes that are not prefixed are assumed to always keep their functionality untouched.
func CreateRoutes(router *fasthttprouter.Router) {

	var v = "/v" + global.ApiVersion

	// Version
	router.GET("/version", api.Version)
	if global.DevelopmentMode {
		router.GET("/log", api.SystemLog)
	}

	// Resources
	router.GET(v + "/resource/get/*filepath", resource.Get)
	router.POST(v + "/resource/upload", resource.Upload)
	router.POST(v + "/resource/delete", resource.Delete)
}
