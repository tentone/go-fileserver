package api

import (
	"github.com/google/logger"
	"github.com/tentone/go-fileserver/database"
	"net/http"
)

// Log API error into the database and system log and set the error message as response of the request.
func LogAPIError(writer http.ResponseWriter, request *http.Request, message string, status int, err error) {
	var route string = string(request.URL.Path)

	// Error details
	var errorDetails string
	if err != nil {
		errorDetails = err.Error()
	}

	// Store log in the database
	var l = database.NewErrorLog(message, errorDetails, status, route)
	err = l.CreateDB(database.DB)
	if err != nil {
		logger.Error("Error trying to log error database." + " [" + err.Error() + "]")
	}

	// Write to the api log
	logger.Error(message + " [" + errorDetails + "]")

	writer.WriteHeader(status)
	_, err = writer.Write([]byte(message))
	if err != nil {
		logger.Error("Error trying to write back HTTP response error." + " [" + err.Error() + "]")
	}
}
