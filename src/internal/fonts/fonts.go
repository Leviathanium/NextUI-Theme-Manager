// src/internal/fonts/fonts.go
// Font operations for the NextUI Theme Manager

package fonts

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"nextui-themes/internal/logging"
)

const (
	// Font file name
	SystemFontName = "chillroundm.ttf"

	// System font path
	SystemFontPath = "/mnt/SDCARD/.system/res/chillroundm.ttf"

	// Backup font name
	BackupFontName = "chillroundm.backup.ttf"
)

// ListFonts returns a list of available fonts
func ListFonts() ([]string, error) {
	var fonts []string

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Get the fonts directory
	fontsDir := filepath.Join(cwd, "Fonts")
	logging.LogDebug("Scanning fonts directory: %s", fontsDir)

	// Check if the directory exists
	_, err = os.Stat(fontsDir)
	if os.IsNotExist(err) {
		logging.LogDebug("Fonts directory does not exist, creating: %s", fontsDir)
		// Create the fonts directory if it doesn't exist
		err = os.MkdirAll(fontsDir, 0755)
		if err != nil {
			logging.LogDebug("Error creating fonts directory: %v", err)
			return nil, fmt.Errorf("error creating fonts directory: %w", err)
		}
		return fonts, nil
	} else if err != nil {
		logging.LogDebug("Error checking fonts directory: %v", err)
		return nil, fmt.Errorf("error checking fonts directory: %w", err)
	}

	// Read the fonts directory
	entries, err := os.ReadDir(fontsDir)
	if err != nil {
		logging.LogDebug("Error reading fonts directory: %v", err)
		return nil, fmt.Errorf("error reading fonts directory: %w", err)
	}

	// Find directories that contain a font file
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			fontPath := filepath.Join(fontsDir, entry.Name(), SystemFontName)
			if _, err := os.Stat(fontPath); err == nil {
				fonts = append(fonts, entry.Name())
			}
		}
	}

	// Add "Default" option if there are other fonts available
	if len(fonts) > 0 {
		fonts = append(fonts, "Restore Default Font")
	}

	logging.LogDebug("Found %d fonts", len(fonts))
	return fonts, nil
}

// GetFontPath returns the path to a font
func GetFontPath(fontName string) (string, error) {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Special case for Default font
	if fontName == "Restore Default Font" {
		backupPath := filepath.Join(filepath.Dir(SystemFontPath), BackupFontName)
		if _, err := os.Stat(backupPath); err == nil {
			return backupPath, nil
		}
		return "", fmt.Errorf("backup font not found")
	}

	// Return the path to the font
	fontPath := filepath.Join(cwd, "Fonts", fontName, SystemFontName)
	logging.LogDebug("Font path: %s", fontPath)

	// Check if the font exists
	_, err = os.Stat(fontPath)
	if err != nil {
		logging.LogDebug("Font not found: %v", err)
		return "", fmt.Errorf("font not found: %w", err)
	}

	return fontPath, nil
}

// ApplyFont applies a font to the system
func ApplyFont(fontName string) error {
	logging.LogDebug("Applying font: %s", fontName)

	// Get the path to the font
	fontPath, err := GetFontPath(fontName)
	if err != nil {
		logging.LogDebug("Error getting font path: %v", err)
		return fmt.Errorf("error getting font path: %w", err)
	}

	// Create backup of the current font if it doesn't exist
	backupPath := filepath.Join(filepath.Dir(SystemFontPath), BackupFontName)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		logging.LogDebug("Creating backup of system font: %s -> %s", SystemFontPath, backupPath)
		if err := CopyFile(SystemFontPath, backupPath); err != nil {
			logging.LogDebug("Error creating backup: %v", err)
			return fmt.Errorf("error creating backup: %w", err)
		}
	}

	// Copy the font to the system
	if err := CopyFile(fontPath, SystemFontPath); err != nil {
		logging.LogDebug("Error copying font to system: %v", err)
		return fmt.Errorf("error copying font to system: %w", err)
	}

	logging.LogDebug("Font applied successfully: %s", fontName)
	return nil
}

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