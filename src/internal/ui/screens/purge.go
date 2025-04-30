// internal/ui/screens/purge.go
package screens

import (
	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// PurgeConfirmScreen displays the purge confirmation screen
func PurgeConfirmScreen() (string, int) {
	logging.LogDebug("Showing purge confirmation screen")
	return ui.ShowConfirmDialog("WARNING: Erase everything?")
}

// HandlePurgeConfirm processes purge confirmation
func HandlePurgeConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandlePurgeConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed purge
		return app.Screens.Purging
	}

	// User cancelled
	return app.Screens.SettingsMenu
}

// PurgingScreen handles the purge process
func PurgingScreen() app.Screen {
	logging.LogDebug("Processing purge")

	// Show purging message and perform operation
	err := ui.ShowMessageWithOperation(
		"Purging...",
		func() error {
			return themes.PurgeAll()
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error during purge: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Purge complete!", "2")
	}

	return app.Screens.SettingsMenu
}