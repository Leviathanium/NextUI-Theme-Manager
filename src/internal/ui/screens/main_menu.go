// src/internal/ui/screens/main_menu.go
// Implementation of the main menu screen

package screens

import (
	"os"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// Update the MainMenuScreen function to include "Themes" as an option
func MainMenuScreen() (string, int) {
	// Menu items without numbers
	menu := []string{
		"Themes",           // Keep this first
		"Customization",
		"Reset",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "NextUI Theme Selector", "--cancel-text", "QUIT")
}

// Update the HandleMainMenu function to handle the "Themes" selection
func HandleMainMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("handleMainMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Themes":
			logging.LogDebug("Selected Themes")
			return app.Screens.ThemesMenu

		case "Customization":
			logging.LogDebug("Selected Customization")
			return app.Screens.CustomizationMenu

		case "Reset":
			logging.LogDebug("Selected Reset")
			return app.Screens.ResetMenu

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.MainMenu
		}

	case 1, 2:
		// User pressed cancel or back
		logging.LogDebug("User cancelled/exited")
		os.Exit(0)
	}

	return app.Screens.MainMenu
}