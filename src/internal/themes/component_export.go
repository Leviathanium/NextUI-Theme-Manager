// src/internal/themes/component_export.go
// Implements export functionality for individual theme components

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
    "regexp"
    "strconv"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"nextui-themes/internal/ui"
)

// CreateDefaultPreviewImage creates a default preview image with text
func CreateDefaultPreviewImage(outputPath string, componentType string) error {
	// Instead of looking for a placeholder image that doesn't exist,
	// we'll create a simple blank one here

	// For now, just log that a preview image is missing and return success
	logging.LogDebug("Creating a blank preview for %s", componentType)

	// Create a blank file as the preview (will show up as blank in the UI)
	// This is preferable to creating empty directories that confuse users
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating blank preview: %w", err)
	}
	f.Close()

	return nil
}

// ExportWallpapers exports current wallpapers as a .bg component package
func ExportWallpapers(name string) error {
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting wallpaper export: %s", name)

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .bg extension
	if !strings.HasSuffix(name, ComponentExtension[ComponentWallpaper]) {
		name = name + ComponentExtension[ComponentWallpaper]
	}

	exportPath := filepath.Join(cwd, "Exports", name)

	// Create directories
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

	// Create component manifest
	manifest, err := CreateComponentManifest(ComponentWallpaper, name)
	if err != nil {
		return fmt.Errorf("error creating wallpaper manifest: %w", err)
	}

	wallpaperManifest := manifest.(*WallpaperManifest)

	// Get system paths for copying wallpapers
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Export wallpapers (similar to exportWallpapers in export.go but with different directory structure)
	// Function to export a wallpaper and add it to the manifest
	exportWallpaper := func(srcPath, relativeDstDir, filename string, metadata map[string]string) error {
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return nil // Skip if source doesn't exist
		}

		dstDir := filepath.Join(exportPath, relativeDstDir)
		dstPath := filepath.Join(dstDir, filename)

		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.DebugFn("Warning: Could not copy wallpaper %s: %v", srcPath, err)
			return err
		}

		// Add to manifest
		themePath := filepath.Join(relativeDstDir, filename)
		wallpaperManifest.PathMappings = append(
			wallpaperManifest.PathMappings,
			PathMapping{
				ThemePath:  themePath,
				SystemPath: srcPath,
				Metadata:   metadata,
			},
		)

		// Add to content list
		if strings.HasPrefix(relativeDstDir, "SystemWallpapers") {
			wallpaperManifest.Content.SystemWallpapers = append(
				wallpaperManifest.Content.SystemWallpapers,
				filename,
			)
		} else {
			wallpaperManifest.Content.CollectionWallpapers = append(
				wallpaperManifest.Content.CollectionWallpapers,
				filename,
			)
		}

		wallpaperManifest.Content.Count++
		logger.DebugFn("Exported wallpaper: %s", dstPath)
		return nil
	}

	// Export root wallpaper
	rootBg := filepath.Join(systemPaths.Root, "bg.png")
	exportWallpaper(rootBg, "SystemWallpapers", "Root.png", map[string]string{
		"SystemName":    "Root",
		"WallpaperType": "Main",
	})

	// Export root media wallpaper
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", "bg.png")
	exportWallpaper(rootMediaBg, "SystemWallpapers", "Root-Media.png", map[string]string{
		"SystemName":    "Root",
		"WallpaperType": "Media",
	})

	// Export Recently Played wallpaper
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
	exportWallpaper(rpBg, "SystemWallpapers", "Recently Played.png", map[string]string{
		"SystemName":    "Recently Played",
		"WallpaperType": "Media",
	})

	// Export Tools wallpaper
	toolsBg := filepath.Join(systemPaths.Tools, ".media", "bg.png")
	exportWallpaper(toolsBg, "SystemWallpapers", "Tools.png", map[string]string{
		"SystemName":    "Tools",
		"WallpaperType": "Media",
	})

	// Export Collections wallpaper
	collectionsBg := filepath.Join(systemPaths.Root, "Collections", ".media", "bg.png")
	exportWallpaper(collectionsBg, "SystemWallpapers", "Collections.png", map[string]string{
		"SystemName":    "Collections",
		"WallpaperType": "Media",
	})

	// Export system wallpapers
	for _, system := range systemPaths.Systems {
		if system.Tag == "" {
			continue // Skip systems without tags
		}

		systemBg := filepath.Join(system.MediaPath, "bg.png")

		// Create filename with system tag
		var filename string
		if strings.Contains(system.Name, fmt.Sprintf("(%s)", system.Tag)) {
			filename = fmt.Sprintf("%s.png", system.Name)
		} else {
			filename = fmt.Sprintf("%s (%s).png", system.Name, system.Tag)
		}

		exportWallpaper(systemBg, "SystemWallpapers", filename, map[string]string{
			"SystemName":    system.Name,
			"SystemTag":     system.Tag,
			"WallpaperType": "System",
		})
	}

	// Export collection wallpapers
	collectionsDir := filepath.Join(systemPaths.Root, "Collections")
	entries, err := os.ReadDir(collectionsDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			collectionName := entry.Name()
			collectionBg := filepath.Join(collectionsDir, collectionName, ".media", "bg.png")

			filename := fmt.Sprintf("%s.png", collectionName)

			exportWallpaper(collectionBg, "CollectionWallpapers", filename, map[string]string{
				"CollectionName": collectionName,
				"WallpaperType":  "Collection",
			})
		}
	}

	// Create preview image (use Recently Played bg or a default)
	previewPath := filepath.Join(exportPath, "preview.png")
	if _, err := os.Stat(rpBg); err == nil {
		// Use Recently Played bg as preview
		if err := CopyFile(rpBg, previewPath); err != nil {
			logger.DebugFn("Warning: Could not copy preview image: %v", err)
			// Create default preview as fallback
			if err := CreateDefaultPreviewImage(previewPath, ComponentWallpaper); err != nil {
				logger.DebugFn("Warning: Could not create default preview: %v", err)
			}
		}
	} else {
		// Create default preview
		if err := CreateDefaultPreviewImage(previewPath, ComponentWallpaper); err != nil {
			logger.DebugFn("Warning: Could not create default preview: %v", err)
		}
	}

	// Write manifest
	if err := WriteComponentManifest(exportPath, wallpaperManifest); err != nil {
		return fmt.Errorf("error writing wallpaper manifest: %w", err)
	}

	logger.DebugFn("Wallpaper export completed: %s", name)

	// Show success message
	ui.ShowMessage(fmt.Sprintf("Wallpapers exported to '%s'", name), "3")

	return nil
}

