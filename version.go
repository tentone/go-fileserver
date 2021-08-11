package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Structure to represent the version of the api
type VersionStruct struct {
	Version   string `json:"version"`
	GitCommit string `json:"commit"`
}

// Version data of the api. Loaded from file updated when the project is built.
var Version VersionStruct = VersionStruct{}

// Read version information from file
func LoadVersion(path string) {
	var err error

	// Read data from file
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		print("Failed to read the version file.", path, err)
	}

	// Unmarshal json data
	var data []byte
	data, err = ioutil.ReadAll(file)
	err = json.Unmarshal(data, &Version)
	if err != nil {
		print("Failed to parse the version file.", path, err)
	}

	print("Loaded version file.")
}
