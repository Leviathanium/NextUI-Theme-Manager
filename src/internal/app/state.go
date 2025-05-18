// internal/app/state.go
package app

import (
	"fmt"
)

// Screen represents a screen in the application
type Screen int

// Screen constants
const (
	ScreenMainMenu Screen = iota
	ScreenThemes
	ScreenSettings
	ScreenAbout
	// Add more screens as needed
)

// Current screen state
var currentScreen Screen

// GetCurrentScreen returns the current screen
func GetCurrentScreen() Screen {
	return currentScreen
}

// SetCurrentScreen sets the current screen
func SetCurrentScreen(screen Screen) {
	LogDebug(fmt.Sprintf("Changing screen from %d to %d", currentScreen, screen))
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
	LogDebug(fmt.Sprintf("Selected item set to: %s", item))
	selectedItem = item
}