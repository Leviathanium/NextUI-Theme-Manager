// internal/themes/import.go
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"thememanager/internal/logging"
	"thememanager/internal/system"

)

// DownloadTheme downloads a theme from the catalog
func DownloadTheme(themeName string) error {
	logging.LogDebug("Downloading theme: %s", themeName)

	// Get source and destination paths
	srcPath, err := GetCatalogThemePath(themeName)
	if err != nil {
		return fmt.Errorf("error getting catalog theme path: %w", err)
	}

	dstPath, err := GetThemePath(themeName)
	if err != nil {
		return fmt.Errorf("error getting theme path: %w", err)
	}

	// Check if source exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("theme not found in catalog: %s", themeName)
	}

	// Check if destination already exists
	if _, err := os.Stat(dstPath); err == nil {
		// Already exists, remove it
		logging.LogDebug("Theme already exists, removing it: %s", dstPath)
		if err := os.RemoveAll(dstPath); err != nil {
			return fmt.Errorf("error removing existing theme: %w", err)
		}
	}

	// Create destination directory
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return fmt.Errorf("error creating theme directory: %w", err)
	}

	// Copy theme files
	if err := CopyDir(srcPath, dstPath); err != nil {
		return fmt.Errorf("error copying theme files: %w", err)
	}

	logging.LogDebug("Theme downloaded successfully: %s", themeName)
	return nil
}

// ApplyTheme applies a theme to the system
func ApplyTheme(themeName string) error {
	logging.LogDebug("Applying theme: %s", themeName)

	// Get theme path
	themePath, err := GetThemePath(themeName)
	if err != nil {
		return fmt.Errorf("error getting theme path: %w", err)
	}

	// Check if theme exists
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return fmt.Errorf("theme not found: %s", themeName)
	}

	// Read manifest
	manifestPath := system.GetThemeManifestPath(themePath)
	manifest, err := ReadThemeManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("error reading theme manifest: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Ensure all necessary media directories exist
	if err := EnsureMediaDirectories(systemPaths); err != nil {
		logging.LogDebug("Warning: Error ensuring media directories: %v", err)
		// Continue anyway
	}

	// Clean up existing components before applying new ones
	if manifest.Content.Backgrounds {
		cleanBackgrounds(systemPaths)
	}

	if manifest.Content.Icons {
		cleanIcons(systemPaths)
	}

	// Apply theme components
	for tag, sysConfig := range manifest.Systems {
		if err := applySystemFiles(themePath, tag, sysConfig); err != nil {
			logging.LogDebug("Warning: Error applying files for system %s: %v", tag, err)
			// Continue with other systems
		}
	}

	// Apply accent settings if present
	if manifest.Content.Accents {
		if err := applyAccentSettings(themePath); err != nil {
			logging.LogDebug("Warning: Error applying accent settings: %v", err)
			// Continue with other components
		}
	}

	logging.LogDebug("Theme applied successfully: %s", themeName)
	return nil
}

// DownloadOverlay downloads an overlay pack from the catalog
func DownloadOverlay(overlayName string) error {
	logging.LogDebug("Downloading overlay: %s", overlayName)

	// Get source and destination paths
	srcPath, err := GetCatalogOverlayPath(overlayName)
	if err != nil {
		return fmt.Errorf("error getting catalog overlay path: %w", err)
	}

	dstPath, err := GetOverlayPath(overlayName)
	if err != nil {
		return fmt.Errorf("error getting overlay path: %w", err)
	}

	// Check if source exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("overlay not found in catalog: %s", overlayName)
	}

	// Check if destination already exists
	if _, err := os.Stat(dstPath); err == nil {
		// Already exists, remove it
		logging.LogDebug("Overlay already exists, removing it: %s", dstPath)
		if err := os.RemoveAll(dstPath); err != nil {
			return fmt.Errorf("error removing existing overlay: %w", err)
		}
	}

	// Create destination directory
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return fmt.Errorf("error creating overlay directory: %w", err)
	}

	// Copy overlay files
	if err := CopyDir(srcPath, dstPath); err != nil {
		return fmt.Errorf("error copying overlay files: %w", err)
	}

	logging.LogDebug("Overlay downloaded successfully: %s", overlayName)
	return nil
}

