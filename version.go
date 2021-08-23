package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

// Structure to represent the version of the api
type VersionStruct struct {
	Version string `json:"version"`
}

// Version data of the api. Loaded from file updated when the project is built.
var Version VersionStruct = VersionStruct{}

// Read version information from file
func LoadVersion(path string) {
	// Read data from file
	var file *os.File
	var err error
	file, err = os.Open(path)
	if err != nil {
		print("Failed to open the version file.", path, err)
	}

	// Unmarshal json data
	var data []byte
	data, err = ioutil.ReadAll(file)
	if err != nil {
		print("Failed to read the version file.", path, err)
	}

	err = json.Unmarshal(data, &Version)
	if err != nil {
		print("Failed to parse the version file.", path, err)
	}

	print("Loaded version file.")
}

// Get the server instance version information, provides a timestamp of the build.
func VersionGetAPI(writer http.ResponseWriter, request *http.Request) {
	var data []byte
	var err error
	data, err = json.Marshal(&Version)
	if err != nil {
		print(writer, request, "Failed to marshal JSON response.", http.StatusBadRequest, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}
