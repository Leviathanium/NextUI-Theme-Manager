// src/internal/themes/component_import.go
// Implements import functionality for individual theme components

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
    "regexp"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"nextui-themes/internal/ui"
)

// ImportComponent dispatches to the appropriate import function based on component type
func ImportComponent(componentPath string) error {
	// First, determine the component type from the extension
	ext := filepath.Ext(componentPath)

	var componentType string
	for cType, cExt := range ComponentExtension {
		if cExt == ext {
			componentType = cType
			break
		}
	}

	if componentType == "" {
		return fmt.Errorf("unknown component type for extension: %s", ext)
	}

	// Update the component's manifest based on its actual content
	if err := UpdateComponentManifest(componentPath); err != nil {
		logging.LogDebug("Warning: Error updating component manifest: %v", err)
		// Continue anyway, as we can still try to import with the existing manifest
	}

	// Dispatch to specific import function
	switch componentType {
	case ComponentWallpaper:
		return ImportWallpapers(componentPath)
	case ComponentIcon:
		return ImportIcons(componentPath)
	case ComponentAccent:
		return ImportAccents(componentPath)
	case ComponentLED:
		return ImportLEDs(componentPath)
	case ComponentFont:
		return ImportFonts(componentPath)
	case ComponentOverlay:
		return ImportOverlays(componentPath)
	default:
		return fmt.Errorf("unhandled component type: %s", componentType)
	}
}

// ImportWallpapers imports a wallpaper component package
func ImportWallpapers(componentPath string) error {
    logger := &Logger{
        DebugFn: logging.LogDebug,
    }

    logger.DebugFn("Starting wallpaper import: %s", componentPath)

    // Load the component manifest
    manifestObj, err := LoadComponentManifest(componentPath)
    if err != nil {
        return fmt.Errorf("error loading wallpaper manifest: %w", err)
    }

    // Ensure it's the right type
    manifest, ok := manifestObj.(*WallpaperManifest)
    if !ok {
        return fmt.Errorf("invalid manifest type for wallpaper component")
    }

    // Get system paths
    systemPaths, err := system.GetSystemPaths()
    if err != nil {
        return fmt.Errorf("error getting system paths: %w", err)
    }

    // Ensure media directories exist
    if err := system.EnsureMediaDirectories(systemPaths); err != nil {
        logger.DebugFn("Warning: Error ensuring media directories: %v", err)
    }

    // IMPORTANT: Always clean up existing wallpapers, even if the component has no wallpapers
    // This allows for "default" packages that clear wallpapers
    if err := cleanupExistingWallpapers(systemPaths, logger); err != nil {
        logger.DebugFn("Warning: Error cleaning up existing wallpapers: %v", err)
    }

    // Import wallpapers based on path mappings
    for _, mapping := range manifest.PathMappings {
        srcPath := filepath.Join(componentPath, mapping.ThemePath)
        dstPath := mapping.SystemPath

        // Copy the file
        if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
            logger.DebugFn("Warning: Failed to copy wallpaper: %v", err)
            // Continue with other files
        }
    }

    // Update global manifest to track this component
    componentName := filepath.Base(componentPath)
    if err := UpdateAppliedComponent(ComponentWallpaper, componentName); err != nil {
        logger.DebugFn("Warning: Failed to update global manifest: %v", err)
    }

    logger.DebugFn("Wallpaper import completed: %s", componentPath)

    // Show success message
    ui.ShowMessage(fmt.Sprintf("Wallpapers from '%s' applied successfully!", manifest.ComponentInfo.Name), "3")

    return nil
}

