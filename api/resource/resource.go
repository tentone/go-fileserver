package resource

import (
	"github.com/valyala/fasthttp"
	"os"
	"strings"
	"godonkey/database"
	"godonkey/global"
	"godonkey/utils"
)

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
	err = os.MkdirAll(global.DataPath+ library, os.ModePerm)
	if err != nil {

		utils.SetErrorResponse(ctx, "Failed to create directory, check the server configuration.", fasthttp.StatusInternalServerError, err)
		return
	}

	var path string = global.DataPath + library + "/" + uuid

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

	database.InsertResource(library, uuid, format, uuid)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func Delete(ctx *fasthttp.RequestCtx) {

	var uuid = string(ctx.FormValue("uuid"))
	var library = string(ctx.FormValue("library"))
	var fname = global.DataPath + library + "/" + uuid

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

	database.RemoveResource(library, uuid)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}
