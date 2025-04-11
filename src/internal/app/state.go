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
	AccentConfirm
	LEDConfirm
	CustomizationMenu
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
	AccentConfirm      Screen
	LEDConfirm         Screen
	CustomizationMenu  Screen
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
	SelectedAccentTheme string
	SelectedLEDTheme    string
}

// Global variables
var (
	Screens  = ScreenEnum{
		MainMenu:           MainMenu,
		ThemeSelection:     ThemeSelection,
		DefaultThemeOptions: DefaultThemeOptions,
		ConfirmScreen:      ConfirmScreen,
		FontSelection:      FontSelection,
		FontPreview:        FontPreview,
		AccentSelection:    AccentSelection,
		LEDSelection:       LEDSelection,
		AccentConfirm:      AccentConfirm,
		LEDConfirm:         LEDConfirm,
		CustomizationMenu:  CustomizationMenu,
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

// GetSelectedAccentTheme returns the selected accent theme
func GetSelectedAccentTheme() string {
	return state.SelectedAccentTheme
}

// SetSelectedAccentTheme sets the selected accent theme
func SetSelectedAccentTheme(theme string) {
	state.SelectedAccentTheme = theme
}

// GetSelectedLEDTheme returns the selected LED theme
func GetSelectedLEDTheme() string {
	return state.SelectedLEDTheme
}

// SetSelectedLEDTheme sets the selected LED theme
func SetSelectedLEDTheme(theme string) {
	state.SelectedLEDTheme = theme
}