package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/buaazp/fasthttprouter"
	"github.com/google/logger"
	"github.com/valyala/fasthttp"
	"math/big"
	"os"
	"time"
	"go-donkey/api"
	"go-donkey/database"
	"go-donkey/global"
)

func main() {

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

	logger.Init("Resources", global.DevelopmentMode, false, file)

	// Read server configuration
	global.LoadConfig()
	database.Load()

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
				var certData, keyCertData, err = GenerateCertificate("resources.unodigital.io", "UNO Digital")
				if err != nil {
					logger.Error("Error generating certificate." + err.Error())
				}

				err = http.ListenAndServeTLSEmbed(global.AddressTLS, certData, keyCertData)

				if err != nil {
					logger.Error("Error starting HTTPS server." + err.Error())
				} else {
					logger.Info("Server HTTPS started and listing at %s.", global.AddressTLS)
				}
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

/// Generate a TLS certificate from host name.
///
/// Should only be used for localhost testing and development purposes.
///
/// Returns Certificate data, Private key data and error if any occurred.
func GenerateCertificate(host string, organization string) ([]byte, []byte, error) {

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return nil, nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{organization},
		},
		NotBefore: time.Now(),
		NotAfter: time.Now().Add(365 * 24 * time.Hour),
		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm: x509.SHA256WithRSA,
		DNSNames: []string{host},
		BasicConstraintsValid: true,
		IsCA: true,
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader, cert, cert, &priv.PublicKey, priv,
	)

	p := pem.EncodeToMemory(
		&pem.Block{
			Type: "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	b := pem.EncodeToMemory(
		&pem.Block{
			Type: "CERTIFICATE",
			Bytes: certBytes,
		},
	)

	return b, p, err
}
