package database

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Used to store in the database logs of the errors from API calls to the platform.
//
// Is created when the user receives a bad response from the server, also show the error ID to the user for debug later.
type ErrorLog struct {
	NumID

	// Date and time of the error log
	Date time.Time `gorm:"column:date" json:"date"`

	// Data send by the user client to the API in the request body
	Received string `gorm:"column:received;type:text" json:"received"`

	// API route requested that caused the problem
	Route string `gorm:"column:route;type:text" json:"route"`

	// Error message (as shown to the user)
	Message string `gorm:"column:message;type:text" json:"message"`

	// Error details that caused the problem
	Error string `gorm:"column:error;type:text" json:"error"`

	// Error code returned to the user
	Code int `gorm:"column:code" json:"code"`
}

func ErrorLogMigrate(db *gorm.DB) {
	db.SingularTable(true)
	db.AutoMigrate(&ErrorLog{})
}

func NewErrorLog(received string, message string, err string, code int, route string) *ErrorLog {

	var log = new(ErrorLog)
	log.Date = time.Now()
	log.Received = received
	log.Message = message
	log.Error = err
	log.Code = code
	log.Route = route

	return log
}