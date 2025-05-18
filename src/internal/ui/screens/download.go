// internal/ui/screens/download.go
package screens

import (
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/ui"
)

// ShowDownloadThemeScreen displays the theme selection screen for downloading
func ShowDownloadThemeScreen() (string, int) {
	app.LogDebug("Showing download theme screen")

	// Check if catalog has been synced - this is a placeholder,
	// will be implemented with actual logic later
	catalogSynced := false

	if !catalogSynced {
		// Catalog not synced, skip to warning screen
		return "", 0
	}

	// This is a placeholder - later we'll implement a gallery view
	// showing all available themes in the catalog

	menuItems := []string{
		"Remote Theme 1",
		"Remote Theme 2",
		"Remote Theme 3",
	}

	return ui.ShowMenu(
		strings.Join(menuItems, "\n"),
		"Select Theme to Download",
		"--cancel-text", "BACK",
	)
}

// HandleDownloadThemeScreen processes the theme selection
func HandleDownloadThemeScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleDownloadThemeScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	// Check if catalog has been synced - this is a placeholder
	catalogSynced := false

	if !catalogSynced {
		// Catalog not synced, go to warning screen
		return app.ScreenCatalogNotSynced
	}

	if exitCode == 0 {
		// User selected a theme
		app.SetSelectedItem(selection)

		// Check if theme is already downloaded - placeholder
		themeDownloaded := false

		if themeDownloaded {
			return app.ScreenThemeDownloaded
		} else {
			return app.ScreenDownloadThemeConfirm
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User cancelled
		return app.ScreenMainMenu
	}

	return app.ScreenDownloadTheme
}

// ShowCatalogNotSyncedScreen displays warning when catalog is not synced
func ShowCatalogNotSyncedScreen() (string, int) {
	app.LogDebug("Showing catalog not synced screen")

	return ui.ShowMessage("Catalog not synced. Please sync the catalog first.", "3")
}

// HandleCatalogNotSyncedScreen processes the warning
func HandleCatalogNotSyncedScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleCatalogNotSyncedScreen called with exitCode: %d", exitCode)

	// Return to main menu after showing warning
	return app.ScreenMainMenu
}

// ShowDownloadThemeConfirmScreen displays the confirmation screen for downloading
func ShowDownloadThemeConfirmScreen() (string, int) {
	app.LogDebug("Showing download theme confirmation screen")

	selectedTheme := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Download theme '" + selectedTheme + "'?")
}

// HandleDownloadThemeConfirmScreen processes the confirmation result
func HandleDownloadThemeConfirmScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleDownloadThemeConfirmScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed - proceed to downloading
		return app.ScreenDownloadingTheme
	} else {
		// User cancelled
		return app.ScreenDownloadTheme
	}
}

// ShowDownloadingThemeScreen displays the theme downloading progress screen
func ShowDownloadingThemeScreen() (string, int) {
	app.LogDebug("Showing downloading theme screen")

	selectedTheme := app.GetSelectedItem()
	return ui.ShowMessage("Downloading theme '" + selectedTheme + "'...", "2")
}

// HandleDownloadingThemeScreen processes the downloading operation
func HandleDownloadingThemeScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleDownloadingThemeScreen called with exitCode: %d", exitCode)

	// Theme download would happen here in the actual implementation
	// For now, just show success message

	return app.ScreenThemeDownloaded
}

// ShowThemeDownloadedScreen displays the theme downloaded success/already exists screen
func ShowThemeDownloadedScreen() (string, int) {
	app.LogDebug("Showing theme downloaded screen")

	selectedTheme := app.GetSelectedItem()

	// Check if the theme was just downloaded or was already present
	justDownloaded := true

	if justDownloaded {
		return ui.ShowMessage("Theme '" + selectedTheme + "' downloaded successfully!", "2")
	} else {
		return ui.ShowMessage("Theme '" + selectedTheme + "' is already downloaded.", "2")
	}
}

// HandleThemeDownloadedScreen processes the success screen
func HandleThemeDownloadedScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleThemeDownloadedScreen called with exitCode: %d", exitCode)

	// Ask if user wants to apply the theme now
	// This would be added in the actual implementation

	// For now, return to download screen
	return app.ScreenDownloadTheme
}