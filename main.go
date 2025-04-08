// main.go for NextUI Theme Manager
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Constants for paths
const (
	THEME_GLOBAL_PATH = "themes/global"
	THEME_DEFAULT_PATH = "themes/default"
	DEFAULT_BG = "themes/default/bg.png"
	LOG_FILE_PATH = "theme-manager.log"
	DEBUG = true // Set to true for verbose logging
)

// AppState tracks the current state of the application
type AppState struct {
	CurrentScreen string
	ThemesList    []string
	Logger        *log.Logger
	SdCardPath    string
	DebugMode     bool
}

// Initialize application state
var appState AppState

// Screens for navigation
const (
	MAIN_MENU = "MAIN_MENU"
	GLOBAL_THEMES = "GLOBAL_THEMES"
	THEME_APPLIED = "THEME_APPLIED"
)

// Initialize the application
func init() {
	// Set the working directory to the PAK directory
	pakDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}

	// Setup logging with timestamp
	logFile, err := os.OpenFile(filepath.Join(pakDir, LOG_FILE_PATH), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}
	appState.Logger = log.New(logFile, "", log.LstdFlags)
	appState.Logger.Println("=== Theme Manager starting ===")
	appState.Logger.Printf("Pak directory: %s", pakDir)

	// Enable debug mode
	appState.DebugMode = DEBUG
	appState.Logger.Printf("Debug mode: %v", appState.DebugMode)

	// Determine the SD card path for NextUI installations
	// Typically this is "/mnt/SDCARD" on the TrimUI brick
	appState.SdCardPath = "/mnt/SDCARD"
	appState.Logger.Printf("Using SD card path: %s", appState.SdCardPath)

	// Verify file system operations by testing write access
	testFilePath := filepath.Join(pakDir, "test_write_access.tmp")
	testContent := []byte("Test file write access: " + time.Now().String())

	appState.Logger.Printf("Testing file system write access to: %s", testFilePath)
	if err := ioutil.WriteFile(testFilePath, testContent, 0644); err != nil {
		appState.Logger.Printf("WARNING: File write test failed: %v", err)
	} else {
		appState.Logger.Printf("File write test successful")
		os.Remove(testFilePath) // Clean up test file
	}

	// Check if minui-list exists
	minuiListPath := filepath.Join(pakDir, "minui-list")
	if _, err := os.Stat(minuiListPath); err != nil {
		appState.Logger.Printf("WARNING: minui-list not found at: %s", minuiListPath)
	} else {
		appState.Logger.Printf("minui-list found at: %s", minuiListPath)
	}

	// Check if minui-presenter exists
	minuiPresenterPath := filepath.Join(pakDir, "minui-presenter")
	if _, err := os.Stat(minuiPresenterPath); err != nil {
		appState.Logger.Printf("WARNING: minui-presenter not found at: %s", minuiPresenterPath)
	} else {
		appState.Logger.Printf("minui-presenter found at: %s", minuiPresenterPath)
	}

	// Check if default theme exists
	defaultBgPath := filepath.Join(pakDir, DEFAULT_BG)
	if _, err := os.Stat(defaultBgPath); err != nil {
		appState.Logger.Printf("WARNING: Default background not found at: %s", defaultBgPath)
	} else {
		appState.Logger.Printf("Default background found at: %s", defaultBgPath)
	}

	// Load available themes
	loadAvailableThemes()

	// Start at the main menu
	appState.CurrentScreen = MAIN_MENU
	appState.Logger.Printf("Initial screen set to: %s", appState.CurrentScreen)
}

