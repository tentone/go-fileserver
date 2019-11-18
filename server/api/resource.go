package api

import (
	"github.com/gorilla/mux"
	"github.com/tentone/godonkey/global"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func ResourceGet(writer http.ResponseWriter, request *http.Request) {
	// Form data
	var variables = mux.Vars(request)
	var library string = variables["library"]
	var uuid string = variables["uuid"]

	var path string = global.Config.Storage.Path + "/" + strings.ToLower(library) + "/" + uuid

	// Read file
	var err error
	var file []byte
	file, err = ioutil.ReadFile(path)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write file
	_, err = writer.Write(file)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/octet-stream; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}

func ResourceUpload(writer http.ResponseWriter, request *http.Request) {
	var uuid = request.FormValue("uuid")
	var library = request.FormValue("library")
	var format = request.FormValue("format")

	var path string = global.Config.Storage.Path + "/" + strings.ToLower(library)

	// Check if path exists and create if necessary
	var err error
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			_, _ = writer.Write([]byte("Failed to create directory to store data."))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// File path
	var fpath string =   path + "/" + strings.ToLower(uuid) + "." + format

	// Read request data
	request.Body = http.MaxBytesReader(writer, request.Body, global.Config.FileServer.MaxUploadSize)
	err = request.ParseMultipartForm(global.Config.FileServer.MaxUploadSize)
	if err != nil {
		_, _ = writer.Write([]byte("Cannot read data from the request form."))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var file multipart.File
	file, _, err = request.FormFile("file")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var data []byte
	data, err = ioutil.ReadAll(file)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create and store file
	var storeFile *os.File
	storeFile, err = os.Create(fpath)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer storeFile.Close()

	_, err = storeFile.Write(data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}