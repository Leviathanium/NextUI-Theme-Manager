// src/internal/icons/icons.go
// Icon management for NextUI Theme Manager

package icons

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
)

// Constants for directory paths
const (
	IconsDir = "Icons"        // Base directory for icon packs
	MediaDir = ".media"       // Media directory for ROM systems
)

// Special icon filenames
const (
	ToolsIcon = "tg5040.png"
	RecentlyPlayedIcon = "Recently Played.png"
	CollectionsIcon = "Collections.png"
)

// IconPack represents an icon pack
type IconPack struct {
	Name        string
	Path        string
	Icons       []Icon
	SpecialIcons map[string]string // Map of special icon names to file paths
}

// Icon represents an individual icon file
type Icon struct {
	Path      string
	Name      string
	SystemTag string // The system tag extracted from the icon filename
}

// ListIconPacks lists available icon packs
func ListIconPacks() ([]string, error) {
	var packs []string

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to icon packs directory
	iconsDir := filepath.Join(cwd, IconsDir)
	logging.LogDebug("Scanning icon packs directory: %s", iconsDir)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(iconsDir, 0755); err != nil {
		logging.LogDebug("Error creating icon packs directory: %v", err)
		return nil, fmt.Errorf("error creating icon packs directory: %w", err)
	}

	// Read the directory
	entries, err := os.ReadDir(iconsDir)
	if err != nil {
		logging.LogDebug("Error reading icon packs directory: %v", err)
		return nil, fmt.Errorf("error reading icon packs directory: %w", err)
	}

	// Find directories that contain icon files
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			packPath := filepath.Join(iconsDir, entry.Name())

			// Check if the directory contains any PNG files
			hasIcons := false
			iconEntries, err := os.ReadDir(packPath)
			if err == nil {
				for _, iconEntry := range iconEntries {
					if !iconEntry.IsDir() && strings.HasSuffix(strings.ToLower(iconEntry.Name()), ".png") {
						hasIcons = true
						break
					}
				}
			}

			if hasIcons {
				packs = append(packs, entry.Name())
			}
		}
	}

	logging.LogDebug("Found %d icon packs", len(packs))
	return packs, nil
}

// LoadIconPack loads the icons from a specific pack
func LoadIconPack(packName string) (*IconPack, error) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to icon pack
	packPath := filepath.Join(cwd, IconsDir, packName)
	logging.LogDebug("Loading icon pack: %s", packPath)

	// Check if the directory exists
	if _, err := os.Stat(packPath); os.IsNotExist(err) {
		logging.LogDebug("Icon pack directory does not exist: %s", packPath)
		return nil, fmt.Errorf("icon pack directory does not exist: %s", packPath)
	}

	// Create a new icon pack
	pack := &IconPack{
		Name: packName,
		Path: packPath,
		SpecialIcons: make(map[string]string),
	}

	// Read the directory
	entries, err := os.ReadDir(packPath)
	if err != nil {
		logging.LogDebug("Error reading icon pack directory: %v", err)
		return nil, fmt.Errorf("error reading icon pack directory: %w", err)
	}

	// Regular expression to extract system tag from filename
	re := regexp.MustCompile(`\((.*?)\)`)

	// Find all PNG files
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
			continue
		}

		// Skip hidden files and macOS metadata files (starting with "._")
		if strings.HasPrefix(entry.Name(), ".") {
			logging.LogDebug("Skipping hidden/metadata file: %s", entry.Name())
			continue
		}

		iconPath := filepath.Join(packPath, entry.Name())

		// Check if this is a special icon
		lowercaseName := strings.ToLower(entry.Name())
		switch lowercaseName {
		case strings.ToLower(ToolsIcon):
			pack.SpecialIcons[ToolsIcon] = iconPath
			logging.LogDebug("Found special Tools icon: %s", iconPath)
		case strings.ToLower(RecentlyPlayedIcon):
			pack.SpecialIcons[RecentlyPlayedIcon] = iconPath
			logging.LogDebug("Found special Recently Played icon: %s", iconPath)
		case strings.ToLower(CollectionsIcon):
			pack.SpecialIcons[CollectionsIcon] = iconPath
			logging.LogDebug("Found special Collections icon: %s", iconPath)
		default:
			// Extract system tag if present
			tag := ""
			matches := re.FindStringSubmatch(entry.Name())
			if len(matches) >= 2 {
				tag = matches[1]

				// Add to the list of icons
				pack.Icons = append(pack.Icons, Icon{
					Path:      iconPath,
					Name:      entry.Name(),
					SystemTag: tag,
				})

				logging.LogDebug("Added icon for system tag '%s': %s", tag, entry.Name())
			} else {
				logging.LogDebug("Skipping icon with no system tag: %s", entry.Name())
			}
		}
	}

	logging.LogDebug("Loaded %d regular icons and %d special icons from pack %s",
		len(pack.Icons), len(pack.SpecialIcons), packName)
	return pack, nil
}

