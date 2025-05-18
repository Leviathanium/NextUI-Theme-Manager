// internal/ui/screens/main_menu.go
package screens

import (
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/ui"
)

// ShowMainMenu displays the main menu screen
func ShowMainMenu() (string, int) {
	app.LogDebug("Showing main menu")

	menuItems := []string{
		"Apply Theme",
		"Download Theme",
		"Sync Catalog",
		"Backups",
		"About",
	}

	return ui.ShowMenu(
		strings.Join(menuItems, "\n"),
		"Theme Manager",
		"--cancel-text", "QUIT",
	)
}

// HandleMainMenu processes main menu selection
func HandleMainMenu(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleMainMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Apply Theme":
			app.LogDebug("Selected Apply Theme")
			return app.ScreenApplyTheme

		case "Download Theme":
			app.LogDebug("Selected Download Theme")
			return app.ScreenDownloadTheme

		case "Sync Catalog":
			app.LogDebug("Selected Sync Catalog")
			return app.ScreenSyncCatalog

		case "Backups":
			app.LogDebug("Selected Backups")
			return app.ScreenBackupsMenu

		case "About":
			app.LogDebug("Selected About")
			return app.ScreenAbout

		default:
			app.LogDebug("Unknown selection: %s", selection)
			return app.ScreenMainMenu
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/back/exit
		// Exit the application
		app.LogDebug("User exited the application")
		return app.ScreenMainMenu
	}

	return app.ScreenMainMenu
}

// ShowAboutScreen displays the about screen
func ShowAboutScreen() (string, int) {
	app.LogDebug("Showing about screen")

	// Display about message
	return ui.ShowMessage(
		"Theme Manager v1.0\n\nA simple application to manage themes for your device.",
		"3",
	)
}

// HandleAboutScreen processes the about screen
func HandleAboutScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleAboutScreen called with exitCode: %d", exitCode)

	// Return to main menu after displaying about screen
	return app.ScreenMainMenu
}