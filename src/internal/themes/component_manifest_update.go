// src/internal/themes/component_manifest_update.go
// Provides functions for automatically updating component manifests based on content

package themes

import (
	"fmt"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// UpdateComponentManifest updates a component's manifest based on its actual content
func UpdateComponentManifest(componentPath string) error {
	// Determine component type from file extension
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

	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Updating manifest for component: %s (type: %s)", componentPath, componentType)

	// Create a system paths instance for reference
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logger.DebugFn("Warning: Error getting system paths: %v", err)
		// Continue anyway, as we can still update most of the manifest
	}

	// First try to load the existing manifest to preserve author information
	var existingAuthor string
	manifestObj, err := LoadComponentManifest(componentPath)
	if err == nil {
		// Extract author from existing manifest based on component type
		switch componentType {
		case ComponentWallpaper:
			if m, ok := manifestObj.(*WallpaperManifest); ok && m.ComponentInfo.Author != "" {
				existingAuthor = m.ComponentInfo.Author
			}
		case ComponentIcon:
			if m, ok := manifestObj.(*IconManifest); ok && m.ComponentInfo.Author != "" {
				existingAuthor = m.ComponentInfo.Author
			}
		case ComponentOverlay:
			if m, ok := manifestObj.(*OverlayManifest); ok && m.ComponentInfo.Author != "" {
				existingAuthor = m.ComponentInfo.Author
			}
		case ComponentFont:
			if m, ok := manifestObj.(*FontManifest); ok && m.ComponentInfo.Author != "" {
				existingAuthor = m.ComponentInfo.Author
			}
		case ComponentAccent:
			if m, ok := manifestObj.(*AccentManifest); ok && m.ComponentInfo.Author != "" {
				existingAuthor = m.ComponentInfo.Author
			}
		case ComponentLED:
			if m, ok := manifestObj.(*LEDManifest); ok && m.ComponentInfo.Author != "" {
				existingAuthor = m.ComponentInfo.Author
			}
		}
	}

	// Dispatch to specific update function based on component type
	var updateErr error
	switch componentType {
	case ComponentWallpaper:
		updateErr = UpdateWallpaperManifest(componentPath, systemPaths, logger)
	case ComponentIcon:
		updateErr = UpdateIconManifest(componentPath, systemPaths, logger)
	case ComponentOverlay:
		updateErr = UpdateOverlayManifest(componentPath, systemPaths, logger)
	case ComponentFont:
		updateErr = UpdateFontManifest(componentPath, logger)
	case ComponentAccent:
		updateErr = UpdateAccentManifest(componentPath, logger)
	case ComponentLED:
		updateErr = UpdateLEDManifest(componentPath, logger)
	default:
		return fmt.Errorf("unhandled component type: %s", componentType)
	}

	// If we had an existing author, restore it after the update
	if existingAuthor != "" && updateErr == nil {
		// Load the updated manifest
		updatedManifest, err := LoadComponentManifest(componentPath)
		if err == nil {
			// Set the author back to the original value
			switch componentType {
			case ComponentWallpaper:
				if m, ok := updatedManifest.(*WallpaperManifest); ok {
					m.ComponentInfo.Author = existingAuthor
					// Write the manifest back
					WriteComponentManifest(componentPath, m)
				}
			case ComponentIcon:
				if m, ok := updatedManifest.(*IconManifest); ok {
					m.ComponentInfo.Author = existingAuthor
					WriteComponentManifest(componentPath, m)
				}
			case ComponentOverlay:
				if m, ok := updatedManifest.(*OverlayManifest); ok {
					m.ComponentInfo.Author = existingAuthor
					WriteComponentManifest(componentPath, m)
				}
			case ComponentFont:
				if m, ok := updatedManifest.(*FontManifest); ok {
					m.ComponentInfo.Author = existingAuthor
					WriteComponentManifest(componentPath, m)
				}
			case ComponentAccent:
				if m, ok := updatedManifest.(*AccentManifest); ok {
					m.ComponentInfo.Author = existingAuthor
					WriteComponentManifest(componentPath, m)
				}
			case ComponentLED:
				if m, ok := updatedManifest.(*LEDManifest); ok {
					m.ComponentInfo.Author = existingAuthor
					WriteComponentManifest(componentPath, m)
				}
			}
		}
	}

	return updateErr
}

