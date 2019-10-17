package file

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	utils2 "godonkey/api/utils"
	"godonkey/global"
	"os"
	"strings"
)

func Get(ctx *fasthttp.RequestCtx) {

	var paths = strings.Split(string(ctx.Path()), "/")
	var path = strings.ToLower(paths[len(paths) - 1])

	fasthttp.ServeFile(ctx, global.DataPath + global.FILES + "/" + path)
}

func Upload(ctx *fasthttp.RequestCtx) {

	// Read form data
	var file, err = ctx.FormFile("file")
	if err != nil {

		utils2.SetErrorResponse(ctx, "No file provided in the form, check the file data.", fasthttp.StatusBadRequest, err)
		return
	}

	var format = string(ctx.FormValue("format"))
	if len(format) == 0 {

		utils2.SetErrorResponse(ctx, "File format is empty or missing.", fasthttp.StatusBadRequest, err)
		return
	}

	//Generate UUID
	var randomUuid, _ = uuid.NewRandom()
	var u string = randomUuid.String()
	var fname string = u + "." + format

	// Make directory
	err = os.MkdirAll(global.DataPath + global.FILES, os.ModePerm)
	if err != nil {

		ctx.Response.SetBodyString("Failed to create directory, check the server configuration.")
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	var path = global.DataPath + global.FILES + "/" + fname

	// Save file to the respective folder
	err = fasthttp.SaveMultipartFile(file, path)
	if err != nil {

		ctx.Response.SetBodyString("Failed to store file, check the file data.")
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {

		ctx.Response.SetBodyString("File not created, check the file data.")
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	var response = struct {
		UUID string `json:"uuid"`
	}{}
	response.UUID = u

	var data []byte

	data, err = json.Marshal(&response)
	if err != nil {

		utils2.SetErrorResponse(ctx, "Error creating JSON response.", fasthttp.StatusInternalServerError, err)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBody(data)
}


func Delete(ctx *fasthttp.RequestCtx) {

	var u = string(ctx.FormValue("uuid"))
	var fname = global.DataPath + global.FILES + "/" + u

	// Check if the file exists
	var _, err = os.Stat(fname)
	if err != nil && os.IsNotExist(err) {

		// Remove file
		err = os.Remove(fname)
		if err != nil {

			utils2.SetErrorResponse(ctx, "Failed to delete the file, may not exist.", fasthttp.StatusInternalServerError, err)
			return
		}
	}
}
