// src/internal/themes/sync.go
// Implementation of theme repository syncing functionality

package themes

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nextui-themes/internal/logging"
	"nextui-themes/internal/ui"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// RepoConfig holds repository configuration information
var RepoConfig struct {
	URL    string
	Branch string
}

// Initialize default repository settings
func init() {
	// Public GitHub repos require the ".git" suffix for anonymous go-git clones/pulls
	RepoConfig.URL = "https://github.com/Leviathanium/NextUI-Themes.git"
	RepoConfig.Branch = "main"
}

// CatalogData represents the structure of the catalog.json file
type CatalogData struct {
	LastUpdated string                                `json:"last_updated"`
	Themes      map[string]CatalogItemInfo            `json:"themes"`
	Components  map[string]map[string]CatalogItemInfo `json:"components"`
}

// CatalogItemInfo represents an item in the catalog
type CatalogItemInfo struct {
	PreviewPath  string `json:"preview_path"`
	ManifestPath string `json:"manifest_path"`
	Author       string `json:"author"`
	Description  string `json:"description"`
	URL          string `json:"URL"`  // Added URL field for ZIP download
}

// SyncOptions contains options for syncing
type SyncOptions struct {
	RepoURL      string
	Branch       string
	LocalDirPath string
	UI           bool // Whether to show UI progress messages
}

// GetDefaultSyncOptions returns default sync options
func GetDefaultSyncOptions() SyncOptions {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		cwd = "."
	}

	return SyncOptions{
		RepoURL:      RepoConfig.URL,
		Branch:       RepoConfig.Branch,
		LocalDirPath: cwd,
		UI:           true,
	}
}

// SetRepoURL updates the repository URL configuration
func SetRepoURL(url string) {
	RepoConfig.URL = url
	logging.LogDebug("Repository URL set to: %s", url)
}

// SetRepoBranch updates the repository branch configuration
func SetRepoBranch(branch string) {
	RepoConfig.Branch = branch
	logging.LogDebug("Repository branch set to: %s", branch)
}

// SyncThemeCatalog syncs the theme catalog from the repository
func SyncThemeCatalog(options SyncOptions) error {
	logging.LogDebug("Starting theme catalog sync from %s", options.RepoURL)

	if options.UI {
		ui.ShowMessage("Syncing theme catalog...", "1")
	}

	// Create directory structure if it doesn't exist
	err := createSyncDirectoryStructure(options.LocalDirPath)
	if err != nil {
		return fmt.Errorf("error creating directory structure: %w", err)
	}

	// First, try to use the HTTP method which is more efficient for this use case
	if err := syncCatalogViaHTTP(options); err != nil {
		logging.LogDebug("HTTP sync failed, falling back to Git: %v", err)

		// If HTTP method fails, fall back to Git
		if err := syncCatalogViaGit(options); err != nil {
			return fmt.Errorf("git sync failed: %w", err)
		}
	}

	if options.UI {
		ui.ShowMessage("Theme catalog sync completed successfully!", "2")
	}

	return nil
}

