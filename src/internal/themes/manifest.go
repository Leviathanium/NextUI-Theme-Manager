// src/internal/themes/manifest.go
// Data structures and functions for theme manifest handling

package themes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ThemeManifest represents the manifest.json file structure
type ThemeManifest struct {
	ThemeInfo struct {
		Name         string    `json:"name"`
		Version      string    `json:"version"`
		Author       string    `json:"author"`         // Author of the theme
		CreationDate time.Time `json:"creation_date"`
		ExportedBy   string    `json:"exported_by"`
	} `json:"theme_info"`
	Content struct {
		Wallpapers struct {
			Present bool `json:"present"`
			Count   int  `json:"count"`
		} `json:"wallpapers"`
		Icons struct {
			Present        bool `json:"present"`
			SystemCount    int  `json:"system_count"`
			ToolCount      int  `json:"tool_count"`
			CollectionCount int `json:"collection_count"`
		} `json:"icons"`
		Overlays struct {
			Present bool     `json:"present"`
			Systems []string `json:"systems"`
		} `json:"overlays"`
		Fonts struct {
			Present      bool `json:"present"`
			OGReplaced   bool `json:"og_replaced"`
			NextReplaced bool `json:"next_replaced"`
		} `json:"fonts"`
		Settings struct {
			AccentsIncluded bool `json:"accents_included"`
			LEDsIncluded    bool `json:"leds_included"`
		} `json:"settings"`
	} `json:"content"`
	PathMappings struct {
		Wallpapers []PathMapping           `json:"wallpapers"`
		Icons      []PathMapping           `json:"icons"`
		Overlays   []PathMapping           `json:"overlays"`
		Fonts      map[string]PathMapping  `json:"fonts"`
		Settings   map[string]PathMapping  `json:"settings"`
	} `json:"path_mappings"`
	AccentColors map[string]string         `json:"accent_colors"`
}

// PathMapping represents a mapping between theme and system paths
type PathMapping struct {
	ThemePath  string `json:"theme_path"`
	SystemPath string `json:"system_path"`
}

// ThemeVersionInfo holds the theme manager version information
type ThemeVersionInfo struct {
	Major int
	Minor int
	Patch int
}

// Current version of the Theme Manager
// This should be updated when releasing new versions
var CurrentVersion = ThemeVersionInfo{
	Major: 1,
	Minor: 0,
	Patch: 0,
}

// GetVersionString returns the current theme manager version as a string
func GetVersionString() string {
	return fmt.Sprintf("Theme Manager v%d.%d.%d",
		CurrentVersion.Major, CurrentVersion.Minor, CurrentVersion.Patch)
}

// WriteManifest writes the manifest to a file in the theme directory
func WriteManifest(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Set creation date, version, author and exported_by
	manifest.ThemeInfo.CreationDate = time.Now()
	manifest.ThemeInfo.Version = "1.0.0"
	manifest.ThemeInfo.Author = "AuthorName" // Default author name as requested
	manifest.ThemeInfo.ExportedBy = GetVersionString()

	// Extract theme name from directory name
	themeName := filepath.Base(themePath)
	manifest.ThemeInfo.Name = strings.TrimSuffix(themeName, ".theme")

	// Convert manifest to JSON
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		logger.Printf("Error creating manifest JSON: %v", err)
		return fmt.Errorf("error creating manifest JSON: %w", err)
	}

	// Write manifest to file
	manifestPath := filepath.Join(themePath, "manifest.json")
	if err := os.WriteFile(manifestPath, manifestJSON, 0644); err != nil {
		logger.Printf("Error writing manifest file: %v", err)
		return fmt.Errorf("error writing manifest file: %w", err)
	}

	logger.Printf("Created manifest file: %s", manifestPath)
	return nil
}

// ValidateTheme validates a theme package and returns its manifest
func ValidateTheme(themePath string, logger *Logger) (*ThemeManifest, error) {
	logger.Printf("Validating theme at: %s", themePath)

	// Check if the theme directory exists
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		logger.Printf("Theme directory does not exist: %s", themePath)
		return nil, fmt.Errorf("theme directory does not exist: %s", themePath)
	}

	// Check for manifest.json
	manifestPath := filepath.Join(themePath, "manifest.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		logger.Printf("Manifest file not found: %s", manifestPath)
		return nil, fmt.Errorf("manifest file not found: %s", manifestPath)
	}

	// Read and parse manifest
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		logger.Printf("Error reading manifest file: %v", err)
		return nil, fmt.Errorf("error reading manifest file: %w", err)
	}

	var manifest ThemeManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		logger.Printf("Error parsing manifest JSON: %v", err)
		return nil, fmt.Errorf("error parsing manifest JSON: %w", err)
	}

	logger.Printf("Theme validation successful, name: %s, version: %s, author: %s",
		manifest.ThemeInfo.Name, manifest.ThemeInfo.Version, manifest.ThemeInfo.Author)

	return &manifest, nil
}