func UpdateWallpaperManifest(componentPath string, systemPaths *system.SystemPaths, logger *Logger) error {
	logger.DebugFn("Updating wallpaper manifest for: %s", componentPath)

	// Get the component name from the path
	componentName := filepath.Base(componentPath)

	// Load existing manifest to preserve component_info
	manifestObj, err := LoadComponentManifest(componentPath)
	if err != nil {
		// If manifest doesn't exist, create a new one
		manifestObj, err = CreateComponentManifest(ComponentWallpaper, componentName)
		if err != nil {
			return fmt.Errorf("error creating wallpaper manifest: %w", err)
		}
	}

	wallpaperManifest, ok := manifestObj.(*WallpaperManifest)
	if !ok {
		return fmt.Errorf("invalid manifest type for wallpaper component")
	}

	// Always update component name to match the directory name
	wallpaperManifest.ComponentInfo.Name = componentName

	// Clear existing content data (but preserve other component_info fields)
	wallpaperManifest.Content.Count = 0
	wallpaperManifest.Content.SystemWallpapers = []string{}
	wallpaperManifest.Content.ListWallpapers = []string{}    // Clear the list wallpapers array
	wallpaperManifest.Content.CollectionWallpapers = []string{}
	wallpaperManifest.PathMappings = []PathMapping{}

	// Check for wallpapers in SystemWallpapers directory
	systemWallpapersDir := filepath.Join(componentPath, "SystemWallpapers")
	if _, err := os.Stat(systemWallpapersDir); err == nil {
		entries, err := os.ReadDir(systemWallpapersDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				fileName := entry.Name()
				filePath := filepath.Join("SystemWallpapers", fileName)

				// Add to content list - only add to SystemWallpapers if it's not a list wallpaper
				if !strings.Contains(fileName, "-list") {
					wallpaperManifest.Content.SystemWallpapers = append(
						wallpaperManifest.Content.SystemWallpapers,
						fileName,
					)
				}

				// Determine system path and metadata
				var systemPath string
				var metadata map[string]string

				// Special case handling for predefined names
				switch strings.TrimSuffix(fileName, ".png") {
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
					re := regexp.MustCompile(`\((.*?)\)`)
					matches := re.FindStringSubmatch(fileName)

					if len(matches) >= 2 {
						systemTag := matches[1]

						// Find matching system by tag
						var systemFound bool
						for _, system := range systemPaths.Systems {
							if system.Tag == systemTag {
								systemPath = filepath.Join(system.MediaPath, "bg.png")
								metadata = map[string]string{
									"SystemName":    system.Name,
									"SystemTag":     systemTag,
									"WallpaperType": "System",
								}
								systemFound = true
								break
							}
						}

						// If system not found, create a default path
						if !systemFound && systemTag != "" {
							baseName := strings.TrimSuffix(fileName, ".png")
							systemName := strings.TrimSuffix(strings.Split(baseName, "(")[0], " ")

							systemPath = filepath.Join(systemPaths.Roms, fmt.Sprintf("%s (%s)", systemName, systemTag), ".media", "bg.png")
							metadata = map[string]string{
								"SystemName":    systemName,
								"SystemTag":     systemTag,
								"WallpaperType": "System",
							}
						}
					}
				}

				// If we determined a system path, add to path mappings
				if systemPath != "" {
					wallpaperManifest.PathMappings = append(
						wallpaperManifest.PathMappings,
						PathMapping{
							ThemePath:  filePath,
							SystemPath: systemPath,
							Metadata:   metadata,
						},
					)
					wallpaperManifest.Content.Count++
					logger.DebugFn("Added system wallpaper to manifest: %s", fileName)
				} else {
					logger.DebugFn("Could not determine system path for wallpaper: %s", fileName)
				}
			}
		}
	}

	// NEW: Check for wallpapers in ListWallpapers directory
	listWallpapersDir := filepath.Join(componentPath, "ListWallpapers")
	if _, err := os.Stat(listWallpapersDir); err == nil {
		entries, err := os.ReadDir(listWallpapersDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				fileName := entry.Name()
				filePath := filepath.Join("ListWallpapers", fileName)

				// Add to ListWallpapers content list
				wallpaperManifest.Content.ListWallpapers = append(
					wallpaperManifest.Content.ListWallpapers,
					fileName,
				)

				// Check if this is a list wallpaper (ends with -list.png)
				baseName := strings.TrimSuffix(fileName, ".png")
				if !strings.HasSuffix(baseName, "-list") {
					logger.DebugFn("List wallpaper doesn't have -list suffix: %s", fileName)
				}

				// Extract system tag
				baseNameWithoutSuffix := strings.TrimSuffix(baseName, "-list")
				re := regexp.MustCompile(`\((.*?)\)`)
				matches := re.FindStringSubmatch(baseNameWithoutSuffix)

				if len(matches) >= 2 {
					systemTag := matches[1]

					// Find matching system by tag
					var systemFound bool
					for _, system := range systemPaths.Systems {
						if system.Tag == systemTag {
							systemPath := filepath.Join(system.MediaPath, "bglist.png")
							metadata := map[string]string{
								"SystemName":    system.Name,
								"SystemTag":     systemTag,
								"WallpaperType": "List",
							}

							wallpaperManifest.PathMappings = append(
								wallpaperManifest.PathMappings,
								PathMapping{
									ThemePath:  filePath,
									SystemPath: systemPath,
									Metadata:   metadata,
								},
							)
							wallpaperManifest.Content.Count++
							logger.DebugFn("Added list wallpaper to manifest: %s", fileName)
							systemFound = true
							break
						}
					}

					// If system not found, create a default path
					if !systemFound && systemTag != "" {
						systemName := strings.TrimSuffix(strings.Split(baseNameWithoutSuffix, "(")[0], " ")
						systemPath := filepath.Join(systemPaths.Roms, fmt.Sprintf("%s (%s)", systemName, systemTag), ".media", "bglist.png")
						metadata := map[string]string{
							"SystemName":    systemName,
							"SystemTag":     systemTag,
							"WallpaperType": "List",
						}

						wallpaperManifest.PathMappings = append(
							wallpaperManifest.PathMappings,
							PathMapping{
								ThemePath:  filePath,
								SystemPath: systemPath,
								Metadata:   metadata,
							},
						)
						wallpaperManifest.Content.Count++
						logger.DebugFn("Added default list wallpaper to manifest: %s", fileName)
					}
				} else {
					logger.DebugFn("Could not determine system for list wallpaper: %s", fileName)
				}
			}
		}
	}

	// Rest of function (for collection wallpapers) remains unchanged...
	// Check for wallpapers in CollectionWallpapers directory
	collectionWallpapersDir := filepath.Join(componentPath, "CollectionWallpapers")
	if _, err := os.Stat(collectionWallpapersDir); err == nil {
		entries, err := os.ReadDir(collectionWallpapersDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				fileName := entry.Name()
				filePath := filepath.Join("CollectionWallpapers", fileName)

				// Add to content list
				wallpaperManifest.Content.CollectionWallpapers = append(
					wallpaperManifest.Content.CollectionWallpapers,
					fileName,
				)

				// Determine collection name and system path
				collectionName := strings.TrimSuffix(fileName, ".png")
				systemPath := filepath.Join(systemPaths.Root, "Collections", collectionName, ".media", "bg.png")

				metadata := map[string]string{
					"CollectionName": collectionName,
					"WallpaperType":  "Collection",
				}

				// Add to path mappings
				wallpaperManifest.PathMappings = append(
					wallpaperManifest.PathMappings,
					PathMapping{
						ThemePath:  filePath,
						SystemPath: systemPath,
						Metadata:   metadata,
					},
				)
				wallpaperManifest.Content.Count++
				logger.DebugFn("Added collection wallpaper to manifest: %s", fileName)
			}
		}
	}

	// Write updated manifest
	return WriteComponentManifest(componentPath, wallpaperManifest)
}

