// src/internal/accents/accents.go
// Accent color management for NextUI Theme Manager

package accents

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"nextui-themes/internal/logging"
)

// ThemeColor represents a set of colors for a theme
type ThemeColor struct {
	Name   string
	Color1 string // Main UI color
	Color2 string // Primary accent color
	Color3 string // Secondary accent color
	Color4 string // List text color
	Color5 string // Selected list text color
	Color6 string // Hint/info color
}

// Settings file path
const (
	SettingsPath = "/mnt/SDCARD/.userdata/shared/minuisettings.txt"
)

// Predefined color themes
var PredefinedThemes = []ThemeColor{
	{
		Name:   "Classic White",
		Color1: "#FFFFFF", // White
		Color2: "#9B2257", // Pink
		Color3: "#1E2329", // Dark Blue
		Color4: "#FFFFFF", // White
		Color5: "#000000", // Black
		Color6: "#FFFFFF", // White
	},
	{
		Name:   "Midnight Blue",
		Color1: "#000044", // Dark Blue
		Color2: "#3366FF", // Bright Blue
		Color3: "#6699FF", // Light Blue
		Color4: "#FFFFFF", // White
		Color5: "#99CCFF", // Very Light Blue
		Color6: "#B3D9FF", // Almost White Blue
	},
	{
		Name:   "Forest Green",
		Color1: "#004400", // Dark Green
		Color2: "#00AA00", // Medium Green
		Color3: "#33FF33", // Bright Green
		Color4: "#FFFFFF", // White
		Color5: "#80FF80", // Light Green
		Color6: "#B3FFB3", // Very Light Green
	},
	{
		Name:   "Ruby Red",
		Color1: "#440000", // Dark Red
		Color2: "#AA0000", // Medium Red
		Color3: "#FF3333", // Bright Red
		Color4: "#FFFFFF", // White
		Color5: "#FF8080", // Light Red
		Color6: "#FFB3B3", // Very Light Red
	},
	{
		Name:   "Royal Purple",
		Color1: "#330066", // Dark Purple
		Color2: "#6600CC", // Medium Purple
		Color3: "#8833FF", // Bright Purple
		Color4: "#FFFFFF", // White
		Color5: "#BB80FF", // Light Purple
		Color6: "#DDB3FF", // Very Light Purple
	},
	{
		Name:   "Sunset Orange",
		Color1: "#442200", // Dark Orange
		Color2: "#AA5500", // Medium Orange
		Color3: "#FF8833", // Bright Orange
		Color4: "#FFFFFF", // White
		Color5: "#FFBB80", // Light Orange
		Color6: "#FFDDB3", // Very Light Orange
	},
	{
		Name:   "Teal Dream",
		Color1: "#004444", // Dark Teal
		Color2: "#00AAAA", // Medium Teal
		Color3: "#33FFFF", // Bright Teal
		Color4: "#FFFFFF", // White
		Color5: "#80FFFF", // Light Teal
		Color6: "#B3FFFF", // Very Light Teal
	},
	{
		Name:   "Monochrome",
		Color1: "#0A0A0A", // Almost Black
		Color2: "#505050", // Dark Gray
		Color3: "#8C8C8C", // Medium Gray
		Color4: "#DCDC4C", // Light Gray
		Color5: "#FFFFFF", // White
		Color6: "#C8C8C8", // Silver
	},
}

// CurrentTheme holds the currently loaded theme settings
var CurrentTheme ThemeColor

// convertHexFormat converts between display format (#RRGGBB) and storage format (0xRRGGBB)
func convertHexFormat(color string, toStorage bool) string {
	if toStorage {
		// Convert from #RRGGBB to 0xRRGGBB
		if strings.HasPrefix(color, "#") {
			return "0x" + color[1:]
		}
		return color // Already in storage format
	} else {
		// Convert from 0xRRGGBB to #RRGGBB
		if strings.HasPrefix(color, "0x") {
			return "#" + color[2:]
		}
		return color // Already in display format
	}
}

// InitAccentColors loads the current accent colors from disk
func InitAccentColors() error {
	logging.LogDebug("Initializing accent colors")

	// Set a default theme name
	CurrentTheme.Name = "Current"

	// Try to load settings from disk
	theme, err := GetCurrentColors()
	if err != nil {
		logging.LogDebug("Error loading current colors: %v, using defaults", err)
		// Initialize with default values
		CurrentTheme.Color1 = "#FFFFFF" // White
		CurrentTheme.Color2 = "#9B2257" // Pink
		CurrentTheme.Color3 = "#1E2329" // Dark Blue
		CurrentTheme.Color4 = "#FFFFFF" // White
		CurrentTheme.Color5 = "#000000" // Black
		CurrentTheme.Color6 = "#FFFFFF" // White
	} else {
		// Convert from storage format to display format
		CurrentTheme.Color1 = convertHexFormat(theme.Color1, false)
		CurrentTheme.Color2 = convertHexFormat(theme.Color2, false)
		CurrentTheme.Color3 = convertHexFormat(theme.Color3, false)
		CurrentTheme.Color4 = convertHexFormat(theme.Color4, false)
		CurrentTheme.Color5 = convertHexFormat(theme.Color5, false)
		CurrentTheme.Color6 = convertHexFormat(theme.Color6, false)
	}

	logging.LogDebug("Current accent colors initialized: %+v", CurrentTheme)
	return nil
}

