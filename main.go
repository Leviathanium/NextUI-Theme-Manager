package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ThemeType represents the type of theme operation
type ThemeType int

const (
	GlobalTheme ThemeType = iota + 1 // Changed from StaticTheme
	DynamicTheme
	CustomTheme                      // Changed from SystemTheme
	DefaultTheme
)

// Screen represents the different UI screens
type Screen int

const (
	MainMenu Screen = iota + 1
	ThemeSelection
	ConfirmScreen
)

// AppState holds the current state of the application
type AppState struct {
	CurrentScreen     Screen
	SelectedThemeType ThemeType
	SelectedTheme     string
}

var appState AppState

func init() {
	// Initialize logger
	err := InitLogger()
	if err != nil {
		// Can't log yet, but we'll handle this in main
		return
	}

	// Log the initialization
	LogDebug("Logger initialized")

	// Create absolute paths based on current directory
	cwd, err := os.Getwd()
	if err != nil {
		LogDebug("Error getting current directory: %v", err)
		return
	}

	// Set up environment variables for the TrimUI brick
	LogDebug("Setting environment variables")

	_ = os.Setenv("DEVICE", "brick")
	_ = os.Setenv("PLATFORM", "tg5040")

	// FIX: Add current directory to PATH instead of replacing it
	existingPath := os.Getenv("PATH")
	newPath := cwd + ":" + existingPath
	_ = os.Setenv("PATH", newPath)
	LogDebug("Updated PATH: %s", newPath)

	_ = os.Setenv("LD_LIBRARY_PATH", "/mnt/SDCARD/.system/tg5040/lib:/usr/trimui/lib")

	// Initialize app state
	appState.CurrentScreen = MainMenu

	// Create theme directories if they don't exist - UPDATED TO CORRECT DIRECTORIES
	LogDebug("Creating theme directories")

	err = os.MkdirAll(filepath.Join(cwd, "Themes", "Global"), 0755)
	if err != nil {
		LogDebug("Error creating Global themes directory: %v", err)
	}

	err = os.MkdirAll(filepath.Join(cwd, "Themes", "Dynamic"), 0755)
	if err != nil {
		LogDebug("Error creating Dynamic themes directory: %v", err)
	}

	err = os.MkdirAll(filepath.Join(cwd, "Themes", "Default"), 0755)
	if err != nil {
		LogDebug("Error creating Default themes directory: %v", err)
	}

	LogDebug("Initialization complete")
}

// mainMenuScreen shows the main menu with theme options
func mainMenuScreen() (string, int) {
	// SIMPLIFIED MENU ITEMS without numbers
	menu := []string{
		"Global Themes",
		"Dynamic Themes",
		"Custom Themes",
		"Default Theme",
	}

	return displayMinUiList(strings.Join(menu, "\n"), "text", "NextUI Theme Selector", "--cancel-text", "QUIT")
}

// handleMainMenu processes the user's selection from the main menu
func handleMainMenu(selection string, exitCode int) {
	LogDebug("handleMainMenu called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected an option
		switch selection {
		case "Global Themes": // Updated from "1. Static Themes"
			LogDebug("Selected Global Themes")
			appState.SelectedThemeType = GlobalTheme
			appState.CurrentScreen = ThemeSelection
		case "Dynamic Themes": // Updated from "2. Dynamic Themes"
			LogDebug("Selected Dynamic Themes")
			appState.SelectedThemeType = DynamicTheme
			appState.CurrentScreen = ThemeSelection
		case "Custom Themes": // Updated from "3. System Themes"
			LogDebug("Selected Custom Themes")
			appState.SelectedThemeType = CustomTheme
			appState.CurrentScreen = ThemeSelection
		case "Default Theme": // Updated from "4. Default MinUI Theme"
			LogDebug("Selected Default Theme")
			appState.SelectedThemeType = DefaultTheme
			appState.CurrentScreen = ConfirmScreen
		default:
			LogDebug("Unknown selection: %s", selection)
		}
	case 1, 2:
		// User pressed cancel or back
		LogDebug("User cancelled/exited")
		os.Exit(0)
	}
}