// UpdateIconManifest updates an icon component's manifest based on its content
func UpdateIconManifest(componentPath string, systemPaths *system.SystemPaths, logger *Logger) error {
	logger.DebugFn("Updating icon manifest for: %s", componentPath)

	// Get the component name from the path
	componentName := filepath.Base(componentPath)

	// Load existing manifest to preserve component_info
	manifestObj, err := LoadComponentManifest(componentPath)
	if err != nil {
		// If manifest doesn't exist, create a new one
		manifestObj, err = CreateComponentManifest(ComponentIcon, componentName)
		if err != nil {
			return fmt.Errorf("error creating icon manifest: %w", err)
		}
	}

	iconManifest, ok := manifestObj.(*IconManifest)
	if !ok {
		return fmt.Errorf("invalid manifest type for icon component")
	}

	// Always update component name to match the directory name
	iconManifest.ComponentInfo.Name = componentName

	// Rest of the function remains the same...

	// Clear existing content data (but preserve component_info)
	iconManifest.Content.SystemCount = 0
	iconManifest.Content.ToolCount = 0
	iconManifest.Content.CollectionCount = 0
	iconManifest.Content.SystemIcons = []string{}
	iconManifest.Content.ToolIcons = []string{}
	iconManifest.Content.CollectionIcons = []string{}
	iconManifest.PathMappings = []PathMapping{}

	// Check for icons in SystemIcons directory
	systemIconsDir := filepath.Join(componentPath, "SystemIcons")
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

				fileName := entry.Name()
				filePath := filepath.Join("SystemIcons", fileName)

				// Add to content list
				iconManifest.Content.SystemIcons = append(
					iconManifest.Content.SystemIcons,
					fileName,
				)

				// Determine system path and metadata
				var systemPath string
				var metadata map[string]string

				// Special case handling for predefined names
				switch fileName {
				case "Recently Played.png":
					systemPath = filepath.Join(systemPaths.Root, ".media", "Recently Played.png")
					metadata = map[string]string{
						"SystemName": "Recently Played",
						"IconType":   "System",
					}

				case "Tools.png":
					// Get parent directory of Tools since Tools path includes tg5040
					toolsParentDir := filepath.Dir(systemPaths.Tools)
					systemPath = filepath.Join(toolsParentDir, ".media", "tg5040.png")
					metadata = map[string]string{
						"SystemName": "Tools",
						"IconType":   "System",
					}

				case "Collections.png":
					systemPath = filepath.Join(systemPaths.Root, ".media", "Collections.png")
					metadata = map[string]string{
						"SystemName": "Collections",
						"IconType":   "System",
					}

				default:
					// Check for system tag in filename
					re := regexp.MustCompile(`\((.*?)\)`)
					matches := re.FindStringSubmatch(fileName)

					if len(matches) >= 2 {
						systemTag := matches[1]

						// First, try to find the exact ROM directory by tag
						// This ensures that we use the actual directory name
						var exactSystemName string
						var exactMatch bool

						for _, system := range systemPaths.Systems {
							if system.Tag == systemTag {
								// Use actual ROM directory name instead of icon file name
								exactSystemName = system.Name
								exactMatch = true
								break
							}
						}

						if exactMatch {
							// Use exact ROM directory name for the icon
							// This ensures the icon will be displayed correctly
							systemPath = filepath.Join(systemPaths.Roms, ".media", exactSystemName+".png")
							metadata = map[string]string{
								"SystemName":     exactSystemName,
								"SystemTag":      systemTag,
								"IconType":       "System",
								"RenameRequired": "true",   // Flag that this icon needs renaming
								"OriginalName":   fileName, // Store original name for identification
							}
						} else {
							// System not found in available systems, use as-is
							systemPath = filepath.Join(systemPaths.Roms, ".media", fileName)
							metadata = map[string]string{
								"SystemName": strings.TrimSuffix(fileName, ".png"),
								"SystemTag":  systemTag,
								"IconType":   "System",
							}
						}
					}
				}

				// If we determined a system path, add to path mappings
				if systemPath != "" {
					iconManifest.PathMappings = append(
						iconManifest.PathMappings,
						PathMapping{
							ThemePath:  filePath,
							SystemPath: systemPath,
							Metadata:   metadata,
						},
					)
					iconManifest.Content.SystemCount++
					logger.DebugFn("Added system icon to manifest: %s", fileName)
				} else {
					logger.DebugFn("Could not determine system path for icon: %s", fileName)
				}
			}
		}
	}

	// Check for icons in ToolIcons directory
	toolIconsDir := filepath.Join(componentPath, "ToolIcons")
	if _, err := os.Stat(toolIconsDir); err == nil {
		entries, err := os.ReadDir(toolIconsDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				fileName := entry.Name()
				filePath := filepath.Join("ToolIcons", fileName)

				// Add to content list
				iconManifest.Content.ToolIcons = append(
					iconManifest.Content.ToolIcons,
					fileName,
				)

				// Determine tool name and system path
				toolName := strings.TrimSuffix(fileName, ".png")
				systemPath := filepath.Join(systemPaths.Tools, toolName, ".media", toolName+".png")

				metadata := map[string]string{
					"ToolName": toolName,
					"IconType": "Tool",
				}

				// Add to path mappings
				iconManifest.PathMappings = append(
					iconManifest.PathMappings,
					PathMapping{
						ThemePath:  filePath,
						SystemPath: systemPath,
						Metadata:   metadata,
					},
				)
				iconManifest.Content.ToolCount++
				logger.DebugFn("Added tool icon to manifest: %s", fileName)
			}
		}
	}

	// Check for icons in CollectionIcons directory
	collectionIconsDir := filepath.Join(componentPath, "CollectionIcons")
	if _, err := os.Stat(collectionIconsDir); err == nil {
		entries, err := os.ReadDir(collectionIconsDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				// Check if file has a PNG extension
				if !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					continue
				}

				fileName := entry.Name()
				filePath := filepath.Join("CollectionIcons", fileName)

				// Add to content list
				iconManifest.Content.CollectionIcons = append(
					iconManifest.Content.CollectionIcons,
					fileName,
				)

				// Determine collection name and system path
				collectionName := strings.TrimSuffix(fileName, ".png")
				systemPath := filepath.Join(systemPaths.Root, "Collections", collectionName, ".media", collectionName+".png")

				metadata := map[string]string{
					"CollectionName": collectionName,
					"IconType":       "Collection",
				}

				// Add to path mappings
				iconManifest.PathMappings = append(
					iconManifest.PathMappings,
					PathMapping{
						ThemePath:  filePath,
						SystemPath: systemPath,
						Metadata:   metadata,
					},
				)
				iconManifest.Content.CollectionCount++
				logger.DebugFn("Added collection icon to manifest: %s", fileName)
			}
		}
	}

	// Write updated manifest
	return WriteComponentManifest(componentPath, iconManifest)
}

