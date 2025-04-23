// src/internal/themes/deconstruction.go
// Implementation of theme deconstruction functionality

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
	"strings"
)

// DeconstructTheme breaks down a theme package into individual component packages
func DeconstructTheme(themeName string) error {
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting theme deconstruction for: %s", themeName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Full path to theme - in the Themes directory
	themePath := filepath.Join(cwd, "Themes", themeName)

	// Validate theme
	manifest, err := ValidateTheme(themePath, logger)
	if err != nil {
		logger.DebugFn("Theme validation failed: %v", err)
		return fmt.Errorf("theme validation failed: %w", err)
	}

	// Track how many components were successfully deconstructed
	componentsDeconstructed := 0

	// Generate export name base from theme name (remove .theme extension)
	exportBaseName := themeName
	if filepath.Ext(exportBaseName) == ".theme" {
		exportBaseName = exportBaseName[:len(exportBaseName)-len(".theme")]
	}

	// Deconstruct wallpapers if present
	if manifest.Content.Wallpapers.Present && manifest.Content.Wallpapers.Count > 0 {
		logger.DebugFn("Deconstructing wallpapers")
		wallpaperName := exportBaseName + ComponentExtension[ComponentWallpaper]

		if err := DeconstructWallpapers(themePath, manifest, wallpaperName, logger); err != nil {
			logger.DebugFn("Warning: Failed to deconstruct wallpapers: %v", err)
		} else {
			componentsDeconstructed++
		}
	}

	// Deconstruct icons if present
	if manifest.Content.Icons.Present &&
	   (manifest.Content.Icons.SystemCount > 0 ||
	    manifest.Content.Icons.ToolCount > 0 ||
	    manifest.Content.Icons.CollectionCount > 0) {
		logger.DebugFn("Deconstructing icons")
		iconName := exportBaseName + ComponentExtension[ComponentIcon]

		if err := DeconstructIcons(themePath, manifest, iconName, logger); err != nil {
			logger.DebugFn("Warning: Failed to deconstruct icons: %v", err)
		} else {
			componentsDeconstructed++
		}
	}

	// Deconstruct overlays if present
	if manifest.Content.Overlays.Present && len(manifest.Content.Overlays.Systems) > 0 {
		logger.DebugFn("Deconstructing overlays")
		overlayName := exportBaseName + ComponentExtension[ComponentOverlay]

		if err := DeconstructOverlays(themePath, manifest, overlayName, logger); err != nil {
			logger.DebugFn("Warning: Failed to deconstruct overlays: %v", err)
		} else {
			componentsDeconstructed++
		}
	}

	// Deconstruct fonts if present
	if manifest.Content.Fonts.Present && (manifest.Content.Fonts.OGReplaced || manifest.Content.Fonts.NextReplaced) {
		logger.DebugFn("Deconstructing fonts")
		fontName := exportBaseName + ComponentExtension[ComponentFont]

		if err := DeconstructFonts(themePath, manifest, fontName, logger); err != nil {
			logger.DebugFn("Warning: Failed to deconstruct fonts: %v", err)
		} else {
			componentsDeconstructed++
		}
	}

	// Deconstruct accent settings if included
	if manifest.Content.Settings.AccentsIncluded {
		logger.DebugFn("Deconstructing accent settings")
		accentName := exportBaseName + ComponentExtension[ComponentAccent]

		if err := DeconstructAccents(themePath, manifest, accentName, logger); err != nil {
			logger.DebugFn("Warning: Failed to deconstruct accent settings: %v", err)
		} else {
			componentsDeconstructed++
		}
	}

	// Deconstruct LED settings if included
	if manifest.Content.Settings.LEDsIncluded {
		logger.DebugFn("Deconstructing LED settings")
		ledName := exportBaseName + ComponentExtension[ComponentLED]

		if err := DeconstructLEDs(themePath, manifest, ledName, logger); err != nil {
			logger.DebugFn("Warning: Failed to deconstruct LED settings: %v", err)
		} else {
			componentsDeconstructed++
		}
	}

	if componentsDeconstructed == 0 {
		return fmt.Errorf("no components were successfully deconstructed from theme: %s", themeName)
	}

	logger.DebugFn("Theme deconstruction completed successfully. %d components extracted.", componentsDeconstructed)

	// Show success message to user
	ui.ShowMessage(fmt.Sprintf("Theme '%s' deconstructed into %d component packages!",
		manifest.ThemeInfo.Name, componentsDeconstructed), "3")

	return nil
}

