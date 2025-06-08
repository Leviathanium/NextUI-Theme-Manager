// internal/ui/wrappers.go
package ui

// Add this import at the top of the file
import (
    "bytes"
    "encoding/json" // New import for JSON handling
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "thememanager/internal/app"
)

// File: src/internal/ui/wrappers.go
// Replace the ShowThemeGallery function with this improved version

// ShowThemeGallery displays a gallery of themes with previews using minui-presenter
func ShowThemeGallery(themes []map[string]string, title string) (string, int) {
    app.LogDebug("Showing theme gallery: %s", title)

    // Get current directory
    cwd, err := os.Getwd()
    if err != nil {
        app.LogDebug("Error getting current directory: %v", err)
        return "", 1
    }

    // Create a temporary file for theme gallery data
    tempFile, err := os.CreateTemp("", "theme-gallery-*")
    if err != nil {
        app.LogDebug("Error creating temp file: %v", err)
        return "", 1
    }
    galleryFilePath := tempFile.Name()
    defer os.Remove(galleryFilePath)

    // Create gallery JSON structure
    gallery := map[string]interface{}{
        "items":    make([]map[string]interface{}, 0),
        "selected": 0,
    }

    // Add each theme to the gallery
    for _, theme := range themes {
        name := theme["name"]
        author := theme["author"]
        previewPath := theme["preview"]
        isValid := theme["is_valid"] == "true"

        // Create item text with theme info
        // Create item text with theme info (no description)
        itemText := name
        if author != "" {
            itemText += " by " + author
        }

        // Add validation warning if theme is invalid
        if !isValid {
            itemText += "\n\n⚠️ WARNING: This theme has missing or invalid fields and can't be applied."
            itemText += "\nUse 'Export Theme' to create a backup before editing the manifest."
        }

        // Create gallery item
        item := map[string]interface{}{
            "text":             itemText,
            "background_color": "#000000",
            "show_pill":        true,
            "alignment":        "top",
        }

        // Add preview image if it exists
        if previewPath != "" {
            // Verify image exists
            if _, err := os.Stat(previewPath); err == nil {
                item["background_image"] = previewPath
            } else {
                app.LogDebug("Warning: Preview image not found at %s", previewPath)
                // No preview available
                item["text"] = itemText + "\n\n(Preview image not found)"
            }
        } else {
            // No preview image, add a note to the text
            item["text"] = itemText + "\n\n(No preview image available)"
        }

        gallery["items"] = append(gallery["items"].([]map[string]interface{}), item)
    }

    // Check if we have any items
    if len(gallery["items"].([]map[string]interface{})) == 0 {
        app.LogDebug("No themes to display in gallery")
        return "", 1
    }

    // Convert gallery data to JSON
    galleryData, err := json.Marshal(gallery)
    if err != nil {
        app.LogDebug("Error creating gallery JSON: %v", err)
        return "", 1
    }

    // Write gallery data to temp file
    err = tempFile.Close()
    if err != nil {
        app.LogDebug("Error closing temp file: %v", err)
        return "", 1
    }

    err = os.WriteFile(galleryFilePath, galleryData, 0644)
    if err != nil {
        app.LogDebug("Error writing gallery data: %v", err)
        return "", 1
    }

    // Create temp file for output
    tempOutFile, err := os.CreateTemp("", "gallery-output-*")
    if err != nil {
        app.LogDebug("Error creating output temp file: %v", err)
        return "", 1
    }
    outputPath := tempOutFile.Name()
    tempOutFile.Close()
    defer os.Remove(outputPath)

    // Path to minui-presenter
    minuiPresenterPath := filepath.Join(cwd, "minui-presenter")

    // Create command
    args := []string{
        "--file", galleryFilePath,
        "--confirm-text", "APPLY",
        "--confirm-show",
        "--cancel-text", "BACK",
        "--cancel-show",
        "--item-key", "items",
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

    // Get selected item index from gallery file
    if exitCode == 0 {
        // Read gallery file to get updated selected index
        updatedGalleryData, err := os.ReadFile(galleryFilePath)
        if err != nil {
            app.LogDebug("Error reading updated gallery data: %v", err)
            return "", 1
        }

        var updatedGallery map[string]interface{}
        err = json.Unmarshal(updatedGalleryData, &updatedGallery)
        if err != nil {
            app.LogDebug("Error parsing updated gallery data: %v", err)
            return "", 1
        }

        selectedIndex := int(updatedGallery["selected"].(float64))
        if selectedIndex >= 0 && selectedIndex < len(themes) {
            // Check if theme is valid before returning
            if themes[selectedIndex]["is_valid"] != "true" {
                app.LogDebug("User attempted to apply invalid theme: %s", themes[selectedIndex]["name"])
                // Show error message
                ShowMessage("Cannot apply theme with missing or invalid fields. Please fix the manifest first.", "3")
                return "", 2 // Return special code to go back to theme list
            }

            // Return the name of the selected theme
            return themes[selectedIndex]["name"], exitCode
        }
    }

    return "", exitCode
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