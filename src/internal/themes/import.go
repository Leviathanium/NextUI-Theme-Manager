// src/internal/themes/import.go
// Simplified implementation of theme import functionality

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

// ImportTheme imports a theme package
// Modify the ImportTheme function in src/internal/themes/import.go
// This updates the ImportTheme function to always clean up existing components
// before applying new ones from the theme pack, matching the behavior of individual component packs.

func ImportTheme(themeName string) error {
	// Create logger
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting theme import for: %s", themeName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Full path to theme - look in Themes directory directly instead of Themes/Imports
	themePath := filepath.Join(cwd, "Themes", themeName)

	// Validate theme
	manifest, err := ValidateTheme(themePath, logger)
	if err != nil {
		logger.DebugFn("Theme validation failed: %v", err)
		return fmt.Errorf("theme validation failed: %w", err)
	}

	// Get system paths BEFORE updating manifest
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logger.DebugFn("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Update manifest based on theme content - now passing systemPaths
	// This is critical for minimal manifests to work properly
	if err := UpdateManifestFromThemeContent(themePath, manifest, systemPaths, logger); err != nil {
		logger.DebugFn("Warning: Error updating manifest from theme content: %v", err)
		// Continue anyway with the original manifest
	}

	// IMPORTANT CHANGE: Always clean up existing components before applying new ones
	// This ensures consistency with how individual component packs work

	// Clean up existing wallpapers (regardless of whether the theme includes them)
	logger.DebugFn("Cleaning up existing wallpapers before theme import")
	if err := cleanupExistingWallpapers(systemPaths, logger); err != nil {
		logger.DebugFn("Warning: Error cleaning up existing wallpapers: %v", err)
		// Continue with import anyway
	}

	// Clean up existing icons (regardless of whether the theme includes them)
	logger.DebugFn("Cleaning up existing icons before theme import")
	if err := cleanupExistingIcons(systemPaths, logger); err != nil {
		logger.DebugFn("Warning: Error cleaning up existing icons: %v", err)
		// Continue with import anyway
	}

	// Clean up existing overlays (regardless of whether the theme includes them)
	logger.DebugFn("Cleaning up existing overlays before theme import")
	if err := cleanupExistingOverlays(systemPaths, logger); err != nil {
		logger.DebugFn("Warning: Error cleaning up existing overlays: %v", err)
		// Continue with import anyway
	}

	// We've removed the conditional cleanup in favor of always cleaning up
	// The old code cleaned up only if the theme didn't include these components:
	// if err := cleanupExistingComponents(manifest, systemPaths, logger); err != nil {
	//     logger.DebugFn("Warning: Error cleaning up existing components: %v", err)
	//     // Continue with import anyway
	// }

	// Apply theme components based on the (now updated) manifest
	if err := importThemeFiles(themePath, manifest, systemPaths, logger); err != nil {
		logger.DebugFn("Error importing theme files: %v", err)
		return fmt.Errorf("error importing theme files: %w", err)
	}

	// Apply accent colors directly from manifest
	if manifest.Content.Settings.AccentsIncluded {
		if err := applyAccentSettings(manifest, logger); err != nil {
			logger.DebugFn("Warning: Error applying accent settings: %v", err)
		}
	}

	// Apply LED settings directly from manifest
	if manifest.Content.Settings.LEDsIncluded {
		if err := applyLEDSettings(manifest, logger); err != nil {
			logger.DebugFn("Warning: Error applying LED settings: %v", err)
		}
	}

	logger.DebugFn("Theme import completed successfully: %s", themeName)

	// Show success message to user
	ui.ShowMessage(fmt.Sprintf("Theme '%s' by %s imported successfully!",
		manifest.ThemeInfo.Name, manifest.ThemeInfo.Author), "3")

	return nil
}

// Note: The cleanupExistingComponents function can be kept for backward compatibility,
// but it's no longer called from ImportTheme since we now always clean up components
// regardless of whether the theme includes them or not.

// importThemeFiles copies all files from the theme to the system based on path mappings
func importThemeFiles(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) error {
	// Ensure media directories exist
	if systemPaths != nil {
		if err := system.EnsureMediaDirectories(systemPaths); err != nil {
			logger.DebugFn("Warning: Failed to ensure media directories: %v", err)
		}
	}

	// Process wallpaper mappings
	for _, mapping := range manifest.PathMappings.Wallpapers {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy wallpaper: %v", err)
			// Continue with other files
		}
	}

	// Process icon mappings with special handling for system icons
	for _, mapping := range manifest.PathMappings.Icons {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Get the icon filename
		iconName := filepath.Base(srcPath)

		// Check if this is a system icon that needs special handling
		if mapping.Metadata != nil && mapping.Metadata["IconType"] == "System" {
			// Use our helper function to get proper destination for system icons
			newDstPath, err := GetSystemIconDestination(srcPath, iconName, dstPath, systemPaths, logger)
			if err != nil {
				logger.DebugFn("Warning: Error determining system icon destination: %v", err)
			} else if newDstPath != dstPath {
				// Update the destination path if it changed
				dstPath = newDstPath
				logger.DebugFn("Updated system icon destination: %s", dstPath)
			}
		}

		// Copy the file to the (possibly renamed) destination
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy icon: %v", err)
			// Continue with other files
		}
	}

	// Process overlay mappings
	for _, mapping := range manifest.PathMappings.Overlays {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy overlay: %v", err)
			// Continue with other files
		}
	}

	// Process font mappings
	for fontType, mapping := range manifest.PathMappings.Fonts {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy font %s: %v", fontType, err)
			// Continue with other files
		}
	}

	// Process settings mappings
	for settingType, mapping := range manifest.PathMappings.Settings {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy setting %s: %v", settingType, err)
			// Continue with other files
		}
	}

	return nil
}