// DeconstructWallpapers extracts wallpapers from a theme package into a standalone component
func DeconstructWallpapers(themePath string, manifest *ThemeManifest, componentName string, logger *Logger) error {
    logger.DebugFn("Extracting wallpapers from theme to component: %s", componentName)

    // Get current directory
    cwd, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("error getting current directory: %w", err)
    }

    // Create export directory path with .bg extension
    if !strings.HasSuffix(componentName, ComponentExtension[ComponentWallpaper]) {
        componentName = componentName + ComponentExtension[ComponentWallpaper]
    }

    // Path where component will be created (in Exports directory)
    exportPath := filepath.Join(cwd, "Exports", componentName)

    // Create directories for the wallpaper component
    dirPaths := []string{
        exportPath,
        filepath.Join(exportPath, "SystemWallpapers"),
        filepath.Join(exportPath, "CollectionWallpapers"),
    }

    for _, dirPath := range dirPaths {
        if err := os.MkdirAll(dirPath, 0755); err != nil {
            return fmt.Errorf("error creating directory %s: %w", dirPath, err)
        }
    }

    // Create minimal component manifest with author from theme
    manifestObj, err := CreateMinimalComponentManifest(ComponentWallpaper, componentName, manifest.ThemeInfo.Author)
    if err != nil {
        return fmt.Errorf("error creating wallpaper manifest: %w", err)
    }

    wallpaperManifest := manifestObj.(*WallpaperManifest)

    // Process each wallpaper mapping from the theme manifest
    // Copy the files but don't populate the component manifest with mappings
    for _, mapping := range manifest.PathMappings.Wallpapers {
        srcPath := filepath.Join(themePath, mapping.ThemePath)

        // Skip non-existent files
        if _, err := os.Stat(srcPath); os.IsNotExist(err) {
            logger.DebugFn("Warning: Source file does not exist: %s", srcPath)
            continue
        }

        // Determine destination path in component package
        // The ThemePath is expected to be like "Wallpapers/SystemWallpapers/Name.png"
        // We strip the initial "Wallpapers/" to get the correct path in our component package
        relativePath := mapping.ThemePath
        if strings.HasPrefix(relativePath, "Wallpapers/") {
            relativePath = relativePath[len("Wallpapers/"):]
        }

        dstPath := filepath.Join(exportPath, relativePath)

        // Ensure destination directory exists
        dstDir := filepath.Dir(dstPath)
        if err := os.MkdirAll(dstDir, 0755); err != nil {
            logger.DebugFn("Warning: Could not create directory %s: %v", dstDir, err)
            continue
        }

        // Copy the file
        if err := CopyFile(srcPath, dstPath); err != nil {
            logger.DebugFn("Warning: Could not copy wallpaper: %v", err)
            continue
        }

        logger.DebugFn("Copied wallpaper: %s", relativePath)
    }

    // Create a preview image - try to use a system wallpaper as preview
    previewPath := filepath.Join(exportPath, "preview.png")

    // Look for Recently Played wallpaper first
    recentlyPlayedPath := filepath.Join(exportPath, "SystemWallpapers", "Recently Played.png")
    if _, err := os.Stat(recentlyPlayedPath); err == nil {
        if err := CopyFile(recentlyPlayedPath, previewPath); err != nil {
            logger.DebugFn("Warning: Could not copy preview image: %v", err)
            // Create default preview as fallback
            if err := CreateDefaultPreviewImage(previewPath, ComponentWallpaper); err != nil {
                logger.DebugFn("Warning: Could not create default preview: %v", err)
            }
        }
    } else {
        // If no Recently Played, try to find any wallpaper
        systemWallpapersDir := filepath.Join(exportPath, "SystemWallpapers")
        entries, err := os.ReadDir(systemWallpapersDir)
        if err == nil && len(entries) > 0 {
            // Use the first wallpaper found
            for _, entry := range entries {
                if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
                    candidatePath := filepath.Join(systemWallpapersDir, entry.Name())
                    if err := CopyFile(candidatePath, previewPath); err != nil {
                        logger.DebugFn("Warning: Could not copy preview image: %v", err)
                    } else {
                        break
                    }
                }
            }
        }

        // If no preview yet, create default
        if _, err := os.Stat(previewPath); os.IsNotExist(err) {
            if err := CreateDefaultPreviewImage(previewPath, ComponentWallpaper); err != nil {
                logger.DebugFn("Warning: Could not create default preview: %v", err)
            }
        }
    }

    // Write the component manifest
    if err := WriteComponentManifest(exportPath, wallpaperManifest); err != nil {
        return fmt.Errorf("error writing wallpaper manifest: %w", err)
    }

    logger.DebugFn("Wallpaper component extraction completed")
    return nil
}

