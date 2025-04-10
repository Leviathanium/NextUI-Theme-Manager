// src/internal/ui/screens/quick_settings.go
// Implementation of the quick settings screen

package screens

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/fonts"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// Settings file paths
const (
	SettingsPath     = "/mnt/SDCARD/.userdata/shared/minuisettings.txt"
	BrickSettingsPath = "/mnt/SDCARD/.userdata/shared/ledsettings_brick.txt"
	NormalSettingsPath = "/mnt/SDCARD/.userdata/shared/ledsettings.txt"
)

// LED effect names
var effectNames = []string{
	"Static",
	"Breathe",
	"Interval Breathe",
	"Static",
	"Blink 1",
	"Blink 2",
	"Blink 3",
	"Rainbow",
	"Twinkle",
	"Fire",
	"Glitter",
	"NeonGlow",
	"Firefly",
	"Aurora",
	"Reactive",
}

// Standard colors for LED and accent options
var standardColors = []string{
	"#FF0000", // Red
	"#00FF00", // Green
	"#0000FF", // Blue
	"#FFFF00", // Yellow
	"#FF00FF", // Magenta
	"#00FFFF", // Cyan
	"#FFFFFF", // White
	"#000000", // Black
	"#FF8800", // Orange
	"#8800FF", // Purple
	"#00AA00", // Dark Green
	"#0000AA", // Dark Blue
	"#AA0000", // Dark Red
	"#9B2257", // Accent Pink (default)
}

// QuickSettingsItem represents a JSON-formatted item for the list
type ListItem struct {
	Name     string   `json:"name"`
	Options  []string `json:"options,omitempty"`
	Selected int      `json:"selected,omitempty"`
}

