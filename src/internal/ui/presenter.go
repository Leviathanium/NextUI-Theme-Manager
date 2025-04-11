// src/internal/ui/presenter.go
// Utilities for using minui-presenter

package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	// Keep track of which item we're showing
	currentIndex := 0

	for {
		// Ensure index is valid
		if currentIndex < 0 {
			currentIndex = len(items) - 1
		} else if currentIndex >= len(items) {
			currentIndex = 0
		}

		// Get current item
		currentItem := items[currentIndex]

		// Create JSON with single item
		jsonData := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"text":             fmt.Sprintf("%s (%d/%d)", currentItem.Text, currentIndex+1, len(items)),
					"background_image": currentItem.BackgroundImage,
					"show_pill":        true,
					"alignment":        "top",
				},
			},
			"selected": 0,
		}

		// Convert to JSON
		jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			logging.LogDebug("Error marshaling JSON: %v", err)
			return "", 1
		}

		// Create a temporary file for the JSON
		tempFile, err := os.CreateTemp("", "gallery-item-*.json")
		if err != nil {
			logging.LogDebug("ERROR: Failed to create temp file: %v", err)
			return "", 1
		}
		jsonPath := tempFile.Name()
		defer os.Remove(jsonPath)

		// Write JSON to temporary file
		if _, err := tempFile.Write(jsonBytes); err != nil {
			logging.LogDebug("ERROR: Failed to write to temp file: %v", err)
			tempFile.Close()
			return "", 1
		}
		tempFile.Close()

		// Execute minui-presenter for this item with navigation buttons
		args := []string{
			"--file", jsonPath,
			"--confirm-text", "SELECT",
			"--confirm-show",
			"--cancel-text", "BACK",
			"--cancel-show",
		}

		// Add action button for "next" navigation
		args = append(args,
			"--action-button", "X",
			"--action-text", "NEXT",
			"--action-show")

		// Add inaction button for "previous" navigation
		args = append(args,
			"--inaction-button", "Y",
			"--inaction-text", "PREV",
			"--inaction-show")

		minuiPresenterPath := filepath.Join(cwd, "minui-presenter")
		cmd := exec.Command(minuiPresenterPath, args...)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		// Run minui-presenter
		err = cmd.Run()
		exitCode := 0
		if err != nil {
			exitCode = cmd.ProcessState.ExitCode()
		}

		// Log stderr output if any
		stderrOutput := stderr.String()
		if stderrOutput != "" {
			logging.LogDebug("stderr: %s", stderrOutput)
		}

		logging.LogDebug("Exit code: %d for item %d: %s", exitCode, currentIndex, currentItem.Text)

		// Handle exit code
		switch exitCode {
		case 0: // User selected THIS item
			logging.LogDebug("User selected item: %s (index: %d)", currentItem.Text, currentIndex)
			return currentItem.Text, 0

		case 2: // User cancelled
			logging.LogDebug("User cancelled")
			return "", 2

		case 4: // Action button (X) - next item
			logging.LogDebug("User pressed NEXT")
			currentIndex++

		case 5: // Inaction button (Y) - previous item
			logging.LogDebug("User pressed PREV")
			currentIndex--

		case 124, 130, 143: // Special exit codes
			return "", exitCode

		default: // Any other exit code, default to next
			logging.LogDebug("Unknown exit code: %d, advancing to next item", exitCode)
			currentIndex++
		}
	}
}