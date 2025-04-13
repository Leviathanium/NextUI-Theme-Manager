// src/internal/themes/import.go
// Implementation of theme import functionality

package themes

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"nextui-themes/internal/ui"
)

// ImportWallpapers imports wallpapers from a theme
func ImportWallpapers(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Skip if no wallpapers present
	if !manifest.Content.Wallpapers.Present {
		logger.Printf("No wallpapers found in theme, skipping wallpaper import")
		return nil
	}

	// Import each wallpaper using path mappings
	for _, mapping := range manifest.PathMappings.Wallpapers {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source wallpaper file not found: %s", srcPath)
			continue
		}

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.Printf("Warning: Could not create directory for wallpaper: %v", err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy wallpaper: %v", err)
		} else {
			logger.Printf("Imported wallpaper: %s -> %s", srcPath, dstPath)
		}
	}

	return nil
}

// ImportIcons imports icons from a theme
func ImportIcons(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Skip if no icons present
	if !manifest.Content.Icons.Present {
		logger.Printf("No icons found in theme, skipping icon import")
		return nil
	}

	// Import each icon using path mappings
	for _, mapping := range manifest.PathMappings.Icons {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source icon file not found: %s", srcPath)
			continue
		}

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.Printf("Warning: Could not create directory for icon: %v", err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy icon: %v", err)
		} else {
			logger.Printf("Imported icon: %s -> %s", srcPath, dstPath)
		}
	}

	return nil
}

// ImportOverlays imports overlays from a theme
func ImportOverlays(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Skip if no overlays present
	if !manifest.Content.Overlays.Present {
		logger.Printf("No overlays found in theme, skipping overlay import")
		return nil
	}

	// Import each overlay using path mappings
	for _, mapping := range manifest.PathMappings.Overlays {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source overlay file not found: %s", srcPath)
			continue
		}

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.Printf("Warning: Could not create directory for overlay: %v", err)
			continue
		}

		// Copy the file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy overlay: %v", err)
		} else {
			logger.Printf("Imported overlay: %s -> %s", srcPath, dstPath)
		}
	}

	return nil
}

// ImportFonts imports fonts from a theme
func ImportFonts(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Skip if no fonts present
	if !manifest.Content.Fonts.Present {
		logger.Printf("No fonts found in theme, skipping font import")
		return nil
	}

	// Process each font mapping
	for fontType, mapping := range manifest.PathMappings.Fonts {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			logger.Printf("Warning: Source font file not found: %s (%s)", srcPath, fontType)
			continue
		}

		// For backup fonts, only copy if they don't already exist
		if strings.Contains(fontType, "backup") {
			if _, err := os.Stat(dstPath); err == nil {
				logger.Printf("Backup font already exists, skipping: %s", dstPath)
				continue
			}
		}

		// For active fonts, create a backup if one doesn't exist
		if fontType == "og_font" {
			// Backup OG font if needed
			backupPath := filepath.Join(filepath.Dir(dstPath), "font2.backup.ttf")
			if _, err := os.Stat(backupPath); os.IsNotExist(err) {
				if _, err := os.Stat(dstPath); err == nil {
					if err := CopyFile(dstPath, backupPath); err != nil {
						logger.Printf("Warning: Could not create backup of OG font: %v", err)
					} else {
						logger.Printf("Created backup of current OG font: %s", backupPath)
					}
				}
			}
		} else if fontType == "next_font" {
			// Backup Next font if needed
			backupPath := filepath.Join(filepath.Dir(dstPath), "font1.backup.ttf")
			if _, err := os.Stat(backupPath); os.IsNotExist(err) {
				if _, err := os.Stat(dstPath); err == nil {
					if err := CopyFile(dstPath, backupPath); err != nil {
						logger.Printf("Warning: Could not create backup of Next font: %v", err)
					} else {
						logger.Printf("Created backup of current Next font: %s", backupPath)
					}
				}
			}
		}

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			logger.Printf("Warning: Could not create directory for font: %v", err)
			continue
		}

		// Copy the font file
		if err := CopyFile(srcPath, dstPath); err != nil {
			logger.Printf("Warning: Could not copy font file: %v", err)
		} else {
			logger.Printf("Imported font file (%s): %s -> %s", fontType, srcPath, dstPath)
		}
	}

	return nil
}

