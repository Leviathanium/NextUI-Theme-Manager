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
	var title string

	// Choose which themes to display based on the selected source
	if app.GetSelectedLEDThemeSource() == app.PresetSource {
		title = "Preset LED Themes"
		for _, theme := range leds.PresetLEDThemes {
			themesList = append(themesList, theme.Name)
		}
	} else {
		title = "Custom LED Themes"
		for _, theme := range leds.CustomLEDThemes {
			themesList = append(themesList, theme.Name)
		}
	}

	if len(themesList) == 0 {
		logging.LogDebug("No LED themes found")
		ui.ShowMessage("No LED themes available in this category.", "3")
		return "", 1
	}

	logging.LogDebug("Displaying %d LED themes", len(themesList))
	return ui.DisplayMinUiList(strings.Join(themesList, "\n"), "text", title)
}

// HandleLEDSelection processes the user's LED theme selection
func HandleLEDSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleLEDSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme - update in-memory state AND apply immediately
		app.SetSelectedLEDTheme(selection)
		logging.LogDebug("Selected LED theme: %s", selection)

		// Update current LED settings in memory
		if err := leds.UpdateCurrentLEDTheme(selection); err != nil {
			logging.LogDebug("Error updating LED theme in memory: %v", err)
			ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			return app.Screens.LEDSelection
		}

		// Apply the LED theme immediately
		logging.LogDebug("Applying LED theme immediately")
		if err := leds.ApplyCurrentLEDSettings(); err != nil {
			logging.LogDebug("Error applying LED settings: %v", err)
			ui.ShowMessage(fmt.Sprintf("Error applying LED theme: %s", err), "3")
		} else {
			ui.ShowMessage(fmt.Sprintf("LED theme '%s' applied successfully!", selection), "3")
		}

		// Return to LED menu
		return app.Screens.LEDMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.LEDMenu
	}

	return app.Screens.LEDSelection
}