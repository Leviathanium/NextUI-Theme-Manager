// src/internal/ui/screens/apply.go
package screens

import (
	"fmt"
	"strings"
    "path/filepath"
    "os"
	"thememanager/internal/app"
	"thememanager/internal/themes"
	"thememanager/internal/ui"
)

// ShowApplyThemeScreen displays the theme selection screen for applying
func ShowApplyThemeScreen() (string, int) {
    app.LogDebug("Showing apply theme screen")

    // Get list of available themes
    themeNames, err := themes.ListThemes()
    if err != nil {
        app.LogDebug("Error listing themes: %v", err)
        ui.ShowMessage(fmt.Sprintf("Error listing themes: %s", err), "3")
        return "", 2  // Exit code 2 signals a controlled return to main menu
    }

    // Check if we have any themes
    if len(themeNames) == 0 {
        app.LogDebug("No themes available")
        // Show message and wait for user acknowledgment
        ui.ShowMessage("No themes available. Please download or import themes first.", "0")
        return "", 2  // Return to main menu with proper exit code
    }

    // Build theme data for gallery
    var themeDataList []map[string]string
    for _, themeName := range themeNames {
        themePath := themes.GetThemePath(themeName)

        // Create theme info
        themeInfo := map[string]string{
            "name": themeName,
            "is_valid": "false", // Default to invalid until proven otherwise
        }

        // Try to read manifest with non-strict validation to display
        manifest, err := themes.ReadManifest(themePath, false)
        if err == nil {
            themeInfo["author"] = manifest.Author
            themeInfo["description"] = manifest.Description

            // Perform validation test
            if themes.IsManifestValid(manifest) {
                themeInfo["is_valid"] = "true"
            } else {
                app.LogDebug("Theme %s has invalid manifest", themeName)
            }
        } else {
            app.LogDebug("Warning: Failed to read manifest for theme %s: %v", themeName, err)
        }

        // Check if preview exists
        previewPath := filepath.Join(themePath, themes.ThemePreviewFile)
        if _, err := os.Stat(previewPath); err == nil {
            themeInfo["preview"] = previewPath
        }

        themeDataList = append(themeDataList, themeInfo)
    }

    // Show theme gallery
    return ui.ShowThemeGallery(themeDataList, "Select Theme to Apply")
}

// HandleApplyThemeScreen processes the theme selection
func HandleApplyThemeScreen(selection string, exitCode int) app.Screen {
    app.LogDebug("HandleApplyThemeScreen called with selection: '%s', exitCode: %d", selection, exitCode)

    if exitCode == 0 && selection != "" {
        // User selected a theme
        app.SetSelectedItem(selection)
        return app.ScreenApplyThemeConfirm
    } else if exitCode == 1 || exitCode == 2 {
        // User cancelled or invalid theme was selected
        return app.ScreenMainMenu
    }

    return app.ScreenApplyTheme
}

// ShowApplyThemeConfirmScreen displays the confirmation screen for applying a theme
func ShowApplyThemeConfirmScreen() (string, int) {
    app.LogDebug("Showing apply theme confirmation screen")

    selectedTheme := app.GetSelectedItem()

    // Get theme path
    themePath := themes.GetThemePath(selectedTheme)

    // Try to read manifest for additional info - non-strict for display
    manifest, err := themes.ReadManifest(themePath, false)

    var confirmMessage string
    if err == nil {
        // Format with manifest details
        confirmMessage = fmt.Sprintf("Apply theme '%s' v%s by %s?",
            manifest.Name, manifest.Version, manifest.Author)
    } else {
        // Simple confirmation without manifest details
        confirmMessage = fmt.Sprintf("Apply theme '%s'?", selectedTheme)
    }

    return ui.ShowConfirmDialog(confirmMessage)
}

// HandleApplyThemeConfirmScreen processes the confirmation result
func HandleApplyThemeConfirmScreen(selection string, exitCode int) app.Screen {
    app.LogDebug("HandleApplyThemeConfirmScreen called with selection: '%s', exitCode: %d", selection, exitCode)

    if exitCode == 0 && selection == "Yes" {
        // User confirmed - proceed to applying
        return app.ScreenApplyingTheme
    } else {
        // User cancelled
        return app.ScreenApplyTheme
    }
}

// ShowApplyingThemeScreen displays the theme applying progress screen
func ShowApplyingThemeScreen() (string, int) {
    app.LogDebug("Showing applying theme screen")

    selectedTheme := app.GetSelectedItem()
    return ui.ShowMessage("Applying theme '" + selectedTheme + "'...", "2")
}

// HandleApplyingThemeScreen processes the applying operation
func HandleApplyingThemeScreen(selection string, exitCode int) app.Screen {
    app.LogDebug("HandleApplyingThemeScreen called with exitCode: %d", exitCode)

    selectedTheme := app.GetSelectedItem()

    // Extract theme name from selection (remove "by Author" part if present)
    themeName := selectedTheme
    if idx := strings.Index(themeName, " by "); idx > 0 {
        themeName = themeName[:idx]
    }

    // Apply the theme
    err := themes.ApplyTheme(themeName)

    if err != nil {
        app.LogDebug("Error applying theme: %v", err)
        ui.ShowMessage(fmt.Sprintf("Error applying theme: %s", err), "3")
    }

    return app.ScreenThemeApplied
}

// ShowThemeAppliedScreen displays the theme applied success screen
func ShowThemeAppliedScreen() (string, int) {
    app.LogDebug("Showing theme applied screen")

    selectedTheme := app.GetSelectedItem()
    return ui.ShowMessage("Theme '" + selectedTheme + "' applied successfully!", "2")
}

// HandleThemeAppliedScreen processes the success screen
func HandleThemeAppliedScreen(selection string, exitCode int) app.Screen {
    app.LogDebug("HandleThemeAppliedScreen called with exitCode: %d", exitCode)

    // Return to main menu after showing success message
    return app.ScreenMainMenu
}