# Makefile for cross-building storage binary
# Valleyman style, server/Redmi/native targets with help

# ---------------------------
# Colors for help output
# ---------------------------
RED    := \033[0;31m
GREEN  := \033[0;32m
YELLOW := \033[0;33m
BLUE   := \033[0;34m
NC     := \033[0m

# ---------------------------
# Default target
# ---------------------------
TARGET ?= native
BIN_DIR := ./bin
BIN     := $(BIN_DIR)/storage
SRC     := ./main.go

# ---------------------------
# Toolchain for NSA320 (server)
# ---------------------------
SERVER_CC  := /ffp/bin/arm-ffp-linux-uclibcgnueabi-gcc
SERVER_CXX := /ffp/bin/arm-ffp-linux-uclibcgnueabi-g++
SERVER_CFLAGS := --sysroot=/ffp
SERVER_LDFLAGS := --sysroot=/ffp

# ---------------------------
# Toolchain for Redmi (ARM64 Android)
# ---------------------------
REDMI_CC  := aarch64-linux-android21-clang
REDMI_CXX := aarch64-linux-android21-clang++
REDMI_CFLAGS := ""
REDMI_LDFLAGS := ""

# ---------------------------
# Help target
# ---------------------------
.PHONY: help
help:
	@echo "$(GREEN)Usage: make [TARGET]$(NC)"
	@echo "Targets:"
	@echo "  $(YELLOW)server$(NC)   - Build for NSA320 (ARMv5, Linux 2.6.31.8, uClibc)"
	@echo "  $(YELLOW)redmi$(NC)    - Build for Redmi (ARM64 Android)"
	@echo "  $(YELLOW)native$(NC)   - Build for current OS/architecture"
	@echo "  $(YELLOW)clean$(NC)    - Remove bin directory"

# ---------------------------
# Build targets
# ---------------------------
.PHONY: all
all: $(TARGET)

.PHONY: server
server: $(BIN)
	@echo "$(BLUE)Built for server (NSA320 ARMv5) -> $(BIN)$(NC)"

.PHONY: redmi
redmi: $(BIN)
	@echo "$(BLUE)Built for Redmi (ARM64 Android) -> $(BIN)$(NC)"

.PHONY: native
native: $(BIN)
	@echo "$(BLUE)Built for native OS/arch -> $(BIN)$(NC)"

$(BIN): $(SRC)
	@mkdir -p $(BIN_DIR)
ifeq ($(TARGET),server)
	@echo "$(YELLOW)Cross-building for NSA320 (ARMv5, uClibc)...$(NC)"
	@CGO_ENABLED=1 \
	GOOS=linux GOARCH=arm GOARM=5 \
	CC=$(SERVER_CC) CXX=$(SERVER_CXX) \
	CGO_CFLAGS="$(SERVER_CFLAGS)" CGO_LDFLAGS="$(SERVER_LDFLAGS)" \
	go build -v -o $(BIN) $(SRC)
endif
ifeq ($(TARGET),redmi)
	@echo "$(YELLOW)Cross-building for Redmi (ARM64 Android)...$(NC)"
	@CGO_ENABLED=1 \
	GOOS=android GOARCH=arm64 \
	CC=$(REDMI_CC) CXX=$(REDMI_CXX) \
	CGO_CFLAGS="$(REDMI_CFLAGS)" CGO_LDFLAGS="$(REDMI_LDFLAGS)" \
	go build -v -o $(BIN) $(SRC)
endif
ifeq ($(TARGET),native)
	@echo "$(YELLOW)Building for current OS/arch...$(NC)"
	@CGO_ENABLED=0 go build -v -o $(BIN) $(SRC)
endif

# ---------------------------
# Clean target
# ---------------------------
.PHONY: clean
clean:
	@echo "$(RED)Removing $(BIN_DIR)...$(NC)"
	@rm -rf $(BIN_DIR)
