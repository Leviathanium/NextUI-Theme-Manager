// main.go - Theme Manager Application
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// Screen identifiers
const (
	// Main menu screens
	ScreenMainMenu = iota
	ScreenThemes
	ScreenOverlays
	ScreenSyncCatalog
	ScreenBackup
	ScreenRevert
	ScreenPurge

	// Theme sub-screens
	ScreenThemeGallery
	ScreenThemeDownloadConfirm
	ScreenThemeDownloading
	ScreenThemeApplyConfirm
	ScreenThemeApplying

	// Overlay sub-screens
	ScreenOverlayGallery
	ScreenOverlayDownloadConfirm
	ScreenOverlayDownloading
	ScreenOverlayApplyConfirm
	ScreenOverlayApplying

	// Backup sub-screens
	ScreenBackupMenu
	ScreenBackupThemeConfirm
	ScreenBackupThemeCreating
	ScreenBackupOverlayConfirm
	ScreenBackupOverlayCreating
	ScreenBackupAutoToggle

	// Revert sub-screens
	ScreenRevertMenu
	ScreenRevertThemeGallery
	ScreenRevertThemeConfirm
	ScreenRevertThemeApplying
	ScreenRevertOverlayGallery
	ScreenRevertOverlayConfirm
	ScreenRevertOverlayApplying

	// Purge screens
	ScreenPurgeConfirm
	ScreenPurging
)

// Global application state
var (
	currentScreen int
	selectedItem string
	autoBackupEnabled bool
)

