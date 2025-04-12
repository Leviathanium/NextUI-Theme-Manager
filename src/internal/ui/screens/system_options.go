// src/internal/ui/screens/system_options.go
// Implementation of the system options menu screen

package screens

import (
	"fmt"
	"strings"
    "nextui-themes/internal/icons"
	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"nextui-themes/internal/ui"
)

// SystemOptionsMenuScreen displays the system selection for system-specific options
func SystemOptionsMenuScreen() (string, int) {
	// Get system paths to find all installed systems
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error detecting systems: %s", err), "3")
		return "", 1
	}

	// Build the list of systems
	var themesList []string

	// Add standard menu items
	themesList = append(themesList, "Root")
	themesList = append(themesList, "Recently Played")
	themesList = append(themesList, "Tools")

	// Add all detected rom systems
	for _, system := range systemPaths.Systems {
		themesList = append(themesList, system.Name)
	}

	if len(themesList) == 0 {
		logging.LogDebug("No systems found")
		ui.ShowMessage("No systems found!", "3")
		themesList = []string{"No systems found"}
	}

	return ui.DisplayMinUiList(strings.Join(themesList, "\n"), "text", "Select System")
}

// HandleSystemOptionsMenu processes the user's system selection
func HandleSystemOptionsMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleSystemOptionsMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a system
		app.SetSelectedSystem(selection)
		return app.Screens.SystemOptionsForSelectedSystem
	case 1, 2:
		// User pressed cancel or back
		return app.Screens.CustomizationMenu
	}

	return app.Screens.SystemOptionsMenu
}

// SystemOptionsForSelectedSystemScreen displays options for the selected system
func SystemOptionsForSelectedSystemScreen() (string, int) {
	selectedSystem := app.GetSelectedSystem()

	// Menu items for the selected system
	menu := []string{
		"Wallpaper",
		"Icon",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", fmt.Sprintf("Options for %s", selectedSystem))
}

// HandleSystemOptionsForSelectedSystem processes the user's selection for a specific system
func HandleSystemOptionsForSelectedSystem(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleSystemOptionsForSelectedSystem called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Wallpaper":
			logging.LogDebug("Selected Wallpaper for system: %s", app.GetSelectedSystem())
			app.SetSelectedThemeType(app.CustomTheme)
			return app.Screens.WallpaperSelection

		case "Icon":
			logging.LogDebug("Selected Icon for system: %s", app.GetSelectedSystem())
			return app.Screens.SystemIconSelection

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.SystemOptionsForSelectedSystem
		}

	case 1, 2:
		// User pressed cancel or back
		app.SetSelectedSystem("") // Clear the selected system
		return app.Screens.SystemOptionsMenu
	}

	return app.Screens.SystemOptionsForSelectedSystem
}

// SystemIconSelectionScreen displays available icon packs for a specific system
func SystemIconSelectionScreen() (string, int) {
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

	logging.LogDebug("Displaying icon selection with %d options for system %s",
	                  len(iconsList), app.GetSelectedSystem())
	return ui.DisplayMinUiList(strings.Join(iconsList, "\n"), "text",
	                          fmt.Sprintf("Select Icon Pack for %s", app.GetSelectedSystem()))
}

// HandleSystemIconSelection processes the user's icon pack selection for a specific system
func HandleSystemIconSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleSystemIconSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an icon pack - proceed to confirmation
		app.SetSelectedIconPack(selection)
		return app.Screens.IconConfirm

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.SystemOptionsForSelectedSystem
	}

	return app.Screens.SystemIconSelection
}