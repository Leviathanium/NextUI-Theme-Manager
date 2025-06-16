// File: src/internal/themes/backup.go
// Complete replacement with updated ReadManifest calls

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
    "strconv"
	"thememanager/internal/app"
)

// ExportTheme exports the current system theme to a new theme package
// with sequential numbering (backup1.theme, backup2.theme, etc.)
func ExportTheme(themeName string) error {
	app.LogDebug("Exporting system theme as backup")

	// Check if system theme directory exists
	if !SystemThemeExists() {
		return fmt.Errorf("system theme directory does not exist: %s", SystemThemeDir)
	}

	// Generate sequential backup name if not provided
	if themeName == "" {
		// Get list of existing backups to determine next number
		backups, err := ListBackups()
		if err != nil {
			return fmt.Errorf("failed to list existing backups: %w", err)
		}

		// Find the highest number used so far
		highestNum := 0
		for _, backup := range backups {
			// Check if the backup name matches our pattern (backup1, backup2, etc.)
			if strings.HasPrefix(backup, "backup") {
				numStr := strings.TrimPrefix(backup, "backup")
				if num, err := strconv.Atoi(numStr); err == nil && num > highestNum {
					highestNum = num
				}
			}
		}

		// Next number is highest + 1
		themeName = fmt.Sprintf("backup%d", highestNum+1)
	}

	// Make sure it has .theme extension
	if !strings.HasSuffix(themeName, ThemeExtension) {
		themeName += ThemeExtension
	}

	// Use Backups directory instead of Themes directory
	exportPath := filepath.Join(app.GetWorkingDir(), "Backups", themeName)

	// Check if backup already exists
	if _, err := os.Stat(exportPath); err == nil {
		// Already exists, ask user for confirmation to overwrite
		// For now, just return an error
		return fmt.Errorf("backup already exists: %s", themeName)
	}

	// Create backup directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Copy all files from system theme to export
	if err := CopyDirectory(SystemThemeDir, exportPath); err != nil {
		return fmt.Errorf("failed to copy system theme to backup: %w", err)
	}

	// Copy tool icons to backup (handle tool icons specially)
	systemToolsMediaPath := "/mnt/SDCARD/Tools/tg5040/.media"
	if _, err := os.Stat(systemToolsMediaPath); err == nil {
		app.LogDebug("Found system tool icons, backing up to theme package")

		// Create Tools/tg5040/.media structure in backup
		backupToolsMediaPath := filepath.Join(exportPath, "Tools", "tg5040", ".media")
		if err := os.MkdirAll(backupToolsMediaPath, 0755); err != nil {
			return fmt.Errorf("failed to create backup tools media directory: %w", err)
		}

		// Copy tool icons
		if err := CopyDirectory(systemToolsMediaPath, backupToolsMediaPath); err != nil {
			app.LogDebug("Warning: Failed to copy tool icons to backup: %v", err)
			// Continue anyway, tool icons are not critical for basic functionality
		} else {
			app.LogDebug("Successfully backed up tool icons")
		}
	} else {
		app.LogDebug("No tool icons found in system, skipping tool icon backup")
	}

	// Copy Collections and Recently Played icons to backup (NEW: handle special icons)
	systemMediaPath := "/mnt/SDCARD/.media"

	// Handle Collections icon
	systemCollectionsPath := filepath.Join(systemMediaPath, "Collections.png")
	if _, err := os.Stat(systemCollectionsPath); err == nil {
		app.LogDebug("Found Collections icon in system, backing up to theme package")

		// Create Collections directory in backup
		backupCollectionsDir := filepath.Join(exportPath, "Collections")
		if err := os.MkdirAll(backupCollectionsDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup Collections directory: %w", err)
		}

		// Copy Collections icon
		backupCollectionsPath := filepath.Join(backupCollectionsDir, "icon.png")
		if err := CopyFile(systemCollectionsPath, backupCollectionsPath); err != nil {
			app.LogDebug("Warning: Failed to copy Collections icon to backup: %v", err)
		} else {
			app.LogDebug("Successfully backed up Collections icon")
		}
	} else {
		app.LogDebug("No Collections icon found in system")
	}

	// Handle Recently Played icon
	systemRecentlyPlayedPath := filepath.Join(systemMediaPath, "Recently Played.png")
	if _, err := os.Stat(systemRecentlyPlayedPath); err == nil {
		app.LogDebug("Found Recently Played icon in system, backing up to theme package")

		// Create Recently Played directory in backup
		backupRecentlyPlayedDir := filepath.Join(exportPath, "Recently Played")
		if err := os.MkdirAll(backupRecentlyPlayedDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup Recently Played directory: %w", err)
		}

		// Copy Recently Played icon
		backupRecentlyPlayedPath := filepath.Join(backupRecentlyPlayedDir, "icon.png")
		if err := CopyFile(systemRecentlyPlayedPath, backupRecentlyPlayedPath); err != nil {
			app.LogDebug("Warning: Failed to copy Recently Played icon to backup: %v", err)
		} else {
			app.LogDebug("Successfully backed up Recently Played icon")
		}
	} else {
		app.LogDebug("No Recently Played icon found in system")
	}

	// Create manifest for backup
	manifest := CreateBackupManifest(strings.TrimSuffix(themeName, ThemeExtension))

	// Write manifest to backup
	if err := WriteManifest(manifest, exportPath); err != nil {
		return fmt.Errorf("failed to write backup manifest: %w", err)
	}

	app.LogDebug("Backup created successfully: %s", themeName)
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

		// Try to read manifest - use non-strict validation for listing
		manifest, err := ReadManifest(themePath, false)

		// Create theme info
		themeInfo := map[string]string{
			"name": themeName,
		}

		// Add manifest info if available
		if err == nil {
			themeInfo["author"] = manifest.Author
			themeInfo["version"] = manifest.Version
			themeInfo["description"] = manifest.Description

			// Check if manifest is valid
			if IsManifestValid(manifest) {
				themeInfo["is_valid"] = "true"
			} else {
				themeInfo["is_valid"] = "false"
			}
		} else {
			app.LogDebug("Warning: Failed to read manifest for theme %s: %v", themeName, err)
			themeInfo["is_valid"] = "false"
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

		// Try to read manifest - use non-strict validation for listing
		manifest, err := ReadManifest(backupPath, false)

		// Create backup info
		backupInfo := map[string]string{
			"name": backupName,
			"is_valid": "true", // Backups are always considered valid for restore
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