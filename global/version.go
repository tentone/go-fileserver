package global

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Version data of the server. Loaded from file updated when the project is built.
var Version VersionStruct = VersionStruct{}

// Read version information from file
func LoadVersion(path string) error {
	var err error

	// Read data from file
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return err
	}

	// Unmarshal json data
	var data []byte
	data, err = ioutil.ReadAll(file)
	err = json.Unmarshal(data, &Version)
	if err != nil {
		return err
	}

	return nil
}

// Structure to represent the version of the server
type VersionStruct struct {
	Version string `json:"version"`
	GitCommit string `json:"commit"`
}
