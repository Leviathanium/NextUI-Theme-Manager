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
	Name         string
	Effect       int
	Color1       string
	Color2       string
	Speed        int
	Brightness   int
	Trigger      int
	Filename     string
	InBrightness int
}

// Settings file paths
const (
	BrickSettingsPath  = "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"
	NormalSettingsPath = "/mnt/SDCARD/.userdata/shared/ledsettings.txt"
	LEDsDir            = "LEDs"      // Directory for external LED theme files
	PresetsDir         = "Presets"   // Subdirectory for preset themes
	CustomDir          = "Custom"    // Subdirectory for custom themes
	StaticEffectValue  = 4           // Using Static Color effect (index 4)
)

// Current LED settings in memory
var (
	CurrentLEDSettings []LightSettings
	CurrentLEDTheme    string
	PresetLEDThemes    []LEDTheme
	CustomLEDThemes    []LEDTheme
)

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
				Name:         "F1 key",
				Effect:       StaticEffectValue,
				Color1:       "#FFFFFF",
				Color2:       "#000000",
				Speed:        1000,
				Brightness:   100,
				Trigger:      1,
				Filename:     "",
				InBrightness: 100,
			},
			{
				Name:         "F2 key",
				Effect:       StaticEffectValue,
				Color1:       "#FFFFFF",
				Color2:       "#000000",
				Speed:        1000,
				Brightness:   100,
				Trigger:      1,
				Filename:     "",
				InBrightness: 100,
			},
			{
				Name:         "Top bar",
				Effect:       StaticEffectValue,
				Color1:       "#FFFFFF",
				Color2:       "#000000",
				Speed:        1000,
				Brightness:   100,
				Trigger:      1,
				Filename:     "",
				InBrightness: 100,
			},
			{
				Name:         "L&R triggers",
				Effect:       StaticEffectValue,
				Color1:       "#FFFFFF",
				Color2:       "#000000",
				Speed:        1000,
				Brightness:   100,
				Trigger:      1,
				Filename:     "",
				InBrightness: 100,
			},
		}
	} else {
		CurrentLEDSettings = settings
		DetermineCurrentTheme()
	}

	// Load external LED theme files
	if err := LoadExternalLEDThemes(); err != nil {
		logging.LogDebug("Warning: Could not load external LED themes: %v", err)
	}

	// Create placeholder files
	if err := CreatePlaceholderFiles(); err != nil {
		logging.LogDebug("Warning: Could not create placeholder files: %v", err)
	}

	logging.LogDebug("LED settings initialized with theme: %s", CurrentLEDTheme)
	return nil
}

// LoadExternalLEDThemes loads LED themes from external files
func LoadExternalLEDThemes() error {
	// Clear the current lists of external themes
	PresetLEDThemes = []LEDTheme{}
	CustomLEDThemes = []LEDTheme{}

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Load preset themes
	presetsDir := filepath.Join(cwd, LEDsDir, PresetsDir)
	logging.LogDebug("Loading preset LED themes from: %s", presetsDir)
	if err := loadThemesFromDir(presetsDir, &PresetLEDThemes); err != nil {
		logging.LogDebug("Warning: Could not load preset LED themes: %v", err)
	}

	// Load custom themes
	customDir := filepath.Join(cwd, LEDsDir, CustomDir)
	logging.LogDebug("Loading custom LED themes from: %s", customDir)
	if err := loadThemesFromDir(customDir, &CustomLEDThemes); err != nil {
		logging.LogDebug("Warning: Could not load custom LED themes: %v", err)
	}

	logging.LogDebug("Loaded %d preset and %d custom LED themes", len(PresetLEDThemes), len(CustomLEDThemes))
	return nil
}

