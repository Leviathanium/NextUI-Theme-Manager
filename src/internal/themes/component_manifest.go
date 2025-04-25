// src/internal/themes/component_manifest.go
// Defines manifest structures for individual theme components

package themes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"nextui-themes/internal/logging"
)

// ComponentType constants
const (
	ComponentWallpaper = "wallpaper"
	ComponentIcon     = "icon"
	ComponentAccent   = "accent"
	ComponentLED      = "led"
	ComponentFont     = "font"
	ComponentOverlay  = "overlay"
)

// ComponentExtension maps component types to their file extensions
var ComponentExtension = map[string]string{
	ComponentWallpaper: ".bg",
	ComponentIcon:     ".icon",
	ComponentAccent:   ".acc",
	ComponentLED:      ".led",
	ComponentFont:     ".font",
	ComponentOverlay:  ".over",
}

// ComponentInfo holds common metadata for all component types
type ComponentInfo struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Version      string    `json:"version"`
	Author       string    `json:"author"`
	CreationDate time.Time `json:"creation_date"`
	ExportedBy   string    `json:"exported_by"`
}

// BaseComponentManifest contains the shared structure for all component manifests
type BaseComponentManifest struct {
	ComponentInfo ComponentInfo `json:"component_info"`
}

// WallpaperManifest for .bg component packages
type WallpaperManifest struct {
	ComponentInfo ComponentInfo `json:"component_info"`
	Content       struct {
		Count                int      `json:"count"`
		SystemWallpapers     []string `json:"system_wallpapers"`
		ListWallpapers       []string `json:"list_wallpapers"`      // New field for list wallpapers
		CollectionWallpapers []string `json:"collection_wallpapers"`
	} `json:"content"`
	PathMappings []PathMapping `json:"path_mappings"`
}

// IconManifest for .icon component packages
type IconManifest struct {
	ComponentInfo ComponentInfo `json:"component_info"`
	Content       struct {
		SystemCount     int      `json:"system_count"`
		ToolCount       int      `json:"tool_count"`
		CollectionCount int      `json:"collection_count"`
		SystemIcons     []string `json:"system_icons"`
		ToolIcons       []string `json:"tool_icons"`
		CollectionIcons []string `json:"collection_icons"`
	} `json:"content"`
	PathMappings []PathMapping `json:"path_mappings"`
}

// AccentManifest for .acc component packages
type AccentManifest struct {
	ComponentInfo ComponentInfo `json:"component_info"`
	AccentColors  struct {
		Color1 string `json:"color1"`
		Color2 string `json:"color2"`
		Color3 string `json:"color3"`
		Color4 string `json:"color4"`
		Color5 string `json:"color5"`
		Color6 string `json:"color6"`
	} `json:"accent_colors"`
}

// LEDManifest for .led component packages
type LEDManifest struct {
	ComponentInfo ComponentInfo `json:"component_info"`
	LEDSettings   struct {
		F1Key      LEDSetting `json:"f1_key"`
		F2Key      LEDSetting `json:"f2_key"`
		TopBar     LEDSetting `json:"top_bar"`
		LRTriggers LEDSetting `json:"lr_triggers"`
	} `json:"led_settings"`
}

// FontManifest for .font component packages
type FontManifest struct {
	ComponentInfo ComponentInfo `json:"component_info"`
	Content       struct {
		OGReplaced   bool `json:"og_replaced"`
		NextReplaced bool `json:"next_replaced"`
	} `json:"content"`
	PathMappings map[string]PathMapping `json:"path_mappings"`
}

// OverlayManifest for .over component packages
type OverlayManifest struct {
	ComponentInfo ComponentInfo `json:"component_info"`
	Content       struct {
		Systems []string `json:"systems"`
	} `json:"content"`
	PathMappings []PathMapping `json:"path_mappings"`
}

