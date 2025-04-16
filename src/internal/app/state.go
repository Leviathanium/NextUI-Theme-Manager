// src/internal/app/state.go
// Simplified application state management for theme import/export

package app

import (
	"nextui-themes/internal/logging"
)

// Screen represents the different UI screens
type Screen int

// Simplified screen constants - only keeping essential screens
const (
	MainMenu Screen = iota + 1
	ThemeImport
	ThemeImportConfirm
	ThemeExport
)

// ScreenEnum holds all available screens
type ScreenEnum struct {
	MainMenu          Screen
	ThemeImport       Screen
	ThemeImportConfirm Screen
	ThemeExport       Screen
}

// AppState holds the current state of the application
type appState struct {
	CurrentScreen Screen
	SelectedTheme string // For theme import/export
}

// Global variables
var (
	Screens = ScreenEnum{
		MainMenu:          MainMenu,
		ThemeImport:       ThemeImport,
		ThemeImportConfirm: ThemeImportConfirm,
		ThemeExport:       ThemeExport,
	}

	state appState
)

// GetCurrentScreen returns the current screen
func GetCurrentScreen() Screen {
	// Ensure we never return an invalid screen value
	if state.CurrentScreen < MainMenu || state.CurrentScreen > ThemeExport {
		logging.LogDebug("WARNING: Invalid current screen value: %d, defaulting to MainMenu", state.CurrentScreen)
		state.CurrentScreen = MainMenu
	}
	return state.CurrentScreen
}

// SetCurrentScreen sets the current screen
func SetCurrentScreen(screen Screen) {
	// Validate screen value before setting
	if screen < MainMenu || screen > ThemeExport {
		logging.LogDebug("WARNING: Attempted to set invalid screen value: %d, using MainMenu instead", screen)
		screen = MainMenu
	}

	// Add explicit debug logging
	logging.LogDebug("Setting current screen from %d to %d", state.CurrentScreen, screen)

	// Set the screen
	state.CurrentScreen = screen

	// Verify the screen was set correctly
	logging.LogDebug("Current screen is now: %d", state.CurrentScreen)
}

// GetSelectedTheme returns the selected theme
func GetSelectedTheme() string {
	return state.SelectedTheme
}

// SetSelectedTheme sets the selected theme
func SetSelectedTheme(theme string) {
	state.SelectedTheme = theme
}