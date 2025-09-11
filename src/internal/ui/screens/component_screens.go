// src/internal/ui/screens/component_screens.go
// Implements UI screens for component management - Updated with Installed/Download distinction

package screens

import (
	"encoding/json"
	"fmt"
	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"nextui-themes/internal/themes"
	"nextui-themes/internal/ui"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func ComponentsMenuScreen() (string, int) {
	// Updated menu without "Deconstruct..."
	menu := []string{
		"Wallpapers",
		"Icons",
		"Accents",
		"Overlays",
		"LEDs",
		"Fonts",
		// "Deconstruct..." option has been removed
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

	// Updated menu options - removed redundant "Sync Catalog" option
	menu := []string{
		"Installed", // Browse locally installed components
		"Download",  // Browse and download components from catalog
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

		// Process based on selected option and component type
		componentType := app.GetSelectedComponentType()

		// If this is overlays, go to system selection first
		if componentType == "Overlays" {
			// Clear any previously selected system tag
			app.SetSelectedSystemTag("")
			return app.Screens.OverlaySystemSelection // New screen for system selection
		} else {
			// For other component types, use existing flow
			switch selection {
			case "Installed":
				return app.Screens.InstalledComponents
			case "Download":
				return app.Screens.DownloadComponents
			case "Export":
				return app.Screens.ExportComponent
			}
		}

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentsMenu
	}

	return app.Screens.ComponentOptions
}

// Modified OverlaySystemSelectionScreen function to fix duplicated system tags
func OverlaySystemSelectionScreen() (string, int) {
	logging.LogDebug("Showing overlay system selection screen")

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Extract system tags and names into a list without duplicating tags
	var systemList []string
	for _, system := range systemPaths.Systems {
		if system.Tag != "" {
			// The system.Name already includes the tag format "Name (TAG)"
			// So we should just use that directly instead of formatting it again
			systemList = append(systemList, system.Name)
		}
	}

	if len(systemList) == 0 {
		logging.LogDebug("No systems with tags found")
		ui.ShowMessage("No systems with tags found", "3")
		return "", 1
	}

	// Sort the list alphabetically
	sort.Strings(systemList)

	componentType := app.GetSelectedComponentType()
	return ui.DisplayMinUiList(strings.Join(systemList, "\n"), "text",
		fmt.Sprintf("Select System for %s", componentType))
}

// HandleOverlaySystemSelection processes the selection from overlay system screen
func HandleOverlaySystemSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleOverlaySystemSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection != "" {
			// Extract system tag from selection "System Name (TAG)"
			re := regexp.MustCompile(`\((.*?)\)`)
			matches := re.FindStringSubmatch(selection)
			if len(matches) >= 2 {
				systemTag := matches[1]
				// Store the selected system tag
				app.SetSelectedSystemTag(systemTag)

				// Now go to the appropriate screen based on the previously selected option
				switch app.GetSelectedComponentOption() {
				case "Installed":
					return app.Screens.InstalledComponents
				case "Download":
					return app.Screens.DownloadComponents
				case "Export":
					return app.Screens.ExportComponent
				}
			}
		}
		// If no valid selection, go back to component options
		return app.Screens.ComponentOptions

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentOptions
	}

	return app.Screens.OverlaySystemSelection
}

