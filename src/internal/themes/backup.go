// internal/themes/backup.go
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
    "time"
    "strings"
	"thememanager/internal/logging"
	"thememanager/internal/system"
)

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

	// Get backup path
	backupPath, err := GetThemeBackupPath(backupName)
	if err != nil {
		return fmt.Errorf("error getting backup path: %w", err)
	}

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Clean up existing theme components
	cleanBackgrounds(systemPaths)
	cleanIcons(systemPaths)

	// Read manifest if it exists
	var manifest *ThemeManifest
	manifestPath := filepath.Join(backupPath, "manifest.yml")
	if FileExists(manifestPath) {
		manifest, err = ReadThemeManifest(manifestPath)
		if err != nil {
			logging.LogDebug("Warning: Error reading backup manifest: %v", err)
			// Continue without manifest
			manifest = nil
		}
	}

	// If we have a manifest, restore using it
	if manifest != nil {
		for tag, sysConfig := range manifest.Systems {
			if err := applySystemFiles(backupPath, tag, sysConfig); err != nil {
				logging.LogDebug("Warning: Error restoring files for system %s: %v", tag, err)
			}
		}

		// Apply accent settings if present
		if manifest.Content.Accents {
			if err := applyAccentSettings(backupPath); err != nil {
				logging.LogDebug("Warning: Error restoring accent settings: %v", err)
			}
		}
	} else {
		// No manifest, restore the old way by copying files from backup structure
		if err := restoreThemeFilesOldWay(backupPath, systemPaths); err != nil {
			return fmt.Errorf("error restoring theme files: %w", err)
		}
	}

	logging.LogDebug("Theme reverted successfully from backup: %s", backupName)
	return nil
}

// RevertOverlayFromBackup restores overlay settings from a backup
func RevertOverlayFromBackup(backupName string) error {
	logging.LogDebug("Reverting overlays from backup: %s", backupName)

	// Get backup path
	backupPath, err := GetOverlayBackupPath(backupName)
	if err != nil {
		return fmt.Errorf("error getting backup path: %w", err)
	}

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Clean up existing overlays
	cleanOverlays(systemPaths)

	// Read manifest if it exists
	var manifest *OverlayManifest
	manifestPath := filepath.Join(backupPath, "manifest.yml")
	if FileExists(manifestPath) {
		manifest, err = ReadOverlayManifest(manifestPath)
		if err != nil {
			logging.LogDebug("Warning: Error reading backup manifest: %v", err)
			// Continue without manifest
			manifest = nil
		}
	}

    // If we have a manifest, do something with it
    if manifest != nil {
        // Use manifest data
        logging.LogDebug("Found manifest with %d systems", len(manifest.Content.Systems))
    }

	// Restore overlays
	overlaysBackupDir := filepath.Join(backupPath, "Overlays")
	if _, err := os.Stat(overlaysBackupDir); err == nil {
		// Get list of system tags from backup
		entries, err := os.ReadDir(overlaysBackupDir)
		if err != nil {
			return fmt.Errorf("error reading overlays backup directory: %w", err)
		}

		// Process each system
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			systemTag := entry.Name()
			srcDir := filepath.Join(overlaysBackupDir, systemTag)
			dstDir := system.GetOverlaySystemPath(systemTag)

			// Ensure destination directory exists
			if err := os.MkdirAll(dstDir, 0755); err != nil {
				logging.LogDebug("Warning: Error creating overlay directory for system %s: %v", systemTag, err)
				continue
			}

			// Copy overlay files
			if err := CopyDir(srcDir, dstDir); err != nil {
				logging.LogDebug("Warning: Error copying overlay files for system %s: %v", systemTag, err)
			}
		}
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

	// Clean up system files
	cleanBackgrounds(systemPaths)
	cleanIcons(systemPaths)
	cleanOverlays(systemPaths)

	// Clean up theme fonts
	cleanFonts()

	// Clean up accent settings
	cleanAccentSettings()

	// Delete all installed themes
	themesDir := filepath.Join(cwd, "Themes")
	if err := CleanDirectory(themesDir); err != nil {
		logging.LogDebug("Warning: Error cleaning themes directory: %v", err)
	}

	// Delete all installed overlays
	overlaysDir := filepath.Join(cwd, "Overlays")
	if err := CleanDirectory(overlaysDir); err != nil {
		logging.LogDebug("Warning: Error cleaning overlays directory: %v", err)
	}

	// Delete all backups
	backupsThemesDir := filepath.Join(cwd, "Backups", "Themes")
	if err := CleanDirectory(backupsThemesDir); err != nil {
		logging.LogDebug("Warning: Error cleaning theme backups directory: %v", err)
	}

	backupsOverlaysDir := filepath.Join(cwd, "Backups", "Overlays")
	if err := CleanDirectory(backupsOverlaysDir); err != nil {
		logging.LogDebug("Warning: Error cleaning overlay backups directory: %v", err)
	}

	logging.LogDebug("Purge completed successfully")
	return nil
}

