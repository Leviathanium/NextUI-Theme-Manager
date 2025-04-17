// src/internal/app/state.go
// Expanded application state management for theme and component management

package app

import (
	"nextui-themes/internal/logging"
)

// Screen represents the different UI screens
type Screen int

// Expanded screen constants to include component-related screens
const (
	MainMenu Screen = iota + 1
	ThemeImport
	ThemeImportConfirm
	ThemeExport
	BrowseThemes
	DownloadThemes
	ComponentsMenu
	ComponentOptions
	BrowseComponents
	DownloadComponents
	ExportComponent
)

// ScreenEnum holds all available screens
type ScreenEnum struct {
	MainMenu          Screen
	ThemeImport       Screen
	ThemeImportConfirm Screen
	ThemeExport       Screen
	BrowseThemes      Screen
	DownloadThemes    Screen
	ComponentsMenu    Screen
	ComponentOptions  Screen
	BrowseComponents  Screen
	DownloadComponents Screen
	ExportComponent   Screen
}

// AppState holds the current state of the application
type appState struct {
	CurrentScreen       Screen
	SelectedTheme       string // For theme import/export
	SelectedComponentType string // For component operations
	SelectedComponentOption string // For component operations
}

// Global variables
var (
	Screens = ScreenEnum{
		MainMenu:           MainMenu,
		ThemeImport:        ThemeImport,
		ThemeImportConfirm: ThemeImportConfirm,
		ThemeExport:        ThemeExport,
		BrowseThemes:       BrowseThemes,
		DownloadThemes:     DownloadThemes,
		ComponentsMenu:     ComponentsMenu,
		ComponentOptions:   ComponentOptions,
		BrowseComponents:   BrowseComponents,
		DownloadComponents: DownloadComponents,
		ExportComponent:    ExportComponent,
	}

	state appState
)

// GetCurrentScreen returns the current screen
func GetCurrentScreen() Screen {
	// Ensure we never return an invalid screen value
	if state.CurrentScreen < MainMenu || state.CurrentScreen > ExportComponent {
		logging.LogDebug("WARNING: Invalid current screen value: %d, defaulting to MainMenu", state.CurrentScreen)
		state.CurrentScreen = MainMenu
	}
	return state.CurrentScreen
}

// SetCurrentScreen sets the current screen
func SetCurrentScreen(screen Screen) {
	// Validate screen value before setting
	if screen < MainMenu || screen > ExportComponent {
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

// GetSelectedComponentType returns the selected component type
func GetSelectedComponentType() string {
	return state.SelectedComponentType
}

// SetSelectedComponentType sets the selected component type
func SetSelectedComponentType(componentType string) {
	state.SelectedComponentType = componentType
}

// GetSelectedComponentOption returns the selected component option
func GetSelectedComponentOption() string {
	return state.SelectedComponentOption
}

// SetSelectedComponentOption sets the selected component option
func SetSelectedComponentOption(option string) {
	state.SelectedComponentOption = option
}