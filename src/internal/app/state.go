// internal/app/state.go
package app

import (
	"thememanager/internal/logging"
)

// Screen represents the different UI screens
type Screen int

// Screen constants defining all possible application screens
const (
	// Main menu and high-level screens
	ScreenMainMenu Screen = iota
	ScreenThemesMenu        // New submenu for Themes
	ScreenOverlaysMenu      // New submenu for Overlays
	ScreenSyncCatalog
	ScreenSettingsMenu      // New Settings menu replacing Backup/Revert/Purge

	// Theme sub-screens
	ScreenInstalledThemes   // New screen for installed themes
	ScreenDownloadThemes    // Renamed (was ThemeGallery)
	ScreenThemeDownloadConfirm
	ScreenThemeDownloading
	ScreenThemeApplyConfirm
	ScreenThemeApplying

	// Overlay sub-screens
	ScreenInstalledOverlays // New screen for installed overlays
	ScreenDownloadOverlays  // Renamed (was OverlayGallery)
	ScreenOverlayDownloadConfirm
	ScreenOverlayDownloading
	ScreenOverlayApplyConfirm
	ScreenOverlayApplying

	// Settings sub-screens
	ScreenRestoreMenu       // Renamed (was RevertMenu)
	ScreenBackupThemeConfirm
	ScreenBackupThemeCreating
	ScreenBackupOverlayConfirm
	ScreenBackupOverlayCreating
	ScreenBackupAutoToggle

	// Restore sub-screens (renamed from Revert)
	ScreenRestoreThemeGallery   // Renamed
	ScreenRestoreThemeConfirm   // Renamed
	ScreenRestoreThemeApplying  // Renamed
	ScreenRestoreOverlayGallery // Renamed
	ScreenRestoreOverlayConfirm // Renamed
	ScreenRestoreOverlayApplying // Renamed

	// Purge screens
	ScreenPurgeConfirm
	ScreenPurging
)

// ScreenEnum provides named access to all screen constants
type ScreenEnum struct {
	// Main menu and high-level screens
	MainMenu     Screen
	ThemesMenu   Screen
	OverlaysMenu Screen
	SyncCatalog  Screen
	SettingsMenu Screen

	// Theme sub-screens
	InstalledThemes       Screen
	DownloadThemes        Screen
	ThemeDownloadConfirm  Screen
	ThemeDownloading      Screen
	ThemeApplyConfirm     Screen
	ThemeApplying         Screen

	// Overlay sub-screens
	InstalledOverlays     Screen
	DownloadOverlays      Screen
	OverlayDownloadConfirm Screen
	OverlayDownloading     Screen
	OverlayApplyConfirm    Screen
	OverlayApplying        Screen

	// Settings sub-screens
	RestoreMenu           Screen
	BackupThemeConfirm    Screen
	BackupThemeCreating   Screen
	BackupOverlayConfirm  Screen
	BackupOverlayCreating Screen
	BackupAutoToggle      Screen

	// Restore sub-screens (renamed from Revert)
	RestoreThemeGallery    Screen
	RestoreThemeConfirm    Screen
	RestoreThemeApplying   Screen
	RestoreOverlayGallery  Screen
	RestoreOverlayConfirm  Screen
	RestoreOverlayApplying Screen

	// Purge screens
	PurgeConfirm          Screen
	Purging               Screen
}

// AppState holds the current state of the application
type appState struct {
	CurrentScreen Screen
	SelectedItem  string
	AutoBackup    bool
}

// Define Screens as global enum for convenient access
var Screens = ScreenEnum{
	// Main menu and high-level screens
	MainMenu:     ScreenMainMenu,
	ThemesMenu:   ScreenThemesMenu,
	OverlaysMenu: ScreenOverlaysMenu,
	SyncCatalog:  ScreenSyncCatalog,
	SettingsMenu: ScreenSettingsMenu,

	// Theme sub-screens
	InstalledThemes:       ScreenInstalledThemes,
	DownloadThemes:        ScreenDownloadThemes,
	ThemeDownloadConfirm:  ScreenThemeDownloadConfirm,
	ThemeDownloading:      ScreenThemeDownloading,
	ThemeApplyConfirm:     ScreenThemeApplyConfirm,
	ThemeApplying:         ScreenThemeApplying,

	// Overlay sub-screens
	InstalledOverlays:     ScreenInstalledOverlays,
	DownloadOverlays:      ScreenDownloadOverlays,
	OverlayDownloadConfirm: ScreenOverlayDownloadConfirm,
	OverlayDownloading:     ScreenOverlayDownloading,
	OverlayApplyConfirm:    ScreenOverlayApplyConfirm,
	OverlayApplying:        ScreenOverlayApplying,

	// Settings sub-screens
	RestoreMenu:           ScreenRestoreMenu,
	BackupThemeConfirm:    ScreenBackupThemeConfirm,
	BackupThemeCreating:   ScreenBackupThemeCreating,
	BackupOverlayConfirm:  ScreenBackupOverlayConfirm,
	BackupOverlayCreating: ScreenBackupOverlayCreating,
	BackupAutoToggle:      ScreenBackupAutoToggle,

	// Restore sub-screens (renamed from Revert)
	RestoreThemeGallery:    ScreenRestoreThemeGallery,
	RestoreThemeConfirm:    ScreenRestoreThemeConfirm,
	RestoreThemeApplying:   ScreenRestoreThemeApplying,
	RestoreOverlayGallery:  ScreenRestoreOverlayGallery,
	RestoreOverlayConfirm:  ScreenRestoreOverlayConfirm,
	RestoreOverlayApplying: ScreenRestoreOverlayApplying,

	// Purge screens
	PurgeConfirm:           ScreenPurgeConfirm,
	Purging:                ScreenPurging,
}

// Global application state
var state appState

// GetCurrentScreen returns the current screen
func GetCurrentScreen() Screen {
	// Validate screen value
	if state.CurrentScreen < ScreenMainMenu || state.CurrentScreen > ScreenPurging {
		logging.LogDebug("WARNING: Invalid current screen value: %d, defaulting to MainMenu", state.CurrentScreen)
		state.CurrentScreen = ScreenMainMenu
	}
	return state.CurrentScreen
}

// SetCurrentScreen sets the current screen
func SetCurrentScreen(screen Screen) {
	// Validate screen value
	if screen < ScreenMainMenu || screen > ScreenPurging {
		logging.LogDebug("WARNING: Attempted to set invalid screen value: %d, using MainMenu instead", screen)
		screen = ScreenMainMenu
	}

	logging.LogDebug("Setting current screen from %d to %d", state.CurrentScreen, screen)
	state.CurrentScreen = screen
}

// GetSelectedItem returns the currently selected item
func GetSelectedItem() string {
	return state.SelectedItem
}

// SetSelectedItem sets the currently selected item
func SetSelectedItem(item string) {
	state.SelectedItem = item
}

// GetAutoBackup returns the auto-backup setting
func GetAutoBackup() bool {
	return state.AutoBackup
}

// SetAutoBackup sets the auto-backup setting
func SetAutoBackup(enabled bool) {
	state.AutoBackup = enabled
}