// ImportSettings imports settings from a theme
func ImportSettings(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Import accent settings if present
	if manifest.Content.Settings.AccentsIncluded {
		if accents, ok := manifest.PathMappings.Settings["accents"]; ok {
			srcPath := filepath.Join(themePath, accents.ThemePath)
			dstPath := accents.SystemPath

			// Check if source file exists
			if _, err := os.Stat(srcPath); os.IsNotExist(err) {
				logger.Printf("Warning: Source accent settings file not found: %s", srcPath)
			} else {
				// Create destination directory
				dstDir := filepath.Dir(dstPath)
				if err := os.MkdirAll(dstDir, 0755); err != nil {
					logger.Printf("Warning: Could not create directory for accent settings: %v", err)
				} else {
					// Copy the file
					if err := CopyFile(srcPath, dstPath); err != nil {
						logger.Printf("Warning: Could not copy accent settings: %v", err)
					} else {
						logger.Printf("Imported accent settings: %s -> %s", srcPath, dstPath)
					}
				}
			}
		}
	}

	// Import LED settings if present
	if manifest.Content.Settings.LEDsIncluded {
		if leds, ok := manifest.PathMappings.Settings["leds"]; ok {
			srcPath := filepath.Join(themePath, leds.ThemePath)
			dstPath := leds.SystemPath

			// Check if source file exists
			if _, err := os.Stat(srcPath); os.IsNotExist(err) {
				logger.Printf("Warning: Source LED settings file not found: %s", srcPath)
			} else {
				// Create destination directory
				dstDir := filepath.Dir(dstPath)
				if err := os.MkdirAll(dstDir, 0755); err != nil {
					logger.Printf("Warning: Could not create directory for LED settings: %v", err)
				} else {
					// Copy the file
					if err := CopyFile(srcPath, dstPath); err != nil {
						logger.Printf("Warning: Could not copy LED settings: %v", err)
					} else {
						logger.Printf("Imported LED settings: %s -> %s", srcPath, dstPath)
					}
				}
			}
		}
	}

	return nil
}

// ImportTheme imports a theme package
func ImportTheme(themeName string) error {
	// Create logging directory if it doesn't exist
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	logsDir := filepath.Join(cwd, "Logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("error creating logs directory: %w", err)
	}

	// Create log file
	logFile, err := os.OpenFile(
		filepath.Join(logsDir, "imports.log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()

	// Create logger
	logger := &Logger{log.New(logFile, "", log.LstdFlags)}
	logger.Printf("Starting theme import for: %s", themeName)

	// Full path to theme - look in Imports directory
	themePath := filepath.Join(cwd, "Themes", "Imports", themeName)

	// Validate theme
	manifest, err := ValidateTheme(themePath, logger)
	if err != nil {
		logger.Printf("Theme validation failed: %v", err)
		return fmt.Errorf("theme validation failed: %w", err)
	}

	// Import theme components
	if err := ImportWallpapers(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing wallpapers: %v", err)
	}

	if err := ImportIcons(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing icons: %v", err)
	}

	if err := ImportOverlays(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing overlays: %v", err)
	}

	if err := ImportFonts(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing fonts: %v", err)
	}

	if err := ImportSettings(themePath, manifest, logger); err != nil {
		logger.Printf("Error importing settings: %v", err)
	}

	logger.Printf("Theme import completed successfully: %s", themeName)

	// Show success message to user
	ui.ShowMessage(fmt.Sprintf("Theme '%s' by %s imported successfully!",
		manifest.ThemeInfo.Name, manifest.ThemeInfo.Author), "5")

	return nil
}