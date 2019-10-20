package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/google/logger"
	"github.com/tentone/godonkey/api"
	"github.com/tentone/godonkey/global"
	"github.com/valyala/fasthttp"
	"os"
	"time"
)

func maina() {

	global.LoadVersion()

	var args = os.Args[1:]

	if len(args) > 0 {
		global.ConfigurationFile = args[0]
	}
	if len(args) > 1 {
		global.LogFile = args[1]
	}

	// Setup logger
	var file, err = os.OpenFile(global.LogFile, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file.")
	}

	logger.Init("GoDonkey", global.DevelopmentMode, true, file)

	// Read server configuration
	global.LoadConfig()

	// Start HTTP server.
	if len(global.Address) > 0 {

		logger.Infof("Starting HTTP server on %q", global.Address)

		go func() {
			var router = fasthttprouter.New()
			api.CreateRoutes(router)

			var http = fasthttp.Server{
				Handler: HandleCORS(router.Handler),
				MaxRequestBodySize: 2147483648, //2 ^ 31 = 2GB
				ReadTimeout: time.Duration(10) * time.Minute,
				WriteTimeout: time.Duration(10) * time.Minute,
				ReadBufferSize: 65536, //2 ^ 16 = 64MB,
				WriteBufferSize: 65536,
				MaxConnsPerIP: 1e5,
				MaxRequestsPerConn: 1e10,
				DisableHeaderNamesNormalizing: false,
				SleepWhenConcurrencyLimitsExceeded: 0,
				NoDefaultServerHeader: true,
				NoDefaultContentType: true,
				ReduceMemoryUsage: false,
				TCPKeepalive: false,
				DisableKeepalive: false,
				Concurrency: 262144,
			}

			var err = http.ListenAndServe(global.Address)

			if err != nil {
				logger.Error("Error starting HTTP server." + err.Error())
			} else {
				logger.Info("Server HTTP started and listing at %s.", global.Address)
			}
		}()
	}

	// Start HTTPS server.
	if len(global.AddressTLS) > 0 {

		logger.Infof("Starting HTTPS server on %q", global.AddressTLS)

		go func() {
			var router = fasthttprouter.New()
			api.CreateRoutes(router)

			var http = fasthttp.Server{
				Handler: HandleCORS(router.Handler),
				MaxRequestBodySize: 2147483648, //2 ^ 31 = 2GB
				ReadTimeout: time.Duration(10) * time.Minute,
				WriteTimeout: time.Duration(10) * time.Minute,
				ReadBufferSize: 65536, //2 ^ 16 = 64MB,
				WriteBufferSize: 65536,
				MaxConnsPerIP: 1e5,
				MaxRequestsPerConn: 1e10,
				DisableHeaderNamesNormalizing: false,
				SleepWhenConcurrencyLimitsExceeded: 0,
				NoDefaultServerHeader: true,
				NoDefaultContentType: true,
				ReduceMemoryUsage: false,
				TCPKeepalive: false,
				DisableKeepalive: false,
				Concurrency: 262144,
			}

			if len(global.CertFileTLS) == 0 || len(global.KeyFileTLS) == 0 {
				/*var certData, keyCertData, err = GenerateCertificate("resources.unodigital.io", "UNO Digital")
				if err != nil {
					logger.Error("Error generating certificate." + err.Error())
				}

				err = http.ListenAndServeTLSEmbed(global.AddressTLS, certData, keyCertData)

				if err != nil {
					logger.Error("Error starting HTTPS server." + err.Error())
				} else {
					logger.Info("Server HTTPS started and listing at %s.", global.AddressTLS)
				}*/
			} else {
				var err = http.ListenAndServeTLS(global.AddressTLS, global.CertFileTLS, global.KeyFileTLS)

				if err != nil {
					logger.Error("Error starting HTTPS server." + err.Error())
				} else {
					logger.Info("Server HTTPS started and listing at %s.", global.AddressTLS)
				}
			}
		}()
	}

	select {}
}

// CORS handler middleware.
// Sets the context response access-control headers.
func HandleCORS(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "HEAD,GET,POST,PUT,DELETE,OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

		handler(ctx)
	})
}
