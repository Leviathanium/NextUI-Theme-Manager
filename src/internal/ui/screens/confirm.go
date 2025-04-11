// src/internal/ui/screens/confirm.go
// Implementation of the confirmation screen

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

// ConfirmScreen asks for confirmation before applying a theme
func ConfirmScreen() (string, int) {
	var message string

	switch app.GetSelectedThemeType() {
	case app.GlobalTheme:
		message = fmt.Sprintf("Apply global background '%s' to all directories?", app.GetSelectedTheme())
	case app.DynamicTheme:
		message = fmt.Sprintf("Apply dynamic theme '%s'?", app.GetSelectedTheme())
	case app.CustomTheme:
		message = fmt.Sprintf("Select theme for '%s'?", app.GetSelectedTheme())
	case app.DefaultTheme:
		if app.GetDefaultAction() == app.OverwriteAction {
			message = "Apply default black theme to all directories?"
		} else {
			message = "Delete all background images from all directories?"
		}
	}

	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleConfirmScreen processes the user's confirmation choice
func HandleConfirmScreen(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleConfirmScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Apply the selected theme
			logging.LogDebug("User confirmed, applying theme")
			applyTheme()
		} else {
			logging.LogDebug("User selected No, returning to previous screen")
			if app.GetSelectedThemeType() == app.DefaultTheme {
				return app.Screens.DefaultThemeOptions
			} else {
				return app.Screens.ThemeSelection
			}
		}
	case 1, 2:
		// User pressed cancel or back
		logging.LogDebug("User cancelled, returning to previous screen")
		if app.GetSelectedThemeType() == app.DefaultTheme {
			return app.Screens.DefaultThemeOptions
		} else {
			return app.Screens.ThemeSelection
		}
	}

	return app.Screens.MainMenu
}

// applyTheme applies the selected theme
func applyTheme() {
	var message string
	var err error

	switch app.GetSelectedThemeType() {
	case app.GlobalTheme:
		// Apply global theme to all directories
		logging.LogDebug("Applying global theme: %s", app.GetSelectedTheme())
		err = themes.ApplyGlobalTheme(app.GetSelectedTheme())
		if err != nil {
			logging.LogDebug("Error applying global theme: %v", err)
			message = fmt.Sprintf("Error: %s", err)
		} else {
			message = fmt.Sprintf("Applied global background: %s", app.GetSelectedTheme())
		}

	case app.DynamicTheme:
		// Skip if "No themes found" is selected
		if app.GetSelectedTheme() == "No themes found" {
			logging.LogDebug("No theme selected")
			message = "No theme selected"
			ui.ShowMessage(message, "3")
			return
		}

		// Apply dynamic theme pack
		logging.LogDebug("Applying dynamic theme: %s", app.GetSelectedTheme())
		err = themes.ApplyDynamicTheme(app.GetSelectedTheme())
		if err != nil {
			logging.LogDebug("Error applying dynamic theme: %v", err)
			message = fmt.Sprintf("Error: %s", err)
		} else {
			message = fmt.Sprintf("Applied dynamic theme: %s", app.GetSelectedTheme())
		}

	case app.CustomTheme:
		// For custom themes, we need to show a theme selection menu first
		// since the themes package can't display UI due to circular dependencies
		cwd, err := os.Getwd()
		if err != nil {
			message = "Error getting current directory"
			ui.ShowMessage(message, "3")
			return
		}

		// Get available themes
		globalThemesPath := filepath.Join(cwd, "Themes", "Global")
		themesList, err := themes.ListGlobalThemes(globalThemesPath)
		if err != nil {
			message = fmt.Sprintf("Error: %s", err)
			ui.ShowMessage(message, "3")
			return
		}

		if len(themesList) == 0 {
			message = "No themes found in Global directory"
			ui.ShowMessage(message, "3")
			return
		}

		// Show theme selection menu
		logging.LogDebug("Displaying theme selection menu for %s", app.GetSelectedTheme())
		themeName, exitCode := ui.DisplayMinUiList(
			strings.Join(themesList, "\n"),
			"text",
			fmt.Sprintf("Select Theme for %s", app.GetSelectedTheme()),
		)

		if exitCode != 0 || themeName == "" {
			message = "Theme selection cancelled"
			ui.ShowMessage(message, "3")
			return
		}

		// Set the selected theme as an environment variable for the theme package to use
		os.Setenv("SELECTED_THEME", themeName)

		// Apply custom theme to specific system
		logging.LogDebug("Applying custom theme to: %s", app.GetSelectedTheme())
		err = themes.ApplyCustomTheme(app.GetSelectedTheme())
		if err != nil {
			logging.LogDebug("Error applying custom theme: %v", err)
			message = fmt.Sprintf("Error: %s", err)
		} else {
			message = fmt.Sprintf("Applied theme to: %s", app.GetSelectedTheme())
		}

	case app.DefaultTheme:
		// Apply default theme based on selected action
		if app.GetDefaultAction() == app.OverwriteAction {
			logging.LogDebug("Applying default theme - overwriting backgrounds")
			err = themes.OverwriteWithDefaultTheme()
			if err != nil {
				logging.LogDebug("Error applying default theme: %v", err)
				message = fmt.Sprintf("Error: %s", err)
			} else {
				message = "Applied default theme"
			}
		} else {
			logging.LogDebug("Applying default theme - deleting backgrounds")
			err = themes.DeleteAllBackgrounds()
			if err != nil {
				logging.LogDebug("Error deleting backgrounds: %v", err)
				message = fmt.Sprintf("Error: %s", err)
			} else {
				message = "Deleted all backgrounds"
			}
		}
	}

	ui.ShowMessage(message, "3")
}