// Complete InstalledComponentsScreen function with system tag filtering
func InstalledComponentsScreen() (string, int) {
	componentType := app.GetSelectedComponentType()
	systemTag := app.GetSelectedSystemTag()

	logging.LogDebug("Showing installed %s components screen", componentType)
	if systemTag != "" {
		logging.LogDebug("Filtering by system tag: %s", systemTag)
	}

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Path to Components directory for the selected type
	componentsDir := filepath.Join(cwd, "Components", componentType)

	// Check if directory exists
	if _, err := os.Stat(componentsDir); os.IsNotExist(err) {
		logging.LogDebug("Components directory not found: %s", componentsDir)
		ui.ShowMessage(fmt.Sprintf("No %s components found.", componentType), "3")
		return "", 1
	}

	// List available components
	entries, err := os.ReadDir(componentsDir)
	if err != nil {
		logging.LogDebug("Error reading components directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Filter for component directories with appropriate extension
	componentExt := ""
	switch componentType {
	case "Wallpapers":
		componentExt = ".bg"
	case "Icons":
		componentExt = ".icon"
	case "Accents":
		componentExt = ".acc"
	case "LEDs":
		componentExt = ".led"
	case "Fonts":
		componentExt = ".font"
	case "Overlays":
		componentExt = ".over"
	}

	var componentList []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), componentExt) {
			componentList = append(componentList, entry.Name())
		}
	}

	// For Overlays, filter by system tag if one is selected
	if componentType == "Overlays" && systemTag != "" {
		var filteredComponentList []string
		for _, compName := range componentList {
			compPath := filepath.Join(componentsDir, compName)

			// Check if this overlay supports the selected system
			if supportsSystem(compPath, systemTag) {
				filteredComponentList = append(filteredComponentList, compName)
			}
		}

		componentList = filteredComponentList
	}

	if len(componentList) == 0 {
		logging.LogDebug("No %s components found", componentType)
		if systemTag != "" {
			ui.ShowMessage(fmt.Sprintf("No installed %s components found for system %s.", componentType, systemTag), "3")
		} else {
			ui.ShowMessage(fmt.Sprintf("No installed %s components found.", componentType), "3")
		}
		return "", 1
	}

	// Get preview images for gallery display
	previewImages := make([]ui.GalleryItem, 0, len(componentList))
	for _, compName := range componentList {
		compPath := filepath.Join(componentsDir, compName)
		previewPath := filepath.Join(compPath, "preview.png")
		manifestPath := filepath.Join(compPath, "manifest.json")

		// Default text in case manifest can't be read
		text := compName

		// Try to read manifest for author info
		if fileExists(manifestPath) {
			if data, err := os.ReadFile(manifestPath); err == nil {
				var manifest map[string]interface{}
				if err := json.Unmarshal(data, &manifest); err == nil {
					if compInfo, ok := manifest["component_info"].(map[string]interface{}); ok {
						if author, ok := compInfo["author"].(string); ok {
							text = fmt.Sprintf("%s by %s", compName, author)
						}
					}
				}
			}
		}

		// Create gallery item with or without preview image
		if fileExists(previewPath) {
			previewImages = append(previewImages, ui.GalleryItem{
				Text:            text,
				BackgroundImage: previewPath,
			})
		} else {
			previewImages = append(previewImages, ui.GalleryItem{
				Text:            text,
				BackgroundImage: "", // No background image
			})
		}
	}

	// Use DisplayImageGallery to display a gallery of preview images
	title := fmt.Sprintf("Installed %s", componentType)
	if systemTag != "" {
		title = fmt.Sprintf("Installed %s for %s", componentType, systemTag)
	}
	selection, exitCode := ui.DisplayImageGallery(previewImages, title)

	// Extract component name from selection (remove author info)
	if selection != "" {
		parts := strings.Split(selection, " by ")
		selection = parts[0]
	}

	logging.LogDebug("Gallery selection: %s, exit code: %d", selection, exitCode)
	return selection, exitCode
}

// HandleInstalledComponents processes the selection of an installed component
func HandleInstalledComponents(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleInstalledComponents called with selection: '%s', exitCode: %d", selection, exitCode)
	componentType := app.GetSelectedComponentType()

	switch exitCode {
	case 0:
		// User selected a component to apply
		if selection != "" {
			// Import/apply the selected component
			componentPath := filepath.Join(app.GetWorkingDir(), "Components", componentType, selection)

			importErr := ui.ShowMessageWithOperation(
				fmt.Sprintf("Applying %s component '%s'...", componentType, selection),
				func() error {
					return themes.ImportComponent(componentPath)
				},
			)

			if importErr != nil {
				logging.LogDebug("Error importing component: %v", importErr)
				ui.ShowMessage(fmt.Sprintf("Error: %s", importErr), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("%s component applied successfully!", componentType), "2")
			}
		}
		return app.Screens.ComponentOptions

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentOptions
	}

	return app.Screens.InstalledComponents
}

