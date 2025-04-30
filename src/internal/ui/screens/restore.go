// internal/ui/screens/restore.go
package screens

import (
	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// RestoreThemeGalleryScreen displays the theme backup gallery screen
func RestoreThemeGalleryScreen() (string, int) {
	logging.LogDebug("Showing theme backup gallery screen")
	return themes.ShowThemeBackupGallery()
}

// HandleRestoreThemeGallery processes theme backup gallery selection
func HandleRestoreThemeGallery(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRestoreThemeGallery called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection != "" {
		// User selected a backup
		app.SetSelectedItem(selection)
		return app.Screens.RestoreThemeConfirm
	}

	// User cancelled or error
	return app.Screens.RestoreMenu
}

// RestoreThemeConfirmScreen displays the theme restore confirmation screen
func RestoreThemeConfirmScreen() (string, int) {
	logging.LogDebug("Showing theme restore confirmation screen")

	selectedBackup := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Restore from theme backup '" + selectedBackup + "'?")
}

// HandleRestoreThemeConfirm processes theme restore confirmation
func HandleRestoreThemeConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRestoreThemeConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed restore
		return app.Screens.RestoreThemeApplying
	}

	// User cancelled
	return app.Screens.RestoreThemeGallery
}

// RestoreThemeApplyingScreen handles the theme restore process
func RestoreThemeApplyingScreen() app.Screen {
	logging.LogDebug("Processing theme restore")

	selectedBackup := app.GetSelectedItem()

	// Show restoring message and perform operation
	err := ui.ShowMessageWithOperation(
		"Restoring from theme backup...",
		func() error {
			return themes.RevertThemeFromBackup(selectedBackup)
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error restoring from backup: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Theme restored successfully!", "2")
	}

	return app.Screens.SettingsMenu
}

// RestoreOverlayGalleryScreen displays the overlay backup gallery screen
func RestoreOverlayGalleryScreen() (string, int) {
	logging.LogDebug("Showing overlay backup gallery screen")
	return themes.ShowOverlayBackupGallery()
}

// HandleRestoreOverlayGallery processes overlay backup gallery selection
func HandleRestoreOverlayGallery(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRestoreOverlayGallery called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection != "" {
		// User selected a backup
		app.SetSelectedItem(selection)
		return app.Screens.RestoreOverlayConfirm
	}

	// User cancelled or error
	return app.Screens.RestoreMenu
}

// RestoreOverlayConfirmScreen displays the overlay restore confirmation screen
func RestoreOverlayConfirmScreen() (string, int) {
	logging.LogDebug("Showing overlay restore confirmation screen")

	selectedBackup := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Restore from overlay backup '" + selectedBackup + "'?")
}

// HandleRestoreOverlayConfirm processes overlay restore confirmation
func HandleRestoreOverlayConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRestoreOverlayConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed restore
		return app.Screens.RestoreOverlayApplying
	}

	// User cancelled
	return app.Screens.RestoreOverlayGallery
}

// RestoreOverlayApplyingScreen handles the overlay restore process
func RestoreOverlayApplyingScreen() app.Screen {
	logging.LogDebug("Processing overlay restore")

	selectedBackup := app.GetSelectedItem()

	// Show restoring message and perform operation
	err := ui.ShowMessageWithOperation(
		"Restoring from overlay backup...",
		func() error {
			return themes.RevertOverlayFromBackup(selectedBackup)
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error restoring from backup: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Overlays restored successfully!", "2")
	}

	return app.Screens.SettingsMenu
}