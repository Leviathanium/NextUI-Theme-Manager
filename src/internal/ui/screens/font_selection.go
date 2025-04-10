// src/internal/ui/screens/font_selection.go
// Implementation of the font selection screen

package screens

import (
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/fonts"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// FontSelectionScreen displays available fonts
func FontSelectionScreen() (string, int) {
	// List available fonts
	fontsList, err := fonts.ListFonts()
	if err != nil {
		logging.LogDebug("Error loading fonts: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error loading fonts: %s", err), "3")
		return "", 1
	}

	if len(fontsList) == 0 {
		logging.LogDebug("No fonts found")
		ui.ShowMessage("No fonts found. Add fonts to the Fonts directory.", "3")
		return "", 1
	}

	logging.LogDebug("Displaying font selection with %d options", len(fontsList))
	return ui.DisplayMinUiList(strings.Join(fontsList, "\n"), "text", "Select Font")
}

// HandleFontSelection processes the user's font selection
func HandleFontSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleFontSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a font
		app.SetSelectedFont(selection)
		return app.Screens.FontPreview
	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.FontSelection
}

// FontPreviewScreen shows a preview of the selected font
func FontPreviewScreen() (string, int) {
	fontName := app.GetSelectedFont()
	logging.LogDebug("Previewing font: %s", fontName)

	// Get the path to the font - just to verify it exists
	_, err := fonts.GetFontPath(fontName)
	if err != nil {
		logging.LogDebug("Error getting font path: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

    // Create preview text
	previewText := []string{
		"Apply font: " + fontName,
		"Cancel",
	}

	// Display the preview without using the custom font
	// This avoids the segmentation fault issues
	logging.LogDebug("Displaying font preview for: %s", fontName)
	return ui.DisplayMinUiList(
		strings.Join(previewText, "\n"),
		"text",
		"Font Preview - " + fontName,
	)
}

// HandleFontPreview processes the user's font preview action
func HandleFontPreview(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleFontPreview called with selection: '%s', exitCode: %d", selection, exitCode)

	// Get selected font name
	fontName := app.GetSelectedFont()

	// Handle segmentation fault from minui-list
	if exitCode < 0 {
		logging.LogDebug("Detected crash in font preview, returning to font selection")
		ui.ShowMessage("Font preview failed. You can still apply the font.", "3")
		return app.Screens.FontSelection
	}

switch exitCode {
    case 0:
        // More robust string comparison - convert to lowercase and check contains
        selectionLower := strings.ToLower(selection)
        if strings.Contains(selectionLower, "apply font") {
            // Apply the font
            logging.LogDebug("Applying font: %s", fontName)
            err := fonts.ApplyFont(fontName)
            if err != nil {
                logging.LogDebug("Error applying font: %v", err)
                ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
            } else {
                // Show success message
                ui.ShowMessage(fmt.Sprintf("Successfully applied font: %s", fontName), "3")

                // Return to main menu after applying font
                return app.Screens.MainMenu
            }
        }
        return app.Screens.FontSelection  // This line should be indented this way
    case 1, 2:
        // User pressed cancel or back
        return app.Screens.FontSelection
    }

    return app.Screens.FontSelection
}