// src/internal/ui/screens/deconstruction_screens.go
// Implementation of theme deconstruction screens

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

// DeconstructionScreen displays a browseable list of available themes to deconstruct
func DeconstructionScreen() (string, int) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Path to Themes directory
	themesDir := filepath.Join(cwd, "Themes")

	// Ensure directory exists
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		logging.LogDebug("Error creating themes directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// List available themes
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		logging.LogDebug("Error reading themes directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Filter for theme directories
	var themes []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".theme") {
			themes = append(themes, entry.Name())
		}
	}

	if len(themes) == 0 {
		logging.LogDebug("No themes found")
		ui.ShowMessage("No themes found in Themes directory", "3")
		return "", 1
	}

	// Get preview images
	previewImages := make([]ui.GalleryItem, 0, len(themes))
	for _, theme := range themes {
		previewPath := filepath.Join(themesDir, theme, "preview.png")

		// Check if preview exists
		if _, err := os.Stat(previewPath); err == nil {
			// Use the preview image
			previewImages = append(previewImages, ui.GalleryItem{
				Text: theme,
				BackgroundImage: previewPath,
			})
		} else {
			// No preview image, just use the theme name
			previewImages = append(previewImages, ui.GalleryItem{
				Text: theme,
				BackgroundImage: "", // No background image
			})
		}
	}

	// Use DisplayImageGallery from presenter.go to display a gallery of preview images
	selection, exitCode := ui.DisplayImageGallery(previewImages, "Select Theme to Deconstruct")

	logging.LogDebug("Gallery selection: %s, exit code: %d", selection, exitCode)
	return selection, exitCode
}

// HandleDeconstruction processes the theme selection for deconstruction
func HandleDeconstruction(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleDeconstruction called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		if selection != "" {
			// Set the selected theme
			app.SetSelectedTheme(selection)
			return app.Screens.DeconstructConfirm
		}
		return app.Screens.ComponentsMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentsMenu
	}

	return app.Screens.Deconstruction
}

// DeconstructConfirmScreen displays a confirmation dialog for theme deconstruction
func DeconstructConfirmScreen() (string, int) {
	themeName := app.GetSelectedTheme()
	message := fmt.Sprintf("Deconstruct theme '%s' into component packages?", themeName)

	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleDeconstructConfirm processes the confirmation for theme deconstruction
func HandleDeconstructConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleDeconstructConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Deconstruct the selected theme with operation message
			themeName := app.GetSelectedTheme()

			deconstructErr := ui.ShowMessageWithOperation(
				fmt.Sprintf("Deconstructing theme '%s'...", themeName),
				func() error {
					return themes.DeconstructTheme(themeName)
				},
			)

			if deconstructErr != nil {
				logging.LogDebug("Error deconstructing theme: %v", deconstructErr)
				ui.ShowMessage(fmt.Sprintf("Error: %s", deconstructErr), "3")
			}
		}
		// Return to components menu
		return app.Screens.ComponentsMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.Deconstruction
	}

	return app.Screens.DeconstructConfirm
}