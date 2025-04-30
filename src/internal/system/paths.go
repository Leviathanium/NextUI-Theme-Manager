// internal/system/paths.go
package system

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Base system paths
const (
	// Root SD card path
	RootPath = "/mnt/SDCARD"

	// ROM directories
	RomsPath = "/mnt/SDCARD/Roms"
	RecentlyPlayedPath = "/mnt/SDCARD/Recently Played"

	// Tool directories
	ToolsPath = "/mnt/SDCARD/Tools/tg5040"
	ToolsParentPath = "/mnt/SDCARD/Tools"

	// Collections
	CollectionsPath = "/mnt/SDCARD/Collections"

	// Media directories for main screens
	RootMediaPath = "/mnt/SDCARD/.media"
	RomsMediaPath = "/mnt/SDCARD/Roms/.media"

	// Overlays
	OverlaysPath = "/mnt/SDCARD/Overlays"

	// System resource path
	SystemResPath = "/mnt/SDCARD/.system/res"

	// Settings
	UserSettingsPath = "/mnt/SDCARD/.userdata/shared"

	// Files
	FontNextPath = "/mnt/SDCARD/.system/res/font1.ttf"
	FontNextBackupPath = "/mnt/SDCARD/.system/res/font1.backup.ttf"
	FontOGPath = "/mnt/SDCARD/.system/res/font2.ttf"
	FontOGBackupPath = "/mnt/SDCARD/.system/res/font2.backup.ttf"
	AccentSettingsPath = "/mnt/SDCARD/.userdata/shared/minuisettings.txt"
	LEDSettingsPath = "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"
)

// File type identifiers
const (
	ThemeExtension = ".theme"
	OverlayExtension = ".over"
)

// Theme structure paths
const (
	// Backgrounds
	ThemeWallpapersDir        = "Wallpapers"
	ThemeSystemWallpapersDir  = "Wallpapers/SystemWallpapers"
	ThemeListWallpapersDir    = "Wallpapers/ListWallpapers"
	ThemeCollectionWallpapersDir = "Wallpapers/CollectionWallpapers"

	// Icons
	ThemeIconsDir             = "Icons"
	ThemeSystemIconsDir       = "Icons/SystemIcons"
	ThemeToolIconsDir         = "Icons/ToolIcons"
	ThemeCollectionIconsDir   = "Icons/CollectionIcons"

	// Others
	ThemeOverlaysDir          = "Overlays"
	ThemeFontsDir             = "Fonts"
)

// PathMapping represents a mapping between theme and system paths
type PathMapping struct {
	ThemePath  string            // Path relative to theme root
	SystemPath string            // Absolute path on system
	Metadata   map[string]string // Additional metadata to aid in mapping
}

// GetSystemMediaPath returns the media path for a given system
func GetSystemMediaPath(systemName string) string {
	return filepath.Join(RomsPath, systemName, ".media")
}

// GetSystemBackgroundPath returns the path for a system's background image
func GetSystemBackgroundPath(systemName string) string {
	return filepath.Join(GetSystemMediaPath(systemName), "bg.png")
}

// GetSystemListBackgroundPath returns the path for a system's list background image
func GetSystemListBackgroundPath(systemName string) string {
	return filepath.Join(GetSystemMediaPath(systemName), "bglist.png")
}

// GetSystemIconPath returns the path for a system icon
func GetSystemIconPath(systemName, systemTag string) string {
	iconName := systemName
	if systemTag != "" && !strings.Contains(systemName, fmt.Sprintf("(%s)", systemTag)) {
		iconName = fmt.Sprintf("%s (%s)", systemName, systemTag)
	}
	return filepath.Join(RomsMediaPath, iconName+".png")
}

// GetRootBackgroundPath returns the path for the root background image
func GetRootBackgroundPath(mediaDir bool) string {
	if mediaDir {
		return filepath.Join(RootMediaPath, "bg.png")
	}
	return filepath.Join(RootPath, "bg.png")
}

// GetToolIconPath returns the path for a tool icon
func GetToolIconPath(toolName string) string {
	return filepath.Join(ToolsPath, toolName, ".media", toolName+".png")
}

// GetCollectionIconPath returns the path for a collection icon
func GetCollectionIconPath(collectionName string) string {
	return filepath.Join(CollectionsPath, collectionName, ".media", collectionName+".png")
}

// GetCollectionBackgroundPath returns the path for a collection background
func GetCollectionBackgroundPath(collectionName string) string {
	return filepath.Join(CollectionsPath, collectionName, ".media", "bg.png")
}

// GetOverlaySystemPath returns the path for system overlays
func GetOverlaySystemPath(systemTag string) string {
	return filepath.Join(OverlaysPath, systemTag)
}

// GetOverlayFilePath returns the path for a specific overlay file
func GetOverlayFilePath(systemTag, fileName string) string {
	return filepath.Join(GetOverlaySystemPath(systemTag), fileName)
}

