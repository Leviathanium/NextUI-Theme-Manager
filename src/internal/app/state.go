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

// ThemeSource represents the source of themes (preset vs custom)
type ThemeSource int

const (
	PresetSource ThemeSource = iota
	CustomSource
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
	AccentMenu
	AccentSelection
	AccentExport
	LEDMenu
	LEDSelection
	LEDExport
	CustomizationMenu
	IconsMenu
	IconSelection
	IconConfirm
	ClearIconsConfirm
	GlobalOptionsMenu
	SystemOptionsMenu
	SystemOptionsForSelectedSystem
	SystemIconSelection
	SystemIconConfirm
	ResetMenu
	WallpaperSelection
	WallpaperConfirm
)

// ScreenEnum holds all available screens
type ScreenEnum struct {
	MainMenu           Screen
	ThemeSelection     Screen
	DefaultThemeOptions Screen
	ConfirmScreen      Screen
	FontSelection      Screen
	FontPreview        Screen
	AccentMenu         Screen
	AccentSelection    Screen
	AccentExport       Screen
	LEDMenu            Screen
	LEDSelection       Screen
	LEDExport          Screen
	CustomizationMenu  Screen
	IconsMenu          Screen
	IconSelection      Screen
	IconConfirm        Screen
	ClearIconsConfirm  Screen
	GlobalOptionsMenu  Screen
	SystemOptionsMenu  Screen
	SystemOptionsForSelectedSystem Screen
	SystemIconSelection Screen
	SystemIconConfirm   Screen
	ResetMenu          Screen
	WallpaperSelection Screen
	WallpaperConfirm   Screen
}

// DefaultThemeAction represents the action to take for default themes
type DefaultThemeAction int

const (
	OverwriteAction DefaultThemeAction = iota
	DeleteAction
)

// AppState holds the current state of the application
type appState struct {
	CurrentScreen           Screen
	SelectedThemeType       ThemeType
	SelectedTheme           string
	DefaultAction           DefaultThemeAction
	SelectedFont            string
	SelectedAccentTheme     string
	SelectedLEDTheme        string
	SelectedAccentThemeSource ThemeSource
	SelectedLEDThemeSource    ThemeSource
	SelectedIconPack        string
	SelectedSystem          string // For system-specific options
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
		AccentMenu:         AccentMenu,
		AccentSelection:    AccentSelection,
		AccentExport:       AccentExport,
		LEDMenu:            LEDMenu,
		LEDSelection:       LEDSelection,
		LEDExport:          LEDExport,
		CustomizationMenu:  CustomizationMenu,
		IconsMenu:          IconsMenu,
		IconSelection:      IconSelection,
		IconConfirm:        IconConfirm,
		ClearIconsConfirm:  ClearIconsConfirm,
		GlobalOptionsMenu:  GlobalOptionsMenu,
		SystemOptionsMenu:  SystemOptionsMenu,
		SystemOptionsForSelectedSystem: SystemOptionsForSelectedSystem,
		SystemIconSelection: SystemIconSelection,
		SystemIconConfirm:   SystemIconConfirm,
		ResetMenu:          ResetMenu,
		WallpaperSelection: WallpaperSelection,
		WallpaperConfirm:   WallpaperConfirm,
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

// GetSelectedAccentThemeSource returns the selected accent theme source
func GetSelectedAccentThemeSource() ThemeSource {
	return state.SelectedAccentThemeSource
}

// SetSelectedAccentThemeSource sets the selected accent theme source
func SetSelectedAccentThemeSource(source ThemeSource) {
	state.SelectedAccentThemeSource = source
}

// GetSelectedLEDThemeSource returns the selected LED theme source
func GetSelectedLEDThemeSource() ThemeSource {
	return state.SelectedLEDThemeSource
}

// SetSelectedLEDThemeSource sets the selected LED theme source
func SetSelectedLEDThemeSource(source ThemeSource) {
	state.SelectedLEDThemeSource = source
}

// GetSelectedIconPack returns the selected icon pack
func GetSelectedIconPack() string {
	return state.SelectedIconPack
}

// SetSelectedIconPack sets the selected icon pack
func SetSelectedIconPack(iconPack string) {
	state.SelectedIconPack = iconPack
}

// GetSelectedSystem returns the selected system
func GetSelectedSystem() string {
	return state.SelectedSystem
}

// SetSelectedSystem sets the selected system
func SetSelectedSystem(system string) {
	state.SelectedSystem = system
}