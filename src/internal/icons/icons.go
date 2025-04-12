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

// IconPack represents an icon pack
type IconPack struct {
	Name  string
	Path  string
	Icons []Icon
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
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
			iconPath := filepath.Join(packPath, entry.Name())

			// Extract system tag if present
			tag := ""
			matches := re.FindStringSubmatch(entry.Name())
			if len(matches) >= 2 {
				tag = matches[1]
			}

			// Add to the list of icons
			pack.Icons = append(pack.Icons, Icon{
				Path:      iconPath,
				Name:      entry.Name(),
				SystemTag: tag,
			})
		}
	}

	logging.LogDebug("Loaded %d icons from pack %s", len(pack.Icons), packName)
	return pack, nil
}

// ApplyIconPack applies the selected icon pack to the system
func ApplyIconPack(packName string) error {
	logging.LogDebug("Applying icon pack: %s", packName)

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

	// Create a map of system tags to icons for faster lookup
	iconsByTag := make(map[string]Icon)
	for _, icon := range pack.Icons {
		if icon.SystemTag != "" {
			iconsByTag[icon.SystemTag] = icon
		}
	}

	// Apply icons to each system
	for _, sysInfo := range systemPaths.Systems {
		// Skip if system has no tag
		if sysInfo.Tag == "" {
			logging.LogDebug("Skipping system with no tag: %s", sysInfo.Name)
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

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

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