// themeSelectionScreen displays available themes based on the selected theme type
func themeSelectionScreen() (string, int) {
	var title string
	var themes []string
	var err error

	// Get current directory for theme paths
	cwd, err := os.Getwd()
	if err != nil {
		LogDebug("Error getting current directory: %v", err)
		return "", 1
	}

	switch appState.SelectedThemeType {
	case GlobalTheme:
		title = "Select Global Theme"

		// Scan global themes directory
		globalDir := filepath.Join(cwd, "Themes", "Global")
		entries, err := os.ReadDir(globalDir)
		if err != nil {
			LogDebug("Error reading Global themes directory: %v", err)
			showMessage(fmt.Sprintf("Error loading global themes: %s", err), "3")
			themes = []string{"No themes found"}
		} else {
			// Find directories that contain a bg.png file
			for _, entry := range entries {
				if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
					bgPath := filepath.Join(globalDir, entry.Name(), "bg.png")
					if _, err := os.Stat(bgPath); err == nil {
						themes = append(themes, entry.Name())
					}
				}
			}

			if len(themes) == 0 {
				LogDebug("No global themes found")
				showMessage("No global themes found. Create one in Themes/Global/", "3")
				themes = []string{"No themes found"}
			}
		}

	case DynamicTheme:
		title = "Select Dynamic Theme"
		// List actual dynamic themes
		themes, err = ListDynamicThemes()
		if err != nil {
			LogDebug("Error loading dynamic themes: %v", err)
			showMessage(fmt.Sprintf("Error loading dynamic themes: %s", err), "3")
			themes = []string{"No themes found"}
		}

		// If no themes found, show a message
		if len(themes) == 0 {
			LogDebug("No dynamic themes found")
			showMessage("No dynamic themes found. Create one in Themes/Dynamic/", "3")
			themes = []string{"No themes found"}
		}
	case CustomTheme:
		title = "Select System"

		// Get system paths to find all installed systems
		systemPaths, err := GetSystemPaths()
		if err != nil {
			LogDebug("Error getting system paths: %v", err)
			showMessage(fmt.Sprintf("Error detecting systems: %s", err), "3")
			return "", 1
		}

		// Add standard menu items
		themes = append(themes, "Root")
		themes = append(themes, "Recently Played")
		themes = append(themes, "Tools")

		// Add all detected rom systems
		for _, system := range systemPaths.Systems {
			themes = append(themes, system.Name)
		}

		if len(themes) == 0 {
			LogDebug("No systems found")
			showMessage("No systems found!", "3")
			themes = []string{"No systems found"}
		}
	}

	LogDebug("Displaying theme selection with %d options", len(themes))
	return displayMinUiList(strings.Join(themes, "\n"), "text", title)
}

// handleThemeSelection processes the user's theme selection
func handleThemeSelection(selection string, exitCode int) {
	LogDebug("handleThemeSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a theme
		appState.SelectedTheme = selection
		appState.CurrentScreen = ConfirmScreen
	case 1, 2:
		// User pressed cancel or back
		appState.CurrentScreen = MainMenu
	}
}

// confirmScreen asks for confirmation before applying a theme
func confirmScreen() (string, int) {
	var message string

	switch appState.SelectedThemeType {
	case GlobalTheme:
		message = fmt.Sprintf("Apply global theme '%s' to all directories?", appState.SelectedTheme)
	case DynamicTheme:
		message = fmt.Sprintf("Apply dynamic theme '%s'?", appState.SelectedTheme)
	case CustomTheme:
		message = fmt.Sprintf("Select theme for '%s'?", appState.SelectedTheme)
	case DefaultTheme:
		message = "Apply default theme to all directories?"
	}

	options := []string{
		"Yes",
		"No",
	}

	return displayMinUiList(strings.Join(options, "\n"), "text", message)
}

// handleConfirmScreen processes the user's confirmation choice
func handleConfirmScreen(selection string, exitCode int) {
	LogDebug("handleConfirmScreen called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		if selection == "Yes" {
			// Apply the selected theme
			LogDebug("User confirmed, applying theme")
			applyTheme()
		} else {
			LogDebug("User selected No, returning to theme selection")
			appState.CurrentScreen = ThemeSelection
		}
	case 1, 2:
		// User pressed cancel or back
		LogDebug("User cancelled, returning to previous screen")
		if appState.SelectedThemeType == DefaultTheme {
			appState.CurrentScreen = MainMenu
		} else {
			appState.CurrentScreen = ThemeSelection
		}
	}
}

// applyTheme applies the selected theme
func applyTheme() {
	var message string
	var err error

	switch appState.SelectedThemeType {
	case GlobalTheme:
		// Apply global theme to all directories
		LogDebug("Applying global theme: %s", appState.SelectedTheme)
		err = ApplyStaticTheme(appState.SelectedTheme)
		if err != nil {
			LogDebug("Error applying global theme: %v", err)
			message = fmt.Sprintf("Error: %s", err)
		} else {
			message = fmt.Sprintf("Applied global theme: %s", appState.SelectedTheme)
		}

	case DynamicTheme:
		// Skip if "No themes found" is selected
		if appState.SelectedTheme == "No themes found" {
			LogDebug("No theme selected")
			message = "No theme selected"
			appState.CurrentScreen = MainMenu
			showMessage(message, "3")
			return
		}

		// Apply dynamic theme pack
		LogDebug("Applying dynamic theme: %s", appState.SelectedTheme)
		err = ApplyDynamicTheme(appState.SelectedTheme)
		if err != nil {
			LogDebug("Error applying dynamic theme: %v", err)
			message = fmt.Sprintf("Error: %s", err)
		} else {
			message = fmt.Sprintf("Applied dynamic theme: %s", appState.SelectedTheme)
		}

	case CustomTheme:
		// Apply custom theme to specific system
		LogDebug("Applying custom theme to: %s", appState.SelectedTheme)
		err = CustomThemeSelection(appState.SelectedTheme)
		if err != nil {
			LogDebug("Error applying custom theme: %v", err)
			message = fmt.Sprintf("Error: %s", err)
		} else {
			message = fmt.Sprintf("Applied theme to: %s", appState.SelectedTheme)
		}

	case DefaultTheme:
		// Apply default theme to all backgrounds
		LogDebug("Applying default theme")
		err = RemoveAllBackgrounds()
		if err != nil {
			LogDebug("Error applying default theme: %v", err)
			message = fmt.Sprintf("Error: %s", err)
		} else {
			message = "Applied default theme"
		}
	}

	showMessage(message, "3")
	appState.CurrentScreen = MainMenu
}

