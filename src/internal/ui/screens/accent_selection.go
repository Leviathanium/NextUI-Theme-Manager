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

	// Add "Apply Changes" option
	themesList = append(themesList, "Apply Changes")

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
		// Check if "Apply Changes" was selected
		if selection == "Apply Changes" {
			// Apply the current theme settings to the system
			logging.LogDebug("Applying current accent settings")
			err := accents.ApplyCurrentTheme()
			if err != nil {
				logging.LogDebug("Error applying accent theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage("Accent settings applied successfully!", "3")
			}
			return app.Screens.MainMenu
		}

		// User selected a theme - update in-memory state but don't apply yet
		app.SetSelectedAccentTheme(selection)
		logging.LogDebug("Selected accent theme: %s", selection)

		// Update the current theme in memory
		if err := accents.UpdateCurrentTheme(selection); err != nil {
			logging.LogDebug("Error updating accent theme in memory: %v", err)
			ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		} else {
			ui.ShowMessage(fmt.Sprintf("Selected theme: %s\nChoose 'Apply Changes' to save", selection), "3")
		}

		// Return to accent selection screen
		return app.Screens.AccentSelection

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.AccentSelection
}