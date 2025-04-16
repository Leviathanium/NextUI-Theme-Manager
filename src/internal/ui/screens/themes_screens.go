// src/internal/ui/screens/themes_screens.go
// Implementation of theme import/export screens - simplified version

package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/themes"
	"nextui-themes/internal/ui"
)

// ThemeImportScreen displays available themes from the Imports directory
func ThemeImportScreen() (string, int) {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Path to Themes/Imports directory
	importsDir := filepath.Join(cwd, "Themes", "Imports")

	// Ensure directory exists
	if err := os.MkdirAll(importsDir, 0755); err != nil {
		logging.LogDebug("Error creating imports directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// List available themes
	entries, err := os.ReadDir(importsDir)
	if err != nil {
		logging.LogDebug("Error reading imports directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Filter for theme directories
	var themesList []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".theme") {
			themesList = append(themesList, entry.Name())
		}
	}

	if len(themesList) == 0 {
		logging.LogDebug("No themes found")
		ui.ShowMessage("No themes found in Imports directory", "3")
		return "", 1
	}

	logging.LogDebug("Found %d themes", len(themesList))
	return ui.DisplayMinUiList(strings.Join(themesList, "\n"), "text", "Select Theme to Import")
}

// HandleThemeImport processes the user's theme selection
func HandleThemeImport(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeImport called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		app.SetSelectedTheme(selection)
		return app.Screens.ThemeImportConfirm

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.ThemeImport
}

// ThemeImportConfirmScreen displays a confirmation dialog for theme import
func ThemeImportConfirmScreen() (string, int) {
	themeName := app.GetSelectedTheme()
	message := fmt.Sprintf("Apply theme '%s'?", themeName)

	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleThemeImportConfirm processes the user's confirmation for theme import
func HandleThemeImportConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeImportConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Import the selected theme
			themeName := app.GetSelectedTheme()
			if err := themes.ImportTheme(themeName); err != nil {
				logging.LogDebug("Error importing theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("Theme '%s' imported successfully!", themeName), "3")
			}
		}
		// Return to main menu
		return app.Screens.MainMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ThemeImport
	}

	return app.Screens.ThemeImportConfirm
}

// ThemeExportScreen displays the theme export confirmation
func ThemeExportScreen() (string, int) {
	// Simple confirmation message
	message := "Export current theme settings?\nThis will create a theme package in Themes/Exports."
	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleThemeExport processes the user's choice to export a theme
func HandleThemeExport(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeExport called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Perform theme export
			if err := themes.ExportTheme(); err != nil {
				logging.LogDebug("Error exporting theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage("Theme exported successfully!", "3")
			}
		}
		// Return to main menu
		return app.Screens.MainMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.ThemeExport
}