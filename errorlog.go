package main

import (
	"time"

	"gorm.io/gorm"
)

// Used to store in the database logs of the errors from API calls to the platform.
//
// Is created when the user receives a bad response from the api, also show the error ID to the user for debug later.
type ErrorLog struct {
	gorm.Model

	// Date and time of the error log
	Date time.Time `gorm:"column:date" json:"date"`

	// API route requested that caused the problem
	Route string `gorm:"column:route;type:text" json:"route"`

	// Error details that caused the problem
	Error string `gorm:"column:error;type:text" json:"error"`

	// Error code returned to the user
	Code int `gorm:"column:code" json:"code"`
}

func ErrorLogMigrate(db *gorm.DB) {
	var err = db.AutoMigrate(&ErrorLog{})
	if err != nil {
		print("Failed to migrate error log table.")
	}
}

func NewErrorLog(message string, err string, code int, route string) *ErrorLog {
	return &ErrorLog{
		Date:  time.Now(),
		Error: err,
		Code:  code,
		Route: route,
	}
}

// Create new log entry in the database, the ID of the object passed is populated.
func (log *ErrorLog) CreateDB(db *gorm.DB) error {
	var conn = db.Create(log)
	if conn.Error != nil {
		return conn.Error
	}

	return nil
}
