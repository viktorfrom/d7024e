package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// DefaultLogFilename is the default filename for the logger
const DefaultLogFilename string = "kademlia.log"

// Logger is a wrapper around the logrus logger
type Logger struct {
	logger *log.Logger
}

// New creates a new logger at the given `logLevel`. If `fileName` is
// given it will output logs to file. Enable `detailed` for method origin in output.
func New(logLevel log.Level, fileName *string, detailed bool) *Logger {
	baseLogger := log.New()

	baseLogger.SetFormatter(&log.TextFormatter{})
	baseLogger.SetLevel(logLevel)
	baseLogger.SetReportCaller(detailed)

	if fileName != nil {
		file, err := os.OpenFile(*fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			baseLogger.SetOutput(file)
		} else {
			baseLogger.Info("Failed to log to file, using default stderr")
		}
	}

	logger := &Logger{baseLogger}

	return logger
}

// Info log info level information
func (logger *Logger) Info(args ...interface{}) {
	logger.logger.Info(args...)
}

// Warning log warning level information
func (logger *Logger) Warn(args ...interface{}) {
	logger.logger.Warning(args...)
}

// Error log warning level information
func (logger *Logger) Error(args ...interface{}) {
	logger.logger.Error(args...)
}