// UpdateManifestFromThemeContent scans a theme directory and updates the manifest
func UpdateManifestFromThemeContent(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) error {
	// Update wallpapers
	if err := updateWallpaperMappings(themePath, manifest, systemPaths, logger); err != nil {
		logger.DebugFn("Warning: Error updating wallpaper mappings: %v", err)
	}

	// Update icons
	if err := updateIconMappings(themePath, manifest, systemPaths, logger); err != nil {
		logger.DebugFn("Warning: Error updating icon mappings: %v", err)
	}

	// Update overlays
	if err := updateOverlayMappings(themePath, manifest, systemPaths, logger); err != nil {
		logger.DebugFn("Warning: Error updating overlay mappings: %v", err)
	}

    // Update fonts
    if err := updateFontMappings(themePath, manifest, systemPaths, logger); err != nil {
        logger.DebugFn("Warning: Error updating font mappings: %v", err)
    }

	// Write updated manifest back to file
	return WriteManifest(themePath, manifest, logger)
}

// updateIconMappings scans icons in the theme and updates manifest mappings
func updateIconMappings(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) error {
	// Create a map of existing mappings for quick lookup
	existingMappings := make(map[string]bool)
	for _, mapping := range manifest.PathMappings.Icons {
		existingMappings[mapping.ThemePath] = true
	}

	// Regular expression to extract system tag from filenames
	tagRegex := regexp.MustCompile(`\((.*?)\)`)

	// Process system icons
	systemIconsDir := filepath.Join(themePath, "Icons", "SystemIcons")
	if _, err := os.Stat(systemIconsDir); err == nil {
		entries, err := os.ReadDir(systemIconsDir)
		if err != nil {
			logger.DebugFn("Warning: Error reading system icons directory: %v", err)
		} else {
			// Map to track which system tags we've already processed
			// This helps us handle duplicate system tags
			processedSystemTags := make(map[string]bool)

			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				themePath := filepath.Join("Icons/SystemIcons", entry.Name())

				// Skip if this file is already in mappings
				if existingMappings[themePath] {
					continue
				}

				// Determine where this file should go based on naming
				var systemPath string
				var metadata map[string]string

				// Special case handling for predefined names
				switch strings.TrimSuffix(entry.Name(), ".png") {
				case "Recently Played":
					systemPath = filepath.Join(systemPaths.Root, ".media", "Recently Played.png")
					metadata = map[string]string{
						"SystemName": "Recently Played",
						"IconType":   "System",
					}

				case "Tools":
					// Get parent directory of Tools path since it includes tg5040
					toolsParentDir := filepath.Dir(systemPaths.Tools)
					systemPath = filepath.Join(toolsParentDir, ".media", "tg5040.png")
					metadata = map[string]string{
						"SystemName": "Tools",
						"IconType":   "System",
					}

				case "Collections":
					systemPath = filepath.Join(systemPaths.Root, ".media", "Collections.png")
					metadata = map[string]string{
						"SystemName": "Collections",
						"IconType":   "System",
					}

				default:
					// Check for system tag in filename
					matches := tagRegex.FindStringSubmatch(entry.Name())
					if len(matches) >= 2 {
						systemTag := matches[1]

						// Skip if we've already processed an icon with this tag and priority handling is enabled
						if processedSystemTags[systemTag] {
							logger.DebugFn("Skipping duplicate system tag '%s' for icon: %s (already processed)",
								systemTag, entry.Name())
							continue
						}
						processedSystemTags[systemTag] = true

						// Full system icon file name
						iconName := entry.Name()

						// Look for existing ROM directory with this tag
						var exactSystemName string
						var matchFound bool

						for _, system := range systemPaths.Systems {
							if system.Tag == systemTag {
								exactSystemName = system.Name
								matchFound = true
								break
							}
						}

						// If no exact match found, use original name
						if !matchFound {
							exactSystemName = strings.TrimSuffix(iconName, ".png")
						}

						// System path based on the icon name (with tag)
						systemPath = filepath.Join(systemPaths.Roms, ".media", iconName)
						metadata = map[string]string{
							"SystemName": strings.TrimSuffix(iconName, ".png"),
							"SystemTag":  systemTag,
							"IconType":   "System",
							// Add match info to metadata
							"MatchFound":      fmt.Sprintf("%v", matchFound),
							"ExactSystemName": exactSystemName,
						}
					}
				}

				// If we determined a system path, add to mappings
				if systemPath != "" {
					manifest.PathMappings.Icons = append(
						manifest.PathMappings.Icons,
						PathMapping{
							ThemePath:  themePath,
							SystemPath: systemPath,
							Metadata:   metadata,
						},
					)
					manifest.Content.Icons.SystemCount++
					logger.DebugFn("Added mapping for system icon: %s -> %s", themePath, systemPath)
				} else {
					logger.DebugFn("Could not determine system path for icon: %s", entry.Name())
				}
			}
		}
	}

	// [Rest of function remains unchanged]
	// Process tool icons
	toolIconsDir := filepath.Join(themePath, "Icons", "ToolIcons")
	if _, err := os.Stat(toolIconsDir); err == nil {
		entries, err := os.ReadDir(toolIconsDir)
		if err != nil {
			logger.DebugFn("Warning: Error reading tool icons directory: %v", err)
		} else {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				themePath := filepath.Join("Icons/ToolIcons", entry.Name())

				// Skip if this file is already in mappings
				if existingMappings[themePath] {
					continue
				}

				// Extract tool name
				toolName := strings.TrimSuffix(entry.Name(), ".png")
				systemPath := filepath.Join(systemPaths.Tools, toolName, ".media", toolName+".png")
				metadata := map[string]string{
					"ToolName": toolName,
					"IconType": "Tool",
				}

				manifest.PathMappings.Icons = append(
					manifest.PathMappings.Icons,
					PathMapping{
						ThemePath:  themePath,
						SystemPath: systemPath,
						Metadata:   metadata,
					},
				)
				manifest.Content.Icons.ToolCount++
				logger.DebugFn("Added mapping for tool icon: %s -> %s", themePath, systemPath)
			}
		}
	}

	// Process collection icons
	collectionIconsDir := filepath.Join(themePath, "Icons", "CollectionIcons")
	if _, err := os.Stat(collectionIconsDir); err == nil {
		entries, err := os.ReadDir(collectionIconsDir)
		if err != nil {
			logger.DebugFn("Warning: Error reading collection icons directory: %v", err)
		} else {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				themePath := filepath.Join("Icons/CollectionIcons", entry.Name())

				// Skip if this file is already in mappings
				if existingMappings[themePath] {
					continue
				}

				// Extract collection name
				collectionName := strings.TrimSuffix(entry.Name(), ".png")
				systemPath := filepath.Join(systemPaths.Root, "Collections", collectionName, ".media", collectionName+".png")
				metadata := map[string]string{
					"CollectionName": collectionName,
					"IconType":       "Collection",
				}

				manifest.PathMappings.Icons = append(
					manifest.PathMappings.Icons,
					PathMapping{
						ThemePath:  themePath,
						SystemPath: systemPath,
						Metadata:   metadata,
					},
				)
				manifest.Content.Icons.CollectionCount++
				logger.DebugFn("Added mapping for collection icon: %s -> %s", themePath, systemPath)
			}
		}
	}

	// If we found any icons, mark icons as present
	if manifest.Content.Icons.SystemCount > 0 ||
		manifest.Content.Icons.ToolCount > 0 ||
		manifest.Content.Icons.CollectionCount > 0 {
		manifest.Content.Icons.Present = true
	}

	return nil
}