// UpdateOverlayManifest updates an overlay component's manifest based on its content
func UpdateOverlayManifest(componentPath string, systemPaths *system.SystemPaths, logger *Logger) error {
	logger.DebugFn("Updating overlay manifest for: %s", componentPath)

	// Get the component name from the path
	componentName := filepath.Base(componentPath)

	// Load existing manifest to preserve component_info
	manifestObj, err := LoadComponentManifest(componentPath)
	if err != nil {
		// If manifest doesn't exist, create a new one
		manifestObj, err = CreateComponentManifest(ComponentOverlay, componentName)
		if err != nil {
			return fmt.Errorf("error creating overlay manifest: %w", err)
		}
	}

	overlayManifest, ok := manifestObj.(*OverlayManifest)
	if !ok {
		return fmt.Errorf("invalid manifest type for overlay component")
	}

	// Always update component name to match the directory name
	overlayManifest.ComponentInfo.Name = componentName

	// Clear existing content data (but preserve component_info)
	overlayManifest.Content.Systems = []string{}
	overlayManifest.PathMappings = []PathMapping{}

	// Check for overlays in Systems directory
	systemsDir := filepath.Join(componentPath, "Systems")
	if _, err := os.Stat(systemsDir); err == nil {
		entries, err := os.ReadDir(systemsDir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
					continue
				}

				systemTag := entry.Name()
				systemOverlaysPath := filepath.Join(systemsDir, systemTag)

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

					filePath := filepath.Join("Systems", systemTag, file.Name())

					// Determine system path
					systemPath := filepath.Join(systemPaths.Root, "Overlays", systemTag, file.Name())

					// Add to manifest
					overlayManifest.PathMappings = append(
						overlayManifest.PathMappings,
						PathMapping{
							ThemePath:  filePath,
							SystemPath: systemPath,
							Metadata: map[string]string{
								"SystemTag":   systemTag,
								"OverlayName": file.Name(),
							},
						},
					)

					hasOverlays = true
					logger.DebugFn("Added overlay to manifest: %s for system %s", file.Name(), systemTag)
				}

				// If this system had overlays, add it to the systems list
				if hasOverlays {
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
		}
	}

	// Write updated manifest
	return WriteComponentManifest(componentPath, overlayManifest)
}

