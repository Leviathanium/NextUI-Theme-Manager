// internal/ui/screens/overlays.go
package screens

import (
	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// OverlayGalleryScreen displays the overlay gallery screen
func OverlayGalleryScreen() (string, int) {
	logging.LogDebug("Showing overlay gallery screen")
	return themes.ShowOverlayGallery()
}

// HandleOverlayGallery processes overlay gallery selection
func HandleOverlayGallery(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleOverlayGallery called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection != "" {
		// User selected an overlay
		app.SetSelectedItem(selection)

		// Check if overlay is already downloaded
		if themes.IsOverlayDownloaded(selection) {
			logging.LogDebug("Overlay already downloaded: %s", selection)
			return app.Screens.OverlayApplyConfirm
		}

		// Not downloaded, ask to download
		logging.LogDebug("Overlay not downloaded: %s", selection)
		return app.Screens.OverlayDownloadConfirm
	}

	// User cancelled or error
	return app.Screens.MainMenu
}

// DownloadOverlaysScreen displays the overlay catalog/download screen
func DownloadOverlaysScreen() (string, int) {
	logging.LogDebug("Showing download overlays screen")
	return themes.ShowOverlayGallery()
}

// HandleDownloadOverlays processes overlay catalog selection
func HandleDownloadOverlays(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleDownloadOverlays called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection != "" {
		// User selected an overlay
		app.SetSelectedItem(selection)

		// Check if overlay is already downloaded
		if themes.IsOverlayDownloaded(selection) {
			logging.LogDebug("Overlay already downloaded: %s", selection)
			return app.Screens.OverlayApplyConfirm
		}

		// Not downloaded, ask to download
		logging.LogDebug("Overlay not downloaded: %s", selection)
		return app.Screens.OverlayDownloadConfirm
	}

	// User cancelled or error
	return app.Screens.OverlaysMenu
}

// OverlayDownloadConfirmScreen displays the overlay download confirmation screen
func OverlayDownloadConfirmScreen() (string, int) {
	logging.LogDebug("Showing overlay download confirmation screen")

	selectedOverlay := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Download overlay '" + selectedOverlay + "'?")
}

// HandleOverlayDownloadConfirm processes overlay download confirmation
func HandleOverlayDownloadConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleOverlayDownloadConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed download
		return app.Screens.OverlayDownloading
	}

	// User cancelled
	return app.Screens.DownloadOverlays
}

// OverlayDownloadingScreen handles the overlay downloading process
func OverlayDownloadingScreen() app.Screen {
	logging.LogDebug("Processing overlay download")

	selectedOverlay := app.GetSelectedItem()

	// Show downloading message and perform operation
	err := ui.ShowMessageWithOperation(
		"Downloading overlay '" + selectedOverlay + "'...",
		func() error {
			return themes.DownloadOverlay(selectedOverlay)
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error downloading overlay: " + err.Error(), "3")
		return app.Screens.OverlayGallery
	}

	ui.ShowMessage("Overlay downloaded successfully!", "2")
	return app.Screens.OverlayApplyConfirm
}

// OverlayApplyConfirmScreen displays the overlay apply confirmation screen
func OverlayApplyConfirmScreen() (string, int) {
	logging.LogDebug("Showing overlay apply confirmation screen")

	selectedOverlay := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Apply overlay '" + selectedOverlay + "'?")
}

// HandleOverlayApplyConfirm processes overlay apply confirmation
func HandleOverlayApplyConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleOverlayApplyConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed apply
		return app.Screens.OverlayApplying
	}

	// User cancelled
	return app.Screens.OverlayGallery
}

// OverlayApplyingScreen handles the overlay applying process
func OverlayApplyingScreen() app.Screen {
	logging.LogDebug("Processing overlay application")

	selectedOverlay := app.GetSelectedItem()

	// Create backup if auto-backup is enabled
	if app.GetAutoBackup() {
		err := themes.CreateOverlayBackup("auto")
		if err != nil {
			logging.LogDebug("Error creating auto-backup: %v", err)
		}
	}

	// Show applying message and perform operation
	err := ui.ShowMessageWithOperation(
		"Applying overlay '" + selectedOverlay + "'...",
		func() error {
			return themes.ApplyOverlay(selectedOverlay)
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error applying overlay: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Overlay applied successfully!", "2")
	}

	return app.Screens.MainMenu
}