// cmd/theme-manager/main.go
package main

import (
	"fmt"
	"os"
	"runtime"

	"thememanager/internal/app"
	"thememanager/internal/ui"
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
		app.LogDebug("Current screen: %d", currentScreen)

		// Process current screen
		switch currentScreen {
		case app.ScreenMainMenu:
			selection, exitCode = showMainMenu()
			nextScreen = handleMainMenu(selection, exitCode)

		case app.ScreenSettings:
			selection, exitCode = showSettingsMenu()
			nextScreen = handleSettingsMenu(selection, exitCode)

		case app.ScreenAbout:
			selection, exitCode = showAboutScreen()
			nextScreen = handleAboutScreen(selection, exitCode)

		default:
			app.LogDebug("Unknown screen: %d, defaulting to main menu", currentScreen)
			nextScreen = app.ScreenMainMenu
		}

		// Update the current screen
		app.SetCurrentScreen(nextScreen)
	}
}

// showMainMenu displays the main menu
func showMainMenu() (string, int) {
	app.LogDebug("Showing main menu")

	menuItems := "Themes\nSettings\nAbout"

	return ui.ShowMenu(
		menuItems,
		"Theme Manager",
		"--cancel-text", "QUIT",
	)
}

// handleMainMenu processes the main menu selection
func handleMainMenu(selection string, exitCode int) app.Screen {
	app.LogDebug("Main menu selection: %s, exit code: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected an option
		switch selection {
		case "Themes":
			return app.ScreenThemes
		case "Settings":
			return app.ScreenSettings
		case "About":
			return app.ScreenAbout
		default:
			return app.ScreenMainMenu
		}
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/back/exit
		// Exit the application
		app.LogDebug("User exited the application")
		os.Exit(0)
	}

	return app.ScreenMainMenu
}

// showSettingsMenu displays the settings menu
func showSettingsMenu() (string, int) {
	app.LogDebug("Showing settings menu")

	menuItems := "Setting 1\nSetting 2\nSetting 3"

	return ui.ShowMenu(
		menuItems,
		"Settings",
		"--cancel-text", "BACK",
	)
}

// handleSettingsMenu processes the settings menu selection
func handleSettingsMenu(selection string, exitCode int) app.Screen {
	app.LogDebug("Settings menu selection: %s, exit code: %d", selection, exitCode)

	if exitCode == 0 {
		// User selected a setting
		ui.ShowMessage("Selected: " + selection, "2")
		return app.ScreenSettings
	} else if exitCode == 1 || exitCode == 2 {
		// User pressed cancel/back
		return app.ScreenMainMenu
	}

	return app.ScreenSettings
}

// showAboutScreen displays the about screen
func showAboutScreen() (string, int) {
	app.LogDebug("Showing about screen")

	ui.ShowMessage("Theme Manager v1.0\nCreated for MinUI", "3")

	// No selection, just return empty string and success code
	return "", 0
}

// handleAboutScreen processes the about screen
func handleAboutScreen(selection string, exitCode int) app.Screen {
	// After showing the about message, return to main menu
	return app.ScreenMainMenu
}