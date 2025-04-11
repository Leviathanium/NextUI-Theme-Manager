// src/internal/ui/presenter.go
// Utilities for using minui-presenter

package ui

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"nextui-themes/internal/logging"
)

// GalleryItem represents an item in a gallery display
type GalleryItem struct {
	Text            string
	BackgroundImage string
}

// DisplayImageGallery displays a gallery of images using minui-presenter
func DisplayImageGallery(items []GalleryItem, title string) (string, int) {
	logging.LogDebug("Displaying image gallery with %d items and title: %s", len(items), title)

	if len(items) == 0 {
		logging.LogDebug("No items to display in gallery")
		ShowMessage("No items to display", "3")
		return "", 1
	}

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return "", 1
	}

	// Create JSON structure with ALL items together
	// This allows minui-presenter to handle navigation internally
	jsonItems := []map[string]interface{}{}

	for _, item := range items {
		jsonItems = append(jsonItems, map[string]interface{}{
			"text":             item.Text,
			"background_image": item.BackgroundImage,
			"show_pill":        true,
			"alignment":        "top",
		})
	}

	jsonData := map[string]interface{}{
		"items":    jsonItems,
		"selected": 0,
	}

	// Convert to JSON
	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		logging.LogDebug("Error marshaling JSON: %v", err)
		return "", 1
	}

	// Create a temporary file for the JSON content
	tempFile, err := os.CreateTemp("", "gallery-*.json")
	if err != nil {
		logging.LogDebug("ERROR: Failed to create temp file: %v", err)
		return "", 1
	}
	jsonPath := tempFile.Name()
	defer os.Remove(jsonPath)

	// Write JSON content to the file
	if _, err := tempFile.Write(jsonBytes); err != nil {
		logging.LogDebug("ERROR: Failed to write to temp file: %v", err)
		tempFile.Close()
		return "", 1
	}
	tempFile.Close()

	// Execute minui-presenter with the JSON file containing all items
	args := []string{
		"--file", jsonPath,
		"--confirm-text", "SELECT",
		"--confirm-show",
		"--cancel-text", "BACK",
		"--cancel-show",
	}

	minuiPresenterPath := filepath.Join(cwd, "minui-presenter")
	cmd := exec.Command(minuiPresenterPath, args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run minui-presenter
	err = cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	// Log any output
	stdoutOutput := stdoutBuf.String()
	stderrOutput := stderrBuf.String()

	if stderrOutput != "" {
		logging.LogDebug("stderr: %s", stderrOutput)
	}

	if stdoutOutput != "" {
		logging.LogDebug("stdout: %s", stdoutOutput)
	}

	logging.LogDebug("Final exit code: %d", exitCode)

	// Handle exit code
	if exitCode == 0 {
		// User selected something
		// Since we can't get the selected index reliably,
		// we'll check stdout but default to the first item
		selectedIndex := 0

		// If the first item was selected, return it
		if selectedIndex >= 0 && selectedIndex < len(items) {
			return items[selectedIndex].Text, 0
		}
		return items[0].Text, 0
	} else if exitCode == 2 {
		// User cancelled
		return "", 2
	}

	// For any other exit code, treat as cancel
	return "", exitCode
}