// updateOverlayMappings scans overlays in the theme and updates manifest mappings
func updateOverlayMappings(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) error {
	// Create a map of existing mappings for quick lookup
	existingMappings := make(map[string]bool)
	for _, mapping := range manifest.PathMappings.Overlays {
		existingMappings[mapping.ThemePath] = true
	}

	// Process overlay directories
	overlaysDir := filepath.Join(themePath, "Overlays")
	if _, err := os.Stat(overlaysDir); os.IsNotExist(err) {
		logger.DebugFn("No Overlays directory found in theme")
		return nil
	}

	// List system directories
	entries, err := os.ReadDir(overlaysDir)
	if err != nil {
		logger.DebugFn("Error reading Overlays directory: %v", err)
		return err
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

		// Process each overlay file
		for _, file := range overlayFiles {
			if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
				continue
			}

			// Only process PNG files
			if !strings.HasSuffix(strings.ToLower(file.Name()), ".png") {
				continue
			}

			themePath := filepath.Join("Overlays", systemTag, file.Name())

			// Skip if already in mappings
			if existingMappings[themePath] {
				continue
			}

			// Determine system path
			systemPath := filepath.Join(systemPaths.Root, "Overlays", systemTag, file.Name())

			// Add to manifest
			manifest.PathMappings.Overlays = append(
				manifest.PathMappings.Overlays,
				PathMapping{
					ThemePath:  themePath,
					SystemPath: systemPath,
					Metadata: map[string]string{
						"SystemTag":   systemTag,
						"OverlayName": file.Name(),
					},
				},
			)

			hasOverlays = true
			logger.DebugFn("Added mapping for overlay %s for system %s", file.Name(), systemTag)
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

	return nil
}

// updateFontMappings scans fonts in the theme and updates manifest mappings
func updateFontMappings(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) error {
    logger.DebugFn("Updating font mappings in theme manifest")

    // Define system paths for fonts - CORRECTED PATHS
    fontSystemPaths := map[string]string{
        "OG":          "/mnt/SDCARD/.system/res/font2.ttf",
        "OG.backup":   "/mnt/SDCARD/.system/res/font2.backup.ttf",
        "Next":        "/mnt/SDCARD/.system/res/font1.ttf",
        "Next.backup": "/mnt/SDCARD/.system/res/font1.backup.ttf",
    }

    // Initialize font mappings map if it doesn't exist
    if manifest.PathMappings.Fonts == nil {
        manifest.PathMappings.Fonts = make(map[string]PathMapping)
    }

    // Check for fonts directory
    fontsDir := filepath.Join(themePath, "Fonts")
    if _, err := os.Stat(fontsDir); os.IsNotExist(err) {
        logger.DebugFn("No Fonts directory found in theme")
        return nil
    }

    // Font files to check for
    fontFiles := []string{
        "OG.ttf",
        "Next.ttf",
        "OG.backup.ttf",
        "Next.backup.ttf",
    }

    // Check for each font file
    for _, fontFile := range fontFiles {
        fontPath := filepath.Join(fontsDir, fontFile)
        if _, err := os.Stat(fontPath); err == nil {
            // Font file exists
            fontName := strings.TrimSuffix(fontFile, ".ttf")

            // Add to manifest
            manifest.PathMappings.Fonts[fontName] = PathMapping{
                ThemePath:  filepath.Join("Fonts", fontFile),
                SystemPath: fontSystemPaths[fontName],
            }

            // Update content flags
            manifest.Content.Fonts.Present = true
            if fontName == "OG" {
                manifest.Content.Fonts.OGReplaced = true
            } else if fontName == "Next" {
                manifest.Content.Fonts.NextReplaced = true
            }

            logger.DebugFn("Added font to manifest: %s", fontName)
        }
    }

    return nil
}

func updateWallpaperMappings(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) error {
	// Create a map of existing mappings for quick lookup
	existingMappings := make(map[string]bool)
	for _, mapping := range manifest.PathMappings.Wallpapers {
		existingMappings[mapping.ThemePath] = true
	}

	// Regular expression to extract system tag from filenames
	tagRegex := regexp.MustCompile(`\((.*?)\)`)

	// Process system wallpapers
	systemWallpapersDir := filepath.Join(themePath, "Wallpapers", "SystemWallpapers")
	if _, err := os.Stat(systemWallpapersDir); err == nil {
		entries, err := os.ReadDir(systemWallpapersDir)
		if err != nil {
			logger.DebugFn("Warning: Error reading system wallpapers directory: %v", err)
		} else {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				themePath := filepath.Join("Wallpapers/SystemWallpapers", entry.Name())

				// Skip if this file is already in mappings
				if existingMappings[themePath] {
					continue
				}

				// Determine where this file should go based on naming
				var systemPath string
				var metadata map[string]string

				// Special case handling for predefined names
				switch strings.TrimSuffix(entry.Name(), ".png") {
				case "Root":
					systemPath = filepath.Join(systemPaths.Root, "bg.png")
					metadata = map[string]string{
						"SystemName":    "Root",
						"WallpaperType": "Main",
					}

				case "Root-Media":
					systemPath = filepath.Join(systemPaths.Root, ".media", "bg.png")
					metadata = map[string]string{
						"SystemName":    "Root",
						"WallpaperType": "Media",
					}

				case "Recently Played":
					systemPath = filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
					metadata = map[string]string{
						"SystemName":    "Recently Played",
						"WallpaperType": "Media",
					}

				case "Tools":
					systemPath = filepath.Join(systemPaths.Tools, ".media", "bg.png")
					metadata = map[string]string{
						"SystemName":    "Tools",
						"WallpaperType": "Media",
					}

				case "Collections":
					systemPath = filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")
					metadata = map[string]string{
						"SystemName":    "Collections",
						"WallpaperType": "Media",
					}

				default:
					// Check for system tag in filename
					matches := tagRegex.FindStringSubmatch(entry.Name())
					if len(matches) >= 2 {
						systemTag := matches[1]

						// Extract the system name without the tag
						fileName := entry.Name()
						baseName := strings.TrimSuffix(fileName, ".png")
						systemName := strings.TrimSuffix(strings.Split(baseName, "(")[0], " ")

						// Find matching system by tag
						var systemFound bool
						for _, system := range systemPaths.Systems {
							if system.Tag == systemTag {
								systemPath = filepath.Join(system.MediaPath, "bg.png")
								metadata = map[string]string{
									"SystemName":    systemName,
									"SystemTag":     systemTag,
									"WallpaperType": "System",
								}
								systemFound = true
								break
							}
						}

						// If system not found in paths, create a default path
						if !systemFound && systemTag != "" {
							systemPath = filepath.Join(systemPaths.Roms, fmt.Sprintf("%s (%s)", systemName, systemTag), ".media", "bg.png")
							metadata = map[string]string{
								"SystemName":    systemName,
								"SystemTag":     systemTag,
								"WallpaperType": "System",
							}
						}
					}
				}

				// If we determined a system path, add to mappings
				if systemPath != "" {
					manifest.PathMappings.Wallpapers = append(
						manifest.PathMappings.Wallpapers,
						PathMapping{
							ThemePath:  themePath,
							SystemPath: systemPath,
							Metadata:   metadata,
						},
					)
					manifest.Content.Wallpapers.Count++
					logger.DebugFn("Added mapping for system wallpaper: %s -> %s", themePath, systemPath)
				} else {
					logger.DebugFn("Could not determine system path for wallpaper: %s", entry.Name())
				}
			}
		}
	}

	// NEW: Process list wallpapers
	listWallpapersDir := filepath.Join(themePath, "Wallpapers", "ListWallpapers")
	if _, err := os.Stat(listWallpapersDir); err == nil {
		entries, err := os.ReadDir(listWallpapersDir)
		if err != nil {
			logger.DebugFn("Warning: Error reading list wallpapers directory: %v", err)
		} else {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				themePath := filepath.Join("Wallpapers/ListWallpapers", entry.Name())

				// Skip if this file is already in mappings
				if existingMappings[themePath] {
					continue
				}

				// Check for system tag in filename and -list suffix
				fileName := entry.Name()
				baseName := strings.TrimSuffix(fileName, ".png")

				// See if this is a list wallpaper (filename ends with -list)
				if !strings.HasSuffix(baseName, "-list") {
					logger.DebugFn("List wallpaper doesn't have -list suffix: %s", fileName)
					continue
				}

				// Remove the -list suffix to get the system name and tag
				baseNameWithoutSuffix := strings.TrimSuffix(baseName, "-list")

				// Extract system tag
				matches := tagRegex.FindStringSubmatch(baseNameWithoutSuffix)
				if len(matches) >= 2 {
					systemTag := matches[1]
					systemName := strings.TrimSuffix(strings.Split(baseNameWithoutSuffix, "(")[0], " ")

					// Find matching system by tag
					var systemFound bool
					for _, system := range systemPaths.Systems {
						if system.Tag == systemTag {
							systemPath := filepath.Join(system.MediaPath, "bglist.png")
							metadata := map[string]string{
								"SystemName":    systemName,
								"SystemTag":     systemTag,
								"WallpaperType": "List",
							}

							manifest.PathMappings.Wallpapers = append(
								manifest.PathMappings.Wallpapers,
								PathMapping{
									ThemePath:  themePath,
									SystemPath: systemPath,
									Metadata:   metadata,
								},
							)
							manifest.Content.Wallpapers.Count++
							manifest.Content.Wallpapers.Present = true
							systemFound = true
							logger.DebugFn("Added mapping for list wallpaper: %s -> %s", themePath, systemPath)
							break
						}
					}

					// If system not found in paths, create a default path
					if !systemFound && systemTag != "" {
						systemPath := filepath.Join(systemPaths.Roms, fmt.Sprintf("%s (%s)", systemName, systemTag), ".media", "bglist.png")
						metadata := map[string]string{
							"SystemName":    systemName,
							"SystemTag":     systemTag,
							"WallpaperType": "List",
						}

						manifest.PathMappings.Wallpapers = append(
							manifest.PathMappings.Wallpapers,
							PathMapping{
								ThemePath:  themePath,
								SystemPath: systemPath,
								Metadata:   metadata,
							},
						)
						manifest.Content.Wallpapers.Count++
						manifest.Content.Wallpapers.Present = true
						logger.DebugFn("Added default mapping for list wallpaper: %s -> %s", themePath, systemPath)
					}
				} else {
					logger.DebugFn("Could not extract system tag from list wallpaper: %s", fileName)
				}
			}
		}
	}

	// Process collection wallpapers
	collectionWallpapersDir := filepath.Join(themePath, "Wallpapers", "CollectionWallpapers")
	if _, err := os.Stat(collectionWallpapersDir); err == nil {
		entries, err := os.ReadDir(collectionWallpapersDir)
		if err != nil {
			logger.DebugFn("Warning: Error reading collection wallpapers directory: %v", err)
		} else {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				themePath := filepath.Join("Wallpapers/CollectionWallpapers", entry.Name())

				// Skip if this file is already in mappings
				if existingMappings[themePath] {
					continue
				}

				// Determine collection name and system path
				collectionName := strings.TrimSuffix(entry.Name(), ".png")
				systemPath := filepath.Join(systemPaths.Root, "Collections", collectionName, ".media", "bg.png")

				metadata := map[string]string{
					"CollectionName": collectionName,
					"WallpaperType":  "Collection",
				}

				// Add to manifest
				manifest.PathMappings.Wallpapers = append(
					manifest.PathMappings.Wallpapers,
					PathMapping{
						ThemePath:  themePath,
						SystemPath: systemPath,
						Metadata:   metadata,
					},
				)
				manifest.Content.Wallpapers.Count++
				logger.DebugFn("Added mapping for collection wallpaper: %s -> %s", themePath, systemPath)
			}
		}
	}

	// If we found any wallpapers, mark wallpapers as present
	if manifest.Content.Wallpapers.Count > 0 {
		manifest.Content.Wallpapers.Present = true
	}

	return nil
}

