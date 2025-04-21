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
		"Deconstruct...", // Added back with ellipsis to indicate it performs an action
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Components")
}

// HandleComponentsMenu processes the component type selection
func HandleComponentsMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleComponentsMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// If selected "Deconstruct...", go directly to deconstruction screen
		if selection == "Deconstruct..." {
			logging.LogDebug("Selected Deconstruct...")
			return app.Screens.Deconstruction
		}

		// Otherwise, set the selected component type and go to options
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

	// Path to catalog.json
	catalogPath := filepath.Join(cwd, "Catalog", "catalog.json")

	// Check if catalog exists
	if _, err := os.Stat(catalogPath); os.IsNotExist(err) {
		logging.LogDebug("Catalog file not found. Ask user to sync first.")
		ui.ShowMessage(fmt.Sprintf("No %s catalog found. Please sync first.", componentType), "3")
		return "", 1
	}

	// Parse the catalog
	data, err := os.ReadFile(catalogPath)
	if err != nil {
		logging.LogDebug("Error reading catalog file: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	var catalog themes.CatalogData
	if err := json.Unmarshal(data, &catalog); err != nil {
		logging.LogDebug("Error parsing catalog JSON: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Map component type to catalog key
	componentTypeMap := map[string]string{
		"Wallpapers": "wallpapers",
		"Icons":      "icons",
		"Accents":    "accents",
		"LEDs":       "leds",
		"Fonts":      "fonts",
		"Overlays":   "overlays",
	}

	catalogType := componentTypeMap[componentType]
	if catalogType == "" {
		logging.LogDebug("Unknown component type: %s", componentType)
		ui.ShowMessage(fmt.Sprintf("Unknown component type: %s", componentType), "3")
		return "", 1
	}

	// Get components of the selected type
	components, exists := catalog.Components[catalogType]
	if !exists || len(components) == 0 {
		logging.LogDebug("No %s found in catalog", componentType)
		ui.ShowMessage(fmt.Sprintf("No %s found in catalog", componentType), "3")
		return "", 1
	}

	// Get preview images
	previewImages := make([]ui.GalleryItem, 0, len(components))
	for compName, compInfo := range components {
		// Get preview path - relative path in catalog needs to be converted to absolute
		previewPath := filepath.Join(cwd, compInfo.PreviewPath)

		// Skip LEDs which don't have preview images
		if componentType == "LEDs" && (previewPath == "" || !fileExists(previewPath)) {
			// Just use the component name as text
			previewImages = append(previewImages, ui.GalleryItem{
				Text: fmt.Sprintf("%s by %s", compName, compInfo.Author),
				BackgroundImage: "", // No background image
			})
			continue
		}

		// Create a GalleryItem for this component
		previewItem := ui.GalleryItem{
			Text:            fmt.Sprintf("%s by %s", compName, compInfo.Author),
			BackgroundImage: previewPath,
		}

		previewImages = append(previewImages, previewItem)
	}

	// Use DisplayImageGallery to display a gallery of preview images
	selection, exitCode := ui.DisplayImageGallery(previewImages, fmt.Sprintf("Browse %s", componentType))

	logging.LogDebug("Gallery selection: %s, exit code: %d", selection, exitCode)

	// Extract component name from selection (remove author info)
	if selection != "" {
		parts := strings.Split(selection, " by ")
		selection = parts[0]
	}

	return selection, exitCode
}

// HandleBrowseComponents processes the component selection
func HandleBrowseComponents(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleBrowseComponents called with selection: '%s', exitCode: %d", selection, exitCode)
	componentType := app.GetSelectedComponentType()

	switch exitCode {
	case 0:
		// User selected a component
		if selection != "" {
			// First, download the component package
			if err := themes.DownloadComponentPackage(componentType, selection); err != nil {
				logging.LogDebug("Error downloading component: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
				return app.Screens.ComponentOptions
			}

			// Import/apply the selected component
			componentPath := filepath.Join(app.GetWorkingDir(), "Components", componentType, selection)
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

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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

	// Path to catalog.json
	catalogPath := filepath.Join(cwd, "Catalog", "catalog.json")

	// Check if catalog exists
	if _, err := os.Stat(catalogPath); os.IsNotExist(err) {
		logging.LogDebug("Catalog file not found. Ask user to sync first.")
		ui.ShowMessage("No theme catalog found. Please sync themes first.", "3")
		return "", 1
	}

	// Parse the catalog
	data, err := os.ReadFile(catalogPath)
	if err != nil {
		logging.LogDebug("Error reading catalog file: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	var catalog themes.CatalogData
	if err := json.Unmarshal(data, &catalog); err != nil {
		logging.LogDebug("Error parsing catalog JSON: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Check if there are themes
	if len(catalog.Themes) == 0 {
		logging.LogDebug("No themes found in catalog")
		ui.ShowMessage("No themes found in catalog", "3")
		return "", 1
	}

	// Get preview images
	previewImages := make([]ui.GalleryItem, 0, len(catalog.Themes))
	for themeName, themeInfo := range catalog.Themes {
		// Get preview path - relative path in catalog needs to be converted to absolute
		previewPath := filepath.Join(cwd, themeInfo.PreviewPath)

		// Create a GalleryItem for this theme
		previewItem := ui.GalleryItem{
			Text:            fmt.Sprintf("%s by %s", themeName, themeInfo.Author),
			BackgroundImage: previewPath,
		}

		previewImages = append(previewImages, previewItem)
	}

	// Use DisplayImageGallery to display a gallery of preview images
	selection, exitCode := ui.DisplayImageGallery(previewImages, "Browse Themes")

	logging.LogDebug("Gallery selection: %s, exit code: %d", selection, exitCode)

	// Extract theme name from selection (remove author info)
	if selection != "" {
		parts := strings.Split(selection, " by ")
		selection = parts[0]
	}

	return selection, exitCode
}

// HandleBrowseThemes processes the theme selection
func HandleBrowseThemes(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleBrowseThemes called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		if selection != "" {
			// First, download the theme package
			if err := themes.DownloadThemePackage(selection); err != nil {
				logging.LogDebug("Error downloading theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
				return app.Screens.MainMenu
			}

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