// loadThemesFromDir loads themes from a specific directory
func loadThemesFromDir(themesDir string, themesList *[]LEDTheme) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		logging.LogDebug("Error creating LED themes directory: %v", err)
		return fmt.Errorf("error creating LED themes directory: %w", err)
	}

	// Read the directory
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		logging.LogDebug("Error reading LED themes directory: %v", err)
		return fmt.Errorf("error reading LED themes directory: %w", err)
	}

	// Process each file in the directory
	for _, entry := range entries {
		// Skip directories and hidden files
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip files that don't have a .txt extension
		if !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		// Skip placeholder files
		if strings.Contains(entry.Name(), "Place-") && strings.Contains(entry.Name(), "-Here") {
			continue
		}

		// Extract theme name (remove .txt extension)
		themeName := strings.TrimSuffix(entry.Name(), ".txt")

		// Load the LED settings from the file
		lightSettings, err := ReadLEDSettingsFile(filepath.Join(themesDir, entry.Name()))
		if err != nil {
			logging.LogDebug("Error reading LED settings file %s: %v", entry.Name(), err)
			continue
		}

		// Create a theme from the settings
		if len(lightSettings) > 0 {
			// Use the first light's settings for the theme
			theme := LEDTheme{
				Name:   themeName,
				Color:  lightSettings[0].Color1,
				Effect: StaticColor, // Always use static color effect (per user request)
			}

			// Add the theme to the list
			*themesList = append(*themesList, theme)
		}
	}

	return nil
}

// ReadLEDSettingsFile reads LED settings from a file
func ReadLEDSettingsFile(filepath string) ([]LightSettings, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open LED settings file: %w", err)
	}
	defer file.Close()

	var settings []LightSettings
	var currentLight *LightSettings

	// Read the file line by line
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
				// Default values
				Effect:       StaticEffectValue, // Always use static color effect
				Speed:        1000,
				Brightness:   100,
				Trigger:      1,
				InBrightness: 100,
				Color2:       "#000000", // Default black for secondary color
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

		// Set the appropriate field - only care about color1 (as per user request)
		if key == "color1" {
			currentLight.Color1 = convertHexFormat(value, false) // Convert to display format
		}
	}

	// Add the last light if there is one
	if currentLight != nil {
		settings = append(settings, *currentLight)
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("error reading LED settings file: %w", scanner.Err())
	}

	return settings, nil
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

// DetermineCurrentTheme tries to determine which theme is closest to the current settings
func DetermineCurrentTheme() {
	// Default
	CurrentLEDTheme = "Custom"

	if len(CurrentLEDSettings) == 0 {
		return
	}

	// Use the first light as reference
	firstLight := CurrentLEDSettings[0]

	// Check if all lights have the same color
	allSame := true
	for i := 1; i < len(CurrentLEDSettings); i++ {
		if CurrentLEDSettings[i].Color1 != firstLight.Color1 {
			allSame = false
			break
		}
	}

	if !allSame {
		return
	}

	// First try to find a match in preset themes
	for _, theme := range PresetLEDThemes {
		if convertHexFormat(theme.Color, true) == convertHexFormat(firstLight.Color1, true) {
			CurrentLEDTheme = theme.Name
			return
		}
	}

	// Then try to find a match in custom themes
	for _, theme := range CustomLEDThemes {
		if convertHexFormat(theme.Color, true) == convertHexFormat(firstLight.Color1, true) {
			CurrentLEDTheme = theme.Name
			return
		}
	}
}

