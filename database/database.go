package database

import (
	"github.com/google/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tentone/godonkey/global"
)

// GORM database object, used to access data and perform operations.
var DB *gorm.DB

// Create database
func Create() {
	Connect()

	// Create databases
	ErrorLogMigrate(DB)
	LibraryMigrate(DB)
	ResourceMigrate(DB)
}

// Connect to the SQL database using the configuration specified.
func Connect() {
	var err error
	DB, err = gorm.Open(global.Config.Database.Dialect, global.Config.Database.ConnectionString)
	if err != nil {
		logger.Fatal("Error connecting to the SQL server.", err.Error())
		return
	}

	Create()
}

// Close the database db, should be called before exiting the application.
func Close() {
	var err = DB.Close()
	if err != nil {
		logger.Fatal("Error closing connection to the SQL server.", err.Error())
		return
	}
}