// ExportIcons exports current icons as a .icon component package
func ExportIcons(name string) error {
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting icon export: %s", name)

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .icon extension
	if !strings.HasSuffix(name, ComponentExtension[ComponentIcon]) {
		name = name + ComponentExtension[ComponentIcon]
	}

	exportPath := filepath.Join(cwd, "Exports", name)

	// Create directories
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
	manifest, err := CreateComponentManifest(ComponentIcon, name)
	if err != nil {
		return fmt.Errorf("error creating icon manifest: %w", err)
	}

	iconManifest := manifest.(*IconManifest)

	// Get system paths for copying icons
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Function to export an icon and add it to the manifest
	exportIcon := func(srcPath, relativeDstDir, filename string, metadata map[string]string) error {
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return nil // Skip if source doesn't exist
		}

		dstDir := filepath.Join(exportPath, relativeDstDir)
		dstPath := filepath.Join(dstDir, filename)

		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.DebugFn("Warning: Could not copy icon %s: %v", srcPath, err)
			return err
		}

		// Add to manifest
		themePath := filepath.Join(relativeDstDir, filename)
		iconManifest.PathMappings = append(
			iconManifest.PathMappings,
			PathMapping{
				ThemePath:  themePath,
				SystemPath: srcPath,
				Metadata:   metadata,
			},
		)

		// Add to content list and update counts
		switch relativeDstDir {
		case "SystemIcons":
			iconManifest.Content.SystemIcons = append(iconManifest.Content.SystemIcons, filename)
			iconManifest.Content.SystemCount++
		case "ToolIcons":
			iconManifest.Content.ToolIcons = append(iconManifest.Content.ToolIcons, filename)
			iconManifest.Content.ToolCount++
		case "CollectionIcons":
			iconManifest.Content.CollectionIcons = append(iconManifest.Content.CollectionIcons, filename)
			iconManifest.Content.CollectionCount++
		}

		logger.DebugFn("Exported icon: %s", dstPath)
		return nil
	}

	// Export Recently Played icon
	rpIcon := filepath.Join(systemPaths.Root, ".media", "Recently Played.png")
	exportIcon(rpIcon, "SystemIcons", "Recently Played.png", map[string]string{
		"SystemName": "Recently Played",
		"IconType":   "System",
	})

	// Export Tools icon
	toolsParentDir := filepath.Dir(systemPaths.Tools)
	toolsIcon := filepath.Join(toolsParentDir, ".media", "tg5040.png")
	exportIcon(toolsIcon, "SystemIcons", "Tools.png", map[string]string{
		"SystemName": "Tools",
		"IconType":   "System",
	})

	// Export Collections icon
	collectionsIcon := filepath.Join(systemPaths.Root, ".media", "Collections.png")
	exportIcon(collectionsIcon, "SystemIcons", "Collections.png", map[string]string{
		"SystemName": "Collections",
		"IconType":   "System",
	})

	// Export system icons
	systemIconsDir := filepath.Join(systemPaths.Roms, ".media")
	if _, err := os.Stat(systemIconsDir); err == nil {
		entries, err := os.ReadDir(systemIconsDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				// Skip special icons we already handled
				if entry.Name() == "Recently Played.png" ||
				   entry.Name() == "Collections.png" ||
				   entry.Name() == "tg5040.png" {
					continue
				}

				// Check for system tag pattern
				tagRegex := regexp.MustCompile(`\((.*?)\)`)
				if !tagRegex.MatchString(entry.Name()) {
					continue
				}

				systemIconPath := filepath.Join(systemIconsDir, entry.Name())

				// Extract system tag for metadata
				matches := tagRegex.FindStringSubmatch(entry.Name())
				systemTag := ""
				if len(matches) >= 2 {
					systemTag = matches[1]
				}

				exportIcon(systemIconPath, "SystemIcons", entry.Name(), map[string]string{
					"SystemName": strings.TrimSuffix(entry.Name(), ".png"),
					"SystemTag":  systemTag,
					"IconType":   "System",
				})
			}
		}
	}

	// Export tool icons
	toolsDir := filepath.Join(systemPaths.Tools)
	toolEntries, err := os.ReadDir(toolsDir)
	if err == nil {
		for _, entry := range toolEntries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			toolName := entry.Name()
			toolIcon := filepath.Join(toolsDir, toolName, ".media", toolName+".png")

			exportIcon(toolIcon, "ToolIcons", fmt.Sprintf("%s.png", toolName), map[string]string{
				"ToolName": toolName,
				"IconType": "Tool",
			})
		}
	}

	// Export collection icons
	collectionsDir := filepath.Join(systemPaths.Root, "Collections")
	colEntries, err := os.ReadDir(collectionsDir)
	if err == nil {
		for _, entry := range colEntries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			collectionName := entry.Name()
			collectionIcon := filepath.Join(collectionsDir, collectionName, ".media", collectionName+".png")

			exportIcon(collectionIcon, "CollectionIcons", fmt.Sprintf("%s.png", collectionName), map[string]string{
				"CollectionName": collectionName,
				"IconType":       "Collection",
			})
		}
	}

	// Create preview image (use a system icon or default)
	previewPath := filepath.Join(exportPath, "preview.png")
	if _, err := os.Stat(collectionsIcon); err == nil {
		// Use Collections icon as preview
		if err := CopyFile(collectionsIcon, previewPath); err != nil {
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

	// Write manifest
	if err := WriteComponentManifest(exportPath, iconManifest); err != nil {
		return fmt.Errorf("error writing icon manifest: %w", err)
	}

	logger.DebugFn("Icon export completed: %s", name)

	// Show success message
	ui.ShowMessage(fmt.Sprintf("Icons exported to '%s'", name), "3")

	return nil
}

// ExportAccents exports current accent settings as a .acc component package
func ExportAccents(name string) error {
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting accent export: %s", name)

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .acc extension
	if !strings.HasSuffix(name, ComponentExtension[ComponentAccent]) {
		name = name + ComponentExtension[ComponentAccent]
	}

	exportPath := filepath.Join(cwd, "Exports", name)

	// Create export directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", exportPath, err)
	}

	// Create component manifest
	manifest, err := CreateComponentManifest(ComponentAccent, name)
	if err != nil {
		return fmt.Errorf("error creating accent manifest: %w", err)
	}

	accentManifest := manifest.(*AccentManifest)

	// Read current accent settings
	settingsPath := "/mnt/SDCARD/.userdata/shared/minuisettings.txt"
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
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

		// Store color values in manifest
		switch key {
		case "color1":
			accentManifest.AccentColors.Color1 = value
		case "color2":
			accentManifest.AccentColors.Color2 = value
		case "color3":
			accentManifest.AccentColors.Color3 = value
		case "color4":
			accentManifest.AccentColors.Color4 = value
		case "color5":
			accentManifest.AccentColors.Color5 = value
		case "color6":
			accentManifest.AccentColors.Color6 = value
		}
	}

	// Create preview image (default for now)
	previewPath := filepath.Join(exportPath, "preview.png")
	if err := CreateDefaultPreviewImage(previewPath, ComponentAccent); err != nil {
		logger.DebugFn("Warning: Could not create default preview: %v", err)
	}

	// Write manifest
	if err := WriteComponentManifest(exportPath, accentManifest); err != nil {
		return fmt.Errorf("error writing accent manifest: %w", err)
	}

	logger.DebugFn("Accent export completed: %s", name)

	// Show success message
	ui.ShowMessage(fmt.Sprintf("Accent colors exported to '%s'", name), "3")

	return nil
}