// UpdateFontManifest updates a font component's manifest based on its content
func UpdateFontManifest(componentPath string, logger *Logger) error {
	logger.DebugFn("Updating font manifest for: %s", componentPath)

	// Get the component name from the path
	componentName := filepath.Base(componentPath)

	// Load existing manifest to preserve component_info
	manifestObj, err := LoadComponentManifest(componentPath)
	if err != nil {
		// If manifest doesn't exist, create a new one
		manifestObj, err = CreateComponentManifest(ComponentFont, componentName)
		if err != nil {
			return fmt.Errorf("error creating font manifest: %w", err)
		}
	}

	fontManifest, ok := manifestObj.(*FontManifest)
	if !ok {
		return fmt.Errorf("invalid manifest type for font component")
	}

	// Always update component name to match the directory name
	fontManifest.ComponentInfo.Name = componentName

	// Clear existing content data (but preserve component_info)
	fontManifest.Content.OGReplaced = false
	fontManifest.Content.NextReplaced = false
	fontManifest.PathMappings = make(map[string]PathMapping)

	// Define system paths for fonts - CORRECTED PATHS
	systemPaths := map[string]string{
		"OG":          "/mnt/SDCARD/.system/res/font2.ttf",
		"OG.backup":   "/mnt/SDCARD/.system/res/font2.backup.ttf", // Corrected extension
		"Next":        "/mnt/SDCARD/.system/res/font1.ttf",
		"Next.backup": "/mnt/SDCARD/.system/res/font1.backup.ttf", // Corrected extension
	}

	// Check for each font file
	fontFiles := []string{
		"OG.ttf",
		"Next.ttf",
		"OG.backup.ttf",
		"Next.backup.ttf",
	}

	for _, fontFile := range fontFiles {
		fontPath := filepath.Join(componentPath, fontFile)
		if _, err := os.Stat(fontPath); err == nil {
			// Font file exists
			fontName := strings.TrimSuffix(fontFile, ".ttf")

			// Add to manifest
			fontManifest.PathMappings[fontName] = PathMapping{
				ThemePath:  fontFile,
				SystemPath: systemPaths[fontName],
			}

			// Update content flags
			if fontName == "OG" {
				fontManifest.Content.OGReplaced = true
			} else if fontName == "Next" {
				fontManifest.Content.NextReplaced = true
			}

			logger.DebugFn("Added font to manifest: %s", fontName)
		}
	}

	// Write updated manifest
	return WriteComponentManifest(componentPath, fontManifest)
}

