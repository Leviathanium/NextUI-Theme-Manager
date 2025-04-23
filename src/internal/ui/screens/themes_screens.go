// src/internal/ui/screens/theme_screens.go
// Implements UI screens for theme browsing and management - Updated with Installed/Download distinction

package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"encoding/json"
	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/themes"
	"nextui-themes/internal/ui"
)

// InstalledThemesScreen displays a browseable list of locally installed themes
func InstalledThemesScreen() (string, int) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Path to Themes directory
	themesDir := filepath.Join(cwd, "Themes")

	// Check if directory exists
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		logging.LogDebug("Themes directory not found: %s", themesDir)
		ui.ShowMessage("No installed themes found.", "3")
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
	var themeList []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".theme") {
			themeList = append(themeList, entry.Name())
		}
	}

	if len(themeList) == 0 {
		logging.LogDebug("No themes found")
		ui.ShowMessage("No installed themes found.", "3")
		return "", 1
	}

	// Get preview images for gallery display
	previewImages := make([]ui.GalleryItem, 0, len(themeList))
	for _, themeName := range themeList {
		themePath := filepath.Join(themesDir, themeName)
		previewPath := filepath.Join(themePath, "preview.png")
		manifestPath := filepath.Join(themePath, "manifest.json")

		// Default text in case manifest can't be read
		text := themeName

		// Try to read manifest for author info
		if fileExists(manifestPath) {
			if data, err := os.ReadFile(manifestPath); err == nil {
				var manifest map[string]interface{}
				if err := json.Unmarshal(data, &manifest); err == nil {
					if themeInfo, ok := manifest["theme_info"].(map[string]interface{}); ok {
						if author, ok := themeInfo["author"].(string); ok {
							text = fmt.Sprintf("%s by %s", themeName, author)
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
	selection, exitCode := ui.DisplayImageGallery(previewImages, "Installed Themes")

	// Extract theme name from selection (remove author info)
	if selection != "" {
		parts := strings.Split(selection, " by ")
		selection = parts[0]
	}

	logging.LogDebug("Gallery selection: %s, exit code: %d", selection, exitCode)
	return selection, exitCode
}

// HandleInstalledThemes processes the selection from installed themes
func HandleInstalledThemes(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleInstalledThemes called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		if selection != "" {
			// Set the selected theme for import/confirm
			app.SetSelectedTheme(selection)
			return app.Screens.ThemeImportConfirm
		}
		return app.Screens.MainMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.InstalledThemes
}

// DownloadThemesScreen displays a browseable list of themes from the catalog
func DownloadThemesScreen() (string, int) {
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
		ui.ShowMessage("No theme catalog found. Please sync catalog first.", "3")
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
		// Check if theme already exists locally
		localThemePath := filepath.Join(cwd, "Themes", themeName)
		alreadyInstalled := fileExists(localThemePath)

		// Get preview path - relative path in catalog needs to be converted to absolute
		previewPath := filepath.Join(cwd, themeInfo.PreviewPath)

		// Create text with installed indicator if needed
		text := fmt.Sprintf("%s by %s", themeName, themeInfo.Author)
		if alreadyInstalled {
			text = "[Installed] " + text
		}

		// Create a GalleryItem for this theme
		previewItem := ui.GalleryItem{
			Text:            text,
			BackgroundImage: previewPath,
		}

		previewImages = append(previewImages, previewItem)
	}

	// Use DisplayImageGallery to display a gallery of preview images
	selection, exitCode := ui.DisplayImageGallery(previewImages, "Download Themes")

	logging.LogDebug("Gallery selection: %s, exit code: %d", selection, exitCode)

	// Extract theme name from selection (remove author info and installed indicator)
	if selection != "" {
		// Remove "[Installed] " prefix if present
		selection = strings.TrimPrefix(selection, "[Installed] ")

		// Split at " by " and take the first part
		parts := strings.Split(selection, " by ")
		selection = parts[0]
	}

	return selection, exitCode
}

// HandleDownloadThemes processes the theme download selection
func HandleDownloadThemes(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleDownloadThemes called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		if selection != "" {
			// Check if theme already exists locally
			cwd := app.GetWorkingDir()
			localThemePath := filepath.Join(cwd, "Themes", selection)

			if !fileExists(localThemePath) {
				// Download the theme package if not already installed
				downloadErr := ui.ShowMessageWithOperation(
					fmt.Sprintf("Downloading theme '%s'...", selection),
					func() error {
						return themes.DownloadThemePackage(selection)
					},
				)

				if downloadErr != nil {
					logging.LogDebug("Error downloading theme: %v", downloadErr)
					ui.ShowMessage(fmt.Sprintf("Error: %s", downloadErr), "3")
					return app.Screens.MainMenu
				}

				// Show success message briefly
				ui.ShowMessage(fmt.Sprintf("Theme '%s' downloaded successfully!", selection), "2")
			} else {
				logging.LogDebug("Theme '%s' already installed, skipping download", selection)
			}

			// Prompt user if they want to apply this theme now
			message := fmt.Sprintf("Apply theme '%s' now?", selection)
			options := []string{
				"Yes",
				"No",
			}
			result, promptCode := ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)

			if promptCode == 0 && result == "Yes" {
				// Apply the theme using the new function
				importErr := ui.ShowMessageWithOperation(
					fmt.Sprintf("Applying theme '%s'...", selection),
					func() error {
						return themes.ImportTheme(selection)
					},
				)

				if importErr != nil {
					logging.LogDebug("Error importing theme: %v", importErr)
					ui.ShowMessage(fmt.Sprintf("Error: %s", importErr), "3")
				} else {
					ui.ShowMessage(fmt.Sprintf("Theme '%s' applied successfully!", selection), "2")
				}
			}
		}
		return app.Screens.MainMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.DownloadThemes
}

// SyncCatalogScreen displays the sync catalog screen
func SyncCatalogScreen() (string, int) {
	// Simple confirmation message
	message := fmt.Sprintf("Sync catalog from %s?\nThis will download the latest theme and component catalog.",
		themes.RepoConfig.URL)
	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleSyncCatalog processes the user's choice to sync the catalog
func HandleSyncCatalog(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleSyncCatalog called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Perform catalog sync
			logging.LogDebug("Starting catalog sync")

			// Get default sync options
			options := themes.GetDefaultSyncOptions()

			// Sync catalog with operation message
			syncErr := ui.ShowMessageWithOperation(
				"Syncing theme catalog...",
				func() error {
					return themes.SyncThemeCatalog(options)
				},
			)

			if syncErr != nil {
				logging.LogDebug("Error syncing catalog: %v", syncErr)
				ui.ShowMessage(fmt.Sprintf("Error: %s", syncErr), "3")
			} else {
				logging.LogDebug("Catalog sync completed successfully")
				ui.ShowMessage("Catalog synced successfully!", "2")
			}
		}
		// Return to main menu
		return app.Screens.MainMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.SyncCatalog
}

// ThemeImportScreen displays available themes from the Themes directory
func ThemeImportScreen() (string, int) {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Path to Themes directory
	importsDir := filepath.Join(cwd, "Themes")

	// Ensure directory exists
	if err := os.MkdirAll(importsDir, 0755); err != nil {
		logging.LogDebug("Error creating themes directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// List available themes
	entries, err := os.ReadDir(importsDir)
	if err != nil {
		logging.LogDebug("Error reading themes directory: %v", err)
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
		ui.ShowMessage("No themes found in Themes directory", "3")
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

			// Use ShowMessageWithOperation for better user feedback
			importErr := ui.ShowMessageWithOperation(
				fmt.Sprintf("Applying theme '%s'...", themeName),
				func() error {
					return themes.ImportTheme(themeName)
				},
			)

			if importErr != nil {
				logging.LogDebug("Error importing theme: %v", importErr)
				ui.ShowMessage(fmt.Sprintf("Error: %s", importErr), "3")
			} else {
				ui.ShowMessage(fmt.Sprintf("Theme '%s' applied successfully!", themeName), "3")
			}
		}
		// Return to main menu
		return app.Screens.MainMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.InstalledThemes
	}

	return app.Screens.ThemeImportConfirm
}

// ThemeExportScreen displays the theme export confirmation
func ThemeExportScreen() (string, int) {
	// Simple confirmation message
	message := "Export current theme settings?\nThis will create a theme package in the Exports directory."
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
			// Perform theme export with operation message
			exportErr := ui.ShowMessageWithOperation(
				"Exporting current theme...",
				func() error {
					return themes.ExportTheme()
				},
			)

			if exportErr != nil {
				logging.LogDebug("Error exporting theme: %v", exportErr)
				ui.ShowMessage(fmt.Sprintf("Error: %s", exportErr), "3")
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