// File: src/internal/ui/screens/backups.go
// Updated version for backup screens with proper ReadManifest calls

package screens

import (
	"fmt"
	"strings"
    "path/filepath"
    "os"
	"thememanager/internal/app"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
	"strconv"
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

	// Let ExportTheme generate the sequential backup name automatically
	// by passing an empty string
	err := themes.ExportTheme("")

	if err != nil {
		app.LogDebug("Error exporting theme: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error exporting theme: %s", err), "3")
		// Set empty selected item on error
		app.SetSelectedItem("")
	} else {
		// Get the list of backups to find the most recent one (highest number)
		backups, listErr := themes.ListBackups()
		if listErr != nil {
			app.LogDebug("Error listing backups after export: %v", listErr)
			app.SetSelectedItem("backup")
		} else {
			// Find the backup with the highest number
			highestNum := 0
			latestBackup := "backup1"
			for _, backup := range backups {
				if strings.HasPrefix(backup, "backup") {
					numStr := strings.TrimPrefix(backup, "backup")
					if num, err := strconv.Atoi(numStr); err == nil && num > highestNum {
						highestNum = num
						latestBackup = backup
					}
				}
			}
			app.SetSelectedItem(latestBackup)
		}
	}

	return app.ScreenThemeExported
}

// ShowThemeExportedScreen displays the theme exported success screen
func ShowThemeExportedScreen() (string, int) {
	app.LogDebug("Showing theme exported screen")

	exportedName := app.GetSelectedItem()
	exportPath := themes.GetBackupPath(exportedName) // Changed to GetBackupPath

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

	// Build theme data for gallery
	var backupDataList []map[string]string
	for _, backupName := range backupNames {
		backupPath := themes.GetBackupPath(backupName)

		// Create backup info
		backupInfo := map[string]string{
			"name": backupName,
			"is_valid": "true", // Backups are always considered valid for restore
		}

		// Try to read manifest - non-strict validation for display
		manifest, err := themes.ReadManifest(backupPath, false)
		if err == nil {
			backupInfo["author"] = manifest.Author
			backupInfo["description"] = manifest.Description
		} else {
			app.LogDebug("Warning: Failed to read manifest for backup %s: %v", backupName, err)
		}

		// Check if preview exists
		previewPath := filepath.Join(backupPath, themes.ThemePreviewFile)
		if _, err := os.Stat(previewPath); err == nil {
			backupInfo["preview"] = previewPath
		}

		backupDataList = append(backupDataList, backupInfo)
	}

	// Show backups gallery
	return ui.ShowThemeGallery(backupDataList, "Select Backup to Restore")
}

// HandleRestoreThemeScreen processes the backup selection
func HandleRestoreThemeScreen(selection string, exitCode int) app.Screen {
	app.LogDebug("HandleRestoreThemeScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	if exitCode == 0 && selection != "" {
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
	backupPath := themes.GetBackupPath(selectedBackup)

	// Try to read manifest for additional info - non-strict for display
	manifest, err := themes.ReadManifest(backupPath, false)

	var confirmMessage string
	if err == nil {
		// Format with manifest details
		confirmMessage = fmt.Sprintf("Restore theme from backup '%s' created by %s?",
			manifest.Name, manifest.Author)
	} else {
		// Simple confirmation without manifest details
		confirmMessage = fmt.Sprintf("Restore theme from backup '%s'?", selectedBackup)
	}

	return ui.ShowConfirmDialog(confirmMessage)
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