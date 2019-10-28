package database

import (
	"fmt"
	"github.com/google/logger"
	"github.com/jinzhu/gorm"
)

// GORM database object
var db *gorm.DB

// Get access to the GORM database object. Used to access data and perform operations.
func Get() *gorm.DB {
	return db
}

func Create() {

}

// Connect to the SQL database using the configuration specified.
func Connect() {
	var err error
	var connection string = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;%s", global.SqlServer, global.SqlUser, global.SqlPassword, global.SqlServerPort, global.SqlDatabase, global.SqlParams)

	logger.Info("Connecting to the database " + connection)

	db, err = gorm.Open("mssql", connection)
	if err != nil {
		logger.Error("Error connecting to the SQL server." + err.Error())
		return
	}

	Create()
}

// Close the database db, should be called before exiting the application.
func Close() {
	var err = db.Close()
	if err != nil {
		logger.Error("Error closing connection to the SQL server.")
		return
	}

	logger.Info("Closed database connection.")
}
