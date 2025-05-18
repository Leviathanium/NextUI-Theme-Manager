package main

import (
	"fmt"
	"os"
	"runtime"

	"thememanager/internal/app"
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
			if app.IsLoggerInitialized() {
				app.LogDebug("PANIC: %v\n\nStack Trace:\n%s\n", r, stackTrace)
			}

			// Exit with error
			os.Exit(1)
		}
	}()

	// Initialize the logger
	defer app.CloseLogger()
	app.LogDebug("Application started")
	app.SetLoggerInitialized()

	// Initialize application
	if err := app.Initialize(); err != nil {
		app.LogDebug("Failed to initialize application: %v", err)
		return
	}

	// Initialize theme system
	if err := themes.EnsureDirectories(); err != nil {
		app.LogDebug("Warning: Failed to ensure theme directories: %v", err)
		// Continue anyway, some functionality might still work
	}

	app.LogDebug("Starting main loop")

	// Set initial screen
	app.SetCurrentScreen(app.ScreenMainMenu)

	// Main application loop
	for {
		var selection string
		var exitCode int
		var nextScreen app.Screen

		// Get current screen
		currentScreen := app.GetCurrentScreen()
		app.LogDebug("Current screen: %s", app.ScreenName(currentScreen))

		// Process current screen
		switch currentScreen {
		// Main Menu
		case app.ScreenMainMenu:
			selection, exitCode = screens.ShowMainMenu()
			nextScreen = screens.HandleMainMenu(selection, exitCode)

		// About Screen
		case app.ScreenAbout:
			selection, exitCode = screens.ShowAboutScreen()
			nextScreen = screens.HandleAboutScreen(selection, exitCode)

		// Apply Theme Screens
		case app.ScreenApplyTheme:
			selection, exitCode = screens.ShowApplyThemeScreen()
			nextScreen = screens.HandleApplyThemeScreen(selection, exitCode)

		case app.ScreenApplyThemeConfirm:
			selection, exitCode = screens.ShowApplyThemeConfirmScreen()
			nextScreen = screens.HandleApplyThemeConfirmScreen(selection, exitCode)

		case app.ScreenApplyingTheme:
			selection, exitCode = screens.ShowApplyingThemeScreen()
			nextScreen = screens.HandleApplyingThemeScreen(selection, exitCode)

		case app.ScreenThemeApplied:
			selection, exitCode = screens.ShowThemeAppliedScreen()
			nextScreen = screens.HandleThemeAppliedScreen(selection, exitCode)

		// Download Theme Screens
		case app.ScreenDownloadTheme:
			selection, exitCode = screens.ShowDownloadThemeScreen()
			nextScreen = screens.HandleDownloadThemeScreen(selection, exitCode)

		case app.ScreenCatalogNotSynced:
			selection, exitCode = screens.ShowCatalogNotSyncedScreen()
			nextScreen = screens.HandleCatalogNotSyncedScreen(selection, exitCode)

		case app.ScreenDownloadThemeConfirm:
			selection, exitCode = screens.ShowDownloadThemeConfirmScreen()
			nextScreen = screens.HandleDownloadThemeConfirmScreen(selection, exitCode)

		case app.ScreenDownloadingTheme:
			selection, exitCode = screens.ShowDownloadingThemeScreen()
			nextScreen = screens.HandleDownloadingThemeScreen(selection, exitCode)

		case app.ScreenThemeDownloaded:
			selection, exitCode = screens.ShowThemeDownloadedScreen()
			nextScreen = screens.HandleThemeDownloadedScreen(selection, exitCode)

		// Sync Catalog Screens
		case app.ScreenSyncCatalog:
			selection, exitCode = screens.ShowSyncCatalogScreen()
			nextScreen = screens.HandleSyncCatalogScreen(selection, exitCode)

		case app.ScreenSyncingCatalog:
			selection, exitCode = screens.ShowSyncingCatalogScreen()
			nextScreen = screens.HandleSyncingCatalogScreen(selection, exitCode)

		case app.ScreenSyncComplete:
			selection, exitCode = screens.ShowSyncCompleteScreen()
			nextScreen = screens.HandleSyncCompleteScreen(selection, exitCode)

		case app.ScreenSyncFailed:
			selection, exitCode = screens.ShowSyncFailedScreen()
			nextScreen = screens.HandleSyncFailedScreen(selection, exitCode)

		// Backups Menu Screens
		case app.ScreenBackupsMenu:
			selection, exitCode = screens.ShowBackupsMenuScreen()
			nextScreen = screens.HandleBackupsMenuScreen(selection, exitCode)

		// Export Theme Screens
		case app.ScreenExportTheme:
			selection, exitCode = screens.ShowExportThemeScreen()
			nextScreen = screens.HandleExportThemeScreen(selection, exitCode)

		case app.ScreenExportingTheme:
			selection, exitCode = screens.ShowExportingThemeScreen()
			nextScreen = screens.HandleExportingThemeScreen(selection, exitCode)

		case app.ScreenThemeExported:
			selection, exitCode = screens.ShowThemeExportedScreen()
			nextScreen = screens.HandleThemeExportedScreen(selection, exitCode)

		// Restore Theme Screens
		case app.ScreenRestoreTheme:
			selection, exitCode = screens.ShowRestoreThemeScreen()
			nextScreen = screens.HandleRestoreThemeScreen(selection, exitCode)

		case app.ScreenRestoreThemeConfirm:
			selection, exitCode = screens.ShowRestoreThemeConfirmScreen()
			nextScreen = screens.HandleRestoreThemeConfirmScreen(selection, exitCode)

		case app.ScreenRestoringTheme:
			selection, exitCode = screens.ShowRestoringThemeScreen()
			nextScreen = screens.HandleRestoringThemeScreen(selection, exitCode)

		case app.ScreenThemeRestored:
			selection, exitCode = screens.ShowThemeRestoredScreen()
			nextScreen = screens.HandleThemeRestoredScreen(selection, exitCode)

		default:
			app.LogDebug("Unknown screen: %s, defaulting to main menu", app.ScreenName(currentScreen))
			nextScreen = app.ScreenMainMenu
		}

		// Update the current screen
		app.SetCurrentScreen(nextScreen)

		// Special case for exit
		if exitCode == 1 || exitCode == 2 && currentScreen == app.ScreenMainMenu {
			app.LogDebug("User exited the application")
			os.Exit(0)
		}
	}
}