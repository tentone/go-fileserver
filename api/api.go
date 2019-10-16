package api

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"go-donkey/global"
	"go-donkey/utils"
)

func SystemLog(ctx *fasthttp.RequestCtx) {

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.SendFile(global.LogFile)
}

func Version(ctx *fasthttp.RequestCtx) {

	var response = struct {

		Version string `json:"version"`
		Api string `json:"api"`
		Commit string `json:"commit"`
	}{}

	response.Version = global.Version
	response.Api = global.ApiVersion
	response.Commit = global.GitCommit

	var data, err = json.Marshal(&response)
	if err != nil {

		utils.SetErrorResponse(ctx, "Error creating JSON response.", fasthttp.StatusInternalServerError, err)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBody(data)
}
