package resource

import (
	"github.com/valyala/fasthttp"
	"godonkey/api/utils"
	"godonkey/global"
	"os"
	"strings"
)

// Generic method to retrieve a resource from the server. Has to used lower case UUID. Used for high bandwidth requests.
// Receives the library and UUID of the resource in the URL.
// The path of the service corresponds directly to the path of the file stored in the server.
func Get(ctx *fasthttp.RequestCtx) {
	var path = string(ctx.Path())
	var paths = strings.Split(path, "/")

	if len(paths) > 0 {
		paths[len(paths) - 1] = strings.ToLower(paths[len(paths) - 1])
	}

	var str = global.DataPath
	for i := 4; i < len(paths); i++ {
		str += paths[i]

		if (i + 1) < len(path) {
			str += "/"
		}
	}

	fasthttp.ServeFile(ctx, str)
}

// Generic resource upload service, that allows the user to specify the UUID of the resource.
// This method should be avoided as much as possible, it does not handle any resource class specific operations (data conversions, compression etc).
// Payload is a multipart-form containing the "file" data, "uuid" and "library" of the resource.
func Upload(ctx *fasthttp.RequestCtx) {

	// Read form data
	var file, err = ctx.FormFile("file")
	if err != nil {

		utils.SetErrorResponse(ctx, "No file provided in the form, check the file data.", fasthttp.StatusBadRequest, err)
		return
	}

	var uuid = strings.ToLower(string(ctx.FormValue("uuid")))
	if len(uuid) == 0 {

		utils.SetErrorResponse(ctx, "UUID is empty or missing.", fasthttp.StatusBadRequest, err)
		return
	}

	var library = strings.ToLower(string(ctx.FormValue("library")))
	if len(library) == 0 {

		utils.SetErrorResponse(ctx, "Library is empty or missing (e.g image, file).", fasthttp.StatusBadRequest, err)
		return
	}

	var format = strings.ToLower(string(ctx.FormValue("format")))
	if len(format) == 0 {

		utils.SetErrorResponse(ctx, "File format is empty or missing.", fasthttp.StatusBadRequest, err)
		return
	}

	// Make directory
	err = os.MkdirAll(global.DataPath + library, os.ModePerm)
	if err != nil {

		utils.SetErrorResponse(ctx, "Failed to create directory, check the server configuration.", fasthttp.StatusInternalServerError, err)
		return
	}

	var path string = global.DataPath + library + "/" + uuid + "." + format

	// Save file to the respective folder
	err = fasthttp.SaveMultipartFile(file, path)
	if err != nil {

		utils.SetErrorResponse(ctx, "Failed to store file, check the file data.", fasthttp.StatusBadRequest, err)
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {

		utils.SetErrorResponse(ctx, "File not created, check the file data.", fasthttp.StatusInternalServerError, err)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

// Generic method to delete a resource from the server using its uuid and library name.
// The files in the data path will be deleted and cannot be recovered.
func Delete(ctx *fasthttp.RequestCtx) {

	var uuid = string(ctx.FormValue("uuid"))
	var library = string(ctx.FormValue("library"))
	var format = string(ctx.FormValue("format"))

	var fname = global.DataPath + library + "/" + uuid + "." + format

	// Check if the file exists
	var _, err = os.Stat(fname)
	if err != nil && os.IsNotExist(err) {
		// Remove file
		err = os.Remove(fname)
		if err != nil {
			utils.SetErrorResponse(ctx, "Failed to delete the file, may not exist.", fasthttp.StatusInternalServerError, err)
			return
		}
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}
