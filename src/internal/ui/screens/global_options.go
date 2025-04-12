// src/internal/ui/screens/global_options.go
// Implementation of the global options menu screen

package screens

import (
	"strings"
    "fmt"
    "os"
    "path/filepath"
    "nextui-themes/internal/themes"
	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// GlobalOptionsMenuScreen displays the global options menu
func GlobalOptionsMenuScreen() (string, int) {
	// Menu items
	menu := []string{
		"Wallpapers",
		"Icon Packs",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Global Options")
}

// HandleGlobalOptionsMenu processes the user's selection from the global options menu
func HandleGlobalOptionsMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleGlobalOptionsMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Wallpapers":
			logging.LogDebug("Selected Wallpapers")
			app.SetSelectedThemeType(app.GlobalTheme)
			return app.Screens.WallpaperSelection

		case "Icon Packs":
			logging.LogDebug("Selected Icon Packs")
			return app.Screens.IconSelection

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.GlobalOptionsMenu
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.CustomizationMenu
	}

	return app.Screens.GlobalOptionsMenu
}

// WallpaperSelectionScreen displays available wallpapers using minui-presenter for a gallery view
func WallpaperSelectionScreen() (string, int) {
	// Get current directory for theme paths
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return "", 1
	}

	// Scan global themes directory for wallpapers
	globalDir := filepath.Join(cwd, "Themes", "Global")
	themesList, err := themes.ListGlobalThemes(globalDir)
	if err != nil {
		logging.LogDebug("Error loading global themes: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error loading global themes: %s", err), "3")
		return "", 1
	}

	if len(themesList) == 0 {
		logging.LogDebug("No global themes found")
		ui.ShowMessage("No global themes found. Create one in Themes/Global/", "3")
		return "", 1
	}

	// For wallpapers, display an image gallery
	return displayGlobalBackgroundsGallery(globalDir, themesList)
}

// HandleWallpaperSelection processes the user's wallpaper selection
func HandleWallpaperSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("handleWallpaperSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a wallpaper - proceed to confirmation
		app.SetSelectedTheme(selection)
		return app.Screens.ConfirmScreen
	case 1, 2:
		// User pressed cancel or back
		if app.GetSelectedSystem() != "" {
			// If we're in a system-specific context, go back to system options
			return app.Screens.SystemOptionsMenu
		}
		// Otherwise, go back to global options
		return app.Screens.GlobalOptionsMenu
	}

	return app.Screens.WallpaperSelection
}