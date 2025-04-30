// internal/themes/export.go
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"thememanager/internal/logging"
	"thememanager/internal/system"
	"thememanager/internal/ui"
)

// CreateThemeExportDirectory creates a new theme directory with sequential naming
func CreateThemeExportDirectory() (string, error) {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to Exports directory
	exportsDir := filepath.Join(cwd, "Exports")

	// Ensure directory exists
	if err := os.MkdirAll(exportsDir, 0755); err != nil {
		return "", fmt.Errorf("error creating exports directory: %w", err)
	}

	// Generate sequential theme name
	themeNumber := 1
	var themePath string

	for {
		themeName := fmt.Sprintf("theme_%d%s", themeNumber, system.ThemeExtension)
		themePath = filepath.Join(exportsDir, themeName)

		if _, err := os.Stat(themePath); os.IsNotExist(err) {
			// Theme directory doesn't exist, we can use this name
			break
		}

		themeNumber++
	}

	// Create the theme directory and subdirectories
	subDirs := []string{
		system.ThemeWallpapersDir + "/SystemWallpapers",
		system.ThemeWallpapersDir + "/ListWallpapers",
		system.ThemeWallpapersDir + "/CollectionWallpapers",
		system.ThemeIconsDir + "/SystemIcons",
		system.ThemeIconsDir + "/ToolIcons",
		system.ThemeIconsDir + "/CollectionIcons",
		system.ThemeOverlaysDir,
		system.ThemeFontsDir,
	}

	for _, dir := range subDirs {
		path := filepath.Join(themePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", fmt.Errorf("error creating theme subdirectory %s: %w", dir, err)
		}
	}

	logging.LogDebug("Created theme export directory: %s", themePath)
	return themePath, nil
}