// CreateMinimalComponentManifest creates a minimal component manifest with just essential information
func CreateMinimalComponentManifest(componentType string, name string, author string) (interface{}, error) {
    // Create basic component info
    info := ComponentInfo{
        Name:         name,
        Type:         componentType,
        Version:      "1.0.0",
        Author:       author,  // Preserve author information
        CreationDate: time.Now(),
        ExportedBy:   GetVersionString(),
    }

    // Create appropriate struct based on component type
    switch componentType {
    // Part of CreateMinimalComponentManifest function - just the relevant case:
    case ComponentWallpaper:
        var manifest WallpaperManifest
        manifest.ComponentInfo = info
        // Initialize content section with empty values
        manifest.Content.Count = 0
        manifest.Content.SystemWallpapers = []string{}
        manifest.Content.ListWallpapers = []string{}     // Initialize the new field
        manifest.Content.CollectionWallpapers = []string{}
        // Leave path_mappings empty
        manifest.PathMappings = []PathMapping{}
        return &manifest, nil

    case ComponentIcon:
        var manifest IconManifest
        manifest.ComponentInfo = info
        // Initialize content section with empty values
        manifest.Content.SystemCount = 0
        manifest.Content.ToolCount = 0
        manifest.Content.CollectionCount = 0
        manifest.Content.SystemIcons = []string{}
        manifest.Content.ToolIcons = []string{}
        manifest.Content.CollectionIcons = []string{}
        // Leave path_mappings empty
        manifest.PathMappings = []PathMapping{}
        return &manifest, nil

    case ComponentAccent:
        var manifest AccentManifest
        manifest.ComponentInfo = info
        // Initialize with default values
        manifest.AccentColors.Color1 = "0xFFFFFF"
        manifest.AccentColors.Color2 = "0x9B2257"
        manifest.AccentColors.Color3 = "0x1E2329"
        manifest.AccentColors.Color4 = "0xFFFFFF"
        manifest.AccentColors.Color5 = "0x000000"
        manifest.AccentColors.Color6 = "0xFFFFFF"
        return &manifest, nil

    case ComponentLED:
        var manifest LEDManifest
        manifest.ComponentInfo = info
        // Initialize with default values
        initLEDSetting := func(led *LEDSetting) {
            led.Effect = 1
            led.Color1 = "0xFFFFFF"
            led.Color2 = "0x000000"
            led.Speed = 1000
            led.Brightness = 100
            led.Trigger = 1
            led.InBrightness = 100
        }
        initLEDSetting(&manifest.LEDSettings.F1Key)
        initLEDSetting(&manifest.LEDSettings.F2Key)
        initLEDSetting(&manifest.LEDSettings.TopBar)
        initLEDSetting(&manifest.LEDSettings.LRTriggers)
        return &manifest, nil

    case ComponentFont:
        var manifest FontManifest
        manifest.ComponentInfo = info
        // Initialize content flags
        manifest.Content.OGReplaced = false
        manifest.Content.NextReplaced = false
        // Initialize path_mappings
        manifest.PathMappings = make(map[string]PathMapping)
        return &manifest, nil

    case ComponentOverlay:
        var manifest OverlayManifest
        manifest.ComponentInfo = info
        // Initialize content
        manifest.Content.Systems = []string{}
        // Leave path_mappings empty
        manifest.PathMappings = []PathMapping{}
        return &manifest, nil

    default:
        return nil, fmt.Errorf("unknown component type: %s", componentType)
    }
}

// CreateMinimalThemeManifest creates a minimal theme manifest with just essential information
func CreateMinimalThemeManifest(themeName string, author string) *ThemeManifest {
    manifest := &ThemeManifest{}

    // Set basic theme info
    manifest.ThemeInfo.Name = themeName
    manifest.ThemeInfo.Version = "1.0.0"
    manifest.ThemeInfo.Author = author
    manifest.ThemeInfo.CreationDate = time.Now()
    manifest.ThemeInfo.ExportedBy = GetVersionString()

    // Initialize content section with default values
    manifest.Content.Wallpapers.Present = false
    manifest.Content.Wallpapers.Count = 0

    manifest.Content.Icons.Present = false
    manifest.Content.Icons.SystemCount = 0
    manifest.Content.Icons.ToolCount = 0
    manifest.Content.Icons.CollectionCount = 0

    manifest.Content.Overlays.Present = false
    manifest.Content.Overlays.Systems = []string{}

    manifest.Content.Fonts.Present = false
    manifest.Content.Fonts.OGReplaced = false
    manifest.Content.Fonts.NextReplaced = false

    manifest.Content.Settings.AccentsIncluded = false
    manifest.Content.Settings.LEDsIncluded = false

    // Initialize path mappings
    manifest.PathMappings.Wallpapers = []PathMapping{}
    manifest.PathMappings.Icons = []PathMapping{}
    manifest.PathMappings.Overlays = []PathMapping{}
    manifest.PathMappings.Fonts = make(map[string]PathMapping)
    manifest.PathMappings.Settings = make(map[string]PathMapping)

    // Initialize default accent colors
    manifest.AccentColors.Color1 = "0xFFFFFF"
    manifest.AccentColors.Color2 = "0x9B2257"
    manifest.AccentColors.Color3 = "0x1E2329"
    manifest.AccentColors.Color4 = "0xFFFFFF"
    manifest.AccentColors.Color5 = "0x000000"
    manifest.AccentColors.Color6 = "0xFFFFFF"

    // Initialize default LED settings
    initLEDSetting := func(led *LEDSetting) {
        led.Effect = 1
        led.Color1 = "0xFFFFFF"
        led.Color2 = "0x000000"
        led.Speed = 1000
        led.Brightness = 100
        led.Trigger = 1
        led.InBrightness = 100
    }
    initLEDSetting(&manifest.LEDSettings.F1Key)
    initLEDSetting(&manifest.LEDSettings.F2Key)
    initLEDSetting(&manifest.LEDSettings.TopBar)
    initLEDSetting(&manifest.LEDSettings.LRTriggers)

    return manifest
}

