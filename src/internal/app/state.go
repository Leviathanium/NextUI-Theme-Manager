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
	ScreenThemes
	ScreenOverlays
	ScreenSyncCatalog
	ScreenBackup
	ScreenRevert
	ScreenPurge

	// Theme sub-screens
	ScreenThemeGallery
	ScreenThemeDownloadConfirm
	ScreenThemeDownloading
	ScreenThemeApplyConfirm
	ScreenThemeApplying

	// Overlay sub-screens
	ScreenOverlayGallery
	ScreenOverlayDownloadConfirm
	ScreenOverlayDownloading
	ScreenOverlayApplyConfirm
	ScreenOverlayApplying

	// Backup sub-screens
	ScreenBackupMenu
	ScreenBackupThemeConfirm
	ScreenBackupThemeCreating
	ScreenBackupOverlayConfirm
	ScreenBackupOverlayCreating
	ScreenBackupAutoToggle

	// Revert sub-screens
	ScreenRevertMenu
	ScreenRevertThemeGallery
	ScreenRevertThemeConfirm
	ScreenRevertThemeApplying
	ScreenRevertOverlayGallery
	ScreenRevertOverlayConfirm
	ScreenRevertOverlayApplying

	// Purge screens
	ScreenPurgeConfirm
	ScreenPurging
)

// ScreenEnum provides named access to all screen constants
type ScreenEnum struct {
	// Main menu and high-level screens
	MainMenu     Screen
	Themes       Screen
	Overlays     Screen
	SyncCatalog  Screen
	Backup       Screen
	Revert       Screen
	Purge        Screen

	// Theme sub-screens
	ThemeGallery          Screen
	ThemeDownloadConfirm  Screen
	ThemeDownloading      Screen
	ThemeApplyConfirm     Screen
	ThemeApplying         Screen

	// Overlay sub-screens
	OverlayGallery         Screen
	OverlayDownloadConfirm Screen
	OverlayDownloading     Screen
	OverlayApplyConfirm    Screen
	OverlayApplying        Screen

	// Backup sub-screens
	BackupMenu             Screen
	BackupThemeConfirm     Screen
	BackupThemeCreating    Screen
	BackupOverlayConfirm   Screen
	BackupOverlayCreating  Screen
	BackupAutoToggle       Screen

	// Revert sub-screens
	RevertMenu             Screen
	RevertThemeGallery     Screen
	RevertThemeConfirm     Screen
	RevertThemeApplying    Screen
	RevertOverlayGallery   Screen
	RevertOverlayConfirm   Screen
	RevertOverlayApplying  Screen

	// Purge screens
	PurgeConfirm           Screen
	Purging                Screen
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
	Themes:       ScreenThemes,
	Overlays:     ScreenOverlays,
	SyncCatalog:  ScreenSyncCatalog,
	Backup:       ScreenBackup,
	Revert:       ScreenRevert,
	Purge:        ScreenPurge,

	// Theme sub-screens
	ThemeGallery:          ScreenThemeGallery,
	ThemeDownloadConfirm:  ScreenThemeDownloadConfirm,
	ThemeDownloading:      ScreenThemeDownloading,
	ThemeApplyConfirm:     ScreenThemeApplyConfirm,
	ThemeApplying:         ScreenThemeApplying,

	// Overlay sub-screens
	OverlayGallery:         ScreenOverlayGallery,
	OverlayDownloadConfirm: ScreenOverlayDownloadConfirm,
	OverlayDownloading:     ScreenOverlayDownloading,
	OverlayApplyConfirm:    ScreenOverlayApplyConfirm,
	OverlayApplying:        ScreenOverlayApplying,

	// Backup sub-screens
	BackupMenu:             ScreenBackupMenu,
	BackupThemeConfirm:     ScreenBackupThemeConfirm,
	BackupThemeCreating:    ScreenBackupThemeCreating,
	BackupOverlayConfirm:   ScreenBackupOverlayConfirm,
	BackupOverlayCreating:  ScreenBackupOverlayCreating,
	BackupAutoToggle:       ScreenBackupAutoToggle,

	// Revert sub-screens
	RevertMenu:             ScreenRevertMenu,
	RevertThemeGallery:     ScreenRevertThemeGallery,
	RevertThemeConfirm:     ScreenRevertThemeConfirm,
	RevertThemeApplying:    ScreenRevertThemeApplying,
	RevertOverlayGallery:   ScreenRevertOverlayGallery,
	RevertOverlayConfirm:   ScreenRevertOverlayConfirm,
	RevertOverlayApplying:  ScreenRevertOverlayApplying,

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