// src/cmd/theme-manager/main.go
// Expanded main entry point for the NextUI Theme Manager application

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
		if currentScreen < app.Screens.MainMenu || currentScreen > app.Screens.ExportComponent {
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
			logging.LogDebug("Main menu returned next screen: %d", nextScreen)

		case app.Screens.ThemeImport:
			logging.LogDebug("Showing theme import selection")
			selection, exitCode = screens.ThemeImportScreen()
			nextScreen = screens.HandleThemeImport(selection, exitCode)
			logging.LogDebug("Theme import returned next screen: %d", nextScreen)

		case app.Screens.ThemeImportConfirm:
			logging.LogDebug("Showing theme import confirmation")
			selection, exitCode = screens.ThemeImportConfirmScreen()
			nextScreen = screens.HandleThemeImportConfirm(selection, exitCode)
			logging.LogDebug("Theme import confirmation returned next screen: %d", nextScreen)

		case app.Screens.ThemeExport:
			logging.LogDebug("Showing theme export screen")
			selection, exitCode = screens.ThemeExportScreen()
			nextScreen = screens.HandleThemeExport(selection, exitCode)
			logging.LogDebug("Theme export returned next screen: %d", nextScreen)

		case app.Screens.BrowseThemes:
			logging.LogDebug("Showing browse themes screen")
			selection, exitCode = screens.BrowseThemesScreen()
			nextScreen = screens.HandleBrowseThemes(selection, exitCode)
			logging.LogDebug("Browse themes returned next screen: %d", nextScreen)

		case app.Screens.DownloadThemes:
			logging.LogDebug("Showing download themes screen")
			selection, exitCode = screens.DownloadThemesScreen()
			nextScreen = screens.HandleDownloadThemes(selection, exitCode)
			logging.LogDebug("Download themes returned next screen: %d", nextScreen)

		case app.Screens.ComponentsMenu:
			logging.LogDebug("Showing components menu screen")
			selection, exitCode = screens.ComponentsMenuScreen()
			nextScreen = screens.HandleComponentsMenu(selection, exitCode)
			logging.LogDebug("Components menu returned next screen: %d", nextScreen)

		case app.Screens.ComponentOptions:
			logging.LogDebug("Showing component options screen")
			selection, exitCode = screens.ComponentOptionsScreen()
			nextScreen = screens.HandleComponentOptions(selection, exitCode)
			logging.LogDebug("Component options returned next screen: %d", nextScreen)

		case app.Screens.BrowseComponents:
			logging.LogDebug("Showing browse components screen")
			selection, exitCode = screens.BrowseComponentsScreen()
			nextScreen = screens.HandleBrowseComponents(selection, exitCode)
			logging.LogDebug("Browse components returned next screen: %d", nextScreen)

		case app.Screens.DownloadComponents:
			logging.LogDebug("Showing download components screen")
			selection, exitCode = screens.DownloadComponentsScreen()
			nextScreen = screens.HandleDownloadComponents(selection, exitCode)
			logging.LogDebug("Download components returned next screen: %d", nextScreen)

		case app.Screens.ExportComponent:
			logging.LogDebug("Showing export component screen")
			selection, exitCode = screens.ExportComponentScreen()
			nextScreen = screens.HandleExportComponent(selection, exitCode)
			logging.LogDebug("Export component returned next screen: %d", nextScreen)

		default:
			logging.LogDebug("Unknown screen type: %d, defaulting to MainMenu", currentScreen)
			nextScreen = app.Screens.MainMenu
		}

		// Add extra debug logging
		logging.LogDebug("Current screen: %d, Next screen: %d", currentScreen, nextScreen)

		// Verify next screen is valid before setting
		if nextScreen < app.Screens.MainMenu || nextScreen > app.Screens.ExportComponent {
			logging.LogDebug("ERROR: Invalid next screen value: %d, defaulting to MainMenu", nextScreen)
			nextScreen = app.Screens.MainMenu
		}

		// Update the current screen - add extra debugging
		logging.LogDebug("Setting next screen to: %d", nextScreen)
		app.SetCurrentScreen(nextScreen)
		logging.LogDebug("Screen set to: %d", app.GetCurrentScreen())
	}
}