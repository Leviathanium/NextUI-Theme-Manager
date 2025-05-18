// internal/app/state.go
package app

import (
	"fmt"
)

// Screen represents a screen in the application
type Screen int

// Screen constants for all application screens
const (
	ScreenMainMenu Screen = iota

	// Apply theme screens
	ScreenApplyTheme
	ScreenApplyThemeConfirm
	ScreenApplyingTheme
	ScreenThemeApplied

	// Download theme screens
	ScreenDownloadTheme
	ScreenDownloadThemeConfirm
	ScreenDownloadingTheme
	ScreenThemeDownloaded
	ScreenCatalogNotSynced

	// Sync catalog screens
	ScreenSyncCatalog
	ScreenSyncingCatalog
	ScreenSyncComplete
	ScreenSyncFailed

	// Backups menu and screens
	ScreenBackupsMenu

	// Export theme screens
	ScreenExportTheme
	ScreenExportThemeConfirm
	ScreenExportingTheme
	ScreenThemeExported

	// Restore theme screens
	ScreenRestoreTheme
	ScreenRestoreThemeConfirm
	ScreenRestoringTheme
	ScreenThemeRestored

	// About screen
	ScreenAbout
)

// ScreenName returns the name of a screen for logging
func ScreenName(screen Screen) string {
	names := map[Screen]string{
		ScreenMainMenu:           "MainMenu",
		ScreenApplyTheme:         "ApplyTheme",
		ScreenApplyThemeConfirm:  "ApplyThemeConfirm",
		ScreenApplyingTheme:      "ApplyingTheme",
		ScreenThemeApplied:       "ThemeApplied",
		ScreenDownloadTheme:      "DownloadTheme",
		ScreenDownloadThemeConfirm: "DownloadThemeConfirm",
		ScreenDownloadingTheme:   "DownloadingTheme",
		ScreenThemeDownloaded:    "ThemeDownloaded",
		ScreenCatalogNotSynced:   "CatalogNotSynced",
		ScreenSyncCatalog:        "SyncCatalog",
		ScreenSyncingCatalog:     "SyncingCatalog",
		ScreenSyncComplete:       "SyncComplete",
		ScreenSyncFailed:         "SyncFailed",
		ScreenBackupsMenu:        "BackupsMenu",
		ScreenExportTheme:        "ExportTheme",
		ScreenExportThemeConfirm: "ExportThemeConfirm",
		ScreenExportingTheme:     "ExportingTheme",
		ScreenThemeExported:      "ThemeExported",
		ScreenRestoreTheme:       "RestoreTheme",
		ScreenRestoreThemeConfirm: "RestoreThemeConfirm",
		ScreenRestoringTheme:     "RestoringTheme",
		ScreenThemeRestored:      "ThemeRestored",
		ScreenAbout:              "About",
	}

	if name, ok := names[screen]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", screen)
}

// Current screen state
var currentScreen Screen

// GetCurrentScreen returns the current screen
func GetCurrentScreen() Screen {
	return currentScreen
}

// SetCurrentScreen sets the current screen
func SetCurrentScreen(screen Screen) {
	LogDebug("Changing screen from %s to %s", ScreenName(currentScreen), ScreenName(screen))
	currentScreen = screen
}

// SelectedItem holds the currently selected item
var selectedItem string

// GetSelectedItem returns the currently selected item
func GetSelectedItem() string {
	return selectedItem
}

// SetSelectedItem sets the currently selected item
func SetSelectedItem(item string) {
	LogDebug("Selected item set to: %s", item)
	selectedItem = item
}