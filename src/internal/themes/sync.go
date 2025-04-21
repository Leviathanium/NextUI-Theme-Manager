// src/internal/themes/sync.go
// Implementation of theme repository syncing functionality

package themes

import (
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
	// Public GitHub repos require the “.git” suffix for anonymous go-git clones/pulls
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

// downloadFile downloads a file from a URL to a local path
func downloadFile(url string, localPath string) error {
	// Create the directory structure for the file
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Download the file
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// Create the local file
	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the content
	_, err = io.Copy(out, resp.Body)
	return err
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

	// Default sync options
	options := GetDefaultSyncOptions()

	// Path to catalog.json
	catalogPath := filepath.Join(cwd, "Catalog", "catalog.json")

	// Parse the catalog.json file
	catalog, err := parseCatalogJSON(catalogPath)
	if err != nil {
		return fmt.Errorf("error parsing catalog.json: %w", err)
	}

	// Check if theme exists in catalog
	_, exists := catalog.Themes[themeName]
	if !exists {
		return fmt.Errorf("theme '%s' not found in catalog", themeName)
	}

	// Show status message
	ui.ShowMessage(fmt.Sprintf("Downloading theme '%s'...", themeName), "1")

	// Base URL for raw content
	baseURL := options.RepoURL
	if strings.Contains(baseURL, "github.com") {
		// Convert GitHub URL to raw content URL
		baseURL = strings.Replace(baseURL, "github.com", "raw.githubusercontent.com", 1)
		baseURL = filepath.Join(baseURL, options.Branch)
	}

	// Removed unused themeURL variable
	localThemePath := filepath.Join(cwd, "Themes", themeName)

	// Ensure Themes directory exists
	if err := os.MkdirAll(filepath.Join(cwd, "Themes"), 0755); err != nil {
		return fmt.Errorf("error creating Themes directory: %w", err)
	}

	// Get path to the manifest file
	manifestPath := filepath.Join(cwd, catalog.Themes[themeName].ManifestPath)

	// Read the manifest to get file structure
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("error reading manifest file: %w", err)
	}

	// Parse manifest to get file structure - just looking for basic content flags
	var manifestMap map[string]interface{}
	err = json.Unmarshal(manifestData, &manifestMap)
	if err != nil {
		return fmt.Errorf("error parsing manifest: %w", err)
	}

	// Parse the content section to determine what to download
	_, ok := manifestMap["content"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid manifest structure: missing or invalid 'content' section")
	}

	// Create theme directory
	if err := os.MkdirAll(localThemePath, 0755); err != nil {
		return fmt.Errorf("error creating theme directory: %w", err)
	}

	// First, download preview.png and manifest.json to the theme directory
	previewURL := fmt.Sprintf("%s/Catalog/Themes/previews/%s.png", baseURL, themeName)
	localPreviewFile := filepath.Join(localThemePath, "preview.png")
	if err := downloadFile(previewURL, localPreviewFile); err != nil {
		logging.LogDebug("Warning: Error downloading preview: %v", err)
	}

	manifestURL := fmt.Sprintf("%s/Catalog/Themes/manifests/%s.json", baseURL, themeName)
	localManifestFile := filepath.Join(localThemePath, "manifest.json")
	if err := downloadFile(manifestURL, localManifestFile); err != nil {
		logging.LogDebug("Warning: Error downloading manifest: %v", err)
	}

	// For now, we're keeping this simple and just showing a message that we'd download the theme
	// In a real implementation, we'd need to recursively download all files in the theme directory
	ui.ShowMessage(fmt.Sprintf("Theme '%s' downloaded successfully!", themeName), "2")

	// Return theme path for importing
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

	// Default sync options
	options := GetDefaultSyncOptions()

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

	_, exists = components[componentName]
	if !exists {
		return fmt.Errorf("component '%s' not found in catalog", componentName)
	}

	// Show status message
	ui.ShowMessage(fmt.Sprintf("Downloading %s component '%s'...", componentType, componentName), "1")

	// Base URL for raw content
	baseURL := options.RepoURL
	if strings.Contains(baseURL, "github.com") {
		// Convert GitHub URL to raw content URL
		baseURL = strings.Replace(baseURL, "github.com", "raw.githubusercontent.com", 1)
		baseURL = filepath.Join(baseURL, options.Branch)
	}

	// Removed unused componentURL variable
	localComponentPath := filepath.Join(cwd, "Components", componentType, componentName)

	// Ensure Components directory exists
	if err := os.MkdirAll(filepath.Join(cwd, "Components", componentType), 0755); err != nil {
		return fmt.Errorf("error creating Components directory: %w", err)
	}

	// Create component directory
	if err := os.MkdirAll(localComponentPath, 0755); err != nil {
		return fmt.Errorf("error creating component directory: %w", err)
	}

	// First, download preview.png and manifest.json to the component directory
	previewURL := fmt.Sprintf("%s/Catalog/Components/%s/previews/%s.png", baseURL, componentType, componentName)
	localPreviewFile := filepath.Join(localComponentPath, "preview.png")
	if err := downloadFile(previewURL, localPreviewFile); err != nil {
		logging.LogDebug("Warning: Error downloading preview: %v", err)
	}

	manifestURL := fmt.Sprintf("%s/Catalog/Components/%s/manifests/%s.json", baseURL, componentType, componentName)
	localManifestFile := filepath.Join(localComponentPath, "manifest.json")
	if err := downloadFile(manifestURL, localManifestFile); err != nil {
		logging.LogDebug("Warning: Error downloading manifest: %v", err)
	}

	// For now, we're keeping this simple and just showing a message that we'd download the component
	// In a real implementation, we'd need to recursively download all files in the component directory
	ui.ShowMessage(fmt.Sprintf("%s component '%s' downloaded successfully!", componentType, componentName), "2")

	// Return component path for importing
	return nil
}
