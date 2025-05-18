// internal/app/logger.go
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

// Logger file handle
var logFile *os.File

// Flag to indicate if logging is enabled
var LoggingEnabled bool = true

// Flag to indicate whether logging is fully initialized
var loggerInitialized int32 = 0

// SetLoggerInitialized marks the logger as fully initialized
func SetLoggerInitialized() {
	atomic.StoreInt32(&loggerInitialized, 1)
}

// IsLoggerInitialized checks if the logger is initialized
func IsLoggerInitialized() bool {
	return atomic.LoadInt32(&loggerInitialized) == 1
}

// Initialize the logger
func init() {
	// If logging is disabled, don't initialize anything
	if !LoggingEnabled {
		return
	}

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		return
	}

	// Create Logs directory if it doesn't exist
	logsDir := filepath.Join(cwd, "Logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Logs directory: %v\n", err)
		return
	}

	// Create log file with timestamp in name
	timestamp := time.Now().Format("20060102_150405")
	logPath := filepath.Join(logsDir, fmt.Sprintf("theme_manager_%s.log", timestamp))
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		return
	}

	logFile = f

	// Log startup information
	LogDebug("=== Theme Manager Started ===")
	LogDebug("Log file: %s", logPath)
}

// CloseLogger closes the log file
func CloseLogger() {
	if !LoggingEnabled || logFile == nil {
		return
	}

	LogDebug("=== Theme Manager Closed ===")
	logFile.Close()
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	// If logging is disabled or log file isn't initialized, return immediately
	if !LoggingEnabled || logFile == nil {
		// If we're very early in initialization, write to stderr
		if len(args) > 0 {
			fmt.Fprintf(os.Stderr, format+"\n", args...)
		} else {
			fmt.Fprintln(os.Stderr, format)
		}
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, message)

	// Write to log file
	_, err := logFile.WriteString(logLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
	}

	// Flush to disk
	logFile.Sync()
}