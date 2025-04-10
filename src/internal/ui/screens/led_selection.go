// src/internal/ui/screens/led_selection.go
// Implementation of the LED settings selection screen

package screens

import (
	"encoding/json"
	"fmt"
	"strings"

	"nextui-themes/internal/app"
	"nextui-themes/internal/leds"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"
)

// LEDSelectionScreen displays LED color and effect settings
func LEDSelectionScreen() (string, int) {
	// Get all unique LED colors from the predefined themes
	var ledColors []string
	colorSet := make(map[string]bool)

	for _, theme := range leds.PredefinedThemes {
		colorHex := "#" + strings.TrimPrefix(theme.Color, "0x")
		if !colorSet[colorHex] {
			colorSet[colorHex] = true
			ledColors = append(ledColors, colorHex)
		}
	}

	// Create the menu items
	var items []map[string]interface{}

	// LED Color item - allows scrolling through all colors
	colorItem := map[string]interface{}{
		"name": "LED Color",
		"options": ledColors,
		"selected": 0, // Default to first color
		"features": map[string]interface{}{
			"draw_arrows": true,  // Show left/right arrows
		},
	}
	items = append(items, colorItem)

	// LED Effect item - allows choosing between Static and Breathing
	effectItem := map[string]interface{}{
		"name": "LED Effect",
		"options": []string{"Static", "Breathing"},
		"selected": 0, // Default to Static
		"features": map[string]interface{}{
			"draw_arrows": true,  // Show left/right arrows
		},
	}
	items = append(items, effectItem)

	// Create Apply button as a separate item
	applyItem := map[string]interface{}{
		"name": "Apply Changes",
		"features": map[string]interface{}{
			"confirm_text": "APPLY",
		},
	}
	items = append(items, applyItem)

	// Store selected indices
	app.SetLEDSelections(0, 0)

	// Create a wrapper object with the items array
	wrapper := map[string]interface{}{
		"items": items,
		"selected": 0,  // Start with LED Color selected
	}

	// Convert to JSON
	jsonData, err := json.Marshal(wrapper)
	if err != nil {
		logging.LogDebug("Error creating JSON: %v", err)
		ui.ShowMessage("Error creating LED settings menu.", "3")
		return "", 1
	}

	logging.LogDebug("Displaying LED settings menu")
	return ui.DisplayMinUiList(string(jsonData), "json", "LED Settings", "--item-key", "items")
}

// HandleLEDSelection processes the user's LED settings selections
func HandleLEDSelection(selection string, exitCode int) app.Screen {
	logging.LogDebug("HandleLEDSelection called with selection: '%s', exitCode: %d", selection, exitCode)

	switch exitCode {
	case 0:
		// User selected "Apply Changes"
		if selection == "Apply Changes" {
			// Get the stored indices
			colorIndex, effectIndex := app.GetLEDSelections()

			// Get unique colors from predefined themes
			var ledColors []string
			colorSet := make(map[string]bool)

			for _, theme := range leds.PredefinedThemes {
				colorHex := "#" + strings.TrimPrefix(theme.Color, "0x")
				if !colorSet[colorHex] {
					colorSet[colorHex] = true
					ledColors = append(ledColors, colorHex)
				}
			}

			// Default fallbacks
			selectedColorHex := "#FFFFFF"
			selectedEffect := leds.Static

			// Get the selected color if index is valid
			if colorIndex >= 0 && colorIndex < len(ledColors) {
				selectedColorHex = ledColors[colorIndex]
			}

			// Get the selected effect
			if effectIndex == 1 {
				selectedEffect = leds.Breathe
			}

			// Convert from #RRGGBB to 0xRRGGBB format
			selectedColor := "0x" + strings.TrimPrefix(selectedColorHex, "#")

			// Create custom theme
			customTheme := &leds.LEDTheme{
				Name: "Custom",
				Color: selectedColor,
				Effect: selectedEffect,
			}

			// Apply the custom theme
			logging.LogDebug("Applying custom LED theme - Color: %s, Effect: %d",
				selectedColor, selectedEffect)

			err := leds.ApplyLEDTheme(customTheme)
			if err != nil {
				logging.LogDebug("Error applying LED theme: %v", err)
				ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
			} else {
				effectName := "Static"
				if selectedEffect == leds.Breathe {
					effectName = "Breathing"
				}
				ui.ShowMessage(fmt.Sprintf("Applied LED settings - %s with %s effect", selectedColor, effectName), "3")
			}
		}
		return app.Screens.MainMenu

	case 1, 2:
		// User pressed cancel or back
		return app.Screens.MainMenu
	}

	return app.Screens.LEDSelection
}