// internal/themes/themes.go
package themes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"thememanager/internal/logging"
	"thememanager/internal/system"
	"thememanager/internal/ui"
)

// Define maximum number of backups to keep
const MaxBackups = 3

// ThemeManifest represents the structure of a manifest.yml file
type ThemeManifest struct {
	Info struct {
		Name        string    `yaml:"name"`
		Author      string    `yaml:"author"`
		Version     string    `yaml:"version"`
		CreatedDate time.Time `yaml:"created_date"`
	} `yaml:"info"`

	Content struct {
		Backgrounds bool `yaml:"backgrounds"`
		Icons       bool `yaml:"icons"`
		Fonts       bool `yaml:"fonts"`
		Accents     bool `yaml:"accents"`
	} `yaml:"content"`

	// Additional fields as needed
}

// OverlayManifest represents the structure of an overlay manifest.yml file
type OverlayManifest struct {
	Info struct {
		Name        string    `yaml:"name"`
		Author      string    `yaml:"author"`
		Version     string    `yaml:"version"`
		CreatedDate time.Time `yaml:"created_date"`
	} `yaml:"info"`

	Content struct {
		Systems []string `yaml:"systems"`
	} `yaml:"content"`
}

// EnsureDirectoryStructure creates all necessary directories for the application
func EnsureDirectoryStructure() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create required directories
	dirs := []string{
		filepath.Join(cwd, "Themes"),
		filepath.Join(cwd, "Overlays"),
		filepath.Join(cwd, "Backups", "Themes"),
		filepath.Join(cwd, "Backups", "Overlays"),
		filepath.Join(cwd, "Catalog", "Themes"),
		filepath.Join(cwd, "Catalog", "Overlays"),
		filepath.Join(cwd, "Logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logging.LogDebug("Error creating directory %s: %v", dir, err)
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	return nil
}

// IsThemeDownloaded checks if a theme is already downloaded
func IsThemeDownloaded(themeName string) bool {
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return false
	}

	themePath := filepath.Join(cwd, "Themes", themeName+".theme")
	_, err = os.Stat(themePath)
	return err == nil
}

// IsOverlayDownloaded checks if an overlay pack is already downloaded
func IsOverlayDownloaded(overlayName string) bool {
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return false
	}

	overlayPath := filepath.Join(cwd, "Overlays", overlayName+".over")
	_, err = os.Stat(overlayPath)
	return err == nil
}

// ShowThemeGallery displays a gallery of available themes
func ShowThemeGallery() (string, int) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read catalog directory for themes
	catalogPath := filepath.Join(cwd, "Catalog", "Themes")
	entries, err := os.ReadDir(catalogPath)
	if err != nil {
		logging.LogDebug("Error reading catalog directory: %v", err)
		ui.ShowMessage("No themes found in catalog. Please sync first.", "3")
		return "", 1
	}

	// Create gallery items
	var galleryItems []ui.GalleryItem
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".theme") {
			themeName := entry.Name()

			// Get preview image path
			previewPath := filepath.Join(catalogPath, themeName, "preview.png")

			// Get theme information from manifest
			manifestPath := filepath.Join(catalogPath, themeName, "manifest.yml")
			manifest, err := readThemeManifest(manifestPath)

			// Default text
			text := themeName

			// Check if theme is already downloaded/installed
			isInstalled := IsThemeDownloaded(themeName)
			if isInstalled {
				text = "[Installed] " + text
			}

			// Add author if available
			if err == nil && manifest.Info.Author != "" {
				text = fmt.Sprintf("%s by %s", text, manifest.Info.Author)
			}

			// Create gallery item
			galleryItems = append(galleryItems, ui.GalleryItem{
				Text:            text,
				BackgroundImage: previewPath,
			})
		}
	}

	// Check if we found any themes
	if len(galleryItems) == 0 {
		ui.ShowMessage("No themes found in catalog. Please sync first.", "3")
		return "", 1
	}

	// Display the gallery
	return ui.DisplayImageGallery(galleryItems, "Available Themes")
}

