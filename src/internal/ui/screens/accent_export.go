// src/internal/ui/screens/accent_export.go
// Implementation of the accent theme export screen

package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"bufio"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/accents"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// AccentExportScreen displays the accent theme export screen
func AccentExportScreen() (string, int) {
	// Prompt for a theme name
	return ui.DisplayMinUiList("Enter a name for the theme", "text", "Export Current Accents")
}

// HandleAccentExport processes the user's accent theme export
func HandleAccentExport(themeName string, exitCode int) app.Screen {
	logging.LogDebug("HandleAccentExport called with theme name: '%s', exitCode: %d", themeName, exitCode)

	switch exitCode {
	case 0:
		if themeName != "" {
			// Create a filename with .txt extension
			fileName := themeName + ".txt"

			// Export directly from system settings file to custom themes directory
			err := exportAccentSettings(fileName)
			if err != nil {
				logging.LogDebug("Error exporting accent settings: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("Accent settings exported as: %s", themeName), "3")

				// Refresh the themes list
				err = accents.LoadExternalAccentThemes()
				if err != nil {
					logging.LogDebug("Error refreshing themes: %v", err)
				}
			}
		} else {
			ui.ShowMessage("Export cancelled: No name provided", "3")
		}
		return app.Screens.AccentMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.AccentMenu
	}

	return app.Screens.AccentMenu
}

// exportAccentSettings exports accent settings directly from system settings file to a theme file
func exportAccentSettings(fileName string) error {
	// System settings file path
	settingsPath := accents.SettingsPath

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Target file path in custom themes directory
	customDir := filepath.Join(cwd, accents.AccentsDir, accents.CustomDir)
	if err := os.MkdirAll(customDir, 0755); err != nil {
		logging.LogDebug("Error creating custom themes directory: %v", err)
		return fmt.Errorf("error creating directory: %w", err)
	}

	targetPath := filepath.Join(customDir, fileName)

	// Check if system settings file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		logging.LogDebug("System settings file not found: %s", settingsPath)
		return fmt.Errorf("system settings file not found: %s", settingsPath)
	}

	// Read color settings from the system settings file
	logging.LogDebug("Reading accent settings from: %s", settingsPath)

	// Open settings file for reading
	settingsFile, err := os.Open(settingsPath)
	if err != nil {
		logging.LogDebug("Error opening settings file: %v", err)
		return fmt.Errorf("error opening settings file: %w", err)
	}
	defer settingsFile.Close()

	// Create target file for writing
	targetFile, err := os.Create(targetPath)
	if err != nil {
		logging.LogDebug("Error creating target file: %v", err)
		return fmt.Errorf("error creating target file: %w", err)
	}
	defer targetFile.Close()

	// Read settings file line by line and extract color settings
	scanner := bufio.NewScanner(settingsFile)
	colorSettings := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Keep only color settings
		if strings.HasPrefix(key, "color") && len(key) == 6 && key[5] >= '1' && key[5] <= '6' {
			colorSettings[key] = value
		}
	}

	if scanner.Err() != nil {
		logging.LogDebug("Error scanning settings file: %v", scanner.Err())
		return fmt.Errorf("error reading settings file: %w", scanner.Err())
	}

	// Write color settings to the target file
	for i := 1; i <= 6; i++ {
		key := fmt.Sprintf("color%d", i)
		if value, ok := colorSettings[key]; ok {
			_, err := fmt.Fprintf(targetFile, "%s=%s\n", key, value)
			if err != nil {
				logging.LogDebug("Error writing to target file: %v", err)
				return fmt.Errorf("error writing to target file: %w", err)
			}
		} else {
			logging.LogDebug("Warning: %s not found in settings file", key)
		}
	}

	logging.LogDebug("Successfully exported accent settings to: %s", targetPath)
	return nil
}