// QuickSettingsScreen displays the quick settings screen
func QuickSettingsScreen() (string, int) {
	// Create a temporary JSON file for the list
	tempFile, err := os.CreateTemp("", "quick-settings-*.json")
	if err != nil {
		logging.LogDebug("Error creating temp file: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}
	defer os.Remove(tempFile.Name())

	// Build the JSON list items
	var items []ListItem

	// Font item
	fonts, err := getFontsList()
	if err != nil {
		logging.LogDebug("Error getting fonts: %v", err)
		fonts = []string{"Default Font"}
	}
	fontItem := ListItem{
		Name:     "Font",
		Options:  fonts,
		Selected: getCurrentFontIndex(fonts),
	}
	items = append(items, fontItem)

	// Accent item
	accentItem := ListItem{
		Name:     "Accent",
		Options:  standardColors,
		Selected: getCurrentAccentIndex(),
	}
	items = append(items, accentItem)

	// LED Color item
	ledColorItem := ListItem{
		Name:     "LED Color",
		Options:  standardColors,
		Selected: getCurrentLEDColorIndex(),
	}
	items = append(items, ledColorItem)

	// LED Effect item
	ledEffectItem := ListItem{
		Name:     "LED Effect",
		Options:  effectNames,
		Selected: getCurrentLEDEffectIndex(),
	}
	items = append(items, ledEffectItem)

	// Convert to JSON and write to temp file
	jsonData := map[string]interface{}{
		"items": items,
	}
	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		logging.LogDebug("Error marshaling JSON: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	_, err = tempFile.Write(jsonBytes)
	if err != nil {
		logging.LogDebug("Error writing to temp file: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Close the file to ensure content is flushed
	tempFile.Close()

	// Display the list using minui-list
	return ui.DisplayMinUiList("", "json", "Quick Settings", "--file", tempFile.Name(), "--item-key", "items")
}

// HandleQuickSettingsScreen processes user interaction with the quick settings screen
func HandleQuickSettingsScreen(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleQuickSettingsScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	// Check if user pressed back or cancel
	if exitCode == 1 || exitCode == 2 {
		return app.Screens.MainMenu
	}

	// Check if user made a selection
	if exitCode == 0 {
		// Parse the JSON result
		var result map[string]interface{}
		err := json.Unmarshal([]byte(selection), &result)
		if err != nil {
			logging.LogDebug("Error parsing JSON result: %v", err)
			ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			return app.Screens.MainMenu
		}

		// Get the selected item and index
		itemsArray, ok := result["items"].([]interface{})
		if !ok || len(itemsArray) == 0 {
			logging.LogDebug("Invalid JSON result format")
			return app.Screens.MainMenu
		}

		// Apply settings based on selections
		applySettings(itemsArray)

		// Return to main menu
		return app.Screens.MainMenu
	}

	// Default case, stay on the same screen
	return app.Screens.QuickSettings
}

// Helper functions for getting current settings

// getFontsList returns the list of available fonts
func getFontsList() ([]string, error) {
	return fonts.ListFonts()
}

// getCurrentFontIndex returns the index of the currently selected font
func getCurrentFontIndex(fontsList []string) int {
	currentFont := app.GetSelectedFont()
	if currentFont == "" {
		return 0
	}

	for i, font := range fontsList {
		if font == currentFont {
			return i
		}
	}
	return 0
}

// getCurrentAccentIndex returns the index of the current accent color
func getCurrentAccentIndex() int {
	// Read minuisettings.txt to get current color2 value
	color := readSettingsValue(SettingsPath, "color2")
	if color == "" {
		return 0 // Default to first color
	}

	// Convert from 0xRRGGBB to #RRGGBB format
	hashColor := "#" + strings.ToUpper(color[2:])

	// Find in standardColors
	for i, c := range standardColors {
		if strings.EqualFold(c, hashColor) {
			return i
		}
	}

	return 0 // Default to first color if not found
}

// getCurrentLEDColorIndex returns the index of the current LED color
func getCurrentLEDColorIndex() int {
	// Determine if we're on a brick device
	isBrick := os.Getenv("DEVICE") == "brick"
	settingsPath := NormalSettingsPath
	if isBrick {
		settingsPath = BrickSettingsPath
	}

	// Get the first light's color1 value
	// In a real implementation, we might want to check if all lights have the same color
	lightSettings := readLEDSettings(settingsPath)
	if len(lightSettings) == 0 {
		return 6 // Default to white
	}

	// Get color1 for the first light
	color := lightSettings[0]["color1"]
	if color == "" {
		return 6 // Default to white
	}

	// Convert from 0xRRGGBB to #RRGGBB format
	hashColor := "#" + strings.ToUpper(color[2:])

	// Find in standardColors
	for i, c := range standardColors {
		if strings.EqualFold(c, hashColor) {
			return i
		}
	}

	return 6 // Default to white if not found
}

// getCurrentLEDEffectIndex returns the index of the current LED effect
func getCurrentLEDEffectIndex() int {
	// Determine if we're on a brick device
	isBrick := os.Getenv("DEVICE") == "brick"
	settingsPath := NormalSettingsPath
	if isBrick {
		settingsPath = BrickSettingsPath
	}

	// Get the first light's effect value
	lightSettings := readLEDSettings(settingsPath)
	if len(lightSettings) == 0 {
		return 0 // Default to static
	}

	// Get effect for the first light
	effectStr := lightSettings[0]["effect"]
	if effectStr == "" {
		return 0 // Default to static
	}

	// Convert to int
	effect, err := strconv.Atoi(effectStr)
	if err != nil {
		return 0 // Default to static
	}

	// LED effects are 1-indexed in the file but 0-indexed in our array
	effectIndex := effect - 1
	if effectIndex < 0 || effectIndex >= len(effectNames) {
		return 0 // Default to static if out of range
	}

	return effectIndex
}

// Helper functions for parsing settings files

// readSettingsValue reads a specific key's value from a settings file
func readSettingsValue(filePath, key string) string {
	file, err := os.Open(filePath)
	if err != nil {
		logging.LogDebug("Error opening settings file %s: %v", filePath, err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			return strings.TrimSpace(parts[1])
		}
	}

	return ""
}

// readLEDSettings reads all settings from a LED settings file
func readLEDSettings(filePath string) []map[string]string {
	file, err := os.Open(filePath)
	if err != nil {
		logging.LogDebug("Error opening LED settings file %s: %v", filePath, err)
		return nil
	}
	defer file.Close()

	var lights []map[string]string
	var currentLight map[string]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Check for light section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if currentLight != nil {
				lights = append(lights, currentLight)
			}
			currentLight = make(map[string]string)
			continue
		}

		// Parse key-value pairs
		if currentLight != nil && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				currentLight[key] = value
			}
		}
	}

	// Add the last light
	if currentLight != nil {
		lights = append(lights, currentLight)
	}

	return lights
}

// Helper functions for applying settings

// applySettings applies all settings based on user selections
func applySettings(items []interface{}) {
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := itemMap["name"].(string)
		if !ok {
			continue
		}

		selected, ok := itemMap["selected"].(float64)
		if !ok {
			continue
		}

		options, ok := itemMap["options"].([]interface{})
		if !ok || int(selected) >= len(options) {
			continue
		}

		selectedOption := options[int(selected)].(string)

		switch name {
		case "Font":
			applyFontSetting(selectedOption)
		case "Accent":
			applyAccentSetting(selectedOption)
		case "LED Color":
			applyLEDColorSetting(selectedOption)
		case "LED Effect":
			applyLEDEffectSetting(int(selected) + 1) // Effects are 1-indexed in the file
		}
	}

	// Show confirmation message
	ui.ShowMessage("Settings applied successfully", "3")
}