// Updated displayMinUiList function that matches the working implementation
func displayMinUiList(list string, format string, title string, extraArgs ...string) (string, int) {
    LogDebug("Displaying minui-list with title: %s", title)
    LogDebug("minui-list content: %s", list)

    // Get current directory
    cwd, err := os.Getwd()
    if err != nil {
        LogDebug("Error getting current directory: %v", err)
        return "", 1
    }

    // Create a temporary file for the list content
    tempFile, err := os.CreateTemp("", "minui-list-input-*")
    if err != nil {
        LogDebug("ERROR: Failed to create temp input file: %v", err)
        return "", 1
    }
    inputPath := tempFile.Name()
    defer os.Remove(inputPath)

    // Write the list content to the temp file
    if _, err := tempFile.WriteString(list); err != nil {
        LogDebug("ERROR: Failed to write to temp input file: %v", err)
        tempFile.Close()
        return "", 1
    }
    tempFile.Close()

    // Create a temporary file for the output
    tempOutFile, err := os.CreateTemp("", "minui-list-output-*")
    if err != nil {
        LogDebug("ERROR: Failed to create temp output file: %v", err)
        return "", 1
    }
    outputPath := tempOutFile.Name()
    tempOutFile.Close()
    defer os.Remove(outputPath)

    // Build the command arguments
    args := []string{"--format", format, "--title", title, "--file", inputPath, "--write-location", outputPath}

    if extraArgs != nil {
        args = append(args, extraArgs...)
    }

    LogDebug("minui-list args: %v", args)

    // Use explicit path to minui-list
    minuiListPath := filepath.Join(cwd, "minui-list")
    cmd := exec.Command(minuiListPath, args...)

    var stderrbuf bytes.Buffer
    cmd.Stderr = &stderrbuf

    // Run the command
    err = cmd.Run()
    exitCode := 0
    if err != nil {
        exitCode = cmd.ProcessState.ExitCode()
        LogDebug("minui-list error: %v", err)
    }

    errValue := stderrbuf.String()
    if errValue != "" {
        LogDebug("stderr: %s", errValue)
    }

    // Read the selection from the output file
    var outValue string
    if exitCode == 0 {
        selectionBytes, err := os.ReadFile(outputPath)
        if err != nil {
            LogDebug("ERROR: Failed to read selection from output file: %v", err)
        } else {
            outValue = strings.TrimSpace(string(selectionBytes))
            LogDebug("Selection read from file: '%s'", outValue)
        }
    }

    LogDebug("minui-list output: '%s', exit code: %d", outValue, exitCode)
    return outValue, exitCode
}

// Update showMessage with better logging
func showMessage(message string, timeout string) {
	LogDebug("Showing message: %s (timeout: %s)", message, timeout)

	args := []string{"--message", message, "--timeout", timeout}
	cmd := exec.Command("minui-presenter", args...)
	err := cmd.Run()

	if err != nil {
		LogDebug("minui-presenter error: %v", err)
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() != 124 {
			log.Fatalf("failed to run minui-presenter: %v", err)
		}
	}
}

func main() {
	defer CloseLogger()

	LogDebug("Application started")

	// Check if minui-list is available
	_, err := exec.LookPath("minui-list")
	if err != nil {
		LogDebug("minui-list not found in PATH: %v", err)
		fmt.Println("Error: minui-list not found in PATH")
		return
	}

	// Check if minui-presenter is available
	_, err = exec.LookPath("minui-presenter")
	if err != nil {
		LogDebug("minui-presenter not found in PATH: %v", err)
		fmt.Println("Error: minui-presenter not found in PATH")
		return
	}

	LogDebug("Starting main loop")

	// Main application loop
	for {
		var selection string
		var exitCode int

		// Log current screen
		LogDebug("Current screen: %d", appState.CurrentScreen)

		switch appState.CurrentScreen {
		case MainMenu:
			LogDebug("Showing main menu")
			selection, exitCode = mainMenuScreen()
			LogDebug("Main menu selection: '%s', exit code: %d", selection, exitCode)
			handleMainMenu(selection, exitCode)

		case ThemeSelection:
			LogDebug("Showing theme selection")
			selection, exitCode = themeSelectionScreen()
			LogDebug("Theme selection: '%s', exit code: %d", selection, exitCode)
			handleThemeSelection(selection, exitCode)

		case ConfirmScreen:
			LogDebug("Showing confirmation screen")
			selection, exitCode = confirmScreen()
			LogDebug("Confirmation: '%s', exit code: %d", selection, exitCode)
			handleConfirmScreen(selection, exitCode)
		}
	}
}