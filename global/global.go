package global

import (
	"encoding/json"
	"github.com/google/logger"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
)

const GEOMETRIES string = "geometries"
const IMAGES string = "images"
const POTREE string = "potree"
const FILES string = "files"

// Indicates if the system is running in development mode.
//
// The may be additional development logs and processing done in development mode for debug.
//
// If true a service to check the system log is available for easier access to data.
var DevelopmentMode bool

// Version of the resource server itself.
//
// Follow the semantic versioning spec.
var Version string

// Version of the resource server API, is only changed on major API changes.
// May be used for service version control.
var ApiVersion string

// GitCommit identifier
var GitCommit string

// Name of the configuration file used to read values.
var ConfigurationFile string

// Path to the server log file.
var LogFile string = "server.log"

// Address for the server to be listed at.
var Address string

var AddressTLS string
var CertFileTLS string
var KeyFileTLS string

// Base path of the data stored in hard disk.
//
// File should always use their identifiers (UUID) as name for fast access (let the file system index everything).
var DataPath string

// URL of the main backend API server.
var ApiServer string

// File of the SQLite local database, used to store resource metadata.
var SqliteDatabase string

// Load the version information from the version file.
func LoadVersion() {

	var file, err = os.Open("version.json")
	if err != nil {
		logger.Fatal("Failed to open version file.")
		return
	}

	var data []byte
	data, err = ioutil.ReadAll(file)
	if err != nil {
		logger.Fatal("Failed to read version file.")
		return
	}

	var result = struct {

		ApiVersion string `json:"api"`
		Version string `json:"version"`
		GitCommit string `json:"commit"`
	}{}

	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		logger.Fatal("Failed to parse version file JSON.")
		return
	}

	ApiVersion = result.ApiVersion
	Version = result.Version
	GitCommit = result.GitCommit

	err = file.Close()
	if err != nil {
		logger.Info("Error closing the version file.")
	}
}

// Read value from the configuration file named ConfigurationFile
// Expects the file to be well formatted with address, dataPath and logFile values.
func LoadConfig() {

	var file, err = os.Open(ConfigurationFile)
	if err != nil {
		logger.Fatal("Failed to open configuration file. " + ConfigurationFile)
		return
	}

	var data, _ = ioutil.ReadAll(file)

	var result = struct {

		Address string `json:"address"`
		AddressTLS string `json:"addressTLS"`
		CertFileTLS string `json:"certFileTLS"`
		KeyFileTLS string `json:"keyFileTLS"`
		DataPath string `json:"dataPath"`
		SqliteDatabase string `json:"sqliteDatabase"`
		DevelopmentMode bool `json:"developmentMode"`
	}{}

	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		logger.Fatal("Failed to parse configuration file JSON. " + ConfigurationFile)
		return
	}

	Address = result.Address

	AddressTLS = result.AddressTLS
	CertFileTLS = result.CertFileTLS
	KeyFileTLS = result.KeyFileTLS

	DataPath = result.DataPath
	ApiServer = result.ApiServer
	SqliteDatabase = result.SqliteDatabase
	DevelopmentMode = result.DevelopmentMode

	logger.Info("Configuration loaded from " + ConfigurationFile + "(Address:" + Address + ", DataPath:" + DataPath + ", APIServer:" + ApiServer + ")")

	err = file.Close()
	if err != nil {
		logger.Info("Error closing the configuration file.")
	}
}
