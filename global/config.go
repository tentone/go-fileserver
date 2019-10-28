package global

const MSSQL string = "mssql"
const SQLITE string = "sqlite"
const MYSQL string = "mysql"
const POSTGRES string = "postgres"

// HTTP Server related configuration parameters
type ServerConfig struct {
	Address string `json:"address"`
	AddressTLS string `json:"addressTLS"`
	CertFileTLS string `json:"certFileTLS"`
	KeyFileTLS string `json:"keyFileTLS"`
}

// Database access configuration
type DatabaseConfig struct {
	Dialect string `json:"dialect"`
	Path string `json:"path;omitempty"` // Used for SQLITE only

}

// General configuration structure, containing all parameters.
type Config struct {
	Server ServerConfig `json:"server"`
	Database DatabaseConfig `json:"database"`

	DataPath string `json:"dataPath"`
	DevelopmentMode bool `json:"developmentMode"`
}
