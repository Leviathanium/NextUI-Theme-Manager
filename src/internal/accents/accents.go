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
		Color1: "0xFFFFFF", // White
		Color2: "0x9B2257", // Pink
		Color3: "0x1E2329", // Dark Blue
		Color4: "0xFFFFFF", // White
		Color5: "0x000000", // Black
		Color6: "0xFFFFFF", // White
	},
	{
		Name:   "Midnight Blue",
		Color1: "0x000044", // Dark Blue
		Color2: "0x3366FF", // Bright Blue
		Color3: "0x6699FF", // Light Blue
		Color4: "0xFFFFFF", // White
		Color5: "0x99CCFF", // Very Light Blue
		Color6: "0xB3D9FF", // Almost White Blue
	},
	{
		Name:   "Forest Green",
		Color1: "0x004400", // Dark Green
		Color2: "0x00AA00", // Medium Green
		Color3: "0x33FF33", // Bright Green
		Color4: "0xFFFFFF", // White
		Color5: "0x80FF80", // Light Green
		Color6: "0xB3FFB3", // Very Light Green
	},
	{
		Name:   "Ruby Red",
		Color1: "0x440000", // Dark Red
		Color2: "0xAA0000", // Medium Red
		Color3: "0xFF3333", // Bright Red
		Color4: "0xFFFFFF", // White
		Color5: "0xFF8080", // Light Red
		Color6: "0xFFB3B3", // Very Light Red
	},
	{
		Name:   "Royal Purple",
		Color1: "0x330066", // Dark Purple
		Color2: "0x6600CC", // Medium Purple
		Color3: "0x8833FF", // Bright Purple
		Color4: "0xFFFFFF", // White
		Color5: "0xBB80FF", // Light Purple
		Color6: "0xDDB3FF", // Very Light Purple
	},
	{
		Name:   "Sunset Orange",
		Color1: "0x442200", // Dark Orange
		Color2: "0xAA5500", // Medium Orange
		Color3: "0xFF8833", // Bright Orange
		Color4: "0xFFFFFF", // White
		Color5: "0xFFBB80", // Light Orange
		Color6: "0xFFDDB3", // Very Light Orange
	},
	{
		Name:   "Teal Dream",
		Color1: "0x004444", // Dark Teal
		Color2: "0x00AAAA", // Medium Teal
		Color3: "0x33FFFF", // Bright Teal
		Color4: "0xFFFFFF", // White
		Color5: "0x80FFFF", // Light Teal
		Color6: "0xB3FFFF", // Very Light Teal
	},
	{
		Name:   "Monochrome",
		Color1: "0x0A0A0A", // Almost Black
		Color2: "0x505050", // Dark Gray
		Color3: "0x8C8C8C", // Medium Gray
		Color4: "0xDCDCDC", // Light Gray
		Color5: "0xFFFFFF", // White
		Color6: "0xC8C8C8", // Silver
	},
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

	// Update with new theme colors
	settings["color1"] = theme.Color1
	settings["color2"] = theme.Color2
	settings["color3"] = theme.Color3
	settings["color4"] = theme.Color4
	settings["color5"] = theme.Color5
	settings["color6"] = theme.Color6

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