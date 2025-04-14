package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// PathMapping represents a mapping between theme and system paths
type PathMapping struct {
	ThemePath  string            `json:"theme_path"`
	SystemPath string            `json:"system_path"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

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
	AccentColors map[string]string            `json:"accent_colors"`
	LEDSettings  map[string]map[string]string `json:"led_settings"`
}

// Create a manifest based on our test_theme
func main() {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Test theme path
	themePath := filepath.Join(cwd, "test_theme")

	fmt.Printf("Creating manifest for test theme at: %s\n", themePath)

	// Create a new manifest
	manifest := createEmptyManifest()

	// Update with theme info
	manifest.ThemeInfo.Name = "Test Theme"
	manifest.ThemeInfo.Author = "Test Author"

	// Set wallpapers
	manifest.Content.Wallpapers.Present = true
	manifest.Content.Wallpapers.Count = 3
	manifest.PathMappings.Wallpapers = []PathMapping{
		{
			ThemePath:  "Wallpapers/SystemWallpapers/Root.png",
			SystemPath: "/mnt/SDCARD/bg.png",
			Metadata: map[string]string{
				"SystemName":    "Root",
				"WallpaperType": "Main",
			},
		},
		{
			ThemePath:  "Wallpapers/SystemWallpapers/Recently Played.png",
			SystemPath: "/mnt/SDCARD/Recently Played/.media/bg.png",
			Metadata: map[string]string{
				"SystemName":    "Recently Played",
				"WallpaperType": "Media",
			},
		},
		{
			ThemePath:  "Wallpapers/CollectionWallpapers/Handhelds.png",
			SystemPath: "/mnt/SDCARD/Collections/Handhelds/.media/bg.png",
			Metadata: map[string]string{
				"CollectionName": "Handhelds",
				"WallpaperType":  "Collection",
			},
		},
	}

	// Set icons
	manifest.Content.Icons.Present = true
	manifest.Content.Icons.SystemCount = 1
	manifest.Content.Icons.ToolCount = 1
	manifest.Content.Icons.CollectionCount = 1
	manifest.PathMappings.Icons = []PathMapping{
		{
			ThemePath:  "Icons/SystemIcons/Collections.png",
			SystemPath: "/mnt/SDCARD/.media/Collections.png",
			Metadata: map[string]string{
				"SystemName": "Collections",
				"SystemTag":  "COLLECTIONS",
				"IconType":   "Special",
			},
		},
		{
			ThemePath:  "Icons/ToolIcons/Tetris.png",
			SystemPath: "/mnt/SDCARD/Tools/.media/Tetris.png",
			Metadata: map[string]string{
				"ToolName": "Tetris",
				"IconType": "Tool",
			},
		},
		{
			ThemePath:  "Icons/CollectionIcons/Favorites.png",
			SystemPath: "/mnt/SDCARD/Collections/.media/Favorites.png",
			Metadata: map[string]string{
				"CollectionName": "Favorites",
				"IconType":       "Collection",
			},
		},
	}

	// Write the manifest to file
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		log.Fatalf("Error creating manifest JSON: %v", err)
	}

	manifestPath := filepath.Join(themePath, "manifest.json")
	err = os.WriteFile(manifestPath, manifestJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing manifest file: %v", err)
	}

	fmt.Printf("Manifest created successfully at: %s\n", manifestPath)
	fmt.Println("Try importing this theme with the theme manager.")
}

// createEmptyManifest creates a new empty manifest with basic metadata
func createEmptyManifest() ThemeManifest {
	var manifest ThemeManifest

	// Initialize metadata sections
	manifest.ThemeInfo.Version = "1.0.0"
	manifest.ThemeInfo.ExportedBy = "Theme Manager v1.0.0"
	manifest.ThemeInfo.CreationDate = time.Now()
	manifest.ThemeInfo.Author = "AuthorName"

	// Initialize content section
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

	// Initialize accent colors
	manifest.AccentColors = make(map[string]string)

	// Initialize LED settings
	manifest.LEDSettings = make(map[string]map[string]string)

	return manifest
}