// UpdateCurrentLEDTheme updates the current LED theme in memory with new values
func UpdateCurrentLEDTheme(themeName string) error {
	logging.LogDebug("Updating current LED theme to: %s", themeName)

	// First look in preset themes
	for _, theme := range PresetLEDThemes {
		if theme.Name == themeName {
			// Update all light settings with the theme's values
			for i := range CurrentLEDSettings {
				CurrentLEDSettings[i].Effect = StaticEffectValue // Always use static color effect
				CurrentLEDSettings[i].Color1 = theme.Color
			}

			CurrentLEDTheme = themeName
			logging.LogDebug("LED theme updated from preset theme: %s", CurrentLEDTheme)
			return nil
		}
	}

	// Then look in custom themes
	for _, theme := range CustomLEDThemes {
		if theme.Name == themeName {
			// Update all light settings with the theme's values
			for i := range CurrentLEDSettings {
				CurrentLEDSettings[i].Effect = StaticEffectValue // Always use static color effect
				CurrentLEDSettings[i].Color1 = theme.Color
			}

			CurrentLEDTheme = themeName
			logging.LogDebug("LED theme updated from custom theme: %s", CurrentLEDTheme)
			return nil
		}
	}

	logging.LogDebug("LED theme not found: %s", themeName)
	return fmt.Errorf("LED theme not found: %s", themeName)
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

// SaveLEDThemeToFile saves the current LED settings to an external file
func SaveLEDThemeToFile(fileName string, isCustom bool) error {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to themes directory
	var themesDir string
	if isCustom {
		themesDir = filepath.Join(cwd, LEDsDir, CustomDir)
	} else {
		themesDir = filepath.Join(cwd, LEDsDir, PresetsDir)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		logging.LogDebug("Error creating LED themes directory: %v", err)
		return fmt.Errorf("error creating LED themes directory: %w", err)
	}

	// Full path to the file
	filePath := filepath.Join(themesDir, fileName)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		logging.LogDebug("Error creating LED theme file: %v", err)
		return fmt.Errorf("error creating LED theme file: %w", err)
	}
	defer file.Close()

	// Get the color from the first light - we only care about color (as per user request)
	if len(CurrentLEDSettings) == 0 {
		return fmt.Errorf("no current LED settings to save")
	}

	color1 := convertHexFormat(CurrentLEDSettings[0].Color1, true)

	// Write just enough information to allow loading the color
	for _, light := range CurrentLEDSettings {
		fmt.Fprintf(file, "[%s]\n", light.Name)
		fmt.Fprintf(file, "color1=%s\n", color1)
		fmt.Fprintf(file, "\n")
	}

	logging.LogDebug("Successfully saved LED theme to file: %s", filePath)
	return nil
}

// CreatePlaceholderFiles creates placeholder files in the theme directories
func CreatePlaceholderFiles() error {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create directories
	presetsDir := filepath.Join(cwd, LEDsDir, PresetsDir)
	customDir := filepath.Join(cwd, LEDsDir, CustomDir)

	if err := os.MkdirAll(presetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create presets directory: %w", err)
	}

	if err := os.MkdirAll(customDir, 0755); err != nil {
		return fmt.Errorf("failed to create custom directory: %w", err)
	}

	// Create placeholder file in custom directory if empty
	entries, err := os.ReadDir(customDir)
	if err != nil {
		return fmt.Errorf("failed to read custom directory: %w", err)
	}

	if len(entries) == 0 {
		placeholderPath := filepath.Join(customDir, "Place-LED-Files-Here.txt")
		file, err := os.Create(placeholderPath)
		if err != nil {
			return fmt.Errorf("failed to create placeholder file: %w", err)
		}

		_, err = file.WriteString("# Place custom LED theme files in this directory\n\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to write placeholder content: %w", err)
		}

		_, err = file.WriteString("# Format should be:\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to write placeholder content: %w", err)
		}

		_, err = file.WriteString("[F1 key]\ncolor1=0xRRGGBB\n\n[F2 key]\ncolor1=0xRRGGBB\n\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to write placeholder content: %w", err)
		}

		_, err = file.WriteString("[Top bar]\ncolor1=0xRRGGBB\n\n[L&R triggers]\ncolor1=0xRRGGBB\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to write placeholder content: %w", err)
		}

		file.Close()
	}

	return nil
}

// Static LED colors for preset themes
var StaticLEDColors = []struct {
	Name  string
	Color string
}{
	{Name: "White", Color: "#FFFFFF"},
	{Name: "Blue", Color: "#3366FF"},
	{Name: "Green", Color: "#00AA00"},
	{Name: "Red", Color: "#AA0000"},
	{Name: "Purple", Color: "#8833FF"},
	{Name: "Orange", Color: "#FF8833"},
	{Name: "Teal", Color: "#00AAAA"},
	{Name: "Pink", Color: "#FF66FF"},
	{Name: "Yellow", Color: "#FFFF33"},
}