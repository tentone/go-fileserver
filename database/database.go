package database

import (
	"github.com/google/logger"
	"github.com/jinzhu/gorm"
	"github.com/tentone/godonkey/global"
)

// GORM database object
var DB *gorm.DB

// Get the GORM database object, used to access data and perform operations.
func Get() *gorm.DB {
	return DB
}

//
func Create() {

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
