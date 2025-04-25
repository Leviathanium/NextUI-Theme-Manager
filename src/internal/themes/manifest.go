// src/internal/themes/manifest.go
// Data structures and functions for theme manifest handling

package themes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ThemeManifest represents the manifest.json file structure
type ThemeManifest struct {
	ThemeInfo struct {
		Name         string    `json:"name"`
		Version      string    `json:"version"`
		Author       string    `json:"author"`
		CreationDate time.Time `json:"creation_date"`
		ExportedBy   string    `json:"exported_by"`
	} `json:"theme_info"`
	Content struct {
		Wallpapers struct {
			Present bool `json:"present"`
			Count   int  `json:"count"`
		} `json:"wallpapers"`
		Icons struct {
			Present         bool `json:"present"`
			SystemCount     int  `json:"system_count"`
			ToolCount       int  `json:"tool_count"`
			CollectionCount int  `json:"collection_count"`
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
		Wallpapers []PathMapping          `json:"wallpapers"`
		Icons      []PathMapping          `json:"icons"`
		Overlays   []PathMapping          `json:"overlays"`
		Fonts      map[string]PathMapping `json:"fonts"`
		Settings   map[string]PathMapping `json:"settings"`
	} `json:"path_mappings"`
	AccentColors struct {
		Color1 string `json:"color1"`
		Color2 string `json:"color2"`
		Color3 string `json:"color3"`
		Color4 string `json:"color4"`
		Color5 string `json:"color5"`
		Color6 string `json:"color6"`
	} `json:"accent_colors"`

	// Add LED settings
	LEDSettings struct {
		F1Key      LEDSetting `json:"f1_key"`
		F2Key      LEDSetting `json:"f2_key"`
		TopBar     LEDSetting `json:"top_bar"`
		LRTriggers LEDSetting `json:"lr_triggers"`
	} `json:"led_settings"`
}

// PathMapping represents a mapping between theme and system paths
type PathMapping struct {
	ThemePath  string            `json:"theme_path"`
	SystemPath string            `json:"system_path"`
	Metadata   map[string]string `json:"metadata,omitempty"` // Additional metadata to aid in matching
}

type LEDSetting struct {
	Effect       int    `json:"effect"`
	Color1       string `json:"color1"`
	Color2       string `json:"color2"`
	Speed        int    `json:"speed"`
	Brightness   int    `json:"brightness"`
	Trigger      int    `json:"trigger"`
	InBrightness int    `json:"in_brightness"`
}

// Logger is a simple wrapper for logging
type Logger struct {
	DebugFn func(format string, args ...interface{})
}

// ThemeVersionInfo holds the theme manager version information
type ThemeVersionInfo struct {
	Major int
	Minor int
	Patch int
}

// Current version of the Theme Manager
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
	// Only set creation date and exported_by
	manifest.ThemeInfo.CreationDate = time.Now()
	manifest.ThemeInfo.ExportedBy = GetVersionString()

	// Only set version if not already set
	if manifest.ThemeInfo.Version == "" {
		manifest.ThemeInfo.Version = "1.0.0"
	}

	// Only set author if not already set
	if manifest.ThemeInfo.Author == "" {
		manifest.ThemeInfo.Author = "AuthorName" // Default author name only if none exists
	}

	// Extract theme name from directory name if not set
	if manifest.ThemeInfo.Name == "" {
		themeName := filepath.Base(themePath)
		manifest.ThemeInfo.Name = themeName
	}

	// Use an encoder that doesn't escape HTML characters
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false) // This prevents & from becoming \u0026
	enc.SetIndent("", "  ")  // Add proper indentation

	if err := enc.Encode(manifest); err != nil {
		logger.DebugFn("Error creating manifest JSON: %v", err)
		return fmt.Errorf("error creating manifest JSON: %w", err)
	}

	// Write manifest to file
	manifestPath := filepath.Join(themePath, "manifest.json")
	if err := os.WriteFile(manifestPath, buf.Bytes(), 0644); err != nil {
		logger.DebugFn("Error writing manifest file: %v", err)
		return fmt.Errorf("error writing manifest file: %w", err)
	}

	logger.DebugFn("Created manifest file: %s", manifestPath)
	return nil
}

// ValidateTheme validates a theme package and returns its manifest
func ValidateTheme(themePath string, logger *Logger) (*ThemeManifest, error) {
	logger.DebugFn("Validating theme at: %s", themePath)

	// Check if the theme directory exists
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		logger.DebugFn("Theme directory does not exist: %s", themePath)
		return nil, fmt.Errorf("theme directory does not exist: %s", themePath)
	}

	// Check for manifest.json
	manifestPath := filepath.Join(themePath, "manifest.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		logger.DebugFn("Manifest file not found: %s", manifestPath)
		return nil, fmt.Errorf("manifest file not found: %s", manifestPath)
	}

	// Read and parse manifest
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		logger.DebugFn("Error reading manifest file: %v", err)
		return nil, fmt.Errorf("error reading manifest file: %w", err)
	}

	var manifest ThemeManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		logger.DebugFn("Error parsing manifest JSON: %v", err)
		return nil, fmt.Errorf("error parsing manifest JSON: %w", err)
	}

	logger.DebugFn("Theme validation successful, name: %s, version: %s, author: %s",
		manifest.ThemeInfo.Name, manifest.ThemeInfo.Version, manifest.ThemeInfo.Author)

	return &manifest, nil
}