// ExportLEDs exports current LED settings as a .led component package
func ExportLEDs(name string) error {
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting LED export: %s", name)

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .led extension
	if !strings.HasSuffix(name, ComponentExtension[ComponentLED]) {
		name = name + ComponentExtension[ComponentLED]
	}

	exportPath := filepath.Join(cwd, "Exports", name)

	// Create export directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", exportPath, err)
	}

	// Create component manifest
	manifest, err := CreateComponentManifest(ComponentLED, name)
	if err != nil {
		return fmt.Errorf("error creating LED manifest: %w", err)
	}

	ledManifest := manifest.(*LEDManifest)

	// Read current LED settings
	settingsPath := "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
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
				currentLED = &ledManifest.LEDSettings.F1Key
			case "F2 key":
				currentLED = &ledManifest.LEDSettings.F2Key
			case "Top bar":
				currentLED = &ledManifest.LEDSettings.TopBar
			case "L&R triggers":
				currentLED = &ledManifest.LEDSettings.LRTriggers
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

	// Note: LEDs don't have a preview image by design

	// Write manifest
	if err := WriteComponentManifest(exportPath, ledManifest); err != nil {
		return fmt.Errorf("error writing LED manifest: %w", err)
	}

	logger.DebugFn("LED export completed: %s", name)

	// Show success message
	ui.ShowMessage(fmt.Sprintf("LED settings exported to '%s'", name), "3")

	return nil
}