// ShowOverlayGallery displays a gallery of available overlay packs
func ShowOverlayGallery() (string, int) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read catalog directory for overlays
	catalogPath := filepath.Join(cwd, "Catalog", "Overlays")
	entries, err := os.ReadDir(catalogPath)
	if err != nil {
		logging.LogDebug("Error reading catalog directory: %v", err)
		ui.ShowMessage("No overlays found in catalog. Please sync first.", "3")
		return "", 1
	}

	// Create gallery items
	var galleryItems []ui.GalleryItem
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".over") {
			overlayName := entry.Name()

			// Get preview image path
			previewPath := filepath.Join(catalogPath, overlayName, "preview.png")

			// Get overlay information from manifest
			manifestPath := filepath.Join(catalogPath, overlayName, "manifest.yml")
			manifest, err := readOverlayManifest(manifestPath)

			// Default text
			text := overlayName

			// Check if overlay is already downloaded/installed
			isInstalled := IsOverlayDownloaded(overlayName)
			if isInstalled {
				text = "[Installed] " + text
			}

			// Add author if available
			if err == nil && manifest.Info.Author != "" {
				text = fmt.Sprintf("%s by %s", text, manifest.Info.Author)
			}

			// Create gallery item
			galleryItems = append(galleryItems, ui.GalleryItem{
				Text:            text,
				BackgroundImage: previewPath,
			})
		}
	}

	// Check if we found any overlays
	if len(galleryItems) == 0 {
		ui.ShowMessage("No overlays found in catalog. Please sync first.", "3")
		return "", 1
	}

	// Display the gallery
	return ui.DisplayImageGallery(galleryItems, "Available Overlays")
}

// ShowThemeBackupGallery displays a gallery of theme backups
func ShowThemeBackupGallery() (string, int) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read backups directory for themes
	backupsPath := filepath.Join(cwd, "Backups", "Themes")
	entries, err := os.ReadDir(backupsPath)
	if err != nil {
		logging.LogDebug("Error reading backups directory: %v", err)
		ui.ShowMessage("No theme backups found.", "3")
		return "", 1
	}

	// Create gallery items
	var galleryItems []ui.GalleryItem
	for _, entry := range entries {
		if entry.IsDir() {
			backupName := entry.Name()

			// Try to find a preview image in the backup
			previewPath := filepath.Join(backupsPath, backupName, "preview.png")
			if _, err := os.Stat(previewPath); os.IsNotExist(err) {
				// No preview, check if there's a screenshot
				previewPath = filepath.Join(backupsPath, backupName, "screenshot.png")
				if _, err := os.Stat(previewPath); os.IsNotExist(err) {
					// No screenshot either, use a blank preview
					previewPath = ""
				}
			}

			// Create gallery item
			galleryItems = append(galleryItems, ui.GalleryItem{
				Text:            backupName,
				BackgroundImage: previewPath,
			})
		}
	}

	// Check if we found any backups
	if len(galleryItems) == 0 {
		ui.ShowMessage("No theme backups found.", "3")
		return "", 1
	}

	// Display the gallery
	return ui.DisplayImageGallery(galleryItems, "Theme Backups")
}

// ShowOverlayBackupGallery displays a gallery of overlay backups
func ShowOverlayBackupGallery() (string, int) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read backups directory for overlays
	backupsPath := filepath.Join(cwd, "Backups", "Overlays")
	entries, err := os.ReadDir(backupsPath)
	if err != nil {
		logging.LogDebug("Error reading backups directory: %v", err)
		ui.ShowMessage("No overlay backups found.", "3")
		return "", 1
	}

	// Create gallery items
	var galleryItems []ui.GalleryItem
	for _, entry := range entries {
		if entry.IsDir() {
			backupName := entry.Name()

			// Try to find a preview image in the backup
			previewPath := filepath.Join(backupsPath, backupName, "preview.png")
			if _, err := os.Stat(previewPath); os.IsNotExist(err) {
				// No preview, use a blank preview
				previewPath = ""
			}

			// Create gallery item
			galleryItems = append(galleryItems, ui.GalleryItem{
				Text:            backupName,
				BackgroundImage: previewPath,
			})
		}
	}

	// Check if we found any backups
	if len(galleryItems) == 0 {
		ui.ShowMessage("No overlay backups found.", "3")
		return "", 1
	}

	// Display the gallery
	return ui.DisplayImageGallery(galleryItems, "Overlay Backups")
}