// updateAccentSettings reads accent settings from file and updates manifest
func updateAccentSettings(themePath string, manifest *ThemeManifest, logger *Logger) error {
	settingsPath := filepath.Join(themePath, "Settings", "minuisettings.txt")

	// Check if settings file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		logger.DebugFn("Accent settings file not found: %s", settingsPath)
		return nil
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

	// We've parsed the settings file directly into the manifest,
	// so we don't need to keep the file path mapping
	delete(manifest.PathMappings.Settings, "accents")

	return nil
}

// updateLEDSettings reads LED settings from file and updates manifest
func updateLEDSettings(themePath string, manifest *ThemeManifest, logger *Logger) error {
	settingsPath := filepath.Join(themePath, "Settings", "ledsettings_brick.txt")

	// Check if settings file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		logger.DebugFn("LED settings file not found: %s", settingsPath)
		return nil
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

	// We've parsed the settings file directly into the manifest,
	// so we don't need to keep the file path mapping
	delete(manifest.PathMappings.Settings, "leds")

	return nil
}

// Update these two functions in src/internal/themes/import.go

// This function should replace the existing applyAccentSettings in import.go
func applyAccentSettings(manifest *ThemeManifest, logger *Logger) error {
    // Get path to settings file
    settingsPath := "/mnt/SDCARD/.userdata/shared/minuisettings.txt"

    // Map of color keys to their values from the manifest
    colorValues := map[string]string{
        "color1": manifest.AccentColors.Color1,
        "color2": manifest.AccentColors.Color2,
        "color3": manifest.AccentColors.Color3,
        "color4": manifest.AccentColors.Color4,
        "color5": manifest.AccentColors.Color5,
        "color6": manifest.AccentColors.Color6,
    }

    // Track which color keys we've seen
    colorKeySeen := map[string]bool{
        "color1": false,
        "color2": false,
        "color3": false,
        "color4": false,
        "color5": false,
        "color6": false,
    }

    // Final settings map and order preservation
    settings := make(map[string]string)
    var keyOrder []string

    // Check if file exists
    fileInfo, err := os.Stat(settingsPath)
    if err == nil && fileInfo.Size() > 0 {
        // File exists and has content
        existingContent, err := os.ReadFile(settingsPath)
        if err == nil {
            // Parse existing settings
            lines := strings.Split(string(existingContent), "\n")
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

                // Check if this is a color key
                if newValue, isColorKey := colorValues[key]; isColorKey {
                    // Use the new color value and mark as seen
                    settings[key] = newValue
                    colorKeySeen[key] = true
                } else {
                    // Keep existing value for non-color keys
                    settings[key] = value
                }

                // Add to ordered keys
                keyOrder = append(keyOrder, key)
            }

            // Add any unseen color keys at the end
            for key, seen := range colorKeySeen {
                if !seen {
                    settings[key] = colorValues[key]
                    keyOrder = append(keyOrder, key)
                }
            }
        } else {
            logger.DebugFn("Warning: Could not read existing settings file: %v", err)
            // Fall through to create new file
        }
    }

    // If we don't have any settings yet (no file or couldn't read it),
    // initialize with just the color keys
    if len(settings) == 0 {
        colorKeys := []string{"color1", "color2", "color3", "color4", "color5", "color6"}
        for _, key := range colorKeys {
            settings[key] = colorValues[key]
            keyOrder = append(keyOrder, key)
        }
    }

    // Build new content with all settings in the preserved order
    var content strings.Builder
    for _, key := range keyOrder {
        if value, exists := settings[key]; exists {
            content.WriteString(fmt.Sprintf("%s=%s\n", key, value))
        }
    }

    // Write updated settings to file
    if err := os.WriteFile(settingsPath, []byte(content.String()), 0644); err != nil {
        return fmt.Errorf("error writing accent settings: %w", err)
    }

    logger.DebugFn("Applied accent settings to %s while preserving other settings and their order", settingsPath)
    return nil
}


