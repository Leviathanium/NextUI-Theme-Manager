// src/internal/ui/screens/led_export.go
// Implementation of the LED theme export screen

package screens

import (
	"fmt"

	"nextui-themes/internal/app"
	"nextui-themes/internal/leds"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// LEDExportScreen displays the LED theme export screen
func LEDExportScreen() (string, int) {
	// Prompt for a theme name
	return ui.DisplayMinUiList("Enter a name for the LED theme", "text", "Export Current LED Settings")
}

// HandleLEDExport processes the user's LED theme export
func HandleLEDExport(themeName string, exitCode int) app.Screen {
	logging.LogDebug("HandleLEDExport called with theme name: '%s', exitCode: %d", themeName, exitCode)

	switch exitCode {
	case 0:
		if themeName != "" {
			// Create a filename with .txt extension
			fileName := themeName + ".txt"

			// Save current LED settings to a file in the custom themes directory
			logging.LogDebug("Exporting current LED settings as: %s", themeName)
			err := leds.SaveLEDThemeToFile(fileName, true) // true = custom
			if err != nil {
				logging.LogDebug("Error exporting LED theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("LED theme exported as: %s", themeName), "3")

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