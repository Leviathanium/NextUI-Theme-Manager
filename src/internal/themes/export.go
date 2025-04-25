// src/internal/themes/export.go
// Simplified implementation of theme export functionality

package themes

import (
	"fmt"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"nextui-themes/internal/ui"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// CreateThemeExportDirectory creates a new theme directory with sequential naming
func CreateThemeExportDirectory() (string, error) {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to Exports directory directly instead of Themes/Exports
	exportsDir := filepath.Join(cwd, "Exports")

	// Ensure directory exists
	if err := os.MkdirAll(exportsDir, 0755); err != nil {
		return "", fmt.Errorf("error creating exports directory: %w", err)
	}

	// Generate sequential theme name
	themeNumber := 1
	var themePath string

	for {
		themeName := fmt.Sprintf("theme_%d.theme", themeNumber)
		themePath = filepath.Join(exportsDir, themeName)

		if _, err := os.Stat(themePath); os.IsNotExist(err) {
			// Theme directory doesn't exist, we can use this name
			break
		}

		themeNumber++
	}

	// Create the theme directory and subdirectories
	subDirs := []string{
		"Wallpapers/SystemWallpapers",
		"Wallpapers/CollectionWallpapers",
		"Icons/SystemIcons",
		"Icons/ToolIcons",
		"Icons/CollectionIcons",
		"Overlays",
		"Fonts",
		// Removed "Settings" directory since we're storing settings directly in manifest.json
	}

	for _, dir := range subDirs {
		path := filepath.Join(themePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", fmt.Errorf("error creating theme subdirectory %s: %w", dir, err)
		}
	}

	return themePath, nil
}

// ExportTheme exports the current theme settings
func ExportTheme() error {
	// Create logger
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting theme export")

	// Create theme directory
	themePath, err := CreateThemeExportDirectory()
	if err != nil {
		logger.DebugFn("Error creating theme directory: %v", err)
		return fmt.Errorf("error creating theme directory: %w", err)
	}

	logger.DebugFn("Created theme directory: %s", themePath)

	// Get theme name from directory
	themeName := filepath.Base(themePath)

	// Try to determine author from global manifest if available
	author := "AuthorName" // Default
	globalManifest, err := LoadGlobalManifest()
	if err == nil && globalManifest != nil {
		// Try to get author from current theme if it exists
		if globalManifest.CurrentTheme != "" {
			// Try to load the theme to get author
			currThemePath := filepath.Join(filepath.Dir(themePath), "..", "Themes", globalManifest.CurrentTheme)
			currManifest, err := ValidateTheme(currThemePath, logger)
			if err == nil && currManifest.ThemeInfo.Author != "" {
				author = currManifest.ThemeInfo.Author
			}
		}
	}

	// Initialize minimal manifest
	manifest := CreateMinimalThemeManifest(themeName, author)

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logger.DebugFn("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Copy the actual files but don't add to manifest content or path_mappings

	// Export wallpapers
	exportWallpapers(themePath, manifest, systemPaths, logger)

	// Export icons
	exportIcons(themePath, manifest, systemPaths, logger)

	// Export overlays
	exportOverlays(themePath, manifest, systemPaths, logger)

	// Export fonts
	exportFonts(themePath, manifest, logger)

	// Read and include accent settings directly in manifest
	if err := readAccentSettingsFromSystem(manifest, logger); err != nil {
		logger.DebugFn("Warning: Could not read accent settings: %v", err)
	}

	// Read and include LED settings directly in manifest
	if err := readLEDSettingsFromSystem(manifest, logger); err != nil {
		logger.DebugFn("Warning: Could not read LED settings: %v", err)
	}

	// Write manifest
	if err := WriteManifest(themePath, manifest, logger); err != nil {
		logger.DebugFn("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	logger.DebugFn("Theme export completed successfully: %s", themePath)

	// Show success message to user
	themeName = filepath.Base(themePath)
	ui.ShowMessage(fmt.Sprintf("Theme exported successfully: %s", themeName), "3")

	return nil
}

// exportWallpapers scans for and exports wallpapers
func exportWallpapers(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) {
	// Initialize wallpaper section
	manifest.Content.Wallpapers.Present = false
	manifest.Content.Wallpapers.Count = 0
	manifest.PathMappings.Wallpapers = []PathMapping{}

	// Check for root wallpaper
	rootBg := filepath.Join(systemPaths.Root, "bg.png")
	if _, err := os.Stat(rootBg); err == nil {
		// Copy to theme
		destPath := filepath.Join(themePath, "Wallpapers", "SystemWallpapers", "Root.png")
		if err := CopyFile(rootBg, destPath); err != nil {
			logger.DebugFn("Warning: Could not copy root bg.png: %v", err)
		} else {
			// Add to manifest
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Root.png",
					SystemPath: rootBg,
					Metadata: map[string]string{
						"SystemName":    "Root",
						"WallpaperType": "Main",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.DebugFn("Exported Root wallpaper to %s", destPath)
		}
	}

	// Check for root media wallpaper
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", "bg.png")
	if _, err := os.Stat(rootMediaBg); err == nil {
		// Copy to theme
		destPath := filepath.Join(themePath, "Wallpapers", "SystemWallpapers", "Root-Media.png")
		if err := CopyFile(rootMediaBg, destPath); err != nil {
			logger.DebugFn("Warning: Could not copy root media bg.png: %v", err)
		} else {
			// Add to manifest
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Root-Media.png",
					SystemPath: rootMediaBg,
					Metadata: map[string]string{
						"SystemName":    "Root",
						"WallpaperType": "Media",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.DebugFn("Exported Root-Media wallpaper to %s", destPath)
		}
	}

	// Check for Recently Played wallpaper
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
	if _, err := os.Stat(rpBg); err == nil {
		// Copy to theme
		destPath := filepath.Join(themePath, "Wallpapers", "SystemWallpapers", "Recently Played.png")
		if err := CopyFile(rpBg, destPath); err != nil {
			logger.DebugFn("Warning: Could not copy Recently Played bg.png: %v", err)
		} else {
			// Add to manifest
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Recently Played.png",
					SystemPath: rpBg,
					Metadata: map[string]string{
						"SystemName":    "Recently Played",
						"WallpaperType": "Media",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.DebugFn("Exported Recently Played wallpaper to %s", destPath)
		}
	}

	// Check for Tools wallpaper
	toolsBg := filepath.Join(systemPaths.Tools, ".media", "bg.png")
	if _, err := os.Stat(toolsBg); err == nil {
		// Copy to theme
		destPath := filepath.Join(themePath, "Wallpapers", "SystemWallpapers", "Tools.png")
		if err := CopyFile(toolsBg, destPath); err != nil {
			logger.DebugFn("Warning: Could not copy Tools bg.png: %v", err)
		} else {
			// Add to manifest
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Tools.png",
					SystemPath: toolsBg,
					Metadata: map[string]string{
						"SystemName":    "Tools",
						"WallpaperType": "Media",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.DebugFn("Exported Tools wallpaper to %s", destPath)
		}
	}

	// Check for Collections wallpaper
	collectionsBg := filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")
	if _, err := os.Stat(collectionsBg); err == nil {
		// Copy to theme
		destPath := filepath.Join(themePath, "Wallpapers", "SystemWallpapers", "Collections.png")
		if err := CopyFile(collectionsBg, destPath); err != nil {
			logger.DebugFn("Warning: Could not copy Collections bg.png: %v", err)
		} else {
			// Add to manifest
			manifest.PathMappings.Wallpapers = append(
				manifest.PathMappings.Wallpapers,
				PathMapping{
					ThemePath:  "Wallpapers/SystemWallpapers/Collections.png",
					SystemPath: collectionsBg,
					Metadata: map[string]string{
						"SystemName":    "Collections",
						"WallpaperType": "Media",
					},
				},
			)
			manifest.Content.Wallpapers.Present = true
			manifest.Content.Wallpapers.Count++
			logger.DebugFn("Exported Collections wallpaper to %s", destPath)
		}
	}

	// Check for system wallpapers
	for _, system := range systemPaths.Systems {
		if system.Tag == "" {
			// Skip systems without tags
			continue
		}

		systemBg := filepath.Join(system.MediaPath, "bg.png")
		if _, err := os.Stat(systemBg); err == nil {
			// Create filename with system tag - check if already contains tag to prevent duplication
			var fileName string
			if strings.Contains(system.Name, fmt.Sprintf("(%s)", system.Tag)) {
				// System name already has the tag, use as is
				fileName = fmt.Sprintf("%s.png", system.Name)
				logger.DebugFn("System name already contains tag: %s", system.Name)
			} else {
				// Add tag to system name
				fileName = fmt.Sprintf("%s (%s).png", system.Name, system.Tag)
				logger.DebugFn("Adding tag to system name: %s (%s)", system.Name, system.Tag)
			}

			destPath := filepath.Join(themePath, "Wallpapers", "SystemWallpapers", fileName)

			if err := CopyFile(systemBg, destPath); err != nil {
				logger.DebugFn("Warning: Could not copy system %s bg.png: %v", system.Name, err)
			} else {
				// Add to manifest
				manifest.PathMappings.Wallpapers = append(
					manifest.PathMappings.Wallpapers,
					PathMapping{
						ThemePath:  "Wallpapers/SystemWallpapers/" + fileName,
						SystemPath: systemBg,
						Metadata: map[string]string{
							"SystemName":    system.Name,
							"SystemTag":     system.Tag,
							"WallpaperType": "System",
						},
					},
				)
				manifest.Content.Wallpapers.Present = true
				manifest.Content.Wallpapers.Count++
				logger.DebugFn("Exported %s wallpaper to %s", system.Name, destPath)
			}
		}
	}

	// Check for collection wallpapers
	collectionsDir := filepath.Join(systemPaths.Root, "Collections")
	entries, err := os.ReadDir(collectionsDir)
	if err != nil {
		logger.DebugFn("Warning: Could not read Collections directory: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		collectionName := entry.Name()
		collectionBg := filepath.Join(collectionsDir, collectionName, ".media", "bg.png")

		if _, err := os.Stat(collectionBg); err == nil {
			// Create filename for collection
			fileName := fmt.Sprintf("%s.png", collectionName)
			destPath := filepath.Join(themePath, "Wallpapers", "CollectionWallpapers", fileName)

			if err := CopyFile(collectionBg, destPath); err != nil {
				logger.DebugFn("Warning: Could not copy collection %s bg.png: %v", collectionName, err)
			} else {
				// Add to manifest
				manifest.PathMappings.Wallpapers = append(
					manifest.PathMappings.Wallpapers,
					PathMapping{
						ThemePath:  "Wallpapers/CollectionWallpapers/" + fileName,
						SystemPath: collectionBg,
						Metadata: map[string]string{
							"CollectionName": collectionName,
							"WallpaperType":  "Collection",
						},
					},
				)
				manifest.Content.Wallpapers.Present = true
				manifest.Content.Wallpapers.Count++
				logger.DebugFn("Exported collection %s wallpaper to %s", collectionName, destPath)
			}
		}
	}
}

// exportIcons scans for and exports icons
func exportIcons(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) {
	// Initialize icon section
	manifest.Content.Icons.Present = false
	manifest.Content.Icons.SystemCount = 0
	manifest.Content.Icons.ToolCount = 0
	manifest.Content.Icons.CollectionCount = 0
	manifest.PathMappings.Icons = []PathMapping{}

	// Export system icons

	// Recently Played icon - in SD_CARD/.media/Recently Played.png
	rpIcon := filepath.Join(systemPaths.Root, ".media", "Recently Played.png")
	if _, err := os.Stat(rpIcon); err == nil {
		destPath := filepath.Join(themePath, "Icons", "SystemIcons", "Recently Played.png")
		if err := CopyFile(rpIcon, destPath); err != nil {
			logger.DebugFn("Warning: Could not copy Recently Played icon: %v", err)
		} else {
			manifest.PathMappings.Icons = append(
				manifest.PathMappings.Icons,
				PathMapping{
					ThemePath:  "Icons/SystemIcons/Recently Played.png",
					SystemPath: rpIcon,
					Metadata: map[string]string{
						"SystemName": "Recently Played",
						"IconType":   "System",
					},
				},
			)
			manifest.Content.Icons.Present = true
			manifest.Content.Icons.SystemCount++
			logger.DebugFn("Exported Recently Played icon to %s", destPath)
		}
	}

	// Tools icon - use parent path of Tools since Tools path includes tg5040
	toolsParentDir := filepath.Dir(systemPaths.Tools) // Gets /mnt/SDCARD/Tools
	toolsIcon := filepath.Join(toolsParentDir, ".media", "tg5040.png")
	if _, err := os.Stat(toolsIcon); err == nil {
		destPath := filepath.Join(themePath, "Icons", "SystemIcons", "Tools.png")
		if err := CopyFile(toolsIcon, destPath); err != nil {
			logger.DebugFn("Warning: Could not copy Tools icon: %v", err)
		} else {
			manifest.PathMappings.Icons = append(
				manifest.PathMappings.Icons,
				PathMapping{
					ThemePath:  "Icons/SystemIcons/Tools.png",
					SystemPath: toolsIcon,
					Metadata: map[string]string{
						"SystemName": "Tools",
						"IconType":   "System",
					},
				},
			)
			manifest.Content.Icons.Present = true
			manifest.Content.Icons.SystemCount++
			logger.DebugFn("Exported Tools icon to %s", destPath)
		}
	}

	// Collections icon - in SD_CARD/.media/Collections.png
	collectionsIcon := filepath.Join(systemPaths.Root, ".media", "Collections.png")
	if _, err := os.Stat(collectionsIcon); err == nil {
		destPath := filepath.Join(themePath, "Icons", "SystemIcons", "Collections.png")
		if err := CopyFile(collectionsIcon, destPath); err != nil {
			logger.DebugFn("Warning: Could not copy Collections icon: %v", err)
		} else {
			manifest.PathMappings.Icons = append(
				manifest.PathMappings.Icons,
				PathMapping{
					ThemePath:  "Icons/SystemIcons/Collections.png",
					SystemPath: collectionsIcon,
					Metadata: map[string]string{
						"SystemName": "Collections",
						"IconType":   "System",
					},
				},
			)
			manifest.Content.Icons.Present = true
			manifest.Content.Icons.SystemCount++
			logger.DebugFn("Exported Collections icon to %s", destPath)
		}
	}

	// System-specific icons - each system has its own icon file in Roms/.media/ with system name and tag
	systemIconsDir := filepath.Join(systemPaths.Roms, ".media")
	if _, err := os.Stat(systemIconsDir); err == nil {
		entries, err := os.ReadDir(systemIconsDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				// Only process icons that match system naming pattern
				// Skip other special icons like Recently Played that we handle separately
				if entry.Name() == "Recently Played.png" ||
					entry.Name() == "Collections.png" ||
					entry.Name() == "tg5040.png" {
					continue
				}

				// Check for system tag pattern
				tagRegex := regexp.MustCompile(`\((.*?)\)`)
				if !tagRegex.MatchString(entry.Name()) {
					logger.DebugFn("Skipping non-system icon: %s", entry.Name())
					continue
				}

				systemIconPath := filepath.Join(systemIconsDir, entry.Name())
				destPath := filepath.Join(themePath, "Icons", "SystemIcons", entry.Name())

				if err := CopyFile(systemIconPath, destPath); err != nil {
					logger.DebugFn("Warning: Could not copy system icon %s: %v", entry.Name(), err)
				} else {
					// Extract system tag for metadata
					matches := tagRegex.FindStringSubmatch(entry.Name())
					systemTag := ""
					if len(matches) >= 2 {
						systemTag = matches[1]
					}

					manifest.PathMappings.Icons = append(
						manifest.PathMappings.Icons,
						PathMapping{
							ThemePath:  "Icons/SystemIcons/" + entry.Name(),
							SystemPath: systemIconPath,
							Metadata: map[string]string{
								"SystemName": strings.TrimSuffix(entry.Name(), ".png"),
								"SystemTag":  systemTag,
								"IconType":   "System",
							},
						},
					)
					manifest.Content.Icons.Present = true
					manifest.Content.Icons.SystemCount++
					logger.DebugFn("Exported system icon %s to %s", entry.Name(), destPath)
				}
			}
		}
	}

	// Tool icons - each tool folder has its own icon.png file
	toolsDir := filepath.Join(systemPaths.Tools)
	toolEntries, err := os.ReadDir(toolsDir)
	if err == nil {
		for _, entry := range toolEntries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			toolName := entry.Name()
			toolIcon := filepath.Join(toolsDir, toolName, ".media", toolName+".png")

			if _, err := os.Stat(toolIcon); err == nil {
				destPath := filepath.Join(themePath, "Icons", "ToolIcons", fmt.Sprintf("%s.png", toolName))

				if err := CopyFile(toolIcon, destPath); err != nil {
					logger.DebugFn("Warning: Could not copy tool %s icon: %v", toolName, err)
				} else {
					manifest.PathMappings.Icons = append(
						manifest.PathMappings.Icons,
						PathMapping{
							ThemePath:  fmt.Sprintf("Icons/ToolIcons/%s.png", toolName),
							SystemPath: toolIcon,
							Metadata: map[string]string{
								"ToolName": toolName,
								"IconType": "Tool",
							},
						},
					)
					manifest.Content.Icons.Present = true
					manifest.Content.Icons.ToolCount++
					logger.DebugFn("Exported tool %s icon to %s", toolName, destPath)
				}
			}
		}
	}

	// Collection icons - each collection has its own icon.png file
	collectionsDir := filepath.Join(systemPaths.Root, "Collections")
	colEntries, err := os.ReadDir(collectionsDir)
	if err == nil {
		for _, entry := range colEntries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			collectionName := entry.Name()
			collectionIcon := filepath.Join(collectionsDir, collectionName, ".media", collectionName+".png")

			if _, err := os.Stat(collectionIcon); err == nil {
				destPath := filepath.Join(themePath, "Icons", "CollectionIcons", fmt.Sprintf("%s.png", collectionName))

				if err := CopyFile(collectionIcon, destPath); err != nil {
					logger.DebugFn("Warning: Could not copy collection %s icon: %v", collectionName, err)
				} else {
					manifest.PathMappings.Icons = append(
						manifest.PathMappings.Icons,
						PathMapping{
							ThemePath:  fmt.Sprintf("Icons/CollectionIcons/%s.png", collectionName),
							SystemPath: collectionIcon,
							Metadata: map[string]string{
								"CollectionName": collectionName,
								"IconType":       "Collection",
							},
						},
					)
					manifest.Content.Icons.Present = true
					manifest.Content.Icons.CollectionCount++
					logger.DebugFn("Exported collection %s icon to %s", collectionName, destPath)
				}
			}
		}
	}
}

// exportOverlays scans for and exports system overlays
func exportOverlays(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) {
	// Initialize overlay section
	manifest.Content.Overlays.Present = false
	manifest.Content.Overlays.Systems = []string{}
	manifest.PathMappings.Overlays = []PathMapping{}

	// Check for overlays directory
	overlaysDir := filepath.Join(systemPaths.Root, "Overlays")
	if _, err := os.Stat(overlaysDir); os.IsNotExist(err) {
		logger.DebugFn("Overlays directory not found: %s", overlaysDir)
		return
	}

	// List system directories in Overlays
	entries, err := os.ReadDir(overlaysDir)
	if err != nil {
		logger.DebugFn("Error reading Overlays directory: %v", err)
		return
	}

	// Process each system's overlays
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		systemTag := entry.Name()
		systemOverlaysPath := filepath.Join(overlaysDir, systemTag)

		// List overlay files for this system
		overlayFiles, err := os.ReadDir(systemOverlaysPath)
		if err != nil {
			logger.DebugFn("Error reading system overlays directory %s: %v", systemTag, err)
			continue
		}

		var hasOverlays bool

		// Copy each overlay file
		for _, file := range overlayFiles {
			if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
				continue
			}

			// Only process PNG files
			if !strings.HasSuffix(strings.ToLower(file.Name()), ".png") {
				continue
			}

			srcPath := filepath.Join(systemOverlaysPath, file.Name())

			// Create destination directories if needed
			destDir := filepath.Join(themePath, "Overlays", systemTag)
			if err := os.MkdirAll(destDir, 0755); err != nil {
				logger.DebugFn("Error creating overlay directory: %v", err)
				continue
			}

			destPath := filepath.Join(destDir, file.Name())

			// Copy the overlay file
			if err := CopyFile(srcPath, destPath); err != nil {
				logger.DebugFn("Warning: Could not copy overlay %s: %v", file.Name(), err)
			} else {
				themePath := filepath.Join("Overlays", systemTag, file.Name())

				// Add to manifest
				manifest.PathMappings.Overlays = append(
					manifest.PathMappings.Overlays,
					PathMapping{
						ThemePath:  themePath,
						SystemPath: srcPath,
						Metadata: map[string]string{
							"SystemTag":   systemTag,
							"OverlayName": file.Name(),
						},
					},
				)

				hasOverlays = true
				logger.DebugFn("Exported overlay %s for system %s", file.Name(), systemTag)
			}
		}

		// If this system had overlays, add it to the systems list
		if hasOverlays {
			manifest.Content.Overlays.Present = true

			// Check if system is already in the list
			var systemExists bool
			for _, sys := range manifest.Content.Overlays.Systems {
				if sys == systemTag {
					systemExists = true
					break
				}
			}

			if !systemExists {
				manifest.Content.Overlays.Systems = append(manifest.Content.Overlays.Systems, systemTag)
			}
		}
	}
}

// exportFonts scans for and exports system fonts
// exportFonts scans for and exports system fonts
func exportFonts(themePath string, manifest *ThemeManifest, logger *Logger) error {
	logger.DebugFn("Exporting fonts")

	// Create Fonts directory if it doesn't exist
	fontsDir := filepath.Join(themePath, "Fonts")
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		logger.DebugFn("Error creating Fonts directory: %v", err)
		return err
	}

	// Initialize fonts section in manifest
	manifest.Content.Fonts.Present = false
	manifest.Content.Fonts.OGReplaced = false
	manifest.Content.Fonts.NextReplaced = false
	manifest.PathMappings.Fonts = make(map[string]PathMapping)

	// Define font paths to check - CORRECTED PATHS
	fontPaths := map[string]string{
		"OG":          "/mnt/SDCARD/.system/res/font2.ttf",
		"OG.backup":   "/mnt/SDCARD/.system/res/font2.backup.ttf", // Corrected extension
		"Next":        "/mnt/SDCARD/.system/res/font1.ttf",
		"Next.backup": "/mnt/SDCARD/.system/res/font1.backup.ttf", // Corrected extension
	}

	// Check and export each font
	for fontName, sourcePath := range fontPaths {
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			logger.DebugFn("Font file not found: %s", sourcePath)
			continue
		}

		dstPath := filepath.Join(fontsDir, fontName+".ttf")

		if err := CopyFile(sourcePath, dstPath); err != nil {
			logger.DebugFn("Warning: Could not copy font %s: %v", fontName, err)
			continue
		}

		// Add to manifest
		manifest.PathMappings.Fonts[fontName] = PathMapping{
			ThemePath:  "Fonts/" + fontName + ".ttf",
			SystemPath: sourcePath,
		}

		// Update content flags
		manifest.Content.Fonts.Present = true
		if fontName == "OG" {
			manifest.Content.Fonts.OGReplaced = true
		} else if fontName == "Next" {
			manifest.Content.Fonts.NextReplaced = true
		}

		logger.DebugFn("Exported font: %s to %s", fontName, dstPath)
	}

	return nil
}

// readAccentSettingsFromSystem reads accent settings from the system and updates the manifest
func readAccentSettingsFromSystem(manifest *ThemeManifest, logger *Logger) error {
	// Path to the accent settings file
	settingsPath := "/mnt/SDCARD/.userdata/shared/minuisettings.txt"

	// Check if settings file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		logger.DebugFn("Accent settings file not found: %s", settingsPath)
		return fmt.Errorf("accent settings file not found: %s", settingsPath)
	}

	// Read settings file
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("error reading accent settings file: %w", err)
	}

	// Parse settings
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Store color values directly in manifest
		switch key {
		case "color1":
			manifest.AccentColors.Color1 = value
		case "color2":
			manifest.AccentColors.Color2 = value
		case "color3":
			manifest.AccentColors.Color3 = value
		case "color4":
			manifest.AccentColors.Color4 = value
		case "color5":
			manifest.AccentColors.Color5 = value
		case "color6":
			manifest.AccentColors.Color6 = value
		}
	}

	// Mark accent colors as included
	manifest.Content.Settings.AccentsIncluded = true
	logger.DebugFn("Read accent settings from system and updated manifest")

	return nil
}