// createSyncDirectoryStructure creates the necessary directory structure for syncing
func createSyncDirectoryStructure(basePath string) error {
	// Create Catalog directory if it doesn't exist
	catalogDir := filepath.Join(basePath, "Catalog")
	if err := os.MkdirAll(catalogDir, 0755); err != nil {
		return fmt.Errorf("error creating Catalog directory: %w", err)
	}

	// Create Themes directory and subdirectories
	themesDir := filepath.Join(catalogDir, "Themes")
	if err := os.MkdirAll(filepath.Join(themesDir, "previews"), 0755); err != nil {
		return fmt.Errorf("error creating Themes/previews directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(themesDir, "manifests"), 0755); err != nil {
		return fmt.Errorf("error creating Themes/manifests directory: %w", err)
	}

	// Create Components directory and subdirectories
	componentsDir := filepath.Join(catalogDir, "Components")

	// Component types
	componentTypes := []string{"Wallpapers", "Icons", "Accents", "LEDs", "Fonts", "Overlays"}

	// Create directories for each component type
	for _, compDirName := range componentTypes {
		compDir := filepath.Join(componentsDir, compDirName)
		if err := os.MkdirAll(filepath.Join(compDir, "previews"), 0755); err != nil {
			return fmt.Errorf("error creating %s/previews directory: %w", compDirName, err)
		}
		if err := os.MkdirAll(filepath.Join(compDir, "manifests"), 0755); err != nil {
			return fmt.Errorf("error creating %s/manifests directory: %w", compDirName, err)
		}
	}

	// Create cache directory for temporary files
	cacheDir := filepath.Join(basePath, ".cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("error creating .cache directory: %w", err)
	}

	return nil
}

// syncCatalogViaHTTP downloads catalog data using HTTP(S) requests - more efficient for small files
func syncCatalogViaHTTP(options SyncOptions) error {
	// Base URL for raw content
	baseURL := options.RepoURL
	if strings.Contains(baseURL, "github.com") {
		// Strip .git so raw.githubusercontent path matches your repo layout
		baseURL = strings.TrimSuffix(baseURL, ".git")
		// Swap domains
		baseURL = strings.Replace(baseURL,
			"github.com", "raw.githubusercontent.com", 1)
		// Manually append the branch (don't use filepath.Join on URLs!)
		baseURL = fmt.Sprintf("%s/%s", baseURL, options.Branch)
	}

	// First download catalog.json
	catalogURL := fmt.Sprintf("%s/Catalog/catalog.json", baseURL)
	logging.LogDebug("Downloading catalog.json from %s", catalogURL)

	// Download the catalog file
	localCatalogPath := filepath.Join(options.LocalDirPath, "Catalog", "catalog.json")
	if err := downloadFile(catalogURL, localCatalogPath); err != nil {
		return fmt.Errorf("error downloading catalog.json: %w", err)
	}

	// Parse the catalog to get list of preview and manifest files
	catalog, err := parseCatalogJSON(localCatalogPath)
	if err != nil {
		return fmt.Errorf("error parsing catalog.json: %w", err)
	}

	// Download all previews and manifests for themes
	for _, themeInfo := range catalog.Themes {
		// Download preview
		if themeInfo.PreviewPath != "" {
			previewURL := fmt.Sprintf("%s/%s", baseURL, themeInfo.PreviewPath)
			localPreviewPath := filepath.Join(options.LocalDirPath, themeInfo.PreviewPath)
			if err := downloadFile(previewURL, localPreviewPath); err != nil {
				logging.LogDebug("Warning: Error downloading preview %s: %v", previewURL, err)
			}
		}

		// Download manifest
		if themeInfo.ManifestPath != "" {
			manifestURL := fmt.Sprintf("%s/%s", baseURL, themeInfo.ManifestPath)
			localManifestPath := filepath.Join(options.LocalDirPath, themeInfo.ManifestPath)
			if err := downloadFile(manifestURL, localManifestPath); err != nil {
				logging.LogDebug("Warning: Error downloading manifest %s: %v", manifestURL, err)
			}
		}
	}

	// Download all previews and manifests for components
	for compType, items := range catalog.Components {
		logging.LogDebug("Processing component type: %s", compType)
		for _, compInfo := range items {
			// Download preview
			if compInfo.PreviewPath != "" {
				previewURL := fmt.Sprintf("%s/%s", baseURL, compInfo.PreviewPath)
				localPreviewPath := filepath.Join(options.LocalDirPath, compInfo.PreviewPath)
				if err := downloadFile(previewURL, localPreviewPath); err != nil {
					logging.LogDebug("Warning: Error downloading preview %s: %v", previewURL, err)
				}
			}

			// Download manifest
			if compInfo.ManifestPath != "" {
				manifestURL := fmt.Sprintf("%s/%s", baseURL, compInfo.ManifestPath)
				localManifestPath := filepath.Join(options.LocalDirPath, compInfo.ManifestPath)
				if err := downloadFile(manifestURL, localManifestPath); err != nil {
					logging.LogDebug("Warning: Error downloading manifest %s: %v", manifestURL, err)
				}
			}
		}
	}

	return nil
}

// Modified downloadFile function for src/internal/themes/sync.go
// Increases timeout and adds better error handling

func downloadFile(url string, localPath string) error {
    // Create the directory structure for the file
    dir := filepath.Dir(localPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("error creating directory %s: %w", dir, err)
    }

    // Create HTTP client with extended timeout (5 minutes instead of 30 seconds)
    client := &http.Client{
        Timeout: 5 * time.Minute,
    }

    // Download the file
    resp, err := client.Get(url)
    if err != nil {
        // Provide more specific error message for timeout
        if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline") {
            return fmt.Errorf("download timed out - try again or download manually: %w", err)
        }
        return fmt.Errorf("download error: %w", err)
    }
    defer resp.Body.Close()

    // Check status code
    if resp.StatusCode != 200 {
        return fmt.Errorf("HTTP error: %s", resp.Status)
    }

    // Create the local file
    out, err := os.Create(localPath)
    if err != nil {
        return fmt.Errorf("error creating local file: %w", err)
    }
    defer out.Close()

    // Copy the content
    _, err = io.Copy(out, resp.Body)

    if err != nil {
        // Clean up partial downloads on error
        out.Close()
        os.Remove(localPath)
        return fmt.Errorf("error during download: %w", err)
    }

    return nil
}