// CreateComponentManifest creates a new component manifest of the specified type
func CreateComponentManifest(componentType string, name string) (interface{}, error) {
    // Create basic component info
    info := ComponentInfo{
        Name:         name,
        Type:         componentType,
        Version:      "1.0.0",
        CreationDate: time.Now(),
        ExportedBy:   GetVersionString(),
    }

    // Set author to blank - we'll preserve existing author when updating
    // Instead of hardcoding "User" here
    info.Author = ""

    // Create manifest based on component type
    switch componentType {
    case ComponentWallpaper:
        var manifest WallpaperManifest
        manifest.ComponentInfo = info
        manifest.Content.Count = 0
        manifest.Content.SystemWallpapers = []string{}
        manifest.Content.CollectionWallpapers = []string{}
        manifest.PathMappings = []PathMapping{}
        return &manifest, nil

    case ComponentIcon:
        var manifest IconManifest
        manifest.ComponentInfo = info
        manifest.Content.SystemCount = 0
        manifest.Content.ToolCount = 0
        manifest.Content.CollectionCount = 0
        manifest.Content.SystemIcons = []string{}
        manifest.Content.ToolIcons = []string{}
        manifest.Content.CollectionIcons = []string{}
        manifest.PathMappings = []PathMapping{}
        return &manifest, nil

    case ComponentAccent:
        var manifest AccentManifest
        manifest.ComponentInfo = info
        return &manifest, nil

    case ComponentLED:
        var manifest LEDManifest
        manifest.ComponentInfo = info
        return &manifest, nil

    case ComponentFont:
        var manifest FontManifest
        manifest.ComponentInfo = info
        manifest.Content.OGReplaced = false
        manifest.Content.NextReplaced = false
        manifest.PathMappings = make(map[string]PathMapping)
        return &manifest, nil

    case ComponentOverlay:
        var manifest OverlayManifest
        manifest.ComponentInfo = info
        manifest.Content.Systems = []string{}
        manifest.PathMappings = []PathMapping{}
        return &manifest, nil

    default:
        return nil, fmt.Errorf("unknown component type: %s", componentType)
    }
}

// WriteComponentManifest writes a component manifest to the specified directory
func WriteComponentManifest(componentPath string, manifest interface{}) error {
	manifestPath := filepath.Join(componentPath, "manifest.json")

	// Convert to JSON
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling component manifest: %w", err)
	}

	// Write to file
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("error writing component manifest: %w", err)
	}

	logging.LogDebug("Saved component manifest to %s", manifestPath)
	return nil
}

// LoadComponentManifest loads a component manifest from the specified directory
func LoadComponentManifest(componentPath string) (interface{}, error) {
	manifestPath := filepath.Join(componentPath, "manifest.json")

	// Check if manifest exists
	_, err := os.Stat(manifestPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("component manifest not found: %s", manifestPath)
	} else if err != nil {
		return nil, fmt.Errorf("error checking component manifest: %w", err)
	}

	// Read manifest file
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("error reading component manifest: %w", err)
	}

	// First, unmarshal as BaseComponentManifest to determine type
	var baseManifest BaseComponentManifest
	if err := json.Unmarshal(data, &baseManifest); err != nil {
		return nil, fmt.Errorf("error parsing component manifest: %w", err)
	}

	// Based on the component type, unmarshal into the appropriate struct
	switch baseManifest.ComponentInfo.Type {
	case ComponentWallpaper:
		var manifest WallpaperManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("error parsing wallpaper manifest: %w", err)
		}
		return &manifest, nil

	case ComponentIcon:
		var manifest IconManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("error parsing icon manifest: %w", err)
		}
		return &manifest, nil

	case ComponentAccent:
		var manifest AccentManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("error parsing accent manifest: %w", err)
		}
		return &manifest, nil

	case ComponentLED:
		var manifest LEDManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("error parsing LED manifest: %w", err)
		}
		return &manifest, nil

	case ComponentFont:
		var manifest FontManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("error parsing font manifest: %w", err)
		}
		return &manifest, nil

	case ComponentOverlay:
		var manifest OverlayManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("error parsing overlay manifest: %w", err)
		}
		return &manifest, nil

	default:
		return nil, fmt.Errorf("unknown component type: %s", baseManifest.ComponentInfo.Type)
	}
}