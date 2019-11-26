package global

import (
	"encoding/json"
	"github.com/google/logger"
	"io/ioutil"
	"os"
)

// Version data of the server. Loaded from file updated when the project is built.
var Version VersionStruct = VersionStruct{}

// Read version information from file
func LoadVersion(path string) {
	var err error

	// Read data from file
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		logger.Fatal("Failed to read the version file.", path, err)
	}

	// Unmarshal json data
	var data []byte
	data, err = ioutil.ReadAll(file)
	err = json.Unmarshal(data, &Version)
	if err != nil {
		logger.Fatal("Failed to parse the version file.", path, err)
	}

	logger.Info("Loaded version file.")
}

// Structure to represent the version of the server
type VersionStruct struct {
	Version string `json:"version"`
	GitCommit string `json:"commit"`
}