// Load available themes from the global themes directory
func loadAvailableThemes() {
	appState.Logger.Println("Loading available themes...")

	// Get current directory (should be the pak directory)
	pakDir, err := os.Getwd()
	if err != nil {
		appState.Logger.Printf("Error getting working directory: %v", err)
		return
	}

	// Get list of directories in the global themes folder
	themesDir := filepath.Join(pakDir, THEME_GLOBAL_PATH)
	appState.Logger.Printf("Looking for themes in: %s", themesDir)

	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		appState.Logger.Printf("Themes directory does not exist: %s", themesDir)
		return
	}

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		appState.Logger.Printf("Error reading themes directory: %v", err)
		return
	}

	// Clear the current themes list
	appState.ThemesList = []string{}

	// Add each directory with a bg.png file to the themes list
	for _, entry := range entries {
		if entry.IsDir() {
			themeName := entry.Name()
			bgPath := filepath.Join(themesDir, themeName, "bg.png")
			appState.Logger.Printf("Checking for theme file: %s", bgPath)

			if _, err := os.Stat(bgPath); err == nil {
				appState.ThemesList = append(appState.ThemesList, themeName)
				appState.Logger.Printf("Found theme: %s", themeName)
			} else {
				appState.Logger.Printf("No bg.png found in theme directory: %s", themeName)
			}
		}
	}

	appState.Logger.Printf("Loaded %d themes", len(appState.ThemesList))
	if len(appState.ThemesList) > 0 {
		appState.Logger.Printf("Available themes: %s", strings.Join(appState.ThemesList, ", "))
	}
}

// Display a menu using minui-list
func displayMinUiList(listContent string, format string, title string, extraArgs ...string) (string, int) {
	appState.Logger.Printf("Displaying menu: %s", title)
	if appState.DebugMode {
		appState.Logger.Printf("Menu content: %s", listContent)
	}

	// Create a temporary file for the list content
	tempFile, err := ioutil.TempFile("", "minui-list-*")
	if err != nil {
		appState.Logger.Printf("ERROR: Failed to create temp file: %v", err)
		return "", 1
	}
	defer os.Remove(tempFile.Name())

	// Write the list content to the temp file
	if _, err := tempFile.WriteString(listContent); err != nil {
		appState.Logger.Printf("ERROR: Failed to write to temp file: %v", err)
		return "", 1
	}
	tempFile.Close()

	// Use the temp file as input for minui-list
	args := []string{"--format", format, "--title", title, "--file", tempFile.Name()}

	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	// Create a file to write the selection to
	tmpOutputFile, err := ioutil.TempFile("", "minui-output-*")
	if err != nil {
		appState.Logger.Printf("ERROR: Failed to create output temp file: %v", err)
		return "", 1
	}
	outputPath := tmpOutputFile.Name()
	tmpOutputFile.Close()
	defer os.Remove(outputPath)

	// Add the --write-location flag to write selection to our temp file
	args = append(args, "--write-location", outputPath)

	if appState.DebugMode {
		appState.Logger.Printf("minui-list args: %s", strings.Join(args, " "))
	}

	cmd := exec.Command("./minui-list", args...)

	var stderrbuf bytes.Buffer
	cmd.Stderr = &stderrbuf

	appState.Logger.Printf("Executing minui-list command...")
	err = cmd.Run()

	exitCode := 0
	if err != nil {
		exitCode = cmd.ProcessState.ExitCode()
		appState.Logger.Printf("Command exited with error code: %d, error: %v", exitCode, err)
	} else {
		exitCode = 0
		appState.Logger.Printf("Command completed successfully with code: %d", exitCode)
	}

	errValue := stderrbuf.String()
	if errValue != "" {
		appState.Logger.Printf("minui-list stderr: %s", errValue)
	}

	// Read the selection from the output file
	var outValue string
	if exitCode == 0 {
		selectionBytes, err := ioutil.ReadFile(outputPath)
		if err != nil {
			appState.Logger.Printf("ERROR: Failed to read selection from output file: %v", err)
		} else {
			outValue = strings.TrimSpace(string(selectionBytes))
			appState.Logger.Printf("Selection read from file: [%s]", outValue)
		}
	}

	appState.Logger.Printf("Menu selection result: [%s] with exit code: %d", outValue, exitCode)
	return outValue, exitCode
}

