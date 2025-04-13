// src/internal/themes/default.go
// Default theme operations

package themes

import (
	"fmt"
	"os"
	"path/filepath"

	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
)

// OverwriteWithDefaultTheme applies the default black theme to all backgrounds
func OverwriteWithDefaultTheme() error {
	logging.LogDebug("Applying default black theme")

	// Get current directory for absolute paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Default black background path - updated to new Wallpapers directory
	defaultBg := filepath.Join(cwd, "Wallpapers", "Default", "bg.png")
	logging.LogDebug("Default theme background path: %s", defaultBg)

	// Check if the default background exists
	_, err = os.Stat(defaultBg)
	if err != nil {
		logging.LogDebug("Default background not found: %v", err)
		return fmt.Errorf("default background not found: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Apply to root .media directory
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", "bg.png")
	logging.LogDebug("Copying to root .media: %s", rootMediaBg)
	if err := CopyFile(defaultBg, rootMediaBg); err != nil {
		logging.LogDebug("Error copying to root .media: %v", err)
		return fmt.Errorf("failed to copy background to root .media: %w", err)
	}

	// Also apply to root directory
	rootBg := filepath.Join(systemPaths.Root, "bg.png")
	logging.LogDebug("Copying to root: %s", rootBg)
	if err := CopyFile(defaultBg, rootBg); err != nil {
		logging.LogDebug("Error copying to root: %v", err)
		return fmt.Errorf("failed to copy background to root: %w", err)
	}

	// Apply to Recently Played
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
	logging.LogDebug("Copying to Recently Played: %s", rpBg)
	if err := CopyFile(defaultBg, rpBg); err != nil {
		logging.LogDebug("Error copying to Recently Played: %v", err)
		return fmt.Errorf("failed to copy background to Recently Played: %w", err)
	}

	// Apply to Tools
	toolsBg := filepath.Join(systemPaths.Tools, ".media", "bg.png")
	logging.LogDebug("Copying to Tools: %s", toolsBg)
	if err := CopyFile(defaultBg, toolsBg); err != nil {
		logging.LogDebug("Error copying to Tools: %v", err)
		return fmt.Errorf("failed to copy background to Tools: %w", err)
	}

	// Apply to all system directories
	for _, system := range systemPaths.Systems {
		systemBg := filepath.Join(system.MediaPath, "bg.png")
		logging.LogDebug("Copying to system %s: %s", system.Name, systemBg)
		if err := CopyFile(defaultBg, systemBg); err != nil {
			logging.LogDebug("Error copying to system %s: %v", system.Name, err)
			return fmt.Errorf("failed to copy background to %s: %w", system.Name, err)
		}
	}

	logging.LogDebug("Default theme applied successfully")
	return nil
}

// DeleteAllBackgrounds deletes all background image files from the system
func DeleteAllBackgrounds() error {
	logging.LogDebug("Deleting all background images")

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Delete from root .media directory
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", "bg.png")
	logging.LogDebug("Removing from root .media: %s", rootMediaBg)
	if err := os.Remove(rootMediaBg); err != nil && !os.IsNotExist(err) {
		logging.LogDebug("Error removing from root .media: %v", err)
		return fmt.Errorf("failed to remove background from root .media: %w", err)
	}

	// Delete from root directory
	rootBg := filepath.Join(systemPaths.Root, "bg.png")
	logging.LogDebug("Removing from root: %s", rootBg)
	if err := os.Remove(rootBg); err != nil && !os.IsNotExist(err) {
		logging.LogDebug("Error removing from root: %v", err)
		return fmt.Errorf("failed to remove background from root: %w", err)
	}

	// Delete from Recently Played
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
	logging.LogDebug("Removing from Recently Played: %s", rpBg)
	if err := os.Remove(rpBg); err != nil && !os.IsNotExist(err) {
		logging.LogDebug("Error removing from Recently Played: %v", err)
		return fmt.Errorf("failed to remove background from Recently Played: %w", err)
	}

	// Delete from Tools
	toolsBg := filepath.Join(systemPaths.Tools, ".media", "bg.png")
	logging.LogDebug("Removing from Tools: %s", toolsBg)
	if err := os.Remove(toolsBg); err != nil && !os.IsNotExist(err) {
		logging.LogDebug("Error removing from Tools: %v", err)
		return fmt.Errorf("failed to remove background from Tools: %w", err)
	}

	// Delete from all system directories
	for _, system := range systemPaths.Systems {
		systemBg := filepath.Join(system.MediaPath, "bg.png")
		logging.LogDebug("Removing from system %s: %s", system.Name, systemBg)
		if err := os.Remove(systemBg); err != nil && !os.IsNotExist(err) {
			logging.LogDebug("Error removing from system %s: %v", system.Name, err)
			return fmt.Errorf("failed to remove background from %s: %w", system.Name, err)
		}
	}

	logging.LogDebug("All background images removed successfully")
	return nil
}