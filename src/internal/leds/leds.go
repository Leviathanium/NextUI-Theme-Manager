// src/internal/leds/leds.go
// LED settings management for NextUI Theme Manager

package leds

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nextui-themes/internal/logging"
)

// LEDEffect represents an LED lighting effect
type LEDEffect int

const (
	Static LEDEffect = iota + 1
	Breathe
	IntervalBreathe
	StaticColor
	Blink1
	Blink2
	Blink3
)

// LEDTheme represents a theme for LED lighting
type LEDTheme struct {
	Name   string
	Color  string
	Effect LEDEffect
}

// LightSettings represents settings for one LED light
type LightSettings struct {
	Name        string
	Effect      int
	Color1      string
	Color2      string
	Speed       int
	Brightness  int
	Trigger     int
	Filename    string
	InBrightness int
}

// Settings file paths
const (
	BrickSettingsPath = "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"
	NormalSettingsPath = "/mnt/SDCARD/.userdata/shared/ledsettings.txt"
)

// Current LED settings in memory
var (
	CurrentLEDSettings []LightSettings
	CurrentLEDTheme string
)

// Predefined LED themes
var PredefinedThemes = []LEDTheme{
	{
		Name:   "White Static",
		Color:  "#FFFFFF",
		Effect: Static,
	},
	{
		Name:   "Blue Static",
		Color:  "#3366FF",
		Effect: Static,
	},
	{
		Name:   "Green Static",
		Color:  "#00AA00",
		Effect: Static,
	},
	{
		Name:   "Red Static",
		Color:  "#AA0000",
		Effect: Static,
	},
	{
		Name:   "Purple Static",
		Color:  "#8833FF",
		Effect: Static,
	},
	{
		Name:   "Orange Static",
		Color:  "#FF8833",
		Effect: Static,
	},
	{
		Name:   "White Breathing",
		Color:  "#FFFFFF",
		Effect: Breathe,
	},
	{
		Name:   "Blue Breathing",
		Color:  "#3366FF",
		Effect: Breathe,
	},
	{
		Name:   "Green Breathing",
		Color:  "#00AA00",
		Effect: Breathe,
	},
	{
		Name:   "Red Breathing",
		Color:  "#AA0000",
		Effect: Breathe,
	},
	{
		Name:   "Purple Breathing",
		Color:  "#8833FF",
		Effect: Breathe,
	},
	{
		Name:   "Orange Breathing",
		Color:  "#FF8833",
		Effect: Breathe,
	},
}

// InitLEDSettings initializes the LED settings from disk
func InitLEDSettings() error {
	logging.LogDebug("Initializing LED settings")

	// Initialize with default values in case we can't read the file
	CurrentLEDTheme = "White Static"

	// Try to load settings from disk
	settings, err := GetCurrentLEDSettings()
	if err != nil {
		logging.LogDebug("Error loading current LED settings: %v, using defaults", err)

		// Create default light settings for a brick device
		CurrentLEDSettings = []LightSettings{
			{
				Name:        "F1 key",
				Effect:      int(Static),
				Color1:      "#FFFFFF",
				Color2:      "#000000",
				Speed:       1000,
				Brightness:  100,
				Trigger:     1,
				Filename:    "",
				InBrightness: 100,
			},
			{
				Name:        "F2 key",
				Effect:      int(Static),
				Color1:      "#FFFFFF",
				Color2:      "#000000",
				Speed:       1000,
				Brightness:  100,
				Trigger:     1,
				Filename:    "",
				InBrightness: 100,
			},
			{
				Name:        "Top bar",
				Effect:      int(Static),
				Color1:      "#FFFFFF",
				Color2:      "#000000",
				Speed:       1000,
				Brightness:  100,
				Trigger:     1,
				Filename:    "",
				InBrightness: 100,
			},
			{
				Name:        "L&R triggers",
				Effect:      int(Static),
				Color1:      "#FFFFFF",
				Color2:      "#000000",
				Speed:       1000,
				Brightness:  100,
				Trigger:     1,
				Filename:    "",
				InBrightness: 100,
			},
		}
	} else {
		CurrentLEDSettings = settings
		DetermineCurrentTheme()
	}

	logging.LogDebug("LED settings initialized with theme: %s", CurrentLEDTheme)
	return nil
}

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

