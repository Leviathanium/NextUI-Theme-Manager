// src/internal/themes/import.go
// Simplified implementation of theme import functionality

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"nextui-themes/internal/ui"
	"strings"
	"regexp"
	"strconv"
)

// ImportTheme imports a theme package
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

    // Full path to theme - look in Imports directory
    themePath := filepath.Join(cwd, "Themes", "Imports", themeName)

    // Validate theme
    manifest, err := ValidateTheme(themePath, logger)
    if err != nil {
        logger.DebugFn("Theme validation failed: %v", err)
        return fmt.Errorf("theme validation failed: %w", err)
    }

    // Update manifest based on theme content
    if err := UpdateManifestFromThemeContent(themePath, manifest, logger); err != nil {
        logger.DebugFn("Warning: Error updating manifest from theme content: %v", err)
        // Continue anyway with the original manifest
    }

    // Get system paths
    systemPaths, err := system.GetSystemPaths()
    if err != nil {
        logger.DebugFn("Error getting system paths: %v", err)
        return fmt.Errorf("error getting system paths: %w", err)
    }

    // Clean up existing components
    if err := cleanupExistingComponents(manifest, systemPaths, logger); err != nil {
        logger.DebugFn("Warning: Error cleaning up existing components: %v", err)
        // Continue with import anyway
    }

    // Apply theme components based on the (now updated) manifest
    if err := importThemeFiles(themePath, manifest, logger); err != nil {
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

// importThemeFiles copies all files from the theme to the system based on path mappings
func importThemeFiles(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logger.DebugFn("Error getting system paths: %v", err)
		// Continue anyway with just the path mappings
	}

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

	// Process icon mappings
	for _, mapping := range manifest.PathMappings.Icons {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
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
func UpdateManifestFromThemeContent(themePath string, manifest *ThemeManifest, logger *Logger) error {
    // Get system paths for mapping
    systemPaths, err := system.GetSystemPaths()
    if err != nil {
        logger.DebugFn("Warning: Error getting system paths: %v", err)
        // Continue anyway - we can still use naming conventions
    }

    // Update wallpapers if present
    if manifest.Content.Wallpapers.Present {
        if err := updateWallpaperMappings(themePath, manifest, systemPaths, logger); err != nil {
            logger.DebugFn("Warning: Error updating wallpaper mappings: %v", err)
        }
    }

    // Update icons if present
    if manifest.Content.Icons.Present {
        if err := updateIconMappings(themePath, manifest, systemPaths, logger); err != nil {
            logger.DebugFn("Warning: Error updating icon mappings: %v", err)
        }
    }

    // Update settings if present
    if manifest.Content.Settings.AccentsIncluded {
        if err := updateAccentSettings(themePath, manifest, logger); err != nil {
            logger.DebugFn("Warning: Error updating accent settings: %v", err)
        }
    }

    if manifest.Content.Settings.LEDsIncluded {
        if err := updateLEDSettings(themePath, manifest, logger); err != nil {
            logger.DebugFn("Warning: Error updating LED settings: %v", err)
        }
    }

    // Write updated manifest back to file
    return WriteManifest(themePath, manifest, logger)
}

// updateIconMappings scans icons in the theme and updates manifest mappings
func updateIconMappings(themePath string, manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) error {
    iconsDir := filepath.Join(themePath, "Icons")

    // Check if directory exists
    if _, err := os.Stat(iconsDir); os.IsNotExist(err) {
        return nil // No icons directory, nothing to update
    }

    // Create a map of existing mappings for quick lookup
    existingMappings := make(map[string]bool)
    for _, mapping := range manifest.PathMappings.Icons {
        existingMappings[mapping.ThemePath] = true
    }

    // Subdirectories to scan
    subdirs := []string{
        "SystemIcons",
        "ToolIcons",
        "CollectionIcons",
    }

    // Regular expression to extract system tag from filenames
    tagRegex := regexp.MustCompile(`\((.*?)\)`)

    for _, subdir := range subdirs {
        fullSubdir := filepath.Join(iconsDir, subdir)

        // Skip if subdirectory doesn't exist
        if _, err := os.Stat(fullSubdir); os.IsNotExist(err) {
            continue
        }

        // List icon files
        entries, err := os.ReadDir(fullSubdir)
        if err != nil {
            logger.DebugFn("Warning: Error reading icons directory %s: %v", subdir, err)
            continue
        }

        for _, entry := range entries {
            if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
                continue
            }

            // Check if file has a PNG extension
            if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
                continue
            }

            themePath := filepath.Join("Icons", subdir, entry.Name())

            // Skip if this file is already in mappings
            if existingMappings[themePath] {
                continue
            }

            // Determine where this file should go based on naming
            var systemPath string
            var metadata map[string]string

            switch subdir {
            case "SystemIcons":
                // Special case handling for predefined names
                switch strings.TrimSuffix(entry.Name(), ".png") {
                case "Recently Played", "Recently-Played":
                    systemPath = filepath.Join(systemPaths.RecentlyPlayed, ".media", "icon.png")
                    metadata = map[string]string{
                        "SystemName": "Recently Played",
                        "IconType": "System",
                    }

                case "Tools":
                    systemPath = filepath.Join(systemPaths.Tools, ".media", "icon.png")
                    metadata = map[string]string{
                        "SystemName": "Tools",
                        "IconType": "System",
                    }

                case "Collections":
                    systemPath = filepath.Join(systemPaths.Root, "Collections", ".media", "icon.png")
                    metadata = map[string]string{
                        "SystemName": "Collections",
                        "IconType": "System",
                    }

                default:
                    // Check for system tag in filename
                    matches := tagRegex.FindStringSubmatch(entry.Name())
                    if len(matches) >= 2 {
                        systemTag := matches[1]
                        systemName := strings.TrimSuffix(strings.Split(entry.Name(), "(")[0], " ")

                        // Find matching system by tag
                        var systemFound bool
                        for _, system := range systemPaths.Systems {
                            if system.Tag == systemTag {
                                systemPath = filepath.Join(system.MediaPath, "icon.png")
                                metadata = map[string]string{
                                    "SystemName": systemName,
                                    "SystemTag": systemTag,
                                    "IconType": "System",
                                }
                                systemFound = true
                                break
                            }
                        }

                        // If system not found in paths, create a default path
                        if !systemFound && systemTag != "" {
                            systemPath = filepath.Join(systemPaths.Roms, fmt.Sprintf("%s (%s)", systemName, systemTag), ".media", "icon.png")
                            metadata = map[string]string{
                                "SystemName": systemName,
                                "SystemTag": systemTag,
                                "IconType": "System",
                            }
                        }
                    }
                }

            case "ToolIcons":
                toolName := strings.TrimSuffix(entry.Name(), ".png")
                systemPath = filepath.Join(systemPaths.Tools, toolName, ".media", "icon.png")
                metadata = map[string]string{
                    "ToolName": toolName,
                    "IconType": "Tool",
                }

            case "CollectionIcons":
                collectionName := strings.TrimSuffix(entry.Name(), ".png")
                systemPath = filepath.Join(systemPaths.Root, "Collections", collectionName, ".media", "icon.png")
                metadata = map[string]string{
                    "CollectionName": collectionName,
                    "IconType": "Collection",
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

                // Update counters based on icon type
                switch subdir {
                case "SystemIcons":
                    manifest.Content.Icons.SystemCount++
                case "ToolIcons":
                    manifest.Content.Icons.ToolCount++
                case "CollectionIcons":
                    manifest.Content.Icons.CollectionCount++
                }

                logger.DebugFn("Added mapping for icon: %s -> %s", themePath, systemPath)
            } else {
                logger.DebugFn("Could not determine system path for icon: %s", entry.Name())
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

// updateWallpaperMappings scans wallpapers in the theme and updates manifest mappings
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
                        "SystemName": "Root",
                        "WallpaperType": "Main",
                    }

                case "Root-Media":
                    systemPath = filepath.Join(systemPaths.Root, ".media", "bg.png")
                    metadata = map[string]string{
                        "SystemName": "Root",
                        "WallpaperType": "Media",
                    }

                case "Recently Played":
                    systemPath = filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
                    metadata = map[string]string{
                        "SystemName": "Recently Played",
                        "WallpaperType": "Media",
                    }

                case "Tools":
                    systemPath = filepath.Join(systemPaths.Tools, ".media", "bg.png")
                    metadata = map[string]string{
                        "SystemName": "Tools",
                        "WallpaperType": "Media",
                    }

                case "Collections":
                    systemPath = filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")
                    metadata = map[string]string{
                        "SystemName": "Collections",
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
                                    "SystemName": systemName,
                                    "SystemTag": systemTag,
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
                                "SystemName": systemName,
                                "SystemTag": systemTag,
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

                // Extract collection name
                collectionName := strings.TrimSuffix(entry.Name(), ".png")
                systemPath := filepath.Join(systemPaths.Root, "Collections", collectionName, ".media", "bg.png")
                metadata := map[string]string{
                    "CollectionName": collectionName,
                    "WallpaperType": "Collection",
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

// applyAccentSettings applies accent color settings from manifest
func applyAccentSettings(manifest *ThemeManifest, logger *Logger) error {
    // Create content for minuisettings.txt
    var content strings.Builder
    content.WriteString(fmt.Sprintf("color1=%s\n", manifest.AccentColors.Color1))
    content.WriteString(fmt.Sprintf("color2=%s\n", manifest.AccentColors.Color2))
    content.WriteString(fmt.Sprintf("color3=%s\n", manifest.AccentColors.Color3))
    content.WriteString(fmt.Sprintf("color4=%s\n", manifest.AccentColors.Color4))
    content.WriteString(fmt.Sprintf("color5=%s\n", manifest.AccentColors.Color5))
    content.WriteString(fmt.Sprintf("color6=%s\n", manifest.AccentColors.Color6))

    // Get path to settings file - FIXED PATH
    settingsPath := "/mnt/SDCARD/.userdata/shared/minuisettings.txt"

    // Write settings to file
    if err := os.WriteFile(settingsPath, []byte(content.String()), 0644); err != nil {
        return fmt.Errorf("error writing accent settings: %w", err)
    }

    logger.DebugFn("Applied accent settings to %s", settingsPath)
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

// cleanupExistingComponents removes existing components that aren't in the new theme
func cleanupExistingComponents(manifest *ThemeManifest, systemPaths *system.SystemPaths, logger *Logger) error {
    logger.DebugFn("Starting cleanup of existing components")

    // If theme doesn't include wallpapers, clean up existing wallpapers
    if !manifest.Content.Wallpapers.Present {
        logger.DebugFn("Theme doesn't include wallpapers - cleaning up existing wallpapers")

        // Root wallpaper
        rootBg := filepath.Join(systemPaths.Root, "bg.png")
        if err := os.Remove(rootBg); err != nil && !os.IsNotExist(err) {
            logger.DebugFn("Warning: Could not remove root wallpaper: %v", err)
        } else {
            logger.DebugFn("Removed root wallpaper: %s", rootBg)
        }

        // Root media wallpaper
        rootMediaBg := filepath.Join(systemPaths.Root, ".media", "bg.png")
        if err := os.Remove(rootMediaBg); err != nil && !os.IsNotExist(err) {
            logger.DebugFn("Warning: Could not remove root media wallpaper: %v", err)
        } else {
            logger.DebugFn("Removed root media wallpaper: %s", rootMediaBg)
        }

        // Recently Played wallpaper
        rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
        if err := os.Remove(rpBg); err != nil && !os.IsNotExist(err) {
            logger.DebugFn("Warning: Could not remove Recently Played wallpaper: %v", err)
        } else {
            logger.DebugFn("Removed Recently Played wallpaper: %s", rpBg)
        }

        // Tools wallpaper
        toolsBg := filepath.Join(systemPaths.Tools, ".media", "bg.png")
        if err := os.Remove(toolsBg); err != nil && !os.IsNotExist(err) {
            logger.DebugFn("Warning: Could not remove Tools wallpaper: %v", err)
        } else {
            logger.DebugFn("Removed Tools wallpaper: %s", toolsBg)
        }

        // Collections wallpaper
        collectionsBg := filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")
        if err := os.Remove(collectionsBg); err != nil && !os.IsNotExist(err) {
            logger.DebugFn("Warning: Could not remove Collections wallpaper: %v", err)
        } else {
            logger.DebugFn("Removed Collections wallpaper: %s", collectionsBg)
        }

        // System wallpapers
        for _, system := range systemPaths.Systems {
            systemBg := filepath.Join(system.MediaPath, "bg.png")
            if err := os.Remove(systemBg); err != nil && !os.IsNotExist(err) {
                logger.DebugFn("Warning: Could not remove %s wallpaper: %v", system.Name, err)
            } else {
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
                } else {
                    logger.DebugFn("Removed %s collection wallpaper: %s", collectionName, collectionBg)
                }
            }
        }
    } else {
        logger.DebugFn("Theme includes wallpapers - keeping existing wallpapers until they're replaced")
    }

    // If theme doesn't include icons, clean up existing icons
    if !manifest.Content.Icons.Present {
        logger.DebugFn("Theme doesn't include icons - cleaning up existing icons")

        // System icons
        for _, system := range systemPaths.Systems {
            systemIcon := filepath.Join(system.MediaPath, "icon.png")
            if err := os.Remove(systemIcon); err != nil && !os.IsNotExist(err) {
                logger.DebugFn("Warning: Could not remove %s icon: %v", system.Name, err)
            } else {
                logger.DebugFn("Removed %s icon: %s", system.Name, systemIcon)
            }
        }

        // Recently Played icon
        rpIcon := filepath.Join(systemPaths.RecentlyPlayed, ".media", "icon.png")
        if err := os.Remove(rpIcon); err != nil && !os.IsNotExist(err) {
            logger.DebugFn("Warning: Could not remove Recently Played icon: %v", err)
        } else {
            logger.DebugFn("Removed Recently Played icon: %s", rpIcon)
        }

        // Tools icon
        toolsIcon := filepath.Join(systemPaths.Tools, ".media", "icon.png")
        if err := os.Remove(toolsIcon); err != nil && !os.IsNotExist(err) {
            logger.DebugFn("Warning: Could not remove Tools icon: %v", err)
        } else {
            logger.DebugFn("Removed Tools icon: %s", toolsIcon)
        }

        // Collections icon
        collectionsIcon := filepath.Join(systemPaths.Root, "Collections", ".media", "icon.png")
        if err := os.Remove(collectionsIcon); err != nil && !os.IsNotExist(err) {
            logger.DebugFn("Warning: Could not remove Collections icon: %v", err)
        } else {
            logger.DebugFn("Removed Collections icon: %s", collectionsIcon)
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
                toolIcon := filepath.Join(toolsDir, toolName, ".media", "icon.png")
                if err := os.Remove(toolIcon); err != nil && !os.IsNotExist(err) {
                    logger.DebugFn("Warning: Could not remove %s tool icon: %v", toolName, err)
                } else {
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
                collectionIcon := filepath.Join(collectionsDir, collectionName, ".media", "icon.png")
                if err := os.Remove(collectionIcon); err != nil && !os.IsNotExist(err) {
                    logger.DebugFn("Warning: Could not remove %s collection icon: %v", collectionName, err)
                } else {
                    logger.DebugFn("Removed %s collection icon: %s", collectionName, collectionIcon)
                }
            }
        }
    } else {
        logger.DebugFn("Theme includes icons - keeping existing icons until they're replaced")
    }

    // If theme doesn't include overlays, clean up existing overlays
    if !manifest.Content.Overlays.Present {
        logger.DebugFn("Theme doesn't include overlays - cleaning up existing overlays")

        // Delete all overlays in the Overlays directory
        overlaysDir := filepath.Join(systemPaths.Root, "Overlays")
        if err := os.RemoveAll(overlaysDir); err != nil && !os.IsNotExist(err) {
            logger.DebugFn("Warning: Could not remove Overlays directory: %v", err)
        } else {
            logger.DebugFn("Removed Overlays directory: %s", overlaysDir)
            // Recreate the directory
            if err := os.MkdirAll(overlaysDir, 0755); err != nil {
                logger.DebugFn("Warning: Could not recreate Overlays directory: %v", err)
            }
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