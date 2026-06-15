// main.go
// Enhanced version that automatically detects connected USB devices on ZyXEL
// All original code preserved as much as possible.

package main

import (
	"fmt"
	"os"
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

// getSize returns size in bytes for a file or directory (recursively)
func getSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if !info.IsDir() {
		return info.Size(), nil
	}
	// For directories, use du -sb (fast and reliable)
	cmd := exec.Command("du", "-sb", path)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(out))
	if len(fields) < 1 {
		return 0, fmt.Errorf("du output parse error")
	}
	size, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}

// formatSize converts bytes to human-readable string (B, KB, MB, GB, TB)
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	switch exp {
	case 0:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(div))
	case 1:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(div))
	case 2:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(div))
	default:
		return fmt.Sprintf("%.1f TB", float64(bytes)/float64(div))
	}
}

// printHelp displays a colorful usage screen with examples
func printHelp() {
	fmt.Printf("%s=== Storage Utility ===%s\n", CYAN, RESET)
	fmt.Printf("%sUsage:%s\n", ORANGE, RESET)
	fmt.Printf("  %s./storage%s                   %sShow overview of all storages (disks, USB, etc.)%s\n", CYAN, RESET, GREEN, RESET)
	fmt.Printf("  %s./storage <path>%s            %sShow size of a file or directory (recursive)%s\n", CYAN, RESET, GREEN, RESET)
	fmt.Printf("  %s./storage -h | --help%s        %sShow this help screen%s\n", CYAN, RESET, GREEN, RESET)
	fmt.Printf("\n%sExamples:%s\n", ORANGE, RESET)
	fmt.Printf("  %s./storage /etc/passwd%s        %sShow size of the passwd file%s\n", CYAN, RESET, GREEN, RESET)
	fmt.Printf("  %s./storage /e-data%s            %sShow total size of the /e-data directory%s\n", CYAN, RESET, GREEN, RESET)
	fmt.Printf("  %s./storage /dev/md0%s           %sShow size of a device (if it is a filesystem)%s\n", CYAN, RESET, GREEN, RESET)
	fmt.Printf("\n%sNote:%s Directory size is computed using %sdu -sb%s, may take a while.\n", ORANGE, RESET, CYAN, RESET)
}

func main() {
	args := os.Args[1:]
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		printHelp()
		return
	}
	if len(args) == 1 {
		path := args[0]
		size, err := getSize(path)
		if err != nil {
			fmt.Printf("%sError: %s%s\n", ORANGE, err, RESET)
			return
		}
		fmt.Printf("%sSize of %s%s:%s %s\n", CYAN, GREEN, path, RESET, formatSize(size))
		return
	}

	// Original storage display
	storages := []Storage{
		// Android / Redmi
		{"/storage/emulated/0", "Internal Storage"},
		{"/storage/65D9-1787", "SD Card"},
		// Android / tablet
		{"/storage/sdcard0", "Internal Storage"},
		// ZyXEL server
		{"/dev/md0", "HDD 1"},
		{"/dev/md1", "HDD 2"},
		// Others
		{"/dev/sda1", "Others"},
	}

	// Append auto-detected ZyXEL USB devices
	usbDevices := detectZyXELUSB()
	storages = append(storages, usbDevices...)

	fmt.Printf("%s=== Storage ===%s\n", CYAN, RESET)
	for _, s := range storages {
		showStorage(s.Path, s.Label)
	}
}
