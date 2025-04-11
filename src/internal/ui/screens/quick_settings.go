// src/internal/ui/screens/quick_settings.go
// Implementation of the quick settings screen

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// QuickSettingsScreen displays quick settings options
func QuickSettingsScreen() (string, int) {
	// Quick settings menu options
	menu := []string{
		"Color Settings",
		"LED Settings",
		"Apply Changes",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Quick Settings")
}

// HandleQuickSettings processes the user's selection from the quick settings menu
func HandleQuickSettings(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleQuickSettings called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Color Settings":
			logging.LogDebug("Selected Color Settings")
			// Get current color values
			r, g, b := app.GetColorSelections()
			// Example of setting modified values
			app.SetColorSelections(r, g, b)
			return app.Screens.QuickSettings

		case "LED Settings":
			logging.LogDebug("Selected LED Settings")
			// Get current LED settings
			brightness, speed := app.GetLEDSelections()
			// Example of setting modified values
			app.SetLEDSelections(brightness, speed)
			return app.Screens.QuickSettings

		case "Apply Changes":
			logging.LogDebug("Selected Apply Changes")
			// Placeholder for applying changes
			return app.Screens.CustomizationMenu

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.QuickSettings
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.CustomizationMenu
	}

	return app.Screens.QuickSettings
}