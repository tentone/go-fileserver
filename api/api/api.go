package api

import (
	"encoding/json"
	"github.com/tentone/godonkey/api/utils"
	"github.com/tentone/godonkey/global"
	"github.com/valyala/fasthttp"
)

// Get the server log, this API method is only reachable when running in development environment.
// Provides direct access to the log file to where the server writes messages in case of error.
// Can be used to easily access log data from the server
func SystemLog(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.SendFile(global.LogFile)
}

// Get the server instance version information, provides a timestamp of the build and the server that it is communicating with.
// Returns an object containing the version of the server and the api version level.
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
