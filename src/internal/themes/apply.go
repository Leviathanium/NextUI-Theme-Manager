
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

	// Copy all files EXCLUDING special directories (Tools/tg5040, Collections, Recently Played)
	excludeSubpaths := []string{"Tools/tg5040", "Collections", "Recently Played"}
	if err := CopyDirectoryExcludingSubpath(themePath, SystemThemeDir, excludeSubpaths); err != nil {
		return fmt.Errorf("failed to copy theme files: %w", err)
	}

	// Handle tool icons specially - copy Tools/tg5040/.media/ to /mnt/SDCARD/Tools/tg5040/.media/
	toolIconsPath := filepath.Join(themePath, "Tools", "tg5040", ".media")
	if _, err := os.Stat(toolIconsPath); err == nil {
		app.LogDebug("Found tool icons in theme, applying to system")

		// Clear existing tool icons
		systemToolsMediaPath := "/mnt/SDCARD/Tools/tg5040/.media"
		if err := ClearToolIcons(systemToolsMediaPath); err != nil {
			app.LogDebug("Warning: Failed to clear existing tool icons: %v", err)
		}

		// Create Tools/.media directory if it doesn't exist
		if err := os.MkdirAll(systemToolsMediaPath, 0755); err != nil {
			return fmt.Errorf("failed to create system tools media directory: %w", err)
		}

		// Copy tool icons
		if err := CopyDirectory(toolIconsPath, systemToolsMediaPath); err != nil {
			return fmt.Errorf("failed to copy tool icons: %w", err)
		}

		app.LogDebug("Successfully applied tool icons to system")
	} else {
		app.LogDebug("No tool icons found in theme")
	}

	// Handle Collections and Recently Played icons specially - copy to /.media/
	systemMediaPath := "/mnt/SDCARD/.media"

	// Create .media directory if it doesn't exist
	if err := os.MkdirAll(systemMediaPath, 0755); err != nil {
		return fmt.Errorf("failed to create system .media directory: %w", err)
	}

	// Handle Collections icon
	collectionsIconPath := filepath.Join(themePath, "Collections", "icon.png")
	if _, err := os.Stat(collectionsIconPath); err == nil {
		app.LogDebug("Found Collections icon in theme, applying to system")
		systemCollectionsPath := filepath.Join(systemMediaPath, "Collections.png")
		if err := CopyFile(collectionsIconPath, systemCollectionsPath); err != nil {
			app.LogDebug("Warning: Failed to copy Collections icon: %v", err)
		} else {
			app.LogDebug("Successfully applied Collections icon")
		}
	} else {
		app.LogDebug("No Collections icon found in theme")
	}

	// Handle Recently Played icon
	recentlyPlayedIconPath := filepath.Join(themePath, "Recently Played", "icon.png")
	if _, err := os.Stat(recentlyPlayedIconPath); err == nil {
		app.LogDebug("Found Recently Played icon in theme, applying to system")
		systemRecentlyPlayedPath := filepath.Join(systemMediaPath, "Recently Played.png")
		if err := CopyFile(recentlyPlayedIconPath, systemRecentlyPlayedPath); err != nil {
			app.LogDebug("Warning: Failed to copy Recently Played icon: %v", err)
		} else {
			app.LogDebug("Successfully applied Recently Played icon")
		}
	} else {
		app.LogDebug("No Recently Played icon found in theme")
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

	// Copy all files from system theme to backup
	if err := CopyDirectory(SystemThemeDir, backupPath); err != nil {
		return fmt.Errorf("failed to copy system theme to backup: %w", err)
	}

	// Copy tool icons to backup (handle tool icons specially)
	systemToolsMediaPath := "/mnt/SDCARD/Tools/tg5040/.media"
	if _, err := os.Stat(systemToolsMediaPath); err == nil {
		app.LogDebug("Found system tool icons, backing up to theme package")

		// Create Tools/tg5040/.media structure in backup
		backupToolsMediaPath := filepath.Join(backupPath, "Tools", "tg5040", ".media")
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
		backupCollectionsDir := filepath.Join(backupPath, "Collections")
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
		backupRecentlyPlayedDir := filepath.Join(backupPath, "Recently Played")
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
	manifest := CreateBackupManifest(strings.TrimSuffix(backupName, ThemeExtension))
	manifest.Description = "Manual backup of system theme"

	// Write manifest to backup
	if err := WriteManifest(manifest, backupPath); err != nil {
		return fmt.Errorf("failed to write backup manifest: %w", err)
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