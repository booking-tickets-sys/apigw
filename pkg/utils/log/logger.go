package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// InitLogger initializes the logger
func InitLogger() error {
	logger = logrus.New()

	// Set log level
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logger.SetLevel(logLevel)

	// Set log format
	format := os.Getenv("LOG_FORMAT")
	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set output
	logger.SetOutput(os.Stdout)

	return nil
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	if logger == nil {
		// Initialize with defaults if not already initialized
		InitLogger()
	}
	return logger
}
