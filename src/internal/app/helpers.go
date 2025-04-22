// src/internal/app/helpers.go
// Helper functions for the application

package app

import (
	"os"
)

// GetWorkingDir returns the current working directory
func GetWorkingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

// GetCatalogDir returns the path to the catalog directory
func GetCatalogDir() string {
	return GetWorkingDir() + "/Catalog"
}

// GetThemesDir returns the path to the themes directory
func GetThemesDir() string {
	return GetWorkingDir() + "/Themes"
}

// GetComponentsDir returns the path to the components directory
func GetComponentsDir() string {
	return GetWorkingDir() + "/Components"
}

// GetExportsDir returns the path to the exports directory
func GetExportsDir() string {
	return GetWorkingDir() + "/Exports"
}