// applyLEDSettings applies LED settings from manifest
func applyLEDSettings(manifest *ThemeManifest, logger *Logger) error {
	// Create content for ledsettings_brick.txt
	var content strings.Builder

	// F1 Key
	content.WriteString("[F1 key]\n")
	content.WriteString(fmt.Sprintf("effect=%d\n", manifest.LEDSettings.F1Key.Effect))
	content.WriteString(fmt.Sprintf("color1=%s\n", manifest.LEDSettings.F1Key.Color1))
	content.WriteString(fmt.Sprintf("color2=%s\n", manifest.LEDSettings.F1Key.Color2))
	content.WriteString(fmt.Sprintf("speed=%d\n", manifest.LEDSettings.F1Key.Speed))
	content.WriteString(fmt.Sprintf("brightness=%d\n", manifest.LEDSettings.F1Key.Brightness))
	content.WriteString(fmt.Sprintf("trigger=%d\n", manifest.LEDSettings.F1Key.Trigger))
	content.WriteString(fmt.Sprintf("inbrightness=%d\n", manifest.LEDSettings.F1Key.InBrightness))
	content.WriteString("\n")

	// F2 Key
	content.WriteString("[F2 key]\n")
	content.WriteString(fmt.Sprintf("effect=%d\n", manifest.LEDSettings.F2Key.Effect))
	content.WriteString(fmt.Sprintf("color1=%s\n", manifest.LEDSettings.F2Key.Color1))
	content.WriteString(fmt.Sprintf("color2=%s\n", manifest.LEDSettings.F2Key.Color2))
	content.WriteString(fmt.Sprintf("speed=%d\n", manifest.LEDSettings.F2Key.Speed))
	content.WriteString(fmt.Sprintf("brightness=%d\n", manifest.LEDSettings.F2Key.Brightness))
	content.WriteString(fmt.Sprintf("trigger=%d\n", manifest.LEDSettings.F2Key.Trigger))
	content.WriteString(fmt.Sprintf("inbrightness=%d\n", manifest.LEDSettings.F2Key.InBrightness))
	content.WriteString("\n")

	// Top bar
	content.WriteString("[Top bar]\n")
	content.WriteString(fmt.Sprintf("effect=%d\n", manifest.LEDSettings.TopBar.Effect))
	content.WriteString(fmt.Sprintf("color1=%s\n", manifest.LEDSettings.TopBar.Color1))
	content.WriteString(fmt.Sprintf("color2=%s\n", manifest.LEDSettings.TopBar.Color2))
	content.WriteString(fmt.Sprintf("speed=%d\n", manifest.LEDSettings.TopBar.Speed))
	content.WriteString(fmt.Sprintf("brightness=%d\n", manifest.LEDSettings.TopBar.Brightness))
	content.WriteString(fmt.Sprintf("trigger=%d\n", manifest.LEDSettings.TopBar.Trigger))
	content.WriteString(fmt.Sprintf("inbrightness=%d\n", manifest.LEDSettings.TopBar.InBrightness))
	content.WriteString("\n")

	// L&R triggers
	content.WriteString("[L&R triggers]\n")
	content.WriteString(fmt.Sprintf("effect=%d\n", manifest.LEDSettings.LRTriggers.Effect))
	content.WriteString(fmt.Sprintf("color1=%s\n", manifest.LEDSettings.LRTriggers.Color1))
	content.WriteString(fmt.Sprintf("color2=%s\n", manifest.LEDSettings.LRTriggers.Color2))
	content.WriteString(fmt.Sprintf("speed=%d\n", manifest.LEDSettings.LRTriggers.Speed))
	content.WriteString(fmt.Sprintf("brightness=%d\n", manifest.LEDSettings.LRTriggers.Brightness))
	content.WriteString(fmt.Sprintf("trigger=%d\n", manifest.LEDSettings.LRTriggers.Trigger))
	content.WriteString(fmt.Sprintf("inbrightness=%d\n", manifest.LEDSettings.LRTriggers.InBrightness))
	content.WriteString("\n")

	// Get path to settings file - FIXED PATH
	settingsPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"

	// Write settings to file
	if err := os.WriteFile(settingsPath, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("error writing LED settings: %w", err)
	}

	logger.DebugFn("Applied LED settings to %s", settingsPath)
	return nil
}

