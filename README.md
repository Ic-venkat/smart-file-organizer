# 📂 Smart File Organizer (SFO)

A powerful, extensible CLI tool written in **Go** to organize your files intelligently. It features a beautiful **interactive TUI** for analyzing your disk usage and a robust rule-based organizer.

## ✨ Features

- **📊 Interactive Scan Mode**: Analyze any directory with a beautiful TUI (Text User Interface).
  - Visualization of file types by count and size.
  - Interactive table with sorting and scrolling.
  - Responsive design that adapts to your terminal window.
- **🚀 Smart Organization**: Automatically moves files into folders based on customizable rules.
- **⚡ Fast**: Built with Go for high performance.
- **🛠 Configurable**: Simple JSON configuration to define your own organization rules.
- **Dry Run**: (Coming Soon) Simulate moves before they happen.
- **Watch Mode**: (Coming Soon) Automatically organize files as they appear.

## 📦 Installation

```bash
# Clone the repository
git clone https://github.com/Ic-venkat/smart-file-organizer.git
cd smart-file-organizer

# Build the binary
go build -o sfo ./cmd/sfo
```

## 🚀 Usage

### 1. Analyze a Directory (Scan)

Before organizing, scan a directory to see what's inside!

```bash
./sfo scan /path/to/directory
```

**Features:**
- Shows Total Files & Total Size.
- Interactive table sorted by file size.
- Shows percentage breakdown.

### 2. Organize Files

Run the organizer to clean up a messy folder.

```bash
./sfo organize /path/to/directory
```

**How it works:**
- Reads your `config.json`.
- Moves files into folders like `Images/`, `Documents/`, etc.
- Handles duplicates by renaming (e.g., `image_123456.jpg`).
- Skips hidden files and directories.

## ⚙️ Configuration

The tool looks for a `config.json` file in the current directory or the executable's directory.

**Default Configuration:**
```json
{
  "Images": ["jpg", "jpeg", "png", "gif", "svg", "webp"],
  "Documents": ["pdf", "doc", "docx", "txt", "md", "xls", "ppt"],
  "Audio": ["mp3", "wav", "aac", "flac"],
  "Video": ["mp4", "mkv", "avi", "mov"],
  "Archives": ["zip", "rar", "tar", "gz"],
  "Code": ["go", "py", "js", "ts", "html", "css", "json", "java"]
}
```

## 🏗 Project Structure

This project follows the **Standard Go Project Layout**:

```
.
├── cmd/
│   └── sfo/          # Main entry point for the CLI
├── internal/
│   ├── organizer/    # Core organization logic
│   ├── scanner/      # Scanning and analytics logic
│   └── config/       # Configuration management (Viper)
└── config.json       # Default configuration rules
```

## 🤝 Contributing

Contributions are welcome! We are working on "Mega" features like:
- **Watch Mode** (Background Daemon)
- **Undo System** (Rollback changes)
- **Content-Aware Sorting** (AI-based)

## 📄 License

MIT License
