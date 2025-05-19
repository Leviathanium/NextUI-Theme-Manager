
// File: src/internal/themes/apply.go
// This is a complete replacement for the file

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"thememanager/internal/app"
)

// ApplyTheme copies a theme from the Themes directory to the system
func ApplyTheme(themeName string) error {
	app.LogDebug("Applying theme: %s", themeName)

	// Get paths
	themePath := GetThemePath(themeName)

	// Check if theme exists
	if !ThemeExists(themeName) {
		return fmt.Errorf("theme does not exist: %s", themeName)
	}

	// Read manifest to verify it's a valid theme - using strict validation
	manifest, err := ReadManifest(themePath, true)
	if err != nil {
		return fmt.Errorf("invalid theme package: %w", err)
	}

	app.LogDebug("Applying theme: %s by %s", manifest.Name, manifest.Author)

	// Check if system theme directory exists
	if !SystemThemeExists() {
		return fmt.Errorf("system theme directory does not exist: %s", SystemThemeDir)
	}

	// First, remove existing theme directory
	if err := RemoveDirectory(SystemThemeDir); err != nil {
		return fmt.Errorf("failed to remove existing theme directory: %w", err)
	}

	// Create theme directory
	if err := os.MkdirAll(SystemThemeDir, 0755); err != nil {
		return fmt.Errorf("failed to create theme directory: %w", err)
	}

	// Copy theme files
	// The actual theme content (not the manifest or preview) is in the "Theme" subdirectory
	themeContentPath := filepath.Join(themePath, "Theme")

	// Check if Theme subdirectory exists
	if _, err := os.Stat(themeContentPath); os.IsNotExist(err) {
		// If Theme subdirectory doesn't exist, assume the full theme package is the content
		themeContentPath = themePath
	}

	// Copy all files
	if err := CopyDirectory(themeContentPath, SystemThemeDir); err != nil {
		return fmt.Errorf("failed to copy theme files: %w", err)
	}

	app.LogDebug("Theme applied successfully: %s", themeName)
	return nil
}

// CreateBackup creates a backup of the current system theme
func CreateBackup(backupName string) error {
	app.LogDebug("Creating theme backup: %s", backupName)

	// Check if system theme directory exists
	if !SystemThemeExists() {
		return fmt.Errorf("system theme directory does not exist: %s", SystemThemeDir)
	}

	// Generate backup path with timestamp if not provided
	if backupName == "" {
		timestamp := time.Now().Format("20060102_150405")
		backupName = fmt.Sprintf(ThemeBackupNameFmt, timestamp)
	}

	// Make sure it has .theme extension
	if !strings.HasSuffix(backupName, ThemeExtension) {
		backupName += ThemeExtension
	}

	// Get backup path
	backupPath := filepath.Join(app.GetWorkingDir(), "Backups", backupName)

	// Ensure backup directory exists
	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Remove existing backup if it exists
	if err := RemoveDirectory(backupPath); err != nil {
		return fmt.Errorf("failed to remove existing backup: %w", err)
	}

	// Create backup directory
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create Theme subdirectory in backup
	themeBackupPath := filepath.Join(backupPath, "Theme")
	if err := os.MkdirAll(themeBackupPath, 0755); err != nil {
		return fmt.Errorf("failed to create Theme subdirectory in backup: %w", err)
	}

	// Copy all files from system theme to backup
	if err := CopyDirectory(SystemThemeDir, themeBackupPath); err != nil {
		return fmt.Errorf("failed to copy system theme to backup: %w", err)
	}

	// Create manifest for backup
	manifest := CreateDefaultManifest(
		strings.TrimSuffix(backupName, ThemeExtension),
		"Theme Manager",
	)
	manifest.Description = "Manual backup of system theme"

	// Write manifest to backup
	if err := WriteManifest(manifest, backupPath); err != nil {
		return fmt.Errorf("failed to write backup manifest: %w", err)
	}

	// Copy current theme preview as backup preview if available
	// First try to find a preview.png in the theme directory
	systemPreviewPath := filepath.Join(SystemThemeDir, ThemePreviewFile)
	if _, err := os.Stat(systemPreviewPath); err == nil {
		// Copy preview to backup
		backupPreviewPath := filepath.Join(backupPath, ThemePreviewFile)
		if err := CopyFile(systemPreviewPath, backupPreviewPath); err != nil {
			app.LogDebug("Warning: Failed to copy theme preview to backup: %v", err)
			// Continue anyway, preview is not critical
		}
	} else {
		// If no preview in theme directory, create a placeholder
		// In a real implementation, we might take a screenshot or generate a preview
		app.LogDebug("No theme preview found, placeholder would be created here")
	}

	app.LogDebug("Backup created successfully: %s", backupName)
	return nil
}

// RestoreBackup restores a theme from a backup in the Backups directory
func RestoreBackup(backupName string) error {
	app.LogDebug("Restoring theme from backup: %s", backupName)

	// Get backup path
	backupPath := GetBackupPath(backupName)

	// Check if backup exists
	if !BackupExists(backupName) {
		return fmt.Errorf("backup does not exist: %s", backupName)
	}

	// Read manifest to verify it's a valid backup (non-strict for backups)
	manifest, err := ReadManifest(backupPath, false)
	if err != nil {
		return fmt.Errorf("invalid backup package: %w", err)
	}

	app.LogDebug("Restoring theme from backup: %s by %s", manifest.Name, manifest.Author)

	// Check if system theme directory exists
	if !SystemThemeExists() {
		// If system theme directory doesn't exist, create it
		if err := os.MkdirAll(SystemThemeDir, 0755); err != nil {
			return fmt.Errorf("failed to create system theme directory: %w", err)
		}
	}

	// First, remove existing theme directory
	if err := RemoveDirectory(SystemThemeDir); err != nil {
		return fmt.Errorf("failed to remove existing theme directory: %w", err)
	}

	// Create system theme directory
	if err := os.MkdirAll(SystemThemeDir, 0755); err != nil {
		return fmt.Errorf("failed to create system theme directory: %w", err)
	}

	// Copy directly from the backup path to the system theme directory
	// without looking for a Theme subdirectory
	if err := CopyDirectory(backupPath, SystemThemeDir); err != nil {
		return fmt.Errorf("failed to copy theme files from backup: %w", err)
	}

	app.LogDebug("Theme restored successfully from backup: %s", backupName)
	return nil
}