// ApplyOverlay applies an overlay pack to the system
func ApplyOverlay(overlayName string) error {
	logging.LogDebug("Applying overlay: %s", overlayName)

	// Get overlay path
	overlayPath, err := GetOverlayPath(overlayName)
	if err != nil {
		return fmt.Errorf("error getting overlay path: %w", err)
	}

	// Check if overlay exists
	if _, err := os.Stat(overlayPath); os.IsNotExist(err) {
		return fmt.Errorf("overlay not found: %s", overlayName)
	}

	// Read manifest
	manifestPath := filepath.Join(overlayPath, "manifest.yml")
	manifest, err := ReadOverlayManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("error reading overlay manifest: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Clean up existing overlays
	cleanOverlays(systemPaths)

	// Apply overlay files
	for _, systemTag := range manifest.Content.Systems {
		srcDir := filepath.Join(overlayPath, "Overlays", systemTag)

		// Skip if directory doesn't exist
		if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			continue
		}

		// Create system overlay directory
		dstDir := system.GetOverlaySystemPath(systemTag)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logging.LogDebug("Warning: Could not create overlay directory for system %s: %v", systemTag, err)
			continue
		}

		// Copy overlay files
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			logging.LogDebug("Warning: Could not read overlay directory for system %s: %v", systemTag, err)
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".png") {
				continue
			}

			srcFile := filepath.Join(srcDir, entry.Name())
			dstFile := filepath.Join(dstDir, entry.Name())

			if err := CopyFile(srcFile, dstFile); err != nil {
				logging.LogDebug("Warning: Could not copy overlay file %s: %v", entry.Name(), err)
			}
		}
	}

	logging.LogDebug("Overlay applied successfully: %s", overlayName)
	return nil
}

// SyncCatalog synchronizes the theme and overlay catalog
func SyncCatalog() error {
	logging.LogDebug("Syncing catalog")

	// In a real implementation, this would download from a repository
	// For now, we'll just create sample content

	// Create sample themes and overlays in catalog
	createSampleCatalogContents()

	logging.LogDebug("Catalog sync completed")
	return nil
}

// Helper functions

// applySystemFiles applies all files for a specific system from a theme
func applySystemFiles(themePath, tag string, config SystemConfig) error {
	logging.LogDebug("Applying files for system: %s", tag)

	for fileType, relPath := range config.Files {
		// Get source file path
		srcPath := filepath.Join(themePath, relPath)
		if !FileExists(srcPath) {
			logging.LogDebug("Source file does not exist: %s", srcPath)
			continue
		}

		// Get destination path from the paths map
		dstPath, ok := config.Paths[fileType+"_path"]
		if !ok {
			logging.LogDebug("No destination path found for file type: %s", fileType)
			continue
		}

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logging.LogDebug("Error creating destination directory %s: %v", dstDir, err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logging.LogDebug("Error copying file from %s to %s: %v", srcPath, dstPath, err)
			continue
		}

		logging.LogDebug("Applied %s: %s -> %s", fileType, srcPath, dstPath)
	}

	return nil
}

// applyAccentSettings applies accent settings from a theme
func applyAccentSettings(themePath string) error {
	// Source file
	srcPath := filepath.Join(themePath, "Settings", "minuisettings.txt")
	if !FileExists(srcPath) {
		return fmt.Errorf("accent settings file not found")
	}

	// Destination file
	dstPath := system.AccentSettingsPath

	// Create destination directory
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("error creating settings directory: %w", err)
	}

	// Copy the file
	if err := CopyFile(srcPath, dstPath); err != nil {
		return fmt.Errorf("error applying accent settings: %w", err)
	}

	logging.LogDebug("Applied accent settings")
	return nil
}

