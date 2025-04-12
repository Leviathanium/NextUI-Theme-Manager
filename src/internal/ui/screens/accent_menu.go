// src/internal/ui/screens/accent_menu.go
// Implementation of the accent color menu screen

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// AccentMenuScreen displays the accent options menu
func AccentMenuScreen() (string, int) {
	// Menu items
	menu := []string{
		"Presets",
		"Custom Accents",
		"Export Current Accents",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Accent Themes")
}

// HandleAccentMenu processes the user's selection from the accent menu
func HandleAccentMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleAccentMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Presets":
			logging.LogDebug("Selected Accent Presets")
			app.SetSelectedAccentThemeSource(app.PresetSource)
			return app.Screens.AccentSelection

		case "Custom Accents":
			logging.LogDebug("Selected Custom Accents")
			app.SetSelectedAccentThemeSource(app.CustomSource)
			return app.Screens.AccentSelection

		case "Export Current Accents":
			logging.LogDebug("Selected Export Current Accents")
			return app.Screens.AccentExport

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.AccentMenu
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.CustomizationMenu
	}

	return app.Screens.AccentMenu
}