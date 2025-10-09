// main.go
// Enhanced version that automatically detects connected USB devices on ZyXEL
// All original code preserved as much as possible.

package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ANSI colors
const (
	ORANGE = "\033[38;5;208m"
	GREEN  = "\033[0;32m"
	CYAN   = "\033[0;36m"
	RESET  = "\033[0m"
)

type Storage struct {
	Path  string
	Label string
}

// Function to display usage bar for given storage path
func showStorage(path, label string) {
	cmd := exec.Command("df", path)
	out, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return
	}

	total, _ := strconv.ParseInt(fields[1], 10, 64)
	used, _ := strconv.ParseInt(fields[2], 10, 64)
	free, _ := strconv.ParseInt(fields[3], 10, 64)

	if total <= 0 {
		return
	}

	totalGB := total / 1024 / 1024
	usedGB := used / 1024 / 1024
	freeMB := free / 1024
	percentUsed := used * 100 / total

	if percentUsed > 100 {
		percentUsed = 100
	} else if percentUsed < 0 {
		percentUsed = 0
	}

	barWidth := 30
	usedChars := int(int64(barWidth) * percentUsed / 100)
	if usedChars < 1 && percentUsed > 0 {
		usedChars = 1
	}
	if usedChars > barWidth {
		usedChars = barWidth
	}
	freeChars := barWidth - usedChars
	if freeChars < 0 {
		freeChars = 0
	}

	usedBar := strings.Repeat("#", usedChars)
	freeBar := strings.Repeat("#", freeChars)

	fmt.Printf("%s%s%s\n", CYAN, label, RESET)
	fmt.Printf("  %s%s%s%s%s  %s%d%% used%s\n",
		ORANGE, usedBar, RESET, GREEN, freeBar, CYAN, percentUsed, RESET)
	fmt.Printf("  Total: %s%d G%s | Used: %s%d G%s | Free: %s%d M%s\n\n",
		CYAN, totalGB, RESET,
		ORANGE, usedGB, RESET,
		GREEN, freeMB, RESET)
}

// Function to detect connected ZyXEL USB devices dynamically
func detectZyXELUSB() []Storage {
	var usbStorages []Storage

	// Try to list /e-data directories (typical for FFP ZyXEL USB mounts)
	cmd := exec.Command("ls", "/e-data")
	out, err := cmd.Output()
	if err != nil {
		return usbStorages
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		path := "/e-data/" + line
		label := "USB: " + line
		usbStorages = append(usbStorages, Storage{Path: path, Label: label})
	}

	return usbStorages
}

func main() {
	storages := []Storage{
		// Android / Redmi
		{"/storage/emulated/0", "Internal Storage"},
		{"/storage/65D9-1787", "SD Card"},
		// Android / tablet
		{"/storage/sdcard0", "Internal Storage"},
		// ZyXEL server
		{"/dev/md0", "HDD 1"},
		{"/dev/md1", "HDD 2"},
	}

	// Append auto-detected ZyXEL USB devices
	usbDevices := detectZyXELUSB()
	storages = append(storages, usbDevices...)

	fmt.Printf("%s=== Storage ===%s\n", CYAN, RESET)
	for _, s := range storages {
		showStorage(s.Path, s.Label)
	}
}
