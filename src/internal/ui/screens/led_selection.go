// src/internal/ui/screens/led_selection.go
// Implementation of the LED settings selection screen

package screens

import (
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/leds"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// LEDSelectionScreen displays available LED themes
func LEDSelectionScreen() (string, int) {
	// Get list of available LED themes
	var themesList []string
	for _, theme := range leds.PredefinedThemes {
		themesList = append(themesList, theme.Name)
	}

	if len(themesList) == 0 {
		logging.LogDebug("No LED themes found")
		ui.ShowMessage("No LED themes available.", "3")
		return "", 1
	}

	logging.LogDebug("Displaying %d LED themes", len(themesList))
	return ui.DisplayMinUiList(strings.Join(themesList, "\n"), "text", "Select LED Theme")
}

// HandleLEDSelection processes the user's LED theme selection
func HandleLEDSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleLEDSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		// Find the selected theme
		var selectedTheme *leds.LEDTheme
		for _, theme := range leds.PredefinedThemes {
			if theme.Name == selection {
				selectedTheme = &theme
				break
			}
		}

		if selectedTheme != nil {
			// Apply the selected theme
			logging.LogDebug("Applying LED theme: %s", selectedTheme.Name)
			err := leds.ApplyLEDTheme(selectedTheme)
			if err != nil {
				logging.LogDebug("Error applying LED theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("Applied LED theme: %s", selectedTheme.Name), "3")
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

	return app.Screens.LEDSelection
}