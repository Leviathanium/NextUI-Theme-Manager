// src/internal/ui/screens/led_selection.go
// Implementation of the LED settings selection screen

package screens

import (
	"encoding/json"
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/leds"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// LEDSelectionScreen displays available LED themes
func LEDSelectionScreen() (string, int) {
	// Create the items array with proper structure
	var items []map[string]interface{}

	for _, theme := range leds.PredefinedThemes {
		// Convert 0xFFFFFF to #FFFFFF format
		color := "#" + strings.TrimPrefix(theme.Color, "0x")

		// Create item with name and color option
		item := map[string]interface{}{
			"name": theme.Name,
			"options": []string{color},
			"selected": 0,
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		logging.LogDebug("No LED themes found")
		ui.ShowMessage("No LED themes available.", "3")
		return "", 1
	}

	// Create a wrapper object with the items array
	wrapper := map[string]interface{}{
		"items": items,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(wrapper)
	if err != nil {
		logging.LogDebug("Error creating JSON: %v", err)
		ui.ShowMessage("Error creating theme list.", "3")
		return "", 1
	}

	logging.LogDebug("Displaying %d LED themes", len(items))
	return ui.DisplayMinUiList(string(jsonData), "json", "Select LED Theme", "--item-key", "items")
}

// HandleLEDSelection processes the user's LED theme selection
func HandleLEDSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleLEDSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme (selection is just the theme name)
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