// GetRecentlyPlayedBackgroundPath returns the path for Recently Played background
func GetRecentlyPlayedBackgroundPath() string {
	return filepath.Join(RecentlyPlayedPath, ".media", "bg.png")
}

// GetRecentlyPlayedIconPath returns the path for Recently Played icon
func GetRecentlyPlayedIconPath() string {
	return filepath.Join(RootMediaPath, "Recently Played.png")
}

// GetToolsBackgroundPath returns the path for Tools background
func GetToolsBackgroundPath() string {
	return filepath.Join(ToolsPath, ".media", "bg.png")
}

// GetToolsIconPath returns the path for Tools icon
func GetToolsIconPath() string {
	return filepath.Join(ToolsParentPath, ".media", "tg5040.png")
}

// GetCollectionsIconPath returns the path for Collections icon
func GetCollectionsIconPath() string {
	return filepath.Join(RootMediaPath, "Collections.png")
}

// GetCollectionsBackgroundPath returns the path for Collections background
func GetCollectionsBackgroundPath() string {
	return filepath.Join(CollectionsPath, ".media", "bg.png")
}

// ExtractSystemTag extracts the system tag from a filename
func ExtractSystemTag(filename string) string {
	re := regexp.MustCompile(`\((.*?)\)`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// GetThemeManifestPath returns the path to a theme's manifest file
func GetThemeManifestPath(themePath string) string {
	return filepath.Join(themePath, "manifest.yml")
}

// GetThemePreviewPath returns the path to a theme's preview image
func GetThemePreviewPath(themePath string) string {
	return filepath.Join(themePath, "preview.png")
}

// GetSystemPathForThemeFile determines the system path for a theme file
func GetSystemPathForThemeFile(themePath string, file string, systemTag string) string {
	// Extract filename
	filename := filepath.Base(file)

	// Handle different file types based on path
	switch {
	case strings.HasPrefix(file, ThemeSystemWallpapersDir):
		// Special cases for system wallpapers
		switch filename {
		case "Root.png":
			return GetRootBackgroundPath(false)
		case "Root-Media.png":
			return GetRootBackgroundPath(true)
		case "Recently Played.png":
			return GetRecentlyPlayedBackgroundPath()
		case "Tools.png":
			return GetToolsBackgroundPath()
		case "Collections.png":
			return GetCollectionsBackgroundPath()
		default:
			// Extract system name from filename
			systemName := strings.TrimSuffix(filename, ".png")
			tag := ExtractSystemTag(systemName)
			if tag == "" {
				tag = systemTag
			}

			// If there's a tag in the filename, use it
			return GetSystemBackgroundPath(systemName)
		}

	case strings.HasPrefix(file, ThemeListWallpapersDir):
		// Handle list backgrounds
		systemName := strings.TrimSuffix(filename, "-list.png")
		tag := ExtractSystemTag(systemName)
		if tag == "" {
			tag = systemTag
		}

		return GetSystemListBackgroundPath(systemName)

	case strings.HasPrefix(file, ThemeSystemIconsDir):
		// Special cases for system icons
		switch filename {
		case "Recently Played.png":
			return GetRecentlyPlayedIconPath()
		case "Tools.png":
			return GetToolsIconPath()
		case "Collections.png":
			return GetCollectionsIconPath()
		default:
			// Extract system name from filename
			systemName := strings.TrimSuffix(filename, ".png")
			tag := ExtractSystemTag(systemName)
			if tag == "" {
				tag = systemTag
			}

			return GetSystemIconPath(systemName, tag)
		}

	case strings.HasPrefix(file, ThemeToolIconsDir):
		// Handle tool icons
		toolName := strings.TrimSuffix(filename, ".png")
		return GetToolIconPath(toolName)

	case strings.HasPrefix(file, ThemeCollectionIconsDir):
		// Handle collection icons
		collectionName := strings.TrimSuffix(filename, ".png")
		return GetCollectionIconPath(collectionName)

	case strings.HasPrefix(file, ThemeCollectionWallpapersDir):
		// Handle collection wallpapers
		collectionName := strings.TrimSuffix(filename, ".png")
		return GetCollectionBackgroundPath(collectionName)

	case strings.HasPrefix(file, ThemeOverlaysDir):
		// Handle overlays - they're organized by system tag subdirectories
		parts := strings.Split(file, "/")
		if len(parts) >= 3 {
			overlaySystemTag := parts[1]
			overlayFileName := parts[2]
			return GetOverlayFilePath(overlaySystemTag, overlayFileName)
		}

	case strings.HasPrefix(file, ThemeFontsDir):
		// Handle fonts
		switch filename {
		case "OG.ttf":
			return FontOGPath
		case "OG.backup.ttf":
			return FontOGBackupPath
		case "Next.ttf":
			return FontNextPath
		case "Next.backup.ttf":
			return FontNextBackupPath
		}
	}

	// If no special case handled the file, return empty string
	return ""
}