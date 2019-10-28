package global

type Config struct {
	Address string `json:"address"`
	AddressTLS string `json:"addressTLS"`
	CertFileTLS string `json:"certFileTLS"`
	KeyFileTLS string `json:"keyFileTLS"`
	DataPath string `json:"dataPath"`
	SqliteDatabase string `json:"sqliteDatabase"`
	DevelopmentMode bool `json:"developmentMode"`
}