// ExportFonts exports current fonts as a .font component package
func ExportFonts(name string) error {
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting font export: %s", name)

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .font extension
	if !strings.HasSuffix(name, ComponentExtension[ComponentFont]) {
		name = name + ComponentExtension[ComponentFont]
	}

	exportPath := filepath.Join(cwd, "Exports", name)

	// Create export directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", exportPath, err)
	}

	// Create component manifest
	manifest, err := CreateComponentManifest(ComponentFont, name)
	if err != nil {
		return fmt.Errorf("error creating font manifest: %w", err)
	}

	fontManifest := manifest.(*FontManifest)

	// Define font paths
	fontPaths := map[string]string{
		"OG":         "/mnt/SDCARD/.userdata/shared/font2.ttf",
		"OG.backup":  "/mnt/SDCARD/.userdata/shared/font2.ttf.bak",
		"Next":       "/mnt/SDCARD/.userdata/shared/font1.ttf",
		"Next.backup": "/mnt/SDCARD/.userdata/shared/font1.ttf.bak",
	}

	// Export each font and update manifest
	for fontName, sourcePath := range fontPaths {
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			logger.DebugFn("Font file not found: %s", sourcePath)
			continue
		}

		dstPath := filepath.Join(exportPath, fontName+".ttf")

		if err := CopyFile(sourcePath, dstPath); err != nil {
			logger.DebugFn("Warning: Could not copy font %s: %v", fontName, err)
			continue
		}

		// Add to manifest
		fontManifest.PathMappings[fontName] = PathMapping{
			ThemePath:  fontName + ".ttf",
			SystemPath: sourcePath,
		}

		// Update content flags
		if fontName == "OG" {
			fontManifest.Content.OGReplaced = true
		} else if fontName == "Next" {
			fontManifest.Content.NextReplaced = true
		}

		logger.DebugFn("Exported font: %s", dstPath)
	}

	// Create preview image (default for now)
	previewPath := filepath.Join(exportPath, "preview.png")
	if err := CreateDefaultPreviewImage(previewPath, ComponentFont); err != nil {
		logger.DebugFn("Warning: Could not create default preview: %v", err)
	}

	// Write manifest
	if err := WriteComponentManifest(exportPath, fontManifest); err != nil {
		return fmt.Errorf("error writing font manifest: %w", err)
	}

	logger.DebugFn("Font export completed: %s", name)

	// Show success message
	ui.ShowMessage(fmt.Sprintf("Fonts exported to '%s'", name), "3")

	return nil
}

