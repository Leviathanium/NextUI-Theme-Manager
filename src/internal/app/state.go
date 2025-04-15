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

// Screen constants with explicit values
const (
	// Main screens
	MainMenu Screen = 1

	// New screens for streamlined UI
	ThemesMenu         Screen = 5
	ComponentsMenu     Screen = 6
	ThemeApplyMenu     Screen = 7
	ThemeExtractMenu   Screen = 8
	ComponentApplyMenu Screen = 9

	// Theme management screens
	ThemeManagementMenu            Screen = 10
	ThemeImportTypeMenu            Screen = 11
	ThemeImportSelection           Screen = 12
	ThemeImportOptions             Screen = 13
	ThemeImportComponentSelection  Screen = 14
	ThemeImportConfirm             Screen = 15
	ThemeExportTypeMenu            Screen = 16
	ThemeExportName                Screen = 17
	ThemeExportOptions             Screen = 18
	ThemeExportComponentSelection  Screen = 19
	ThemeExportConfirm             Screen = 20
	ThemeConvertSelection          Screen = 21
	ThemeConvertOptions            Screen = 22
	ThemeConvertComponentSelection Screen = 23
	ThemeConvertConfirm            Screen = 24

	// Customization screens
	CustomizationMenu              Screen = 30
	GlobalOptionsMenu              Screen = 31
	SystemOptionsMenu              Screen = 32
	SystemOptionsForSelectedSystem Screen = 33
	WallpaperSelection             Screen = 34
	WallpaperConfirm               Screen = 35

	// Font screens
	FontSelection Screen = 40
	FontList      Screen = 41
	FontPreview   Screen = 42

	// Accent screens
	AccentMenu      Screen = 50
	AccentSelection Screen = 51
	AccentExport    Screen = 52

	// LED screens
	LEDMenu      Screen = 60
	LEDSelection Screen = 61
	LEDExport    Screen = 62

	// Icon screens
	IconsMenu           Screen = 70
	IconSelection       Screen = 71
	IconConfirm         Screen = 72
	ClearIconsConfirm   Screen = 73
	SystemIconSelection Screen = 74
	SystemIconConfirm   Screen = 75

	// Other screens
	ThemeSelection      Screen = 80
	DefaultThemeOptions Screen = 81
	ConfirmScreen       Screen = 82
	ResetMenu           Screen = 83
)

// ScreenEnum holds all available screens
type ScreenEnum struct {
	// Main screens
	MainMenu Screen

	// New screens for streamlined UI
	ThemesMenu         Screen
	ComponentsMenu     Screen
	ThemeApplyMenu     Screen
	ThemeExtractMenu   Screen
	ComponentApplyMenu Screen

	// Theme management screens
	ThemeManagementMenu            Screen
	ThemeImportTypeMenu            Screen
	ThemeImportSelection           Screen
	ThemeImportOptions             Screen
	ThemeImportComponentSelection  Screen
	ThemeImportConfirm             Screen
	ThemeExportTypeMenu            Screen
	ThemeExportName                Screen
	ThemeExportOptions             Screen
	ThemeExportComponentSelection  Screen
	ThemeExportConfirm             Screen
	ThemeConvertSelection          Screen
	ThemeConvertOptions            Screen
	ThemeConvertComponentSelection Screen
	ThemeConvertConfirm            Screen

	// Customization screens
	CustomizationMenu              Screen
	GlobalOptionsMenu              Screen
	SystemOptionsMenu              Screen
	SystemOptionsForSelectedSystem Screen
	WallpaperSelection             Screen
	WallpaperConfirm               Screen

	// Font screens
	FontSelection Screen
	FontList      Screen
	FontPreview   Screen

	// Accent screens
	AccentMenu      Screen
	AccentSelection Screen
	AccentExport    Screen

	// LED screens
	LEDMenu      Screen
	LEDSelection Screen
	LEDExport    Screen

	// Icon screens
	IconsMenu           Screen
	IconSelection       Screen
	IconConfirm         Screen
	ClearIconsConfirm   Screen
	SystemIconSelection Screen
	SystemIconConfirm   Screen

	// Other screens
	ThemeSelection      Screen
	DefaultThemeOptions Screen
	ConfirmScreen       Screen
	ResetMenu           Screen
}