// Show a message using minui-presenter
func showMessage(message string, timeout string) {
	appState.Logger.Printf("Showing message: %s (timeout: %s)", message, timeout)

	// Get the current directory
	pakDir, _ := os.Getwd()

	// Use explicit path to minui-presenter
	minuiPresenterPath := filepath.Join(pakDir, "minui-presenter")

	args := []string{"--message", message, "--timeout", timeout}
	cmd := exec.Command(minuiPresenterPath, args...)

	var stderrbuf bytes.Buffer
	cmd.Stderr = &stderrbuf

	err := cmd.Run()
	if err != nil {
		appState.Logger.Printf("Error showing message: %v", err)
		if stderrbuf.Len() > 0 {
			appState.Logger.Printf("minui-presenter stderr: %s", stderrbuf.String())
		}
	}
}

// Show loading message and keep it displayed until canceled
func showLoadingScreen() (func(), error) {
	appState.Logger.Printf("Showing loading message")

	// Get the current directory
	pakDir, _ := os.Getwd()

	// Use explicit path to minui-presenter
	minuiPresenterPath := filepath.Join(pakDir, "minui-presenter")

	// Create a context with cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Start minui-presenter with indefinite timeout
	cmd := exec.CommandContext(ctx, minuiPresenterPath,
		"--message", "Applying theme...",
		"--timeout", "-1", // Indefinite timeout
	)

	// Start the message in background
	if err := cmd.Start(); err != nil {
		appState.Logger.Printf("ERROR: Failed to show loading message: %v", err)
		cancel()
		return nil, err
	}

	appState.Logger.Printf("Loading message displayed with PID: %d", cmd.Process.Pid)

	// Return a function that will stop the loading message
	cleanup := func() {
		appState.Logger.Printf("Stopping loading message")
		cancel()
		if err := cmd.Wait(); err != nil {
			appState.Logger.Printf("WARNING: Loading message exited with error: %v", err)
		}
		appState.Logger.Printf("Loading message stopped")
	}

	return cleanup, nil
}