// DeconstructIcons extracts icons from a theme package into a standalone component
func DeconstructIcons(themePath string, manifest *ThemeManifest, componentName string, logger *Logger) error {
	logger.DebugFn("Extracting icons from theme to component: %s", componentName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .icon extension
	if !strings.HasSuffix(componentName, ComponentExtension[ComponentIcon]) {
		componentName = componentName + ComponentExtension[ComponentIcon]
	}

	// Path where component will be created (in Exports directory)
	exportPath := filepath.Join(cwd, "Exports", componentName)

	// Create directories for the icon component
	dirPaths := []string{
		exportPath,
		filepath.Join(exportPath, "SystemIcons"),
		filepath.Join(exportPath, "ToolIcons"),
		filepath.Join(exportPath, "CollectionIcons"),
	}

	for _, dirPath := range dirPaths {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dirPath, err)
		}
	}

	// Create component manifest
	manifestObj, err := CreateComponentManifest(ComponentIcon, componentName)
	if err != nil {
		return fmt.Errorf("error creating icon manifest: %w", err)
	}

	iconManifest := manifestObj.(*IconManifest)

	// Copy over author information from theme manifest
	iconManifest.ComponentInfo.Author = manifest.ThemeInfo.Author

	// Process each icon mapping from the theme manifest
	for _, mapping := range manifest.PathMappings.Icons {
		srcPath := filepath.Join(themePath, mapping.ThemePath)

		// Skip non-existent files
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.DebugFn("Warning: Source file does not exist: %s", srcPath)
			continue
		}

		// Determine destination path in component package
		// The ThemePath is expected to be like "Icons/SystemIcons/Name.png"
		// We strip the initial "Icons/" to get the correct path in our component package
		relativePath := mapping.ThemePath
		if strings.HasPrefix(relativePath, "Icons/") {
			relativePath = relativePath[len("Icons/"):]
		}

		dstPath := filepath.Join(exportPath, relativePath)

		// Ensure destination directory exists
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.DebugFn("Warning: Could not create directory %s: %v", dstDir, err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.DebugFn("Warning: Could not copy icon: %v", err)
			continue
		}

		// Add to component manifest
		iconManifest.PathMappings = append(
			iconManifest.PathMappings,
			PathMapping{
				ThemePath:  relativePath,
				SystemPath: mapping.SystemPath,
				Metadata:   mapping.Metadata,
			},
		)

		// Update component manifest content section and counters
		filename := filepath.Base(relativePath)
		if strings.HasPrefix(relativePath, "SystemIcons/") {
			iconManifest.Content.SystemIcons = append(
				iconManifest.Content.SystemIcons,
				filename,
			)
			iconManifest.Content.SystemCount++
		} else if strings.HasPrefix(relativePath, "ToolIcons/") {
			iconManifest.Content.ToolIcons = append(
				iconManifest.Content.ToolIcons,
				filename,
			)
			iconManifest.Content.ToolCount++
		} else if strings.HasPrefix(relativePath, "CollectionIcons/") {
			iconManifest.Content.CollectionIcons = append(
				iconManifest.Content.CollectionIcons,
				filename,
			)
			iconManifest.Content.CollectionCount++
		}

		logger.DebugFn("Copied icon: %s", relativePath)
	}

	// Create a preview image - try to use a system icon as preview
	previewPath := filepath.Join(exportPath, "preview.png")
	if iconManifest.Content.SystemCount > 0 ||
	   iconManifest.Content.ToolCount > 0 ||
	   iconManifest.Content.CollectionCount > 0 {
		var previewSource string

		// Try to find a good candidate for the preview
		// First try Collections icon since it usually has a good, representative icon
		for _, mapping := range iconManifest.PathMappings {
			if strings.HasSuffix(mapping.ThemePath, "Collections.png") {
				previewSource = filepath.Join(exportPath, mapping.ThemePath)
				break
			}
		}

		// If no preview found yet, use the first icon
		if previewSource == "" && len(iconManifest.PathMappings) > 0 {
			previewSource = filepath.Join(exportPath, iconManifest.PathMappings[0].ThemePath)
		}

		// Copy the preview
		if previewSource != "" {
			if err := CopyFile(previewSource, previewPath); err != nil {
				logger.DebugFn("Warning: Could not copy preview image: %v", err)

				// Create default preview as fallback
				if err := CreateDefaultPreviewImage(previewPath, ComponentIcon); err != nil {
					logger.DebugFn("Warning: Could not create default preview: %v", err)
				}
			}
		} else {
			// Create default preview
			if err := CreateDefaultPreviewImage(previewPath, ComponentIcon); err != nil {
				logger.DebugFn("Warning: Could not create default preview: %v", err)
			}
		}
	} else {
		// No icons found, create a default preview
		if err := CreateDefaultPreviewImage(previewPath, ComponentIcon); err != nil {
			logger.DebugFn("Warning: Could not create default preview: %v", err)
		}
	}

	// Write the component manifest
	if err := WriteComponentManifest(exportPath, iconManifest); err != nil {
		return fmt.Errorf("error writing icon manifest: %w", err)
	}

	logger.DebugFn("Icon component extraction completed: %d system, %d tool, %d collection icons",
		iconManifest.Content.SystemCount,
		iconManifest.Content.ToolCount,
		iconManifest.Content.CollectionCount)
	return nil
}

