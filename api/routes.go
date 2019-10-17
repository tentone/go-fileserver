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

	// swagger:route GET /version Meta version
	//
	// Get the server instance version information, provides a timestamp of the build and the server that it is communicating with.
	// Returns an object containing the version of the server and the api version level.
	//
	// Responses:
	//   200: versionInfo
	router.GET("/version", Version)


	// swagger:route GET /log Meta log
	//
	// Get the server log, this API method is only reachable when running in development environment.
	// Provides direct access to the log file to where the server writes messages in case of error.
	// Can be used to easily access log data from the server
	//
	// Responses:
	//    200: successInfo
	if global.DevelopmentMode {
		router.GET("/log", SystemLog)
	}

	// swagger:route GET /v1/resource/get/uuid.format Resources getResource
	//
	// Generic method to retrieve a resource from the server. Has to used lower case UUID. Used for high bandwidth requests.
	// Receives the library and UUID of the resource in the URL.
	// The path of the service corresponds directly to the path of the file stored in the server.
	//
	// Responses:
	//    200: successInfo
	router.ServeFiles(v + "/resource/get/*filepath", global.DataPath)

	// swagger:route POST /v1/resource/upload Resources uploadResource
	//
	// Generic resource upload service, that allows the user to specify the UUID of the resource.
	// This method should be avoided as much as possible, it does not handle any resource class specific operations (data conversions, compression etc).
	// Payload is a multipart-form containing the "file" data, "uuid" and "library" of the resource.
	//
	// Responses:
	//    200: uuidInfo
	router.POST(v + "/resource/upload", resource.Upload)

	// swagger:route POST /v1/resource/delete Resources deleteResource
	//
	// Generic method to delete a resource from the server using its uuid and library name.
	// The files in the data path will be deleted and cannot be recovered.
	//
	// Responses:
	//    200: successInfo
	router.POST(v + "/resource/delete", resource.Delete)

	// swagger:route GET /v1/resource/get/uuid.format Resources getImage
	//
	// Image get service, used to obtain the image file initially uploaded.
	// The UUID and the format of the image are specified in the URL.
	//
	// Responses:
	//    200: successInfo
	router.GET(v + "/image/get/*filepath", image.Get)

	// swagger:route POST /v1/image/upload Images uploadImage
	//
	// Image upload service used to send and store image in binary format.
	// Payload is a multipart-form containing the file data and file format.
	// These files are stored using their UUID and file extension.
	// Response is a JSON with the UUID of the file stored e.g {"uuid":85d88263-5381-4857-9c84-43a91a739f5f}
	//
	// Responses:
	//    200: uuidInfo
	router.POST(v + "/image/upload", image.Upload)

	// swagger:route POST /v1/image/upload/base64 Images uploadImageBase64
	//
	// Image upload service using base64 encoded data to store the file.
	// Receives the data to be stored in a file should be avoided as much as possible since base64 encoded data is bigger.
	// Received a JSON with the "file" encoded as base64, and the corresponding "format".
	// Response is a JSON with the UUID of the file stored e.g {"uuid" 85d88263-5381-4857-9c84-43a91a739f5f}
	//
	//
	// Responses:
	//    200: uuidInfo
	router.POST(v + "/image/upload/base64", image.UploadBase64)

	// swagger:route GET /v1/file/get/uuid.format Files getFile
	//
	// Get a file from the server using its UUID and its format.
	//
	// Responses:
	//    200: successInfo
	router.GET(v + "/file/get/*filepath", file.Get)

	// swagger:route POST /v1/file/upload Files uploadFile
	//
	// Generic file upload service.
	// Payload is a multipart-form containing the "file" data and "format" of the file.
	//
	// Responses:
	//    200: uuidInfo
	router.POST(v + "/file/upload", file.Upload)

	// swagger:route POST /v1/file/delete Files deleteFile
	//
	// Delete a file from its UUID.
	//
	// The files in the data path will be deleted and cannot be recovered.
	//
	// Responses:
	//    200: successInfo
	router.POST(v + "/file/delete", resource.Delete)

	// swagger:route GET /v1/potree/get Pointclouds getPotree
	//
	// Get a file from a potree pointcloud directory. These pointclouds are stored in multiple files inside of a folder.
	// The cloud.js file provides an index of all the subfolders available each folder contains a node of a spacially indexed octree structure and the index of the next level.
	//
	// Responses:
	//    200: successInfo
	router.ServeFiles(v + "/potree/get/*filepath", global.DataPath + global.POTREE)

	// OK message.
	// swagger:response successInfo
	type SuccessInfo struct {}

	// Message describing the error.
	// swagger:response errorInfo
	type ErrorInfo struct {}

	// Returns a JSON with the UUID of the file stored.
	// swagger:response uuidInfo
	type UuidInfo struct {}

	// Returns a JSON indicating the version information.
	// swagger:response versionInfo
	type VersionInfo struct {}
}
