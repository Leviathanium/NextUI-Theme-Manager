// internal/ui/screens/backup.go
package screens

import (
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// BackupMenuScreen displays the backup menu screen
func BackupMenuScreen() (string, int) {
	logging.LogDebug("Showing backup menu screen")

	menuItems := []string{
		"Backup Theme",
		"Backup Overlays",
		"Auto-Backup",
	}

	return ui.DisplayMinUiList(
		strings.Join(menuItems, "\n"),
		"text",
		"Backup Menu",
	)
}

// HandleBackupMenu processes backup menu selection
func HandleBackupMenu(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleBackupMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Backup Theme":
			return app.Screens.BackupThemeConfirm

		case "Backup Overlays":
			return app.Screens.BackupOverlayConfirm

		case "Auto-Backup":
			return app.Screens.BackupAutoToggle
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/back
		return app.Screens.MainMenu
	}

	return app.Screens.BackupMenu
}

// BackupThemeConfirmScreen displays the theme backup confirmation screen
func BackupThemeConfirmScreen() (string, int) {
	logging.LogDebug("Showing theme backup confirmation screen")
	return ui.ShowConfirmDialog("Create theme backup?")
}

// HandleBackupThemeConfirm processes theme backup confirmation
func HandleBackupThemeConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleBackupThemeConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed backup
		return app.Screens.BackupThemeCreating
	}

	// User cancelled
	return app.Screens.SettingsMenu
}

// BackupThemeCreatingScreen handles the theme backup creation process
func BackupThemeCreatingScreen() app.Screen {
	logging.LogDebug("Processing theme backup creation")

	// Show creating message and perform operation
	err := ui.ShowMessageWithOperation(
		"Creating theme backup...",
		func() error {
			return themes.CreateThemeBackup("manual")
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error creating backup: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Theme backup created successfully!", "2")
	}

	return app.Screens.SettingsMenu
}

// BackupOverlayConfirmScreen displays the overlay backup confirmation screen
func BackupOverlayConfirmScreen() (string, int) {
	logging.LogDebug("Showing overlay backup confirmation screen")
	return ui.ShowConfirmDialog("Create overlay backup?")
}

// HandleBackupOverlayConfirm processes overlay backup confirmation
func HandleBackupOverlayConfirm(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleBackupOverlayConfirm called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed backup
		return app.Screens.BackupOverlayCreating
	}

	// User cancelled
	return app.Screens.SettingsMenu
}

// BackupOverlayCreatingScreen handles the overlay backup creation process
func BackupOverlayCreatingScreen() app.Screen {
	logging.LogDebug("Processing overlay backup creation")

	// Show creating message and perform operation
	err := ui.ShowMessageWithOperation(
		"Creating overlay backup...",
		func() error {
			return themes.CreateOverlayBackup("manual")
		},
	)

	// Check result
	if err != nil {
		ui.ShowMessage("Error creating backup: " + err.Error(), "3")
	} else {
		ui.ShowMessage("Overlay backup created successfully!", "2")
	}

	return app.Screens.SettingsMenu
}

// BackupAutoToggleScreen displays the auto-backup toggle screen
func BackupAutoToggleScreen() (string, int) {
	logging.LogDebug("Showing auto-backup toggle screen")

	// Pass current auto-backup status to set default selection
	return ui.ShowConfirmDialog("Enable Auto-Backup?", app.GetAutoBackup())
}

// HandleBackupAutoToggle processes auto-backup toggle
func HandleBackupAutoToggle(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleBackupAutoToggle called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// Toggle auto-backup setting
		enabled := selection == "Yes"
		app.SetAutoBackup(enabled)

		// Show confirmation message
		if enabled {
			ui.ShowMessage("Auto-backup enabled", "2")
		} else {
			ui.ShowMessage("Auto-backup disabled", "2")
		}
	}

	// Return to settings menu
	return app.Screens.SettingsMenu
}