// Complete DownloadComponentsScreen function with system tag filtering
func DownloadComponentsScreen() (string, int) {
	componentType := app.GetSelectedComponentType()
	systemTag := app.GetSelectedSystemTag()

	logging.LogDebug("Showing download %s components screen", componentType)
	if systemTag != "" {
		logging.LogDebug("Filtering by system tag: %s", systemTag)
	}

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
		ui.ShowMessage(fmt.Sprintf("No %s catalog found. Please sync catalog first.", componentType), "3")
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

	// For Overlays, filter by system tag if one is selected
	if componentType == "Overlays" && systemTag != "" {
		// Filter components by system tag
		filteredComponents := make(map[string]themes.CatalogItemInfo)

		for compName, compInfo := range components {
			// Try to load manifest from catalog path
			manifestPath := filepath.Join(cwd, compInfo.ManifestPath)
			if fileExists(manifestPath) {
				data, err := os.ReadFile(manifestPath)
				if err == nil {
					var manifest map[string]interface{}
					if err := json.Unmarshal(data, &manifest); err == nil {
						// Check if manifest has content.systems that includes our tag
						if content, hasContent := manifest["content"].(map[string]interface{}); hasContent {
							if systems, hasSystems := content["systems"].([]interface{}); hasSystems {
								for _, sys := range systems {
									if sysTag, ok := sys.(string); ok && sysTag == systemTag {
										filteredComponents[compName] = compInfo
										break
									}
								}
							}
						}
					}
				}
			}
		}

		components = filteredComponents
	}

	if len(components) == 0 {
		if systemTag != "" {
			logging.LogDebug("No %s components found for system tag %s", componentType, systemTag)
			ui.ShowMessage(fmt.Sprintf("No %s found in catalog for system %s", componentType, systemTag), "3")
		} else {
			logging.LogDebug("No %s components found in catalog", componentType)
			ui.ShowMessage(fmt.Sprintf("No %s found in catalog", componentType), "3")
		}
		return "", 1
	}

	// Get preview images
	previewImages := make([]ui.GalleryItem, 0, len(components))
	for compName, compInfo := range components {
		// Check if component already exists locally
		localComponentPath := filepath.Join(cwd, "Components", componentType, compName)
		alreadyInstalled := fileExists(localComponentPath)

		// Get preview path - relative path in catalog needs to be converted to absolute
		previewPath := filepath.Join(cwd, compInfo.PreviewPath)

		// Skip LEDs which don't have preview images
		if componentType == "LEDs" && (previewPath == "" || !fileExists(previewPath)) {
			// Just use the component name as text with installed indicator
			text := fmt.Sprintf("%s by %s", compName, compInfo.Author)
			if alreadyInstalled {
				text = "[Installed] " + text
			}

			previewImages = append(previewImages, ui.GalleryItem{
				Text:            text,
				BackgroundImage: "", // No background image
			})
			continue
		}

		// Create a GalleryItem for this component with installed indicator
		text := fmt.Sprintf("%s by %s", compName, compInfo.Author)
		if alreadyInstalled {
			text = "[Installed] " + text
		}

		previewItem := ui.GalleryItem{
			Text:            text,
			BackgroundImage: previewPath,
		}

		previewImages = append(previewImages, previewItem)
	}

	// Use DisplayImageGallery to display a gallery of preview images
	title := fmt.Sprintf("Download %s", componentType)
	if systemTag != "" {
		title = fmt.Sprintf("Download %s for %s", componentType, systemTag)
	}
	selection, exitCode := ui.DisplayImageGallery(previewImages, title)

	logging.LogDebug("Gallery selection: %s, exit code: %d", selection, exitCode)

	// Extract component name from selection (remove author info and installed indicator)
	if selection != "" {
		// Remove "[Installed] " prefix if present
		selection = strings.TrimPrefix(selection, "[Installed] ")

		// Split at " by " and take the first part
		parts := strings.Split(selection, " by ")
		selection = parts[0]
	}

	return selection, exitCode
}

