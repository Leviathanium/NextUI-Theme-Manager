// src/internal/ui/screens/customization_menu.go
// Implementation of the customization menu screen

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// CustomizationMenuScreen displays the customization menu options
func CustomizationMenuScreen() (string, int) {
	// Menu items
	menu := []string{
		"Global Options",
		"System Options",
		"Accents",
		"LEDs",
		"Fonts",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Customization")
}

// HandleCustomizationMenu processes the user's selection from the customization menu
func HandleCustomizationMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleCustomizationMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Global Options":
			logging.LogDebug("Selected Global Options")
			return app.Screens.GlobalOptionsMenu

		case "System Options":
			logging.LogDebug("Selected System Options")
			return app.Screens.SystemOptionsMenu

		case "Accents":
			logging.LogDebug("Selected Accents")
			return app.Screens.AccentMenu

		case "LEDs":
			logging.LogDebug("Selected LEDs")
			return app.Screens.LEDMenu

		case "Fonts":
			logging.LogDebug("Selected Fonts")
			return app.Screens.FontSelection

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.CustomizationMenu
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.CustomizationMenu
}