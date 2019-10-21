package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MAX_UPLOAD_SIZE = 32 * 200000 * 1024
const DATA_PATH = "./data/"

func main() {

	var router = Create()

	var cfg *tls.Config = &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	var server http.Server = http.Server{
		Addr: ":9090",
		Handler: router,
		TLSConfig: cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// Generate certificate for localhost
	var cert, privateKey []byte
	var err error
	cert, privateKey, err = GenerateCertificate("localhost", "UNO Digital")
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

type Route struct {
	Verb string
	Path string
	HandlerFunction http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"GET", "/resource/get/{library}/{uuid}", ResourceGet},
	Route{"POST", "/resource/upload", ResourceUpload},
}

func Create() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunction
		router.Methods(route.Verb).Path(route.Path).Handler(handler)
	}

	return router
}

func ResourceGet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/octet-stream; charset=UTF-8")

	// Form data
	var variables = mux.Vars(request)
	var uuid string = variables["uuid"]
	var library string = variables["library"]
	var fileLocation string = filepath.Join(DATA_PATH, library, uuid)

	// Read file
	var err error
	var file []byte
	file, err = ioutil.ReadFile(fileLocation)

	if err != nil {
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(200)
	_, _ = writer.Write(file)
}

func ResourceUpload(writer http.ResponseWriter, request *http.Request) {
	var uuid = request.FormValue("uuid")
	var library = request.FormValue("library")

	var fileLocation = filepath.Join(DATA_PATH, strings.ToLower(library), uuid)

	request.Body = http.MaxBytesReader(writer, request.Body, MAX_UPLOAD_SIZE)

	var err error
	err = request.ParseMultipartForm(MAX_UPLOAD_SIZE);
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		err = os.MkdirAll(fileLocation, 0755)
		if err != nil {
			panic(err)
		}
	}

	file, _, err := request.FormFile("file")
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	newFile, err := os.Create(fileLocation)
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil {
		writer.WriteHeader(500)
		return
	}
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
