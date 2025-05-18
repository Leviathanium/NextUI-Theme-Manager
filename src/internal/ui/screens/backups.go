// Fixes for src/internal/ui/screens/backups.go

package screens

import (
	"fmt"
	"strings"
	"time"

	"thememanager/internal/app"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// ShowBackupsMenuScreen displays the backups menu screen
func ShowBackupsMenuScreen() (string, int) {
	app.LogDebug("Showing backups menu screen")

	menuItems := []string{
		"Export Theme",
		"Restore Theme",
	}

	return ui.ShowMenu(
		strings.Join(menuItems, "\n"),
		"Backups Menu",
		"--cancel-text", "BACK",
	)
}

// HandleBackupsMenuScreen processes the backups menu selection
func HandleBackupsMenuScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleBackupsMenuScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Export Theme":
			return app.ScreenExportTheme

		case "Restore Theme":
			return app.ScreenRestoreTheme

		default:
			return app.ScreenBackupsMenu
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User cancelled
		return app.ScreenMainMenu
	}

	return app.ScreenBackupsMenu
}

// ShowExportThemeScreen displays the export theme confirmation screen
func ShowExportThemeScreen() (string, int) {
	app.LogDebug("Showing export theme screen")

	// Check if system theme directory exists
	if !themes.SystemThemeExists() {
		ui.ShowMessage("No system theme found to export.", "3")
		return "", 1
	}

	return ui.ShowConfirmDialog("Export current system theme as backup?")
}

// HandleExportThemeScreen processes the export theme confirmation
func HandleExportThemeScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleExportThemeScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed - proceed to exporting
		return app.ScreenExportingTheme
	} else {
		// User cancelled
		return app.ScreenBackupsMenu
	}
}

// ShowExportingThemeScreen displays the export progress screen
func ShowExportingThemeScreen() (string, int) {
	app.LogDebug("Showing exporting theme screen")

	return ui.ShowMessage("Exporting current theme...", "2")
}

// HandleExportingThemeScreen processes the export operation
func HandleExportingThemeScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleExportingThemeScreen called with exitCode: %d", exitCode)

	// Generate timestamp for export name
	timestamp := time.Now().Format("20060102_150405")
	exportName := fmt.Sprintf("backup_%s", timestamp)

	// Export the theme
	err := themes.ExportTheme(exportName)

	if err != nil {
		app.LogDebug("Error exporting theme: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error exporting theme: %s", err), "3")
	} else {
		// Set the exported name for success screen
		app.SetSelectedItem(exportName)
	}

	return app.ScreenThemeExported
}

// ShowThemeExportedScreen displays the theme exported success screen
func ShowThemeExportedScreen() (string, int) {
	app.LogDebug("Showing theme exported screen")

	exportedName := app.GetSelectedItem()
	exportPath := themes.GetThemePath(exportedName)

	return ui.ShowMessage(
		fmt.Sprintf("Theme exported successfully to:\n%s", exportPath),
		"2",
	)
}

// HandleThemeExportedScreen processes the theme exported screen
func HandleThemeExportedScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleThemeExportedScreen called with exitCode: %d", exitCode)

	// Return to backups menu after showing success message
	return app.ScreenBackupsMenu
}

// ShowRestoreThemeScreen displays the theme selection screen for restoring
func ShowRestoreThemeScreen() (string, int) {
	app.LogDebug("Showing restore theme screen")

	// Get list of available backups
	backupNames, err := themes.ListBackups()
	if err != nil {
		app.LogDebug("Error listing backups: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error listing backups: %s", err), "3")
		return "", 1
	}

	// Check if we have any backups
	if len(backupNames) == 0 {
		ui.ShowMessage("No backups available. Please create a backup first.", "3")
		return "", 1
	}

	// Build menu items
	var menuItems []string
	for _, backupName := range backupNames {
		// Try to get author from manifest
		backupPath := themes.GetBackupPath(backupName)
		manifest, err := themes.ReadManifest(backupPath)

		if err == nil && manifest.Description != "" {
			menuItems = append(menuItems, fmt.Sprintf("%s - %s", backupName, manifest.Description))
		} else {
			menuItems = append(menuItems, backupName)
		}
	}

	return ui.ShowMenu(
		strings.Join(menuItems, "\n"),
		"Select Backup to Restore",
		"--cancel-text", "BACK",
	)
}

// HandleRestoreThemeScreen processes the backup selection
func HandleRestoreThemeScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleRestoreThemeScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected a backup
		app.SetSelectedItem(selection)
		return app.ScreenRestoreThemeConfirm
	} else if exitCode == 1 || exitCode == 2 {
		// User cancelled
		return app.ScreenBackupsMenu
	}

	return app.ScreenRestoreTheme
}

// ShowRestoreThemeConfirmScreen displays the confirmation screen for restoring
func ShowRestoreThemeConfirmScreen() (string, int) {
	app.LogDebug("Showing restore theme confirmation screen")

	selectedBackup := app.GetSelectedItem()
	return ui.ShowConfirmDialog("Restore theme from backup '" + selectedBackup + "'?")
}

// HandleRestoreThemeConfirmScreen processes the confirmation result
func HandleRestoreThemeConfirmScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleRestoreThemeConfirmScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection == "Yes" {
		// User confirmed - proceed to restoring
		return app.ScreenRestoringTheme
	} else {
		// User cancelled
		return app.ScreenRestoreTheme
	}
}

// ShowRestoringThemeScreen displays the theme restoring progress screen
func ShowRestoringThemeScreen() (string, int) {
	app.LogDebug("Showing restoring theme screen")

	selectedBackup := app.GetSelectedItem()
	return ui.ShowMessage("Restoring theme from '" + selectedBackup + "'...", "2")
}

// HandleRestoringThemeScreen processes the restoring operation
func HandleRestoringThemeScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleRestoringThemeScreen called with exitCode: %d", exitCode)

	selectedBackup := app.GetSelectedItem()

	// Extract backup name from selection (remove description part if present)
	backupName := selectedBackup
	if idx := strings.Index(backupName, " - "); idx > 0 {
		backupName = backupName[:idx]
	}

	// Restore the backup
	err := themes.RestoreBackup(backupName)

	if err != nil {
		app.LogDebug("Error restoring backup: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error restoring backup: %s", err), "3")
	}

	return app.ScreenThemeRestored
}

// ShowThemeRestoredScreen displays the theme restored success screen
func ShowThemeRestoredScreen() (string, int) {
	app.LogDebug("Showing theme restored screen")

	selectedBackup := app.GetSelectedItem()
	return ui.ShowMessage("Theme restored successfully from '" + selectedBackup + "'!", "2")
}

// HandleThemeRestoredScreen processes the success screen
func HandleThemeRestoredScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleThemeRestoredScreen called with exitCode: %d", exitCode)

	// Return to main menu after showing success message
	return app.ScreenMainMenu
}