// ExportTheme exports the current theme settings
func ExportTheme() error {
	logging.LogDebug("Starting theme export")

	// Create theme directory
	themePath, err := CreateThemeExportDirectory()
	if err != nil {
		logging.LogDebug("Error creating theme directory: %v", err)
		return fmt.Errorf("error creating theme directory: %w", err)
	}

	// Get theme name from directory
	themeName := filepath.Base(themePath)
	themeName = strings.TrimSuffix(themeName, system.ThemeExtension)

	// Create manifest
	manifest := CreateEmptyThemeManifest(themeName, "ThemeManager")

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Export components

	// Export wallpapers and update manifest
	if err := exportWallpapers(themePath, manifest, systemPaths); err != nil {
		logging.LogDebug("Warning: Error exporting wallpapers: %v", err)
		// Continue with other components
	} else {
		manifest.Content.Backgrounds = true
	}

	// Export icons and update manifest
	if err := exportIcons(themePath, manifest, systemPaths); err != nil {
		logging.LogDebug("Warning: Error exporting icons: %v", err)
		// Continue with other components
	} else {
		manifest.Content.Icons = true
	}

	// Export fonts and update manifest
	if err := exportFonts(themePath, manifest); err != nil {
		logging.LogDebug("Warning: Error exporting fonts: %v", err)
		// Continue with other components
	} else {
		manifest.Content.Fonts = true
	}

	// Export accent settings and update manifest
	if err := exportAccentSettings(themePath, manifest); err != nil {
		logging.LogDebug("Warning: Error exporting accent settings: %v", err)
		// Continue with other components
	} else {
		manifest.Content.Accents = true
	}

	// Write manifest to file
	manifestPath := system.GetThemeManifestPath(themePath)
	if err := WriteThemeManifest(manifest, manifestPath); err != nil {
		logging.LogDebug("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	// Create a preview image
	createThemePreview(themePath, systemPaths)

	logging.LogDebug("Theme export completed successfully: %s", themePath)

	// Show success message to user
	ui.ShowMessage(fmt.Sprintf("Theme exported successfully: %s", themeName), "3")

	return nil
}

// Helper functions

// exportWallpapers exports wallpaper images to the theme
func exportWallpapers(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths) error {
	logging.LogDebug("Exporting wallpapers")

	// Export root background
	rootBgPath := filepath.Join(system.RootPath, "bg.png")
	if FileExists(rootBgPath) {
		destPath := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers", "Root.png")
		if err := CopyFile(rootBgPath, destPath); err != nil {
			logging.LogDebug("Warning: Could not copy root background: %v", err)
		} else {
			// Update system info in manifest for root
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
		}
	}

	// Export root media background
	rootMediaBgPath := filepath.Join(system.RootMediaPath, "bg.png")
	if FileExists(rootMediaBgPath) {
		destPath := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers", "Root-Media.png")
		if err := CopyFile(rootMediaBgPath, destPath); err != nil {
			logging.LogDebug("Warning: Could not copy root media background: %v", err)
		} else {
			// Update system info in manifest
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
		}
	}

	// Export Recently Played background
	rpBgPath := system.GetRecentlyPlayedBackgroundPath()
	if FileExists(rpBgPath) {
		destPath := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers", "Recently Played.png")
		if err := CopyFile(rpBgPath, destPath); err != nil {
			logging.LogDebug("Warning: Could not copy Recently Played background: %v", err)
		} else {
			// Update system info in manifest
			rpSystem := SystemConfig{
				DisplayName: "Recently Played",
				Files: map[string]string{
					"menu_bg": "Wallpapers/SystemWallpapers/Recently Played.png",
				},
				Paths: map[string]string{
					"menu_bg_path": rpBgPath,
				},
			}
			manifest.Systems["RecentlyPlayed"] = rpSystem
		}
	}

	// Export Tools background
	toolsBgPath := system.GetToolsBackgroundPath()
	if FileExists(toolsBgPath) {
		destPath := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers", "Tools.png")
		if err := CopyFile(toolsBgPath, destPath); err != nil {
			logging.LogDebug("Warning: Could not copy Tools background: %v", err)
		} else {
			// Update system info in manifest
			toolsSystem := SystemConfig{
				DisplayName: "Tools",
				Files: map[string]string{
					"menu_bg": "Wallpapers/SystemWallpapers/Tools.png",
				},
				Paths: map[string]string{
					"menu_bg_path": toolsBgPath,
				},
			}
			manifest.Systems["Tools"] = toolsSystem
		}
	}

	// Export Collections background
	collectionsBgPath := system.GetCollectionsBackgroundPath()
	if FileExists(collectionsBgPath) {
		destPath := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers", "Collections.png")
		if err := CopyFile(collectionsBgPath, destPath); err != nil {
			logging.LogDebug("Warning: Could not copy Collections background: %v", err)
		} else {
			// Update system info in manifest
			collectionsSystem := SystemConfig{
				DisplayName: "Collections",
				Files: map[string]string{
					"menu_bg": "Wallpapers/SystemWallpapers/Collections.png",
				},
				Paths: map[string]string{
					"menu_bg_path": collectionsBgPath,
				},
			}
			manifest.Systems["Collections"] = collectionsSystem
		}
	}

	// Export system backgrounds and list backgrounds
	for _, sysInfo := range systemPaths.Systems {
		if sysInfo.Tag == "" {
			continue // Skip systems without tags
		}

		// System background
		systemBgPath := system.GetSystemBackgroundPath(sysInfo.Name)
		if FileExists(systemBgPath) {
			fileName := fmt.Sprintf("%s (%s).png", sysInfo.Name, sysInfo.Tag)
			destPath := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers", fileName)

			if err := CopyFile(systemBgPath, destPath); err != nil {
				logging.LogDebug("Warning: Could not copy system background for %s: %v", sysInfo.Name, err)
			} else {
				// Update manifest with system info
				UpdateManifestWithSystemInfo(manifest, &sysInfo, themePath)
			}
		}

		// System list background
		systemListBgPath := system.GetSystemListBackgroundPath(sysInfo.Name)
		if FileExists(systemListBgPath) {
			fileName := fmt.Sprintf("%s-list (%s).png", sysInfo.Name, sysInfo.Tag)
			destPath := filepath.Join(themePath, system.ThemeWallpapersDir, "ListWallpapers", fileName)

			if err := CopyFile(systemListBgPath, destPath); err != nil {
				logging.LogDebug("Warning: Could not copy system list background for %s: %v", sysInfo.Name, err)
			} else {
				// Update system info in manifest if it exists
				if sysConfig, ok := manifest.Systems[sysInfo.Tag]; ok {
					sysConfig.Files["list_bg"] = fmt.Sprintf("Wallpapers/ListWallpapers/%s-list (%s).png", sysInfo.Name, sysInfo.Tag)
					sysConfig.Paths["list_bg_path"] = systemListBgPath
					manifest.Systems[sysInfo.Tag] = sysConfig
				}
			}
		}
	}

	// Export collection wallpapers
	collectionsDir := filepath.Join(system.CollectionsPath)
	if entries, err := os.ReadDir(collectionsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			collectionName := entry.Name()
			collectionBgPath := system.GetCollectionBackgroundPath(collectionName)

			if FileExists(collectionBgPath) {
				destPath := filepath.Join(themePath, system.ThemeWallpapersDir, "CollectionWallpapers", collectionName+".png")

				if err := CopyFile(collectionBgPath, destPath); err != nil {
					logging.LogDebug("Warning: Could not copy collection background for %s: %v", collectionName, err)
				} else {
					// Create or update collection in manifest
					collectionSystem := SystemConfig{
						DisplayName: collectionName,
						Files: map[string]string{
							"collection_bg": fmt.Sprintf("Wallpapers/CollectionWallpapers/%s.png", collectionName),
						},
						Paths: map[string]string{
							"collection_bg_path": collectionBgPath,
						},
					}

					colKey := "Collection_" + collectionName
					manifest.Systems[colKey] = collectionSystem
				}
			}
		}
	}

	return nil
}

// exportIcons exports icon images to the theme
func exportIcons(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths) error {
	logging.LogDebug("Exporting icons")

	// Export Recently Played icon
	rpIconPath := system.GetRecentlyPlayedIconPath()
	if FileExists(rpIconPath) {
		destPath := filepath.Join(themePath, system.ThemeIconsDir, "SystemIcons", "Recently Played.png")
		if err := CopyFile(rpIconPath, destPath); err != nil {
			logging.LogDebug("Warning: Could not copy Recently Played icon: %v", err)
		} else {
			// Update system info in manifest
			if rpSystem, ok := manifest.Systems["RecentlyPlayed"]; ok {
				rpSystem.Files["menu_icon"] = "Icons/SystemIcons/Recently Played.png"
				rpSystem.Paths["menu_icon_path"] = rpIconPath
				manifest.Systems["RecentlyPlayed"] = rpSystem
			}
		}
	}

	// Export Tools icon
	toolsIconPath := system.GetToolsIconPath()
	if FileExists(toolsIconPath) {
		destPath := filepath.Join(themePath, system.ThemeIconsDir, "SystemIcons", "Tools.png")
		if err := CopyFile(toolsIconPath, destPath); err != nil {
			logging.LogDebug("Warning: Could not copy Tools icon: %v", err)
		} else {
			// Update system info in manifest
			if toolsSystem, ok := manifest.Systems["Tools"]; ok {
				toolsSystem.Files["menu_icon"] = "Icons/SystemIcons/Tools.png"
				toolsSystem.Paths["menu_icon_path"] = toolsIconPath
				manifest.Systems["Tools"] = toolsSystem
			}
		}
	}

	// Export Collections icon
	collectionsIconPath := system.GetCollectionsIconPath()
	if FileExists(collectionsIconPath) {
		destPath := filepath.Join(themePath, system.ThemeIconsDir, "SystemIcons", "Collections.png")
		if err := CopyFile(collectionsIconPath, destPath); err != nil {
			logging.LogDebug("Warning: Could not copy Collections icon: %v", err)
		} else {
			// Update system info in manifest
			if colSystem, ok := manifest.Systems["Collections"]; ok {
				colSystem.Files["menu_icon"] = "Icons/SystemIcons/Collections.png"
				colSystem.Paths["menu_icon_path"] = collectionsIconPath
				manifest.Systems["Collections"] = colSystem
			}
		}
	}

	// Export system icons
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

			// Extract system tag
			iconName := entry.Name()
			systemName := strings.TrimSuffix(iconName, ".png")
			systemTag := system.ExtractSystemTag(systemName)

			if systemTag == "" {
				logging.LogDebug("Skipping system icon without tag: %s", iconName)
				continue
			}

			// Source path
			iconPath := filepath.Join(systemIconsDir, iconName)

			// Destination path
			destPath := filepath.Join(themePath, system.ThemeIconsDir, "SystemIcons", iconName)

			if err := CopyFile(iconPath, destPath); err != nil {
				logging.LogDebug("Warning: Could not copy system icon for %s: %v", systemName, err)
			} else {
				// Update system info in manifest if it exists
				if sysConfig, ok := manifest.Systems[systemTag]; ok {
					sysConfig.Files["menu_icon"] = fmt.Sprintf("Icons/SystemIcons/%s", iconName)
					sysConfig.Paths["menu_icon_path"] = iconPath
					manifest.Systems[systemTag] = sysConfig
				}
			}
		}
	}

	// Export tool icons
	toolsDir := system.ToolsPath
	if entries, err := os.ReadDir(toolsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			toolName := entry.Name()
			toolIconPath := system.GetToolIconPath(toolName)

			if FileExists(toolIconPath) {
				destPath := filepath.Join(themePath, system.ThemeIconsDir, "ToolIcons", toolName+".png")

				if err := CopyFile(toolIconPath, destPath); err != nil {
					logging.LogDebug("Warning: Could not copy tool icon for %s: %v", toolName, err)
				} else {
					// Create or update tool in manifest
					toolSystem := SystemConfig{
						DisplayName: toolName,
						Files: map[string]string{
							"tool_icon": fmt.Sprintf("Icons/ToolIcons/%s.png", toolName),
						},
						Paths: map[string]string{
							"tool_icon_path": toolIconPath,
						},
					}

					toolKey := "Tool_" + toolName
					manifest.Systems[toolKey] = toolSystem
				}
			}
		}
	}

	// Export collection icons
	collectionsDir := system.CollectionsPath
	if entries, err := os.ReadDir(collectionsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			collectionName := entry.Name()
			collectionIconPath := system.GetCollectionIconPath(collectionName)

			if FileExists(collectionIconPath) {
				destPath := filepath.Join(themePath, system.ThemeIconsDir, "CollectionIcons", collectionName+".png")

				if err := CopyFile(collectionIconPath, destPath); err != nil {
					logging.LogDebug("Warning: Could not copy collection icon for %s: %v", collectionName, err)
				} else {
					// Update collection in manifest if it exists
					colKey := "Collection_" + collectionName
					if colSystem, ok := manifest.Systems[colKey]; ok {
						colSystem.Files["collection_icon"] = fmt.Sprintf("Icons/CollectionIcons/%s.png", collectionName)
						colSystem.Paths["collection_icon_path"] = collectionIconPath
						manifest.Systems[colKey] = colSystem
					}
				}
			}
		}
	}

	return nil
}

