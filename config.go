package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Global configuration of the api.
var Config ConfigStruct = ConfigStruct{}

// Read configuration from file
func LoadConfig(path string) {
	var err error

	// Read data from file
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		print("Failed to read the configuration file.", path, err)
	}

	// Unmarshal json data
	var data []byte
	data, err = ioutil.ReadAll(file)
	err = json.Unmarshal(data, &Config)
	if err != nil {
		print("Failed to parse the configuration file.", path, err)
	}

	print("Loaded configuration file.")
}

// General configuration structure, containing all parameters.
type ConfigStruct struct {
	Server      ServerConfig     `json:"api"`
	Database    DatabaseConfig   `json:"database"`
	FileServer  FileServerConfig `json:"fileServer"`
	Storage     StorageConfig    `json:"storage"`
	Development bool             `json:"development"`
}

// HTTP Server related configuration parameters
type ServerConfig struct {
	Address         string `json:"address"`
	AddressTLS      string `json:"addressTLS,omitempty"`
	CertFileTLS     string `json:"certFileTLS,omitempty"`
	KeyFileTLS      string `json:"keyFileTLS,omitempty"`
	GenerateCertTLS bool   `json:"generateCertTLS,omitempty"`
}

// Database access configuration
type DatabaseConfig struct {
	Dialect          string `json:"dialect"`
	ConnectionString string `json:"connectionString,omitempty"`
}

// File api specific configuration
type FileServerConfig struct {
	MaxUploadSize int64 `json:"maxUploadSize"`
}

// Storage configuration
type StorageConfig struct {
	Mode string `json:"mode"`
	Path string `json:"path"`
}
