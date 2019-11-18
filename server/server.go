package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/google/logger"
	"github.com/tentone/godonkey/global"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

// Structure represents the server entry point.
//
// Stores pointers to all configurations required
type Server struct {}

// Start the server using data from the configurations structures provided.
func (s Server) Start() {
	var router = RouterCreate()
	var err error

	// Generate certificate for localhost
	if global.Config.Server.GenerateCertTLS {
		var cert, privateKey []byte
		cert, privateKey, err = GenerateCertificate("localhost", "Local")
		if err != nil {
			logger.Fatal("Failed to generate certificate file.")
		}

		// Write certificate files
		_ = ioutil.WriteFile(global.Config.Server.CertFileTLS, cert, 0644)
		_ = ioutil.WriteFile(global.Config.Server.KeyFileTLS, privateKey, 0644)
	}

	// Start the server with TLS, since we are running HTTP/2 it must be run with TLS.
	if len(global.Config.Server.AddressTLS) > 0 {
		var server http.Server = http.Server{
			Addr: global.Config.Server.AddressTLS,
			Handler: router,
		}

		err = server.ListenAndServeTLS(global.Config.Server.CertFileTLS, global.Config.Server.KeyFileTLS)
		if err != nil {
			logger.Fatal("Failed to start HTTPS/H2 server.")
		}
	}

	// Start the server using HTTP
	if len(global.Config.Server.Address) > 0 {
		var server http.Server = http.Server{
			Addr: global.Config.Server.Address,
			Handler: router,
		}
		err = server.ListenAndServe()
		if err != nil {
			logger.Fatal("Failed to start HTTP server.")
		}
	}
}

/// Generate a TLS certificate from host name.
///
/// Should only be used for localhost testing and development purposes.
///
/// Returns Certificate data, Private key data and error if any occurred.
func GenerateCertificate(host string, organization string) ([]byte, []byte, error) {
	var err error
	var priv *rsa.PrivateKey
	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	var serialNumberLimit *big.Int
	var serialNumber *big.Int
	serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	var cert *x509.Certificate
	cert = &x509.Certificate{
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

	var certBytes []byte
	certBytes, err = x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	var p []byte
	p = pem.EncodeToMemory(
		&pem.Block{
			Type: "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	var b []byte
	b = pem.EncodeToMemory(
		&pem.Block{
			Type: "CERTIFICATE",
			Bytes: certBytes,
		},
	)

	return b, p, nil
}

