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
	VERSION = "1.4"
	AUTHOR  = "PhateValleyman"
	EMAIL   = "Jonas.Ned@outlook.com"
)

type Storage struct {
	Path  string
	Label string
}

func getAvailablePaths() []Storage {
	storages := []Storage{
		{"/i-data/3776680e", "HDD 1"},
		{"/i-data/43d6b8c3", "HDD 2"},
		{"/storage/emulated/0", "Internal Storage"},
		{"/storage/65D9-1787", "SD Card"},
		{"/storage/sdcard0", "Internal Storage"},
		{"/dev/md0", "HDD 1"},
		{"/dev/md1", "HDD 2"},
		{"/dev/sda1", "Others"},
		{"/", "Root"},
	}
	storages = append(storages, detectZyXELUSB()...)
	storages = append(storages, detectEnvironments()...)

	var active []Storage
	seen := make(map[string]bool)
	for _, s := range storages {
		if seen[s.Path] {
			continue
		}
		if _, err := os.Stat(s.Path); err == nil {
			active = append(active, s)
			seen[s.Path] = true
		}
	}
	return active
}

func detectEnvironments() []Storage {
	var envStorages []Storage
	if os.Getenv("CLOUD_SHELL") == "true" {
		envStorages = append(envStorages, Storage{Path: "/home", Label: "Cloud Shell Home"})
	}
	if os.Getenv("CODESPACES") == "true" {
		envStorages = append(envStorages, Storage{Path: "/workspaces", Label: "Codespace Workspace"})
	}
	return envStorages
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
	fmt.Printf("%s=== Storage Utility v%s ===%s\n", CYAN, VERSION, RESET)
	fmt.Printf("%sUsage:%s\n", ORANGE, RESET)
	fmt.Printf("  %s./storage%s                   Show storage overview (HDD, USB, Cloud, Root)\n", CYAN, RESET)
	fmt.Printf("  %s./storage <path1> [path2]%s   Show size of specific files/directories\n", CYAN, RESET)
	fmt.Printf("  %s./storage -v|--version%s      Show version and author info\n", CYAN, RESET)
	fmt.Printf("  %s./storage -c|--completion%s   Generate bash completion script\n", CYAN, RESET)
	fmt.Printf("  %s./storage -h|--help%s         Show this help screen\n", CYAN, RESET)
	fmt.Printf("\n%sFeatures:%s\n", ORANGE, RESET)
	fmt.Printf("  - Auto-detection of HDD 1/2 (/dev/md0, /dev/md1)\n")
	fmt.Printf("  - Auto-detection of USB drives in /e-data\n")
	fmt.Printf("  - Cloud environment awareness (Cloud Shell, Codespaces)\n")
	fmt.Printf("  - Dynamic bash completion for active drives\n")
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
	binaryName := "storage"
	if len(os.Args) > 0 {
		binaryName = filepath.Base(os.Args[0])
	}
	fmt.Printf(`# bash completion for %s
_%s_completion() {
	local cur prev opts paths
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
	opts="-h --help -v --version -c --completion"

	# Get dynamic paths from the binary itself
	paths=$( "${COMP_WORDS[0]}" --list-paths 2>/dev/null )

	if [[ ${cur} == -* ]] ; then
		COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
		return 0
	fi

	# Suggest known paths and standard files/dirs
	COMPREPLY=( $(compgen -W "${paths}" -- "${cur}") )
	COMPREPLY+=( $(compgen -f -- "${cur}") )

	return 0
}
complete -o filenames -F _%s_completion %s
`, binaryName, binaryName, binaryName, binaryName)
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
		case "--list-paths":
			active := getAvailablePaths()
			for i, s := range active {
				fmt.Print(s.Path)
				if i < len(active)-1 {
					fmt.Print(" ")
				}
			}
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

	active := getAvailablePaths()

	fmt.Printf("%s=== Storage Overview ===%s\n\n", CYAN, RESET)
	
	if len(active) == 0 {
		fmt.Printf("%sNo active storage paths found.%s\n", ORANGE, RESET)
		return
	}

	for _, s := range active {
		showStorage(s.Path, s.Label)
	}
}
