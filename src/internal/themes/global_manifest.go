// src/internal/themes/global_manifest.go
// Provides global manifest functionality to track applied components

package themes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"nextui-themes/internal/logging"
)

// GlobalManifest represents the root-level manifest tracking all applied components
type GlobalManifest struct {
	LastUpdated       time.Time `json:"last_updated"`
	CurrentTheme      string    `json:"current_theme,omitempty"` // Name of the full theme if applied
	AppliedComponents struct {
		Wallpapers string `json:"wallpapers,omitempty"` // Name of applied wallpaper package
		Icons      string `json:"icons,omitempty"`      // Name of applied icon package
		Accents    string `json:"accents,omitempty"`    // Name of applied accent package
		LEDs       string `json:"leds,omitempty"`       // Name of applied LED package
		Fonts      string `json:"fonts,omitempty"`      // Name of applied font package
		Overlays   string `json:"overlays,omitempty"`   // Name of applied overlay package
	} `json:"applied_components"`
	ApplicationInfo struct {
		Version   string `json:"version"`
		BuildDate string `json:"build_date"`
	} `json:"application_info"`
}

// GetGlobalManifestPath returns the path to the global manifest file
func GetGlobalManifestPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	return filepath.Join(cwd, "manifest.json"), nil
}

// LoadGlobalManifest loads the global manifest from disk, or creates a new one if it doesn't exist
func LoadGlobalManifest() (*GlobalManifest, error) {
	manifestPath, err := GetGlobalManifestPath()
	if err != nil {
		return nil, err
	}

	// Check if manifest exists
	_, err = os.Stat(manifestPath)
	if os.IsNotExist(err) {
		// Create a new manifest
		manifest := &GlobalManifest{
			LastUpdated: time.Now(),
			ApplicationInfo: struct {
				Version   string `json:"version"`
				BuildDate string `json:"build_date"`
			}{
				Version:   GetVersionString(),
				BuildDate: time.Now().Format("2006-01-02"),
			},
		}

		// Save the new manifest
		if err := SaveGlobalManifest(manifest); err != nil {
			return nil, fmt.Errorf("error saving new global manifest: %w", err)
		}

		return manifest, nil
	} else if err != nil {
		return nil, fmt.Errorf("error checking global manifest: %w", err)
	}

	// Load existing manifest
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("error reading global manifest: %w", err)
	}

	var manifest GlobalManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("error parsing global manifest: %w", err)
	}

	return &manifest, nil
}

// SaveGlobalManifest saves the global manifest to disk
func SaveGlobalManifest(manifest *GlobalManifest) error {
	manifestPath, err := GetGlobalManifestPath()
	if err != nil {
		return err
	}

	// Update timestamp
	manifest.LastUpdated = time.Now()

	// Convert to JSON
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling global manifest: %w", err)
	}

	// Write to file
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("error writing global manifest: %w", err)
	}

	logging.LogDebug("Saved global manifest to %s", manifestPath)
	return nil
}

// UpdateAppliedComponent updates the global manifest with the newly applied component
func UpdateAppliedComponent(componentType string, componentName string) error {
	manifest, err := LoadGlobalManifest()
	if err != nil {
		return err
	}

	// Update the appropriate component field
	switch componentType {
	case "wallpaper":
		manifest.AppliedComponents.Wallpapers = componentName
	case "icon":
		manifest.AppliedComponents.Icons = componentName
	case "accent":
		manifest.AppliedComponents.Accents = componentName
	case "led":
		manifest.AppliedComponents.LEDs = componentName
	case "font":
		manifest.AppliedComponents.Fonts = componentName
	case "overlay":
		manifest.AppliedComponents.Overlays = componentName
	case "theme":
		manifest.CurrentTheme = componentName
		// Don't clear component fields when applying a full theme
		// They serve as a record of the last specific component packages applied
	default:
		return fmt.Errorf("unknown component type: %s", componentType)
	}

	return SaveGlobalManifest(manifest)
}

// GetAppliedComponent returns the name of the currently applied component of the specified type
func GetAppliedComponent(componentType string) (string, error) {
	manifest, err := LoadGlobalManifest()
	if err != nil {
		return "", err
	}

	switch componentType {
	case "wallpaper":
		return manifest.AppliedComponents.Wallpapers, nil
	case "icon":
		return manifest.AppliedComponents.Icons, nil
	case "accent":
		return manifest.AppliedComponents.Accents, nil
	case "led":
		return manifest.AppliedComponents.LEDs, nil
	case "font":
		return manifest.AppliedComponents.Fonts, nil
	case "overlay":
		return manifest.AppliedComponents.Overlays, nil
	case "theme":
		return manifest.CurrentTheme, nil
	default:
		return "", fmt.Errorf("unknown component type: %s", componentType)
	}
}