// HandleDownloadComponents processes the component download selection
func HandleDownloadComponents(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleDownloadComponents called with selection: '%s', exitCode: %d", selection, exitCode)
	componentType := app.GetSelectedComponentType()

	switch exitCode {
	case 0:
		// User selected a component
		if selection != "" {
			// Check if component already exists locally
			cwd := app.GetWorkingDir()
			localComponentPath := filepath.Join(cwd, "Components", componentType, selection)

			if !fileExists(localComponentPath) {
				// Download the component package if not already installed
				if err := themes.DownloadComponentPackage(componentType, selection); err != nil {
					logging.LogDebug("Error downloading component: %v", err)
					ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
					return app.Screens.ComponentOptions
				}
			} else {
				logging.LogDebug("Component '%s' already installed, skipping download", selection)
			}

			// Prompt user if they want to apply this component now
			message := fmt.Sprintf("Apply %s '%s' now?", componentType, selection)
			options := []string{
				"Yes",
				"No",
			}
			result, promptCode := ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)

			// Inside HandleDownloadComponents where component is applied:
			if promptCode == 0 && result == "Yes" {
				// Import/apply the selected component with operation message
				componentPath := filepath.Join(app.GetWorkingDir(), "Components", componentType, selection)

				importErr := ui.ShowMessageWithOperation(
					fmt.Sprintf("Applying %s component '%s'...", componentType, selection),
					func() error {
						return themes.ImportComponent(componentPath)
					},
				)

				if importErr != nil {
					logging.LogDebug("Error importing component: %v", importErr)
					ui.ShowMessage(fmt.Sprintf("Error: %s", importErr), "3")
				} else {
					ui.ShowMessage(fmt.Sprintf("%s component applied successfully!", componentType), "2")
				}
			}
		}
		return app.Screens.ComponentOptions

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentOptions
	}

	return app.Screens.DownloadComponents
}

// Modified ExportComponentScreen function to properly display success messages
func ExportComponentScreen() (string, int) {
	componentType := app.GetSelectedComponentType()
	systemTag := app.GetSelectedSystemTag()

	// Generate export name with timestamp to ensure uniqueness
	timestamp := time.Now().Format("20060102_150405")
	var exportName string

	if componentType == "Overlays" && systemTag != "" {
		// Include system tag in export name for system-specific overlay exports
		exportName = fmt.Sprintf("%s_%s_%s", strings.ToLower(componentType), systemTag, timestamp)
	} else {
		exportName = fmt.Sprintf("%s_%s", strings.ToLower(componentType), timestamp)
	}

	// Map component type to export function
	exportFunctions := map[string]func(string) error{
		"Wallpapers": themes.ExportWallpapers,
		"Icons":      themes.ExportIcons,
		"Accents":    themes.ExportAccents,
		"Fonts":      themes.ExportFonts,
		"LEDs":       themes.ExportLEDs,
	}

	// For overlays with a system tag, use the new function
	var exportFunc func(string) error
	var exportErr error

	if componentType == "Overlays" {
		if systemTag != "" {
			// Use system-specific overlay export function
			exportErr = ui.ShowMessageWithOperation(
				fmt.Sprintf("Exporting %s component for system %s...", componentType, systemTag),
				func() error {
					return themes.ExportOverlaysForSystem(exportName, systemTag)
				},
			)
		} else {
			// Use general overlay export function
			exportErr = ui.ShowMessageWithOperation(
				fmt.Sprintf("Exporting %s component...", componentType),
				func() error {
					return themes.ExportOverlays(exportName)
				},
			)
		}
	} else {
		// Get the export function for other component types
		exportFunc, _ = exportFunctions[componentType]
		if exportFunc == nil {
			logging.LogDebug("Unknown component type: %s", componentType)
			ui.ShowMessage(fmt.Sprintf("Unknown component type: %s", componentType), "3")
			return "", 1
		}

		// Export the component with operation message
		exportErr = ui.ShowMessageWithOperation(
			fmt.Sprintf("Exporting %s component...", componentType),
			func() error {
				return exportFunc(exportName)
			},
		)
	}

	if exportErr != nil {
		logging.LogDebug("Error exporting component: %v", exportErr)
		ui.ShowMessage(fmt.Sprintf("Error: %s", exportErr), "3")
		return "", 1
	}

	// Show success message
	if componentType == "Overlays" && systemTag != "" {
		ui.ShowMessage(fmt.Sprintf("%s component for system %s exported successfully!", componentType, systemTag), "3")
	} else {
		ui.ShowMessage(fmt.Sprintf("%s component exported successfully!", componentType), "3")
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

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Helper function to check if an overlay supports a specific system tag
func supportsSystem(overlayPath string, systemTag string) bool {
	// First method: Check if the Systems directory contains the system tag
	systemDir := filepath.Join(overlayPath, "Systems", systemTag)
	if _, err := os.Stat(systemDir); err == nil {
		return true
	}

	// Second method: Check the manifest.json
	manifestObj, err := themes.LoadComponentManifest(overlayPath)
	if err != nil {
		return false
	}

	if manifest, ok := manifestObj.(*themes.OverlayManifest); ok {
		for _, tag := range manifest.Content.Systems {
			if tag == systemTag {
				return true
			}
		}
	}

	return false
}