// cleanBackgrounds removes existing background images
func cleanBackgrounds(systemPaths *system.SystemPaths) error {
	logging.LogDebug("Cleaning up existing backgrounds")

	// Clean root backgrounds
	if FileExists(system.GetRootBackgroundPath(false)) {
		os.Remove(system.GetRootBackgroundPath(false))
	}

	if FileExists(system.GetRootBackgroundPath(true)) {
		os.Remove(system.GetRootBackgroundPath(true))
	}

	// Clean Recently Played background
	if FileExists(system.GetRecentlyPlayedBackgroundPath()) {
		os.Remove(system.GetRecentlyPlayedBackgroundPath())
	}

	// Clean Tools background
	if FileExists(system.GetToolsBackgroundPath()) {
		os.Remove(system.GetToolsBackgroundPath())
	}

	// Clean Collections background
	if FileExists(system.GetCollectionsBackgroundPath()) {
		os.Remove(system.GetCollectionsBackgroundPath())
	}

	// Clean system backgrounds
	for _, sysInfo := range systemPaths.Systems {
		// Background
		bgPath := system.GetSystemBackgroundPath(sysInfo.Name)
		if FileExists(bgPath) {
			os.Remove(bgPath)
		}

		// List background
		listBgPath := system.GetSystemListBackgroundPath(sysInfo.Name)
		if FileExists(listBgPath) {
			os.Remove(listBgPath)
		}
	}

	// Clean collection backgrounds
	collectionsDir := system.CollectionsPath
	if entries, err := os.ReadDir(collectionsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			collectionName := entry.Name()
			bgPath := system.GetCollectionBackgroundPath(collectionName)
			if FileExists(bgPath) {
				os.Remove(bgPath)
			}
		}
	}

	return nil
}

// cleanIcons removes existing icon images
func cleanIcons(systemPaths *system.SystemPaths) error {
	logging.LogDebug("Cleaning up existing icons")

	// Clean special icons
	if FileExists(system.GetRecentlyPlayedIconPath()) {
		os.Remove(system.GetRecentlyPlayedIconPath())
	}

	if FileExists(system.GetToolsIconPath()) {
		os.Remove(system.GetToolsIconPath())
	}

	if FileExists(system.GetCollectionsIconPath()) {
		os.Remove(system.GetCollectionsIconPath())
	}

	// Clean system icons in Roms/.media
	systemIconsDir := system.RomsMediaPath
	if entries, err := os.ReadDir(systemIconsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || !strings.HasSuffix(entry.Name(), ".png") {
				continue
			}

			// Skip special icons we already handled
			if entry.Name() == "Recently Played.png" ||
			   entry.Name() == "Collections.png" ||
			   entry.Name() == "tg5040.png" {
				continue
			}

			iconPath := filepath.Join(systemIconsDir, entry.Name())
			os.Remove(iconPath)
		}
	}

	// Clean tool icons
	toolsDir := system.ToolsPath
	if entries, err := os.ReadDir(toolsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			toolName := entry.Name()
			iconPath := system.GetToolIconPath(toolName)
			if FileExists(iconPath) {
				os.Remove(iconPath)
			}
		}
	}

	// Clean collection icons
	collectionsDir := system.CollectionsPath
	if entries, err := os.ReadDir(collectionsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			collectionName := entry.Name()
			iconPath := system.GetCollectionIconPath(collectionName)
			if FileExists(iconPath) {
				os.Remove(iconPath)
			}
		}
	}

	return nil
}

