// internal/ui/screens/sync.go
package screens

import (
	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// SyncCatalogScreen handles the catalog synchronization process
func SyncCatalogScreen() app.Screen {
	logging.LogDebug("Processing catalog synchronization")

	// Show syncing message and perform operation
	err := ui.ShowMessageWithOperation(
		"Syncing catalog...",
		func() error {
			return themes.SyncCatalog()
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error syncing catalog: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Catalog synced successfully!", "2")
	}

	return app.Screens.MainMenu
}