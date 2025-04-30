// internal/themes/gallery.go
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"thememanager/internal/logging"
	"thememanager/internal/system"
	"thememanager/internal/ui"
)

// ShowThemeGallery displays a gallery of available themes
func ShowThemeGallery() (string, int) {
	logging.LogDebug("Showing theme gallery")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read catalog directory for themes
	catalogPath := filepath.Join(cwd, "Catalog", "Themes")
	entries, err := os.ReadDir(catalogPath)
	if err != nil {
		logging.LogDebug("Error reading catalog directory: %v", err)
		ui.ShowMessage("No themes found in catalog. Please sync first.", "3")
		return "", 1
	}

	// Create gallery items
	var galleryItems []ui.GalleryItem
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), system.ThemeExtension) {
			themeName := entry.Name()
			themeName = strings.TrimSuffix(themeName, system.ThemeExtension)

			// Get preview image path
			previewPath := filepath.Join(catalogPath, entry.Name(), "preview.png")

			// Get theme information from manifest
			manifestPath := filepath.Join(catalogPath, entry.Name(), "manifest.yml")
			manifest, err := ReadThemeManifest(manifestPath)

			// Default text
			text := themeName

			// Check if theme is already downloaded/installed
			isInstalled := IsThemeDownloaded(themeName)
			if isInstalled {
				text = "[Installed] " + text
			}

			// Add author if available
			if err == nil && manifest.Info.Author != "" {
				text = fmt.Sprintf("%s by %s", text, manifest.Info.Author)
			}

			// Create gallery item
			galleryItems = append(galleryItems, ui.GalleryItem{
				Text:            text,
				BackgroundImage: previewPath,
			})
		}
	}

	// Check if we found any themes
	if len(galleryItems) == 0 {
		ui.ShowMessage("No themes found in catalog. Please sync first.", "3")
		return "", 1
	}

	// Display the gallery
	return ui.DisplayImageGallery(galleryItems, "Available Themes")
}

// ShowOverlayGallery displays a gallery of available overlay packs
func ShowOverlayGallery() (string, int) {
	logging.LogDebug("Showing overlay gallery")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read catalog directory for overlays
	catalogPath := filepath.Join(cwd, "Catalog", "Overlays")
	entries, err := os.ReadDir(catalogPath)
	if err != nil {
		logging.LogDebug("Error reading catalog directory: %v", err)
		ui.ShowMessage("No overlays found in catalog. Please sync first.", "3")
		return "", 1
	}

	// Create gallery items
	var galleryItems []ui.GalleryItem
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), system.OverlayExtension) {
			overlayName := entry.Name()
			overlayName = strings.TrimSuffix(overlayName, system.OverlayExtension)

			// Get preview image path
			previewPath := filepath.Join(catalogPath, entry.Name(), "preview.png")

			// Get overlay information from manifest
			manifestPath := filepath.Join(catalogPath, entry.Name(), "manifest.yml")
			manifest, err := ReadOverlayManifest(manifestPath)

			// Default text
			text := overlayName

			// Check if overlay is already downloaded/installed
			isInstalled := IsOverlayDownloaded(overlayName)
			if isInstalled {
				text = "[Installed] " + text
			}

			// Add author if available
			if err == nil && manifest.Info.Author != "" {
				text = fmt.Sprintf("%s by %s", text, manifest.Info.Author)
			}

			// Create gallery item
			galleryItems = append(galleryItems, ui.GalleryItem{
				Text:            text,
				BackgroundImage: previewPath,
			})
		}
	}

	// Check if we found any overlays
	if len(galleryItems) == 0 {
		ui.ShowMessage("No overlays found in catalog. Please sync first.", "3")
		return "", 1
	}

	// Display the gallery
	return ui.DisplayImageGallery(galleryItems, "Available Overlays")
}

// ShowThemeBackupGallery displays a gallery of theme backups
func ShowThemeBackupGallery() (string, int) {
	logging.LogDebug("Showing theme backup gallery")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read backups directory for themes
	backupsPath := filepath.Join(cwd, "Backups", "Themes")
	entries, err := os.ReadDir(backupsPath)
	if err != nil {
		logging.LogDebug("Error reading backups directory: %v", err)
		ui.ShowMessage("No theme backups found.", "3")
		return "", 1
	}

	// Create gallery items
	var galleryItems []ui.GalleryItem
	for _, entry := range entries {
		if entry.IsDir() {
			backupName := entry.Name()

			// Try to find a preview image in the backup
			previewPath := filepath.Join(backupsPath, backupName, "preview.png")
			if _, err := os.Stat(previewPath); os.IsNotExist(err) {
				// No preview, check if there's a screenshot
				previewPath = filepath.Join(backupsPath, backupName, "screenshot.png")
				if _, err := os.Stat(previewPath); os.IsNotExist(err) {
					// No screenshot either, use a blank preview
					previewPath = ""
				}
			}

			// Create gallery item
			galleryItems = append(galleryItems, ui.GalleryItem{
				Text:            backupName,
				BackgroundImage: previewPath,
			})
		}
	}

	// Check if we found any backups
	if len(galleryItems) == 0 {
		ui.ShowMessage("No theme backups found.", "3")
		return "", 1
	}

	// Display the gallery
	return ui.DisplayImageGallery(galleryItems, "Theme Backups")
}

// ShowOverlayBackupGallery displays a gallery of overlay backups
func ShowOverlayBackupGallery() (string, int) {
	logging.LogDebug("Showing overlay backup gallery")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		ui.ShowMessage(fmt.Sprintf("Error: %s", err), "3")
		return "", 1
	}

	// Read backups directory for overlays
	backupsPath := filepath.Join(cwd, "Backups", "Overlays")
	entries, err := os.ReadDir(backupsPath)
	if err != nil {
		logging.LogDebug("Error reading backups directory: %v", err)
		ui.ShowMessage("No overlay backups found.", "3")
		return "", 1
	}

	// Create gallery items
	var galleryItems []ui.GalleryItem
	for _, entry := range entries {
		if entry.IsDir() {
			backupName := entry.Name()

			// Try to find a preview image in the backup
			previewPath := filepath.Join(backupsPath, backupName, "preview.png")
			if _, err := os.Stat(previewPath); os.IsNotExist(err) {
				// No preview, use a blank preview
				previewPath = ""
			}

			// Create gallery item
			galleryItems = append(galleryItems, ui.GalleryItem{
				Text:            backupName,
				BackgroundImage: previewPath,
			})
		}
	}

	// Check if we found any backups
	if len(galleryItems) == 0 {
		ui.ShowMessage("No overlay backups found.", "3")
		return "", 1
	}

	// Display the gallery
	return ui.DisplayImageGallery(galleryItems, "Overlay Backups")
}