// readLEDSettingsFromSystem reads LED settings from the system and updates the manifest
func readLEDSettingsFromSystem(manifest *ThemeManifest, logger *Logger) error {
	// Path to the LED settings file
	settingsPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"

	// Check if settings file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		logger.DebugFn("LED settings file not found: %s", settingsPath)
		return fmt.Errorf("LED settings file not found: %s", settingsPath)
	}

	// Read settings file
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("error reading LED settings file: %w", err)
	}

	// Parse settings
	lines := strings.Split(string(content), "\n")
	var currentLED *LEDSetting
	var currentSection string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for section header [X]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]

			// Determine which LED setting to update
			switch currentSection {
			case "F1 key":
				currentLED = &manifest.LEDSettings.F1Key
			case "F2 key":
				currentLED = &manifest.LEDSettings.F2Key
			case "Top bar":
				currentLED = &manifest.LEDSettings.TopBar
			case "L&R triggers":
				currentLED = &manifest.LEDSettings.LRTriggers
			default:
				currentLED = nil
			}
			continue
		}

		// If not in a valid section, skip
		if currentLED == nil {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Update LED setting based on key
		switch key {
		case "effect":
			currentLED.Effect, _ = strconv.Atoi(value)
		case "color1":
			currentLED.Color1 = value
		case "color2":
			currentLED.Color2 = value
		case "speed":
			currentLED.Speed, _ = strconv.Atoi(value)
		case "brightness":
			currentLED.Brightness, _ = strconv.Atoi(value)
		case "trigger":
			currentLED.Trigger, _ = strconv.Atoi(value)
		case "inbrightness":
			currentLED.InBrightness, _ = strconv.Atoi(value)
		}
	}

	// Mark LED settings as included
	manifest.Content.Settings.LEDsIncluded = true
	logger.DebugFn("Read LED settings from system and updated manifest")

	return nil
}
