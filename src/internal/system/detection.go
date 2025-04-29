// internal/system/detection.go
package system

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SystemInfo represents information about an installed system
type SystemInfo struct {
	Name      string // Full name including tag
	Tag       string // Just the tag (e.g., "GBA")
	Path      string // Full path to the system directory
	MediaPath string // Path to the .media directory
}

// SystemPaths contains paths for standard system directories
type SystemPaths struct {
	Root           string
	RecentlyPlayed string
	Tools          string
	Roms           string
	Systems        []SystemInfo
}

// GetSystemPaths returns the paths to all system directories
func GetSystemPaths() (*SystemPaths, error) {
	// Define base paths
	rootPath := "/mnt/SDCARD"
	recentlyPath := filepath.Join(rootPath, "Recently Played")
	toolsPath := filepath.Join(rootPath, "Tools/tg5040")
	romsPath := filepath.Join(rootPath, "Roms")

	// Create the result structure
	systemPaths := &SystemPaths{
		Root:           rootPath,
		RecentlyPlayed: recentlyPath,
		Tools:          toolsPath,
		Roms:           romsPath,
		Systems:        []SystemInfo{},
	}

	// Scan for ROM system directories
	romsDirs, err := os.ReadDir(romsPath)
	if err != nil {
		return nil, err
	}

	// Regular expression to extract system tag from directory name
	re := regexp.MustCompile(`\((.*?)\)`)

	for _, dir := range romsDirs {
		if dir.IsDir() && dir.Name() != ".media" && !strings.HasPrefix(dir.Name(), ".") {
			systemPath := filepath.Join(romsPath, dir.Name())
			mediaPath := filepath.Join(systemPath, ".media")

			// Extract system tag from directory name if present
			tag := ""
			matches := re.FindStringSubmatch(dir.Name())
			if len(matches) >= 2 {
				tag = matches[1]
			}

			// Add to our list of systems
			systemPaths.Systems = append(systemPaths.Systems, SystemInfo{
				Name:      dir.Name(),
				Tag:       tag,
				Path:      systemPath,
				MediaPath: mediaPath,
			})
		}
	}

	return systemPaths, nil
}

// EnsureMediaDirectories ensures that all necessary .media directories exist
func EnsureMediaDirectories(paths *SystemPaths) error {
	// Ensure Root .media directory
	rootMediaPath := filepath.Join(paths.Root, ".media")
	if err := os.MkdirAll(rootMediaPath, 0755); err != nil {
		return err
	}

	// Ensure Recently Played .media directory
	rpMediaPath := filepath.Join(paths.RecentlyPlayed, ".media")
	if err := os.MkdirAll(rpMediaPath, 0755); err != nil {
		return err
	}

	// Ensure Tools .media directory
	toolsMediaPath := filepath.Join(paths.Tools, ".media")
	if err := os.MkdirAll(toolsMediaPath, 0755); err != nil {
		return err
	}

	// Ensure System .media directories
	for _, system := range paths.Systems {
		if err := os.MkdirAll(system.MediaPath, 0755); err != nil {
			return err
		}
	}

	// Ensure Overlays directory
	overlaysPath := filepath.Join(paths.Root, "Overlays")
	if err := os.MkdirAll(overlaysPath, 0755); err != nil {
		return err
	}

	return nil
}

// CreateOverlaySystemDirectory creates a system directory in the Overlays folder
func CreateOverlaySystemDirectory(paths *SystemPaths, systemTag string) error {
	overlaySystemPath := filepath.Join(paths.Root, "Overlays", systemTag)
	return os.MkdirAll(overlaySystemPath, 0755)
}

// GetSystemByTag returns a system info by its tag
func GetSystemByTag(paths *SystemPaths, tag string) *SystemInfo {
	for _, system := range paths.Systems {
		if system.Tag == tag {
			return &system
		}
	}
	return nil
}

// GetAllSystemTags returns a list of all system tags
func GetAllSystemTags(paths *SystemPaths) []string {
	var tags []string
	for _, system := range paths.Systems {
		if system.Tag != "" {
			tags = append(tags, system.Tag)
		}
	}
	return tags
}