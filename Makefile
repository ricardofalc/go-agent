# Binary names
BINARY_NAME=fusion-agent
WINDOWS_BINARY=$(BINARY_NAME).exe
LINUX_BINARY=$(BINARY_NAME)

# Version info
VERSION=1.0.0
BUILD_TIME=$(shell date -u +"%Y-%m-%d_%H:%M:%S")

.PHONY: all windows linux clean

all: windows linux

# Build for Windows
windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -ldflags '-X "main.Version=$(VERSION)" -X "main.BuildTime=$(BUILD_TIME)"' -o bin/windows_amd64/$(WINDOWS_BINARY)
	GOOS=windows GOARCH=386 go build -ldflags '-X "main.Version=$(VERSION)" -X "main.BuildTime=$(BUILD_TIME)"' -o bin/windows_386/$(WINDOWS_BINARY)

# Build for Linux
linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -ldflags '-X "main.Version=$(VERSION)" -X "main.BuildTime=$(BUILD_TIME)"' -o bin/linux_amd64/$(LINUX_BINARY)
	GOOS=linux GOARCH=arm64 go build -ldflags '-X "main.Version=$(VERSION)" -X "main.BuildTime=$(BUILD_TIME)"' -o bin/linux_arm64/$(LINUX_BINARY)

# Clean built binaries
clean:
	@echo "Cleaning..."
	rm -rf bin/