// exportFonts exports fonts to the theme
func exportFonts(themePath string, manifest *ThemeManifest) error {
	logging.LogDebug("Exporting fonts")

	// Font paths to check
	fontPaths := map[string]string{
		"OG":          system.FontOGPath,
		"OG.backup":   system.FontOGBackupPath,
		"Next":        system.FontNextPath,
		"Next.backup": system.FontNextBackupPath,
	}

	fontDestDir := filepath.Join(themePath, system.ThemeFontsDir)
	if err := os.MkdirAll(fontDestDir, 0755); err != nil {
		return fmt.Errorf("error creating fonts directory: %w", err)
	}

	// Check and export each font
	for fontName, srcPath := range fontPaths {
		if FileExists(srcPath) {
			destPath := filepath.Join(fontDestDir, fontName+".ttf")

			if err := CopyFile(srcPath, destPath); err != nil {
				logging.LogDebug("Warning: Could not copy font %s: %v", fontName, err)
			} else {
				// Update font info in manifest
				fontSystem := SystemConfig{
					DisplayName: "Fonts",
					Files: map[string]string{
						fontName: fmt.Sprintf("Fonts/%s.ttf", fontName),
					},
					Paths: map[string]string{
						fontName + "_path": srcPath,
					},
				}

				// Add or update the font system entry
				if fs, ok := manifest.Systems["Fonts"]; ok {
					fs.Files[fontName] = fmt.Sprintf("Fonts/%s.ttf", fontName)
					fs.Paths[fontName+"_path"] = srcPath
					manifest.Systems["Fonts"] = fs
				} else {
					manifest.Systems["Fonts"] = fontSystem
				}
			}
		}
	}

	return nil
}

