// src/internal/ui/common.go
// Common UI utilities for the application

package ui

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
    "time"
	"nextui-themes/internal/logging"
)

// Selection represents the result of a UI interaction
type Selection struct {
	Value    string
	Code     int
	Error    error
}

// ShowMessageWithOperation displays a message while performing an operation,
// then cleans up and returns any error from the operation
func ShowMessageWithOperation(message string, operation func() error) error {
	logging.LogDebug("Showing message with operation: %s", message)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Use explicit path to minui-presenter
	minuiPresenterPath := filepath.Join(cwd, "minui-presenter")

	// Start presenter with negative timeout (will stay until killed)
	cmd := exec.Command(minuiPresenterPath, "--message", message, "--timeout", "-1")

	// Start in background
	if err := cmd.Start(); err != nil {
		logging.LogDebug("Error starting minui-presenter: %v", err)
		return err
	}

	// Ensure the process gets killed when we're done
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
			logging.LogDebug("Killed minui-presenter process")
		}
	}()

	// Run the operation
	operationErr := operation()

	// Small delay to make sure the message is visible for at least a moment
	// even if the operation is very fast
	time.Sleep(500 * time.Millisecond)

	return operationErr
}

// DisplayMinUiList displays a list of items using minui-list
func DisplayMinUiList(list string, format string, title string, extraArgs ...string) (string, int) {
	logging.LogDebug("Displaying minui-list with title: %s", title)
	logging.LogDebug("minui-list content: %s", list)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return "", 1
	}

	// Create a temporary file for the list content
	tempFile, err := os.CreateTemp("", "minui-list-input-*")
	if err != nil {
		logging.LogDebug("ERROR: Failed to create temp input file: %v", err)
		return "", 1
	}
	inputPath := tempFile.Name()
	defer os.Remove(inputPath)

	// Write the list content to the temp file
	if _, err := tempFile.WriteString(list); err != nil {
		logging.LogDebug("ERROR: Failed to write to temp input file: %v", err)
		tempFile.Close()
		return "", 1
	}
	tempFile.Close()

	// Create a temporary file for the output
	tempOutFile, err := os.CreateTemp("", "minui-list-output-*")
	if err != nil {
		logging.LogDebug("ERROR: Failed to create temp output file: %v", err)
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

	logging.LogDebug("minui-list args: %v", args)

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
		logging.LogDebug("minui-list error: %v", err)
	}

	errValue := stderrbuf.String()
	if errValue != "" {
		logging.LogDebug("stderr: %s", errValue)
	}

	// Read the selection from the output file
	var outValue string
	if exitCode == 0 {
		selectionBytes, err := os.ReadFile(outputPath)
		if err != nil {
			logging.LogDebug("ERROR: Failed to read selection from output file: %v", err)
		} else {
			outValue = strings.TrimSpace(string(selectionBytes))
			logging.LogDebug("Selection read from file: '%s'", outValue)
		}
	}

	logging.LogDebug("minui-list output: '%s', exit code: %d", outValue, exitCode)
	return outValue, exitCode
}

// ShowMessage displays a message using minui-presenter
func ShowMessage(message string, timeout string) {
	logging.LogDebug("Showing message: %s (timeout: %s)", message, timeout)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return
	}

	// Use explicit path to minui-presenter
	minuiPresenterPath := filepath.Join(cwd, "minui-presenter")

	args := []string{"--message", message, "--timeout", timeout}
	cmd := exec.Command(minuiPresenterPath, args...)
	err = cmd.Run()

	if err != nil {
		logging.LogDebug("minui-presenter error: %v", err)
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() != 124 {
			fmt.Printf("Failed to run minui-presenter: %v\n", err)
		}
	}
}