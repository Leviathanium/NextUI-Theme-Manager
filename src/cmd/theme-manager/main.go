// cmd/theme-manager/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"thememanager/internal/app"
	"thememanager/internal/logging"
	"thememanager/internal/themes"
	"thememanager/internal/ui/screens"
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

	logging.LogDebug("Starting main loop")

	// Set initial screen
	app.SetCurrentScreen(app.Screens.MainMenu)

	// Main application loop
	for {
		var selection string
		var exitCode int
		var nextScreen app.Screen

		// Log current screen
		currentScreen := app.GetCurrentScreen()
		logging.LogDebug("Current screen: %d", currentScreen)

		// Process current screen
		switch currentScreen {
		// Main menu
		case app.Screens.MainMenu:
			selection, exitCode = screens.MainMenuScreen()
			nextScreen = screens.HandleMainMenu(selection, exitCode)

		// Theme screens
		case app.Screens.Themes, app.Screens.ThemeGallery:
			selection, exitCode = screens.ThemeGalleryScreen()
			nextScreen = screens.HandleThemeGallery(selection, exitCode)

		case app.Screens.ThemeDownloadConfirm:
			selection, exitCode = screens.ThemeDownloadConfirmScreen()
			nextScreen = screens.HandleThemeDownloadConfirm(selection, exitCode)

		case app.Screens.ThemeDownloading:
			nextScreen = screens.ThemeDownloadingScreen()

		case app.Screens.ThemeApplyConfirm:
			selection, exitCode = screens.ThemeApplyConfirmScreen()
			nextScreen = screens.HandleThemeApplyConfirm(selection, exitCode)

		case app.Screens.ThemeApplying:
			nextScreen = screens.ThemeApplyingScreen()

		// Overlay screens
		case app.Screens.Overlays, app.Screens.OverlayGallery:
			selection, exitCode = screens.OverlayGalleryScreen()
			nextScreen = screens.HandleOverlayGallery(selection, exitCode)

		case app.Screens.OverlayDownloadConfirm:
			selection, exitCode = screens.OverlayDownloadConfirmScreen()
			nextScreen = screens.HandleOverlayDownloadConfirm(selection, exitCode)

		case app.Screens.OverlayDownloading:
			nextScreen = screens.OverlayDownloadingScreen()

		case app.Screens.OverlayApplyConfirm:
			selection, exitCode = screens.OverlayApplyConfirmScreen()
			nextScreen = screens.HandleOverlayApplyConfirm(selection, exitCode)

		case app.Screens.OverlayApplying:
			nextScreen = screens.OverlayApplyingScreen()

		// Sync catalog screen
		case app.Screens.SyncCatalog:
			nextScreen = screens.SyncCatalogScreen()

		// Backup screens
		case app.Screens.Backup, app.Screens.BackupMenu:
			selection, exitCode = screens.BackupMenuScreen()
			nextScreen = screens.HandleBackupMenu(selection, exitCode)

		case app.Screens.BackupThemeConfirm:
			selection, exitCode = screens.BackupThemeConfirmScreen()
			nextScreen = screens.HandleBackupThemeConfirm(selection, exitCode)

		case app.Screens.BackupThemeCreating:
			nextScreen = screens.BackupThemeCreatingScreen()

		case app.Screens.BackupOverlayConfirm:
			selection, exitCode = screens.BackupOverlayConfirmScreen()
			nextScreen = screens.HandleBackupOverlayConfirm(selection, exitCode)

		case app.Screens.BackupOverlayCreating:
			nextScreen = screens.BackupOverlayCreatingScreen()

		case app.Screens.BackupAutoToggle:
			selection, exitCode = screens.BackupAutoToggleScreen()
			nextScreen = screens.HandleBackupAutoToggle(selection, exitCode)

		// Revert screens
		case app.Screens.Revert, app.Screens.RevertMenu:
			selection, exitCode = screens.RevertMenuScreen()
			nextScreen = screens.HandleRevertMenu(selection, exitCode)

		case app.Screens.RevertThemeGallery:
			selection, exitCode = screens.RevertThemeGalleryScreen()
			nextScreen = screens.HandleRevertThemeGallery(selection, exitCode)

		case app.Screens.RevertThemeConfirm:
			selection, exitCode = screens.RevertThemeConfirmScreen()
			nextScreen = screens.HandleRevertThemeConfirm(selection, exitCode)

		case app.Screens.RevertThemeApplying:
			nextScreen = screens.RevertThemeApplyingScreen()

		case app.Screens.RevertOverlayGallery:
			selection, exitCode = screens.RevertOverlayGalleryScreen()
			nextScreen = screens.HandleRevertOverlayGallery(selection, exitCode)

		case app.Screens.RevertOverlayConfirm:
			selection, exitCode = screens.RevertOverlayConfirmScreen()
			nextScreen = screens.HandleRevertOverlayConfirm(selection, exitCode)

		case app.Screens.RevertOverlayApplying:
			nextScreen = screens.RevertOverlayApplyingScreen()

		// Purge screens
		case app.Screens.Purge, app.Screens.PurgeConfirm:
			selection, exitCode = screens.PurgeConfirmScreen()
			nextScreen = screens.HandlePurgeConfirm(selection, exitCode)

		case app.Screens.Purging:
			nextScreen = screens.PurgingScreen()

		default:
			logging.LogDebug("Unknown screen: %d, defaulting to main menu", currentScreen)
			nextScreen = app.Screens.MainMenu
		}

		// Validate next screen
		if nextScreen < app.Screens.MainMenu || int(nextScreen) > int(app.Screens.Purging) {
			logging.LogDebug("Invalid next screen: %d, defaulting to main menu", nextScreen)
			nextScreen = app.Screens.MainMenu
		}

		// Update the current screen
		app.SetCurrentScreen(nextScreen)
	}
}