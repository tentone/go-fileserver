package database

import (
	"database/sql"
	"github.com/google/logger"
	_ "github.com/mattn/go-sqlite3"
	"godonkey/global"
)

var database *sql.DB

// Load SQLite database file, and try to create the data structure if it is not present.
//
// The path for the database file is specified in the configuration file.
func Load() {

	var err error

	database, err = sql.Open("sqlite3", global.SqliteDatabase)
	if err != nil {
		logger.Error("Error opening the sqlite database file.")
		return
	}

	CreateStructure()
}

// Run a query against this database, the function checks internally for error in the SQL query.
//
// Print them into the application log for debugging.
func Query(query string) sql.Result {

	var statement, err = database.Prepare (query)
	if err != nil {
		logger.Error("Error in the SQL query. " + query + " " + err.Error())
		return nil
	}

	var result sql.Result

	result, err = statement.Exec()
	if err != nil {
		logger.Error("Error running the SQL query. " + query + " " + err.Error())
		return nil
	}

	return result
}

// Create the database structure, create tables for every type of resource.
//
// All tables use the same base layout.
func CreateStructure() {

	var tables = [4]string{global.IMAGES, global.FILES, global.GEOMETRIES, global.POTREE}

	for i := 0; i < len(tables); i++ {

		Query("CREATE TABLE IF NOT EXISTS '" + tables[i] + "' (uuid VARCHAR(36) UNIQUE PRIMARY KEY, format TEXT, filename TEXT);")
	}
}

// Add a resource entry to the database.
func InsertResource(library string, uuid string, format string, filename string) {

	logger.Info("Resource added " + library + ", uuid " + uuid + ", format " + format)
	Query("INSERT INTO " + library + " (uuid, format, filename) VALUES ('" + uuid + "', '" + format + "', '" + filename + "');")
}

// Remove a resource from the database.
func RemoveResource(library string, uuid string) {

	logger.Info("Resource deleted " + library + ", uuid " + uuid)
	Query("DELETE FROM " + library + " WHERE uuid = " + uuid + ";")
}
