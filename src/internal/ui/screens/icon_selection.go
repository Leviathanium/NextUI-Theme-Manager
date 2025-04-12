// src/internal/ui/screens/icon_selection.go
// Implementation of the icon pack selection screen

package screens

import (
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/icons"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// IconSelectionScreen displays available icon packs
func IconSelectionScreen() (string, int) {
	// List available icon packs
	iconsList, err := icons.ListIconPacks()
	if err != nil {
		logging.LogDebug("Error loading icon packs: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error loading icon packs: %s", err), "3")
		return "", 1
	}

	if len(iconsList) == 0 {
		logging.LogDebug("No icon packs found")
		ui.ShowMessage("No icon packs found. Add icon packs to the Icons directory.", "3")
		return "", 1
	}

	logging.LogDebug("Displaying icon selection with %d options", len(iconsList))
	return ui.DisplayMinUiList(strings.Join(iconsList, "\n"), "text", "Select Icon Pack")
}

// HandleIconSelection processes the user's icon pack selection
func HandleIconSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleIconSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an icon pack - proceed to confirmation
		app.SetSelectedIconPack(selection)
		return app.Screens.IconConfirm

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.CustomizationMenu
	}

	return app.Screens.IconSelection
}

// IconConfirmScreen asks for confirmation before applying an icon pack
func IconConfirmScreen() (string, int) {
	message := fmt.Sprintf("Apply icon pack '%s' to all systems?", app.GetSelectedIconPack())

	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleIconConfirm processes the user's confirmation choice for icon pack
func HandleIconConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleIconConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Apply the selected icon pack
			logging.LogDebug("User confirmed, applying icon pack")

			// Get the selected icon pack name
			iconPack := app.GetSelectedIconPack()

			// Apply the icon pack
			err := icons.ApplyIconPack(iconPack)
			if err != nil {
				logging.LogDebug("Error applying icon pack: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("Applied icon pack: %s", iconPack), "3")
			}
		} else {
			logging.LogDebug("User selected No, returning to icon selection screen")
			return app.Screens.IconSelection
		}
	case 1, 2:
		// User pressed cancel or back
		logging.LogDebug("User cancelled, returning to icon selection screen")
		return app.Screens.IconSelection
	}

	return app.Screens.CustomizationMenu
}