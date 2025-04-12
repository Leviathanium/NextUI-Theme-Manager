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
		"Systems",
		"Accents",
		"LED Quick Settings",
		"Icons",
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
		case "Systems":
			logging.LogDebug("Selected Systems (System Backgrounds)")
			app.SetSelectedThemeType(app.CustomTheme)
			return app.Screens.ThemeSelection

		case "Accents":
			logging.LogDebug("Selected Accents")
			return app.Screens.AccentMenu

		case "LED Quick Settings":
			logging.LogDebug("Selected LED Quick Settings")
			return app.Screens.LEDMenu

		case "Icons":
			logging.LogDebug("Selected Icons")
			return app.Screens.IconSelection

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