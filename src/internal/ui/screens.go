// internal/ui/screen.go
package ui

import (
	"strings"
)

// ShowConfirmDialog displays a confirmation dialog with Yes/No options
// The optional defaultYes parameter can be used to set the default selection
func ShowConfirmDialog(message string, defaultYes ...bool) (string, int) {
	options := []string{
		"Yes",
		"No",
	}

	// Optional parameter to set default selection
	extraArgs := []string{}
	if len(defaultYes) > 0 && defaultYes[0] {
		extraArgs = append(extraArgs, "--selected", "0")
	} else {
		extraArgs = append(extraArgs, "--selected", "1") // Default to "No" for safety
	}

	return DisplayMinUiList(
		strings.Join(options, "\n"),
		"text",
		message,
		extraArgs...,
	)
}