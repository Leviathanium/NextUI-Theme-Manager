# Makefile for Theme Manager

# Variables
BINARY_NAME=theme-manager
GO=go
GOBUILD=$(GO) build
GOINSTALL=$(GO) install
GOCLEAN=$(GO) clean
GOGET=$(GO) get

# Main build target
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Clean the binary
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Install dependencies
.PHONY: deps
deps:
	$(GOGET) gopkg.in/yaml.v3

# Build and run
.PHONY: run
run: build
	./$(BINARY_NAME)

# Create a release package
.PHONY: release
release: build
	mkdir -p release/$(BINARY_NAME)
	cp $(BINARY_NAME) release/$(BINARY_NAME)/
	cp minui-list release/$(BINARY_NAME)/
	cp minui-presenter release/$(BINARY_NAME)/
	cp -r Catalog release/$(BINARY_NAME)/
	mkdir -p release/$(BINARY_NAME)/Themes
	mkdir -p release/$(BINARY_NAME)/Overlays
	mkdir -p release/$(BINARY_NAME)/Backups/Themes
	mkdir -p release/$(BINARY_NAME)/Backups/Overlays
	cd release && zip -r $(BINARY_NAME).zip $(BINARY_NAME)

# Utility target to create sample catalog for testing
.PHONY: sample-catalog
sample-catalog:
	mkdir -p Catalog/Themes/RetroWave.theme
	mkdir -p Catalog/Themes/NeonFuture.theme
	mkdir -p Catalog/Themes/Minimal.theme
	mkdir -p Catalog/Overlays/RetroSystem.over
	mkdir -p Catalog/Overlays/ArcadeCabinet.over
	mkdir -p Catalog/Overlays/Transparent.over
	cp example-manifest.yml Catalog/Themes/RetroWave.theme/manifest.yml
	cp example-overlay-manifest.yml Catalog/Overlays/RetroSystem.over/manifest.yml