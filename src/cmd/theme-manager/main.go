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
        currentScreen := app.GetCurrentScreen()
        logging.LogDebug("Current screen: %d", currentScreen)

        // Ensure screen value is valid
        if currentScreen < app.Screens.MainMenu || currentScreen > app.Screens.WallpaperConfirm {
            logging.LogDebug("CRITICAL ERROR: Invalid screen value: %d, resetting to MainMenu", currentScreen)
            app.SetCurrentScreen(app.Screens.MainMenu)
            continue
        }

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

		case app.Screens.ResetMenu:
			logging.LogDebug("Showing reset menu")
			selection, exitCode = screens.ResetMenuScreen()
			nextScreen = screens.HandleResetMenu(selection, exitCode)

		case app.Screens.ConfirmScreen:
			logging.LogDebug("Showing confirmation screen")
			selection, exitCode = screens.ConfirmScreen()
			nextScreen = screens.HandleConfirmScreen(selection, exitCode)

		case app.Screens.FontSelection:
			logging.LogDebug("Showing font slot selection")
			selection, exitCode = screens.FontSelectionScreen()
			nextScreen = screens.HandleFontSelection(selection, exitCode)

		case app.Screens.FontList:
			logging.LogDebug("Showing font list")
			selection, exitCode = screens.FontListScreen()
			nextScreen = screens.HandleFontList(selection, exitCode)

		case app.Screens.FontPreview:
			logging.LogDebug("Showing font preview")
			selection, exitCode = screens.FontPreviewScreen()
			nextScreen = screens.HandleFontPreview(selection, exitCode)

		case app.Screens.CustomizationMenu:
			logging.LogDebug("Showing customization menu")
			selection, exitCode = screens.CustomizationMenuScreen()
			nextScreen = screens.HandleCustomizationMenu(selection, exitCode)

		case app.Screens.GlobalOptionsMenu:
			logging.LogDebug("Showing global options menu")
			selection, exitCode = screens.GlobalOptionsMenuScreen()
			nextScreen = screens.HandleGlobalOptionsMenu(selection, exitCode)

		case app.Screens.SystemOptionsMenu:
			logging.LogDebug("Showing system options menu")
			selection, exitCode = screens.SystemOptionsMenuScreen()
			nextScreen = screens.HandleSystemOptionsMenu(selection, exitCode)

		case app.Screens.SystemOptionsForSelectedSystem:
			logging.LogDebug("Showing options for selected system")
			selection, exitCode = screens.SystemOptionsForSelectedSystemScreen()
			nextScreen = screens.HandleSystemOptionsForSelectedSystem(selection, exitCode)

		case app.Screens.WallpaperSelection:
			logging.LogDebug("Showing wallpaper selection")
			selection, exitCode = screens.WallpaperSelectionScreen()
			nextScreen = screens.HandleWallpaperSelection(selection, exitCode)

		case app.Screens.AccentMenu:
			logging.LogDebug("Showing accent menu")
			selection, exitCode = screens.AccentMenuScreen()
			nextScreen = screens.HandleAccentMenu(selection, exitCode)

		case app.Screens.AccentSelection:
			logging.LogDebug("Showing accent selection")
			selection, exitCode = screens.AccentSelectionScreen()
			nextScreen = screens.HandleAccentSelection(selection, exitCode)

		case app.Screens.AccentExport:
			logging.LogDebug("Handling accent export")
			selection, exitCode = screens.AccentExportScreen()
			nextScreen = screens.HandleAccentExport(selection, exitCode)

		case app.Screens.LEDMenu:
			logging.LogDebug("Showing LED menu")
			selection, exitCode = screens.LEDMenuScreen()
			nextScreen = screens.HandleLEDMenu(selection, exitCode)

		case app.Screens.LEDSelection:
			logging.LogDebug("Showing LED selection")
			selection, exitCode = screens.LEDSelectionScreen()
			nextScreen = screens.HandleLEDSelection(selection, exitCode)

		case app.Screens.LEDExport:
			logging.LogDebug("Handling LED export")
			selection, exitCode = screens.LEDExportScreen()
			nextScreen = screens.HandleLEDExport(selection, exitCode)

		case app.Screens.IconsMenu:
			logging.LogDebug("Showing icons menu")
			selection, exitCode = screens.IconsMenuScreen()
			nextScreen = screens.HandleIconsMenu(selection, exitCode)

		case app.Screens.IconSelection:
			logging.LogDebug("Showing icon selection")
			selection, exitCode = screens.IconSelectionScreen()
			nextScreen = screens.HandleIconSelection(selection, exitCode)

        case app.Screens.SystemIconSelection:
            logging.LogDebug("Showing system icon selection")
            selection, exitCode = screens.SystemIconSelectionScreen()
            nextScreen = screens.HandleSystemIconSelection(selection, exitCode)
            logging.LogDebug("System icon selection returned next screen: %d", nextScreen)

        case app.Screens.SystemIconConfirm:
            logging.LogDebug("Showing system icon confirmation")
            selection, exitCode = screens.SystemIconConfirmScreen()
            nextScreen = screens.HandleSystemIconConfirm(selection, exitCode)
            logging.LogDebug("System icon confirmation returned next screen: %d", nextScreen)

		case app.Screens.IconConfirm:
			logging.LogDebug("Showing icon confirmation")
			selection, exitCode = screens.IconConfirmScreen()
			nextScreen = screens.HandleIconConfirm(selection, exitCode)

		case app.Screens.ClearIconsConfirm:
			logging.LogDebug("Showing clear icons confirmation")
			selection, exitCode = screens.ClearIconsConfirmScreen()
			nextScreen = screens.HandleClearIconsConfirm(selection, exitCode)


        default:
            logging.LogDebug("Unhandled screen type: %d, defaulting to MainMenu", currentScreen)
            nextScreen = app.Screens.MainMenu
        }

        // Verify next screen is valid before setting
        if nextScreen < app.Screens.MainMenu || nextScreen > app.Screens.WallpaperConfirm {
            logging.LogDebug("ERROR: Invalid next screen value: %d, defaulting to MainMenu", nextScreen)
            nextScreen = app.Screens.MainMenu
        }

        // Update the current screen
        logging.LogDebug("Setting next screen to: %d", nextScreen)
        app.SetCurrentScreen(nextScreen)

	}
}