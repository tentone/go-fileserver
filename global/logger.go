package global

import (
	"github.com/google/logger"
	"os"
)

// Start the server logger, create the log file and start its services.
func StartLogger(path string) {
	var file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file " + path)
	}

	logger.Init("GoDonkey", true, true, file)
}