// ImportIcons imports an icon component package
func ImportIcons(componentPath string) error {
    logger := &Logger{
        DebugFn: logging.LogDebug,
    }

    logger.DebugFn("Starting icon import: %s", componentPath)

    // Load the component manifest
    manifestObj, err := LoadComponentManifest(componentPath)
    if err != nil {
        return fmt.Errorf("error loading icon manifest: %w", err)
    }

    // Ensure it's the right type
    manifest, ok := manifestObj.(*IconManifest)
    if !ok {
        return fmt.Errorf("invalid manifest type for icon component")
    }

    // Get system paths
    systemPaths, err := system.GetSystemPaths()
    if err != nil {
        return fmt.Errorf("error getting system paths: %w", err)
    }

    // Ensure media directories exist
    if err := system.EnsureMediaDirectories(systemPaths); err != nil {
        logger.DebugFn("Warning: Error ensuring media directories: %v", err)
    }

    // IMPORTANT: Always clean up existing icons, even if the component has no icons
    // This allows for "default" packages that clear icons
    if err := cleanupExistingIcons(systemPaths, logger); err != nil {
        logger.DebugFn("Warning: Error cleaning up existing icons: %v", err)
    }

    // Import icons based on path mappings
    for _, mapping := range manifest.PathMappings {
        srcPath := filepath.Join(componentPath, mapping.ThemePath)
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

    // Update global manifest to track this component
    componentName := filepath.Base(componentPath)
    if err := UpdateAppliedComponent(ComponentIcon, componentName); err != nil {
        logger.DebugFn("Warning: Failed to update global manifest: %v", err)
    }

    logger.DebugFn("Icon import completed: %s", componentPath)

    // Show success message to user
    ui.ShowMessage(fmt.Sprintf("Icons from '%s' applied successfully!", manifest.ComponentInfo.Name), "3")

    return nil
}

// ImportAccents imports an accent component package
func ImportAccents(componentPath string) error {
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting accent import: %s", componentPath)

	// Load the component manifest
	manifestObj, err := LoadComponentManifest(componentPath)
	if err != nil {
		return fmt.Errorf("error loading accent manifest: %w", err)
	}

	// Ensure it's the right type
	manifest, ok := manifestObj.(*AccentManifest)
	if !ok {
		return fmt.Errorf("invalid manifest type for accent component")
	}

	// Apply accent settings directly from manifest
	settingsPath := "/mnt/SDCARD/.userdata/shared/minuisettings.txt"

	// Create content for minuisettings.txt
	var content strings.Builder
	content.WriteString(fmt.Sprintf("color1=%s\n", manifest.AccentColors.Color1))
	content.WriteString(fmt.Sprintf("color2=%s\n", manifest.AccentColors.Color2))
	content.WriteString(fmt.Sprintf("color3=%s\n", manifest.AccentColors.Color3))
	content.WriteString(fmt.Sprintf("color4=%s\n", manifest.AccentColors.Color4))
	content.WriteString(fmt.Sprintf("color5=%s\n", manifest.AccentColors.Color5))
	content.WriteString(fmt.Sprintf("color6=%s\n", manifest.AccentColors.Color6))

	// Write settings to file
	if err := os.WriteFile(settingsPath, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("error writing accent settings: %w", err)
	}

	// Update global manifest to track this component
	componentName := filepath.Base(componentPath)
	if err := UpdateAppliedComponent(ComponentAccent, componentName); err != nil {
		logger.DebugFn("Warning: Failed to update global manifest: %v", err)
	}

	logger.DebugFn("Accent import completed: %s", componentPath)

	// Show success message
	ui.ShowMessage(fmt.Sprintf("Accent colors from '%s' applied successfully!", manifest.ComponentInfo.Name), "3")

	return nil
}

// ImportLEDs imports a LED component package
func ImportLEDs(componentPath string) error {
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting LED import: %s", componentPath)

	// Load the component manifest
	manifestObj, err := LoadComponentManifest(componentPath)
	if err != nil {
		return fmt.Errorf("error loading LED manifest: %w", err)
	}

	// Ensure it's the right type
	manifest, ok := manifestObj.(*LEDManifest)
	if !ok {
		return fmt.Errorf("invalid manifest type for LED component")
	}

	// Apply LED settings from manifest
	settingsPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"

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

	// Write settings to file
	if err := os.WriteFile(settingsPath, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("error writing LED settings: %w", err)
	}

	// Update global manifest to track this component
	componentName := filepath.Base(componentPath)
	if err := UpdateAppliedComponent(ComponentLED, componentName); err != nil {
		logger.DebugFn("Warning: Failed to update global manifest: %v", err)
	}

	logger.DebugFn("LED import completed: %s", componentPath)

	// Show success message
	ui.ShowMessage(fmt.Sprintf("LED settings from '%s' applied successfully!", manifest.ComponentInfo.Name), "3")

	return nil
}

// ImportFonts imports a font component package
func ImportFonts(componentPath string) error {
    logger := &Logger{
        DebugFn: logging.LogDebug,
    }

    logger.DebugFn("Starting font import: %s", componentPath)

    // Load the component manifest
    manifestObj, err := LoadComponentManifest(componentPath)
    if err != nil {
        return fmt.Errorf("error loading font manifest: %w", err)
    }

    // Ensure it's the right type
    manifest, ok := manifestObj.(*FontManifest)
    if !ok {
        return fmt.Errorf("invalid manifest type for font component")
    }

    // Import fonts based on path mappings
    for fontName, mapping := range manifest.PathMappings {
        srcPath := filepath.Join(componentPath, mapping.ThemePath)
        dstPath := mapping.SystemPath

        // Skip if source file doesn't exist
        if _, err := os.Stat(srcPath); os.IsNotExist(err) {
            logger.DebugFn("Font file doesn't exist: %s", srcPath)
            continue
        }

        // Only create backups for the main font files, not for backup files
        if !strings.Contains(fontName, "backup") && !strings.Contains(dstPath, "backup") {
            // If destination exists and we don't have a backup, create one
            if _, err := os.Stat(dstPath); err == nil {
                // Determine correct backup path format
                var backupPath string
                if strings.HasSuffix(dstPath, "font1.ttf") {
                    backupPath = "/mnt/SDCARD/.system/res/font1.backup.ttf"
                } else if strings.HasSuffix(dstPath, "font2.ttf") {
                    backupPath = "/mnt/SDCARD/.system/res/font2.backup.ttf"
                } else {
                    backupPath = dstPath + ".backup.ttf"  // Fallback
                }

                if _, err := os.Stat(backupPath); os.IsNotExist(err) {
                    if err := CopyFile(dstPath, backupPath); err != nil {
                        logger.DebugFn("Warning: Failed to create font backup for %s: %v", fontName, err)
                    } else {
                        logger.DebugFn("Created backup for font %s: %s", fontName, backupPath)
                    }
                }
            }
        }

        // Copy the font file
        if err := CopyFile(srcPath, dstPath); err != nil {
            logger.DebugFn("Warning: Failed to copy font %s: %v", fontName, err)
        } else {
            logger.DebugFn("Imported font %s to %s", fontName, dstPath)
        }
    }

    // Update global manifest to track this component
    componentName := filepath.Base(componentPath)
    if err := UpdateAppliedComponent(ComponentFont, componentName); err != nil {
        logger.DebugFn("Warning: Failed to update global manifest: %v", err)
    }

    logger.DebugFn("Font import completed: %s", componentPath)

    // Show success message
    ui.ShowMessage(fmt.Sprintf("Fonts from '%s' applied successfully!", manifest.ComponentInfo.Name), "3")

    return nil
}

// ImportOverlays imports an overlay component package
func ImportOverlays(componentPath string) error {
    logger := &Logger{
        DebugFn: logging.LogDebug,
    }

    logger.DebugFn("Starting overlay import: %s", componentPath)

    // Load the component manifest
    manifestObj, err := LoadComponentManifest(componentPath)
    if err != nil {
        return fmt.Errorf("error loading overlay manifest: %w", err)
    }

    // Ensure it's the right type
    manifest, ok := manifestObj.(*OverlayManifest)
    if !ok {
        return fmt.Errorf("invalid manifest type for overlay component")
    }

    // Get system paths
    systemPaths, err := system.GetSystemPaths()
    if err != nil {
        return fmt.Errorf("error getting system paths: %w", err)
    }

    // Create overlays directory if it doesn't exist
    overlaysDir := filepath.Join(systemPaths.Root, "Overlays")
    if err := os.MkdirAll(overlaysDir, 0755); err != nil {
        return fmt.Errorf("error creating overlays directory: %w", err)
    }

    // IMPORTANT: Always clean up existing overlays, even if the component has no overlays
    // This allows for "default" packages that clear overlays
    if err := cleanupExistingOverlays(systemPaths, logger); err != nil {
        logger.DebugFn("Warning: Error cleaning up existing overlays: %v", err)
    }

    // For each system in the package, create the directory
    for _, systemTag := range manifest.Content.Systems {
        systemDir := filepath.Join(overlaysDir, systemTag)
        if err := os.MkdirAll(systemDir, 0755); err != nil {
            logger.DebugFn("Warning: Failed to create system overlay directory %s: %v", systemTag, err)
            continue
        }
    }

    // Import overlays based on path mappings
    for _, mapping := range manifest.PathMappings {
        srcPath := filepath.Join(componentPath, mapping.ThemePath)
        dstPath := mapping.SystemPath

        // Copy the file
        if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
            logger.DebugFn("Warning: Failed to copy overlay: %v", err)
            // Continue with other files
        }
    }

    // Update global manifest to track this component
    componentName := filepath.Base(componentPath)
    if err := UpdateAppliedComponent(ComponentOverlay, componentName); err != nil {
        logger.DebugFn("Warning: Failed to update global manifest: %v", err)
    }

    logger.DebugFn("Overlay import completed: %s", componentPath)

    // Show success message
    ui.ShowMessage(fmt.Sprintf("Overlays from '%s' applied successfully!", manifest.ComponentInfo.Name), "3")

    return nil
}

