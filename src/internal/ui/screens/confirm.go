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
		// This is for system-specific wallpaper application
		systemName := app.GetSelectedSystem()
		if systemName != "" {
			message = fmt.Sprintf("Apply wallpaper '%s' to %s?", app.GetSelectedTheme(), systemName)
		} else {
			message = fmt.Sprintf("Select theme for '%s'?", app.GetSelectedTheme())
		}
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
				return app.Screens.ResetMenu
			} else if app.GetSelectedThemeType() == app.GlobalTheme {
				return app.Screens.GlobalOptionsMenu
			} else if app.GetSelectedThemeType() == app.CustomTheme && app.GetSelectedSystem() != "" {
				return app.Screens.SystemOptionsForSelectedSystem
			} else {
				return app.Screens.ThemeSelection
			}
		}
	case 1, 2:
		// User pressed cancel or back
		logging.LogDebug("User cancelled, returning to previous screen")
		if app.GetSelectedThemeType() == app.DefaultTheme {
			return app.Screens.ResetMenu
		} else if app.GetSelectedThemeType() == app.GlobalTheme {
			return app.Screens.GlobalOptionsMenu
		} else if app.GetSelectedThemeType() == app.CustomTheme && app.GetSelectedSystem() != "" {
			return app.Screens.SystemOptionsForSelectedSystem
		} else {
			return app.Screens.ThemeSelection
		}
	}

	// Default return to main menu
	if app.GetSelectedThemeType() == app.GlobalTheme {
		return app.Screens.GlobalOptionsMenu
	} else if app.GetSelectedThemeType() == app.CustomTheme && app.GetSelectedSystem() != "" {
		return app.Screens.SystemOptionsForSelectedSystem
	} else {
		return app.Screens.MainMenu
	}
}

// WallpaperConfirmScreen asks for confirmation before applying a wallpaper
func WallpaperConfirmScreen() (string, int) {
	var message string

	// Check if we're applying to a specific system
	if app.GetSelectedSystem() != "" {
		message = fmt.Sprintf("Apply wallpaper '%s' to %s?",
		                     app.GetSelectedTheme(), app.GetSelectedSystem())
	} else {
		message = fmt.Sprintf("Apply wallpaper '%s' to all directories?",
		                     app.GetSelectedTheme())
	}

	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleWallpaperConfirm processes the user's confirmation for a wallpaper
func HandleWallpaperConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleWallpaperConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Apply the wallpaper
			if app.GetSelectedSystem() != "" {
				// System-specific wallpaper
				applySystemWallpaper()
			} else {
				// Global wallpaper
				applyGlobalWallpaper()
			}
		}

		// Return to appropriate screen based on context
		if app.GetSelectedSystem() != "" {
			return app.Screens.SystemOptionsForSelectedSystem
		} else {
			return app.Screens.GlobalOptionsMenu
		}
	case 1, 2:
		// User pressed cancel or back
		if app.GetSelectedSystem() != "" {
			return app.Screens.SystemOptionsForSelectedSystem
		} else {
			return app.Screens.GlobalOptionsMenu
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

		// Get available themes from Wallpapers directory
		wallpapersDir := filepath.Join(cwd, "Wallpapers")
		themesList, err := themes.ListGlobalThemes(wallpapersDir)
		if err != nil {
			message = fmt.Sprintf("Error: %s", err)
			ui.ShowMessage(message, "3")
			return
		}

		if len(themesList) == 0 {
			message = "No themes found in Wallpapers directory"
			ui.ShowMessage(message, "3")
			return
		}

		// Apply to specific system if one is selected
		systemName := app.GetSelectedSystem()
		if systemName != "" {
			// Apply custom theme to specific system
			logging.LogDebug("Applying custom theme to: %s", systemName)
			err = themes.ApplyCustomTheme(systemName, app.GetSelectedTheme())
			if err != nil {
				logging.LogDebug("Error applying custom theme: %v", err)
				message = fmt.Sprintf("Error: %s", err)
			} else {
				message = fmt.Sprintf("Applied theme to: %s", systemName)
			}
		} else {
			// No system selected - display error
			logging.LogDebug("No system selected for custom theme")
			message = "No system selected for custom theme"
		}

	case app.DefaultTheme:
		// Only delete option now, as overwrite option has been removed
		logging.LogDebug("Applying default theme - deleting backgrounds")
		err = themes.DeleteAllBackgrounds()
		if err != nil {
			logging.LogDebug("Error deleting backgrounds: %v", err)
			message = fmt.Sprintf("Error: %s", err)
		} else {
			message = "Deleted all backgrounds"
		}
	}

	ui.ShowMessage(message, "3")
}

// applySystemWallpaper applies a wallpaper to a specific system
func applySystemWallpaper() {
	selectedSystem := app.GetSelectedSystem()
	selectedTheme := app.GetSelectedTheme()

	logging.LogDebug("Applying wallpaper '%s' to system: %s", selectedTheme, selectedSystem)

	// Apply the custom theme to the specific system
	err := themes.ApplyCustomTheme(selectedSystem, selectedTheme)
	if err != nil {
		logging.LogDebug("Error applying wallpaper: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
	} else {
		ui.ShowMessage(fmt.Sprintf("Applied wallpaper '%s' to %s", selectedTheme, selectedSystem), "3")
	}
}

// applyGlobalWallpaper applies a wallpaper to all systems
func applyGlobalWallpaper() {
	selectedTheme := app.GetSelectedTheme()

	logging.LogDebug("Applying global wallpaper: %s", selectedTheme)

	// Apply the global theme to all directories
	err := themes.ApplyGlobalTheme(selectedTheme)
	if err != nil {
		logging.LogDebug("Error applying global wallpaper: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
	} else {
		ui.ShowMessage(fmt.Sprintf("Applied global wallpaper: %s", selectedTheme), "3")
	}
}