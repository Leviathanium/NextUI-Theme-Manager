// internal/ui/screens/revert.go
package screens

import (
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// RevertMenuScreen displays the revert menu screen
func RevertMenuScreen() (string, int) {
	logging.LogDebug("Showing revert menu screen")

	menuItems := []string{
		"Revert Theme",
		"Revert Overlays",
	}

	return ui.DisplayMinUiList(
		strings.Join(menuItems, "\n"),
		"text",
		"Revert Menu",
	)
}

// HandleRevertMenu processes revert menu selection
func HandleRevertMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRevertMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Revert Theme":
			return app.Screens.RevertThemeGallery

		case "Revert Overlays":
			return app.Screens.RevertOverlayGallery
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/back
		return app.Screens.MainMenu
	}

	return app.Screens.RevertMenu
}

// RevertThemeGalleryScreen displays the theme backup gallery screen
func RevertThemeGalleryScreen() (string, int) {
	logging.LogDebug("Showing theme backup gallery screen")
	return themes.ShowThemeBackupGallery()
}

// HandleRevertThemeGallery processes theme backup gallery selection
func HandleRevertThemeGallery(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRevertThemeGallery called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection != "" {
		// User selected a backup
		app.SetSelectedItem(selection)
		return app.Screens.RevertThemeConfirm
	}

	// User cancelled or error
	return app.Screens.RevertMenu
}

// RevertThemeConfirmScreen displays the theme revert confirmation screen
func RevertThemeConfirmScreen() (string, int) {
	logging.LogDebug("Showing theme revert confirmation screen")

	selectedBackup := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Revert to theme backup '" + selectedBackup + "'?")
}

// HandleRevertThemeConfirm processes theme revert confirmation
func HandleRevertThemeConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRevertThemeConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed revert
		return app.Screens.RevertThemeApplying
	}

	// User cancelled
	return app.Screens.RevertThemeGallery
}

// RevertThemeApplyingScreen handles the theme revert process
func RevertThemeApplyingScreen() app.Screen {
	logging.LogDebug("Processing theme revert")

	selectedBackup := app.GetSelectedItem()

	// Show reverting message and perform operation
	err := ui.ShowMessageWithOperation(
		"Reverting from theme backup...",
		func() error {
			return themes.RevertThemeFromBackup(selectedBackup)
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error reverting from backup: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Theme reverted successfully!", "2")
	}

	return app.Screens.MainMenu
}

// RevertOverlayGalleryScreen displays the overlay backup gallery screen
func RevertOverlayGalleryScreen() (string, int) {
	logging.LogDebug("Showing overlay backup gallery screen")
	return themes.ShowOverlayBackupGallery()
}

// HandleRevertOverlayGallery processes overlay backup gallery selection
func HandleRevertOverlayGallery(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRevertOverlayGallery called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection != "" {
		// User selected a backup
		app.SetSelectedItem(selection)
		return app.Screens.RevertOverlayConfirm
	}

	// User cancelled or error
	return app.Screens.RevertMenu
}

// RevertOverlayConfirmScreen displays the overlay revert confirmation screen
func RevertOverlayConfirmScreen() (string, int) {
	logging.LogDebug("Showing overlay revert confirmation screen")

	selectedBackup := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Revert to overlay backup '" + selectedBackup + "'?")
}

// HandleRevertOverlayConfirm processes overlay revert confirmation
func HandleRevertOverlayConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRevertOverlayConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed revert
		return app.Screens.RevertOverlayApplying
	}

	// User cancelled
	return app.Screens.RevertOverlayGallery
}

// RevertOverlayApplyingScreen handles the overlay revert process
func RevertOverlayApplyingScreen() app.Screen {
	logging.LogDebug("Processing overlay revert")

	selectedBackup := app.GetSelectedItem()

	// Show reverting message and perform operation
	err := ui.ShowMessageWithOperation(
		"Reverting from overlay backup...",
		func() error {
			return themes.RevertOverlayFromBackup(selectedBackup)
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error reverting from backup: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Overlays reverted successfully!", "2")
	}

	return app.Screens.MainMenu
}