// Helper functions for cleanup

// cleanupExistingWallpapers removes existing wallpapers before applying new ones
func cleanupExistingWallpapers(systemPaths *system.SystemPaths, logger *Logger) error {
	logger.DebugFn("Cleaning up existing wallpapers")

	// Root wallpaper
	rootBg := filepath.Join(systemPaths.Root, "bg.png")
	if err := os.Remove(rootBg); err != nil && !os.IsNotExist(err) {
		logger.DebugFn("Warning: Could not remove root wallpaper: %v", err)
	} else if err == nil {
		logger.DebugFn("Removed root wallpaper: %s", rootBg)
	}

	// Root media wallpaper
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", "bg.png")
	if err := os.Remove(rootMediaBg); err != nil && !os.IsNotExist(err) {
		logger.DebugFn("Warning: Could not remove root media wallpaper: %v", err)
	} else if err == nil {
		logger.DebugFn("Removed root media wallpaper: %s", rootMediaBg)
	}

	// Recently Played wallpaper
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
	if err := os.Remove(rpBg); err != nil && !os.IsNotExist(err) {
		logger.DebugFn("Warning: Could not remove Recently Played wallpaper: %v", err)
	} else if err == nil {
		logger.DebugFn("Removed Recently Played wallpaper: %s", rpBg)
	}

	// Tools wallpaper
	toolsBg := filepath.Join(systemPaths.Tools, ".media", "bg.png")
	if err := os.Remove(toolsBg); err != nil && !os.IsNotExist(err) {
		logger.DebugFn("Warning: Could not remove Tools wallpaper: %v", err)
	} else if err == nil {
		logger.DebugFn("Removed Tools wallpaper: %s", toolsBg)
	}

	// Collections wallpaper
	collectionsBg := filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")
	if err := os.Remove(collectionsBg); err != nil && !os.IsNotExist(err) {
		logger.DebugFn("Warning: Could not remove Collections wallpaper: %v", err)
	} else if err == nil {
		logger.DebugFn("Removed Collections wallpaper: %s", collectionsBg)
	}

	// System wallpapers
	for _, system := range systemPaths.Systems {
		systemBg := filepath.Join(system.MediaPath, "bg.png")
		if err := os.Remove(systemBg); err != nil && !os.IsNotExist(err) {
			logger.DebugFn("Warning: Could not remove %s wallpaper: %v", system.Name, err)
		} else if err == nil {
			logger.DebugFn("Removed %s wallpaper: %s", system.Name, systemBg)
		}
	}

	// Collection wallpapers
	collectionsDir := filepath.Join(systemPaths.Root, "Collections")
	entries, err := os.ReadDir(collectionsDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			collectionName := entry.Name()
			collectionBg := filepath.Join(collectionsDir, collectionName, ".media", "bg.png")
			if err := os.Remove(collectionBg); err != nil && !os.IsNotExist(err) {
				logger.DebugFn("Warning: Could not remove %s collection wallpaper: %v", collectionName, err)
			} else if err == nil {
				logger.DebugFn("Removed %s collection wallpaper: %s", collectionName, collectionBg)
			}
		}
	}

	return nil
}

