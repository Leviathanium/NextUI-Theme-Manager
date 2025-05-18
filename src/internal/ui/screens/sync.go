// internal/ui/screens/sync.go
package screens

import (
	"thememanager/internal/app"
	"thememanager/internal/ui"
)

// ShowSyncCatalogScreen displays the sync catalog confirmation screen
func ShowSyncCatalogScreen() (string, int) {
	app.LogDebug("Showing sync catalog screen")

	return ui.ShowConfirmDialog("Sync theme catalog from repository?")
}

// HandleSyncCatalogScreen processes the sync catalog confirmation
func HandleSyncCatalogScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleSyncCatalogScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed - proceed to syncing
		return app.ScreenSyncingCatalog
	} else {
		// User cancelled
		return app.ScreenMainMenu
	}
}

// ShowSyncingCatalogScreen displays the sync progress screen
func ShowSyncingCatalogScreen() (string, int) {
	app.LogDebug("Showing syncing catalog screen")

	return ui.ShowMessage("Syncing theme catalog...", "2")
}

// HandleSyncingCatalogScreen processes the sync operation
func HandleSyncingCatalogScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleSyncingCatalogScreen called with exitCode: %d", exitCode)

	// Catalog sync would happen here in the actual implementation
	// For now, simulate success or failure

	// Simulated success for now
	syncSuccess := true

	if syncSuccess {
		return app.ScreenSyncComplete
	} else {
		return app.ScreenSyncFailed
	}
}

// ShowSyncCompleteScreen displays the sync complete success screen
func ShowSyncCompleteScreen() (string, int) {
	app.LogDebug("Showing sync complete screen")

	return ui.ShowMessage("Catalog synchronized successfully!", "2")
}

// HandleSyncCompleteScreen processes the sync complete screen
func HandleSyncCompleteScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleSyncCompleteScreen called with exitCode: %d", exitCode)

	// Return to main menu after showing success message
	return app.ScreenMainMenu
}

// ShowSyncFailedScreen displays the sync failed error screen
func ShowSyncFailedScreen() (string, int) {
	app.LogDebug("Showing sync failed screen")

	return ui.ShowMessage("Failed to synchronize catalog. Please check your internet connection and try again.", "3")
}

// HandleSyncFailedScreen processes the sync failed screen
func HandleSyncFailedScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleSyncFailedScreen called with exitCode: %d", exitCode)

	// Return to main menu after showing error message
	return app.ScreenMainMenu
}