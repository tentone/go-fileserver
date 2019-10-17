package api

import (
	"github.com/buaazp/fasthttprouter"
	"godonkey/api/file"
	"godonkey/api/image"
	"godonkey/api/resource"
	"godonkey/global"
)

// Create all the routes for service access using fast http router.
// Attention some routes may be prefixed by a version number and other are not.
// Routes that are not prefixed are assumed to always keep their functionality untouched.
func CreateRoutes(router *fasthttprouter.Router) {

	var v = "/v" + global.ApiVersion


	router.GET("/version", Version)
	if global.DevelopmentMode {
		router.GET("/log", SystemLog)
	}

	router.ServeFiles(v + "/resource/get/*filepath", global.DataPath)
	router.POST(v + "/resource/upload", resource.Upload)

	// swagger:route POST /v1/resource/delete Resources deleteResource
	//
	// Generic method to delete a resource from the server using its uuid and library name.
	// The files in the data path will be deleted and cannot be recovered.
	router.POST(v + "/resource/delete", resource.Delete)

	// swagger:route GET /v1/resource/get/uuid.format Resources getImage
	//
	// Image get service, used to obtain the image file initially uploaded.
	// The UUID and the format of the image are specified in the URL.
	router.GET(v + "/image/get/*filepath", image.Get)

	// swagger:route POST /v1/image/upload Images uploadImage
	//
	// Image upload service used to send and store image in binary format.
	// Payload is a multipart-form containing the file data and file format.
	// These files are stored using their UUID and file extension.
	// Response is a JSON with the UUID of the file stored e.g {"uuid":85d88263-5381-4857-9c84-43a91a739f5f}
	router.POST(v + "/image/upload", image.Upload)

	// swagger:route POST /v1/image/upload/base64 Images uploadImageBase64
	//
	// Image upload service using base64 encoded data to store the file.
	// Receives the data to be stored in a file should be avoided as much as possible since base64 encoded data is bigger.
	// Received a JSON with the "file" encoded as base64, and the corresponding "format".
	// Response is a JSON with the UUID of the file stored e.g {"uuid" 85d88263-5381-4857-9c84-43a91a739f5f}
	router.POST(v + "/image/upload/base64", image.UploadBase64)

	// swagger:route GET /v1/file/get/uuid.format Files getFile
	//
	// Get a file from the server using its UUID and its format.
	router.GET(v + "/file/get/*filepath", file.Get)

	// swagger:route POST /v1/file/upload Files uploadFile
	//
	// Generic file upload service.
	// Payload is a multipart-form containing the "file" data and "format" of the file.
	router.POST(v + "/file/upload", file.Upload)

	// swagger:route POST /v1/file/delete Files deleteFile
	//
	// Delete a file from its UUID.
	//
	// The files in the data path will be deleted and cannot be recovered.
	router.POST(v + "/file/delete", resource.Delete)

	// swagger:route GET /v1/potree/get Pointclouds getPotree
	//
	// Get a file from a potree pointcloud directory. These pointclouds are stored in multiple files inside of a folder.
	// The cloud.js file provides an index of all the subfolders available each folder contains a node of a spacially indexed octree structure and the index of the next level.
	router.ServeFiles(v + "/potree/get/*filepath", global.DataPath + global.POTREE)
}
