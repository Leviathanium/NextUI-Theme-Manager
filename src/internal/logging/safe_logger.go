// src/internal/logging/safe_logger.go
// Safe logging functions that won't crash during initialization

package logging

import (
	"fmt"
	"sync/atomic"
)

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