// IsBrickDevice checks if the current device is a TrimUI Brick
func IsBrickDevice() bool {
	device := os.Getenv("DEVICE")
	return device == "brick"
}

// GetSettingsPath returns the appropriate settings path based on device
func GetSettingsPath() string {
	if IsBrickDevice() {
		return BrickSettingsPath
	}
	return NormalSettingsPath
}

// GetLightCount returns the number of light locations for the current device
func GetLightCount() int {
	if IsBrickDevice() {
		return 4 // F1 key, F2 key, Top bar, L&R triggers
	}
	return 2 // Joysticks, Logo
}

// GetCurrentLEDSettings reads the current LED settings from disk
func GetCurrentLEDSettings() ([]LightSettings, error) {
	settingsPath := GetSettingsPath()
	lightCount := GetLightCount()

	logging.LogDebug("Reading LED settings from: %s", settingsPath)

	// Check if the file exists
	_, err := os.Stat(settingsPath)
	if os.IsNotExist(err) {
		logging.LogDebug("LED settings file does not exist: %s", settingsPath)
		return nil, fmt.Errorf("LED settings file not found: %s", settingsPath)
	}

	// Read the file
	file, err := os.Open(settingsPath)
	if err != nil {
		logging.LogDebug("Error opening LED settings file: %v", err)
		return nil, fmt.Errorf("failed to open LED settings file: %w", err)
	}
	defer file.Close()

	// Parse the file
	var settings []LightSettings
	var currentLight *LightSettings

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Extract light name
			lightName := line[1 : len(line)-1]

			// Add the previous light settings if they exist
			if currentLight != nil {
				settings = append(settings, *currentLight)
			}

			// Start a new light settings
			currentLight = &LightSettings{
				Name: lightName,
			}
			continue
		}

		// Skip if we're not in a section yet
		if currentLight == nil {
			continue
		}

		// Parse key-value pair
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set the appropriate field
		switch key {
		case "effect":
			effect := 0
			fmt.Sscanf(value, "%d", &effect)
			currentLight.Effect = effect
		case "color1":
			currentLight.Color1 = convertHexFormat(value, false) // Convert to display format
		case "color2":
			currentLight.Color2 = convertHexFormat(value, false) // Convert to display format
		case "speed":
			speed := 0
			fmt.Sscanf(value, "%d", &speed)
			currentLight.Speed = speed
		case "brightness":
			brightness := 0
			fmt.Sscanf(value, "%d", &brightness)
			currentLight.Brightness = brightness
		case "trigger":
			trigger := 0
			fmt.Sscanf(value, "%d", &trigger)
			currentLight.Trigger = trigger
		case "filename":
			currentLight.Filename = value
		case "inbrightness":
			inbrightness := 0
			fmt.Sscanf(value, "%d", &inbrightness)
			currentLight.InBrightness = inbrightness
		}
	}

	// Add the last light if there is one
	if currentLight != nil {
		settings = append(settings, *currentLight)
	}

	// Check if we got all expected lights
	if len(settings) != lightCount {
		logging.LogDebug("Warning: Expected %d lights but found %d", lightCount, len(settings))
	}

	return settings, nil
}

// DetermineCurrentTheme tries to determine which predefined theme is closest to the current settings
func DetermineCurrentTheme() {
	// Default
	CurrentLEDTheme = "Custom"

	if len(CurrentLEDSettings) == 0 {
		return
	}

	// Use the first light as reference
	firstLight := CurrentLEDSettings[0]

	// Check if all lights have the same effect and color
	allSame := true
	for i := 1; i < len(CurrentLEDSettings); i++ {
		if CurrentLEDSettings[i].Effect != firstLight.Effect ||
		   CurrentLEDSettings[i].Color1 != firstLight.Color1 {
			allSame = false
			break
		}
	}

	if !allSame {
		return
	}

	// Find matching predefined theme
	for _, theme := range PredefinedThemes {
		if int(theme.Effect) == firstLight.Effect &&
		   convertHexFormat(theme.Color, true) == convertHexFormat(firstLight.Color1, true) {
			CurrentLEDTheme = theme.Name
			return
		}
	}
}

