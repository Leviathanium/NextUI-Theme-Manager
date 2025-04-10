// src/internal/app/state.go
// Application state management

package app

import (
	"fmt"
	"strconv"
)

// ThemeType represents the type of theme operation
type ThemeType int

const (
	GlobalTheme ThemeType = iota + 1
	DynamicTheme
	CustomTheme
	DefaultTheme
)

// Screen represents the different UI screens
type Screen int

const (
	MainMenu Screen = iota + 1
	ThemeSelection
	DefaultThemeOptions
	ConfirmScreen
	FontSelection
	FontPreview
	AccentSelection
	LEDSelection
	QuickSettings // Added new screen type
)

// ScreenEnum holds all available screens
type ScreenEnum struct {
	MainMenu           Screen
	ThemeSelection     Screen
	DefaultThemeOptions Screen
	ConfirmScreen      Screen
	FontSelection      Screen
	FontPreview        Screen
	AccentSelection    Screen
	LEDSelection       Screen
	QuickSettings      Screen // Added new screen
}

// DefaultThemeAction represents the action to take for default themes
type DefaultThemeAction int

const (
	OverwriteAction DefaultThemeAction = iota
	DeleteAction
)

// AppState holds the current state of the application
type appState struct {
	CurrentScreen      Screen
	SelectedThemeType  ThemeType
	SelectedTheme      string
	DefaultAction      DefaultThemeAction
	SelectedFont       string
	ColorSelections    map[string]string // Added for accent color selections
	LEDSelections      map[string]string // Added for LED settings selections
}

// Global variables
var (
	Screens  = ScreenEnum{
		MainMenu:          MainMenu,
		ThemeSelection:    ThemeSelection,
		DefaultThemeOptions: DefaultThemeOptions,
		ConfirmScreen:     ConfirmScreen,
		FontSelection:     FontSelection,
		FontPreview:       FontPreview,
		AccentSelection:   AccentSelection,
		LEDSelection:      LEDSelection,
		QuickSettings:     QuickSettings, // Added new screen
	}

	state appState
)

// Initialize appState
func init() {
	// Initialize maps
	state.ColorSelections = make(map[string]string)
	state.LEDSelections = make(map[string]string)
}

// GetCurrentScreen returns the current screen
func GetCurrentScreen() Screen {
	return state.CurrentScreen
}

// SetCurrentScreen sets the current screen
func SetCurrentScreen(screen Screen) {
	state.CurrentScreen = screen
}

// GetSelectedThemeType returns the selected theme type
func GetSelectedThemeType() ThemeType {
	return state.SelectedThemeType
}

// SetSelectedThemeType sets the selected theme type
func SetSelectedThemeType(themeType ThemeType) {
	state.SelectedThemeType = themeType
}

// GetSelectedTheme returns the selected theme
func GetSelectedTheme() string {
	return state.SelectedTheme
}

// SetSelectedTheme sets the selected theme
func SetSelectedTheme(theme string) {
	state.SelectedTheme = theme
}

// GetDefaultAction returns the default theme action
func GetDefaultAction() DefaultThemeAction {
	return state.DefaultAction
}

// SetDefaultAction sets the default theme action
func SetDefaultAction(action DefaultThemeAction) {
	state.DefaultAction = action
}

// GetSelectedFont returns the selected font
func GetSelectedFont() string {
	return state.SelectedFont
}

// SetSelectedFont sets the selected font
func SetSelectedFont(font string) {
	state.SelectedFont = font
}

// SetColorSelections sets the color selections
func SetColorSelections(color1, color2, color3 int) {
	// Store color selections as indices
	state.ColorSelections = map[string]string{
		"color1": fmt.Sprintf("%d", color1),
		"color2": fmt.Sprintf("%d", color2),
		"color3": fmt.Sprintf("%d", color3),
	}
}

// GetColorSelections returns the color selections
func GetColorSelections() (int, int, int) {
	// Get color selections as indices with defaults if not set
	color1, _ := strconv.Atoi(state.ColorSelections["color1"])
	color2, _ := strconv.Atoi(state.ColorSelections["color2"])
	color3, _ := strconv.Atoi(state.ColorSelections["color3"])
	return color1, color2, color3
}

// SetLEDSelections sets the LED selections
func SetLEDSelections(effect, color int) {
	// Store LED selections
	state.LEDSelections = map[string]string{
		"effect": fmt.Sprintf("%d", effect),
		"color":  fmt.Sprintf("%d", color),
	}
}

// GetLEDSelections returns the LED selections
func GetLEDSelections() (int, int) {
	// Get LED selections with defaults if not set
	effect, _ := strconv.Atoi(state.LEDSelections["effect"])
	color, _ := strconv.Atoi(state.LEDSelections["color"])
	return effect, color
}