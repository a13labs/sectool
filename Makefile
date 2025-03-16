# Makefile for Go project

# Variables
GOCMD = go
GOBUILD = $(GOCMD) build 
GOTEST = $(GOCMD) test
BINARY_NAME = sectool
BUILD_DIR = build

# Targets and Commands
all: clean build-linux build-windows build-mac-intel build-mac-arm test

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-extldflags=-lm" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -v .

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -v .

build-mac-intel:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-extldflags=-lm" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -v .

build-mac-arm:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="-extldflags=-lm" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 -v .

test:
	$(GOTEST) -ldflags="-extldflags=-lm" -v ./internal/...
	$(GOTEST) -ldflags="-extldflags=-lm" -v ./cmd/...

debug:
	$(GOBUILD) -ldflags="-extldflags=-lm" -o $(BUILD_DIR)/$(BINARY_NAME)-debug -gcflags="all=-N -l" -v .

clean:
	rm -rf $(BUILD_DIR)