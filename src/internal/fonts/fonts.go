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
	// Font base directory
	SystemFontPath = "/mnt/SDCARD/.system/res/"

	// OG Font (ID 0) - Actually uses font2.ttf
	OGFontName = "font2.ttf"
	OGFontPath = SystemFontPath + OGFontName
	OGFontBackupName = "font2.backup.ttf"

	// Next Font (ID 1) - Actually uses font1.ttf
	NextFontName = "font1.ttf"
	NextFontPath = SystemFontPath + NextFontName
	NextFontBackupName = "font1.backup.ttf"

	// Legacy font (kept for backward compatibility)
	LegacyFontName = "chillroundm.ttf"
	LegacyFontPath = SystemFontPath + LegacyFontName
	LegacyFontBackupName = "chillroundm.backup.ttf"
)

// FontType represents which system font slot to replace
type FontType int

const (
	OGFont FontType = iota
	NextFont
	LegacyFont
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

	// Find font files
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			// Check if it's a TTF or OTF font file
			if strings.HasSuffix(strings.ToLower(entry.Name()), ".ttf") ||
			   strings.HasSuffix(strings.ToLower(entry.Name()), ".otf") {
				fonts = append(fonts, entry.Name())
			}
		}
	}

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

	// Special case for Restore options
	if fontName == "Restore OG Font" {
		backupPath := filepath.Join(filepath.Dir(OGFontPath), OGFontBackupName)
		if _, err := os.Stat(backupPath); err == nil {
			return backupPath, nil
		}
		return "", fmt.Errorf("OG font backup not found")
	} else if fontName == "Restore Next Font" {
		backupPath := filepath.Join(filepath.Dir(NextFontPath), NextFontBackupName)
		if _, err := os.Stat(backupPath); err == nil {
			return backupPath, nil
		}
		return "", fmt.Errorf("Next font backup not found")
	} else if fontName == "Restore Legacy Font" {
		backupPath := filepath.Join(filepath.Dir(LegacyFontPath), LegacyFontBackupName)
		if _, err := os.Stat(backupPath); err == nil {
			return backupPath, nil
		}
		return "", fmt.Errorf("Legacy font backup not found")
	}

	// Return the path to the font
	fontPath := filepath.Join(cwd, "Fonts", fontName)
	logging.LogDebug("Font path: %s", fontPath)

	// Check if the font exists
	_, err = os.Stat(fontPath)
	if err != nil {
		logging.LogDebug("Font not found: %v", err)
		return "", fmt.Errorf("font not found: %w", err)
	}

	return fontPath, nil
}

// BackupExists checks if a backup of the specified font type exists
func BackupExists(fontType FontType) bool {
	var backupPath string

	switch fontType {
	case OGFont:
		backupPath = filepath.Join(filepath.Dir(OGFontPath), OGFontBackupName)
	case NextFont:
		backupPath = filepath.Join(filepath.Dir(NextFontPath), NextFontBackupName)
	case LegacyFont:
		backupPath = filepath.Join(filepath.Dir(LegacyFontPath), LegacyFontBackupName)
	}

	_, err := os.Stat(backupPath)
	exists := err == nil
	logging.LogDebug("Checking if backup exists for font type %d at %s: %v", fontType, backupPath, exists)
	return exists
}

// GetFontTypeFromName determines which font type to restore based on the font name
func GetFontTypeFromName(fontName string) (FontType, bool) {
	switch fontName {
	case "Restore OG Font":
		return OGFont, true
	case "Restore Next Font":
		return NextFont, true
	default:
		return OGFont, false // Not a restore operation
	}
}

