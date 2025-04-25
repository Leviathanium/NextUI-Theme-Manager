// src/internal/app/state.go
// Expanded application state management for theme and component management

package app

import (
	"nextui-themes/internal/logging"
)

// Screen represents the different UI screens
type Screen int

// Expanded screen constants to include installed themes and components screens
const (
	MainMenu Screen = iota + 1
	InstalledThemes      // NEW: Browse local themes
	DownloadThemes       // RENAMED: Browse downloadable themes from catalog
	SyncCatalog          // RENAMED: Sync the themes catalog
	ThemeImport
	ThemeImportConfirm
	ThemeExport
	ComponentsMenu
	ComponentOptions
	InstalledComponents  // NEW: Browse locally installed components
	DownloadComponents   // RENAMED: Browse downloadable components from catalog
	SyncComponents       // UNCHANGED but renamed for clarity
	ExportComponent
	Deconstruction
	DeconstructConfirm
	OverlaySystemSelection  // New screen for system tag selection
)

// ScreenEnum holds all available screens
type ScreenEnum struct {
	MainMenu           Screen
	InstalledThemes    Screen // NEW
	DownloadThemes     Screen // RENAMED
	SyncCatalog        Screen // RENAMED
	ThemeImport        Screen
	ThemeImportConfirm Screen
	ThemeExport        Screen
	ComponentsMenu     Screen
	ComponentOptions   Screen
	InstalledComponents Screen // NEW
	DownloadComponents Screen // RENAMED
	SyncComponents     Screen
	ExportComponent    Screen
	Deconstruction     Screen
	DeconstructConfirm Screen
    OverlaySystemSelection Screen // New screen for system tag selection

}

// AppState holds the current state of the application
type appState struct {
	CurrentScreen           Screen
	SelectedTheme           string // For theme import/export
	SelectedComponentType   string // For component operations
	SelectedComponentOption string // For component operations
	SelectedSystemTag       string // New field for system tag selection
}

// Global variables
var (
	Screens = ScreenEnum{
		MainMenu:           MainMenu,
		InstalledThemes:    InstalledThemes,    // NEW
		DownloadThemes:     DownloadThemes,     // RENAMED
		SyncCatalog:        SyncCatalog,        // RENAMED
		ThemeImport:        ThemeImport,
		ThemeImportConfirm: ThemeImportConfirm,
		ThemeExport:        ThemeExport,
		ComponentsMenu:     ComponentsMenu,
		ComponentOptions:   ComponentOptions,
		InstalledComponents: InstalledComponents, // NEW
		DownloadComponents: DownloadComponents,   // RENAMED
		SyncComponents:     SyncComponents,
		ExportComponent:    ExportComponent,
		Deconstruction:     Deconstruction,
		DeconstructConfirm: DeconstructConfirm,
		OverlaySystemSelection: OverlaySystemSelection, // Add new screen
	}

	state appState
)

// Replace with:
func GetCurrentScreen() Screen {
    // Ensure we never return an invalid screen value
    if state.CurrentScreen < MainMenu || state.CurrentScreen > OverlaySystemSelection {
        logging.LogDebug("WARNING: Invalid current screen value: %d, defaulting to MainMenu", state.CurrentScreen)
        state.CurrentScreen = MainMenu
    }
    return state.CurrentScreen
}


// Replace with:
func SetCurrentScreen(screen Screen) {
    // Validate screen value before setting
    if screen < MainMenu || screen > OverlaySystemSelection {
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

// Add new getter/setter functions for SelectedSystemTag
// GetSelectedSystemTag returns the selected system tag
func GetSelectedSystemTag() string {
	return state.SelectedSystemTag
}

// SetSelectedSystemTag sets the selected system tag
func SetSelectedSystemTag(tag string) {
	state.SelectedSystemTag = tag
}