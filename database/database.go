package database

import (
	"github.com/google/logger"
	_ "github.com/tentone/go-fileserver/global"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GORM database object, used to access data and perform operations.
var DB *gorm.DB

// Create database
func Create() {
	var err error = Connect()
	if err != nil {
		logger.Fatal("Error connecting to the database.", err.Error())
	}

	logger.Info("Migrating database structure.")

	// Create databases
	ErrorLogMigrate(DB)
	LibraryMigrate(DB)
	ResourceMigrate(DB)
}

// Connect to the SQL database using the configuration specified.
func Connect() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{}) // global.Config.Database.Dialect, global.Config.Database.ConnectionString)
	if err != nil {
		return err
	}

	return nil
}

// Close the database db, should be called before exiting the application.
func Close() {
	var err = DB.Close()
	if err != nil {
		logger.Fatal("Error closing connection to the SQL server.", err.Error())
		return
	}
}
