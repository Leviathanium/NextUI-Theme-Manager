// src/internal/ui/screens/customization_menu.go
// Implementation of the customization menu screen

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// CustomizationMenuScreen shows customization options
func CustomizationMenuScreen() (string, int) {
	// Customization submenu options
	menu := []string{
		"System Backgrounds",
		"Fonts",
		"Accents",
		"LED Quick Settings",
		// "Quick Settings" option removed
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Customization Options")
}

// HandleCustomizationMenu processes the user's selection from the customization menu
func HandleCustomizationMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleCustomizationMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "System Backgrounds":
			logging.LogDebug("Selected System Backgrounds")
			app.SetSelectedThemeType(app.CustomTheme)
			return app.Screens.ThemeSelection

		case "Fonts":
			logging.LogDebug("Selected Fonts")
			return app.Screens.FontSelection

		case "Accents":
			logging.LogDebug("Selected Accents")
			return app.Screens.AccentSelection

		case "LED Quick Settings":
			logging.LogDebug("Selected LED Quick Settings")
			return app.Screens.LEDSelection

		// "Quick Settings" case removed

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