func main() {
	// Recover from panics
	defer func() {
		if r := recover(); r != nil {
			// Get stack trace
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			stackTrace := string(buf[:n])

			// Log the panic
			fmt.Fprintf(os.Stderr, "PANIC: %v\n\nStack Trace:\n%s\n", r, stackTrace)

			// Also try to log to file if possible
			if logging.IsLoggerInitialized() {
				logging.LogDebug("PANIC: %v\n\nStack Trace:\n%s\n", r, stackTrace)
			}

			// Exit with error
			os.Exit(1)
		}
	}()

	// Initialize the logger
	defer logging.CloseLogger()

	logging.LogDebug("Application started")
	logging.SetLoggerInitialized() // Explicitly mark logger as initialized

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return
	}

	// Check if minui-list exists in the application directory
	minuiListPath := filepath.Join(cwd, "minui-list")
	_, err = os.Stat(minuiListPath)
	if err != nil {
		logging.LogDebug("minui-list not found at %s: %v", minuiListPath, err)
		return
	}

	// Check if minui-presenter exists in the application directory
	minuiPresenterPath := filepath.Join(cwd, "minui-presenter")
	_, err = os.Stat(minuiPresenterPath)
	if err != nil {
		logging.LogDebug("minui-presenter not found at %s: %v", minuiPresenterPath, err)
		return
	}

	// Initialize application
	if err := app.Initialize(); err != nil {
		logging.LogDebug("Failed to initialize application: %v", err)
		return
	}

	// Create necessary directory structure
	if err := themes.EnsureDirectoryStructure(); err != nil {
		logging.LogDebug("Warning: Could not create theme directories: %v", err)
	}

	// Load auto-backup setting
	autoBackupEnabled = app.GetAutoBackupSetting()

	logging.LogDebug("Starting main loop")

	// Set initial screen
	currentScreen = ScreenMainMenu

	// Main application loop
	for {
		var selection string
		var exitCode int
		var nextScreen int

		// Log current screen
		logging.LogDebug("Current screen: %d", currentScreen)

		// Process current screen
		switch currentScreen {
		case ScreenMainMenu:
			selection, exitCode = showMainMenu()
			nextScreen = handleMainMenu(selection, exitCode)

		// Theme screens
		case ScreenThemes:
			nextScreen = ScreenThemeGallery
		case ScreenThemeGallery:
			selection, exitCode = themes.ShowThemeGallery()
			nextScreen = handleThemeGallery(selection, exitCode)
		case ScreenThemeDownloadConfirm:
			selection, exitCode = ui.ShowConfirmDialog(fmt.Sprintf("Download theme '%s'?", selectedItem))
			nextScreen = handleThemeDownloadConfirm(selection, exitCode)
		case ScreenThemeDownloading:
			err := ui.ShowMessageWithOperation(
				fmt.Sprintf("Downloading theme '%s'...", selectedItem),
				func() error {
					return themes.DownloadTheme(selectedItem)
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error downloading theme: %s", err), "3")
				nextScreen = ScreenThemeGallery
			} else {
				ui.ShowMessage("Theme downloaded successfully!", "2")
				nextScreen = ScreenThemeApplyConfirm
			}
		case ScreenThemeApplyConfirm:
			selection, exitCode = ui.ShowConfirmDialog(fmt.Sprintf("Apply theme '%s'?", selectedItem))
			nextScreen = handleThemeApplyConfirm(selection, exitCode)
		case ScreenThemeApplying:
			// Create backup if auto-backup is enabled
			if autoBackupEnabled {
				err := themes.CreateThemeBackup("auto")
				if err != nil {
					logging.LogDebug("Error creating auto-backup: %v", err)
				}
			}

			// Apply theme with operation message
			err := ui.ShowMessageWithOperation(
				fmt.Sprintf("Applying theme '%s'...", selectedItem),
				func() error {
					return themes.ApplyTheme(selectedItem)
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error applying theme: %s", err), "3")
			} else {
				ui.ShowMessage("Theme applied successfully!", "2")
			}
			nextScreen = ScreenMainMenu

		// Overlay screens
		case ScreenOverlays:
			nextScreen = ScreenOverlayGallery
		case ScreenOverlayGallery:
			selection, exitCode = themes.ShowOverlayGallery()
			nextScreen = handleOverlayGallery(selection, exitCode)
		case ScreenOverlayDownloadConfirm:
			selection, exitCode = ui.ShowConfirmDialog(fmt.Sprintf("Download overlays '%s'?", selectedItem))
			nextScreen = handleOverlayDownloadConfirm(selection, exitCode)
		case ScreenOverlayDownloading:
			err := ui.ShowMessageWithOperation(
				fmt.Sprintf("Downloading overlays '%s'...", selectedItem),
				func() error {
					return themes.DownloadOverlay(selectedItem)
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error downloading overlays: %s", err), "3")
				nextScreen = ScreenOverlayGallery
			} else {
				ui.ShowMessage("Overlays downloaded successfully!", "2")
				nextScreen = ScreenOverlayApplyConfirm
			}
		case ScreenOverlayApplyConfirm:
			selection, exitCode = ui.ShowConfirmDialog(fmt.Sprintf("Apply overlays '%s'?", selectedItem))
			nextScreen = handleOverlayApplyConfirm(selection, exitCode)
		case ScreenOverlayApplying:
			// Create backup if auto-backup is enabled
			if autoBackupEnabled {
				err := themes.CreateOverlayBackup("auto")
				if err != nil {
					logging.LogDebug("Error creating auto-backup: %v", err)
				}
			}

			// Apply overlay with operation message
			err := ui.ShowMessageWithOperation(
				fmt.Sprintf("Applying overlays '%s'...", selectedItem),
				func() error {
					return themes.ApplyOverlay(selectedItem)
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error applying overlays: %s", err), "3")
			} else {
				ui.ShowMessage("Overlays applied successfully!", "2")
			}
			nextScreen = ScreenMainMenu

		// Sync catalog screen
		case ScreenSyncCatalog:
			err := ui.ShowMessageWithOperation(
				"Syncing catalog...",
				func() error {
					return themes.SyncCatalog()
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error syncing catalog: %s", err), "3")
			} else {
				ui.ShowMessage("Catalog synced successfully!", "2")
			}
			nextScreen = ScreenMainMenu

		// Backup screens
		case ScreenBackup:
			nextScreen = ScreenBackupMenu
		case ScreenBackupMenu:
			selection, exitCode = showBackupMenu()
			nextScreen = handleBackupMenu(selection, exitCode)
		case ScreenBackupThemeConfirm:
			selection, exitCode = ui.ShowConfirmDialog("Create theme backup?")
			nextScreen = handleBackupThemeConfirm(selection, exitCode)
		case ScreenBackupThemeCreating:
			err := ui.ShowMessageWithOperation(
				"Creating theme backup...",
				func() error {
					return themes.CreateThemeBackup("manual")
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error creating backup: %s", err), "3")
			} else {
				ui.ShowMessage("Theme backup created successfully!", "2")
			}
			nextScreen = ScreenBackupMenu
		case ScreenBackupOverlayConfirm:
			selection, exitCode = ui.ShowConfirmDialog("Create overlay backup?")
			nextScreen = handleBackupOverlayConfirm(selection, exitCode)
		case ScreenBackupOverlayCreating:
			err := ui.ShowMessageWithOperation(
				"Creating overlay backup...",
				func() error {
					return themes.CreateOverlayBackup("manual")
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error creating backup: %s", err), "3")
			} else {
				ui.ShowMessage("Overlay backup created successfully!", "2")
			}
			nextScreen = ScreenBackupMenu
		case ScreenBackupAutoToggle:
			selection, exitCode = ui.ShowConfirmDialog("Enable Auto-Backup?", autoBackupEnabled)
			nextScreen = handleBackupAutoToggle(selection, exitCode)

		// Revert screens
		case ScreenRevert:
			nextScreen = ScreenRevertMenu
		case ScreenRevertMenu:
			selection, exitCode = showRevertMenu()
			nextScreen = handleRevertMenu(selection, exitCode)
		case ScreenRevertThemeGallery:
			selection, exitCode = themes.ShowThemeBackupGallery()
			nextScreen = handleRevertThemeGallery(selection, exitCode)
		case ScreenRevertThemeConfirm:
			selection, exitCode = ui.ShowConfirmDialog(fmt.Sprintf("Revert to theme backup '%s'?", selectedItem))
			nextScreen = handleRevertThemeConfirm(selection, exitCode)
		case ScreenRevertThemeApplying:
			err := ui.ShowMessageWithOperation(
				"Reverting from theme backup...",
				func() error {
					return themes.RevertThemeFromBackup(selectedItem)
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error reverting from backup: %s", err), "3")
			} else {
				ui.ShowMessage("Theme reverted successfully!", "2")
			}
			nextScreen = ScreenMainMenu
		case ScreenRevertOverlayGallery:
			selection, exitCode = themes.ShowOverlayBackupGallery()
			nextScreen = handleRevertOverlayGallery(selection, exitCode)
		case ScreenRevertOverlayConfirm:
			selection, exitCode = ui.ShowConfirmDialog(fmt.Sprintf("Revert to overlay backup '%s'?", selectedItem))
			nextScreen = handleRevertOverlayConfirm(selection, exitCode)
		case ScreenRevertOverlayApplying:
			err := ui.ShowMessageWithOperation(
				"Reverting from overlay backup...",
				func() error {
					return themes.RevertOverlayFromBackup(selectedItem)
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error reverting from backup: %s", err), "3")
			} else {
				ui.ShowMessage("Overlays reverted successfully!", "2")
			}
			nextScreen = ScreenMainMenu

		// Purge screens
		case ScreenPurge:
			nextScreen = ScreenPurgeConfirm
		case ScreenPurgeConfirm:
			selection, exitCode = ui.ShowConfirmDialog("WARNING: Erase everything?")
			nextScreen = handlePurgeConfirm(selection, exitCode)
		case ScreenPurging:
			err := ui.ShowMessageWithOperation(
				"Purging...",
				func() error {
					return themes.PurgeAll()
				},
			)
			if err != nil {
				ui.ShowMessage(fmt.Sprintf("Error during purge: %s", err), "3")
			} else {
				ui.ShowMessage("Purge complete!", "2")
			}
			nextScreen = ScreenMainMenu

		default:
			logging.LogDebug("Unknown screen: %d, defaulting to main menu", currentScreen)
			nextScreen = ScreenMainMenu
		}

		// Update the current screen
		currentScreen = nextScreen
	}
}

// showMainMenu displays the main menu
func showMainMenu() (string, int) {
	menuItems := []string{
		"Themes",
		"Overlays",
		"Sync Catalog",
		"Backup",
		"Revert",
		"Purge",
	}

	return ui.DisplayMinUiList(
		strings.Join(menuItems, "\n"),
		"text",
		"Theme Manager",
		"--cancel-text", "QUIT",
	)
}

// handleMainMenu processes main menu selection
func handleMainMenu(selection string, exitCode int) int {
	if exitCode == 0 {
		switch selection {
		case "Themes":
			return ScreenThemes
		case "Overlays":
			return ScreenOverlays
		case "Sync Catalog":
			return ScreenSyncCatalog
		case "Backup":
			return ScreenBackup
		case "Revert":
			return ScreenRevert
		case "Purge":
			return ScreenPurge
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/quit button
		os.Exit(0)
	}

	return ScreenMainMenu
}

// handleThemeGallery processes theme gallery selection
func handleThemeGallery(selection string, exitCode int) int {
	if exitCode == 0 && selection != "" {
		selectedItem = selection

		// Check if theme is already downloaded
		if themes.IsThemeDownloaded(selection) {
			return ScreenThemeApplyConfirm
		}

		// Not downloaded, ask to download
		return ScreenThemeDownloadConfirm
	}

	// Return to main menu on cancel
	return ScreenMainMenu
}

// handleThemeDownloadConfirm processes theme download confirmation
func handleThemeDownloadConfirm(selection string, exitCode int) int {
	if exitCode == 0 && selection == "Yes" {
		return ScreenThemeDownloading
	}

	// Return to gallery on cancel
	return ScreenThemeGallery
}

// handleThemeApplyConfirm processes theme apply confirmation
func handleThemeApplyConfirm(selection string, exitCode int) int {
	if exitCode == 0 && selection == "Yes" {
		return ScreenThemeApplying
	}

	// Return to gallery on cancel
	return ScreenThemeGallery
}

// handleOverlayGallery processes overlay gallery selection
func handleOverlayGallery(selection string, exitCode int) int {
	if exitCode == 0 && selection != "" {
		selectedItem = selection

		// Check if overlay is already downloaded
		if themes.IsOverlayDownloaded(selection) {
			return ScreenOverlayApplyConfirm
		}

		// Not downloaded, ask to download
		return ScreenOverlayDownloadConfirm
	}

	// Return to main menu on cancel
	return ScreenMainMenu
}

// handleOverlayDownloadConfirm processes overlay download confirmation
func handleOverlayDownloadConfirm(selection string, exitCode int) int {
	if exitCode == 0 && selection == "Yes" {
		return ScreenOverlayDownloading
	}

	// Return to gallery on cancel
	return ScreenOverlayGallery
}

// handleOverlayApplyConfirm processes overlay apply confirmation
func handleOverlayApplyConfirm(selection string, exitCode int) int {
	if exitCode == 0 && selection == "Yes" {
		return ScreenOverlayApplying
	}

	// Return to gallery on cancel
	return ScreenOverlayGallery
}

// showBackupMenu displays the backup menu
func showBackupMenu() (string, int) {
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

// handleBackupMenu processes backup menu selection
func handleBackupMenu(selection string, exitCode int) int {
	if exitCode == 0 {
		switch selection {
		case "Backup Theme":
			return ScreenBackupThemeConfirm
		case "Backup Overlays":
			return ScreenBackupOverlayConfirm
		case "Auto-Backup":
			return ScreenBackupAutoToggle
		}
	} else if exitCode == 1 || exitCode == 2 {
		// Return to main menu on cancel
		return ScreenMainMenu
	}

	return ScreenBackupMenu
}

// handleBackupThemeConfirm processes theme backup confirmation
func handleBackupThemeConfirm(selection string, exitCode int) int {
	if exitCode == 0 && selection == "Yes" {
		return ScreenBackupThemeCreating
	}

	// Return to backup menu on cancel
	return ScreenBackupMenu
}

// handleBackupOverlayConfirm processes overlay backup confirmation
func handleBackupOverlayConfirm(selection string, exitCode int) int {
	if exitCode == 0 && selection == "Yes" {
		return ScreenBackupOverlayCreating
	}

	// Return to backup menu on cancel
	return ScreenBackupMenu
}

// handleBackupAutoToggle processes auto-backup toggle
func handleBackupAutoToggle(selection string, exitCode int) int {
	if exitCode == 0 {
		// Toggle auto-backup setting
		autoBackupEnabled = selection == "Yes"
		app.SetAutoBackupSetting(autoBackupEnabled)

		// Show confirmation message
		if autoBackupEnabled {
			ui.ShowMessage("Auto-backup enabled", "2")
		} else {
			ui.ShowMessage("Auto-backup disabled", "2")
		}
	}

	// Return to backup menu
	return ScreenBackupMenu
}

// showRevertMenu displays the revert menu
func showRevertMenu() (string, int) {
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

// handleRevertMenu processes revert menu selection
func handleRevertMenu(selection string, exitCode int) int {
	if exitCode == 0 {
		switch selection {
		case "Revert Theme":
			return ScreenRevertThemeGallery
		case "Revert Overlays":
			return ScreenRevertOverlayGallery
		}
	} else if exitCode == 1 || exitCode == 2 {
		// Return to main menu on cancel
		return ScreenMainMenu
	}

	return ScreenRevertMenu
}

// handleRevertThemeGallery processes theme backup gallery selection
func handleRevertThemeGallery(selection string, exitCode int) int {
	if exitCode == 0 && selection != "" {
		selectedItem = selection
		return ScreenRevertThemeConfirm
	}

	// Return to revert menu on cancel
	return ScreenRevertMenu
}

// handleRevertThemeConfirm processes theme revert confirmation
func handleRevertThemeConfirm(selection string, exitCode int) int {
	if exitCode == 0 && selection == "Yes" {
		return ScreenRevertThemeApplying
	}

	// Return to backup gallery on cancel
	return ScreenRevertThemeGallery
}

// handleRevertOverlayGallery processes overlay backup gallery selection
func handleRevertOverlayGallery(selection string, exitCode int) int {
	if exitCode == 0 && selection != "" {
		selectedItem = selection
		return ScreenRevertOverlayConfirm
	}

	// Return to revert menu on cancel
	return ScreenRevertMenu
}

// handleRevertOverlayConfirm processes overlay revert confirmation
func handleRevertOverlayConfirm(selection string, exitCode int) int {
	if exitCode == 0 && selection == "Yes" {
		return ScreenRevertOverlayApplying
	}

	// Return to backup gallery on cancel
	return ScreenRevertOverlayGallery
}

// handlePurgeConfirm processes purge confirmation
func handlePurgeConfirm(selection string, exitCode int) int {
	if exitCode == 0 && selection == "Yes" {
		return ScreenPurging
	}

	// Return to main menu on cancel
	return ScreenMainMenu
}