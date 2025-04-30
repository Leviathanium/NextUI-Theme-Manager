// internal/ui/screens/themes.go
package screens

import (
	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// ThemeGalleryScreen displays the theme gallery screen
func ThemeGalleryScreen() (string, int) {
	logging.LogDebug("Showing theme gallery screen")
	return themes.ShowThemeGallery()
}

// HandleThemeGallery processes theme gallery selection
func HandleThemeGallery(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeGallery called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection != "" {
		// User selected a theme
		app.SetSelectedItem(selection)

		// Check if theme is already downloaded
		if themes.IsThemeDownloaded(selection) {
			logging.LogDebug("Theme already downloaded: %s", selection)
			return app.Screens.ThemeApplyConfirm
		}

		// Not downloaded, ask to download
		logging.LogDebug("Theme not downloaded: %s", selection)
		return app.Screens.ThemeDownloadConfirm
	}

	// User cancelled or error
	return app.Screens.MainMenu
}

// ThemeDownloadConfirmScreen displays the theme download confirmation screen
func ThemeDownloadConfirmScreen() (string, int) {
	logging.LogDebug("Showing theme download confirmation screen")

	selectedTheme := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Download theme '" + selectedTheme + "'?")
}

// HandleThemeDownloadConfirm processes theme download confirmation
func HandleThemeDownloadConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeDownloadConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed download
		return app.Screens.ThemeDownloading
	}

	// User cancelled
	return app.Screens.ThemeGallery
}

// ThemeDownloadingScreen handles the theme downloading process
func ThemeDownloadingScreen() app.Screen {
	logging.LogDebug("Processing theme download")

	selectedTheme := app.GetSelectedItem()

	// Show downloading message and perform operation
	err := ui.ShowMessageWithOperation(
		"Downloading theme '" + selectedTheme + "'...",
		func() error {
			return themes.DownloadTheme(selectedTheme)
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error downloading theme: " + err.Error(), "3")
		return app.Screens.ThemeGallery
	}

	ui.ShowMessage("Theme downloaded successfully!", "2")
	return app.Screens.ThemeApplyConfirm
}

// ThemeApplyConfirmScreen displays the theme apply confirmation screen
func ThemeApplyConfirmScreen() (string, int) {
	logging.LogDebug("Showing theme apply confirmation screen")

	selectedTheme := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Apply theme '" + selectedTheme + "'?")
}

// HandleThemeApplyConfirm processes theme apply confirmation
func HandleThemeApplyConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleThemeApplyConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed apply
		return app.Screens.ThemeApplying
	}

	// User cancelled
	return app.Screens.ThemeGallery
}

// ThemeApplyingScreen handles the theme applying process
func ThemeApplyingScreen() app.Screen {
	logging.LogDebug("Processing theme application")

	selectedTheme := app.GetSelectedItem()

	// Create backup if auto-backup is enabled
	if app.GetAutoBackup() {
		err := themes.CreateThemeBackup("auto")
		if err != nil {
			logging.LogDebug("Error creating auto-backup: %v", err)
		}
	}

	// Show applying message and perform operation
	err := ui.ShowMessageWithOperation(
		"Applying theme '" + selectedTheme + "'...",
		func() error {
			return themes.ApplyTheme(selectedTheme)
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error applying theme: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Theme applied successfully!", "2")
	}

	return app.Screens.MainMenu
}