// internal/themes/manifest.go
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
	"thememanager/internal/logging"
	"thememanager/internal/system"
)

// ThemeManifest represents the structure of a theme manifest.yml file
type ThemeManifest struct {
	Info struct {
		Name        string    `yaml:"name"`
		Author      string    `yaml:"author"`
		Version     string    `yaml:"version"`
		CreatedDate time.Time `yaml:"created_date"`
	} `yaml:"info"`

	Content struct {
		SystemsCount int  `yaml:"systems_count"`
		Backgrounds  bool `yaml:"backgrounds"`
		Icons        bool `yaml:"icons"`
		Fonts        bool `yaml:"fonts"`
		Accents      bool `yaml:"accents"`
	} `yaml:"content"`

	Systems map[string]SystemConfig `yaml:"systems"`
}

// SystemConfig represents the configuration for a specific system in the manifest
type SystemConfig struct {
	DisplayName string            `yaml:"display_name"`
	Files       map[string]string `yaml:"files"`
	Paths       map[string]string `yaml:"paths"`
}

// OverlayManifest represents the structure of an overlay manifest.yml file
type OverlayManifest struct {
	Info struct {
		Name        string    `yaml:"name"`
		Author      string    `yaml:"author"`
		Version     string    `yaml:"version"`
		CreatedDate time.Time `yaml:"created_date"`
	} `yaml:"info"`

	Content struct {
		Systems []string `yaml:"systems"`
	} `yaml:"content"`

	// Additional overlay-specific fields can be added here
}

// ReadThemeManifest reads and parses a theme manifest.yml file
func ReadThemeManifest(path string) (*ThemeManifest, error) {
	logging.LogDebug("Reading theme manifest from: %s", path)
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest file: %w", err)
	}

	var manifest ThemeManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("error parsing manifest: %w", err)
	}

	logging.LogDebug("Successfully read manifest for theme: %s by %s", 
		manifest.Info.Name, manifest.Info.Author)
	return &manifest, nil
}

// ReadOverlayManifest reads and parses an overlay manifest.yml file
func ReadOverlayManifest(path string) (*OverlayManifest, error) {
	logging.LogDebug("Reading overlay manifest from: %s", path)
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest file: %w", err)
	}

	var manifest OverlayManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("error parsing manifest: %w", err)
	}

	logging.LogDebug("Successfully read manifest for overlay: %s by %s", 
		manifest.Info.Name, manifest.Info.Author)
	return &manifest, nil
}

// WriteThemeManifest writes a theme manifest to disk
func WriteThemeManifest(manifest *ThemeManifest, path string) error {
	logging.LogDebug("Writing theme manifest to: %s", path)
	
	// Ensure created date is set
	if manifest.Info.CreatedDate.IsZero() {
		manifest.Info.CreatedDate = time.Now()
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("error marshaling manifest: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing manifest file: %w", err)
	}

	logging.LogDebug("Successfully wrote manifest for theme: %s", manifest.Info.Name)
	return nil
}

// WriteOverlayManifest writes an overlay manifest to disk
func WriteOverlayManifest(manifest *OverlayManifest, path string) error {
	logging.LogDebug("Writing overlay manifest to: %s", path)
	
	// Ensure created date is set
	if manifest.Info.CreatedDate.IsZero() {
		manifest.Info.CreatedDate = time.Now()
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("error marshaling manifest: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing manifest file: %w", err)
	}

	logging.LogDebug("Successfully wrote manifest for overlay: %s", manifest.Info.Name)
	return nil
}

// CreateEmptyThemeManifest creates a new empty theme manifest with basic information
func CreateEmptyThemeManifest(name, author string) *ThemeManifest {
	manifest := &ThemeManifest{}
	manifest.Info.Name = name
	manifest.Info.Author = author
	manifest.Info.Version = "1.0.0"
	manifest.Info.CreatedDate = time.Now()
	
	// Initialize maps
	manifest.Systems = make(map[string]SystemConfig)
	
	return manifest
}

// CreateEmptyOverlayManifest creates a new empty overlay manifest with basic information
func CreateEmptyOverlayManifest(name, author string) *OverlayManifest {
	manifest := &OverlayManifest{}
	manifest.Info.Name = name
	manifest.Info.Author = author
	manifest.Info.Version = "1.0.0"
	manifest.Info.CreatedDate = time.Now()
	
	return manifest
}

// UpdateManifestWithSystemInfo updates a theme manifest with information about a system
func UpdateManifestWithSystemInfo(manifest *ThemeManifest, systemInfo *system.SystemInfo, themeRoot string) {
	if systemInfo.Tag == "" {
		logging.LogDebug("Skipping system with no tag: %s", systemInfo.Name)
		return
	}
	
	// Create or update system entry
	config, exists := manifest.Systems[systemInfo.Tag]
	if !exists {
		config = SystemConfig{
			DisplayName: systemInfo.Name,
			Files: make(map[string]string),
			Paths: make(map[string]string),
		}
	}
	
	// Define relative paths inside theme
	menuBgFile := filepath.Join("Backgrounds/SystemBackgrounds", fmt.Sprintf("%s (%s).png", systemInfo.Name, systemInfo.Tag))
	listBgFile := filepath.Join("Backgrounds/ListBackgrounds", fmt.Sprintf("%s-list (%s).png", systemInfo.Name, systemInfo.Tag))
	menuIconFile := filepath.Join("Icons/SystemIcons", fmt.Sprintf("%s (%s).png", systemInfo.Name, systemInfo.Tag))
	
	// Define system paths
	menuBgPath := filepath.Join(systemInfo.MediaPath, "bg.png")
	listBgPath := filepath.Join(systemInfo.MediaPath, "bglist.png")
	menuIconPath := filepath.Join(system.RomsMediaPath, fmt.Sprintf("%s (%s).png", systemInfo.Name, systemInfo.Tag))
	
	// Update config
	config.Files["menu_bg"] = menuBgFile
	config.Files["list_bg"] = listBgFile
	config.Files["menu_icon"] = menuIconFile
	
	config.Paths["menu_bg_path"] = menuBgPath
	config.Paths["list_bg_path"] = listBgPath
	config.Paths["menu_icon_path"] = menuIconPath
	
	// Update manifest
	manifest.Systems[systemInfo.Tag] = config
	manifest.Content.SystemsCount = len(manifest.Systems)
}

// ValidateThemeManifest checks if a manifest is valid and contains required fields
func ValidateThemeManifest(manifest *ThemeManifest) error {
	if manifest.Info.Name == "" {
		return fmt.Errorf("theme name is required")
	}
	
	if manifest.Info.Author == "" {
		return fmt.Errorf("theme author is required")
	}
	
	if manifest.Info.Version == "" {
		return fmt.Errorf("theme version is required")
	}
	
	return nil
}

// ValidateOverlayManifest checks if an overlay manifest is valid and contains required fields
func ValidateOverlayManifest(manifest *OverlayManifest) error {
	if manifest.Info.Name == "" {
		return fmt.Errorf("overlay name is required")
	}
	
	if manifest.Info.Author == "" {
		return fmt.Errorf("overlay author is required")
	}
	
	if manifest.Info.Version == "" {
		return fmt.Errorf("overlay version is required")
	}
	
	return nil
}