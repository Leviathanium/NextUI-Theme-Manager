// src/internal/themes/global.go
// Global theme operations

package themes

// Add at the top of global.go, after existing imports
import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"  // Add this for regex support
	"strings"

	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
)

// ListGlobalThemes returns a list of available global themes
func ListGlobalThemes(globalThemesDir string) ([]string, error) {
	var themes []string

	// Check if the directory exists
	_, err := os.Stat(globalThemesDir)
	if os.IsNotExist(err) {
		logging.LogDebug("Global themes directory does not exist: %s", globalThemesDir)
		return nil, fmt.Errorf("global themes directory does not exist: %s", globalThemesDir)
	} else if err != nil {
		logging.LogDebug("Error checking global themes directory: %v", err)
		return nil, fmt.Errorf("error checking themes directory: %w", err)
	}

	// Read the directory
	entries, err := os.ReadDir(globalThemesDir)
	if err != nil {
		logging.LogDebug("Error reading global themes directory: %v", err)
		return nil, fmt.Errorf("error reading themes directory: %w", err)
	}

	// First check for new style PNG files directly in the directory
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			// Check if it's a PNG file
			if strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				// Remove the .png extension to get the theme name
				themeName := strings.TrimSuffix(entry.Name(), ".png")
				themes = append(themes, themeName)
			}
		}
	}

	// If we found new style themes, return them
	if len(themes) > 0 {
		logging.LogDebug("Found %d global themes (new style)", len(themes))
		return themes, nil
	}

	// No new style themes found, try old directory style
	logging.LogDebug("No new style themes found, checking old directory style")

	// Find directories that contain a bg.png file (old style)
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			bgPath := filepath.Join(globalThemesDir, entry.Name(), "bg.png")
			if _, err := os.Stat(bgPath); err == nil {
				themes = append(themes, entry.Name())
			}
		}
	}

	logging.LogDebug("Found %d global themes (old style)", len(themes))
	return themes, nil
}

// ApplyGlobalTheme applies a global theme to all directories
func ApplyGlobalTheme(themeName string) error {
	logging.LogDebug("Applying global theme: %s", themeName)

	// Get current directory for absolute paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}
	logging.LogDebug("Current directory: %s", cwd)

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Ensure all media directories exist
	if err := system.EnsureMediaDirectories(systemPaths); err != nil {
		logging.LogDebug("Error ensuring media directories: %v", err)
		return fmt.Errorf("error ensuring media directories: %w", err)
	}

	// Check new style first - if a PNG file exists in Wallpapers root
	srcBg := filepath.Join(cwd, "Wallpapers", themeName + ".png")

	// If not found, check old directory style
	if _, err := os.Stat(srcBg); os.IsNotExist(err) {
		oldStylePath := filepath.Join(cwd, "Wallpapers", themeName, "bg.png")
		if _, err := os.Stat(oldStylePath); err == nil {
			srcBg = oldStylePath
			logging.LogDebug("Using old style directory wallpaper: %s", srcBg)
		} else {
			// Neither new nor old style found
			logging.LogDebug("Theme background not found in either format: %s", themeName)
			return fmt.Errorf("theme background not found: %s", themeName)
		}
	}

	logging.LogDebug("Theme background path: %s", srcBg)

	// Apply to root .media directory
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", "bg.png")
	logging.LogDebug("Copying to root .media: %s", rootMediaBg)
	if err := CopyFile(srcBg, rootMediaBg); err != nil {
		logging.LogDebug("Error copying to root .media: %v", err)
		return fmt.Errorf("failed to copy background to root .media: %w", err)
	}

	// Also apply to root directory (NextUI sometimes looks for bg.png in the root)
	rootBg := filepath.Join(systemPaths.Root, "bg.png")
	logging.LogDebug("Copying to root: %s", rootBg)
	if err := CopyFile(srcBg, rootBg); err != nil {
		logging.LogDebug("Error copying to root: %v", err)
		return fmt.Errorf("failed to copy background to root: %w", err)
	}

	// Apply to Recently Played
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
	logging.LogDebug("Copying to Recently Played: %s", rpBg)
	if err := CopyFile(srcBg, rpBg); err != nil {
		logging.LogDebug("Error copying to Recently Played: %v", err)
		return fmt.Errorf("failed to copy background to Recently Played: %w", err)
	}

	// Apply to Tools
	toolsBg := filepath.Join(systemPaths.Tools, ".media", "bg.png")
	logging.LogDebug("Copying to Tools: %s", toolsBg)
	if err := CopyFile(srcBg, toolsBg); err != nil {
		logging.LogDebug("Error copying to Tools: %v", err)
		return fmt.Errorf("failed to copy background to Tools: %w", err)
	}

	// Apply to all system directories
	for _, system := range systemPaths.Systems {
		systemBg := filepath.Join(system.MediaPath, "bg.png")
		logging.LogDebug("Copying to system %s: %s", system.Name, systemBg)
		if err := CopyFile(srcBg, systemBg); err != nil {
			logging.LogDebug("Error copying to system %s: %v", system.Name, err)
			return fmt.Errorf("failed to copy background to %s: %w", system.Name, err)
		}
	}

	logging.LogDebug("Global theme applied successfully")
	return nil
}

