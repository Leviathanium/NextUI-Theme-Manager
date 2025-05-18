// src/internal/ui/screens/apply.go
package screens

import (
	"fmt"
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// ShowApplyThemeScreen displays the theme selection screen for applying
func ShowApplyThemeScreen() (string, int) {
	app.LogDebug("Showing apply theme screen")

	// Get list of available themes
	themeNames, err := themes.ListThemes()
	if err != nil {
		app.LogDebug("Error listing themes: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error listing themes: %s", err), "3")
		return "", 2  // Exit code 2 signals a controlled return to main menu
	}

	// Check if we have any themes
	if len(themeNames) == 0 {
		app.LogDebug("No themes available")
		// Show message and wait for user acknowledgment
		ui.ShowMessage("No themes available. Please download or import themes first.", "0")
		return "", 2  // Return to main menu with proper exit code
	}

	// Build menu items
	var menuItems []string
	for _, themeName := range themeNames {
		// Try to get author from manifest
		themePath := themes.GetThemePath(themeName)
		manifest, err := themes.ReadManifest(themePath)

		if err == nil && manifest.Author != "" {
			menuItems = append(menuItems, fmt.Sprintf("%s by %s", themeName, manifest.Author))
		} else {
			menuItems = append(menuItems, themeName)
		}
	}

	return ui.ShowMenu(
		strings.Join(menuItems, "\n"),
		"Select Theme to Apply",
		"--cancel-text", "BACK",
	)
}

// HandleApplyThemeScreen processes the theme selection
func HandleApplyThemeScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleApplyThemeScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected a theme
		app.SetSelectedItem(selection)
		return app.ScreenApplyThemeConfirm
	} else if exitCode == 1 || exitCode == 2 {
		// User cancelled
		return app.ScreenMainMenu
	}

	return app.ScreenApplyTheme
}

// ShowApplyThemeConfirmScreen displays the confirmation screen for applying a theme
func ShowApplyThemeConfirmScreen() (string, int) {
	app.LogDebug("Showing apply theme confirmation screen")

	selectedTheme := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Apply theme '" + selectedTheme + "'?")
}

// HandleApplyThemeConfirmScreen processes the confirmation result
func HandleApplyThemeConfirmScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleApplyThemeConfirmScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed - proceed to applying
		return app.ScreenApplyingTheme
	} else {
		// User cancelled
		return app.ScreenApplyTheme
	}
}

// ShowApplyingThemeScreen displays the theme applying progress screen
func ShowApplyingThemeScreen() (string, int) {
	app.LogDebug("Showing applying theme screen")

	selectedTheme := app.GetSelectedItem()
	return ui.ShowMessage("Applying theme '" + selectedTheme + "'...", "2")
}

// HandleApplyingThemeScreen processes the applying operation
func HandleApplyingThemeScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleApplyingThemeScreen called with exitCode: %d", exitCode)

	selectedTheme := app.GetSelectedItem()

	// Extract theme name from selection (remove "by Author" part if present)
	themeName := selectedTheme
	if idx := strings.Index(themeName, " by "); idx > 0 {
		themeName = themeName[:idx]
	}

	// Apply the theme
	err := themes.ApplyTheme(themeName)

	if err != nil {
		app.LogDebug("Error applying theme: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error applying theme: %s", err), "3")
	}

	return app.ScreenThemeApplied
}

// ShowThemeAppliedScreen displays the theme applied success screen
func ShowThemeAppliedScreen() (string, int) {
	app.LogDebug("Showing theme applied screen")

	selectedTheme := app.GetSelectedItem()
	return ui.ShowMessage("Theme '" + selectedTheme + "' applied successfully!", "2")
}

// HandleThemeAppliedScreen processes the success screen
func HandleThemeAppliedScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleThemeAppliedScreen called with exitCode: %d", exitCode)

	// Return to main menu after showing success message
	return app.ScreenMainMenu
}