// Helper functions

// backupCurrentTheme backs up the current theme settings
func backupCurrentTheme(backupPath string, systemPaths *system.SystemPaths) error {
	logging.LogDebug("Backing up current theme settings to: %s", backupPath)

	// Create required directories
	os.MkdirAll(filepath.Join(backupPath, system.ThemeWallpapersDir, "SystemWallpapers"), 0755)
	os.MkdirAll(filepath.Join(backupPath, system.ThemeWallpapersDir, "ListWallpapers"), 0755)
	os.MkdirAll(filepath.Join(backupPath, system.ThemeWallpapersDir, "CollectionWallpapers"), 0755)
	os.MkdirAll(filepath.Join(backupPath, system.ThemeIconsDir, "SystemIcons"), 0755)
	os.MkdirAll(filepath.Join(backupPath, system.ThemeIconsDir, "ToolIcons"), 0755)
	os.MkdirAll(filepath.Join(backupPath, system.ThemeIconsDir, "CollectionIcons"), 0755)
	os.MkdirAll(filepath.Join(backupPath, system.ThemeFontsDir), 0755)
	os.MkdirAll(filepath.Join(backupPath, "Settings"), 0755)

	// Create manifest
	manifest := CreateEmptyThemeManifest("Backup", "ThemeManager")

	// Backup root background
	rootBgPath := filepath.Join(system.RootPath, "bg.png")
	if FileExists(rootBgPath) {
		destPath := filepath.Join(backupPath, system.ThemeWallpapersDir, "SystemWallpapers", "Root.png")
		CopyFile(rootBgPath, destPath)

		// Update manifest
		rootSystem := SystemConfig{
			DisplayName: "Root",
			Files: map[string]string{
				"menu_bg": "Wallpapers/SystemWallpapers/Root.png",
			},
			Paths: map[string]string{
				"menu_bg_path": rootBgPath,
			},
		}
		manifest.Systems["Root"] = rootSystem
		manifest.Content.Backgrounds = true
	}

	// Backup root media background
	rootMediaBgPath := filepath.Join(system.RootMediaPath, "bg.png")
	if FileExists(rootMediaBgPath) {
		destPath := filepath.Join(backupPath, system.ThemeWallpapersDir, "SystemWallpapers", "Root-Media.png")
		CopyFile(rootMediaBgPath, destPath)

		// Update manifest
		if rootSystem, ok := manifest.Systems["Root"]; ok {
			rootSystem.Files["media_bg"] = "Wallpapers/SystemWallpapers/Root-Media.png"
			rootSystem.Paths["media_bg_path"] = rootMediaBgPath
			manifest.Systems["Root"] = rootSystem
		} else {
			rootSystem := SystemConfig{
				DisplayName: "Root",
				Files: map[string]string{
					"media_bg": "Wallpapers/SystemWallpapers/Root-Media.png",
				},
				Paths: map[string]string{
					"media_bg_path": rootMediaBgPath,
				},
			}
			manifest.Systems["Root"] = rootSystem
		}
		manifest.Content.Backgrounds = true
	}

	// Back up system backgrounds and icons
	for _, sysInfo := range systemPaths.Systems {
		if sysInfo.Tag == "" {
			continue // Skip systems without tags
		}

		// Check for background image
		bgPath := system.GetSystemBackgroundPath(sysInfo.Name)
		listBgPath := system.GetSystemListBackgroundPath(sysInfo.Name)
		iconPath := system.GetSystemIconPath(sysInfo.Name, sysInfo.Tag)

		// Create system config
		sysConfig := SystemConfig{
			DisplayName: sysInfo.Name,
			Files:       make(map[string]string),
			Paths:       make(map[string]string),
		}

		// Background
		if FileExists(bgPath) {
			destPath := filepath.Join(backupPath, system.ThemeWallpapersDir, "SystemWallpapers",
				fmt.Sprintf("%s (%s).png", sysInfo.Name, sysInfo.Tag))
			CopyFile(bgPath, destPath)

			sysConfig.Files["menu_bg"] = fmt.Sprintf("Wallpapers/SystemWallpapers/%s (%s).png",
				sysInfo.Name, sysInfo.Tag)
			sysConfig.Paths["menu_bg_path"] = bgPath
			manifest.Content.Backgrounds = true
		}

		// List background
		if FileExists(listBgPath) {
			destPath := filepath.Join(backupPath, system.ThemeWallpapersDir, "ListWallpapers",
				fmt.Sprintf("%s-list (%s).png", sysInfo.Name, sysInfo.Tag))
			CopyFile(listBgPath, destPath)

			sysConfig.Files["list_bg"] = fmt.Sprintf("Wallpapers/ListWallpapers/%s-list (%s).png",
				sysInfo.Name, sysInfo.Tag)
			sysConfig.Paths["list_bg_path"] = listBgPath
			manifest.Content.Backgrounds = true
		}

		// Icon
		if FileExists(iconPath) {
			destPath := filepath.Join(backupPath, system.ThemeIconsDir, "SystemIcons",
				fmt.Sprintf("%s (%s).png", sysInfo.Name, sysInfo.Tag))
			CopyFile(iconPath, destPath)

			sysConfig.Files["menu_icon"] = fmt.Sprintf("Icons/SystemIcons/%s (%s).png",
				sysInfo.Name, sysInfo.Tag)
			sysConfig.Paths["menu_icon_path"] = iconPath
			manifest.Content.Icons = true
		}

		// Add to manifest if we have any files
		if len(sysConfig.Files) > 0 {
			manifest.Systems[sysInfo.Tag] = sysConfig
		}
	}

	// Back up special icons
	// Recently Played
	rpIconPath := system.GetRecentlyPlayedIconPath()
	if FileExists(rpIconPath) {
		destPath := filepath.Join(backupPath, system.ThemeIconsDir, "SystemIcons", "Recently Played.png")
		CopyFile(rpIconPath, destPath)

		// Update manifest
		rpSystem := SystemConfig{
			DisplayName: "Recently Played",
			Files: map[string]string{
				"menu_icon": "Icons/SystemIcons/Recently Played.png",
			},
			Paths: map[string]string{
				"menu_icon_path": rpIconPath,
			},
		}
		manifest.Systems["RecentlyPlayed"] = rpSystem
		manifest.Content.Icons = true
	}

	// Tools
	toolsIconPath := system.GetToolsIconPath()
	if FileExists(toolsIconPath) {
		destPath := filepath.Join(backupPath, system.ThemeIconsDir, "SystemIcons", "Tools.png")
		CopyFile(toolsIconPath, destPath)

		// Update manifest
		toolsSystem := SystemConfig{
			DisplayName: "Tools",
			Files: map[string]string{
				"menu_icon": "Icons/SystemIcons/Tools.png",
			},
			Paths: map[string]string{
				"menu_icon_path": toolsIconPath,
			},
		}
		manifest.Systems["Tools"] = toolsSystem
		manifest.Content.Icons = true
	}

	// Collections
	collectionsIconPath := system.GetCollectionsIconPath()
	if FileExists(collectionsIconPath) {
		destPath := filepath.Join(backupPath, system.ThemeIconsDir, "SystemIcons", "Collections.png")
		CopyFile(collectionsIconPath, destPath)

		// Update manifest
		colSystem := SystemConfig{
			DisplayName: "Collections",
			Files: map[string]string{
				"menu_icon": "Icons/SystemIcons/Collections.png",
			},
			Paths: map[string]string{
				"menu_icon_path": collectionsIconPath,
			},
		}
		manifest.Systems["Collections"] = colSystem
		manifest.Content.Icons = true
	}

	// Back up fonts
	fontPaths := map[string]string{
		"OG":          system.FontOGPath,
		"OG.backup":   system.FontOGBackupPath,
		"Next":        system.FontNextPath,
		"Next.backup": system.FontNextBackupPath,
	}

	fontSystem := SystemConfig{
		DisplayName: "Fonts",
		Files:       make(map[string]string),
		Paths:       make(map[string]string),
	}

	for fontName, srcPath := range fontPaths {
		if FileExists(srcPath) {
			destPath := filepath.Join(backupPath, system.ThemeFontsDir, fontName+".ttf")
			CopyFile(srcPath, destPath)

			fontSystem.Files[fontName] = fmt.Sprintf("Fonts/%s.ttf", fontName)
			fontSystem.Paths[fontName+"_path"] = srcPath
			manifest.Content.Fonts = true
		}
	}

	if len(fontSystem.Files) > 0 {
		manifest.Systems["Fonts"] = fontSystem
	}

	// Back up accent settings
	accentSettingsPath := system.AccentSettingsPath
	if FileExists(accentSettingsPath) {
		destPath := filepath.Join(backupPath, "Settings", "minuisettings.txt")
		CopyFile(accentSettingsPath, destPath)

		// Update manifest
		settingsSystem := SystemConfig{
			DisplayName: "Settings",
			Files: map[string]string{
				"accent_settings": "Settings/minuisettings.txt",
			},
			Paths: map[string]string{
				"accent_settings_path": accentSettingsPath,
			},
		}
		manifest.Systems["Settings"] = settingsSystem
		manifest.Content.Accents = true
	}

	// Write manifest
	manifestPath := filepath.Join(backupPath, "manifest.yml")
	err := WriteThemeManifest(manifest, manifestPath)
	if err != nil {
		logging.LogDebug("Warning: Error writing backup manifest: %v", err)
	}

	// Create backup screenshot
	createBackupScreenshot(backupPath)

	return nil
}

