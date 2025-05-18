// internal/themes/manifest.go
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
	"thememanager/internal/app"
)

// ThemeManifest represents the structure of a theme's manifest.yml file
type ThemeManifest struct {
	Name        string    `yaml:"name"`
	Author      string    `yaml:"author"`
	Version     string    `yaml:"version"`
	Description string    `yaml:"description"`
	CreatedDate time.Time `yaml:"created_date"`
	UpdatedDate time.Time `yaml:"updated_date"`
	Tags        []string  `yaml:"tags,omitempty"`
}

// ReadManifest reads and parses a theme's manifest file
func ReadManifest(themePath string) (*ThemeManifest, error) {
	manifestPath := filepath.Join(themePath, ThemeManifestFile)

	app.LogDebug("Reading manifest from %s", manifestPath)

	// Check if manifest file exists
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("manifest file does not exist: %w", err)
	}

	// Read the manifest file
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	// Parse YAML
	var manifest ThemeManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Validate manifest
	if manifest.Name == "" {
		return nil, fmt.Errorf("manifest is missing required 'name' field")
	}

	app.LogDebug("Successfully read manifest for theme %s by %s", manifest.Name, manifest.Author)
	return &manifest, nil
}

// WriteManifest writes a theme manifest to file
func WriteManifest(manifest *ThemeManifest, themePath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(themePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	manifestPath := filepath.Join(themePath, ThemeManifestFile)

	app.LogDebug("Writing manifest to %s", manifestPath)

	// Update timestamp
	manifest.UpdatedDate = time.Now()

	// Convert to YAML
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to convert manifest to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	app.LogDebug("Successfully wrote manifest for theme %s", manifest.Name)
	return nil
}

// CreateDefaultManifest creates a new manifest with default values
func CreateDefaultManifest(name, author string) *ThemeManifest {
	now := time.Now()

	return &ThemeManifest{
		Name:        name,
		Author:      author,
		Version:     "1.0.0",
		Description: "A theme for the device",
		CreatedDate: now,
		UpdatedDate: now,
		Tags:        []string{},
	}
}

// ReadOrCreateManifest reads a manifest from a theme path,
// or creates a default one if it doesn't exist
func ReadOrCreateManifest(themePath, themeName, author string) (*ThemeManifest, error) {
	// Try to read existing manifest
	manifest, err := ReadManifest(themePath)
	if err == nil {
		return manifest, nil
	}

	// Create default manifest if file doesn't exist
	if os.IsNotExist(err) {
		app.LogDebug("Creating default manifest for %s", themeName)
		manifest = CreateDefaultManifest(themeName, author)

		// Write the new manifest
		if err := WriteManifest(manifest, themePath); err != nil {
			return nil, fmt.Errorf("failed to write new manifest: %w", err)
		}

		return manifest, nil
	}

	// Return original error if not a "not exists" error
	return nil, err
}

// GetThemeInfo returns a formatted string with theme information
func GetThemeInfo(manifest *ThemeManifest) string {
	info := fmt.Sprintf("%s v%s\nBy: %s", manifest.Name, manifest.Version, manifest.Author)

	if manifest.Description != "" {
		info += fmt.Sprintf("\n\n%s", manifest.Description)
	}

	if len(manifest.Tags) > 0 {
		info += fmt.Sprintf("\n\nTags: %s", formatTags(manifest.Tags))
	}

	return info
}

// Helper function to format tags
func formatTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	result := tags[0]
	for i := 1; i < len(tags); i++ {
		result += ", " + tags[i]
	}

	return result
}