// DeconstructOverlays extracts overlays from a theme package into a standalone component
func DeconstructOverlays(themePath string, manifest *ThemeManifest, componentName string, logger *Logger) error {
	logger.DebugFn("Extracting overlays from theme to component: %s", componentName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .over extension
	if !strings.HasSuffix(componentName, ComponentExtension[ComponentOverlay]) {
		componentName = componentName + ComponentExtension[ComponentOverlay]
	}

	// Path where component will be created (in Exports directory)
	exportPath := filepath.Join(cwd, "Exports", componentName)

	// Create the root directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", exportPath, err)
	}

	// Create the Systems directory
	systemsDir := filepath.Join(exportPath, "Systems")
	if err := os.MkdirAll(systemsDir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", systemsDir, err)
	}

	// Create component manifest
	manifestObj, err := CreateComponentManifest(ComponentOverlay, componentName)
	if err != nil {
		return fmt.Errorf("error creating overlay manifest: %w", err)
	}

	overlayManifest := manifestObj.(*OverlayManifest)

	// Copy over author information from theme manifest
	overlayManifest.ComponentInfo.Author = manifest.ThemeInfo.Author

	// Copy over the systems list
	overlayManifest.Content.Systems = append([]string{}, manifest.Content.Overlays.Systems...)

	// Process each overlay mapping from the theme manifest
	for _, mapping := range manifest.PathMappings.Overlays {
		srcPath := filepath.Join(themePath, mapping.ThemePath)

		// Skip non-existent files
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.DebugFn("Warning: Source file does not exist: %s", srcPath)
			continue
		}

		// Determine destination path in component package
		// The ThemePath is expected to be like "Overlays/SYSTEM/file.png"
		relativePath := mapping.ThemePath
		if strings.HasPrefix(relativePath, "Overlays/") {
			relativePath = "Systems/" + relativePath[len("Overlays/"):]
		}

		dstPath := filepath.Join(exportPath, relativePath)

		// Ensure destination directory exists
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.DebugFn("Warning: Could not create directory %s: %v", dstDir, err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.DebugFn("Warning: Could not copy overlay: %v", err)
			continue
		}

		// Add to component manifest
		overlayManifest.PathMappings = append(
			overlayManifest.PathMappings,
			PathMapping{
				ThemePath:  relativePath,
				SystemPath: mapping.SystemPath,
				Metadata:   mapping.Metadata,
			},
		)

		logger.DebugFn("Copied overlay: %s", relativePath)
	}

	// Create a default preview image
	previewPath := filepath.Join(exportPath, "preview.png")
	if err := CreateDefaultPreviewImage(previewPath, ComponentOverlay); err != nil {
		logger.DebugFn("Warning: Could not create default preview: %v", err)
	}

	// Write the component manifest
	if err := WriteComponentManifest(exportPath, overlayManifest); err != nil {
		return fmt.Errorf("error writing overlay manifest: %w", err)
	}

	logger.DebugFn("Overlay component extraction completed for %d systems", len(overlayManifest.Content.Systems))
	return nil
}