// UpdateAccentManifest validates and updates an accent component's manifest
func UpdateAccentManifest(componentPath string, logger *Logger) error {
	logger.DebugFn("Updating accent manifest for: %s", componentPath)

	// Get the component name from the path
	componentName := filepath.Base(componentPath)

	// Load existing manifest to preserve component_info
	manifestObj, err := LoadComponentManifest(componentPath)
	if err != nil {
		// If manifest doesn't exist, create a new one
		manifestObj, err = CreateComponentManifest(ComponentAccent, componentName)
		if err != nil {
			return fmt.Errorf("error creating accent manifest: %w", err)
		}
	}

	accentManifest, ok := manifestObj.(*AccentManifest)
	if !ok {
		return fmt.Errorf("invalid manifest type for accent component")
	}

	// Always update component name to match the directory name
	accentManifest.ComponentInfo.Name = componentName

	// For accent settings, we mainly validate the manifest as the data is stored in it
	// Ensure required color fields are present
	if accentManifest.AccentColors.Color1 == "" {
		accentManifest.AccentColors.Color1 = "0xFFFFFF" // Default white
	}
	if accentManifest.AccentColors.Color2 == "" {
		accentManifest.AccentColors.Color2 = "0x9B2257" // Default accent
	}
	if accentManifest.AccentColors.Color3 == "" {
		accentManifest.AccentColors.Color3 = "0x1E2329" // Default secondary
	}
	if accentManifest.AccentColors.Color4 == "" {
		accentManifest.AccentColors.Color4 = "0xFFFFFF" // Default list text
	}
	if accentManifest.AccentColors.Color5 == "" {
		accentManifest.AccentColors.Color5 = "0x000000" // Default selected text
	}
	if accentManifest.AccentColors.Color6 == "" {
		accentManifest.AccentColors.Color6 = "0xFFFFFF" // Default hint text
	}

	// Write updated manifest
	return WriteComponentManifest(componentPath, accentManifest)
}

