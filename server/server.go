package server

import (
	"github.com/tentone/godonkey/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func ServerStart() {
	var router = RouterCreate()

	/*var cfg *tls.Config = &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}*/

	var server http.Server = http.Server{
		Addr: ":9090",
		Handler: router,
		//TLSConfig: cfg,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// Generate certificate for localhost
	var cert, privateKey []byte
	var err error
	cert, privateKey, err = utils.GenerateCertificate("localhost", "Local")
	if err != nil {
		log.Fatal("Failed to generate certificate file.")
		return
	}

	// Write certificate files
	_ = ioutil.WriteFile("./local.crt", cert, 0644)
	_ = ioutil.WriteFile("./local.key", privateKey, 0644)

	// Start the server with TLS, since we are running HTTP/2 it must be run with TLS.
	err = server.ListenAndServeTLS("local.crt", "local.key")
	if err != nil {
		log.Fatal("Failed to start server.")
		return
	}
}