// Copy a file from source to destination
func copyFile(src, dst string) error {
	appState.Logger.Printf("Copying file from %s to %s", src, dst)

	// Verify source file exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		errMsg := fmt.Sprintf("Source file does not exist: %s", src)
		appState.Logger.Printf("ERROR: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		errMsg := fmt.Sprintf("Error opening source file: %v", err)
		appState.Logger.Printf("ERROR: %s", errMsg)
		return fmt.Errorf(errMsg)
	}
	defer sourceFile.Close()

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	appState.Logger.Printf("Creating destination directory if needed: %s", dstDir)

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		errMsg := fmt.Sprintf("Error creating destination directory: %v", err)
		appState.Logger.Printf("ERROR: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	// Create destination file
	appState.Logger.Printf("Creating destination file: %s", dst)
	destFile, err := os.Create(dst)
	if err != nil {
		errMsg := fmt.Sprintf("Error creating destination file: %v", err)
		appState.Logger.Printf("ERROR: %s", errMsg)
		return fmt.Errorf(errMsg)
	}
	defer destFile.Close()

	// Copy the contents
	bytesWritten, err := io.Copy(destFile, sourceFile)
	if err != nil {
		errMsg := fmt.Sprintf("Error copying file content: %v", err)
		appState.Logger.Printf("ERROR: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	appState.Logger.Printf("Successfully copied %d bytes from %s to %s", bytesWritten, src, dst)
	return nil
}

// Apply a theme by copying bg.png to all required locations
func applyTheme(themeName string) error {
	appState.Logger.Printf("=== Applying theme: %s ===", themeName)

	// Start loading screen and get cleanup function
	stopLoading, err := showLoadingScreen()
	if err != nil {
		appState.Logger.Printf("WARNING: Could not show loading screen: %v", err)
		// Continue without loading screen if it fails
	}

	// Ensure we stop the loading screen when we're done
	if stopLoading != nil {
		defer stopLoading()
	}

	var srcPath string
	pakDir, _ := os.Getwd()

	if themeName == "Default Theme" {
		srcPath = filepath.Join(pakDir, DEFAULT_BG)
	} else {
		srcPath = filepath.Join(pakDir, THEME_GLOBAL_PATH, themeName, "bg.png")
	}

	appState.Logger.Printf("Source background path: %s", srcPath)

	// Verify the source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		errMsg := fmt.Sprintf("Theme file does not exist: %s", srcPath)
		appState.Logger.Printf("ERROR: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	// List of standard locations to copy the background to
	standardLocations := []string{
		filepath.Join(appState.SdCardPath, "bg.png"),                        // Root
		filepath.Join(appState.SdCardPath, "Tools", "tg5040", ".media", "bg.png"),   // Tools
		filepath.Join(appState.SdCardPath, "Recently Played", ".media", "bg.png"),   // Recently Played
	}

	appState.Logger.Printf("Will copy to %d standard locations", len(standardLocations))

	// Track success/failure counts
	successCount := 0
	failureCount := 0

	// Copy to standard locations
	for _, location := range standardLocations {
		appState.Logger.Printf("Copying to standard location: %s", location)

		if err := copyFile(srcPath, location); err != nil {
			appState.Logger.Printf("WARNING: Could not copy to %s: %v", location, err)
			failureCount++
		} else {
			successCount++
		}
	}

	// Find all ROM directories and create .media folders with bg.png
	romsDir := filepath.Join(appState.SdCardPath, "Roms")
	appState.Logger.Printf("Looking for ROM directories in: %s", romsDir)

	if _, err := os.Stat(romsDir); os.IsNotExist(err) {
		appState.Logger.Printf("WARNING: Roms directory does not exist: %s", romsDir)
	} else {
		if entries, err := os.ReadDir(romsDir); err == nil {
			appState.Logger.Printf("Found %d entries in Roms directory", len(entries))

			for _, entry := range entries {
				if entry.IsDir() {
					romDir := filepath.Join(romsDir, entry.Name())
					appState.Logger.Printf("Processing ROM directory: %s", romDir)

					// Create .media directory if it doesn't exist
					mediaDir := filepath.Join(romDir, ".media")
					appState.Logger.Printf("Creating media directory: %s", mediaDir)

					if err := os.MkdirAll(mediaDir, 0755); err != nil {
						appState.Logger.Printf("WARNING: Could not create media directory %s: %v", mediaDir, err)
						failureCount++
						continue
					}

					// Copy bg.png to the .media directory
					dstPath := filepath.Join(mediaDir, "bg.png")
					appState.Logger.Printf("Copying to ROM media location: %s", dstPath)

					if err := copyFile(srcPath, dstPath); err != nil {
						appState.Logger.Printf("WARNING: Could not copy to %s: %v", dstPath, err)
						failureCount++
					} else {
						successCount++
					}
				}
			}
		} else {
			appState.Logger.Printf("WARNING: Could not access Roms directory: %v", err)
		}
	}

	appState.Logger.Printf("Theme application complete: %d successful copies, %d failures", successCount, failureCount)

	if failureCount > 0 {
		return fmt.Errorf("%d failures occurred during theme application", failureCount)
	}

	return nil
}

// Display the main menu
func mainMenuScreen() (string, int) {
	appState.Logger.Printf("=== Entering main menu screen ===")

	menuItems := []string{
		"Global Themes",
		"Default Theme",
	}

	menuText := strings.Join(menuItems, "\n")
	appState.Logger.Printf("Showing main menu with %d options", len(menuItems))

	result, exitCode := displayMinUiList(menuText, "text", "NextUI Theme Manager", "--cancel-text", "QUIT")
	appState.Logger.Printf("Main menu selection: [%s] with exit code: %d", result, exitCode)

	return result, exitCode
}

// Display the global themes menu
func globalThemesScreen() (string, int) {
	appState.Logger.Printf("=== Entering global themes screen ===")

	if len(appState.ThemesList) == 0 {
		appState.Logger.Printf("No themes found in themes list")
		showMessage("No themes found!", "3")
		return "", 404
	}

	appState.Logger.Printf("Showing themes list with %d themes", len(appState.ThemesList))
	menuText := strings.Join(appState.ThemesList, "\n")

	result, exitCode := displayMinUiList(menuText, "text", "Available Themes", "--cancel-text", "BACK")
	appState.Logger.Printf("Theme selection: [%s] with exit code: %d", result, exitCode)

	return result, exitCode
}

// Main application loop
func main() {
	appState.Logger.Println("=== Theme Manager main loop started ===")

	for {
		var selection string
		var exitCode int

		appState.Logger.Printf("Current screen: %s", appState.CurrentScreen)

		switch appState.CurrentScreen {
		case MAIN_MENU:
			appState.Logger.Printf("Processing MAIN_MENU screen")
			selection, exitCode = mainMenuScreen()

			switch {
			case exitCode == 0:
				appState.Logger.Printf("User selected: [%s]", selection)

				// Handle menu selection
				switch selection {
				case "Global Themes":
					appState.Logger.Printf("Transitioning to GLOBAL_THEMES screen")
					appState.CurrentScreen = GLOBAL_THEMES

				case "Default Theme":
					appState.Logger.Printf("Applying Default theme")
					err := applyTheme(selection)
					if err != nil {
						errMsg := fmt.Sprintf("Error applying default theme: %v", err)
						appState.Logger.Printf("ERROR: %s", errMsg)
						showMessage("Error applying default theme: " + err.Error(), "3")
					} else {
						appState.Logger.Printf("Default theme applied successfully")
						showMessage("Default theme applied successfully!", "2")
					}

				default:
					appState.Logger.Printf("WARNING: Unknown selection in main menu: [%s]", selection)
				}

			case exitCode == 1, exitCode == 2:
				// User pressed cancel/back/quit
				appState.Logger.Println("User quit the application")
				os.Exit(0)

			default:
				appState.Logger.Printf("WARNING: Unexpected exit code from main menu: %d", exitCode)
			}

		case GLOBAL_THEMES:
			appState.Logger.Printf("Processing GLOBAL_THEMES screen")
			selection, exitCode = globalThemesScreen()

			switch {
			case exitCode == 0:
				// User selected a theme
				appState.Logger.Printf("User selected theme: [%s]", selection)

				err := applyTheme(selection)
				if err != nil {
					errMsg := fmt.Sprintf("Error applying theme: %v", err)
					appState.Logger.Printf("ERROR: %s", errMsg)
					showMessage("Error applying theme: " + err.Error(), "3")
				} else {
					appState.Logger.Printf("Theme [%s] applied successfully", selection)
					showMessage("Theme applied successfully!", "2")
				}

				appState.Logger.Printf("Returning to MAIN_MENU screen")
				appState.CurrentScreen = MAIN_MENU

			case exitCode == 1, exitCode == 2:
				// User pressed back
				appState.Logger.Printf("User pressed back from themes list")
				appState.CurrentScreen = MAIN_MENU

			case exitCode == 404:
				// No themes found
				appState.Logger.Printf("No themes found, returning to MAIN_MENU")
				appState.CurrentScreen = MAIN_MENU

			default:
				appState.Logger.Printf("WARNING: Unexpected exit code from themes menu: %d", exitCode)
				appState.CurrentScreen = MAIN_MENU
			}

		default:
			appState.Logger.Printf("WARNING: Unknown screen state: %s, resetting to MAIN_MENU", appState.CurrentScreen)
			appState.CurrentScreen = MAIN_MENU
		}
	}
}