// parseCatalogJSON parses the catalog.json file
func parseCatalogJSON(path string) (*CatalogData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var catalog CatalogData
	err = json.Unmarshal(data, &catalog)
	if err != nil {
		return nil, err
	}

	return &catalog, nil
}

// syncCatalogViaGit syncs the theme catalog using Git
func syncCatalogViaGit(options SyncOptions) error {
	logging.LogDebug("Syncing theme catalog via Git from %s", options.RepoURL)

	repoPath := filepath.Join(options.LocalDirPath, ".git")
	var r *git.Repository

	// Check if repo already exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		// Clone the repository
		r, _ = git.PlainClone(options.LocalDirPath, false, &git.CloneOptions{
			URL:           options.RepoURL,
			Progress:      nil,
			ReferenceName: plumbing.NewBranchReferenceName(options.Branch),
			SingleBranch:  true,
			NoCheckout:    false,
			Depth:         1,
		})
	} else {
		// Open existing repository
		r, err = git.PlainOpen(options.LocalDirPath)
		if err != nil {
			return fmt.Errorf("error opening repository: %w", err)
		}

		// Pull latest changes
		w, err := r.Worktree()
		if err != nil {
			return fmt.Errorf("error getting worktree: %w", err)
		}

		err = w.Pull(&git.PullOptions{
			RemoteName:    "origin",
			ReferenceName: plumbing.NewBranchReferenceName(options.Branch),
			SingleBranch:  true,
			Force:         true,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("error pulling repository: %w", err)
		}
	}

	logging.LogDebug("Git sync completed successfully")
	return nil
}