// cleanupExistingIcons removes existing icons before applying new ones
func cleanupExistingIcons(systemPaths *system.SystemPaths, logger *Logger) error {
	logger.DebugFn("Cleaning up existing icons")

	// System icons in Roms/.media directory
	romsMediaDir := filepath.Join(systemPaths.Roms, ".media")
	if _, err := os.Stat(romsMediaDir); !os.IsNotExist(err) {
		entries, err := os.ReadDir(romsMediaDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				// Skip non-system icons
				tagRegex := regexp.MustCompile(`\((.*?)\)`)
				if !tagRegex.MatchString(entry.Name()) &&
				   entry.Name() != "Recently Played.png" &&
				   entry.Name() != "Collections.png" &&
				   entry.Name() != "tg5040.png" {
					continue
				}

				systemIcon := filepath.Join(romsMediaDir, entry.Name())
				if err := os.Remove(systemIcon); err != nil && !os.IsNotExist(err) {
					logger.DebugFn("Warning: Could not remove system icon %s: %v", entry.Name(), err)
				} else if err == nil {
					logger.DebugFn("Removed system icon: %s", systemIcon)
				}
			}
		}
	}

	// Root media directory for special icons
	rootMediaDir := filepath.Join(systemPaths.Root, ".media")
	if _, err := os.Stat(rootMediaDir); !os.IsNotExist(err) {
		// Recently Played icon
		rpIcon := filepath.Join(rootMediaDir, "Recently Played.png")
		if err := os.Remove(rpIcon); err != nil && !os.IsNotExist(err) {
			logger.DebugFn("Warning: Could not remove Recently Played icon: %v", err)
		} else if err == nil {
			logger.DebugFn("Removed Recently Played icon: %s", rpIcon)
		}

		// Collections icon
		collectionsIcon := filepath.Join(rootMediaDir, "Collections.png")
		if err := os.Remove(collectionsIcon); err != nil && !os.IsNotExist(err) {
			logger.DebugFn("Warning: Could not remove Collections icon: %v", err)
		} else if err == nil {
			logger.DebugFn("Removed Collections icon: %s", collectionsIcon)
		}
	}

	// Tools icon - use parent path of Tools since Tools path includes tg5040
	toolsParentDir := filepath.Dir(systemPaths.Tools) // Gets /mnt/SDCARD/Tools
	toolsMediaDir := filepath.Join(toolsParentDir, ".media")
	if _, err := os.Stat(toolsMediaDir); !os.IsNotExist(err) {
		toolsIcon := filepath.Join(toolsMediaDir, "tg5040.png")
		if err := os.Remove(toolsIcon); err != nil && !os.IsNotExist(err) {
			logger.DebugFn("Warning: Could not remove Tools icon: %v", err)
		} else if err == nil {
			logger.DebugFn("Removed Tools icon: %s", toolsIcon)
		}
	}

	// Tool icons
	toolsDir := filepath.Join(systemPaths.Tools)
	toolEntries, err := os.ReadDir(toolsDir)
	if err == nil {
		for _, entry := range toolEntries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			toolName := entry.Name()
			toolMediaDir := filepath.Join(toolsDir, toolName, ".media")

			if _, err := os.Stat(toolMediaDir); os.IsNotExist(err) {
				continue
			}

			toolIcon := filepath.Join(toolMediaDir, toolName + ".png")
			if err := os.Remove(toolIcon); err != nil && !os.IsNotExist(err) {
				logger.DebugFn("Warning: Could not remove %s tool icon: %v", toolName, err)
			} else if err == nil {
				logger.DebugFn("Removed %s tool icon: %s", toolName, toolIcon)
			}
		}
	}

	// Collection icons
	collectionsDir := filepath.Join(systemPaths.Root, "Collections")
	colEntries, err := os.ReadDir(collectionsDir)
	if err == nil {
		for _, entry := range colEntries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			collectionName := entry.Name()
			collectionMediaDir := filepath.Join(collectionsDir, collectionName, ".media")

			if _, err := os.Stat(collectionMediaDir); os.IsNotExist(err) {
				continue
			}

			collectionIcon := filepath.Join(collectionMediaDir, collectionName + ".png")
			if err := os.Remove(collectionIcon); err != nil && !os.IsNotExist(err) {
				logger.DebugFn("Warning: Could not remove %s collection icon: %v", collectionName, err)
			} else if err == nil {
				logger.DebugFn("Removed %s collection icon: %s", collectionName, collectionIcon)
			}
		}
	}

	return nil
}

