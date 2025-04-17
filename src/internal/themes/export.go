// src/internal/themes/export.go
// Simplified implementation of theme export functionality

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"nextui-themes/internal/ui"
	"strings"
)

// CreateThemeExportDirectory creates a new theme directory with sequential naming
func CreateThemeExportDirectory() (string, error) {
    // Get the current directory
    cwd, err := os.Getwd()
    if err != nil {
        return "", fmt.Errorf("error getting current directory: %w", err)
    }

    // Path to Themes/Exports directory
    exportsDir := filepath.Join(cwd, "Themes", "Exports")

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
        "Settings",
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

	// Initialize manifest
	manifest := &ThemeManifest{}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logger.DebugFn("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Scan and export various components
	// This is a simplified version that just focuses on wallpapers as an example
	exportWallpapers(themePath, manifest, systemPaths, logger)

	// Write manifest
	if err := WriteManifest(themePath, manifest, logger); err != nil {
		logger.DebugFn("Error writing manifest: %v", err)
		return fmt.Errorf("error writing manifest: %w", err)
	}

	logger.DebugFn("Theme export completed successfully: %s", themePath)

	// Show success message to user
	themeName := filepath.Base(themePath)
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
                        "SystemName": "Root",
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
                        "SystemName": "Root",
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
                        "SystemName": "Recently Played",
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
                        "SystemName": "Tools",
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
                        "SystemName": "Collections",
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
                            "SystemName": system.Name,
                            "SystemTag": system.Tag,
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
                            "WallpaperType": "Collection",
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