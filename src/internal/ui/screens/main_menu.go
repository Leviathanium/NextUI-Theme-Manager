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

func MainMenuScreen() (string, int) {
	// Updated menu items with "Deconstruct" added
	menu := []string{
		"Installed Themes",
		"Download Themes",
		"Sync Catalog",
		"Components",
		"Deconstruct", // Added the Deconstruct option to main menu (without ellipsis)
		"Export",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "NextUI Theme Manager", "--cancel-text", "QUIT")
}

func HandleMainMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleMainMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Installed Themes":
			logging.LogDebug("Selected Installed Themes")
			return app.Screens.InstalledThemes

		case "Download Themes":
			logging.LogDebug("Selected Download Themes")
			return app.Screens.DownloadThemes

		case "Sync Catalog":
			logging.LogDebug("Selected Sync Catalog")
			return app.Screens.SyncCatalog

		case "Components":
			logging.LogDebug("Selected Components")
			return app.Screens.ComponentsMenu

		case "Deconstruct": // Add handling for the new main menu option
			logging.LogDebug("Selected Deconstruct")
			return app.Screens.Deconstruction

		case "Export":
			logging.LogDebug("Selected Export")
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
