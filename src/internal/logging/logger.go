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

// InitLogger initializes the debug log file
func init() {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		return
	}

	// Create log file
	logPath := filepath.Join(cwd, "theme_manager.log")
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
	if logFile != nil {
		LogDebug("=== Theme Manager Closed ===")
		logFile.Close()
	}
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	if logFile == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, message)

	logFile.WriteString(logLine)
	logFile.Sync()
}