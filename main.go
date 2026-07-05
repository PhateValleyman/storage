// main.go
// Enhanced version that automatically detects connected USB devices on ZyXEL
// All original code preserved as much as possible.

package main

// Import required packages
import (
	// Standard format package
	"fmt"
	// Standard OS package
	"os"
	// Package for executing external commands
	"os/exec"
	// Package for string conversions
	"strconv"
	// Package for string manipulation
	"strings"
)

// ANSI colors constants
const (
	// Orange color for warnings or used space
	ORANGE = "\033[38;5;208m"
	// Green color for free space
	GREEN = "\033[0;32m"
	// Cyan color for labels and headers
	CYAN = "\033[0;36m"
	// Reset color to default
	RESET = "\033[0m"
)

// Version metadata constants
const (
	// Current version of the utility
	VERSION = "1.1"
	// Author of the utility
	AUTHOR = "PhateValleyman"
	// Email address for contact
	EMAIL = "Jonas.Ned@outlook.com"
)

// Storage struct defines path and label for storage devices
type Storage struct {
	// Filesystem path for the storage
	Path string
	// Display label for the storage
	Label string
}

// Function to display usage bar for given storage path
func showStorage(path, label string) {
	// Create command to run df on specified path
	cmd := exec.Command("df", path)
	// Execute command and capture output
	out, err := cmd.Output()
	// Check for errors during execution
	if err != nil {
		// Exit function if error occurs
		return
	}

	// Split output into lines
	lines := strings.Split(string(out), "\n")
	// Check if output has enough lines
	if len(lines) < 2 {
		// Exit function if output is too short
		return
	}

	// Extract fields from the second line of df output
	fields := strings.Fields(lines[1])
	// Check if there are enough fields
	if len(fields) < 4 {
		// Exit function if fields are missing
		return
	}

	// Parse total blocks from the second field
	total, _ := strconv.ParseInt(fields[1], 10, 64)
	// Parse used blocks from the third field
	used, _ := strconv.ParseInt(fields[2], 10, 64)
	// Parse free blocks from the fourth field
	free, _ := strconv.ParseInt(fields[3], 10, 64)

	// Avoid division by zero if total is zero
	if total <= 0 {
		// Exit function if total is zero
		return
	}

	// Calculate total space in Gigabytes
	totalGB := total / 1024 / 1024
	// Calculate used space in Gigabytes
	usedGB := used / 1024 / 1024
	// Calculate free space in Megabytes
	freeMB := free / 1024

	// Calculate percentage of used space
	percentUsed := used * 100 / total
	// Clamp percentage to maximum 100
	if percentUsed > 100 {
		// Set to 100 if exceeded
		percentUsed = 100
	// Clamp percentage to minimum 0
	} else if percentUsed < 0 {
		// Set to 0 if negative
		percentUsed = 0
	}

	// Set width of the usage bar
	barWidth := 30
	// Calculate number of characters for used space
	usedChars := int(int64(barWidth) * percentUsed / 100)
	// Ensure at least one character if used space is greater than zero
	if usedChars < 1 && percentUsed > 0 {
		// Set to 1
		usedChars = 1
	}
	// Clamp used characters to bar width
	if usedChars > barWidth {
		// Set to barWidth
		usedChars = barWidth
	}

	// Calculate number of characters for free space
	freeChars := barWidth - usedChars
	// Ensure free characters is not negative
	if freeChars < 0 {
		// Set to 0
		freeChars = 0
	}

	// Create string for used part of the bar
	usedBar := strings.Repeat("#", usedChars)
	// Create string for free part of the bar
	freeBar := strings.Repeat("#", freeChars)

	// Print storage label with cyan color
	fmt.Printf("%s%s%s\n", CYAN, label, RESET)
	// Print the usage bar with orange and green colors
	fmt.Printf("  %s%s%s%s%s  %s%d%% used%s\n",
		ORANGE, usedBar, RESET, GREEN, freeBar, CYAN, percentUsed, RESET)
	// Print detailed space information
	fmt.Printf("  Total: %s%d G%s | Used: %s%d G%s | Free: %s%d M%s\n\n",
		CYAN, totalGB, RESET,
		ORANGE, usedGB, RESET,
		GREEN, freeMB, RESET)
}

