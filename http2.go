package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

func ListenAndServeTLS(srv *http.Server, certPEMBlock, keyPEMBlock []byte) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	return srv.Serve(tlsListener)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted connections.
// It's used by ListenAndServe and ListenAndServeTLS so dead TCP connections (e.g. closing laptop mid-download) eventually go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func main() {
	var server http.Server
	server = http.Server{Addr: ":9090", Handler: http.HandlerFunc(handle)}

	var cert, privateKey []byte
	var err error
	cert, privateKey, err = GenerateCertificate("localhost", "unodigital")
	if err != nil {
		log.Fatal("Failed to generate certificate file.")
		return
	}

	// Start the server with TLS, since we are running HTTP/2 it must be run with TLS.
	log.Printf("Serving on https://0.0.0.0:8000")

	err = ListenAndServeTLS(&server, cert, privateKey)
	if err != nil {
		log.Fatal("Failed to start server.")
		return
	}
}