// UpdateLEDManifest validates and updates an LED component's manifest
func UpdateLEDManifest(componentPath string, logger *Logger) error {
	logger.DebugFn("Updating LED manifest for: %s", componentPath)

	// Get the component name from the path
	componentName := filepath.Base(componentPath)

	// Load existing manifest to preserve component_info
	manifestObj, err := LoadComponentManifest(componentPath)
	if err != nil {
		// If manifest doesn't exist, create a new one
		manifestObj, err = CreateComponentManifest(ComponentLED, componentName)
		if err != nil {
			return fmt.Errorf("error creating LED manifest: %w", err)
		}
	}

	ledManifest, ok := manifestObj.(*LEDManifest)
	if !ok {
		return fmt.Errorf("invalid manifest type for LED component")
	}

	// Always update component name to match the directory name
	ledManifest.ComponentInfo.Name = componentName

	// For LED settings, we mainly validate the manifest as the data is stored in it
	// Ensure all LED settings are initialized with defaults if missing

	// Default settings for an LED
	initLEDSetting := func(led *LEDSetting) {
		if led.Effect == 0 {
			led.Effect = 1 // Default effect
		}
		if led.Color1 == "" {
			led.Color1 = "0xFFFFFF" // Default white
		}
		if led.Color2 == "" {
			led.Color2 = "0x000000" // Default black
		}
		if led.Speed == 0 {
			led.Speed = 1000 // Default speed
		}
		if led.Brightness == 0 {
			led.Brightness = 100 // Default brightness
		}
		if led.Trigger == 0 {
			led.Trigger = 1 // Default trigger
		}
		if led.InBrightness == 0 {
			led.InBrightness = 100 // Default info brightness
		}
	}

	// Initialize all LED settings
	initLEDSetting(&ledManifest.LEDSettings.F1Key)
	initLEDSetting(&ledManifest.LEDSettings.F2Key)
	initLEDSetting(&ledManifest.LEDSettings.TopBar)
	initLEDSetting(&ledManifest.LEDSettings.LRTriggers)

	// Write updated manifest
	return WriteComponentManifest(componentPath, ledManifest)
}
