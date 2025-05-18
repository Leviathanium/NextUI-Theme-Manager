// internal/themes/common.go
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"thememanager/internal/app"
)

// Constants for file extensions and paths
const (
    ThemeExtension      = ".theme"
    ThemeManifestFile   = "manifest.yml"
    ThemePreviewFile    = "preview.png"
    SystemThemeDir      = "/mnt/SDCARD/Theme"
    ThemeBackupNameFmt  = "backup_%s"  // Format string for backup names
)

// EnsureDirectories creates all required directories for the application
func EnsureDirectories() error {
	app.LogDebug("Ensuring required directories exist")

	cwd := app.GetWorkingDir()

	// List of directories to ensure
	dirs := []string{
		filepath.Join(cwd, "Themes"),    // For downloaded themes
		filepath.Join(cwd, "Backups"),   // For theme backups
		filepath.Join(cwd, "Catalog"),   // For catalog syncing
	}

	// Create each directory if it doesn't exist
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			app.LogDebug("Error creating directory %s: %v", dir, err)
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	app.LogDebug("All required directories ensured")
	return nil
}

// ListThemes returns a list of all themes in the Themes directory
func ListThemes() ([]string, error) {
	app.LogDebug("Listing themes in Themes directory")

	themesDir := filepath.Join(app.GetWorkingDir(), "Themes")

	// Check if themes directory exists
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		// Create it if it doesn't exist
		if err := os.MkdirAll(themesDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create Themes directory: %w", err)
		}
		// Return empty list since directory was just created
		return []string{}, nil
	}

	// Read directory
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read Themes directory: %w", err)
	}

	// Filter for theme directories
	var themes []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ThemeExtension) {
			// Remove extension from name
			themeName := strings.TrimSuffix(entry.Name(), ThemeExtension)
			themes = append(themes, themeName)
		}
	}

	app.LogDebug("Found %d themes", len(themes))
	return themes, nil
}

// ListBackups returns a list of all theme backups in the Backups directory
func ListBackups() ([]string, error) {
	app.LogDebug("Listing theme backups in Backups directory")

	backupsDir := filepath.Join(app.GetWorkingDir(), "Backups")

	// Check if backups directory exists
	if _, err := os.Stat(backupsDir); os.IsNotExist(err) {
		// Create it if it doesn't exist
		if err := os.MkdirAll(backupsDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create Backups directory: %w", err)
		}
		// Return empty list since directory was just created
		return []string{}, nil
	}

	// Read directory
	entries, err := os.ReadDir(backupsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read Backups directory: %w", err)
	}

	// Filter for theme backup directories
	var backups []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ThemeExtension) {
			// Remove extension from name
			backupName := strings.TrimSuffix(entry.Name(), ThemeExtension)
			backups = append(backups, backupName)
		}
	}

	app.LogDebug("Found %d backups", len(backups))
	return backups, nil
}

// GetThemePath returns the full path to a theme package
func GetThemePath(themeName string) string {
	if !strings.HasSuffix(themeName, ThemeExtension) {
		themeName += ThemeExtension
	}
	return filepath.Join(app.GetWorkingDir(), "Themes", themeName)
}

// GetBackupPath returns the full path to a theme backup
func GetBackupPath(backupName string) string {
	if !strings.HasSuffix(backupName, ThemeExtension) {
		backupName += ThemeExtension
	}
	return filepath.Join(app.GetWorkingDir(), "Backups", backupName)
}

// GetCatalogPath returns the full path to the catalog directory
func GetCatalogPath() string {
	return filepath.Join(app.GetWorkingDir(), "Catalog")
}

// GetThemePreviewPath returns the path to a theme's preview image
func GetThemePreviewPath(themeName string) string {
	themePath := GetThemePath(themeName)
	return filepath.Join(themePath, ThemePreviewFile)
}

// GetBackupPreviewPath returns the path to a backup's preview image
func GetBackupPreviewPath(backupName string) string {
	backupPath := GetBackupPath(backupName)
	return filepath.Join(backupPath, ThemePreviewFile)
}

// ThemeExists checks if a theme exists in the Themes directory
func ThemeExists(themeName string) bool {
	themePath := GetThemePath(themeName)
	_, err := os.Stat(themePath)
	return err == nil
}

// BackupExists checks if a backup exists in the Backups directory
func BackupExists(backupName string) bool {
	backupPath := GetBackupPath(backupName)
	_, err := os.Stat(backupPath)
	return err == nil
}

// SystemThemeExists checks if the system theme directory exists
func SystemThemeExists() bool {
	_, err := os.Stat(SystemThemeDir)
	return err == nil
}

// CopyDirectory recursively copies a directory
func CopyDirectory(src, dst string) error {
	app.LogDebug("Copying directory %s to %s", src, dst)

	// Get source file info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error getting source directory info: %w", err)
	}

	// Create destination directory with same permissions
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("error creating destination directory: %w", err)
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("error reading source directory: %w", err)
	}

	// Process each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := CopyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a single file
func CopyFile(src, dst string) error {
	app.LogDebug("Copying file %s to %s", src, dst)

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting source file info: %w", err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy content
	bufSize := 1024 * 1024 // 1MB buffer
	buf := make([]byte, bufSize)

	for {
		n, err := srcFile.Read(buf)
		if err != nil && err.Error() != "EOF" {
			return fmt.Errorf("error reading source file: %w", err)
		}
		if n == 0 {
			break
		}

		if _, err := dstFile.Write(buf[:n]); err != nil {
			return fmt.Errorf("error writing to destination file: %w", err)
		}
	}

	return nil
}

// RemoveDirectory recursively removes a directory and all its contents
func RemoveDirectory(path string) error {
	app.LogDebug("Removing directory %s", path)

	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Directory doesn't exist, nothing to do
		return nil
	}

	// Remove directory and all contents
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("error removing directory %s: %w", path, err)
	}

	return nil
}