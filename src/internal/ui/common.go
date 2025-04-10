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

	"nextui-themes/internal/logging"
)

// Selection represents the result of a UI interaction
type Selection struct {
	Value    string
	Code     int
	Error    error
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

// DisplayMinUiListWithFont displays a list of items using minui-list with a custom font
func DisplayMinUiListWithFont(list string, format string, title string, fontPath string) (string, int) {
	logging.LogDebug("Displaying minui-list with title: %s and font: %s", title, fontPath)
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
	args := []string{
		"--format", format,
		"--title", title,
		"--file", inputPath,
		"--write-location", outputPath,
		"--font-default", fontPath,
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

	args := []string{"--message", message, "--timeout", timeout}
	cmd := exec.Command("minui-presenter", args...)
	err := cmd.Run()

	if err != nil {
		logging.LogDebug("minui-presenter error: %v", err)
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() != 124 {
			fmt.Printf("Failed to run minui-presenter: %v\n", err)
		}
	}
}