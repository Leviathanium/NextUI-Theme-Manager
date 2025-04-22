// src/internal/ui/screens/sync_screens.go
// Implementation of theme and component sync screens

package screens

import (
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/themes"
	"nextui-themes/internal/ui"
)

// SyncComponentsScreen displays the sync components screen
func SyncComponentsScreen() (string, int) {
	componentType := app.GetSelectedComponentType()

	// Simple confirmation message
	message := fmt.Sprintf("Sync %s catalog from %s?\nThis will download the latest component catalog.",
		componentType, themes.RepoConfig.URL)
	options := []string{
		"Yes",
		"No",
	}

	return ui.DisplayMinUiList(strings.Join(options, "\n"), "text", message)
}

// HandleSyncComponents processes the user's choice to sync components
func HandleSyncComponents(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleSyncComponents called with selection: '%s', exitCode: %d", selection, exitCode)
	componentType := app.GetSelectedComponentType()

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Perform component sync
			logging.LogDebug("Starting component catalog sync for %s", componentType)

			// Get default sync options
			options := themes.GetDefaultSyncOptions()

			// Sync catalog
			if err := themes.SyncThemeCatalog(options); err != nil {
				logging.LogDebug("Error syncing component catalog: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				logging.LogDebug("Component catalog sync completed successfully")
				ui.ShowMessage(fmt.Sprintf("%s catalog synced successfully!", componentType), "2")
			}
		}
		// Return to component options
		return app.Screens.ComponentOptions

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.ComponentOptions
	}

	return app.Screens.SyncComponents
}