// backupCurrentOverlays backs up the current overlay settings
func backupCurrentOverlays(backupPath string, systemPaths *system.SystemPaths) error {
	logging.LogDebug("Backing up current overlay settings to: %s", backupPath)

	// Create overlays directory
	overlaysBackupDir := filepath.Join(backupPath, "Overlays")
	os.MkdirAll(overlaysBackupDir, 0755)

	// Create manifest
	manifest := CreateEmptyOverlayManifest("Backup", "ThemeManager")
	var systemTags []string

	// Get overlays directory
	overlaysDir := system.OverlaysPath

	// Check if it exists
	if _, err := os.Stat(overlaysDir); os.IsNotExist(err) {
		logging.LogDebug("No overlays directory found")
		return nil
	}

	// Read directory
	entries, err := os.ReadDir(overlaysDir)
	if err != nil {
		return fmt.Errorf("error reading overlays directory: %w", err)
	}

	// Backup each system's overlays
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		systemTag := entry.Name()
		systemDir := filepath.Join(overlaysDir, systemTag)
		backupSystemDir := filepath.Join(overlaysBackupDir, systemTag)

		// Create backup directory for this system
		os.MkdirAll(backupSystemDir, 0755)

		// Check if there are any overlay files
		files, err := os.ReadDir(systemDir)
		if err != nil || len(files) == 0 {
			continue
		}

		// Copy overlay files
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".png") {
				continue
			}

			srcPath := filepath.Join(systemDir, file.Name())
			dstPath := filepath.Join(backupSystemDir, file.Name())

			if err := CopyFile(srcPath, dstPath); err != nil {
				logging.LogDebug("Warning: Could not backup overlay file %s: %v", file.Name(), err)
				continue
			}
		}

		// Add to system tags
		systemTags = append(systemTags, systemTag)
	}

	// Update manifest
	manifest.Content.Systems = systemTags

	// Write manifest
	manifestPath := filepath.Join(backupPath, "manifest.yml")
	err = WriteOverlayManifest(manifest, manifestPath)
	if err != nil {
		logging.LogDebug("Warning: Error writing backup manifest: %v", err)
	}

	// Create backup screenshot
	createBackupScreenshot(backupPath)

	return nil
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

// restoreThemeFilesOldWay restores theme files from a backup without using a manifest
func restoreThemeFilesOldWay(backupPath string, systemPaths *system.SystemPaths) error {
	logging.LogDebug("Restoring theme files using old method")

	// TODO: Implement old way of restoring theme files if needed
	// This would restore based on the directory structure rather than the manifest

	return nil
}

// cleanFonts removes custom fonts and restores defaults
func cleanFonts() error {
	logging.LogDebug("Cleaning up custom fonts")

	// TODO: Implement font restoration logic
	// This would restore system default fonts

	return nil
}

// cleanAccentSettings restores default accent settings
func cleanAccentSettings() error {
	logging.LogDebug("Cleaning up accent settings")

	// TODO: Implement accent settings restoration logic
	// This would restore system default accent settings

	return nil
}

// createBackupScreenshot creates a screenshot for the backup
func createBackupScreenshot(backupPath string) {
	// In a real implementation, this would take a screenshot of the current state
	// For now, create an empty file
	CreateEmptyFile(filepath.Join(backupPath, "screenshot.png"))
}