// Function to detect connected ZyXEL USB devices dynamically
func detectZyXELUSB() []Storage {
	// Initialize slice to hold storage information
	var usbStorages []Storage

	// Create command to list directories in /e-data
	cmd := exec.Command("ls", "/e-data")
	// Execute command and capture output
	out, err := cmd.Output()
	// Check for errors during execution
	if err != nil {
		// Return empty slice on error
		return usbStorages
	}

	// Split output into lines and trim whitespace
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// Iterate through each line in output
	for _, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			// Continue to next line
			continue
		}
		// Construct full path for the USB device
		path := "/e-data/" + line
		// Construct label for the USB device
		label := "USB: " + line
		// Append storage info to the slice
		usbStorages = append(usbStorages, Storage{Path: path, Label: label})
	}

	// Return detected USB storages
	return usbStorages
}

// getSize returns size in bytes for a file or directory (recursively)
func getSize(path string) (int64, error) {
	// Get information about the file or directory
	info, err := os.Stat(path)
	// Check for errors getting file info
	if err != nil {
		// Return zero and error if stat fails
		return 0, err
	}
	// Check if the path is not a directory
	if !info.IsDir() {
		// Return file size directly
		return info.Size(), nil
	}

	// Create command to run du for directory size in bytes
	cmd := exec.Command("du", "-sb", path)
	// Execute command and capture output
	out, err := cmd.Output()
	// Check for errors during execution
	if err != nil {
		// Return zero and error if du fails
		return 0, err
	}

	// Split output into fields
	fields := strings.Fields(string(out))
	// Check if output has enough fields
	if len(fields) < 1 {
		// Return error if output cannot be parsed
		return 0, fmt.Errorf("du output parse error")
	}

	// Parse size from the first field of du output
	size, err := strconv.ParseInt(fields[0], 10, 64)
	// Check for parsing errors
	if err != nil {
		// Return zero and error if parsing fails
		return 0, err
	}
	// Return the calculated size
	return size, nil
}

// formatSize converts bytes to human-readable string
func formatSize(bytes int64) string {
	// Define base unit for calculations
	const unit = 1024

	// Return bytes directly if less than one kilobyte
	if bytes < unit {
		// Return formatted string with B unit
		return fmt.Sprintf("%d B", bytes)
	}

	// Initialize divisor and exponent for unit conversion
	div, exp := int64(unit), 0
	// Loop to determine appropriate unit scale
	for n := bytes / unit; n >= unit; n /= unit {
		// Increase divisor by unit factor
		div *= unit
		// Increment exponent
		exp++
	}

	// Switch based on exponent to choose unit
	switch exp {
	// Case for kilobytes
	case 0:
		// Return size in KB
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(div))
	// Case for megabytes
	case 1:
		// Return size in MB
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(div))
	// Case for gigabytes
	case 2:
		// Return size in GB
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(div))
	// Default case for terabytes
	default:
		// Return size in TB
		return fmt.Sprintf("%.1f TB", float64(bytes)/float64(div))
	}
}

