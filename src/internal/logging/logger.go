// internal/logging/logger.go
package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

var logFile *os.File

// Flag to indicate whether logging is enabled
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

// LogSafe logs a message safely, using fmt.Printf if the logger isn't initialized yet
func LogSafe(format string, args ...interface{}) {
	if IsLoggerInitialized() {
		LogDebug(format, args...)
	} else {
		// Fall back to standard output if logger isn't ready
		fmt.Printf("EARLY LOG: "+format+"\n", args...)
	}
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

    // Create log file in the Logs directory
    logPath := filepath.Join(logsDir, "theme_manager.log")
    f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
        return
    }

    logFile = f

    // Log startup information
    LogDebug("=== Theme Manager Started ===")
    LogDebug("Current directory: %s", cwd)
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
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, message)

	logFile.WriteString(logLine)
	logFile.Sync()
}