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
	CustomizationMenu // Added for customization menu
	// QuickSettings screen has been removed
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
	CustomizationMenu Screen
	// QuickSettings field removed
}

// DefaultThemeAction represents the action to take for default themes
type DefaultThemeAction int

const (
	OverwriteAction DefaultThemeAction = iota
	DeleteAction
)

// AppState holds the current state of the application
type appState struct {
	CurrentScreen        Screen
	SelectedThemeType    ThemeType
	SelectedTheme        string
	DefaultAction        DefaultThemeAction
	SelectedFont         string
	ColorR               int // For color selections
	ColorG               int
	ColorB               int
	LEDBrightness        int // For LED selections
	LEDSpeed             int
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
		CustomizationMenu: CustomizationMenu,
		// QuickSettings initialization removed
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

// GetColorSelections returns the color RGB values
func GetColorSelections() (int, int, int) {
	return state.ColorR, state.ColorG, state.ColorB
}

// SetColorSelections sets the color RGB values
func SetColorSelections(r int, g int, b int) {
	state.ColorR = r
	state.ColorG = g
	state.ColorB = b
}

// GetLEDSelections returns the LED brightness and speed
func GetLEDSelections() (int, int) {
	return state.LEDBrightness, state.LEDSpeed
}

// SetLEDSelections sets the LED brightness and speed
func SetLEDSelections(brightness int, speed int) {
	state.LEDBrightness = brightness
	state.LEDSpeed = speed
}