// exportAccentSettings exports accent settings to the theme
func exportAccentSettings(themePath string, manifest *ThemeManifest) error {
	logging.LogDebug("Exporting accent settings")

	// Accent settings path
	accentSettingsPath := system.AccentSettingsPath
	if !FileExists(accentSettingsPath) {
		logging.LogDebug("Accent settings file not found: %s", accentSettingsPath)
		return fmt.Errorf("accent settings file not found")
	}

	// Read settings file
	content, err := os.ReadFile(accentSettingsPath)
	if err != nil {
		return fmt.Errorf("error reading accent settings: %w", err)
	}

	// Create directory for settings
	settingsDir := filepath.Join(themePath, "Settings")
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		return fmt.Errorf("error creating settings directory: %w", err)
	}

	// Write settings to file
	destPath := filepath.Join(settingsDir, "minuisettings.txt")
	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("error writing accent settings: %w", err)
	}

	// Update settings info in manifest
	settingsSystem := SystemConfig{
		DisplayName: "Settings",
		Files: map[string]string{
			"accent_settings": "Settings/minuisettings.txt",
		},
		Paths: map[string]string{
			"accent_settings_path": accentSettingsPath,
		},
	}

	// Add or update the settings system entry
	if ss, ok := manifest.Systems["Settings"]; ok {
		ss.Files["accent_settings"] = "Settings/minuisettings.txt"
		ss.Paths["accent_settings_path"] = accentSettingsPath
		manifest.Systems["Settings"] = ss
	} else {
		manifest.Systems["Settings"] = settingsSystem
	}

	return nil
}

