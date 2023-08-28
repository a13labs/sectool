# Makefile for Go project

# Variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
BINARY_NAME = sectool
BUILD_DIR = build

# Targets and Commands
all: clean build test

build:
	mkdir -p ${BUILD_DIR}
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v .

test:
	$(GOTEST) -v ./internal/...

clean:
	rm -rf $(BUILD_DIR)