// DeconstructFonts extracts fonts from a theme package into a standalone component
func DeconstructFonts(themePath string, manifest *ThemeManifest, componentName string, logger *Logger) error {
	logger.DebugFn("Extracting fonts from theme to component: %s", componentName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .font extension
	if !strings.HasSuffix(componentName, ComponentExtension[ComponentFont]) {
		componentName = componentName + ComponentExtension[ComponentFont]
	}

	// Path where component will be created (in Exports directory)
	exportPath := filepath.Join(cwd, "Exports", componentName)

	// Create the root directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", exportPath, err)
	}

	// Create component manifest
	manifestObj, err := CreateComponentManifest(ComponentFont, componentName)
	if err != nil {
		return fmt.Errorf("error creating font manifest: %w", err)
	}

	fontManifest := manifestObj.(*FontManifest)

	// Copy over author information from theme manifest
	fontManifest.ComponentInfo.Author = manifest.ThemeInfo.Author

	// Copy content flags
	fontManifest.Content.OGReplaced = manifest.Content.Fonts.OGReplaced
	fontManifest.Content.NextReplaced = manifest.Content.Fonts.NextReplaced

	// Process fonts from the theme
	fontPaths := []string{
		"Fonts/OG.ttf",
		"Fonts/Next.ttf",
		"Fonts/OG.backup.ttf",
		"Fonts/Next.backup.ttf",
	}

	fontNames := []string{
		"OG",
		"Next",
		"OG.backup",
		"Next.backup",
	}

    // Define system paths for fonts - CORRECTED PATHS
    systemPaths := map[string]string{
        "OG":          "/mnt/SDCARD/.system/res/font2.ttf",
        "OG.backup":   "/mnt/SDCARD/.system/res/font2.backup.ttf",  // Corrected extension
        "Next":        "/mnt/SDCARD/.system/res/font1.ttf",
        "Next.backup": "/mnt/SDCARD/.system/res/font1.backup.ttf",  // Corrected extension
    }

	for i, fontPath := range fontPaths {
		fontName := fontNames[i]

		srcPath := filepath.Join(themePath, fontPath)
		dstPath := filepath.Join(exportPath, fontName + ".ttf")

		// Skip non-existent files
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.DebugFn("Font file does not exist in theme: %s", srcPath)
			continue
		}

		// Copy the font file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.DebugFn("Warning: Could not copy font %s: %v", fontName, err)
			continue
		}

		// Add to component manifest
		fontManifest.PathMappings[fontName] = PathMapping{
			ThemePath:  fontName + ".ttf",
			SystemPath: systemPaths[fontName],
		}

		logger.DebugFn("Copied font: %s", fontName)
	}

	// Create a default preview image
	previewPath := filepath.Join(exportPath, "preview.png")
	if err := CreateDefaultPreviewImage(previewPath, ComponentFont); err != nil {
		logger.DebugFn("Warning: Could not create default preview: %v", err)
	}

	// Write the component manifest
	if err := WriteComponentManifest(exportPath, fontManifest); err != nil {
		return fmt.Errorf("error writing font manifest: %w", err)
	}

	logger.DebugFn("Font component extraction completed")
	return nil
}

