// internal/ui/presenter.go
package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
    "time"
	"thememanager/internal/logging"
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

	// Add any extra arguments
	if len(extraArgs) > 0 {
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
			outValue = string(selectionBytes)
			outValue = removeNewlines(outValue)
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

// Helper function to remove newlines from a string
func removeNewlines(s string) string {
	return string(bytes.TrimRight([]byte(s), "\r\n"))
}