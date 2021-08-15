package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// GORM database object, used to access data and perform operations.
var db *gorm.DB

// Get access to the GORM database object. Used to access data and perform operations.
func Get() *gorm.DB {
	return db
}

// Create database
func Create() {
	var err error = Connect()
	if err != nil {
		print("Error connecting to the database.", err.Error())
	}

	print("Migrating database structure.")

	Initialize()
}

// Create database tables
func Initialize() {
	LogMigrate(db)
	LibraryMigrate(db)
	ResourceMigrate(db)
}

// Connect to the SQL database using the configuration specified.
func Connect() error {
	var logg = logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 10 * time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)

	var cfg = &gorm.Config{
		CreateBatchSize:                          1000,
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableAutomaticPing:                     true,
		FullSaveAssociations:                     true,
		NamingStrategy:                           schema.NamingStrategy{TablePrefix: "", SingularTable: true},
		Logger:                                   logg,
	}

	var err error
	db, err = gorm.Open(sqlite.Open("database.db"), cfg)
	if err != nil {
		panic("failed to connect database")
	}

	return nil
}

// Close the database db, should be called before exiting the application.
func Close() {
	var sql, err = db.DB()
	if err != nil {
		print("Error getting SQL db.")
		return
	}

	err = sql.Close()
	if err != nil {
		print("Error closing connection to the SQL server.")
		return
	}
}
