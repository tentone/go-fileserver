package image

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"os"
	"strings"
	"go-donkey/database"
	"go-donkey/global"
	"go-donkey/utils"
)

func Get(ctx *fasthttp.RequestCtx) {

	var paths = strings.Split(string(ctx.Path()), "/")
	var path = strings.ToLower(paths[len(paths) - 1])

	fasthttp.ServeFile(ctx, global.DataPath + global.IMAGES + "/" + path)
}

func Upload(ctx *fasthttp.RequestCtx) {

	var file, err = ctx.FormFile("file")
	if err != nil {

		utils.SetErrorResponse(ctx, "No file provided in the form, check the file data.", fasthttp.StatusBadRequest, err)
		return
	}

	var format = strings.ToLower(string(ctx.FormValue("format")))

	if len(format) == 0 {

		utils.SetErrorResponse(ctx, "File format is empty or missing.", fasthttp.StatusBadRequest, err)
		return
	}

	logger.Info("Format received: " + format + ", length: " + string(len(format)))

	//Generate UUID
	var randomUuid, _ = uuid.NewRandom()
	var uuid = randomUuid.String()
	var fname = uuid + "." + format

	//Make directory
	err = os.MkdirAll(global.DataPath+ global.IMAGES, os.ModePerm)
	if err != nil {

		utils.SetErrorResponse(ctx, "Failed to create directory, check the server configuration.", fasthttp.StatusInternalServerError, err)
		return
	}

	var path = global.DataPath + global.IMAGES + "/" + fname

	//Save file data
	err = fasthttp.SaveMultipartFile(file, path)
	if err != nil {

		utils.SetErrorResponse(ctx, "Failed to store file, check the file data.", fasthttp.StatusBadRequest, err)
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {

		utils.SetErrorResponse(ctx, "File not created, check the file data.", fasthttp.StatusInternalServerError, err)
		return
	}

	database.InsertResource(global.IMAGES, uuid, format, fname)

	var response = struct {

		UUID string `json:"uuid"`
	}{}
	response.UUID = uuid

	var data []byte

	data, err = json.Marshal(&response)
	if err != nil {

		utils.SetErrorResponse(ctx, "Error creating JSON response.", fasthttp.StatusInternalServerError, err)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBody(data)
}

func UploadBase64(ctx *fasthttp.RequestCtx) {

	var received = struct {

		Format string `json:"format"`
		File string `json:"file"`
	}{}

	var err = json.Unmarshal(ctx.PostBody(), &received)
	if err != nil {

		utils.SetErrorResponse(ctx, "Error parsing JSON data received.", fasthttp.StatusBadRequest, err)
		return
	}

	var randomUuid, _ = uuid.NewRandom()
	var uuid = randomUuid.String()
	var format = received.Format

	var data []byte

	//Convert base64 data
	data, err = base64.StdEncoding.DecodeString(received.File)
	if err != nil {

		utils.SetErrorResponse(ctx, "Failed to decode base64 data.", fasthttp.StatusInternalServerError, err)
		return
	}

	//Make directory
	var path = global.DataPath + global.IMAGES + "/" + uuid + "." + format
	err = os.MkdirAll(global.DataPath+ global.IMAGES, os.ModePerm)
	if err != nil {

		utils.SetErrorResponse(ctx, "Failed to create directory.", fasthttp.StatusInternalServerError, err)
		return
	}

	//Write file
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {

		utils.SetErrorResponse(ctx, "Failed to write file to the server.", fasthttp.StatusInternalServerError, err)
		return
	}

	database.InsertResource(global.IMAGES, uuid, format, uuid)

	var response = struct {

		UUID string `json:"uuid"`
	}{}
	response.UUID = uuid

	var dataResponse []byte
	dataResponse, err = json.Marshal(&response)
	if err != nil {

		utils.SetErrorResponse(ctx, "Error creating JSON response.", fasthttp.StatusInternalServerError, err)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBody(dataResponse)
}
