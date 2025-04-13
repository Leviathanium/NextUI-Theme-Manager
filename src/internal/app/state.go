// src/internal/app/state.go
// Application state management

package app

import (

	"nextui-themes/internal/logging"

)

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

// The Screen constants should be defined like this:
const (
    MainMenu Screen = iota + 1
    ThemeSelection
    DefaultThemeOptions
    ConfirmScreen
    FontSelection
    FontList
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
    // New screens for theme management
    ThemesMenu          // This should be 27
    ThemeImport         // This should be 28
    ThemeImportConfirm  // This should be 29
    ThemeExport         // This should be 30
)

// ScreenEnum holds all available screens
type ScreenEnum struct {
	MainMenu           Screen
	ThemeSelection     Screen
	DefaultThemeOptions Screen
	ConfirmScreen      Screen
	FontSelection      Screen
	FontList           Screen
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
	// New screen fields for theme management
	ThemesMenu         Screen
	ThemeImport        Screen
	ThemeImportConfirm Screen
	ThemeExport        Screen
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
	SelectedFontSlot        string // Which font slot to modify (OG, Next, Legacy)
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
		FontList:           FontList,
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
        ThemesMenu:         ThemesMenu,
        ThemeImport:        ThemeImport,
        ThemeImportConfirm: ThemeImportConfirm,
        ThemeExport:        ThemeExport,
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

// GetSelectedFontSlot returns the selected font slot
func GetSelectedFontSlot() string {
	return state.SelectedFontSlot
}

// SetSelectedFontSlot sets the selected font slot
func SetSelectedFontSlot(slot string) {
	state.SelectedFontSlot = slot
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