// internal/ui/screens/overlays_menu.go
package screens

import (
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/ui"
)

// OverlaysMenuScreen displays the overlays submenu screen
func OverlaysMenuScreen() (string, int) {
	logging.LogDebug("Showing overlays menu screen")

	menuItems := []string{
		"Installed Overlays",
		"Download Overlays",
	}

	return ui.DisplayMinUiList(
		strings.Join(menuItems, "\n"),
		"text",
		"Overlays Menu",
	)
}

// HandleOverlaysMenu processes overlays menu selection
func HandleOverlaysMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleOverlaysMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Installed Overlays":
			logging.LogDebug("Selected Installed Overlays")
			return app.Screens.InstalledOverlays

		case "Download Overlays":
			logging.LogDebug("Selected Download Overlays")
			return app.Screens.DownloadOverlays

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.OverlaysMenu
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/back
		return app.Screens.MainMenu
	}

	return app.Screens.OverlaysMenu
}

// InstalledOverlaysScreen displays the installed overlays screen
func InstalledOverlaysScreen() (string, int) {
	logging.LogDebug("Showing installed overlays screen")
	return themes.ShowInstalledOverlays()
}

// HandleInstalledOverlays processes installed overlays selection
func HandleInstalledOverlays(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleInstalledOverlays called with selection: '%s', exitCode: %d", selection, exitCode)

	// When an overlay is selected, we'll set it as the current selection and proceed to application
	if exitCode == 0 && selection != "" {
		app.SetSelectedItem(selection)
		return app.Screens.OverlayApplyConfirm
	}

	// User pressed back
	return app.Screens.OverlaysMenu
}