// DownloadTheme downloads a theme from the catalog
func DownloadTheme(themeName string) error {
	logging.LogDebug("Downloading theme: %s", themeName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Source and destination paths
	srcPath := filepath.Join(cwd, "Catalog", "Themes", themeName+".theme")
	dstPath := filepath.Join(cwd, "Themes", themeName+".theme")

	// Check if source exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("theme not found in catalog: %s", themeName)
	}

	// Create destination directory
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return fmt.Errorf("error creating theme directory: %w", err)
	}

	// Copy theme files
	if err := copyDir(srcPath, dstPath); err != nil {
		return fmt.Errorf("error copying theme files: %w", err)
	}

	logging.LogDebug("Theme downloaded successfully: %s", themeName)
	return nil
}

// DownloadOverlay downloads an overlay pack from the catalog
func DownloadOverlay(overlayName string) error {
	logging.LogDebug("Downloading overlay: %s", overlayName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Source and destination paths
	srcPath := filepath.Join(cwd, "Catalog", "Overlays", overlayName+".over")
	dstPath := filepath.Join(cwd, "Overlays", overlayName+".over")

	// Check if source exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("overlay not found in catalog: %s", overlayName)
	}

	// Create destination directory
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return fmt.Errorf("error creating overlay directory: %w", err)
	}

	// Copy overlay files
	if err := copyDir(srcPath, dstPath); err != nil {
		return fmt.Errorf("error copying overlay files: %w", err)
	}

	logging.LogDebug("Overlay downloaded successfully: %s", overlayName)
	return nil
}

