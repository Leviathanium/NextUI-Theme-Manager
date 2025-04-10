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
)

// LEDTheme represents a theme for LED lighting
type LEDTheme struct {
	Name   string
	Color  string
	Effect LEDEffect
}

// Settings file paths
const (
	BrickSettingsPath = "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"
	NormalSettingsPath = "/mnt/SDCARD/.userdata/shared/ledsettings.txt"
)

// Predefined LED themes
var PredefinedThemes = []LEDTheme{
	{
		Name:   "White Static",
		Color:  "0xFFFFFF",
		Effect: Static,
	},
	{
		Name:   "Blue Static",
		Color:  "0x3366FF",
		Effect: Static,
	},
	{
		Name:   "Green Static",
		Color:  "0x00AA00",
		Effect: Static,
	},
	{
		Name:   "Red Static",
		Color:  "0xAA0000",
		Effect: Static,
	},
	{
		Name:   "Purple Static",
		Color:  "0x8833FF",
		Effect: Static,
	},
	{
		Name:   "Orange Static",
		Color:  "0xFF8833",
		Effect: Static,
	},
	{
		Name:   "White Breathing",
		Color:  "0xFFFFFF",
		Effect: Breathe,
	},
	{
		Name:   "Blue Breathing",
		Color:  "0x3366FF",
		Effect: Breathe,
	},
	{
		Name:   "Green Breathing",
		Color:  "0x00AA00",
		Effect: Breathe,
	},
	{
		Name:   "Red Breathing",
		Color:  "0xAA0000",
		Effect: Breathe,
	},
	{
		Name:   "Purple Breathing",
		Color:  "0x8833FF",
		Effect: Breathe,
	},
	{
		Name:   "Orange Breathing",
		Color:  "0xFF8833",
		Effect: Breathe,
	},
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

// ApplyLEDTheme applies the LED theme to all lights
func ApplyLEDTheme(theme *LEDTheme) error {
	settingsPath := GetSettingsPath()
	lightCount := GetLightCount()

	logging.LogDebug("Applying LED theme: %s to %s", theme.Name, settingsPath)

	// Check if settings directory exists
	settingsDir := filepath.Dir(settingsPath)
	if _, err := os.Stat(settingsDir); os.IsNotExist(err) {
		logging.LogDebug("Creating settings directory: %s", settingsDir)
		if err := os.MkdirAll(settingsDir, 0755); err != nil {
			logging.LogDebug("Failed to create settings directory: %v", err)
			return fmt.Errorf("failed to create settings directory: %w", err)
		}
	}

	// Read existing settings or create new file
	settingsExists := true
	_, err := os.Stat(settingsPath)
	if os.IsNotExist(err) {
		settingsExists = false
		logging.LogDebug("Settings file does not exist, will create new one")
	}

	// Light settings to modify
	type lightSettings struct {
		name        string
		effect      int
		color1      string
		color2      string
		speed       int
		brightness  int
		trigger     int
		filename    string
		inbrightness int
	}

	// Default values for new settings
	lights := make([]lightSettings, lightCount)

	if IsBrickDevice() {
		// Brick device light names
		lightNames := []string{"F1 key", "F2 key", "Top bar", "L&R triggers"}
		for i := 0; i < lightCount; i++ {
			lights[i] = lightSettings{
				name:        lightNames[i],
				effect:      int(theme.Effect),
				color1:      theme.Color,
				color2:      "0x000000",
				speed:       1000,
				brightness:  100,
				trigger:     1,
				filename:    "",
				inbrightness: 100,
			}
		}
	} else {
		// Normal device light names
		lightNames := []string{"Joysticks", "Logo"}
		for i := 0; i < lightCount; i++ {
			lights[i] = lightSettings{
				name:        lightNames[i],
				effect:      int(theme.Effect),
				color1:      theme.Color,
				color2:      "0x000000",
				speed:       1000,
				brightness:  100,
				trigger:     1,
				filename:    "",
				inbrightness: 100,
			}
		}
	}

	// If settings file exists, read existing values
	if settingsExists {
		file, err := os.Open(settingsPath)
		if err != nil {
			logging.LogDebug("Error opening settings file: %v", err)
			return fmt.Errorf("failed to open settings file: %w", err)
		}

		scanner := bufio.NewScanner(file)
		lightIndex := -1

		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)

			// Check for light section header
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				name := line[1 : len(line)-1]

				// Find matching light
				for i, light := range lights {
					if light.name == name {
						lightIndex = i
						break
					}
				}
				continue
			}

			// Skip if not in a valid light section
			if lightIndex < 0 || lightIndex >= len(lights) {
				continue
			}

            // Parse settings
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			_ = strings.TrimSpace(parts[1]) // Use blank identifier to explicitly ignore the value

			// Only update color1 and effect, preserve other settings
			switch key {
			case "effect":
				// Always override with theme effect
				lights[lightIndex].effect = int(theme.Effect)
			case "color1":
				// Always override with theme color
				lights[lightIndex].color1 = theme.Color
			}
		}

		file.Close()
	}

	// Write updated settings
	tempFile := settingsPath + ".tmp"
	outFile, err := os.Create(tempFile)
	if err != nil {
		logging.LogDebug("Error creating temp settings file: %v", err)
		return fmt.Errorf("failed to create temp settings file: %w", err)
	}

	// Write each light section
	for _, light := range lights {
		fmt.Fprintf(outFile, "[%s]\n", light.name)
		fmt.Fprintf(outFile, "effect=%d\n", light.effect)
		fmt.Fprintf(outFile, "color1=%s\n", light.color1)
		fmt.Fprintf(outFile, "color2=%s\n", light.color2)
		fmt.Fprintf(outFile, "speed=%d\n", light.speed)
		fmt.Fprintf(outFile, "brightness=%d\n", light.brightness)
		fmt.Fprintf(outFile, "trigger=%d\n", light.trigger)
		fmt.Fprintf(outFile, "filename=%s\n", light.filename)
		fmt.Fprintf(outFile, "inbrightness=%d\n", light.inbrightness)
		fmt.Fprintf(outFile, "\n")
	}

	outFile.Close()

	// Replace original file with updated one
	err = os.Rename(tempFile, settingsPath)
	if err != nil {
		logging.LogDebug("Error replacing settings file: %v", err)
		return fmt.Errorf("failed to update settings file: %w", err)
	}

	logging.LogDebug("Successfully applied LED theme")
	return nil
}