// cleanOverlays removes all existing overlay images
func cleanOverlays(systemPaths *system.SystemPaths) error {
	logging.LogDebug("Cleaning up existing overlays")

	// Get overlays directory
	overlaysDir := system.OverlaysPath

	// Check if it exists
	if _, err := os.Stat(overlaysDir); os.IsNotExist(err) {
		// Create it if it doesn't exist
		os.MkdirAll(overlaysDir, 0755)
		return nil
	}

	// Read directory
	entries, err := os.ReadDir(overlaysDir)
	if err != nil {
		return fmt.Errorf("error reading overlays directory: %w", err)
	}

	// Remove each system's overlays
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		systemTag := entry.Name()
		systemDir := filepath.Join(overlaysDir, systemTag)

		// Remove all files in directory
		if err := CleanDirectory(systemDir); err != nil {
			logging.LogDebug("Warning: Could not clean overlay directory for system %s: %v", systemTag, err)
		}
	}

	return nil
}

// createSampleCatalogContents creates sample themes and overlays for testing
func createSampleCatalogContents() {
	logging.LogDebug("Creating sample catalog contents")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return
	}

	// Catalog paths
	catalogThemesDir := filepath.Join(cwd, "Catalog", "Themes")
	catalogOverlaysDir := filepath.Join(cwd, "Catalog", "Overlays")

	// Create sample themes
	sampleThemes := []struct {
		name   string
		author string
	}{
		{"RetroWave", "PixelArtist"},
		{"NeonFuture", "SynthDesigner"},
		{"Minimal", "CleanUI"},
	}

	for _, theme := range sampleThemes {
		themePath := filepath.Join(catalogThemesDir, theme.name+system.ThemeExtension)

		// Skip if it already exists
		if _, err := os.Stat(themePath); err == nil {
			continue
		}

		// Create directory
		os.MkdirAll(themePath, 0755)

		// Create empty subdirectories
		os.MkdirAll(filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers"), 0755)
		os.MkdirAll(filepath.Join(themePath, system.ThemeWallpapersDir, "ListWallpapers"), 0755)
		os.MkdirAll(filepath.Join(themePath, system.ThemeIconsDir, "SystemIcons"), 0755)
		os.MkdirAll(filepath.Join(themePath, system.ThemeOverlaysDir), 0755)

		// Create manifest
		manifest := CreateEmptyThemeManifest(theme.name, theme.author)
		manifest.Content.Backgrounds = true
		manifest.Content.Icons = true

		// Write manifest
		manifestPath := filepath.Join(themePath, "manifest.yml")
		if err := WriteThemeManifest(manifest, manifestPath); err != nil {
			logging.LogDebug("Error writing sample theme manifest: %v", err)
		}

		// Create empty preview file
		previewPath := filepath.Join(themePath, "preview.png")
		CreateEmptyFile(previewPath)
	}

	// Create sample overlays
	sampleOverlays := []struct {
		name    string
		author  string
		systems []string
	}{
		{"RetroSystem", "OverlayMaker", []string{"SNES", "NES", "GBA"}},
		{"ArcadeCabinet", "CabinetDesigner", []string{"MAME", "CPS1", "CPS2"}},
		{"Transparent", "MinimalDesigner", []string{"ALL"}},
	}

	for _, overlay := range sampleOverlays {
		overlayPath := filepath.Join(catalogOverlaysDir, overlay.name+system.OverlayExtension)

		// Skip if it already exists
		if _, err := os.Stat(overlayPath); err == nil {
			continue
		}

		// Create directory
		os.MkdirAll(overlayPath, 0755)

		// Create system subdirectories
		for _, sys := range overlay.systems {
			os.MkdirAll(filepath.Join(overlayPath, system.ThemeOverlaysDir, sys), 0755)
		}

		// Create manifest
		manifest := CreateEmptyOverlayManifest(overlay.name, overlay.author)
		manifest.Content.Systems = overlay.systems

		// Write manifest
		manifestPath := filepath.Join(overlayPath, "manifest.yml")
		if err := WriteOverlayManifest(manifest, manifestPath); err != nil {
			logging.LogDebug("Error writing sample overlay manifest: %v", err)
		}

		// Create empty preview file
		previewPath := filepath.Join(overlayPath, "preview.png")
		CreateEmptyFile(previewPath)
	}

	logging.LogDebug("Sample catalog contents created")
}