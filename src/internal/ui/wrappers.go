// internal/ui/wrappers.go
package ui

// Add this import at the top of the file
import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "encoding/json"
    "thememanager/internal/app"
)

// ShowThemeGallery displays a gallery of themes with previews using manual navigation
func ShowThemeGallery(themes []map[string]string, title string) (string, int) {
    app.LogDebug("Showing theme gallery: %s", title)

    // Check if we have any themes
    if len(themes) == 0 {
        app.LogDebug("No themes to display in gallery")
        return "", 1
    }

    // Get current directory
    cwd, err := os.Getwd()
    if err != nil {
        app.LogDebug("Error getting current directory: %v", err)
        return "", 1
    }

    // Path to minui-presenter
    minuiPresenterPath := filepath.Join(cwd, "minui-presenter")

    currentIndex := 0

    for {
        theme := themes[currentIndex]
        name := theme["name"]
        author := theme["author"]
        previewPath := theme["preview"]
        isValid := theme["is_valid"] == "true"

        app.LogDebug("Showing theme %d/%d: %s", currentIndex+1, len(themes), name)

        // Create item text with theme info
        itemText := name
        if author != "" {
            itemText += " by " + author
        }

        // Add theme counter
        itemText += fmt.Sprintf("\n\nTheme %d of %d", currentIndex+1, len(themes))

        // Add validation warning if theme is invalid
        if !isValid {
            itemText += "\n\n⚠️ WARNING: This theme has missing or invalid fields and can't be applied."
        }

        // Create a temporary file for single theme data
        tempFile, err := os.CreateTemp("", "single-theme-*")
        if err != nil {
            app.LogDebug("Error creating temp file: %v", err)
            return "", 1
        }
        tempFilePath := tempFile.Name()
        defer os.Remove(tempFilePath)

        // Create single-item gallery structure
        singleTheme := map[string]interface{}{
            "items": []map[string]interface{}{
                {
                    "text":             itemText,
                    "background_color": "#000000",
                    "show_pill":        true,
                    "alignment":        "top",
                },
            },
            "selected": 0,
        }

        // Add preview image if available
        if previewPath != "" {
            if _, err := os.Stat(previewPath); err == nil {
                singleTheme["items"].([]map[string]interface{})[0]["background_image"] = previewPath
            } else {
                app.LogDebug("Warning: Preview image not found at %s", previewPath)
                // Add note about missing preview
                singleTheme["items"].([]map[string]interface{})[0]["text"] = itemText + "\n\n(Preview image not found)"
            }
        } else {
            // Add note about no preview
            singleTheme["items"].([]map[string]interface{})[0]["text"] = itemText + "\n\n(No preview image available)"
        }

        // Convert to JSON and write to temp file
        themeData, err := json.Marshal(singleTheme)
        if err != nil {
            app.LogDebug("Error creating theme JSON: %v", err)
            return "", 1
        }

        err = tempFile.Close()
        if err != nil {
            app.LogDebug("Error closing temp file: %v", err)
            return "", 1
        }

        err = os.WriteFile(tempFilePath, themeData, 0644)
        if err != nil {
            app.LogDebug("Error writing theme data: %v", err)
            return "", 1
        }

        // Build command arguments
        args := []string{
            "--file", tempFilePath,
            "--item-key", "items",
            "--confirm-text", "APPLY",
            "--confirm-show",
            "--cancel-text", "BACK",
            "--cancel-show",
        }

        // Add navigation buttons if we have multiple themes
        if len(themes) > 1 {
            args = append(args, "--action-button", "X")
            args = append(args, "--action-text", "NEXT")
            args = append(args, "--action-show")

            args = append(args, "--inaction-button", "Y")
            args = append(args, "--inaction-text", "PREV")
            args = append(args, "--inaction-show")
        }

        // Run minui-presenter
        cmd := exec.Command(minuiPresenterPath, args...)

        var stderrBuf bytes.Buffer
        cmd.Stderr = &stderrBuf

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

        app.LogDebug("minui-presenter exit code: %d", exitCode)

        // Handle the result
        switch exitCode {
        case 0: // Confirm/Apply button pressed
            if !isValid {
                app.LogDebug("User attempted to apply invalid theme: %s", name)
                ShowMessage("Cannot apply theme with missing or invalid fields. Please fix the manifest first.", "3")
                continue // Stay in the loop, don't exit
            }
            app.LogDebug("User selected theme: %s", name)
            return name, 0

        case 2: // Cancel/Back button pressed
            app.LogDebug("User cancelled theme selection")
            return "", 2

        case 4: // Action button (NEXT) pressed
            app.LogDebug("User pressed NEXT")
            currentIndex = (currentIndex + 1) % len(themes) // Wrap around to first
            continue

        case 5: // Inaction button (PREV) pressed
            app.LogDebug("User pressed PREV")
            currentIndex = (currentIndex - 1 + len(themes)) % len(themes) // Wrap around to last
            continue

        default:
            app.LogDebug("Unexpected exit code: %d", exitCode)
            return "", exitCode
        }
    }
}

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
// Modified ShowMessage function to better handle timeouts
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
			// Handle timeout exit code (124) as a normal exit
			if exitCode == 124 {
				app.LogDebug("Message timed out (normal behavior)")
				// Convert timeout to 0 (success) for consistent handling
				exitCode = 0
			}
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

// Modified ShowConfirmDialog function without the --selected parameter
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

	// Write Yes/No options - place "No" first so it's selected by default
	_, err = tempFile.WriteString("No\nYes")
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

	// Create command - removed --selected parameter
	args := []string{
		"--format", "text",
		"--title", message,
		"--file", inputPath,
		"--write-location", outputPath,
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