// cleanupExistingOverlays removes existing overlays before applying new ones
func cleanupExistingOverlays(systemPaths *system.SystemPaths, logger *Logger) error {
    logger.DebugFn("Cleaning up existing overlays")

    // Check for overlays directory
    overlaysDir := filepath.Join(systemPaths.Root, "Overlays")
    if _, err := os.Stat(overlaysDir); os.IsNotExist(err) {
        logger.DebugFn("Overlays directory not found, nothing to clean up")
        return nil
    }

    // List system directories in Overlays
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

        // Remove each overlay file
        for _, file := range overlayFiles {
            if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
                continue
            }

            // Only process PNG files
            if !strings.HasSuffix(strings.ToLower(file.Name()), ".png") {
                continue
            }

            overlayPath := filepath.Join(systemOverlaysPath, file.Name())
            if err := os.Remove(overlayPath); err != nil && !os.IsNotExist(err) {
                logger.DebugFn("Warning: Could not remove overlay %s: %v", file.Name(), err)
            } else if err == nil {
                logger.DebugFn("Removed overlay: %s", overlayPath)
            }
        }

        // Check if system directory is now empty and remove if so
        remainingFiles, _ := os.ReadDir(systemOverlaysPath)
        if len(remainingFiles) == 0 {
            if err := os.Remove(systemOverlaysPath); err != nil {
                logger.DebugFn("Warning: Could not remove empty system overlay directory %s: %v", systemTag, err)
            } else {
                logger.DebugFn("Removed empty system overlay directory: %s", systemOverlaysPath)
            }
        }
    }

    return nil
}