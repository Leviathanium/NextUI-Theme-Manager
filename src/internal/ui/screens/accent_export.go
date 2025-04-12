// src/internal/ui/screens/accent_export.go
// Implementation of the accent theme export screen

package screens

import (
	"fmt"

	"nextui-themes/internal/app"
	"nextui-themes/internal/accents"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// AccentExportScreen displays the accent theme export screen
func AccentExportScreen() (string, int) {
	// Prompt for a theme name
	return ui.DisplayMinUiList("Enter a name for the theme", "text", "Export Current Theme")
}

// HandleAccentExport processes the user's accent theme export
func HandleAccentExport(themeName string, exitCode int) app.Screen {
	logging.LogDebug("HandleAccentExport called with theme name: '%s', exitCode: %d", themeName, exitCode)

	switch exitCode {
	case 0:
		if themeName != "" {
			// Create a filename with .txt extension
			fileName := themeName + ".txt"

			// Save the current theme to a file in the custom themes directory
			logging.LogDebug("Exporting current theme as: %s", themeName)
			err := accents.SaveThemeToFile(&accents.CurrentTheme, fileName, true) // true = custom
			if err != nil {
				logging.LogDebug("Error exporting theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error exporting theme: %s", err), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("Theme exported as: %s", themeName), "3")

				// Refresh the themes list
				err = accents.LoadExternalAccentThemes()
				if err != nil {
					logging.LogDebug("Error refreshing themes: %v", err)
				}
			}
		} else {
			ui.ShowMessage("Export cancelled: No name provided", "3")
		}
		return app.Screens.AccentMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.AccentMenu
	}

	return app.Screens.AccentMenu
}