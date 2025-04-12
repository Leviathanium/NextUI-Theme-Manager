// src/internal/ui/screens/led_menu.go
// Implementation of the LED settings menu screen

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// LEDMenuScreen displays the LED options menu
func LEDMenuScreen() (string, int) {
	// Menu items
	menu := []string{
		"Presets",
		"Custom LEDs",
		"Export Current LEDs",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "LED Themes")
}

// HandleLEDMenu processes the user's selection from the LED menu
func HandleLEDMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleLEDMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Presets":
			logging.LogDebug("Selected LED Presets")
			app.SetSelectedLEDThemeSource(app.PresetSource)
			return app.Screens.LEDSelection

		case "Custom LEDs":
			logging.LogDebug("Selected Custom LEDs")
			app.SetSelectedLEDThemeSource(app.CustomSource)
			return app.Screens.LEDSelection

		case "Export Current LEDs":
			logging.LogDebug("Selected Export Current LEDs")
			return app.Screens.LEDExport

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.LEDMenu
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.CustomizationMenu
	}

	return app.Screens.LEDMenu
}