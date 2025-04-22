// src/internal/logging/logger.go
// Logging functionality for the application

package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File

// Enables or disables logging functionality
var LoggingEnabled bool = true  // Set to false to disable all logging

// InitLogger initializes the debug log file
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

    // Mark logger as initialized
    SetLoggerInitialized()
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