// UpdateCurrentLEDTheme updates the current LED theme in memory with new values
func UpdateCurrentLEDTheme(themeName string) error {
	logging.LogDebug("Updating current LED theme to: %s", themeName)

	// Find the selected theme
	var selectedTheme *LEDTheme
	for _, theme := range PredefinedThemes {
		if theme.Name == themeName {
			selectedTheme = &theme
			break
		}
	}

	if selectedTheme == nil {
		return fmt.Errorf("LED theme not found: %s", themeName)
	}

	// Update all light settings with the theme's values
	for i := range CurrentLEDSettings {
		CurrentLEDSettings[i].Effect = int(selectedTheme.Effect)
		CurrentLEDSettings[i].Color1 = selectedTheme.Color
	}

	CurrentLEDTheme = themeName

	logging.LogDebug("LED theme updated in memory: %s", CurrentLEDTheme)
	return nil
}

// ApplyLEDTheme applies the LED theme to all lights
func ApplyLEDTheme(theme *LEDTheme) error {
	// First update the current settings in memory
	if err := UpdateCurrentLEDTheme(theme.Name); err != nil {
		return err
	}

	return ApplyCurrentLEDSettings()
}

// ApplyCurrentLEDSettings applies the current in-memory LED settings to the system
func ApplyCurrentLEDSettings() error {
	settingsPath := GetSettingsPath()

	logging.LogDebug("Applying current LED settings to: %s", settingsPath)

	// Check if settings directory exists
	settingsDir := filepath.Dir(settingsPath)
	if _, err := os.Stat(settingsDir); os.IsNotExist(err) {
		logging.LogDebug("Creating settings directory: %s", settingsDir)
		if err := os.MkdirAll(settingsDir, 0755); err != nil {
			logging.LogDebug("Failed to create settings directory: %v", err)
			return fmt.Errorf("failed to create settings directory: %w", err)
		}
	}

	// Write settings to temp file first
	tempFile := settingsPath + ".tmp"
	outFile, err := os.Create(tempFile)
	if err != nil {
		logging.LogDebug("Error creating temp settings file: %v", err)
		return fmt.Errorf("failed to create temp settings file: %w", err)
	}

	// Write each light section
	for _, light := range CurrentLEDSettings {
		// Convert color values to storage format
		color1 := convertHexFormat(light.Color1, true)
		color2 := convertHexFormat(light.Color2, true)

		fmt.Fprintf(outFile, "[%s]\n", light.Name)
		fmt.Fprintf(outFile, "effect=%d\n", light.Effect)
		fmt.Fprintf(outFile, "color1=%s\n", color1)
		fmt.Fprintf(outFile, "color2=%s\n", color2)
		fmt.Fprintf(outFile, "speed=%d\n", light.Speed)
		fmt.Fprintf(outFile, "brightness=%d\n", light.Brightness)
		fmt.Fprintf(outFile, "trigger=%d\n", light.Trigger)
		fmt.Fprintf(outFile, "filename=%s\n", light.Filename)
		fmt.Fprintf(outFile, "inbrightness=%d\n", light.InBrightness)
		fmt.Fprintf(outFile, "\n")
	}

	outFile.Close()

	// Replace original file with updated one
	err = os.Rename(tempFile, settingsPath)
	if err != nil {
		logging.LogDebug("Error replacing settings file: %v", err)
		return fmt.Errorf("failed to update settings file: %w", err)
	}

	logging.LogDebug("Successfully applied LED settings")
	return nil
}