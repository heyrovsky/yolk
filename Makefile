# Paths
BUILD_CONFIG := $(shell pwd)/config/build-config.json
CONFIG_PACKAGE := github.com/heyrovsky/yolk/common/config
OUTPUT_BINARY := bin/yolk

# Extract Git hash
GIT_HASH := $(shell git describe --tags --always --dirty | sed 's/^v//')



BUILD_FLAGS        := -trimpath -buildvcs=false -ldflags "\
	-w -s \
	-X $(CONFIG_PACKAGE).EXTERNAL_VERSION=$(VERSION) \
	-X $(CONFIG_PACKAGE).EXTERNAL_APP_NAME=$(BINARY_NAME)"
PLATFORMS          := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64
