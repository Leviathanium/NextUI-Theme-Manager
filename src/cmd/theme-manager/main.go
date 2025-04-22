// src/cmd/theme-manager/main.go
// Main entry point for the NextUI Theme Manager application

package main

import (
	"os"
	"path/filepath"
	"runtime"
	"fmt"
	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui/screens"
	"nextui-themes/internal/themes"
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

	// Create theme directory structure
	if err := themes.EnsureThemeDirectoryStructure(); err != nil {
		logging.LogDebug("Warning: Could not create theme directories: %v", err)
	}

	logging.LogDebug("Starting main loop")

	// Main application loop
	for {
		var selection string
		var exitCode int
		var nextScreen app.Screen

		// Log current screen
		currentScreen := app.GetCurrentScreen()
		logging.LogDebug("Current screen: %d", currentScreen)

		// Ensure screen value is valid
		if currentScreen < app.Screens.MainMenu || currentScreen > app.Screens.DeconstructConfirm {
			logging.LogDebug("CRITICAL ERROR: Invalid screen value: %d, resetting to MainMenu", currentScreen)
			app.SetCurrentScreen(app.Screens.MainMenu)
			continue
		}

		// Process current screen
		switch currentScreen {
		case app.Screens.MainMenu:
			logging.LogDebug("Showing main menu")
			selection, exitCode = screens.MainMenuScreen()
			nextScreen = screens.HandleMainMenu(selection, exitCode)

		// New and renamed screens
		case app.Screens.InstalledThemes:
			logging.LogDebug("Showing installed themes screen")
			selection, exitCode = screens.InstalledThemesScreen()
			nextScreen = screens.HandleInstalledThemes(selection, exitCode)

		case app.Screens.DownloadThemes:
			logging.LogDebug("Showing download themes screen")
			selection, exitCode = screens.DownloadThemesScreen()
			nextScreen = screens.HandleDownloadThemes(selection, exitCode)

		case app.Screens.SyncCatalog:
			logging.LogDebug("Showing sync catalog screen")
			selection, exitCode = screens.SyncCatalogScreen()
			nextScreen = screens.HandleSyncCatalog(selection, exitCode)

		// Original screens
		case app.Screens.ThemeImport:
			logging.LogDebug("Showing theme import selection")
			selection, exitCode = screens.ThemeImportScreen()
			nextScreen = screens.HandleThemeImport(selection, exitCode)

		case app.Screens.ThemeImportConfirm:
			logging.LogDebug("Showing theme import confirmation")
			selection, exitCode = screens.ThemeImportConfirmScreen()
			nextScreen = screens.HandleThemeImportConfirm(selection, exitCode)

		case app.Screens.ThemeExport:
			logging.LogDebug("Showing theme export screen")
			selection, exitCode = screens.ThemeExportScreen()
			nextScreen = screens.HandleThemeExport(selection, exitCode)

		case app.Screens.ComponentsMenu:
			logging.LogDebug("Showing components menu screen")
			selection, exitCode = screens.ComponentsMenuScreen()
			nextScreen = screens.HandleComponentsMenu(selection, exitCode)

		case app.Screens.ComponentOptions:
			logging.LogDebug("Showing component options screen")
			selection, exitCode = screens.ComponentOptionsScreen()
			nextScreen = screens.HandleComponentOptions(selection, exitCode)

		// New component screens
		case app.Screens.InstalledComponents:
			logging.LogDebug("Showing installed components screen")
			selection, exitCode = screens.InstalledComponentsScreen()
			nextScreen = screens.HandleInstalledComponents(selection, exitCode)

		case app.Screens.DownloadComponents:
			logging.LogDebug("Showing download components screen")
			selection, exitCode = screens.DownloadComponentsScreen()
			nextScreen = screens.HandleDownloadComponents(selection, exitCode)

		// Original component-related screens
		case app.Screens.SyncComponents:
			logging.LogDebug("Showing sync components screen")
			selection, exitCode = screens.SyncComponentsScreen()
			nextScreen = screens.HandleSyncComponents(selection, exitCode)

		case app.Screens.ExportComponent:
			logging.LogDebug("Showing export component screen")
			selection, exitCode = screens.ExportComponentScreen()
			nextScreen = screens.HandleExportComponent(selection, exitCode)

		case app.Screens.Deconstruction:
			logging.LogDebug("Showing deconstruction screen")
			selection, exitCode = screens.DeconstructionScreen()
			nextScreen = screens.HandleDeconstruction(selection, exitCode)

		case app.Screens.DeconstructConfirm:
			logging.LogDebug("Showing deconstruction confirmation screen")
			selection, exitCode = screens.DeconstructConfirmScreen()
			nextScreen = screens.HandleDeconstructConfirm(selection, exitCode)

		default:
			logging.LogDebug("Unknown screen type: %d, defaulting to MainMenu", currentScreen)
			nextScreen = app.Screens.MainMenu
		}

		// Add extra debug logging
		logging.LogDebug("Current screen: %d, Next screen: %d", currentScreen, nextScreen)

		// Verify next screen is valid before setting
		if nextScreen < app.Screens.MainMenu || nextScreen > app.Screens.DeconstructConfirm {
			logging.LogDebug("ERROR: Invalid next screen value: %d, defaulting to MainMenu", nextScreen)
			nextScreen = app.Screens.MainMenu
		}

		// Update the current screen
		logging.LogDebug("Setting next screen to: %d", nextScreen)
		app.SetCurrentScreen(nextScreen)
		logging.LogDebug("Screen set to: %d", app.GetCurrentScreen())
	}
}