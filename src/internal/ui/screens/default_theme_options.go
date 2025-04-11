// src/internal/ui/screens/default_theme_options.go
// Implementation of the default theme options screen

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// DefaultThemeOptionsScreen displays options for the Default Theme
func DefaultThemeOptionsScreen() (string, int) {
	options := []string{
		"Overwrite all backgrounds with black",
		"Delete all backgrounds",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", "Reset Options") // Updated title
}

// HandleDefaultThemeOptions processes user selection for default theme options
func HandleDefaultThemeOptions(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleDefaultThemeOptions called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Overwrite all backgrounds with black" {
			logging.LogDebug("Selected to overwrite backgrounds with black")
			app.SetDefaultAction(app.OverwriteAction)
			return app.Screens.ConfirmScreen
		} else if selection == "Delete all backgrounds" {
			logging.LogDebug("Selected to delete all backgrounds")
			app.SetDefaultAction(app.DeleteAction)
			return app.Screens.ConfirmScreen
		}
	case 1, 2:
		return app.Screens.MainMenu
	}

	return app.Screens.DefaultThemeOptions
}