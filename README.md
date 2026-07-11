# Storage Monitor (Go)

A lightweight **Go** program that displays storage usage across multiple devices — including **Android**, **tablets**, **ZyXEL NAS** systems, and **Cloud Environments** (Google Cloud Shell, GitHub Codespaces) — with automatic detection of **USB drives** connected to ZyXEL NAS (`/e-data`).

This project preserves the original code structure while extending it with dynamic USB detection and cloud awareness.

---

## 📦 Features

- **v1.4 Highlights:**
  - **Dynamic Bash Completion:** Real-time suggestions for active drives and system paths.
  - **Cloud Awareness:** Auto-detects Google Cloud Shell and GitHub Codespaces storage.
  - **HDD Detection:** Specifically tracks ZyXEL HDD 1 and HDD 2 via `/i-data/` mount points.
- Displays formatted storage information for:
  - **Redmi** internal and SD card paths
  - **Android tablet** storage
  - **ZyXEL NAS** RAID volumes (`/dev/md0`, `/dev/md1`)
  - **Automatically detected USB drives** in `/e-data`
- Uses native `syscall.Statfs` to gather disk stats (efficient and portable).
- Shows color-coded bar charts with usage percentage.
- Designed to run on **low-resource environments** (ARMv5 ZyXEL NAS).
- ANSI color support for clean, readable terminal output.

---

## 🧠 How It Works

1. A list of known storage paths is defined (Android + ZyXEL HDDs + Cloud Home).
2. The function `detectZyXELUSB()` lists the `/e-data` directory to find all connected USB devices automatically.
3. Environment variables (`CLOUD_SHELL`, `CODESPACES`) are checked to add cloud-specific workspaces.
4. Active paths are filtered for existence and uniqueness.
5. For each storage path, disk statistics (Total, Used, Free) are retrieved.
6. The information is displayed as a colored usage bar and numerical summary.

---

## ⚙️ Requirements

- **Go 1.17+** (Builds tested with **Go 1.19.3**)
- Works on:
  - Linux (x86 / ARM / ARMv5 / ARMv7)
  - Android (via Termux)
  - ZyXEL NAS with FFP (fonz_fun_plug)
  - Google Cloud Shell / GitHub Codespaces

---

## 🚀 Usage

### Run directly:
```bash
go run main.go
```

### Build binary:
```bash
go build -o storage main.go
```

### Shell Completion:
Generate and source the completion script to get TAB-suggestions for your disks:
```bash
source <(./storage -c)
```

---

## 🛠 Build System (Makefile)

The included `Makefile` supports cross-compilation for various targets:

- `make zyxel` - Build for ARMv5 (ZyXEL NAS)
- `make redmi` - Build for Android ARM64
- `make native` - Build for current host
- `make install` - Detect environment and install to bin path

---

## 🧑‍💻 Author

Jonáš Nedvědický (PhateValleyman)
Maintained as part of embedded utilities for ARM-based systems.
