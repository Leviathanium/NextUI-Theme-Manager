// internal/ui/wrappers.go
package ui

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"thememanager/internal/app"
)

// ShowMenu displays a menu using minui-list
func ShowMenu(menuItems string, title string, extraArgs ...string) (string, int) {
	app.LogDebug("Showing menu: %s", title)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		app.LogDebug("Error getting current directory: %v", err)
		return "", 1
	}

	// Create a temporary file for menu items
	tempFile, err := os.CreateTemp("", "menu-items-*")
	if err != nil {
		app.LogDebug("Error creating temp file: %v", err)
		return "", 1
	}
	inputPath := tempFile.Name()
	defer os.Remove(inputPath)

	// Write menu items to temp file
	_, err = tempFile.WriteString(menuItems)
	if err != nil {
		app.LogDebug("Error writing to temp file: %v", err)
		tempFile.Close()
		return "", 1
	}
	tempFile.Close()

	// Create temp file for output
	tempOutFile, err := os.CreateTemp("", "menu-output-*")
	if err != nil {
		app.LogDebug("Error creating output temp file: %v", err)
		return "", 1
	}
	outputPath := tempOutFile.Name()
	tempOutFile.Close()
	defer os.Remove(outputPath)

	// Build command arguments
	args := []string{
		"--format", "text",
		"--title", title,
		"--file", inputPath,
		"--write-location", outputPath,
	}

	// Add any extra arguments
	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	// Path to minui-list
	minuiListPath := filepath.Join(cwd, "minui-list")

	// Create command
	cmd := exec.Command(minuiListPath, args...)

	// Capture stderr
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	// Run the command
	err = cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			app.LogDebug("Error running minui-list: %v", err)
			return "", 1
		}
	}

	// Log stderr if any
	if stderrBuf.Len() > 0 {
		app.LogDebug("minui-list stderr: %s", stderrBuf.String())
	}

	// Read selection from output file if successful
	var selection string
	if exitCode == 0 {
		selectionBytes, err := os.ReadFile(outputPath)
		if err != nil {
			app.LogDebug("Error reading selection from output file: %v", err)
		} else {
			selection = string(selectionBytes)
			// Trim newlines
			selection = strings.TrimRight(selection, "\r\n")
			app.LogDebug("Selection: %s", selection)
		}
	}

	app.LogDebug("Menu exit code: %d", exitCode)
	return selection, exitCode
}

// ShowMessage displays a message using minui-presenter
func ShowMessage(message string, timeout string) (string, int) {
	app.LogDebug("Showing message: %s (timeout: %s)", message, timeout)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		app.LogDebug("Error getting current directory: %v", err)
		return "", 1
	}

	// Path to minui-presenter
	minuiPresenterPath := filepath.Join(cwd, "minui-presenter")

	// Create command
	args := []string{
		"--message", message,
		"--timeout", timeout,
		"--confirm-text", "OK",
		"--confirm-show",
	}

	cmd := exec.Command(minuiPresenterPath, args...)

	// Capture stderr
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	// Run the command
	err = cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			app.LogDebug("Error running minui-presenter: %v", err)
			return "", 1
		}
	}

	// Log stderr if any
	if stderrBuf.Len() > 0 {
		app.LogDebug("minui-presenter stderr: %s", stderrBuf.String())
	}

	app.LogDebug("Message exit code: %d", exitCode)
	return "", exitCode
}

// ShowConfirmDialog displays a yes/no confirmation dialog
func ShowConfirmDialog(message string) (string, int) {
	app.LogDebug("Showing confirm dialog: %s", message)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		app.LogDebug("Error getting current directory: %v", err)
		return "", 1
	}

	// Create temp file for options
	tempFile, err := os.CreateTemp("", "confirm-options-*")
	if err != nil {
		app.LogDebug("Error creating temp file: %v", err)
		return "", 1
	}
	inputPath := tempFile.Name()
	defer os.Remove(inputPath)

	// Write Yes/No options
	_, err = tempFile.WriteString("Yes\nNo")
	if err != nil {
		app.LogDebug("Error writing to temp file: %v", err)
		tempFile.Close()
		return "", 1
	}
	tempFile.Close()

	// Create temp file for output
	tempOutFile, err := os.CreateTemp("", "confirm-output-*")
	if err != nil {
		app.LogDebug("Error creating output temp file: %v", err)
		return "", 1
	}
	outputPath := tempOutFile.Name()
	tempOutFile.Close()
	defer os.Remove(outputPath)

	// Path to minui-list
	minuiListPath := filepath.Join(cwd, "minui-list")

	// Create command
	args := []string{
		"--format", "text",
		"--title", message,
		"--file", inputPath,
		"--write-location", outputPath,
		"--selected", "1", // Default to "No" for safety
	}

	cmd := exec.Command(minuiListPath, args...)

	// Capture stderr
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	// Run the command
	err = cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			app.LogDebug("Error running confirm dialog: %v", err)
			return "", 1
		}
	}

	// Log stderr if any
	if stderrBuf.Len() > 0 {
		app.LogDebug("confirm dialog stderr: %s", stderrBuf.String())
	}

	// Read selection from output file if successful
	var selection string
	if exitCode == 0 {
		selectionBytes, err := os.ReadFile(outputPath)
		if err != nil {
			app.LogDebug("Error reading selection from output file: %v", err)
		} else {
			selection = string(selectionBytes)
			// Trim newlines
			selection = strings.TrimRight(selection, "\r\n")
			app.LogDebug("Confirmation selection: %s", selection)
		}
	}

	app.LogDebug("Confirm dialog exit code: %d", exitCode)
	return selection, exitCode
}