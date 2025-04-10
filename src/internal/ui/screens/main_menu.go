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

// MainMenuScreen shows the main menu with theme options
func MainMenuScreen() (string, int) {
	// Menu items without numbers
	menu := []string{
		"Global Themes",
		"Dynamic Themes",
		"Custom Themes",
		"Default Theme",
		"Fonts",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "NextUI Theme Selector", "--cancel-text", "QUIT")
}

// HandleMainMenu processes the user's selection from the main menu
func HandleMainMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("handleMainMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Global Themes":
			logging.LogDebug("Selected Global Themes")
			app.SetSelectedThemeType(app.GlobalTheme)
			return app.Screens.ThemeSelection

		case "Dynamic Themes":
			logging.LogDebug("Selected Dynamic Themes")
			app.SetSelectedThemeType(app.DynamicTheme)
			return app.Screens.ThemeSelection

		case "Custom Themes":
			logging.LogDebug("Selected Custom Themes")
			app.SetSelectedThemeType(app.CustomTheme)
			return app.Screens.ThemeSelection

		case "Default Theme":
			logging.LogDebug("Selected Default Theme")
			app.SetSelectedThemeType(app.DefaultTheme)
			return app.Screens.DefaultThemeOptions

		case "Fonts":
			logging.LogDebug("Selected Fonts")
			return app.Screens.FontSelection

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