// ApplyIconPack applies the selected icon pack to all systems or a specific system
func ApplyIconPack(packName string) error {
	return ApplyIconPackToSystem(packName, "")
}

// ApplyIconPackToSystem applies the selected icon pack to a specific system
// If systemName is empty, applies to all systems
func ApplyIconPackToSystem(packName string, systemName string) error {
	logging.LogDebug("Applying icon pack: %s to system: %s", packName, systemName)

	// If not applying to a specific system, first delete all existing icons to ensure clean application
	if systemName == "" {
		if err := DeleteAllIcons(); err != nil {
			logging.LogDebug("Warning: Failed to delete existing icons: %v", err)
			// Continue anyway, as we'll overwrite icons
		}
	}

	// Load the icon pack
	pack, err := LoadIconPack(packName)
	if err != nil {
		logging.LogDebug("Error loading icon pack: %v", err)
		return fmt.Errorf("error loading icon pack: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Create media directory if it doesn't exist
	mediaPath := filepath.Join(systemPaths.Roms, MediaDir)
	if err := os.MkdirAll(mediaPath, 0755); err != nil {
		logging.LogDebug("Error creating media directory: %v", err)
		return fmt.Errorf("error creating media directory: %w", err)
	}

	// If applying to all systems, apply special icons
	if systemName == "" {
		if err := applySpecialIcons(pack, systemPaths); err != nil {
			logging.LogDebug("Error applying special icons: %v", err)
			// Continue with regular icons
		}
	} else if systemName == "Root" {
		// Handle special case for Root
		if collectionsIconPath, exists := pack.SpecialIcons[CollectionsIcon]; exists {
			rootMediaPath := filepath.Join(systemPaths.Root, MediaDir)
			if err := os.MkdirAll(rootMediaPath, 0755); err != nil {
				logging.LogDebug("Error creating root media directory: %v", err)
				return fmt.Errorf("error creating root media directory: %w", err)
			}

			destPath := filepath.Join(rootMediaPath, CollectionsIcon)
			logging.LogDebug("Copying Collections icon: %s -> %s", collectionsIconPath, destPath)

			if err := CopyFile(collectionsIconPath, destPath); err != nil {
				logging.LogDebug("Error copying Collections icon: %v", err)
				return fmt.Errorf("error copying Collections icon: %w", err)
			}
		}
	} else if systemName == "Recently Played" {
		// Handle special case for Recently Played
		if rpIconPath, exists := pack.SpecialIcons[RecentlyPlayedIcon]; exists {
			rootMediaPath := filepath.Join(systemPaths.Root, MediaDir)
			if err := os.MkdirAll(rootMediaPath, 0755); err != nil {
				logging.LogDebug("Error creating root media directory: %v", err)
				return fmt.Errorf("error creating root media directory: %w", err)
			}

			destPath := filepath.Join(rootMediaPath, RecentlyPlayedIcon)
			logging.LogDebug("Copying Recently Played icon: %s -> %s", rpIconPath, destPath)

			if err := CopyFile(rpIconPath, destPath); err != nil {
				logging.LogDebug("Error copying Recently Played icon: %v", err)
				return fmt.Errorf("error copying Recently Played icon: %w", err)
			}
		}
	} else if systemName == "Tools" {
		// Handle special case for Tools
		if toolsIconPath, exists := pack.SpecialIcons[ToolsIcon]; exists {
			// Tools icon goes in Tools/.media/tg5040.png
			// Need to go up one level from systemPaths.Tools to get base Tools directory
			toolsBaseDir := filepath.Dir(systemPaths.Tools) // Gets Tools directory without tg5040
			toolsMediaPath := filepath.Join(toolsBaseDir, MediaDir)

			if err := os.MkdirAll(toolsMediaPath, 0755); err != nil {
				logging.LogDebug("Error creating Tools media directory: %v", err)
				return fmt.Errorf("error creating Tools media directory: %w", err)
			}

			destPath := filepath.Join(toolsMediaPath, ToolsIcon)
			logging.LogDebug("Copying Tools icon: %s -> %s", toolsIconPath, destPath)

			if err := CopyFile(toolsIconPath, destPath); err != nil {
				logging.LogDebug("Error copying Tools icon: %v", err)
				return fmt.Errorf("error copying Tools icon: %w", err)
			}
		}
	}

	// Create a map of system tags to icons for faster lookup
	iconsByTag := make(map[string]Icon)
	for _, icon := range pack.Icons {
		if icon.SystemTag != "" {
			iconsByTag[icon.SystemTag] = icon
		}
	}

	// Apply icons to each system or just the specified system
	for _, sysInfo := range systemPaths.Systems {
		// Skip if system has no tag
		if sysInfo.Tag == "" {
			logging.LogDebug("Skipping system with no tag: %s", sysInfo.Name)
			continue
		}

		// If applying to a specific system, check if this is the one
		if systemName != "" && systemName != sysInfo.Name {
			continue
		}

		// Look for an icon with a matching tag
		icon, exists := iconsByTag[sysInfo.Tag]
		if !exists {
			logging.LogDebug("No icon found for system tag: %s", sysInfo.Tag)
			continue
		}

		// Create destination path with renamed icon to match system folder name
		destPath := filepath.Join(mediaPath, sysInfo.Name + ".png")
		logging.LogDebug("Copying icon for system %s: %s -> %s", sysInfo.Name, icon.Path, destPath)

		// Copy and rename the icon
		if err := CopyFile(icon.Path, destPath); err != nil {
			logging.LogDebug("Error copying icon: %v", err)
			return fmt.Errorf("error copying icon for %s: %w", sysInfo.Name, err)
		}
	}

	logging.LogDebug("Icon pack applied successfully")
	return nil
}

// applySpecialIcons applies the special icons (Tools, Recently Played, Collections)
func applySpecialIcons(pack *IconPack, systemPaths *system.SystemPaths) error {
	// Get root media directory for Collections and Recently Played icons
	rootMediaPath := filepath.Join(systemPaths.Root, MediaDir)
	if err := os.MkdirAll(rootMediaPath, 0755); err != nil {
		logging.LogDebug("Error creating root media directory: %v", err)
		return fmt.Errorf("error creating root media directory: %w", err)
	}

	// Apply Tools icon if available
	if toolsIconPath, exists := pack.SpecialIcons[ToolsIcon]; exists {
		// Tools icon goes in Tools/.media/tg5040.png
		// Need to go up one level from systemPaths.Tools to get base Tools directory
		// since systemPaths.Tools is actually Tools/tg5040
		toolsBaseDir := filepath.Dir(systemPaths.Tools) // Gets Tools directory without tg5040
		toolsMediaPath := filepath.Join(toolsBaseDir, MediaDir)

		if err := os.MkdirAll(toolsMediaPath, 0755); err != nil {
			logging.LogDebug("Error creating Tools media directory: %v", err)
			return fmt.Errorf("error creating Tools media directory: %w", err)
		}

		// Copy the icon
		destPath := filepath.Join(toolsMediaPath, ToolsIcon)
		logging.LogDebug("Copying Tools icon: %s -> %s", toolsIconPath, destPath)

		if err := CopyFile(toolsIconPath, destPath); err != nil {
			logging.LogDebug("Error copying Tools icon: %v", err)
			return fmt.Errorf("error copying Tools icon: %w", err)
		}

		logging.LogDebug("Successfully copied Tools icon")
	}

	// Apply Recently Played icon if available
	if rpIconPath, exists := pack.SpecialIcons[RecentlyPlayedIcon]; exists {
		// Recently Played icon goes in root/.media/Recently Played.png
		destPath := filepath.Join(rootMediaPath, RecentlyPlayedIcon)
		logging.LogDebug("Copying Recently Played icon: %s -> %s", rpIconPath, destPath)

		if err := CopyFile(rpIconPath, destPath); err != nil {
			logging.LogDebug("Error copying Recently Played icon: %v", err)
			return fmt.Errorf("error copying Recently Played icon: %w", err)
		}

		logging.LogDebug("Successfully copied Recently Played icon")
	}

	// Apply Collections icon if available
	if collectionsIconPath, exists := pack.SpecialIcons[CollectionsIcon]; exists {
		// Collections icon goes in root/.media/Collections.png
		destPath := filepath.Join(rootMediaPath, CollectionsIcon)
		logging.LogDebug("Copying Collections icon: %s -> %s", collectionsIconPath, destPath)

		if err := CopyFile(collectionsIconPath, destPath); err != nil {
			logging.LogDebug("Error copying Collections icon: %v", err)
			return fmt.Errorf("error copying Collections icon: %w", err)
		}

		logging.LogDebug("Successfully copied Collections icon")
	}

	return nil
}

// CopyFile copies a file from src to dst with additional validation
func CopyFile(src, dst string) error {
	// Add better logging and validation
	logging.LogDebug("Starting file copy: %s -> %s", src, dst)

	// First verify source file exists and is readable
	srcInfo, err := os.Stat(src)
	if err != nil {
		logging.LogDebug("Error checking source file: %v", err)
		return fmt.Errorf("failed to access source file: %w", err)
	}

	// Log warning for suspiciously small files (likely metadata)
	if srcInfo.Size() < 100 {
		logging.LogDebug("WARNING: Source file is suspiciously small (%d bytes), may be a metadata file", srcInfo.Size())
	}

	// Verify source file isn't a directory
	if srcInfo.IsDir() {
		logging.LogDebug("Error: Source is a directory, not a file")
		return fmt.Errorf("source is a directory, not a file")
	}

	// Create the destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		logging.LogDebug("Error creating destination directory: %v", err)
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		logging.LogDebug("Error opening source file: %v", err)
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Try to read a small amount of data to verify file is readable
	testBuf := make([]byte, 10)
	_, err = srcFile.Read(testBuf)
	if err != nil && err != io.EOF {
		logging.LogDebug("Error reading from source file: %v", err)
		return fmt.Errorf("failed to read from source file: %w", err)
	}

	// Reset to beginning of file
	_, err = srcFile.Seek(0, 0)
	if err != nil {
		logging.LogDebug("Error seeking to start of file: %v", err)
		return fmt.Errorf("failed to seek in source file: %w", err)
	}

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		logging.LogDebug("Error creating destination file: %v", err)
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy contents
	bytesCopied, err := io.Copy(dstFile, srcFile)
	if err != nil {
		logging.LogDebug("Error copying file contents: %v", err)
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Verify we copied some data
	if bytesCopied == 0 {
		logging.LogDebug("Warning: Zero bytes copied from %s to %s", src, dst)
	}

	logging.LogDebug("Successfully copied %d bytes from %s to %s", bytesCopied, src, dst)
	return nil
}

// CreatePlaceholderFile creates a placeholder file in the icons directory
func CreatePlaceholderFile() error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create icons directory if it doesn't exist
	iconsDir := filepath.Join(cwd, IconsDir)
	if err := os.MkdirAll(iconsDir, 0755); err != nil {
		logging.LogDebug("Error creating icons directory: %v", err)
		return fmt.Errorf("error creating icons directory: %w", err)
	}

	// Check if directory is empty
	entries, err := os.ReadDir(iconsDir)
	if err != nil {
		logging.LogDebug("Error reading icons directory: %v", err)
		return fmt.Errorf("error reading icons directory: %w", err)
	}

	if len(entries) == 0 {
		// Create a placeholder README file
		placeholderPath := filepath.Join(iconsDir, "README.txt")
		file, err := os.Create(placeholderPath)
		if err != nil {
			logging.LogDebug("Error creating placeholder file: %v", err)
			return fmt.Errorf("error creating placeholder file: %w", err)
		}

		// Write instructions
		_, err = file.WriteString("# Icon Packs\n\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("error writing to placeholder file: %w", err)
		}

		_, err = file.WriteString("Place your icon packs in subdirectories here.\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("error writing to placeholder file: %w", err)
		}

		_, err = file.WriteString("Each icon should be named with the system tag in parentheses, e.g., '(SNES).png'.\n")
		if err != nil {
			file.Close()
			return fmt.Errorf("error writing to placeholder file: %w", err)
		}

		file.Close()
	}

	return nil
}

// DeleteAllIcons removes all icons from the Roms/.media directory
func DeleteAllIcons() error {
	logging.LogDebug("Deleting all system icons")

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Delete regular system icons
	count := 0

	// Path to media directory where system icons are stored
	mediaPath := filepath.Join(systemPaths.Roms, MediaDir)

	// Check if the directory exists
	if _, err := os.Stat(mediaPath); os.IsNotExist(err) {
		logging.LogDebug("Media directory doesn't exist, nothing to delete: %s", mediaPath)
	} else {
		// Read directory contents
		entries, err := os.ReadDir(mediaPath)
		if err != nil {
			logging.LogDebug("Error reading media directory: %v", err)
			return fmt.Errorf("error reading media directory: %w", err)
		}

		// Delete all PNG files (icons) except special ones
		for _, entry := range entries {
			// Skip directories and non-PNG files
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				continue
			}

			// Skip special files that aren't system icons
			if entry.Name() == "bg.png" {
				logging.LogDebug("Skipping background file: %s", entry.Name())
				continue
			}

			// Delete the file
			filePath := filepath.Join(mediaPath, entry.Name())
			logging.LogDebug("Deleting icon: %s", filePath)

			if err := os.Remove(filePath); err != nil {
				logging.LogDebug("Error deleting icon: %v", err)
				// Continue anyway to delete as many as possible
			} else {
				count++
			}
		}
	}

	// Delete special icons too
	if err := DeleteSpecialIcons(systemPaths); err != nil {
		logging.LogDebug("Error deleting special icons: %v", err)
		// Continue anyway
	} else {
		// Count will be incremented in DeleteSpecialIcons
	}

	logging.LogDebug("Successfully deleted %d icons", count)
	return nil
}

// DeleteSpecialIcons removes the special icons (Tools, Recently Played, Collections)
func DeleteSpecialIcons(systemPaths *system.SystemPaths) error {
	count := 0

	// Delete Collections and Recently Played icons from root/.media
	rootMediaPath := filepath.Join(systemPaths.Root, MediaDir)
	if _, err := os.Stat(rootMediaPath); !os.IsNotExist(err) {
		// Delete Collections icon
		collectionsPath := filepath.Join(rootMediaPath, CollectionsIcon)
		if _, err := os.Stat(collectionsPath); !os.IsNotExist(err) {
			logging.LogDebug("Deleting Collections icon: %s", collectionsPath)
			if err := os.Remove(collectionsPath); err != nil {
				logging.LogDebug("Error deleting Collections icon: %v", err)
				// Continue with other deletions
			} else {
				count++
			}
		}

		// Delete Recently Played icon
		rpPath := filepath.Join(rootMediaPath, RecentlyPlayedIcon)
		if _, err := os.Stat(rpPath); !os.IsNotExist(err) {
			logging.LogDebug("Deleting Recently Played icon: %s", rpPath)
			if err := os.Remove(rpPath); err != nil {
				logging.LogDebug("Error deleting Recently Played icon: %v", err)
				// Continue with other deletions
			} else {
				count++
			}
		}
	}

	// Delete Tools icon from Tools/.media
	toolsBaseDir := filepath.Dir(systemPaths.Tools) // Get Tools directory without tg5040
	toolsMediaPath := filepath.Join(toolsBaseDir, MediaDir)
	if _, err := os.Stat(toolsMediaPath); !os.IsNotExist(err) {
		toolsIconPath := filepath.Join(toolsMediaPath, ToolsIcon)
		if _, err := os.Stat(toolsIconPath); !os.IsNotExist(err) {
			logging.LogDebug("Deleting Tools icon: %s", toolsIconPath)
			if err := os.Remove(toolsIconPath); err != nil {
				logging.LogDebug("Error deleting Tools icon: %v", err)
				// Continue with other deletions
			} else {
				count++
			}
		}
	}

	logging.LogDebug("Successfully deleted %d special icons", count)
	return nil
}