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

// Connect to the SQL database using the configuration specified.
func ConnectDatabase() error {
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
		DisableAutomaticPing:                     false,
		FullSaveAssociations:                     true,
		NamingStrategy:                           schema.NamingStrategy{TablePrefix: "", SingularTable: true},
		Logger:                                   logg,
	}

	var err error
	db, err = gorm.Open(sqlite.Open("database.db"), cfg)
	if err != nil {
		return err
	}

	return nil
}

// Close the database db, should be called before exiting the application.
func Close() error {
	sql, err := db.DB()
	if err != nil {
		print("Error getting SQL db.")
		return err
	}

	err = sql.Close()
	if err != nil {
		print("Error closing connection to the SQL server.")
		return err
	}

	return nil
}
