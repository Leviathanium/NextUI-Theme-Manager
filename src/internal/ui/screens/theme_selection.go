// src/internal/ui/screens/theme_selection.go
// Implementation of the theme selection screen

package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"nextui-themes/internal/themes"
	"nextui-themes/internal/ui"
)

// ThemeSelectionScreen displays available themes based on the selected theme type
func ThemeSelectionScreen() (string, int) {
	var title string
	var themesList []string
	var err error

	// Get current directory for theme paths
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return "", 1
	}

	switch app.GetSelectedThemeType() {
	case app.GlobalTheme:
		title = "Select Global Background"

		// Scan global themes directory
		globalDir := filepath.Join(cwd, "Themes", "Global")
		themesList, err = themes.ListGlobalThemes(globalDir)
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

		// For Global Backgrounds, display an image gallery instead of a text list
		return displayGlobalBackgroundsGallery(globalDir, themesList)

	case app.DynamicTheme:
		title = "Select Dynamic Theme"
		// List actual dynamic themes
		themesList, err = themes.ListDynamicThemes()
		if err != nil {
			logging.LogDebug("Error loading dynamic themes: %v", err)
			ui.ShowMessage(fmt.Sprintf("Error loading dynamic themes: %s", err), "3")
			themesList = []string{"No themes found"}
		}

		// If no themes found, show a message
		if len(themesList) == 0 {
			logging.LogDebug("No dynamic themes found")
			ui.ShowMessage("No dynamic themes found. Create one in Themes/Dynamic/", "3")
			themesList = []string{"No themes found"}
		}

	case app.CustomTheme:
		title = "Select System"

		// Get system paths to find all installed systems
		systemPaths, err := system.GetSystemPaths()
		if err != nil {
			logging.LogDebug("Error getting system paths: %v", err)
			ui.ShowMessage(fmt.Sprintf("Error detecting systems: %s", err), "3")
			return "", 1
		}

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
	}

	// For non-global themes (Dynamic, Custom), use the standard list view
	logging.LogDebug("Displaying theme selection with %d options", len(themesList))
	return ui.DisplayMinUiList(strings.Join(themesList, "\n"), "text", title)
}

// displayGlobalBackgroundsGallery displays available global backgrounds as an image gallery
func displayGlobalBackgroundsGallery(globalDir string, themeNames []string) (string, int) {
	logging.LogDebug("Preparing global backgrounds gallery with %d themes", len(themeNames))

	// Create gallery items from theme list
	var galleryItems []ui.GalleryItem
	for _, themeName := range themeNames {
		// First check for new style PNG file directly in the directory
		bgPath := filepath.Join(globalDir, themeName + ".png")
		if _, err := os.Stat(bgPath); err != nil {
			// If not found, check old style directory with bg.png
			oldStylePath := filepath.Join(globalDir, themeName, "bg.png")
			if _, err := os.Stat(oldStylePath); err == nil {
				bgPath = oldStylePath
			} else {
				// Neither style found, skip this theme
				logging.LogDebug("No image found for theme %s", themeName)
				continue
			}
		}

		// Add to gallery items
		galleryItems = append(galleryItems, ui.GalleryItem{
			Text:            themeName,
			BackgroundImage: bgPath,
		})
	}

	// Display gallery
	return ui.DisplayImageGallery(galleryItems, "Select Global Background")
}

// HandleThemeSelection processes the user's theme selection
func HandleThemeSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("handleThemeSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme - proceed to confirmation
		app.SetSelectedTheme(selection)
		return app.Screens.ConfirmScreen
	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.ThemeSelection
}