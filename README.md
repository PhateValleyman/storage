# Storage Monitor (Go)

A lightweight **Go** program that displays storage usage across multiple devices — including **Android**, **tablets**, and **ZyXEL NAS** systems — with automatic detection of **USB drives** connected to ZyXEL NAS (`/e-data`).

This project preserves the original code structure while extending it with dynamic USB detection on ZyXEL devices running **FFP (fonz_fun_plug)**.

---

## 📦 Features

- Displays formatted storage information for:
  - **Redmi** internal and SD card paths
  - **Android tablet** storage
  - **ZyXEL NAS** RAID volumes (`/dev/md0`, `/dev/md1`)
  - **Automatically detected USB drives** in `/e-data`
- Uses `df` to gather disk stats (works on Linux and embedded NAS systems)
- Shows color-coded bar charts with usage percentage
- Designed to run on **low-resource environments** (ARMv5 ZyXEL NAS)
- ANSI color support for clean, readable terminal output

---

## 🧠 How It Works

1. A static list of known storage paths is defined (Android + ZyXEL HDDs).
2. The function `detectZyXELUSB()` lists the `/e-data` directory to find all connected USB devices automatically.
3. Each found mount is appended to the main storage list.
4. For each storage path, `df` is executed to retrieve capacity and usage data.
5. The information is displayed as a colored usage bar and numerical summary.

---

## ⚙️ Requirements

- **Go 1.19+**  
- Works on:
  - Linux (x86 / ARM / ARMv5 / ARMv7)
  - Android (via Termux)
  - ZyXEL NAS with FFP (fonz_fun_plug)

Ensure that:
- `/e-data` exists (for ZyXEL USB devices)
- `df` and `ls` commands are available in your system PATH

---

## 🚀 Usage

### Run directly:
```bash
go run main.go

Build binary:

go build -o storage main.go

Run the binary:

./storage


---

📁 Directory Detection Logic (ZyXEL USB)

The ZyXEL USB devices are mounted under /e-data/ with directories named after the device UUID or label.
For example:

/e-data/af5af4bc-aed3-4497-9663-9e2c60bbd5cb
/e-data/68E0-DADB

The program dynamically lists these entries and adds them to the monitored storages list, labeling them as USB: <directory>.


---

🧩 Future Enhancements

Auto-detect USB labels via blkid

Shorten UUID display (e.g., USB: af5af)

Include network mounts (/mnt/nfs, /mnt/smb)

JSON output mode for integration with dashboards



---

🧑‍💻 Author

Jonáš Nedvědický (PhateValleyman)
Maintained as part of embedded utilities for ARM-based systems.


---



