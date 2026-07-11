package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"strings"
)

const (
	ORANGE = "\033[38;5;208m"
	GREEN  = "\033[0;32m"
	CYAN   = "\033[0;36m"
	RESET  = "\033[0m"
)

const (
	VERSION = "1.2"
	AUTHOR  = "PhateValleyman"
	EMAIL   = "Jonas.Ned@outlook.com"
)

type Storage struct {
	Path  string
	Label string
}

func showStorage(path, label string) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return
	}

	total := int64(stat.Blocks) * int64(stat.Bsize)
	free := int64(stat.Bavail) * int64(stat.Bsize)
	used := total - (int64(stat.Bfree) * int64(stat.Bsize))

	if total <= 0 {
		return
	}

	percentUsed := (used * 100) / total
	if percentUsed > 100 {
		percentUsed = 100
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

	usedBar := strings.Repeat("#", usedChars)
	freeBar := strings.Repeat("#", freeChars)

	fmt.Printf("%s%s%s (%s)\n", CYAN, label, RESET, path)
	fmt.Printf("  %s%s%s%s%s  %s%d%% used%s\n",
		ORANGE, usedBar, RESET, GREEN, freeBar, CYAN, percentUsed, RESET)
	fmt.Printf("  Total: %s%s%s | Used: %s%s%s | Free: %s%s%s\n\n",
		CYAN, formatSize(total), RESET,
		ORANGE, formatSize(used), RESET,
		GREEN, formatSize(free), RESET)
}

func detectZyXELUSB() []Storage {
	var usbStorages []Storage
	entries, err := os.ReadDir("/e-data")
	if err != nil {
		return usbStorages
	}

	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			path := "/e-data/" + name
			label := "USB: " + name
			usbStorages = append(usbStorages, Storage{Path: path, Label: label})
		}
	}
	return usbStorages
}

func getSize(path string) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	var size int64
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

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
	
	val := float64(bytes) / float64(div)
	switch exp {
	case 0:
		return fmt.Sprintf("%.1f KB", val)
	case 1:
		return fmt.Sprintf("%.1f MB", val)
	case 2:
		return fmt.Sprintf("%.1f GB", val)
	default:
		return fmt.Sprintf("%.1f TB", val)
	}
}

func printHelp() {
	fmt.Printf("%s=== Storage Utility ===%s\n", CYAN, RESET)
	fmt.Printf("%sUsage:%s\n", ORANGE, RESET)
	fmt.Printf("  %s./storage%s                   Show storage overview\n", CYAN, RESET)
	fmt.Printf("  %s./storage <path1> [path2]%s   Show size of files/directories\n", CYAN, RESET)
	fmt.Printf("  %s./storage -v|--version%s      Show version\n", CYAN, RESET)
	fmt.Printf("  %s./storage -c|--completion%s   Generate bash completion\n", CYAN, RESET)
	fmt.Printf("  %s./storage -h|--help%s         Show help\n", CYAN, RESET)
	fmt.Printf("\n%sExamples:%s\n", ORANGE, RESET)
	fmt.Printf("  %s./storage /etc%s\n", CYAN, RESET)
	fmt.Printf("  %s./storage /etc /var /home%s\n", CYAN, RESET)
}

func printVersion() {
	fmt.Printf("%sstorage v%s%s\n", CYAN, VERSION, RESET)
	fmt.Printf("%sby %s%s\n", ORANGE, AUTHOR, RESET)
	fmt.Printf("%s%s%s\n", GREEN, EMAIL, RESET)
}

func printCompletion() {
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

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help":
			printHelp()
			return
		case "-v", "--version":
			printVersion()
			return
		case "-c", "--completion":
			printCompletion()
			return
		}
		for _, path := range args {
			size, err := getSize(path)
			if err != nil {
				fmt.Printf("%sError for %s: %v%s\n", ORANGE, path, err, RESET)
				continue
			}
			fmt.Printf("%s%s: %s%s\n", CYAN, path, formatSize(size), RESET)
		}
		return
	}

	storages := []Storage{
		{"/storage/emulated/0", "Internal Storage"},
		{"/storage/65D9-1787", "SD Card"},
		{"/storage/sdcard0", "Internal Storage"},
		{"/dev/md0", "HDD 1"},
		{"/dev/md1", "HDD 2"},
		{"/dev/sda1", "Others"},
		{"/", "Root"},
	}

	usbDevices := detectZyXELUSB()
	storages = append(storages, usbDevices...)

	fmt.Printf("%s=== Storage Overview ===%s\n\n", CYAN, RESET)
	
	activeFound := false
	for _, s := range storages {
		// Check if path exists before trying to get stats
		if _, err := os.Stat(s.Path); err == nil {
			showStorage(s.Path, s.Label)
			activeFound = true
		}
	}
	
	if !activeFound {
		fmt.Printf("%sNo active storage paths found.%s\n", ORANGE, RESET)
	}
}