func cleanupExistingComponents(manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) error {
	logger.DebugFn("Starting cleanup of existing components")

	// If theme doesn't include wallpapers, clean up existing wallpapers
	if !manifest.Content.Wallpapers.Present {
		logger.DebugFn("Theme doesn't include wallpapers - cleaning up existing wallpapers")

		// Call the updated cleanupExistingWallpapers function which now handles bglist.png
		if err := cleanupExistingWallpapers(systemPaths, logger); err != nil {
			logger.DebugFn("Warning: Error cleaning up existing wallpapers: %v", err)
		}
	} else {
		logger.DebugFn("Theme includes wallpapers - keeping existing wallpapers until they're replaced")
	}

	// If theme doesn't include icons, clean up existing icons
	if !manifest.Content.Icons.Present {
		logger.DebugFn("Theme doesn't include icons - cleaning up existing icons")
		if err := cleanupExistingIcons(systemPaths, logger); err != nil {
			logger.DebugFn("Warning: Error cleaning up existing icons: %v", err)
		}
	} else {
		logger.DebugFn("Theme includes icons - keeping existing icons until they're replaced")
	}

	// If theme doesn't include overlays, clean up existing overlays
	if !manifest.Content.Overlays.Present {
		logger.DebugFn("Theme doesn't include overlays - cleaning up existing overlays")
		if err := cleanupExistingOverlays(systemPaths, logger); err != nil {
			logger.DebugFn("Warning: Error cleaning up existing overlays: %v", err)
		}
	} else {
		logger.DebugFn("Theme includes overlays - keeping existing overlays until they're replaced")
	}

	logger.DebugFn("Completed cleanup of existing components")
	return nil
}

// copyMappedFile copies a file from source to destination with appropriate checks
func copyMappedFile(srcPath, dstPath string, logger *Logger) error {
	// Check if source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		logger.DebugFn("Source file does not exist: %s", srcPath)
		return fmt.Errorf("source file does not exist: %s", srcPath)
	}

	// Create destination directory
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		logger.DebugFn("Failed to create destination directory: %v", err)
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy the file
	if err := CopyFile(srcPath, dstPath); err != nil {
		logger.DebugFn("Failed to copy file: %v", err)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	logger.DebugFn("Copied file: %s -> %s", srcPath, dstPath)
	return nil
}
