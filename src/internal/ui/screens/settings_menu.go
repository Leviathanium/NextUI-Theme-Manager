// internal/ui/screens/settings_menu.go
package screens

import (
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/ui"
)

// SettingsMenuScreen displays the settings menu screen
func SettingsMenuScreen() (string, int) {
	logging.LogDebug("Showing settings menu screen")

	menuItems := []string{
		"Restore Backup",
		"Backup Theme",
		"Backup Overlays",
		"Auto-Backup",
		"Purge",
	}

	return ui.DisplayMinUiList(
		strings.Join(menuItems, "\n"),
		"text",
		"Settings Menu",
	)
}

// HandleSettingsMenu processes settings menu selection
func HandleSettingsMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleSettingsMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Restore Backup":
			logging.LogDebug("Selected Restore Backup")
			return app.Screens.RestoreMenu

		case "Backup Theme":
			logging.LogDebug("Selected Backup Theme")
			return app.Screens.BackupThemeConfirm

		case "Backup Overlays":
			logging.LogDebug("Selected Backup Overlays")
			return app.Screens.BackupOverlayConfirm

		case "Auto-Backup":
			logging.LogDebug("Selected Auto-Backup")
			return app.Screens.BackupAutoToggle

		case "Purge":
			logging.LogDebug("Selected Purge")
			return app.Screens.PurgeConfirm

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.SettingsMenu
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/back
		return app.Screens.MainMenu
	}

	return app.Screens.SettingsMenu
}

// RestoreMenuScreen displays the restore menu screen
func RestoreMenuScreen() (string, int) {
	logging.LogDebug("Showing restore menu screen")

	menuItems := []string{
		"Restore Theme",
		"Restore Overlays",
	}

	return ui.DisplayMinUiList(
		strings.Join(menuItems, "\n"),
		"text",
		"Restore Menu",
	)
}

// HandleRestoreMenu processes restore menu selection
func HandleRestoreMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleRestoreMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Restore Theme":
			logging.LogDebug("Selected Restore Theme")
			return app.Screens.RestoreThemeGallery

		case "Restore Overlays":
			logging.LogDebug("Selected Restore Overlays")
			return app.Screens.RestoreOverlayGallery

		default:
			logging.LogDebug("Unknown selection: %s", selection)
			return app.Screens.RestoreMenu
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/back
		return app.Screens.SettingsMenu
	}

	return app.Screens.RestoreMenu
}