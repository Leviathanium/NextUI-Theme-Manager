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

// FontSelectionScreen displays the font target selection screen with clearer options
func FontSelectionScreen() (string, int) {
	// Create menu with replace and restore options
	var menu []string

	// Add replace options
	menu = append(menu, "Replace OG Font")
	menu = append(menu, "Replace Next Font")

	// Add restore options, if available
	if fonts.BackupExists(fonts.OGFont) {
		menu = append(menu, "Restore OG Font")
	}

	if fonts.BackupExists(fonts.NextFont) {
		menu = append(menu, "Restore Next Font")
	}

	return ui.DisplayMinUiList(strings.Join(menu, "\n"), "text", "Font Options")
}

// HandleFontSelection processes the user's font action selection
func HandleFontSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleFontSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// Check if this is a restore option
		if strings.HasPrefix(selection, "Restore") {
			// Handle restore options directly
			var fontType fonts.FontType

			if strings.Contains(selection, "OG") {
				fontType = fonts.OGFont
			} else {
				fontType = fonts.NextFont
			}

			err := fonts.RestoreFont(fontType)
			if err != nil {
				logging.LogDebug("Error restoring font: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				// Show success message
				ui.ShowMessage(fmt.Sprintf("Successfully restored original %s",
					strings.TrimPrefix(selection, "Restore ")), "3")
			}

			return app.Screens.FontSelection
		} else {
			// For replace options, move to font list
			app.SetSelectedFontSlot(selection)
			return app.Screens.FontList
		}
	case 1, 2:
		// User pressed cancel or back
		return app.Screens.CustomizationMenu
	}

	return app.Screens.FontSelection
}

// FontListScreen shows the list of available fonts for the selected slot
func FontListScreen() (string, int) {
	// List available fonts
	fontsList, err := fonts.ListFonts()
	if err != nil {
		logging.LogDebug("Error loading fonts: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error loading fonts: %s", err), "3")
		return "", 1
	}

	if len(fontsList) == 0 {
		logging.LogDebug("No fonts found")
		ui.ShowMessage("No fonts found. Add TTF or OTF fonts to the Fonts directory.", "3")
		return "", 1
	}

	// Remove restore options from the fonts list
	var filteredFonts []string
	for _, font := range fontsList {
		if !strings.HasPrefix(font, "Restore") {
			filteredFonts = append(filteredFonts, font)
		}
	}

	if len(filteredFonts) == 0 {
		ui.ShowMessage("No custom fonts found. Add TTF or OTF fonts to the Fonts directory.", "3")
		return "", 1
	}

	selectedSlot := app.GetSelectedFontSlot()
	logging.LogDebug("Displaying font selection with %d options for: %s",
		len(filteredFonts), selectedSlot)

	return ui.DisplayMinUiList(strings.Join(filteredFonts, "\n"), "text",
		fmt.Sprintf("Select Font for %s", selectedSlot))
}

// HandleFontList processes the user's font selection
func HandleFontList(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleFontList called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected a font
		app.SetSelectedFont(selection)
		return app.Screens.FontPreview
	case 1, 2:
		// User pressed cancel or back
		return app.Screens.FontSelection
	}

	return app.Screens.FontList
}

// FontPreviewScreen shows a preview of the selected font with simplified text
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

    // Create preview text with simplified options
	previewText := []string{
		fmt.Sprintf("Apply %s", fontName),
		"Cancel",
	}

	// Display the preview without using the custom font
	// This avoids the segmentation fault issues
	logging.LogDebug("Displaying font preview for: %s", fontName)
	return ui.DisplayMinUiList(
		strings.Join(previewText, "\n"),
		"text",
		fmt.Sprintf("Font Preview - %s", fontName),
	)
}

// HandleFontPreview processes the user's font preview action with simplified messages
func HandleFontPreview(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleFontPreview called with selection: '%s', exitCode: %d", selection, exitCode)

	// Get selected font and slot
	fontName := app.GetSelectedFont()
	fontSlot := app.GetSelectedFontSlot()

	// Handle segmentation fault from minui-list
	if exitCode < 0 {
		logging.LogDebug("Detected crash in font preview, returning to font selection")
		ui.ShowMessage("Font preview failed. You can still apply the font.", "3")
		return app.Screens.FontList
	}

	switch exitCode {
	case 0:
		// More robust string comparison - convert to lowercase and check contains
		selectionLower := strings.ToLower(selection)
		if strings.Contains(selectionLower, "apply") {
			// Apply the font to the selected slot
			logging.LogDebug("Applying font: %s to slot: %s", fontName, fontSlot)

			// Determine which font type to modify
			var fontType fonts.FontType
			if strings.Contains(fontSlot, "OG") {
				fontType = fonts.OGFont
			} else {
				fontType = fonts.NextFont // Default to Next font if not OG
			}

			err := fonts.ApplyFont(fontName, fontType)
			if err != nil {
				logging.LogDebug("Error applying font: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				// Show simplified success message
				ui.ShowMessage(fmt.Sprintf("Applied %s successfully", fontName), "3")
			}
		}
		return app.Screens.FontSelection
	case 1, 2:
		// User pressed cancel or back
		return app.Screens.FontList
	}

	return app.Screens.FontSelection
}