// ApplyFont applies a font to the specified system font slot
func ApplyFont(fontName string, fontType FontType) error {
	logging.LogDebug("===== START ApplyFont =====")
	logging.LogDebug("Applying font: %s to slot type: %d", fontName, fontType)

	// Store the original font type to ensure consistency throughout the function
	originalFontType := fontType

	// Check if this is a restoration
	isRestore := false
	if fontType, isRestore = GetFontTypeFromName(fontName); isRestore {
		logging.LogDebug("This is a restore operation, calling RestoreFont")
		return RestoreFont(fontType)
	}

	// Get the path to the font
	fontPath, err := GetFontPath(fontName)
	if err != nil {
		logging.LogDebug("Error getting font path: %v", err)
		return fmt.Errorf("error getting font path: %w", err)
	}

	logging.LogDebug("Font file path: %s", fontPath)
	logging.LogDebug("Using font type: %d (0=OG, 1=Next)", originalFontType)

	// Determine target path based on font type
	var targetPath string
	var backupPath string
	var fontTypeName string

	// Set paths based on font type - explicitly log each branch
	switch originalFontType {
	case OGFont:
		targetPath = OGFontPath
		backupPath = filepath.Join(filepath.Dir(OGFontPath), OGFontBackupName)
		fontTypeName = "OG"
		logging.LogDebug("OG Font case: targetPath=%s, backupPath=%s", targetPath, backupPath)

		// Create backup of OG font if it doesn't exist
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			logging.LogDebug("Creating backup of OG font: %s -> %s", OGFontPath, backupPath)
			if err := CopyFile(OGFontPath, backupPath); err != nil {
				logging.LogDebug("Error creating OG font backup: %v", err)
				return fmt.Errorf("error creating backup: %w", err)
			}
			logging.LogDebug("OG font backup created successfully")
		} else {
			logging.LogDebug("OG font backup already exists at: %s", backupPath)
		}

	case NextFont:
		targetPath = NextFontPath
		backupPath = filepath.Join(filepath.Dir(NextFontPath), NextFontBackupName)
		fontTypeName = "Next"
		logging.LogDebug("Next Font case: targetPath=%s, backupPath=%s", targetPath, backupPath)

		// Create backup of Next font if it doesn't exist
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			logging.LogDebug("Creating backup of Next font: %s -> %s", NextFontPath, backupPath)
			if err := CopyFile(NextFontPath, backupPath); err != nil {
				logging.LogDebug("Error creating Next font backup: %v", err)
				return fmt.Errorf("error creating backup: %w", err)
			}
			logging.LogDebug("Next font backup created successfully")
		} else {
			logging.LogDebug("Next font backup already exists at: %s", backupPath)
		}

	default:
		logging.LogDebug("Unsupported font type: %d", originalFontType)
		return fmt.Errorf("unsupported font type: %d", originalFontType)
	}

	// Verify paths right before copying
	logging.LogDebug("FINAL CHECK - About to copy font:")
	logging.LogDebug("  From: %s", fontPath)
	logging.LogDebug("  To: %s", targetPath)
	logging.LogDebug("  Font Type: %s (%d)", fontTypeName, originalFontType)

	// Copy the font to the system - use explicit paths for maximum clarity
	var destPath string
	if originalFontType == OGFont {
		destPath = OGFontPath
	} else if originalFontType == NextFont {
		destPath = NextFontPath
	}

	logging.LogDebug("Using explicit destPath: %s", destPath)

	if err := CopyFile(fontPath, destPath); err != nil {
		logging.LogDebug("Error copying font to system: %v", err)
		return fmt.Errorf("error copying font to system: %w", err)
	}

	logging.LogDebug("Font '%s' applied successfully to %s font (type: %d)",
		fontName, fontTypeName, originalFontType)
	logging.LogDebug("===== END ApplyFont =====")
	return nil
}

// RestoreFont restores the original font for a particular font slot
func RestoreFont(fontType FontType) error {
	var backupPath, targetPath string
	var typeName string

	switch fontType {
	case OGFont:
		backupPath = filepath.Join(filepath.Dir(OGFontPath), OGFontBackupName)
		targetPath = OGFontPath
		typeName = "OG"
	case NextFont:
		backupPath = filepath.Join(filepath.Dir(NextFontPath), NextFontBackupName)
		targetPath = NextFontPath
		typeName = "Next"
	case LegacyFont:
		backupPath = filepath.Join(filepath.Dir(LegacyFontPath), LegacyFontBackupName)
		targetPath = LegacyFontPath
		typeName = "Legacy"
	}

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		logging.LogDebug("Backup file doesn't exist for %s font: %s", typeName, backupPath)
		return fmt.Errorf("backup file doesn't exist for %s font", typeName)
	}

	// Restore from backup
	if err := CopyFile(backupPath, targetPath); err != nil {
		logging.LogDebug("Error restoring %s font: %v", typeName, err)
		return fmt.Errorf("error restoring %s font: %w", typeName, err)
	}

	logging.LogDebug("Successfully restored %s font", typeName)
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