// ApplyTheme applies a theme to the system
func ApplyTheme(themeName string) error {
	logging.LogDebug("Applying theme: %s", themeName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Theme path
	themePath := filepath.Join(cwd, "Themes", themeName+".theme")

	// Check if theme exists
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return fmt.Errorf("theme not found: %s", themeName)
	}

	// Read manifest to determine what to apply
	manifestPath := filepath.Join(themePath, "manifest.yml")
	manifest, err := readThemeManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("error reading theme manifest: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Apply backgrounds if present
	if manifest.Content.Backgrounds {
		if err := applyBackgrounds(themePath, systemPaths); err != nil {
			logging.LogDebug("Error applying backgrounds: %v", err)
			// Continue with other components
		}
	}

	// Apply icons if present
	if manifest.Content.Icons {
		if err := applyIcons(themePath, systemPaths); err != nil {
			logging.LogDebug("Error applying icons: %v", err)
			// Continue with other components
		}
	}

	// Apply fonts if present
	if manifest.Content.Fonts {
		if err := applyFonts(themePath); err != nil {
			logging.LogDebug("Error applying fonts: %v", err)
			// Continue with other components
		}
	}

	// Apply accents if present
	if manifest.Content.Accents {
		if err := applyAccents(themePath); err != nil {
			logging.LogDebug("Error applying accents: %v", err)
			// Continue with other components
		}
	}

	logging.LogDebug("Theme applied successfully: %s", themeName)
	return nil
}

// ApplyOverlay applies an overlay pack to the system
func ApplyOverlay(overlayName string) error {
	logging.LogDebug("Applying overlay: %s", overlayName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Overlay path
	overlayPath := filepath.Join(cwd, "Overlays", overlayName+".over")

	// Check if overlay exists
	if _, err := os.Stat(overlayPath); os.IsNotExist(err) {
		return fmt.Errorf("overlay not found: %s", overlayName)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Apply overlays
	if err := applyOverlays(overlayPath, systemPaths); err != nil {
		return fmt.Errorf("error applying overlays: %w", err)
	}

	logging.LogDebug("Overlay applied successfully: %s", overlayName)
	return nil
}

// SyncCatalog synchronizes the theme and overlay catalog
func SyncCatalog() error {
	logging.LogDebug("Syncing catalog")

	// In a real implementation, this would download from a repository
	// For now, we'll just create a placeholder method

	// Create sample themes and overlays in catalog if none exist
	createSampleCatalogContents()

	logging.LogDebug("Catalog sync completed")
	return nil
}

// CreateThemeBackup creates a backup of the current theme settings
func CreateThemeBackup(backupType string) error {
	logging.LogDebug("Creating theme backup (type: %s)", backupType)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create backup name with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s_%s", backupType, timestamp)
	backupPath := filepath.Join(cwd, "Backups", "Themes", backupName)

	// Create backup directory
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("error creating backup directory: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Back up the current theme settings
	if err := backupCurrentTheme(backupPath, systemPaths); err != nil {
		return fmt.Errorf("error backing up theme: %w", err)
	}

	// Maintain maximum number of backups
	if err := pruneOldBackups("Themes", MaxBackups); err != nil {
		logging.LogDebug("Warning: Error pruning old backups: %v", err)
	}

	logging.LogDebug("Theme backup created: %s", backupName)
	return nil
}

// CreateOverlayBackup creates a backup of the current overlay settings
func CreateOverlayBackup(backupType string) error {
	logging.LogDebug("Creating overlay backup (type: %s)", backupType)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create backup name with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s_%s", backupType, timestamp)
	backupPath := filepath.Join(cwd, "Backups", "Overlays", backupName)

	// Create backup directory
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("error creating backup directory: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Back up the current overlay settings
	if err := backupCurrentOverlays(backupPath, systemPaths); err != nil {
		return fmt.Errorf("error backing up overlays: %w", err)
	}

	// Maintain maximum number of backups
	if err := pruneOldBackups("Overlays", MaxBackups); err != nil {
		logging.LogDebug("Warning: Error pruning old backups: %v", err)
	}

	logging.LogDebug("Overlay backup created: %s", backupName)
	return nil
}

// RevertThemeFromBackup restores theme settings from a backup
func RevertThemeFromBackup(backupName string) error {
	logging.LogDebug("Reverting theme from backup: %s", backupName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Backup path
	backupPath := filepath.Join(cwd, "Backups", "Themes", backupName)

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Apply backup
	if err := restoreThemeFromBackup(backupPath, systemPaths); err != nil {
		return fmt.Errorf("error restoring theme: %w", err)
	}

	logging.LogDebug("Theme reverted successfully from backup: %s", backupName)
	return nil
}

// RevertOverlayFromBackup restores overlay settings from a backup
func RevertOverlayFromBackup(backupName string) error {
	logging.LogDebug("Reverting overlays from backup: %s", backupName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Backup path
	backupPath := filepath.Join(cwd, "Backups", "Overlays", backupName)

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Apply backup
	if err := restoreOverlaysFromBackup(backupPath, systemPaths); err != nil {
		return fmt.Errorf("error restoring overlays: %w", err)
	}

	logging.LogDebug("Overlays reverted successfully from backup: %s", backupName)
	return nil
}

// PurgeAll removes all themes, overlays, and backups
func PurgeAll() error {
	logging.LogDebug("Purging all theme and overlay data")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Remove all theme-related files from the system
	if err := cleanAllThemeFiles(systemPaths); err != nil {
		logging.LogDebug("Error cleaning theme files: %v", err)
		// Continue with other operations
	}

	// Remove all overlay-related files from the system
	if err := cleanAllOverlayFiles(systemPaths); err != nil {
		logging.LogDebug("Error cleaning overlay files: %v", err)
		// Continue with other operations
	}

	// Delete all installed themes
	themesDir := filepath.Join(cwd, "Themes")
	if err := os.RemoveAll(themesDir); err != nil {
		logging.LogDebug("Error removing themes directory: %v", err)
		// Continue with other operations
	}
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		logging.LogDebug("Error recreating themes directory: %v", err)
	}

	// Delete all installed overlays
	overlaysDir := filepath.Join(cwd, "Overlays")
	if err := os.RemoveAll(overlaysDir); err != nil {
		logging.LogDebug("Error removing overlays directory: %v", err)
		// Continue with other operations
	}
	if err := os.MkdirAll(overlaysDir, 0755); err != nil {
		logging.LogDebug("Error recreating overlays directory: %v", err)
	}

	// Delete all backups
	backupsDir := filepath.Join(cwd, "Backups")
	if err := os.RemoveAll(backupsDir); err != nil {
		logging.LogDebug("Error removing backups directory: %v", err)
		// Continue with other operations
	}
	if err := os.MkdirAll(filepath.Join(backupsDir, "Themes"), 0755); err != nil {
		logging.LogDebug("Error recreating theme backups directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(backupsDir, "Overlays"), 0755); err != nil {
		logging.LogDebug("Error recreating overlay backups directory: %v", err)
	}

	logging.LogDebug("Purge completed successfully")
	return nil
}

// Helper functions

// readThemeManifest reads and parses a theme manifest.yml file
func readThemeManifest(path string) (*ThemeManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest file: %w", err)
	}

	var manifest ThemeManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("error parsing manifest: %w", err)
	}

	return &manifest, nil
}

// readOverlayManifest reads and parses an overlay manifest.yml file
func readOverlayManifest(path string) (*OverlayManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest file: %w", err)
	}

	var manifest OverlayManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("error parsing manifest: %w", err)
	}

	return &manifest, nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	// Get file info
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory with same permissions
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy directory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get file info
	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

// Implementation of theme application functions

// applyBackgrounds applies background images from a theme
func applyBackgrounds(themePath string, systemPaths *system.SystemPaths) error {
	// Clean existing backgrounds first
	cleanBackgrounds(systemPaths)

	// Apply main menu backgrounds
	mainMenuBgDir := filepath.Join(themePath, "Backgrounds", "MainMenuBackgrounds")
	if _, err := os.Stat(mainMenuBgDir); err == nil {
		if err := applyMainMenuBackgrounds(mainMenuBgDir, systemPaths); err != nil {
			return err
		}
	}

	// Apply system backgrounds
	systemBgDir := filepath.Join(themePath, "Backgrounds", "SystemBackgrounds")
	if _, err := os.Stat(systemBgDir); err == nil {
		if err := applySystemBackgrounds(systemBgDir, systemPaths); err != nil {
			return err
		}
	}

	return nil
}

// applyMainMenuBackgrounds applies backgrounds for the main menu
func applyMainMenuBackgrounds(srcDir string, systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	// This would copy bg.png files to their appropriate destinations
	return nil
}

// applySystemBackgrounds applies bglist.png backgrounds for system lists
func applySystemBackgrounds(srcDir string, systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	// This would copy bglist.png files to their appropriate destinations
	return nil
}

// applyIcons applies icons from a theme
func applyIcons(themePath string, systemPaths *system.SystemPaths) error {
	// Clean existing icons first
	cleanIcons(systemPaths)

	// Apply main menu icons
	mainMenuIconsDir := filepath.Join(themePath, "Icons", "MainMenuIcons")
	if _, err := os.Stat(mainMenuIconsDir); err == nil {
		if err := applyMainMenuIcons(mainMenuIconsDir, systemPaths); err != nil {
			return err
		}
	}

	// Apply collection icons
	collectionIconsDir := filepath.Join(themePath, "Icons", "CollectionIcons")
	if _, err := os.Stat(collectionIconsDir); err == nil {
		if err := applyCollectionIcons(collectionIconsDir, systemPaths); err != nil {
			return err
		}
	}

	// Apply tool icons
	toolIconsDir := filepath.Join(themePath, "Icons", "ToolIcons")
	if _, err := os.Stat(toolIconsDir); err == nil {
		if err := applyToolIcons(toolIconsDir, systemPaths); err != nil {
			return err
		}
	}

	return nil
}

// applyMainMenuIcons applies icons for the main menu
func applyMainMenuIcons(srcDir string, systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	// This would copy icon files to their appropriate destinations
	return nil
}

// applyCollectionIcons applies icons for collections
func applyCollectionIcons(srcDir string, systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	return nil
}

// applyToolIcons applies icons for tools
func applyToolIcons(srcDir string, systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	return nil
}

// applyFonts applies fonts from a theme
func applyFonts(themePath string) error {
	// Implementation details would go here
	return nil
}

// applyAccents applies accent colors from a theme
func applyAccents(themePath string) error {
	// Implementation details would go here
	return nil
}

// applyOverlays applies overlay images from an overlay pack
func applyOverlays(overlayPath string, systemPaths *system.SystemPaths) error {
	// Clean existing overlays first
	cleanOverlays(systemPaths)

	// Implementation details would go here
	return nil
}

// Backup and restore functions

// backupCurrentTheme backs up the current theme settings
func backupCurrentTheme(backupPath string, systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	return nil
}

// backupCurrentOverlays backs up the current overlay settings
func backupCurrentOverlays(backupPath string, systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	return nil
}

// restoreThemeFromBackup restores theme settings from a backup
func restoreThemeFromBackup(backupPath string, systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	return nil
}

// restoreOverlaysFromBackup restores overlay settings from a backup
func restoreOverlaysFromBackup(backupPath string, systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	return nil
}

// Cleaning functions

// cleanBackgrounds removes all background images
func cleanBackgrounds(systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	return nil
}

// cleanIcons removes all icon images
func cleanIcons(systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	return nil
}

// cleanOverlays removes all overlay images
func cleanOverlays(systemPaths *system.SystemPaths) error {
	// Implementation details would go here
	return nil
}

// cleanAllThemeFiles removes all theme-related files
func cleanAllThemeFiles(systemPaths *system.SystemPaths) error {
	if err := cleanBackgrounds(systemPaths); err != nil {
		return err
	}

	if err := cleanIcons(systemPaths); err != nil {
		return err
	}

	// Clean fonts and accents

	return nil
}

// cleanAllOverlayFiles removes all overlay-related files
func cleanAllOverlayFiles(systemPaths *system.SystemPaths) error {
	return cleanOverlays(systemPaths)
}

// pruneOldBackups maintains the maximum number of backups
func pruneOldBackups(backupType string, maxBackups int) error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Backups directory
	backupsDir := filepath.Join(cwd, "Backups", backupType)

	// Read backups directory
	entries, err := os.ReadDir(backupsDir)
	if err != nil {
		return fmt.Errorf("error reading backups directory: %w", err)
	}

	// Sort backups by creation time (oldest first)
	type backupInfo struct {
		name string
		time time.Time
	}

	var backups []backupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			backups = append(backups, backupInfo{
				name: entry.Name(),
				time: info.ModTime(),
			})
		}
	}

	// Sort by modification time (oldest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].time.Before(backups[j].time)
	})

	// Remove oldest backups if we have more than the maximum
	if len(backups) > maxBackups {
		for i := 0; i < len(backups)-maxBackups; i++ {
			backupPath := filepath.Join(backupsDir, backups[i].name)
			if err := os.RemoveAll(backupPath); err != nil {
				logging.LogDebug("Error removing old backup: %v", err)
				// Continue with other backups
			} else {
				logging.LogDebug("Removed old backup: %s", backups[i].name)
			}
		}
	}

	return nil
}

// createSampleCatalogContents creates sample themes and overlays in the catalog
func createSampleCatalogContents() {
	// This is a placeholder for actual sync logic that would download from a repository
	// For testing purposes, it would create sample themes and overlays
}