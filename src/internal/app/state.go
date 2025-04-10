// src/internal/app/state.go
// Application state management

package app

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
)

// ScreenEnum holds all available screens
type ScreenEnum struct {
	MainMenu          Screen
	ThemeSelection    Screen
	DefaultThemeOptions Screen
	ConfirmScreen     Screen
	FontSelection     Screen
	FontPreview       Screen
	AccentSelection   Screen
	LEDSelection      Screen
}

// DefaultThemeAction represents the action to take for default themes
type DefaultThemeAction int

const (
	OverwriteAction DefaultThemeAction = iota
	DeleteAction
)

// ColorSelections holds the currently selected accent color indices
type colorSelections struct {
	PrimaryIndex   int
	SecondaryIndex int
	TextIndex      int
}

// LEDSelections holds the currently selected LED settings
type ledSelections struct {
	ColorIndex  int
	EffectIndex int
}


// Update the appState struct definition to add the LED selections
type appState struct {
	CurrentScreen      Screen
	SelectedThemeType  ThemeType
	SelectedTheme      string
	DefaultAction      DefaultThemeAction
	SelectedFont       string
	ColorSelections    colorSelections
	LEDSelections      ledSelections // Add this new field
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
	}

	state appState
)

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

// SetColorSelections stores the selected color indices
func SetColorSelections(primary, secondary, text int) {
	state.ColorSelections.PrimaryIndex = primary
	state.ColorSelections.SecondaryIndex = secondary
	state.ColorSelections.TextIndex = text
}

// GetColorSelections retrieves the selected color indices
func GetColorSelections() (int, int, int) {
	return state.ColorSelections.PrimaryIndex,
	       state.ColorSelections.SecondaryIndex,
	       state.ColorSelections.TextIndex
}

// SetLEDSelections stores the selected LED color and effect indices
func SetLEDSelections(color, effect int) {
	state.LEDSelections.ColorIndex = color
	state.LEDSelections.EffectIndex = effect
}

// GetLEDSelections retrieves the selected LED color and effect indices
func GetLEDSelections() (int, int) {
	return state.LEDSelections.ColorIndex, state.LEDSelections.EffectIndex
}