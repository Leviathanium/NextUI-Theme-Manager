// src/internal/ui/screens/accent_selection.go
// Implementation of the accent color selection screen

package screens

import (
	"encoding/json"
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/accents"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// AccentSelectionScreen displays available accent color themes
func AccentSelectionScreen() (string, int) {
	// Create the items array with proper structure
	var items []map[string]interface{}

	for _, theme := range accents.PredefinedThemes {
		// Convert 0xFFFFFF to #FFFFFF format
		color := "#" + strings.TrimPrefix(theme.Color2, "0x")

		// Create item with name and color option
		item := map[string]interface{}{
			"name": theme.Name,
			"options": []string{color},
			"selected": 0,
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		logging.LogDebug("No accent themes found")
		ui.ShowMessage("No accent themes available.", "3")
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

	logging.LogDebug("Displaying %d accent themes", len(items))
	return ui.DisplayMinUiList(string(jsonData), "json", "Select Accent Theme", "--item-key", "items")
}

// HandleAccentSelection processes the user's accent theme selection
func HandleAccentSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleAccentSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme (selection is just the theme name)
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