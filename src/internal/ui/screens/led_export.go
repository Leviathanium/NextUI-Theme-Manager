// src/internal/ui/screens/led_export.go
// Implementation of the LED theme export screen

package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"bufio"
	"strings"
	"regexp"
	"strconv"

	"nextui-themes/internal/app"
	"nextui-themes/internal/leds"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// LEDExportScreen exports the current LED settings with a sequential name
func LEDExportScreen() (string, int) {
	// Generate sequential file name (LEDs_1, LEDs_2, etc.)
	fileName := generateSequentialLEDFileName()

	// Export directly using the generated file name
	err := exportLEDSettings(fileName)
	if err != nil {
		logging.LogDebug("Error exporting LED settings: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
	} else {
		ui.ShowMessage(fmt.Sprintf("LED settings exported as: %s", fileName), "3")

		// Refresh the LED themes list
		err = leds.LoadExternalLEDThemes()
		if err != nil {
			logging.LogDebug("Error refreshing LED themes: %v", err)
		}
	}

	return "", 0 // Return to LED menu after exporting
}

// HandleLEDExport processes the LED theme export
func HandleLEDExport(selection string, exitCode int) app.Screen {
	// Simply return to the LED menu after export
	return app.Screens.LEDMenu
}

// generateSequentialLEDFileName generates a sequential file name (LEDs_1, LEDs_2, etc.)
func generateSequentialLEDFileName() string {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return "LEDs_1.txt"
	}

	// Custom LEDs directory path
	customDir := filepath.Join(cwd, leds.LEDsDir, leds.CustomDir)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(customDir, 0755); err != nil {
		logging.LogDebug("Error creating custom LEDs directory: %v", err)
		return "LEDs_1.txt"
	}

	// Read the directory
	entries, err := os.ReadDir(customDir)
	if err != nil {
		logging.LogDebug("Error reading custom LEDs directory: %v", err)
		return "LEDs_1.txt"
	}

	// Find the highest number in existing "LEDs_X.txt" files
	highestNum := 0
	regex := regexp.MustCompile(`^LEDs_(\d+)\.txt$`)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		matches := regex.FindStringSubmatch(entry.Name())
		if len(matches) == 2 {
			num, err := strconv.Atoi(matches[1])
			if err == nil && num > highestNum {
				highestNum = num
			}
		}
	}

	// Generate new file name with the next number
	return fmt.Sprintf("LEDs_%d.txt", highestNum+1)
}

// exportLEDSettings exports LED settings to a file
func exportLEDSettings(fileName string) error {
	// System LED settings file path
	settingsPath := leds.GetSettingsPath()

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Target file path in custom LEDs directory
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