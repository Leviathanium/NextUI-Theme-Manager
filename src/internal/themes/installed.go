// internal/themes/installed.go
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"thememanager/internal/logging"
	"thememanager/internal/system"
	"thememanager/internal/ui"
)

// ShowInstalledThemes displays a list of installed themes
func ShowInstalledThemes() (string, int) {
	logging.LogDebug("Showing installed themes")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read themes directory
	themesPath := filepath.Join(cwd, "Themes")
	entries, err := os.ReadDir(themesPath)
	if err != nil {
		logging.LogDebug("Error reading themes directory: %v", err)
		ui.ShowMessage("No installed themes found.", "3")
		return "", 1
	}

	// Build list of theme names
	var themesList strings.Builder
	themeCount := 0

	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), system.ThemeExtension) {
			themeName := entry.Name()
			themeName = strings.TrimSuffix(themeName, system.ThemeExtension)

			// Try to get author from manifest
			manifestPath := filepath.Join(themesPath, entry.Name(), "manifest.yml")
			manifest, err := ReadThemeManifest(manifestPath)

			if err == nil && manifest.Info.Author != "" {
				themesList.WriteString(fmt.Sprintf("%s by %s\n", themeName, manifest.Info.Author))
			} else {
				themesList.WriteString(fmt.Sprintf("%s\n", themeName))
			}

			themeCount++
		}
	}

	// Check if we found any themes
	if themeCount == 0 {
		ui.ShowMessage("No installed themes found.", "3")
		return "", 1
	}

	// Display the list
	return ui.DisplayMinUiList(
		themesList.String(),
		"text",
		"Installed Themes",
		"--cancel-text", "BACK",
	)
}

// ShowInstalledOverlays displays a list of installed overlays
func ShowInstalledOverlays() (string, int) {
	logging.LogDebug("Showing installed overlays")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read overlays directory
	overlaysPath := filepath.Join(cwd, "Overlays")
	entries, err := os.ReadDir(overlaysPath)
	if err != nil {
		logging.LogDebug("Error reading overlays directory: %v", err)
		ui.ShowMessage("No installed overlays found.", "3")
		return "", 1
	}

	// Build list of overlay names
	var overlaysList strings.Builder
	overlayCount := 0

	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), system.OverlayExtension) {
			overlayName := entry.Name()
			overlayName = strings.TrimSuffix(overlayName, system.OverlayExtension)

			// Try to get author from manifest
			manifestPath := filepath.Join(overlaysPath, entry.Name(), "manifest.yml")
			manifest, err := ReadOverlayManifest(manifestPath)

			if err == nil && manifest.Info.Author != "" {
				overlaysList.WriteString(fmt.Sprintf("%s by %s\n", overlayName, manifest.Info.Author))
			} else {
				overlaysList.WriteString(fmt.Sprintf("%s\n", overlayName))
			}

			overlayCount++
		}
	}

	// Check if we found any overlays
	if overlayCount == 0 {
		ui.ShowMessage("No installed overlays found.", "3")
		return "", 1
	}

	// Display the list
	return ui.DisplayMinUiList(
		overlaysList.String(),
		"text",
		"Installed Overlays",
		"--cancel-text", "BACK",
	)
}