// src/internal/ui/screens/reset_menu.go
// Implementation of the reset menu screen

package screens

import (
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// ResetMenuScreen displays options for resetting various aspects of the UI
func ResetMenuScreen() (string, int) {
	options := []string{
		"Overwrite all backgrounds with black",
		"Delete all backgrounds",
		"Delete all icons",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", "Reset Options")
}

// HandleResetMenu processes user selection for reset options
func HandleResetMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleResetMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Overwrite all backgrounds with black" {
			logging.LogDebug("Selected to overwrite backgrounds with black")
			app.SetSelectedThemeType(app.DefaultTheme)
			app.SetDefaultAction(app.OverwriteAction)
			return app.Screens.ConfirmScreen
		} else if selection == "Delete all backgrounds" {
			logging.LogDebug("Selected to delete all backgrounds")
			app.SetSelectedThemeType(app.DefaultTheme)
			app.SetDefaultAction(app.DeleteAction)
			return app.Screens.ConfirmScreen
		} else if selection == "Delete all icons" {
			logging.LogDebug("Selected to delete all icons")
			return app.Screens.ClearIconsConfirm
		}
	case 1, 2:
		return app.Screens.MainMenu
	}

	return app.Screens.ResetMenu
}