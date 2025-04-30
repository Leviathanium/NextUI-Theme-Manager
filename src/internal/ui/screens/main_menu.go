// internal/ui/screens/main_menu.go
package screens

import (
	"os"
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/ui"
)

// MainMenuScreen displays the main menu screen
func MainMenuScreen() (string, int) {
	logging.LogDebug("Showing main menu")

	menuItems := []string{
		"Themes",
		"Overlays",
		"Sync Catalog",
		"Backup",
		"Revert",
		"Purge",
	}

	return ui.DisplayMinUiList(
		strings.Join(menuItems, "\n"),
		"text",
		"Theme Manager",
		"--cancel-text", "QUIT",
	)
}

// HandleMainMenu processes main menu selection
func HandleMainMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleMainMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Themes":
			logging.LogDebug("Selected Themes")
			return app.Screens.Themes

		case "Overlays":
			logging.LogDebug("Selected Overlays")
			return app.Screens.Overlays

		case "Sync Catalog":
			logging.LogDebug("Selected Sync Catalog")
			return app.Screens.SyncCatalog

		case "Backup":
			logging.LogDebug("Selected Backup")
			return app.Screens.Backup

		case "Revert":
			logging.LogDebug("Selected Revert")
			return app.Screens.Revert

		case "Purge":
			logging.LogDebug("Selected Purge")
			return app.Screens.Purge

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.MainMenu
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel or exit
		logging.LogDebug("User quit the application")
		os.Exit(0)
	}

	return app.Screens.MainMenu
}