// ApplyCustomTheme applies a custom theme to a specific system
func ApplyCustomTheme(systemName string, themeName string) error {
	logging.LogDebug("Applying custom theme to system: %s, theme: %s", systemName, themeName)

	// Get current directory for absolute paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Try different possible wallpaper paths
	var srcBg string

	// Check for system-specific wallpaper in new format first
	if strings.HasPrefix(systemName, "(") && strings.HasSuffix(systemName, ")") {
		// This is a tag-only name, look for direct tag match in SystemWallpapers
		tagOnly := systemName[1 : len(systemName)-1]

		// First check if there's an exact match with this tag
		for _, sys := range systemPaths.Systems {
			if sys.Tag == tagOnly {
				specificWallpaper := filepath.Join(cwd, "Wallpapers", "SystemWallpapers", fmt.Sprintf("%s (%s).png", sys.Name, sys.Tag))
				if _, err := os.Stat(specificWallpaper); err == nil {
					srcBg = specificWallpaper
					logging.LogDebug("Found system-specific wallpaper for tag %s: %s", tagOnly, srcBg)
					break
				}
			}
		}

		// If not found, check old-style directory structure with tag
		if srcBg == "" {
			oldStylePath := filepath.Join(cwd, "Wallpapers", themeName, "Systems", systemName, "bg.png")
			if _, err := os.Stat(oldStylePath); err == nil {
				srcBg = oldStylePath
				logging.LogDebug("Using old style tag directory wallpaper: %s", srcBg)
			}
		}
	} else {
		// Look for system-specific wallpaper by name
		var systemTag string

		// Try to match system by name to get its tag
		for _, sys := range systemPaths.Systems {
			if sys.Name == systemName {
				systemTag = sys.Tag
				break
			}
		}

		// Check for system-specific wallpaper using full name and tag
		if systemTag != "" {
			specificWallpaper := filepath.Join(cwd, "Wallpapers", "SystemWallpapers", fmt.Sprintf("%s (%s).png", systemName, systemTag))
			if _, err := os.Stat(specificWallpaper); err == nil {
				srcBg = specificWallpaper
				logging.LogDebug("Found system-specific wallpaper: %s", srcBg)
			}
		}

		// For special systems like Root, Recently Played, etc.
		if srcBg == "" && (systemName == "Root" || systemName == "Recently Played" ||
		                   systemName == "Tools" || systemName == "Collections") {
			specificWallpaper := filepath.Join(cwd, "Wallpapers", "SystemWallpapers", systemName + ".png")
			if _, err := os.Stat(specificWallpaper); err == nil {
				srcBg = specificWallpaper
				logging.LogDebug("Found specific wallpaper for %s: %s", systemName, srcBg)
			}
		}

		// Check for collection-specific wallpaper
		if srcBg == "" && !strings.Contains(systemName, "(") {
			// This might be a collection
			collWallpaper := filepath.Join(cwd, "Wallpapers", "CollectionWallpapers", systemName + ".png")
			if _, err := os.Stat(collWallpaper); err == nil {
				srcBg = collWallpaper
				logging.LogDebug("Found collection-specific wallpaper: %s", srcBg)
			}
		}
	}

	// If we still don't have a source, try the theme's main wallpaper
	if srcBg == "" {
		// Try new style theme (direct PNG file)
		newStylePath := filepath.Join(cwd, "Wallpapers", themeName + ".png")
		if _, err := os.Stat(newStylePath); err == nil {
			srcBg = newStylePath
			logging.LogDebug("Using theme's main wallpaper (new style): %s", srcBg)
		} else {
			// Try old style theme (directory with bg.png)
			oldStylePath := filepath.Join(cwd, "Wallpapers", themeName, "bg.png")
			if _, err := os.Stat(oldStylePath); err == nil {
				srcBg = oldStylePath
				logging.LogDebug("Using theme's main wallpaper (old style): %s", srcBg)
			}
		}
	}

	// If we still don't have a source, try old directories
	if srcBg == "" {
		// Check old style path with system directory
		if strings.HasPrefix(systemName, "(") && strings.HasSuffix(systemName, ")") {
			oldStylePath := filepath.Join(cwd, "Wallpapers", themeName, "Systems", systemName, "bg.png")
			if _, err := os.Stat(oldStylePath); err == nil {
				srcBg = oldStylePath
				logging.LogDebug("Using old style tag directory wallpaper: %s", srcBg)
			}
		} else if systemName == "Root" {
			oldStylePath := filepath.Join(cwd, "Wallpapers", themeName, "Root", "bg.png")
			if _, err := os.Stat(oldStylePath); err == nil {
				srcBg = oldStylePath
				logging.LogDebug("Using old style root directory wallpaper: %s", srcBg)
			}
		} else if systemName == "Recently Played" {
			oldStylePath := filepath.Join(cwd, "Wallpapers", themeName, "Recently Played", "bg.png")
			if _, err := os.Stat(oldStylePath); err == nil {
				srcBg = oldStylePath
				logging.LogDebug("Using old style recently played directory wallpaper: %s", srcBg)
			}
		} else if systemName == "Tools" {
			oldStylePath := filepath.Join(cwd, "Wallpapers", themeName, "Tools", "bg.png")
			if _, err := os.Stat(oldStylePath); err == nil {
				srcBg = oldStylePath
				logging.LogDebug("Using old style tools directory wallpaper: %s", srcBg)
			}
		} else if systemName == "Collections" {
			oldStylePath := filepath.Join(cwd, "Wallpapers", themeName, "Collections", "bg.png")
			if _, err := os.Stat(oldStylePath); err == nil {
				srcBg = oldStylePath
				logging.LogDebug("Using old style collections directory wallpaper: %s", srcBg)
			}
		}
	}

	// Final check to make sure we have a source
	if srcBg == "" {
		logging.LogDebug("Could not find a wallpaper source for system %s in theme %s", systemName, themeName)
		return fmt.Errorf("wallpaper source not found for system %s in theme %s", systemName, themeName)
	}

	logging.LogDebug("Theme background path: %s", srcBg)

	// Determine target directory
	var targetPath string
	var targetMediaPath string

	if systemName == "Root" {
		targetPath = systemPaths.Root
		targetMediaPath = filepath.Join(targetPath, ".media")

		// Ensure media directory exists
		if err := os.MkdirAll(targetMediaPath, 0755); err != nil {
			logging.LogDebug("Error creating media directory: %v", err)
			return fmt.Errorf("failed to create media directory: %w", err)
		}

		// Apply background to .media directory
		dstBg := filepath.Join(targetMediaPath, "bg.png")
		if err := CopyFile(srcBg, dstBg); err != nil {
			logging.LogDebug("Error copying background: %v", err)
			return fmt.Errorf("failed to copy background: %w", err)
		}

		// Also copy to the root directory itself
		rootBg := filepath.Join(targetPath, "bg.png")
		if err := CopyFile(srcBg, rootBg); err != nil {
			logging.LogDebug("Error copying to root: %v", err)
			return fmt.Errorf("failed to copy background to root: %w", err)
		}

	} else if systemName == "Recently Played" {
		targetPath = systemPaths.RecentlyPlayed
		targetMediaPath = filepath.Join(targetPath, ".media")

		// Ensure media directory exists
		if err := os.MkdirAll(targetMediaPath, 0755); err != nil {
			logging.LogDebug("Error creating media directory: %v", err)
			return fmt.Errorf("failed to create media directory: %w", err)
		}

		// Apply background
		dstBg := filepath.Join(targetMediaPath, "bg.png")
		if err := CopyFile(srcBg, dstBg); err != nil {
			logging.LogDebug("Error copying background: %v", err)
			return fmt.Errorf("failed to copy background: %w", err)
		}

	} else if systemName == "Tools" {
		targetPath = systemPaths.Tools
		targetMediaPath = filepath.Join(targetPath, ".media")

		// Ensure media directory exists
		if err := os.MkdirAll(targetMediaPath, 0755); err != nil {
			logging.LogDebug("Error creating media directory: %v", err)
			return fmt.Errorf("failed to create media directory: %w", err)
		}

		// Apply background
		dstBg := filepath.Join(targetMediaPath, "bg.png")
		if err := CopyFile(srcBg, dstBg); err != nil {
			logging.LogDebug("Error copying background: %v", err)
			return fmt.Errorf("failed to copy background: %w", err)
		}

	} else if systemName == "Collections" {
		targetPath = filepath.Join(systemPaths.Root, "Collections")
		targetMediaPath = filepath.Join(targetPath, ".media")

		// Ensure media directory exists
		if err := os.MkdirAll(targetMediaPath, 0755); err != nil {
			logging.LogDebug("Error creating media directory: %v", err)
			return fmt.Errorf("failed to create media directory: %w", err)
		}

		// Apply background
		dstBg := filepath.Join(targetMediaPath, "bg.png")
		if err := CopyFile(srcBg, dstBg); err != nil {
			logging.LogDebug("Error copying background: %v", err)
			return fmt.Errorf("failed to copy background: %w", err)
		}
	} else {
		// Check if this is a collection name without explicit "Collections" prefix
		collPath := filepath.Join(systemPaths.Root, "Collections", systemName)
		if _, err := os.Stat(collPath); err == nil {
			// This is a collection directory
			targetPath = collPath
			targetMediaPath = filepath.Join(targetPath, ".media")

			// Ensure media directory exists
			if err := os.MkdirAll(targetMediaPath, 0755); err != nil {
				logging.LogDebug("Error creating collection media directory: %v", err)
				return fmt.Errorf("failed to create collection media directory: %w", err)
			}

			// Apply background
			dstBg := filepath.Join(targetMediaPath, "bg.png")
			if err := CopyFile(srcBg, dstBg); err != nil {
				logging.LogDebug("Error copying background: %v", err)
				return fmt.Errorf("failed to copy background: %w", err)
			}
		} else {
			// Find the system in our list - try multiple matching approaches
			found := false
			var matchedSystem system.SystemInfo

			// First, try exact match with system name
			for _, sys := range systemPaths.Systems {
				if sys.Name == systemName {
					matchedSystem = sys
					found = true
					logging.LogDebug("Found exact name match for system: %s", sys.Name)
					break
				}
			}

			// If not found by name, try matching by tag
			if !found {
				// Check if the systemName is a tag in parentheses
				if strings.HasPrefix(systemName, "(") && strings.HasSuffix(systemName, ")") {
					tagOnly := systemName[1 : len(systemName)-1]

					// Search for a system with this tag
					for _, sys := range systemPaths.Systems {
						if sys.Tag == tagOnly {
							matchedSystem = sys
							found = true
							logging.LogDebug("Found system by tag in parentheses: %s (Tag: %s)", sys.Name, tagOnly)
							break
						}
					}
				} else {
					// Check if the systemName contains a tag that we can extract
					re := regexp.MustCompile(`\((.*?)\)`)
					matches := re.FindStringSubmatch(systemName)
					if len(matches) >= 2 {
						extractedTag := matches[1]
						logging.LogDebug("Extracted tag from system name: %s", extractedTag)

						// Search for a system with this tag
						for _, sys := range systemPaths.Systems {
							if sys.Tag == extractedTag {
								matchedSystem = sys
								found = true
								logging.LogDebug("Found system by extracted tag: %s (Tag: %s)", sys.Name, extractedTag)
								break
							}
						}
					}
				}
			}

			if !found {
				logging.LogDebug("System not found by any matching method: %s", systemName)
				return fmt.Errorf("system not found: %s", systemName)
			}

			targetPath = matchedSystem.Path
			targetMediaPath = matchedSystem.MediaPath

			// Ensure media directory exists
			if err := os.MkdirAll(targetMediaPath, 0755); err != nil {
				logging.LogDebug("Error creating media directory: %v", err)
				return fmt.Errorf("failed to create media directory: %w", err)
			}

			// Apply background
			dstBg := filepath.Join(targetMediaPath, "bg.png")
			if err := CopyFile(srcBg, dstBg); err != nil {
				logging.LogDebug("Error copying background: %v", err)
				return fmt.Errorf("failed to copy background: %w", err)
			}
		}
	}

	logging.LogDebug("Custom theme applied successfully for %s", systemName)
	return nil
}

// DisplayThemeSelectionList is a helper function for theme selection
func DisplayThemeSelectionList(themes []string, title string) (string, int) {
	// Implement using UI functions from common.go
	// This is a simple wrapper to avoid circular imports
	return displayMinUiList(strings.Join(themes, "\n"), "text", title)
}

// Placeholder implementation - will be replaced by actual UI function
func displayMinUiList(list string, format string, title string, extraArgs ...string) (string, int) {
	// This is a placeholder that would be replaced by the actual UI function
	// In a real implementation, this would either be provided by an interface
	// or we would need to restructure to avoid circular dependencies
	return "", 0
}