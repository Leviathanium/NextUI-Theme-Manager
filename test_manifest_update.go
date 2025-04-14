package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nextui-themes/internal/themes"
)

func main() {
	// Current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Test theme path
	themePath := filepath.Join(cwd, "test_theme")

	fmt.Printf("Updating manifest for test theme at: %s\n", themePath)

	// Update the manifest
	if err := themes.UpdateThemeManifest(themePath); err != nil {
		log.Fatalf("Error updating manifest: %v", err)
	}

	// Read the resulting manifest file
	manifestPath := filepath.Join(themePath, "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Fatalf("Error reading updated manifest: %v", err)
	}

	// Print the manifest
	fmt.Println("\nUpdated manifest:")
	fmt.Println(string(manifestData))

	fmt.Println("\nManifest update test completed successfully!")
}
