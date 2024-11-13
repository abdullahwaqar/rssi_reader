GO := go

OUTPUT_DIR := bin

SOURCES := $(wildcard *.go)

BINARY_NAME := rssi_reader

# ARM compilation settings
GOARCH := arm
GOARM := 7
GOOS := linux

# Default target
all: arm

# Create output directory
$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

# Compile for ARM
arm: $(OUTPUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)_armv7 $(SOURCES)

# Clean build artifacts
clean:
	rm -rf $(OUTPUT_DIR)

.PHONY: all arm clean
