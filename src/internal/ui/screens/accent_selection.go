// src/internal/ui/screens/accent_selection.go
// Implementation of the accent color selection screen

package screens

import (
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/accents"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// AccentSelectionScreen displays available accent color themes
func AccentSelectionScreen() (string, int) {
	// Get list of available accent themes
	var themesList []string
	for _, theme := range accents.PredefinedThemes {
		themesList = append(themesList, theme.Name)
	}

	if len(themesList) == 0 {
		logging.LogDebug("No accent themes found")
		ui.ShowMessage("No accent themes available.", "3")
		return "", 1
	}

	logging.LogDebug("Displaying %d accent themes", len(themesList))
	return ui.DisplayMinUiList(strings.Join(themesList, "\n"), "text", "Select Accent Theme")
}

// HandleAccentSelection processes the user's accent theme selection
func HandleAccentSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleAccentSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		// Find the selected theme
		var selectedTheme *accents.ThemeColor
		for _, theme := range accents.PredefinedThemes {
			if theme.Name == selection {
				selectedTheme = &theme
				break
			}
		}

		if selectedTheme != nil {
			// Apply the selected theme
			logging.LogDebug("Applying accent theme: %s", selectedTheme.Name)
			err := accents.ApplyThemeColors(selectedTheme)
			if err != nil {
				logging.LogDebug("Error applying accent theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("Applied accent theme: %s", selectedTheme.Name), "3")
			}
		} else {
			logging.LogDebug("Selected theme not found: %s", selection)
			ui.ShowMessage("Selected theme not found.", "3")
		}

		return app.Screens.MainMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.AccentSelection
}