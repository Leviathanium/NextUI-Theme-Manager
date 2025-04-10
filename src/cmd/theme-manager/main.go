// src/cmd/theme-manager/main.go
// Main entry point for the NextUI Theme Manager application

package main

import (
	"os"
	"path/filepath"

	"nextui-themes/internal/app"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui/screens"
)
func main() {
	// Initialize the logger
	defer logging.CloseLogger()

	logging.LogDebug("Application started")

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

	logging.LogDebug("Starting main loop")

	// Main application loop
	for {
		var selection string
		var exitCode int
		var nextScreen app.Screen

		// Log current screen
		logging.LogDebug("Current screen: %d", app.GetCurrentScreen())

		// Process current screen
		switch app.GetCurrentScreen() {
		case app.Screens.MainMenu:
			logging.LogDebug("Showing main menu")
			selection, exitCode = screens.MainMenuScreen()
			nextScreen = screens.HandleMainMenu(selection, exitCode)

		case app.Screens.ThemeSelection:
			logging.LogDebug("Showing theme selection")
			selection, exitCode = screens.ThemeSelectionScreen()
			nextScreen = screens.HandleThemeSelection(selection, exitCode)

		case app.Screens.DefaultThemeOptions:
			logging.LogDebug("Showing default theme options")
			selection, exitCode = screens.DefaultThemeOptionsScreen()
			nextScreen = screens.HandleDefaultThemeOptions(selection, exitCode)

        case app.Screens.ConfirmScreen:
			logging.LogDebug("Showing confirmation screen")
			selection, exitCode = screens.ConfirmScreen()
			nextScreen = screens.HandleConfirmScreen(selection, exitCode)

		case app.Screens.FontSelection:
			logging.LogDebug("Showing font selection")
			selection, exitCode = screens.FontSelectionScreen()
			nextScreen = screens.HandleFontSelection(selection, exitCode)

		case app.Screens.FontPreview:
			logging.LogDebug("Showing font preview")
			selection, exitCode = screens.FontPreviewScreen()
			nextScreen = screens.HandleFontPreview(selection, exitCode)

		case app.Screens.AccentSelection:
			logging.LogDebug("Showing accent selection")
			selection, exitCode = screens.AccentSelectionScreen()
			nextScreen = screens.HandleAccentSelection(selection, exitCode)

		case app.Screens.LEDSelection:
			logging.LogDebug("Showing LED selection")
			selection, exitCode = screens.LEDSelectionScreen()
			nextScreen = screens.HandleLEDSelection(selection, exitCode)

        }

		// Update the current screen
		app.SetCurrentScreen(nextScreen)
	}
}