// DeconstructAccents extracts accent settings from a theme package into a standalone component
func DeconstructAccents(themePath string, manifest *ThemeManifest, componentName string, logger *Logger) error {
	logger.DebugFn("Extracting accent settings from theme to component: %s", componentName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .acc extension
	if !strings.HasSuffix(componentName, ComponentExtension[ComponentAccent]) {
		componentName = componentName + ComponentExtension[ComponentAccent]
	}

	// Path where component will be created (in Exports directory)
	exportPath := filepath.Join(cwd, "Exports", componentName)

	// Create the root directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", exportPath, err)
	}

	// Create component manifest
	manifestObj, err := CreateComponentManifest(ComponentAccent, componentName)
	if err != nil {
		return fmt.Errorf("error creating accent manifest: %w", err)
	}

	accentManifest := manifestObj.(*AccentManifest)

	// Copy over author information from theme manifest
	accentManifest.ComponentInfo.Author = manifest.ThemeInfo.Author

	// Copy accent colors from theme manifest
	accentManifest.AccentColors.Color1 = manifest.AccentColors.Color1
	accentManifest.AccentColors.Color2 = manifest.AccentColors.Color2
	accentManifest.AccentColors.Color3 = manifest.AccentColors.Color3
	accentManifest.AccentColors.Color4 = manifest.AccentColors.Color4
	accentManifest.AccentColors.Color5 = manifest.AccentColors.Color5
	accentManifest.AccentColors.Color6 = manifest.AccentColors.Color6

	// Create a default preview image
	previewPath := filepath.Join(exportPath, "preview.png")
	if err := CreateDefaultPreviewImage(previewPath, ComponentAccent); err != nil {
		logger.DebugFn("Warning: Could not create default preview: %v", err)
	}

	// Write the component manifest
	if err := WriteComponentManifest(exportPath, accentManifest); err != nil {
		return fmt.Errorf("error writing accent manifest: %w", err)
	}

	logger.DebugFn("Accent settings component extraction completed")
	return nil
}

// DeconstructLEDs extracts LED settings from a theme package into a standalone component
func DeconstructLEDs(themePath string, manifest *ThemeManifest, componentName string, logger *Logger) error {
	logger.DebugFn("Extracting LED settings from theme to component: %s", componentName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .led extension
	if !strings.HasSuffix(componentName, ComponentExtension[ComponentLED]) {
		componentName = componentName + ComponentExtension[ComponentLED]
	}

	// Path where component will be created (in Exports directory)
	exportPath := filepath.Join(cwd, "Exports", componentName)

	// Create the root directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", exportPath, err)
	}

	// Create component manifest
	manifestObj, err := CreateComponentManifest(ComponentLED, componentName)
	if err != nil {
		return fmt.Errorf("error creating LED manifest: %w", err)
	}

	ledManifest := manifestObj.(*LEDManifest)

	// Copy over author information from theme manifest
	ledManifest.ComponentInfo.Author = manifest.ThemeInfo.Author

	// Copy LED settings from theme manifest
	ledManifest.LEDSettings.F1Key = manifest.LEDSettings.F1Key
	ledManifest.LEDSettings.F2Key = manifest.LEDSettings.F2Key
	ledManifest.LEDSettings.TopBar = manifest.LEDSettings.TopBar
	ledManifest.LEDSettings.LRTriggers = manifest.LEDSettings.LRTriggers

	// Note: LEDs don't have a preview image by design

	// Write the component manifest
	if err := WriteComponentManifest(exportPath, ledManifest); err != nil {
		return fmt.Errorf("error writing LED manifest: %w", err)
	}

	logger.DebugFn("LED settings component extraction completed")
	return nil
}