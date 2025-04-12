// src/internal/ui/screens/led_export.go
// Implementation of the LED theme export screen

package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"bufio"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/leds"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// LEDExportScreen displays the LED theme export screen
func LEDExportScreen() (string, int) {
	// Prompt for a theme name
	return ui.DisplayMinUiList("Enter a name for the LED theme", "text", "Export Current LEDs")
}

// HandleLEDExport processes the user's LED theme export
func HandleLEDExport(themeName string, exitCode int) app.Screen {
	logging.LogDebug("HandleLEDExport called with theme name: '%s', exitCode: %d", themeName, exitCode)

	switch exitCode {
	case 0:
		if themeName != "" {
			// Create a filename with .txt extension
			fileName := themeName + ".txt"

			// Export directly from system LED settings file to custom themes directory
			err := exportLEDSettings(fileName)
			if err != nil {
				logging.LogDebug("Error exporting LED settings: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("LED settings exported as: %s", themeName), "3")

				// Refresh the LED themes list
				err = leds.LoadExternalLEDThemes()
				if err != nil {
					logging.LogDebug("Error refreshing LED themes: %v", err)
				}
			}
		} else {
			ui.ShowMessage("Export cancelled: No name provided", "3")
		}
		return app.Screens.LEDMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.LEDMenu
	}

	return app.Screens.LEDMenu
}

// exportLEDSettings exports LED settings directly from system settings file to a theme file
func exportLEDSettings(fileName string) error {
	// System LED settings file path
	settingsPath := leds.GetSettingsPath()

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Target file path in custom themes directory
	customDir := filepath.Join(cwd, leds.LEDsDir, leds.CustomDir)
	if err := os.MkdirAll(customDir, 0755); err != nil {
		logging.LogDebug("Error creating custom LED themes directory: %v", err)
		return fmt.Errorf("error creating directory: %w", err)
	}

	targetPath := filepath.Join(customDir, fileName)

	// Check if system LED settings file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		logging.LogDebug("System LED settings file not found: %s", settingsPath)
		return fmt.Errorf("system LED settings file not found: %s", settingsPath)
	}

	logging.LogDebug("Reading LED settings from: %s", settingsPath)

	// Open settings file for reading
	settingsFile, err := os.Open(settingsPath)
	if err != nil {
		logging.LogDebug("Error opening LED settings file: %v", err)
		return fmt.Errorf("error opening LED settings file: %w", err)
	}
	defer settingsFile.Close()

	// Create target file for writing
	targetFile, err := os.Create(targetPath)
	if err != nil {
		logging.LogDebug("Error creating target file: %v", err)
		return fmt.Errorf("error creating target file: %w", err)
	}
	defer targetFile.Close()

	// For LED themes, we just want to keep the section headers and color1 values
	// since we're only interested in static colors
	scanner := bufio.NewScanner(settingsFile)
	var currentSection string
	var colors = make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check if this is a section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line
			continue
		}

		// If we have a current section and this is a color1 setting, store it
		if currentSection != "" && strings.HasPrefix(line, "color1=") {
			colors[currentSection] = line
		}
	}

	if scanner.Err() != nil {
		logging.LogDebug("Error scanning LED settings file: %v", scanner.Err())
		return fmt.Errorf("error reading LED settings file: %w", scanner.Err())
	}

	// Now write the section headers and color1 values to the target file
	for section, color := range colors {
		_, err := fmt.Fprintf(targetFile, "%s\n%s\n\n", section, color)
		if err != nil {
			logging.LogDebug("Error writing to target file: %v", err)
			return fmt.Errorf("error writing to target file: %w", err)
		}
	}

	logging.LogDebug("Successfully exported LED settings to: %s", targetPath)
	return nil
}