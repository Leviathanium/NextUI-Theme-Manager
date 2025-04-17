// src/internal/ui/screens/component_screens.go
// Implements UI screens for component management

package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/themes"
	"nextui-themes/internal/ui"
)

// ComponentsMenuScreen displays the menu of component types
func ComponentsMenuScreen() (string, int) {
	menu := []string{
		"Wallpapers",
		"Icons",
		"Accents",
		"Overlays",
		"LEDs",
		"Fonts",
		// "Deconstruction" removed from this menu
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Components")
}

// HandleComponentsMenu processes the component type selection
func HandleComponentsMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentsMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Set the selected component type
		app.SetSelectedComponentType(selection)
		return app.Screens.ComponentOptions

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.ComponentsMenu
}

// ComponentOptionsScreen displays options for the selected component type
func ComponentOptionsScreen() (string, int) {
	componentType := app.GetSelectedComponentType()

	menu := []string{
		"Browse",
		"Download",
		"Export",
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", componentType)
}

// HandleComponentOptions processes the component option selection
func HandleComponentOptions(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentOptions called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Set the selected option
		app.SetSelectedComponentOption(selection)

		// Process based on selected option
		switch selection {
		case "Browse":
			return app.Screens.BrowseComponents
		case "Download":
			return app.Screens.DownloadComponents
		case "Export":
			return app.Screens.ExportComponent
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentsMenu
	}

	return app.Screens.ComponentOptions
}

// BrowseComponentsScreen displays a browseable list of available components of the selected type
func BrowseComponentsScreen() (string, int) {
	componentType := app.GetSelectedComponentType()

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Map component type name to type constant and extension
	componentTypeMap := map[string]string{
		"Wallpapers": themes.ComponentWallpaper,
		"Icons":      themes.ComponentIcon,
		"Accents":    themes.ComponentAccent,
		"Overlays":   themes.ComponentOverlay,
		"LEDs":       themes.ComponentLED,
		"Fonts":      themes.ComponentFont,
	}

	// Map component type to subdirectory name
	componentDirMap := map[string]string{
		"Wallpapers": "Wallpapers",
		"Icons":      "Icons",
		"Accents":    "Accents",
		"Overlays":   "Overlays",
		"LEDs":       "LEDs",
		"Fonts":      "Fonts",
	}

	typeConstant := componentTypeMap[componentType]
	extension := themes.ComponentExtension[typeConstant]
	componentSubDir := componentDirMap[componentType]

	// Path to component directory where components are stored (inside Components directory)
	componentsDir := filepath.Join(cwd, "Components", componentSubDir)

	// Ensure directory exists
	if err := os.MkdirAll(componentsDir, 0755); err != nil {
		logging.LogDebug("Error creating components directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// List available components of the selected type
	entries, err := os.ReadDir(componentsDir)
	if err != nil {
		logging.LogDebug("Error reading components directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Filter for components of the selected type
	var components []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), extension) {
			components = append(components, entry.Name())
		}
	}

	if len(components) == 0 {
		logging.LogDebug("No %s found", componentType)
		ui.ShowMessage(fmt.Sprintf("No %s found in Components/%s directory", componentType, componentSubDir), "3")
		return "", 1
	}

	// Get preview images
	previewImages := make([]ui.GalleryItem, 0, len(components))
	for _, component := range components {
		previewPath := filepath.Join(componentsDir, component, "preview.png")

		// Skip LEDs which don't have preview images
		if typeConstant == themes.ComponentLED {
			// Just use the component name as text
			previewImages = append(previewImages, ui.GalleryItem{
				Text: component,
				BackgroundImage: "", // No background image
			})
			continue
		}

		// Check if preview exists
		if _, err := os.Stat(previewPath); err == nil {
			// Use the preview image
			previewImages = append(previewImages, ui.GalleryItem{
				Text: component,
				BackgroundImage: previewPath,
			})
		} else {
			// No preview image, just use the component name
			previewImages = append(previewImages, ui.GalleryItem{
				Text: component,
				BackgroundImage: "", // No background image
			})
		}
	}

	// Use DisplayImageGallery from presenter.go to display a gallery of preview images
	selection, exitCode := ui.DisplayImageGallery(previewImages, fmt.Sprintf("Browse %s", componentType))

	logging.LogDebug("Gallery selection: %s, exit code: %d", selection, exitCode)
	return selection, exitCode
}

// HandleBrowseComponents processes the component selection
func HandleBrowseComponents(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleBrowseComponents called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a component
		if selection != "" {
			// Get the component type and map to directory
			componentType := app.GetSelectedComponentType()
			componentDirMap := map[string]string{
				"Wallpapers": "Wallpapers",
				"Icons":      "Icons",
				"Accents":    "Accents",
				"Overlays":   "Overlays",
				"LEDs":       "LEDs",
				"Fonts":      "Fonts",
			}
			componentSubDir := componentDirMap[componentType]

			// Get the component path
			cwd, err := os.Getwd()
			if err != nil {
				logging.LogDebug("Error getting current directory: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
				return app.Screens.ComponentOptions
			}

			componentPath := filepath.Join(cwd, "Components", componentSubDir, selection)

			// Import/apply the selected component
			if err := themes.ImportComponent(componentPath); err != nil {
				logging.LogDebug("Error importing component: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			}
		}
		return app.Screens.ComponentOptions

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentOptions
	}

	return app.Screens.BrowseComponents
}

// ExportComponentScreen prompts for component name and exports the selected component type
func ExportComponentScreen() (string, int) {
	componentType := app.GetSelectedComponentType()

	// Generate export name with timestamp to ensure uniqueness
	timestamp := time.Now().Format("20060102_150405")
	exportName := fmt.Sprintf("%s_%s", strings.ToLower(componentType), timestamp)

	// Map component type to export function
	exportFunctions := map[string]func(string) error{
		"Wallpapers": themes.ExportWallpapers,
		"Icons":      themes.ExportIcons,
		"Accents":    themes.ExportAccents,
		"Overlays":   themes.ExportOverlays,
		"LEDs":       themes.ExportLEDs,
		"Fonts":      themes.ExportFonts,
	}

	// Get the export function
	exportFunc, ok := exportFunctions[componentType]
	if !ok {
		logging.LogDebug("Unknown component type: %s", componentType)
		ui.ShowMessage(fmt.Sprintf("Unknown component type: %s", componentType), "3")
		return "", 1
	}

	// Export the component
	if err := exportFunc(exportName); err != nil {
		logging.LogDebug("Error exporting component: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Return to component options screen
	return "", 0
}

// HandleExportComponent processes the export component result
func HandleExportComponent(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleExportComponent called with exitCode: %d", exitCode)

	// Always return to component options screen
	return app.Screens.ComponentOptions
}

// BrowseThemesScreen displays a browseable list of available themes
func BrowseThemesScreen() (string, int) {
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
	selection, exitCode := ui.DisplayImageGallery(previewImages, "Browse Themes")

	logging.LogDebug("Gallery selection: %s, exit code: %d", selection, exitCode)
	return selection, exitCode
}

// HandleBrowseThemes processes the theme selection
func HandleBrowseThemes(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleBrowseThemes called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		if selection != "" {
			// Set the selected theme
			app.SetSelectedTheme(selection)
			return app.Screens.ThemeImportConfirm
		}
		return app.Screens.MainMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.BrowseThemes
}

// DownloadThemesScreen and DownloadComponentsScreen are placeholders
// for future implementation of theme/component downloading

// DownloadThemesScreen is a placeholder for theme downloading
func DownloadThemesScreen() (string, int) {
	ui.ShowMessage("Theme downloading not implemented yet", "3")
	return "", 0
}

// HandleDownloadThemes processes the theme download result
func HandleDownloadThemes(selection string, exitCode int) app.Screen {
	return app.Screens.MainMenu
}

// DownloadComponentsScreen is a placeholder for component downloading
func DownloadComponentsScreen() (string, int) {
	ui.ShowMessage("Component downloading not implemented yet", "3")
	return "", 0
}

// HandleDownloadComponents processes the component download result
func HandleDownloadComponents(selection string, exitCode int) app.Screen {
	return app.Screens.ComponentOptions
}