// createThemePreview creates a preview image for the theme
func createThemePreview(themePath string, systemPaths *system.SystemPaths) {
	// For now, just use the root background as the preview if it exists
	rootBgPath := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers", "Root.png")
	if FileExists(rootBgPath) {
		previewPath := filepath.Join(themePath, "preview.png")
		CopyFile(rootBgPath, previewPath)
		return
	}

	// Try root media background
	rootMediaBgPath := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers", "Root-Media.png")
	if FileExists(rootMediaBgPath) {
		previewPath := filepath.Join(themePath, "preview.png")
		CopyFile(rootMediaBgPath, previewPath)
		return
	}

	// Try recently played background
	rpBgPath := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers", "Recently Played.png")
	if FileExists(rpBgPath) {
		previewPath := filepath.Join(themePath, "preview.png")
		CopyFile(rpBgPath, previewPath)
		return
	}

	// Try first system background
	wallpapersDir := filepath.Join(themePath, system.ThemeWallpapersDir, "SystemWallpapers")
	if entries, err := os.ReadDir(wallpapersDir); err == nil && len(entries) > 0 {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".png") {
				previewPath := filepath.Join(themePath, "preview.png")
				CopyFile(filepath.Join(wallpapersDir, entry.Name()), previewPath)
				return
			}
		}
	}

	// Create an empty preview file if nothing else is available
	CreateEmptyFile(filepath.Join(themePath, "preview.png"))
}