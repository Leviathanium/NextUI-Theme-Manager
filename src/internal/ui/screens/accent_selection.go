// src/internal/ui/screens/accent_selection.go
// Implementation of the accent color selection screen with individual color settings

package screens

import (
	"encoding/json"
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/accents"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// AccentSelectionScreen displays individual accent color settings
func AccentSelectionScreen() (string, int) {
	// First, get current accent colors
	currentColors, err := accents.GetCurrentColors()
	if err != nil {
		logging.LogDebug("Error reading current colors: %v", err)
		// Use default colors if we can't read current ones
		currentColors = &accents.ThemeColor{
			Color1: "0xFFFFFF",
			Color2: "0x9B2257",
			Color4: "0xFFFFFF",
		}
	}

	// We'll extract all unique colors from the predefined themes
	var primaryColors []string
	var secondaryColors []string
	var textColors []string

	// Track the index of the current colors in the options arrays
	primaryColorIndex := 0
	secondaryColorIndex := 0
	textColorIndex := 0

	// Create maps to track unique colors
	primarySet := make(map[string]bool)
	secondarySet := make(map[string]bool)
	textSet := make(map[string]bool)

	// Process all predefined themes
	for _, theme := range accents.PredefinedThemes {
		// Primary color (Color1)
		primaryHex := "#" + strings.TrimPrefix(theme.Color1, "0x")
		if !primarySet[primaryHex] {
			primarySet[primaryHex] = true
			primaryColors = append(primaryColors, primaryHex)
		}

		// Secondary color (Color2)
		secondaryHex := "#" + strings.TrimPrefix(theme.Color2, "0x")
		if !secondarySet[secondaryHex] {
			secondarySet[secondaryHex] = true
			secondaryColors = append(secondaryColors, secondaryHex)
		}

		// Text color (Color4)
		textHex := "#" + strings.TrimPrefix(theme.Color4, "0x")
		if !textSet[textHex] {
			textSet[textHex] = true
			textColors = append(textColors, textHex)
		}
	}

	// Find indices of current colors
	currentPrimaryHex := "#" + strings.TrimPrefix(currentColors.Color1, "0x")
	currentSecondaryHex := "#" + strings.TrimPrefix(currentColors.Color2, "0x")
	currentTextHex := "#" + strings.TrimPrefix(currentColors.Color4, "0x")

	for i, color := range primaryColors {
		if color == currentPrimaryHex {
			primaryColorIndex = i
			break
		}
	}

	for i, color := range secondaryColors {
		if color == currentSecondaryHex {
			secondaryColorIndex = i
			break
		}
	}

	for i, color := range textColors {
		if color == currentTextHex {
			textColorIndex = i
			break
		}
	}

	logging.LogDebug("Available colors - Primary: %d, Secondary: %d, Text: %d",
		len(primaryColors), len(secondaryColors), len(textColors))

	// Create the menu items
	var items []map[string]interface{}

	// Primary Color item
	primaryItem := map[string]interface{}{
		"name": "Primary Color",
		"options": primaryColors,
		"selected": primaryColorIndex,
		"features": map[string]interface{}{
			"draw_arrows": true,  // Show left/right arrows
		},
	}
	items = append(items, primaryItem)

	// Secondary Color item
	secondaryItem := map[string]interface{}{
		"name": "Secondary Color",
		"options": secondaryColors,
		"selected": secondaryColorIndex,
		"features": map[string]interface{}{
			"draw_arrows": true,  // Show left/right arrows
		},
	}
	items = append(items, secondaryItem)

	// Text Color item
	textItem := map[string]interface{}{
		"name": "Text Color",
		"options": textColors,
		"selected": textColorIndex,
		"features": map[string]interface{}{
			"draw_arrows": true,  // Show left/right arrows
		},
	}
	items = append(items, textItem)

	// Create Apply button as a separate item
	applyItem := map[string]interface{}{
		"name": "Apply Changes",
		"features": map[string]interface{}{
			"confirm_text": "APPLY",
		},
	}
	items = append(items, applyItem)

	// Store selected indices for use when applying changes
	app.SetColorSelections(primaryColorIndex, secondaryColorIndex, textColorIndex)

	// Create a wrapper object with the items array
	wrapper := map[string]interface{}{
		"items": items,
		"selected": 0,  // Start with Primary Color selected
	}

	// Convert to JSON
	jsonData, err := json.Marshal(wrapper)
	if err != nil {
		logging.LogDebug("Error creating JSON: %v", err)
		ui.ShowMessage("Error creating accent settings menu.", "3")
		return "", 1
	}

	logging.LogDebug("Displaying accent settings menu")
	return ui.DisplayMinUiList(string(jsonData), "json", "Accent Settings", "--item-key", "items")
}

// HandleAccentSelection processes the user's accent settings selections
func HandleAccentSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleAccentSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected "Apply Changes"
		if selection == "Apply Changes" {
			// Get the stored indices
			primaryIndex, secondaryIndex, textIndex := app.GetColorSelections()

			// Create a custom theme with the selected colors
			selectedPrimaryHex := "#FFFFFF"    // Default fallback
			selectedSecondaryHex := "#9B2257"  // Default fallback
			selectedTextHex := "#FFFFFF"       // Default fallback

			// Get unique colors from predefined themes
			var primaryColors []string
			var secondaryColors []string
			var textColors []string

			primarySet := make(map[string]bool)
			secondarySet := make(map[string]bool)
			textSet := make(map[string]bool)

			for _, theme := range accents.PredefinedThemes {
				// Primary color (Color1)
				primaryHex := "#" + strings.TrimPrefix(theme.Color1, "0x")
				if !primarySet[primaryHex] {
					primarySet[primaryHex] = true
					primaryColors = append(primaryColors, primaryHex)
				}

				// Secondary color (Color2)
				secondaryHex := "#" + strings.TrimPrefix(theme.Color2, "0x")
				if !secondarySet[secondaryHex] {
					secondarySet[secondaryHex] = true
					secondaryColors = append(secondaryColors, secondaryHex)
				}

				// Text color (Color4)
				textHex := "#" + strings.TrimPrefix(theme.Color4, "0x")
				if !textSet[textHex] {
					textSet[textHex] = true
					textColors = append(textColors, textHex)
				}
			}

			// Get the selected colors if indices are valid
			if primaryIndex >= 0 && primaryIndex < len(primaryColors) {
				selectedPrimaryHex = primaryColors[primaryIndex]
			}

			if secondaryIndex >= 0 && secondaryIndex < len(secondaryColors) {
				selectedSecondaryHex = secondaryColors[secondaryIndex]
			}

			if textIndex >= 0 && textIndex < len(textColors) {
				selectedTextHex = textColors[textIndex]
			}

			// Convert from #RRGGBB to 0xRRGGBB format
			selectedPrimaryColor := "0x" + strings.TrimPrefix(selectedPrimaryHex, "#")
			selectedSecondaryColor := "0x" + strings.TrimPrefix(selectedSecondaryHex, "#")
			selectedTextColor := "0x" + strings.TrimPrefix(selectedTextHex, "#")

			// Create custom theme
			customTheme := &accents.ThemeColor{
				Name: "Custom",
				Color1: selectedPrimaryColor,
				Color2: selectedSecondaryColor,
				Color4: selectedTextColor,
				// Other colors will be set to defaults or copied from existing theme
			}

			// Apply the custom theme
			logging.LogDebug("Applying custom accent theme with colors - Primary: %s, Secondary: %s, Text: %s",
				selectedPrimaryColor, selectedSecondaryColor, selectedTextColor)

			err := accents.ApplyThemeColors(customTheme)
			if err != nil {
				logging.LogDebug("Error applying custom theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				ui.ShowMessage("Applied custom accent colors", "3")
			}
		}
		return app.Screens.MainMenu

	case 4:
		// Handle option change event
		// Unfortunately, minui-list doesn't tell us what changed, so we can't update our stored selections
		// This is a limitation of the current design
		return app.Screens.AccentSelection

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.AccentSelection
}