// printHelp displays usage screen
func printHelp() {
	// Print header with cyan color
	fmt.Printf("%s=== Storage Utility ===%s\n", CYAN, RESET)
	// Print usage label with orange color
	fmt.Printf("%sUsage:%s\n", ORANGE, RESET)

	// Print basic usage for overview
	fmt.Printf("  %s./storage%s                   Show storage overview\n", CYAN, RESET)
	// Print usage for multiple paths
	fmt.Printf("  %s./storage <path1> [path2]%s   Show size of files/directories\n", CYAN, RESET)
	// Print usage for version flag
	fmt.Printf("  %s./storage -v|--version%s      Show version\n", CYAN, RESET)
	// Print usage for completion flag
	fmt.Printf("  %s./storage -c|--completion%s   Generate bash completion\n", CYAN, RESET)
	// Print usage for help flag
	fmt.Printf("  %s./storage -h|--help%s         Show help\n", CYAN, RESET)

	// Print examples label with orange color
	fmt.Printf("\n%sExamples:%s\n", ORANGE, RESET)
	// Print example for single path
	fmt.Printf("  %s./storage /etc%s\n", CYAN, RESET)
	// Print example for multiple paths
	fmt.Printf("  %s./storage /etc /var /home%s\n", CYAN, RESET)

	// Print technical note about du usage
	fmt.Printf("\n%sNote:%s du -sb is used for directory sizes\n", ORANGE, RESET)
}

// printVersion shows colored version info
func printVersion() {
	// Print version information
	fmt.Printf("%sstorage v%s%s\n", CYAN, VERSION, RESET)
	// Print author information
	fmt.Printf("%sby %s%s\n", ORANGE, AUTHOR, RESET)
	// Print email information
	fmt.Printf("%s%s%s\n", GREEN, EMAIL, RESET)
}

// printCompletion generates bash completion script
func printCompletion() {
	// Print the bash completion script
	fmt.Printf(`# bash completion for storage

_storage_completion() {
	local cur prev
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"

	local opts="-h --help -v --version -c --completion"

	if [[ ${COMP_CWORD} -eq 1 ]]; then
		COMPREPLY=( $(compgen -W "${opts}" -f -- "$cur") )
		return
	fi

	COMPREPLY=( $(compgen -f -- "$cur") )
}

complete -F _storage_completion storage
`)
}

// main entry point of the application
func main() {
	// Capture command line arguments excluding program name
	args := os.Args[1:]

	// Check if there are any arguments provided
	if len(args) > 0 {
		// Evaluate the first argument for flags
		switch args[0] {
		// Match help flags
		case "-h", "--help":
			// Display help screen
			printHelp()
			// Terminate program
			return
		// Match version flags
		case "-v", "--version":
			// Display version info
			printVersion()
			// Terminate program
			return
		// Match completion flags
		case "-c", "--completion":
			// Display completion script
			printCompletion()
			// Terminate program
			return
		}

		// Iterate through all provided path arguments
		for _, path := range args {
			// Get size for the current path
			size, err := getSize(path)
			// Check if size calculation failed
			if err != nil {
				// Print error message
				fmt.Printf("%sError for %s: %s%s\n", ORANGE, path, err, RESET)
				// Continue to next path
				continue
			}
			// Print the path and its formatted size
			fmt.Printf("%s%s: %s%s\n", CYAN, path, formatSize(size), RESET)
		}
		// Terminate program after processing paths
		return
	}

	// Initialize default storage paths and labels
	storages := []Storage{
		// Path for internal storage on some systems
		{"/storage/emulated/0", "Internal Storage"},
		// Path for SD card
		{"/storage/65D9-1787", "SD Card"},
		// Alternative path for internal storage
		{"/storage/sdcard0", "Internal Storage"},
		// Path for first HDD device
		{"/dev/md0", "HDD 1"},
		// Path for second HDD device
		{"/dev/md1", "HDD 2"},
		// Path for other partitions
		{"/dev/sda1", "Others"},
	}

	// Dynamically detect USB devices on ZyXEL
	usbDevices := detectZyXELUSB()
	// Merge detected USB devices into the storages list
	storages = append(storages, usbDevices...)

	// Print main storage header
	fmt.Printf("%s=== Storage ===%s\n", CYAN, RESET)
	// Iterate through all storage entries
	for _, s := range storages {
		// Display storage info and usage bar
		showStorage(s.Path, s.Label)
	}
}
