// src/internal/themes/common.go
// Common utilities for theme operations

package themes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"nextui-themes/internal/logging"
)

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	logging.LogDebug("Copying %s to %s", src, dst)

	// Create the destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		logging.LogDebug("Error creating directory %s: %v", dstDir, err)
		return fmt.Errorf("failed to create directory %s: %w", dstDir, err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		logging.LogDebug("Error opening source file: %v", err)
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		logging.LogDebug("Error creating destination file: %v", err)
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		logging.LogDebug("Error copying file: %v", err)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	logging.LogDebug("Successfully copied %d bytes", bytes)
	return nil
}