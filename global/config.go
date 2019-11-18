package global

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const MSSQL string = "mssql"
const SQLITE string = "sqlite"
const MYSQL string = "mysql"
const POSTGRES string = "postgres"

const FILE string = "file"
const FTP string = "ftp"

// Global configuration of the server.
var Config ConfigStruct = ConfigStruct{}

// Read configuration from file
func LoadConfig(path string) (e error) {
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
	err = json.Unmarshal(data, &Config)
	if err != nil {
		return err
	}

	return nil
}

// General configuration structure, containing all parameters.
type ConfigStruct struct {
	Server ServerConfig `json:"server"`
	Database DatabaseConfig `json:"database"`
	FileServer FileServerConfig `json:"fileServer"`
	Storage StorageConfig `json:"storage"`
	DevelopmentMode bool `json:"developmentMode"`
}

// HTTP Server related configuration parameters
type ServerConfig struct {
	Address string `json:"address,omitempty"`
	AddressTLS string `json:"addressTLS,omitempty"`
	CertFileTLS string `json:"certFileTLS,omitempty"`
	KeyFileTLS string `json:"keyFileTLS,omitempty"`
	GenerateCertTLS bool `json:"generateCertTLS,omitempty"`
}

// Database access configuration
type DatabaseConfig struct {
	Dialect string `json:"dialect"`
	Path string `json:"path,omitempty"` // Used for SQLITE only
	ConnectionString string `json:"connectionString,omitempty"`
}

// File server specific configuration
type FileServerConfig struct {
	MaxUploadSize int `json:"maxUploadSize"`
}

// Storage configuration
type StorageConfig struct {
	Mode string `json:"mode"`
	Path string `json:"path"`
}
