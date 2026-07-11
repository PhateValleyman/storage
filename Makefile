# Makefile for Storage Monitor project
# Supports building for ZyXEL NAS, Redmi (Android / Termux), and native Linux.

# Colors
YELLOW := \033[1;33m
CYAN := \033[0;36m
GREEN := \033[0;32m
RED := \033[0;31m
RESET := \033[0m

# Binary name and output folder
BINARY := storage
BIN_DIR := ./bin

# Go build flags
# Using GOTOOLCHAIN to ensure compatibility with older ARMv5 if needed
GOTOOLCHAIN := go1.19.3
BUILD_FLAGS := -trimpath -ldflags="-s -w"

# Default target
.DEFAULT_GOAL := help

# Architecture definitions
ZyXEL_GOARCH := arm
ZyXEL_GOARM := 5
ZyXEL_GOOS := linux

Redmi_GOARCH := arm64
Redmi_GOARM := 8
Redmi_GOOS := android

Termux_GOARCH := arm64
Termux_GOARM := 8
Termux_GOOS := linux

Native_GOARCH := $(shell go env GOARCH)
Native_GOOS := $(shell go env GOOS)

# Ensure bin directory exists
$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

# Help target
help:
	@echo ""
	@echo "$(CYAN)Storage Monitor Build System v1.4$(RESET)"
	@echo ""
	@echo "$(YELLOW)Usage:$(RESET)"
	@echo "  make [target]"
	@echo ""
	@echo "$(YELLOW)Targets:$(RESET)"
	@echo "  $(GREEN)help$(RESET)       - Show this help message"
	@echo "  $(GREEN)zyxel$(RESET)      - Build binary for ZyXEL NAS (ARMv5, statically linked)"
	@echo "  $(GREEN)redmi$(RESET)      - Build binary for Redmi (Android/Termux ARM64)"
	@echo "  $(GREEN)native$(RESET)     - Build binary for native host system ($(Native_GOOS)/$(Native_GOARCH))"
	@echo "  $(GREEN)install$(RESET)    - Auto-detect system and install appropriate binary"
	@echo "  $(GREEN)clean$(RESET)      - Remove build artifacts"
	@echo ""

# ZyXEL target
zyxel: $(BIN_DIR)
	@echo "$(CYAN)[*] Building for ZyXEL NAS (ARMv5 soft-float)$(RESET)"
	GOTOOLCHAIN=$(GOTOOLCHAIN) CGO_ENABLED=0 GOOS=$(ZyXEL_GOOS) GOARCH=$(ZyXEL_GOARCH) GOARM=$(ZyXEL_GOARM) go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(BINARY)-zyxel main.go
	@echo "$(GREEN)✔ Build complete: $(BIN_DIR)/$(BINARY)-zyxel$(RESET)"

# Redmi target
redmi: $(BIN_DIR)
	@echo "$(CYAN)[*] Building for Redmi (Android / ARM64)$(RESET)"
	GOOS=$(Redmi_GOOS) GOARCH=$(Redmi_GOARCH) GOARM=$(Redmi_GOARM) go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(BINARY)-redmi main.go
	@echo "$(GREEN)✔ Build complete: $(BIN_DIR)/$(BINARY)-redmi$(RESET)"

# Native build
native: $(BIN_DIR)
	@echo "$(CYAN)[*] Building for native system ($(Native_GOOS)/$(Native_GOARCH))$(RESET)"
	GOOS=$(Native_GOOS) GOARCH=$(Native_GOARCH) go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(BINARY)-native main.go
	@echo "$(GREEN)✔ Build complete: $(BIN_DIR)/$(BINARY)-native$(RESET)"

# Install target (auto-detect platform)
install:
	@echo "$(CYAN)[*] Detecting target system...$(RESET)"
	@if [ -d "/ffp" ]; then \
		echo "$(CYAN)[+] Detected ZyXEL NAS (FFP environment)$(RESET)"; \
		$(MAKE) zyxel >/dev/null; \
		cp $(BIN_DIR)/$(BINARY)-zyxel /ffp/bin/$(BINARY); \
		chmod +x /ffp/bin/$(BINARY); \
		echo "$(GREEN)✔ Installed to /ffp/bin/$(BINARY)$(RESET)"; \
	elif [ -d "/data/data/com.termux/files/usr" ]; then \
		echo "$(CYAN)[+] Detected Termux/Android environment$(RESET)"; \
		GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(BINARY)-termux main.go; \
		cp $(BIN_DIR)/$(BINARY)-termux /data/data/com.termux/files/usr/bin/$(BINARY); \
		chmod +x /data/data/com.termux/files/usr/bin/$(BINARY); \
		echo "$(GREEN)✔ Installed to /data/data/com.termux/files/usr/bin/$(BINARY)$(RESET)"; \
	else \
		echo "$(CYAN)[+] Detected native Linux system$(RESET)"; \
		$(MAKE) native >/dev/null; \
		sudo cp $(BIN_DIR)/$(BINARY)-native /usr/local/bin/$(BINARY); \
		echo "$(GREEN)✔ Installed to /usr/local/bin/$(BINARY)$(RESET)"; \
	fi

# Clean target
clean:
	@echo "$(CYAN)[*] Cleaning build artifacts...$(RESET)"
	rm -rf $(BIN_DIR)
	@echo "$(GREEN)✔ Clean complete.$(RESET)"
