package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	var router = NewRouter()

	var server http.Server
	server = http.Server{Addr: ":9090", Handler: router}

	var cert, privateKey []byte
	var err error
	cert, privateKey, err = GenerateCertificate("localhost", "unodigital")
	if err != nil {
		log.Fatal("Failed to generate certificate file.")
		return
	}

	// Start the server with TLS, since we are running HTTP/2 it must be run with TLS.
	log.Printf("Serving...")

	err = ListenAndServeTLS(&server, cert, privateKey)
	if err != nil {
		log.Fatal("Failed to start server.")
		return
	}
}

func ListenAndServeTLS(srv *http.Server, certPEMBlock, keyPEMBlock []byte) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	//TODO <CHECK THIS>
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


type Route struct {
	Name            string
	Verb            string
	Path            string
	HandlerFunction http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"GetGeneric", "GET", "/scene/{uuidscene}/{library}/{uuidresource}", GenericGet},
	Route{"UploadResource", "POST", "/upload", UploadResource},
	Route{"Ping", "GET", "ping", Ping},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunction
		handler = LogAccess(handler, route.Name)
		router.Methods(route.Verb).Path(route.Path).Name(route.Name).Handler(handler)
	}

	return router
}

func LogAccess(handler http.Handler, name string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//start := time.Now()
			handler.ServeHTTP(w, r)
			//log.Printf("%s \t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
		},
	)
}

func ProcessLog(action string, api string, message string, actionTime time.Time) {
	//log.Printf(" PROCESS LOG : %s \t%s\t%s\t%s", action, api, message, time.Since(actionTime))
}

func ErrorLog(w http.ResponseWriter, action string, api string, message string, err error, actionTime time.Time) {
	//log.Printf(" PROCESS LOG : \n\t ACTION : %s \n\t API : %s \n\t MESSAGE : %s \n\t TIME:%s \n\t ERROR :%s", action, api, message, time.Since(actionTime), err)
	w.WriteHeader(500)

}

const maxUploadSize = 32 * 200000 * 1024
const basePath = ".\\data\\"

func getFile(fileLocation string, file chan<- []byte, err chan<- error) {
	f, e := ioutil.ReadFile(fileLocation)
	ProcessLog("GetResource ", "Method", "File - size"+strconv.Itoa(len(f)), time.Now())
	err <- e
	file <- f
}

func GenericGet(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ProcessLog("GetResource ", r.RequestURI, "Started", start)

	w.Header().Set("Content-Type", "application/octet-stream; charset=UTF-8")

	file := make(chan []byte)

	variables := mux.Vars(r)
	sceneId := variables["uuidscene"]
	resourceId := variables["uuidresource"]
	library := variables["library"]

	fileLocation := filepath.Join(basePath, sceneId, library, resourceId)
	ProcessLog("FILE LOCATION", "", fileLocation, start)

	err := make(chan error)

	go getFile(fileLocation, file, err)

	if <-err != nil {
		ErrorLog(w, "UploadResource ", r.RequestURI, "GetResource ERROR", <-err, start)
	}

	w.WriteHeader(200)
	w.Write([]byte(<-file))

	ProcessLog("GetResource ", r.RequestURI, "Finished", start)
}

func UploadResource(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ProcessLog("UploadResource ", r.RequestURI, "START Upload", start)

	sceneId := r.FormValue("scene")
	resourceId := r.FormValue("resource")
	library := r.FormValue("library")

	scenePath := filepath.Join(basePath, sceneId, strings.ToLower(library))
	fileLocation := filepath.Join(basePath, sceneId, strings.ToLower(library), resourceId)

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		ErrorLog(w, "UploadResource ", r.RequestURI, "ParseMultipartForm ERROR", err, start)
	}

	if _, err := os.Stat(scenePath); os.IsNotExist(err) {
		err = os.MkdirAll(scenePath, 0755)
		if err != nil {
			panic(err)
		}
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		ErrorLog(w, "UploadResource ", r.RequestURI, "FormFile ERROR", err, start)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		ErrorLog(w, "UploadResource ", r.RequestURI, "ReadAll ERROR", err, start)
	}

	newFile, err := os.Create(fileLocation)
	if err != nil {
		ErrorLog(w, "UploadResource ", r.RequestURI, "ERROR : Could not be placed at: "+filepath.Join(basePath, sceneId), err, start)
	}
	defer newFile.Close()

	start2 := time.Now()

	ProcessLog("WRITE START ", r.RequestURI, "START EXECUTION", start2)
	if _, err := newFile.Write(fileBytes); err != nil {
		ErrorLog(w, "UploadResource ", r.RequestURI, "Write ERROR", err, start)
	}
	ProcessLog("WRITE END ", r.RequestURI, "END EXECUTION", start2)

	if err == nil {
		_, _ = w.Write([]byte("Uploaded successfully"))
	}

	ProcessLog("UploadResource ", r.RequestURI, "END Upload", start)

}

func Ping(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode("Pong")
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