// ExportOverlays exports current overlays as a .over component package
func ExportOverlays(name string) error {
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting overlay export: %s", name)

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create export directory path with .over extension
	if !strings.HasSuffix(name, ComponentExtension[ComponentOverlay]) {
		name = name + ComponentExtension[ComponentOverlay]
	}

	exportPath := filepath.Join(cwd, "Exports", name)

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
	manifest, err := CreateComponentManifest(ComponentOverlay, name)
	if err != nil {
		return fmt.Errorf("error creating overlay manifest: %w", err)
	}

	overlayManifest := manifest.(*OverlayManifest)

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Check for overlays directory
	overlaysDir := filepath.Join(systemPaths.Root, "Overlays")
	if _, err := os.Stat(overlaysDir); os.IsNotExist(err) {
		logger.DebugFn("Overlays directory not found: %s", overlaysDir)
		return fmt.Errorf("overlays directory not found: %s", overlaysDir)
	}

	// List system directories in Overlays
	entries, err := os.ReadDir(overlaysDir)
	if err != nil {
		logger.DebugFn("Error reading Overlays directory: %v", err)
		return fmt.Errorf("error reading overlays directory: %w", err)
	}

	// Process each system's overlays
	hasOverlays := false
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		systemTag := entry.Name()
		systemOverlaysPath := filepath.Join(overlaysDir, systemTag)
		exportSystemDir := filepath.Join(systemsDir, systemTag)

		// Create system directory in export
		if err := os.MkdirAll(exportSystemDir, 0755); err != nil {
			logger.DebugFn("Error creating system overlay directory: %v", err)
			continue
		}

		// List overlay files for this system
		overlayFiles, err := os.ReadDir(systemOverlaysPath)
		if err != nil {
			logger.DebugFn("Error reading system overlays directory %s: %v", systemTag, err)
			continue
		}

		var systemHasOverlays bool

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
			dstPath := filepath.Join(exportSystemDir, file.Name())

			// Copy the overlay file
			if err := CopyFile(srcPath, dstPath); err != nil {
				logger.DebugFn("Warning: Could not copy overlay %s: %v", file.Name(), err)
				continue
			}

			// Add to manifest
			themePath := filepath.Join("Systems", systemTag, file.Name())
			overlayManifest.PathMappings = append(
				overlayManifest.PathMappings,
				PathMapping{
					ThemePath:  themePath,
					SystemPath: srcPath,
					Metadata: map[string]string{
						"SystemTag":   systemTag,
						"OverlayName": file.Name(),
					},
				},
			)

			systemHasOverlays = true
			hasOverlays = true
			logger.DebugFn("Exported overlay %s for system %s", file.Name(), systemTag)
		}

		// If this system had overlays, add it to the systems list
		if systemHasOverlays {
			// Check if system is already in the list
			var systemExists bool
			for _, sys := range overlayManifest.Content.Systems {
				if sys == systemTag {
					systemExists = true
					break
				}
			}

			if !systemExists {
				overlayManifest.Content.Systems = append(overlayManifest.Content.Systems, systemTag)
			}
		}
	}

	if !hasOverlays {
		return fmt.Errorf("no overlays found to export")
	}

	// Create preview image (default for now)
	previewPath := filepath.Join(exportPath, "preview.png")
	if err := CreateDefaultPreviewImage(previewPath, ComponentOverlay); err != nil {
		logger.DebugFn("Warning: Could not create default preview: %v", err)
	}

	// Write manifest
	if err := WriteComponentManifest(exportPath, overlayManifest); err != nil {
		return fmt.Errorf("error writing overlay manifest: %w", err)
	}

	logger.DebugFn("Overlay export completed: %s", name)

	// Show success message
	ui.ShowMessage(fmt.Sprintf("Overlays exported to '%s'", name), "3")

	return nil
}

// Helper function to ensure component directories exist for importing
func EnsureComponentDirectories() error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create main Components directory
	componentsDir := filepath.Join(cwd, "Components")
	if err := os.MkdirAll(componentsDir, 0755); err != nil {
		return fmt.Errorf("error creating Components directory: %w", componentsDir, err)
	}

	// Component subdirectories to create
	directories := []string{
		filepath.Join(componentsDir, "Wallpapers"),
		filepath.Join(componentsDir, "Icons"),
		filepath.Join(componentsDir, "Accents"),
		filepath.Join(componentsDir, "Overlays"),
		filepath.Join(componentsDir, "LEDs"),
		filepath.Join(componentsDir, "Fonts"),
	}

	// Create each directory
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	return nil
}