// DefaultThemeAction represents the action to take for default themes
type DefaultThemeAction int

const (
	OverwriteAction DefaultThemeAction = iota
	DeleteAction
)

// AppState holds the current state of the application
type appState struct {
	CurrentScreen             Screen
	SelectedThemeType         ThemeType
	SelectedTheme             string
	DefaultAction             DefaultThemeAction
	SelectedFont              string
	SelectedFontSlot          string // Which font slot to modify (OG, Next, Legacy)
	SelectedAccentTheme       string
	SelectedLEDTheme          string
	SelectedAccentThemeSource ThemeSource
	SelectedLEDThemeSource    ThemeSource
	SelectedIconPack          string
	SelectedSystem            string // For system-specific options
}

// Global variables
var (
	Screens = ScreenEnum{
		// Main screens
		MainMenu: MainMenu,

		// New screens for streamlined UI
		ThemesMenu:         ThemesMenu,
		ComponentsMenu:     ComponentsMenu,
		ThemeApplyMenu:     ThemeApplyMenu,
		ThemeExtractMenu:   ThemeExtractMenu,
		ComponentApplyMenu: ComponentApplyMenu,

		// Theme management screens
		ThemeManagementMenu:            ThemeManagementMenu,
		ThemeImportTypeMenu:            ThemeImportTypeMenu,
		ThemeImportSelection:           ThemeImportSelection,
		ThemeImportOptions:             ThemeImportOptions,
		ThemeImportComponentSelection:  ThemeImportComponentSelection,
		ThemeImportConfirm:             ThemeImportConfirm,
		ThemeExportTypeMenu:            ThemeExportTypeMenu,
		ThemeExportName:                ThemeExportName,
		ThemeExportOptions:             ThemeExportOptions,
		ThemeExportComponentSelection:  ThemeExportComponentSelection,
		ThemeExportConfirm:             ThemeExportConfirm,
		ThemeConvertSelection:          ThemeConvertSelection,
		ThemeConvertOptions:            ThemeConvertOptions,
		ThemeConvertComponentSelection: ThemeConvertComponentSelection,
		ThemeConvertConfirm:            ThemeConvertConfirm,

		// Customization screens
		CustomizationMenu:              CustomizationMenu,
		GlobalOptionsMenu:              GlobalOptionsMenu,
		SystemOptionsMenu:              SystemOptionsMenu,
		SystemOptionsForSelectedSystem: SystemOptionsForSelectedSystem,
		WallpaperSelection:             WallpaperSelection,
		WallpaperConfirm:               WallpaperConfirm,

		// Font screens
		FontSelection: FontSelection,
		FontList:      FontList,
		FontPreview:   FontPreview,

		// Accent screens
		AccentMenu:      AccentMenu,
		AccentSelection: AccentSelection,
		AccentExport:    AccentExport,

		// LED screens
		LEDMenu:      LEDMenu,
		LEDSelection: LEDSelection,
		LEDExport:    LEDExport,

		// Icon screens
		IconsMenu:           IconsMenu,
		IconSelection:       IconSelection,
		IconConfirm:         IconConfirm,
		ClearIconsConfirm:   ClearIconsConfirm,
		SystemIconSelection: SystemIconSelection,
		SystemIconConfirm:   SystemIconConfirm,

		// Other screens
		ThemeSelection:      ThemeSelection,
		DefaultThemeOptions: DefaultThemeOptions,
		ConfirmScreen:       ConfirmScreen,
		ResetMenu:           ResetMenu,
	}

	state appState
)

// GetCurrentScreen returns the current screen
func GetCurrentScreen() Screen {
	// Ensure we never return an invalid screen value
	if state.CurrentScreen < MainMenu || state.CurrentScreen > ResetMenu {
		logging.LogDebug("WARNING: Invalid current screen value: %d, defaulting to MainMenu", state.CurrentScreen)
		state.CurrentScreen = MainMenu
	}
	return state.CurrentScreen
}

// SetCurrentScreen sets the current screen
func SetCurrentScreen(screen Screen) {
	// Validate screen value before setting
	if screen < MainMenu || screen > ResetMenu {
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
