# Storage Usage Monitor

A lightweight Go utility that displays storage usage information with a visual progress bar and color-coded output for easy monitoring.

## Features

- 📊 Visual storage usage representation with progress bars
- 🎨 Color-coded output for better readability
- 📱 Supports multiple storage devices (internal, external, RAID arrays)
- 🔢 Displays total, used, and free space in appropriate units
- 🖥️ Clean, formatted output for terminal viewing

## Supported Storage Types

- Android internal storage
- SD cards
- RAID arrays (md devices)
- Any mountable storage device

## Installation

1. Ensure you have Go installed on your system
2. Clone or download this repository
3. Build the executable:
   ```bash
   CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -trimpath -v -x -o ./bin/storage ./main.go
```

Usage

Run the compiled binary:

```bash
./storage
```

Output Example

```
=== Úložiště ===
Interní úložiště
  ##############################  85% used
  Total: 64 G | Used: 54 G | Free: 10240 M

SD karta
  ##################           60% used
  Total: 128 G | Used: 77 G | Free: 51200 M
```

Customization

To monitor different storage paths, edit the storages slice in the main() function:

```go
storages := []Storage{
    {"/path/to/your/storage", "Custom Label"},
    {"/another/path", "Another Label"},
}
```

Requirements

· Go 1.11 or higher
· Linux/Unix-like system with df command available
· Terminal that supports ANSI color codes