// GetCurrentColors reads the current theme colors from the settings file
func GetCurrentColors() (*ThemeColor, error) {
	logging.LogDebug("Reading current accent colors from: %s", SettingsPath)

	// Check if the file exists
	_, err := os.Stat(SettingsPath)
	if os.IsNotExist(err) {
		logging.LogDebug("Settings file does not exist: %s", SettingsPath)
		return nil, fmt.Errorf("settings file not found: %s", SettingsPath)
	}

	// Read the file
	file, err := os.Open(SettingsPath)
	if err != nil {
		logging.LogDebug("Error opening settings file: %v", err)
		return nil, fmt.Errorf("failed to open settings file: %w", err)
	}
	defer file.Close()

	// Parse the file
	colors := &ThemeColor{
		Name: "Current",
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "color1":
			colors.Color1 = value
		case "color2":
			colors.Color2 = value
		case "color3":
			colors.Color3 = value
		case "color4":
			colors.Color4 = value
		case "color5":
			colors.Color5 = value
		case "color6":
			colors.Color6 = value
		}
	}

	if scanner.Err() != nil {
		logging.LogDebug("Error scanning settings file: %v", scanner.Err())
		return nil, fmt.Errorf("error reading settings file: %w", scanner.Err())
	}

	return colors, nil
}

// UpdateCurrentTheme updates the current theme in memory with new values
func UpdateCurrentTheme(themeName string) error {
	logging.LogDebug("Updating current theme to: %s", themeName)

	// Find the selected theme
	var found bool
	for _, theme := range PredefinedThemes {
		if theme.Name == themeName {
			CurrentTheme.Name = theme.Name
			CurrentTheme.Color1 = theme.Color1
			CurrentTheme.Color2 = theme.Color2
			CurrentTheme.Color3 = theme.Color3
			CurrentTheme.Color4 = theme.Color4
			CurrentTheme.Color5 = theme.Color5
			CurrentTheme.Color6 = theme.Color6
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("theme not found: %s", themeName)
	}

	logging.LogDebug("Theme updated in memory: %+v", CurrentTheme)
	return nil
}

// ApplyThemeColors applies the specified theme colors to the system
func ApplyThemeColors(theme *ThemeColor) error {
	logging.LogDebug("Applying theme colors: %s", theme.Name)

	// Read the current settings file
	_, err := os.Stat(SettingsPath)
	if os.IsNotExist(err) {
		logging.LogDebug("Settings file does not exist, creating: %s", SettingsPath)

		// Create parent directories if needed
		err = os.MkdirAll(filepath.Dir(SettingsPath), 0755)
		if err != nil {
			logging.LogDebug("Error creating parent directories: %v", err)
			return fmt.Errorf("failed to create settings directory: %w", err)
		}

		// Create an empty file
		file, err := os.Create(SettingsPath)
		if err != nil {
			logging.LogDebug("Error creating settings file: %v", err)
			return fmt.Errorf("failed to create settings file: %w", err)
		}
		file.Close()
	}

	// Read existing settings
	settings := make(map[string]string)

	file, err := os.Open(SettingsPath)
	if err != nil {
		logging.LogDebug("Error opening settings file: %v", err)
		return fmt.Errorf("failed to open settings file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		settings[key] = value
	}
	file.Close()

	// Update with new theme colors - convert from display format to storage format
	settings["color1"] = convertHexFormat(theme.Color1, true)
	settings["color2"] = convertHexFormat(theme.Color2, true)
	settings["color3"] = convertHexFormat(theme.Color3, true)
	settings["color4"] = convertHexFormat(theme.Color4, true)
	settings["color5"] = convertHexFormat(theme.Color5, true)
	settings["color6"] = convertHexFormat(theme.Color6, true)

	// Write back to file
	tempFile := SettingsPath + ".tmp"
	outFile, err := os.Create(tempFile)
	if err != nil {
		logging.LogDebug("Error creating temp settings file: %v", err)
		return fmt.Errorf("failed to create temp settings file: %w", err)
	}

	// Write each line back to the file
	for key, value := range settings {
		_, err := fmt.Fprintf(outFile, "%s=%s\n", key, value)
		if err != nil {
			outFile.Close()
			os.Remove(tempFile)
			logging.LogDebug("Error writing to settings file: %v", err)
			return fmt.Errorf("failed to write settings: %w", err)
		}
	}

	outFile.Close()

	// Replace the original file with the new one
	err = os.Rename(tempFile, SettingsPath)
	if err != nil {
		logging.LogDebug("Error replacing settings file: %v", err)
		return fmt.Errorf("failed to update settings file: %w", err)
	}

	logging.LogDebug("Successfully applied theme colors")
	return nil
}

// GetColorPreviewText formats a color value for display
func GetColorPreviewText(colorName string, colorValue string) string {
	return fmt.Sprintf("%s: %s", colorName, colorValue)
}

// ApplyCurrentTheme applies the current in-memory theme to the system
func ApplyCurrentTheme() error {
	logging.LogDebug("Applying current in-memory theme")
	return ApplyThemeColors(&CurrentTheme)
}