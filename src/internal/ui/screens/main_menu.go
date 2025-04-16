// src/internal/ui/screens/main_menu.go
// Implementation of the main menu screen - simplified for theme management only

package screens

import (
	"os"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// MainMenuScreen displays the main menu with only theme options
func MainMenuScreen() (string, int) {
	// Menu items - simplified to just theme management
	menu := []string{
		"Import Theme",
		"Export Current Settings",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "NextUI Theme Manager", "--cancel-text", "QUIT")
}

// HandleMainMenu processes the user's selection from the main menu
func HandleMainMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleMainMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Import Theme":
			logging.LogDebug("Selected Import Theme")
			return app.Screens.ThemeImport

		case "Export Current Settings":
			logging.LogDebug("Selected Export Current Settings")
			return app.Screens.ThemeExport

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