// DownloadThemePackage downloads a specific theme package from the repository
func DownloadThemePackage(themeName string) error {
	logging.LogDebug("Downloading theme package: %s", themeName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Check if the theme already exists locally
	localThemePath := filepath.Join(cwd, "Themes", themeName)
	if _, err := os.Stat(localThemePath); err == nil {
		logging.LogDebug("Theme '%s' already exists locally, skipping download", themeName)
		return nil
	}

	// Path to catalog.json
	catalogPath := filepath.Join(cwd, "Catalog", "catalog.json")

	// Parse the catalog.json file
	catalog, err := parseCatalogJSON(catalogPath)
	if err != nil {
		return fmt.Errorf("error parsing catalog.json: %w", err)
	}

	// Check if theme exists in catalog
	themeInfo, exists := catalog.Themes[themeName]
	if !exists {
		return fmt.Errorf("theme '%s' not found in catalog", themeName)
	}

	// Ensure the URL field exists
	if themeInfo.URL == "" {
		return fmt.Errorf("no download URL found for theme '%s'", themeName)
	}

	// Show status message
	ui.ShowMessage(fmt.Sprintf("Downloading theme '%s'...", themeName), "1")

	// Create cache directory
	cacheDir := filepath.Join(cwd, ".cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("error creating cache directory: %w", err)
	}

	// Create a temporary file for the ZIP
	zipPath := filepath.Join(cacheDir, fmt.Sprintf("%s.zip", themeName))

	// Download the ZIP file
	if err := downloadFile(themeInfo.URL, zipPath); err != nil {
		return fmt.Errorf("error downloading theme ZIP: %w", err)
	}

	// Ensure Themes directory exists
	themesDir := filepath.Join(cwd, "Themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		return fmt.Errorf("error creating Themes directory: %w", err)
	}

	// Extract the ZIP file
	if err := extractZipFile(zipPath, localThemePath); err != nil {
		return fmt.Errorf("error extracting theme ZIP: %w", err)
	}

	// Clean up the ZIP file
	if err := os.Remove(zipPath); err != nil {
		logging.LogDebug("Warning: Failed to remove temporary ZIP file: %v", err)
	}

	ui.ShowMessage(fmt.Sprintf("Theme '%s' downloaded successfully!", themeName), "2")
	return nil
}

// DownloadComponentPackage downloads a specific component package from the repository
func DownloadComponentPackage(componentType, componentName string) error {
	logging.LogDebug("Downloading component package: %s - %s", componentType, componentName)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Check if the component already exists locally
	localComponentPath := filepath.Join(cwd, "Components", componentType, componentName)
	if _, err := os.Stat(localComponentPath); err == nil {
		logging.LogDebug("Component '%s' already exists locally, skipping download", componentName)
		return nil
	}

	// Path to catalog.json
	catalogPath := filepath.Join(cwd, "Catalog", "catalog.json")

	// Parse the catalog.json file
	catalog, err := parseCatalogJSON(catalogPath)
	if err != nil {
		return fmt.Errorf("error parsing catalog.json: %w", err)
	}

	// Map component type to catalog key
	componentTypeMap := map[string]string{
		"Wallpapers": "wallpapers",
		"Icons":      "icons",
		"Accents":    "accents",
		"LEDs":       "leds",
		"Fonts":      "fonts",
		"Overlays":   "overlays",
	}

	catalogType := componentTypeMap[componentType]
	if catalogType == "" {
		return fmt.Errorf("unknown component type: %s", componentType)
	}

	// Check if component exists in catalog
	components, exists := catalog.Components[catalogType]
	if !exists {
		return fmt.Errorf("component type '%s' not found in catalog", componentType)
	}

	componentInfo, exists := components[componentName]
	if !exists {
		return fmt.Errorf("component '%s' not found in catalog", componentName)
	}

	// Ensure the URL field exists
	if componentInfo.URL == "" {
		return fmt.Errorf("no download URL found for component '%s'", componentName)
	}

	// Show status message
	ui.ShowMessage(fmt.Sprintf("Downloading %s component '%s'...", componentType, componentName), "1")

	// Create cache directory
	cacheDir := filepath.Join(cwd, ".cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("error creating cache directory: %w", err)
	}

	// Create a temporary file for the ZIP
	zipPath := filepath.Join(cacheDir, fmt.Sprintf("%s.zip", componentName))

	// Download the ZIP file
	if err := downloadFile(componentInfo.URL, zipPath); err != nil {
		return fmt.Errorf("error downloading component ZIP: %w", err)
	}

	// Ensure Components directory exists
	componentsDir := filepath.Join(cwd, "Components", componentType)
	if err := os.MkdirAll(componentsDir, 0755); err != nil {
		return fmt.Errorf("error creating Components directory: %w", err)
	}

	// Extract the ZIP file
	if err := extractZipFile(zipPath, localComponentPath); err != nil {
		return fmt.Errorf("error extracting component ZIP: %w", err)
	}

	// Clean up the ZIP file
	if err := os.Remove(zipPath); err != nil {
		logging.LogDebug("Warning: Failed to remove temporary ZIP file: %v", err)
	}

	ui.ShowMessage(fmt.Sprintf("%s component '%s' downloaded successfully!", componentType, componentName), "2")
	return nil
}

// extractZipFile extracts a ZIP archive to the specified destination directory
func extractZipFile(zipPath, destDir string) error {
	logging.LogDebug("Extracting ZIP file %s to %s", zipPath, destDir)

	// Open the ZIP file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("error opening ZIP file: %w", err)
	}
	defer reader.Close()

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("error creating destination directory: %w", err)
	}

	// Analyze the ZIP structure to detect common root directories
	// This helps prevent issues like Theme/Theme nesting
	rootDirs := make(map[string]int)
	totalFiles := 0

	for _, file := range reader.File {
		// Skip __MACOSX directories and hidden files
		if strings.Contains(file.Name, "__MACOSX") || strings.HasPrefix(filepath.Base(file.Name), ".") {
			continue
		}

		totalFiles++

		// Get the top-level directory from the path
		pathParts := strings.Split(file.Name, "/")
		if len(pathParts) > 1 && pathParts[0] != "" {
			rootDirs[pathParts[0]]++
		}
	}

	// Check if all files are in a single root directory
	var commonRoot string
	destBaseName := filepath.Base(destDir)

	// Find if there's a single common root directory that contains all files
	for dir, count := range rootDirs {
		if count == totalFiles || (float64(count)/float64(totalFiles) > 0.9) {
			commonRoot = dir
			break
		}
	}

	logging.LogDebug("ZIP analysis - Total files: %d, Common root: %s, Dest dir: %s",
		totalFiles, commonRoot, destBaseName)

	// Extract each file in the ZIP archive
	for _, file := range reader.File {
		// Skip __MACOSX directories and hidden files
		if strings.Contains(file.Name, "__MACOSX") || strings.HasPrefix(filepath.Base(file.Name), ".") {
			continue
		}

		// Determine the target path for extraction
		var targetPath string

		// If there's a common root that matches the destination directory name or ends with the same extension,
		// strip it to avoid Theme/Theme nesting
		if commonRoot != "" && (commonRoot == destBaseName ||
			(strings.HasSuffix(commonRoot, filepath.Ext(destBaseName)) &&
			 strings.HasSuffix(destBaseName, filepath.Ext(destBaseName)))) {

			if strings.HasPrefix(file.Name, commonRoot+"/") {
				// Strip the common root to avoid nesting
				relativePath := strings.TrimPrefix(file.Name, commonRoot+"/")

				// Special case: if we have an entry for just the directory itself (resulting in empty path)
				// Skip this entry as we've already created the destination directory
				if relativePath == "" {
					logging.LogDebug("Skipping root directory entry: %s", file.Name)
					continue
				}

				targetPath = filepath.Join(destDir, relativePath)
				logging.LogDebug("Stripping common root from: %s to: %s", file.Name, relativePath)
			} else {
				// Normal file, not in common root
				targetPath = filepath.Join(destDir, file.Name)
				logging.LogDebug("File doesn't have common root prefix: %s", file.Name)
			}
		} else {
			// No common root or it doesn't match destination - extract normally
			targetPath = filepath.Join(destDir, file.Name)
			logging.LogDebug("Normal extraction for: %s", file.Name)
		}

		// Check for directory traversal attacks - only for non-empty paths
		// The destDir path itself is always safe since we create it explicitly
		cleanPath := filepath.Clean(targetPath)
		if cleanPath != destDir && !strings.HasPrefix(cleanPath, destDir+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", file.Name)
		}

		// Handle directories
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %w", targetPath, err)
			}
			continue
		}

		// Create the directory structure for the file
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("error creating directory structure: %w", err)
		}

		// Open the file in the ZIP
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening file in ZIP: %w", err)
		}

		// Create the destination file
		outFile, err := os.Create(targetPath)
		if err != nil {
			rc.Close()
			return fmt.Errorf("error creating file %s: %w", targetPath, err)
		}

		// Copy the content
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return fmt.Errorf("error extracting file %s: %w", targetPath, err)
		}
	}

	logging.LogDebug("ZIP extraction completed successfully")
	return nil
}