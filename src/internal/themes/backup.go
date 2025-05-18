// internal/themes/backup.go
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"thememanager/internal/app"
)

// ExportTheme exports the current system theme to a new theme package
func ExportTheme(themeName string) error {
	app.LogDebug("Exporting system theme to: %s", themeName)

	// Check if system theme directory exists
	if !SystemThemeExists() {
		return fmt.Errorf("system theme directory does not exist: %s", SystemThemeDir)
	}

	// Generate theme name with timestamp if not provided
	if themeName == "" {
		timestamp := time.Now().Format("20060102_150405")
		themeName = fmt.Sprintf("export_%s", timestamp)
	}

	// Make sure it has .theme extension
	if !strings.HasSuffix(themeName, ThemeExtension) {
		themeName += ThemeExtension
	}

	// Get export path
	exportPath := filepath.Join(app.GetWorkingDir(), "Themes", themeName)

	// Check if theme already exists
	if _, err := os.Stat(exportPath); err == nil {
		// Already exists, ask user for confirmation to overwrite
		// For now, just return an error
		return fmt.Errorf("theme already exists: %s", themeName)
	}

	// Create theme directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("failed to create theme directory: %w", err)
	}

	// Create Theme subdirectory in export
	themeExportPath := filepath.Join(exportPath, "Theme")
	if err := os.MkdirAll(themeExportPath, 0755); err != nil {
		return fmt.Errorf("failed to create Theme subdirectory in export: %w", err)
	}

	// Copy all files from system theme to export
	if err := CopyDirectory(SystemThemeDir, themeExportPath); err != nil {
		return fmt.Errorf("failed to copy system theme to export: %w", err)
	}

	// Create manifest for export
	manifest := CreateDefaultManifest(
		strings.TrimSuffix(themeName, ThemeExtension),
		"Theme Manager",
	)
	manifest.Description = "Exported system theme"

	// Write manifest to export
	if err := WriteManifest(manifest, exportPath); err != nil {
		return fmt.Errorf("failed to write export manifest: %w", err)
	}

	// Copy current theme preview as export preview if available
	// First try to find a preview.png in the theme directory
	systemPreviewPath := filepath.Join(SystemThemeDir, ThemePreviewFile)
	if _, err := os.Stat(systemPreviewPath); err == nil {
		// Copy preview to export
		exportPreviewPath := filepath.Join(exportPath, ThemePreviewFile)
		if err := CopyFile(systemPreviewPath, exportPreviewPath); err != nil {
			app.LogDebug("Warning: Failed to copy theme preview to export: %v", err)
			// Continue anyway, preview is not critical
		}
	} else {
		// If no preview in theme directory, create a placeholder
		// In a real implementation, we might take a screenshot or generate a preview
		app.LogDebug("No theme preview found, placeholder would be created here")
	}

	app.LogDebug("Theme exported successfully: %s", themeName)
	return nil
}

// CleanupOldBackups removes old backups to prevent using too much space
func CleanupOldBackups(maxBackups int) error {
	app.LogDebug("Cleaning up old backups, keeping max %d", maxBackups)

	// Get list of backups
	backups, err := ListBackups()
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	// If we have fewer backups than the max, do nothing
	if len(backups) <= maxBackups {
		app.LogDebug("No cleanup needed, only have %d backups", len(backups))
		return nil
	}

	// Get file info for each backup to sort by date
	type backupInfo struct {
		name    string
		modTime time.Time
	}

	var backupsList []backupInfo

	for _, backupName := range backups {
		backupPath := GetBackupPath(backupName)

		// Get file info
		info, err := os.Stat(backupPath)
		if err != nil {
			app.LogDebug("Warning: Failed to get info for backup %s: %v", backupName, err)
			continue
		}

		backupsList = append(backupsList, backupInfo{
			name:    backupName,
			modTime: info.ModTime(),
		})
	}

	// Sort backups by modification time (oldest first)
	sort.Slice(backupsList, func(i, j int) bool {
		return backupsList[i].modTime.Before(backupsList[j].modTime)
	})

	// Remove oldest backups until we're under the limit
	numToRemove := len(backupsList) - maxBackups

	for i := 0; i < numToRemove; i++ {
		backupToRemove := backupsList[i].name
		backupPath := GetBackupPath(backupToRemove)

		app.LogDebug("Removing old backup: %s", backupToRemove)

		if err := RemoveDirectory(backupPath); err != nil {
			app.LogDebug("Warning: Failed to remove backup %s: %v", backupToRemove, err)
			// Continue with other backups
		}
	}

	app.LogDebug("Cleanup complete, removed %d old backups", numToRemove)
	return nil
}

// GetThemeList returns a list of themes with their information
func GetThemeList() ([]map[string]string, error) {
	app.LogDebug("Getting theme list with information")

	themeNames, err := ListThemes()
	if err != nil {
		return nil, fmt.Errorf("failed to list themes: %w", err)
	}

	var themeList []map[string]string

	for _, themeName := range themeNames {
		themePath := GetThemePath(themeName)

		// Try to read manifest
		manifest, err := ReadManifest(themePath)

		// Create theme info
		themeInfo := map[string]string{
			"name": themeName,
		}

		// Add manifest info if available
		if err == nil {
			themeInfo["author"] = manifest.Author
			themeInfo["version"] = manifest.Version
			themeInfo["description"] = manifest.Description
		} else {
			app.LogDebug("Warning: Failed to read manifest for theme %s: %v", themeName, err)
		}

		// Check if preview exists
		previewPath := filepath.Join(themePath, ThemePreviewFile)
		if _, err := os.Stat(previewPath); err == nil {
			themeInfo["preview"] = previewPath
		}

		themeList = append(themeList, themeInfo)
	}

	return themeList, nil
}

// GetBackupList returns a list of backups with their information
func GetBackupList() ([]map[string]string, error) {
	app.LogDebug("Getting backup list with information")

	backupNames, err := ListBackups()
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	var backupList []map[string]string

	for _, backupName := range backupNames {
		backupPath := GetBackupPath(backupName)

		// Try to read manifest
		manifest, err := ReadManifest(backupPath)

		// Create backup info
		backupInfo := map[string]string{
			"name": backupName,
		}

		// Add manifest info if available
		if err == nil {
			backupInfo["author"] = manifest.Author
			backupInfo["version"] = manifest.Version
			backupInfo["description"] = manifest.Description
		} else {
			app.LogDebug("Warning: Failed to read manifest for backup %s: %v", backupName, err)
		}

		// Check if preview exists
		previewPath := filepath.Join(backupPath, ThemePreviewFile)
		if _, err := os.Stat(previewPath); err == nil {
			backupInfo["preview"] = previewPath
		}

		backupList = append(backupList, backupInfo)
	}

	return backupList, nil
}