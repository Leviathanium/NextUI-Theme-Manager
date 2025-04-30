// internal/ui/screens/themes_menu.go
package screens

import (
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/ui"
)

// ThemesMenuScreen displays the themes submenu screen
func ThemesMenuScreen() (string, int) {
	logging.LogDebug("Showing themes menu screen")

	menuItems := []string{
		"Installed Themes",
		"Download Themes",
	}

	return ui.DisplayMinUiList(
		strings.Join(menuItems, "\n"),
		"text",
		"Themes Menu",
	)
}

// HandleThemesMenu processes themes menu selection
func HandleThemesMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemesMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Installed Themes":
			logging.LogDebug("Selected Installed Themes")
			return app.Screens.InstalledThemes

		case "Download Themes":
			logging.LogDebug("Selected Download Themes")
			return app.Screens.DownloadThemes

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.ThemesMenu
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/back
		return app.Screens.MainMenu
	}

	return app.Screens.ThemesMenu
}

// InstalledThemesScreen displays the installed themes screen
func InstalledThemesScreen() (string, int) {
	logging.LogDebug("Showing installed themes screen")
	return themes.ShowInstalledThemes()
}

// HandleInstalledThemes processes installed themes selection
func HandleInstalledThemes(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleInstalledThemes called with selection: '%s', exitCode: %d", selection, exitCode)

	// When a theme is selected, we'll set it as the current selection and proceed to application
	if exitCode == 0 && selection != "" {
		app.SetSelectedItem(selection)
		return app.Screens.ThemeApplyConfirm
	}

	// User pressed back
	return app.Screens.ThemesMenu
}