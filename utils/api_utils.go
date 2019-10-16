package utils

import (
	"github.com/google/logger"
	"github.com/valyala/fasthttp"
)

func SetErrorResponse(ctx  *fasthttp.RequestCtx, message string, status int, err error) {

	if err != nil {
		message += " " + err.Error()
	}

	logger.Error(message)
	ctx.Response.SetBodyString(message)
	ctx.SetStatusCode(status)
}