// File: src/internal/themes/manifest.go
// This is a complete replacement for the file

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"strings"
    "sort"
	"gopkg.in/yaml.v3"
	"thememanager/internal/app"
)

// ThemeManifest represents the structure of a theme's manifest.yml file
type ThemeManifest struct {
	Name          string    `yaml:"name"`
	Author        string    `yaml:"author"`
	Description   string    `yaml:"description"`
	RepositoryURL string    `yaml:"repository_url"`
	Commit        string    `yaml:"commit"`
	Branch        string    `yaml:"branch"`
	Device        string    `yaml:"device"`
	Systems       []string  `yaml:"systems"`
	Version       string    `yaml:"version"`
	CreatedDate   time.Time `yaml:"created_date"`
	UpdatedDate   time.Time `yaml:"updated_date"`
	Tags          []string  `yaml:"tags,omitempty"`
}

// AutoDetectSystems scans the system theme directory and returns a sorted list of systems
func AutoDetectSystems() []string {
	app.LogDebug("Auto-detecting systems from theme directory")

	// Check if system theme directory exists
	if !SystemThemeExists() {
		app.LogDebug("System theme directory does not exist, returning empty systems list")
		return []string{}
	}

	// Read the system theme directory
	entries, err := os.ReadDir(SystemThemeDir)
	if err != nil {
		app.LogDebug("Error reading system theme directory: %v", err)
		return []string{}
	}

	var systems []string
	for _, entry := range entries {
		// Only include directories, skip files
		if entry.IsDir() {
			systems = append(systems, entry.Name())
		}
	}

	// Sort alphabetically
	sort.Strings(systems)

	app.LogDebug("Auto-detected systems: %v", systems)
	return systems
}

// CreateBackupManifest creates a manifest specifically for backup themes with template placeholders
func CreateBackupManifest(name string) *ThemeManifest {
	now := time.Now()

	// Auto-detect systems from the theme directory
	detectedSystems := AutoDetectSystems()

	// If no systems detected, provide a reasonable default
	if len(detectedSystems) == 0 {
		detectedSystems = []string{"all"}
	}

	return &ThemeManifest{
		Name:          name,
		Author:        "System",
		Version:       "1.0.0",
		Description:   "Exported system theme",
		CreatedDate:   now,
		UpdatedDate:   now,
		Tags:          []string{},
		RepositoryURL: "https://github.com/[username]/[repo]",
		Commit:        "[commit-hash-will-go-here]",
		Branch:        "main",
		Device:        "brick",
		Systems:       detectedSystems,
	}
}

// ReadManifest reads and parses a theme's manifest file
// Use strictValidation=true when applying themes, false when just displaying
func ReadManifest(themePath string, strictValidation bool) (*ThemeManifest, error) {
	manifestPath := filepath.Join(themePath, ThemeManifestFile)

	app.LogDebug("Reading manifest from %s (strict validation: %v)", manifestPath, strictValidation)

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

	// Always require name at minimum
	if manifest.Name == "" {
		return nil, fmt.Errorf("manifest is missing required 'name' field")
	}

	// Only validate all fields if strict validation is enabled
	if strictValidation {
		if err := ValidateManifest(&manifest); err != nil {
			return nil, err
		}
	}

	app.LogDebug("Successfully read manifest for theme %s by %s", manifest.Name, manifest.Author)
	return &manifest, nil
}

// ValidateManifest performs full validation of all required fields
func ValidateManifest(manifest *ThemeManifest) error {
	var missingFields []string

	if manifest.Name == "" {
		missingFields = append(missingFields, "name")
	}

	if manifest.Author == "" {
		missingFields = append(missingFields, "author")
	}

	if manifest.Description == "" {
		missingFields = append(missingFields, "description")
	}

	if manifest.RepositoryURL == "" {
		missingFields = append(missingFields, "repository_url")
	} else if !strings.HasPrefix(manifest.RepositoryURL, "https://github.com/") {
		return fmt.Errorf("repository_url must be a GitHub URL starting with https://github.com/")
	}

	if manifest.Commit == "" {
		missingFields = append(missingFields, "commit")
	}

	if manifest.Branch == "" {
		missingFields = append(missingFields, "branch")
	}

	if manifest.Device == "" {
		missingFields = append(missingFields, "device")
	}

	if manifest.Systems == nil || len(manifest.Systems) == 0 {
		missingFields = append(missingFields, "systems")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("manifest is missing required fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

// IsManifestValid checks if a manifest is valid without returning detailed errors
func IsManifestValid(manifest *ThemeManifest) bool {
	return ValidateManifest(manifest) == nil
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
		Name:          name,
		Author:        author,
		Version:       "1.0.0",
		Description:   "A theme for the device",
		CreatedDate:   now,
		UpdatedDate:   now,
		Tags:          []string{},
		RepositoryURL: "https://github.com/username/repository",
		Commit:        "main",
		Branch:        "main",
		Device:        "generic",
		Systems:       []string{"all"},
	}
}

// ReadOrCreateManifest reads a manifest from a theme path,
// or creates a default one if it doesn't exist
func ReadOrCreateManifest(themePath, themeName, author string) (*ThemeManifest, error) {
	// Try to read existing manifest (with non-strict validation)
	manifest, err := ReadManifest(themePath, false)
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

	if manifest.Systems != nil && len(manifest.Systems) > 0 {
		info += fmt.Sprintf("\n\nSystems: %s", formatSystems(manifest.Systems))
	}

	if manifest.Device != "" {
		info += fmt.Sprintf("\n\nDevice: %s", manifest.Device)
	}

	if manifest.RepositoryURL != "" {
		info += fmt.Sprintf("\nRepository: %s", manifest.RepositoryURL)
	}

	if manifest.Branch != "" {
		info += fmt.Sprintf("\nBranch: %s", manifest.Branch)
	}

	if len(manifest.Tags) > 0 {
		info += fmt.Sprintf("\n\nTags: %s", formatTags(manifest.Tags))
	}

	return info
}

// Helper function to format systems
func formatSystems(systems []string) string {
	if len(systems) == 0 {
		return ""
	}

	result := systems[0]
	for i := 1; i < len(systems); i++ {
		result += ", " + systems[i]
	}

	return result
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

// GenerateManifestTemplate returns a template string for a manifest
func GenerateManifestTemplate() string {
	return `# Theme Manifest Template
name: My Theme Name               # Required: Display name of the theme
author: Author Name               # Required: Creator of the theme
description: A detailed description of this theme  # Required: What this theme is about
version: 1.0.0                    # Version of this theme

# GitHub repository information
repository_url: https://github.com/username/repository  # Required: GitHub URL
commit: abc123def456              # Required: Commit hash
branch: main                      # Required: Branch name

# Platform information
device: brick                     # Required: Device this theme is for
systems:                          # Required: List of systems this theme covers
  - system1
  - system2
  - system3

# Optional metadata
tags:                            # Optional: Tags for categorization
  - tag1
  - tag2
`
}