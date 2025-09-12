PAK_NAME := $(shell jq -r .name pak.json)
PAK_TYPE := $(shell jq -r .type pak.json)
PAK_FOLDER := $(shell echo $(PAK_TYPE) | cut -c1)$(shell echo $(PAK_TYPE) | tr '[:upper:]' '[:lower:]' | cut -c2-)s

clean:
	rm -rf dist src/theme-manager

build:
	cd src && env GOOS=linux GOARCH=arm64 go build -o theme-manager cmd/theme-manager/main.go

release: build
	mkdir -p "dist/$(PAK_NAME).pak"
	$(MAKE) bump-version
	cp -R src/theme-manager resources/launch.sh README.md LICENSE pak.json resources/minui-list resources/minui-presenter "dist/$(PAK_NAME).pak"
	cd "dist/$(PAK_NAME).pak" && zip -r "../$(PAK_NAME).pak.zip" "."
	zip -r "dist/$(PAK_NAME).pak.zip" pak.json
	ls -lah dist

bump-version:
ifeq ($(RELEASE_VERSION),)
	error "RELEASE_VERSION is not set"
endif
	jq '.version = "$(RELEASE_VERSION)"' pak.json > pak.json.tmp
	mv pak.json.tmp pak.json
