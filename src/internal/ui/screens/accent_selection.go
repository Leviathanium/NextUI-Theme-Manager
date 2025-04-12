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
	var title string

	// Choose which themes to display based on the selected source
	if app.GetSelectedAccentThemeSource() == app.PresetSource {
		title = "Preset Accent Themes"
		for _, theme := range accents.PresetThemes {
			themesList = append(themesList, theme.Name)
		}
	} else {
		title = "Custom Accent Themes"
		for _, theme := range accents.CustomThemes {
			themesList = append(themesList, theme.Name)
		}
	}

	if len(themesList) == 0 {
		logging.LogDebug("No accent themes found")
		ui.ShowMessage("No accent themes available in this category.", "3")
		return "", 1
	}

	logging.LogDebug("Displaying %d accent themes", len(themesList))
	return ui.DisplayMinUiList(strings.Join(themesList, "\n"), "text", title)
}

// HandleAccentSelection processes the user's accent theme selection
func HandleAccentSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleAccentSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme - update in-memory state AND apply immediately
		app.SetSelectedAccentTheme(selection)
		logging.LogDebug("Selected accent theme: %s", selection)

		// Update the current theme in memory
		if err := accents.UpdateCurrentTheme(selection); err != nil {
			logging.LogDebug("Error updating accent theme in memory: %v", err)
			ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			return app.Screens.AccentSelection
		}

		// Apply the theme immediately
		logging.LogDebug("Applying accent theme immediately")
		if err := accents.ApplyCurrentTheme(); err != nil {
			logging.LogDebug("Error applying accent theme: %v", err)
			ui.ShowMessage(fmt.Sprintf("Error applying theme: %s", err), "3")
		} else {
			ui.ShowMessage(fmt.Sprintf("Theme '%s' applied successfully!", selection), "3")
		}

		// Return to accent menu
		return app.Screens.AccentMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.AccentMenu
	}

	return app.Screens.AccentSelection
}