// applyFontSetting applies the selected font
func applyFontSetting(fontName string) {
	// Save the selected font to app state
	app.SetSelectedFont(fontName)

	// Apply the font
	err := fonts.ApplyFont(fontName)
	if err != nil {
		logging.LogDebug("Error applying font: %v", err)
		return
	}

	logging.LogDebug("Applied font: %s", fontName)
}

// applyAccentSetting applies the selected accent color
func applyAccentSetting(colorHex string) {
	// Convert from #RRGGBB to 0xRRGGBB format
	colorValue := "0x" + colorHex[1:]

	// Update settings file
	updateSettingsFile(SettingsPath, "color2", colorValue)

	// Update app state (keeping existing colors for color1 and color3)
	// Find index of selected color in standardColors
	selectedIndex := 0
	for i, c := range standardColors {
		if strings.EqualFold(c, colorHex) {
			selectedIndex = i
			break
		}
	}

	// Get current color selections from app state
	color1, _, color3 := app.GetColorSelections()

	// Update color2 value and store back to app state
	app.SetColorSelections(color1, selectedIndex, color3)

	logging.LogDebug("Applied accent color: %s (index: %d)", colorValue, selectedIndex)
}

// applyLEDColorSetting applies the selected LED color to all lights
func applyLEDColorSetting(colorHex string) {
	// Convert from #RRGGBB to 0xRRGGBB format
	colorValue := "0x" + colorHex[1:]

	// Determine if we're on a brick device
	isBrick := os.Getenv("DEVICE") == "brick"
	settingsPath := NormalSettingsPath
	if isBrick {
		settingsPath = BrickSettingsPath
	}

	// Update all lights in the file
	updateLEDSettingsFile(settingsPath, "color1", colorValue)

	// Find index of selected color in standardColors
	selectedIndex := 0
	for i, c := range standardColors {
		if strings.EqualFold(c, colorHex) {
			selectedIndex = i
			break
		}
	}

	// Get current LED effect from app state
	effect, _ := app.GetLEDSelections()

	// Update color value and store back to app state
	app.SetLEDSelections(effect, selectedIndex)

	logging.LogDebug("Applied LED color: %s (index: %d)", colorValue, selectedIndex)
}

// applyLEDEffectSetting applies the selected LED effect to all lights
func applyLEDEffectSetting(effectIndex int) {
	// Determine if we're on a brick device
	isBrick := os.Getenv("DEVICE") == "brick"
	settingsPath := NormalSettingsPath
	if isBrick {
		settingsPath = BrickSettingsPath
	}

	// Update all lights in the file
	updateLEDSettingsFile(settingsPath, "effect", strconv.Itoa(effectIndex))

	// Get current LED color from app state
	_, color := app.GetLEDSelections()

	// Update effect value and store back to app state
	app.SetLEDSelections(effectIndex-1, color) // Store 0-based index in app state

	logging.LogDebug("Applied LED effect: %d", effectIndex)
}

// Helper functions for updating settings files

// updateSettingsFile updates a specific key's value in a settings file
func updateSettingsFile(filePath, key, value string) {
	// Read all lines from the file
	file, err := os.Open(filePath)
	if err != nil {
		logging.LogDebug("Error opening settings file %s: %v", filePath, err)
		return
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	// Update the value
	found := false
	for i, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			lines[i] = key + "=" + value
			found = true
			break
		}
	}

	// If key not found, add it
	if !found {
		lines = append(lines, key + "=" + value)
	}

	// Write the file back
	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		logging.LogDebug("Error writing settings file %s: %v", filePath, err)
	}
}

// updateLEDSettingsFile updates a specific key's value for all lights in a LED settings file
func updateLEDSettingsFile(filePath, key, value string) {
	// Read the entire file
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		logging.LogDebug("Error reading LED settings file %s: %v", filePath, err)
		return
	}

	// Split into lines and process
	lines := strings.Split(string(fileBytes), "\n")
	inLightSection := false

	for i, line := range lines {
		trimLine := strings.TrimSpace(line)

		// Check for light section header
		if strings.HasPrefix(trimLine, "[") && strings.HasSuffix(trimLine, "]") {
			inLightSection = true
			continue
		}

		// Process key-value pairs within a light section
		if inLightSection && strings.Contains(trimLine, "=") {
			parts := strings.SplitN(trimLine, "=", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
				lines[i] = key + "=" + value
			}
		}

		// Reset section flag on empty line
		if trimLine == "" {
			inLightSection = false
		}
	}

	// Write the file back
	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		logging.LogDebug("Error writing LED settings file %s: %v", filePath, err)
	}

	// Reload LED settings
	err = os.WriteFile("/tmp/update_leds", []byte("1"), 0644)
	if err != nil {
		logging.LogDebug("Error triggering LED update: %v", err)
	}
}