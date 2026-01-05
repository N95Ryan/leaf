# 🌱 Leaf

> An elegant, blazingly fast markdown note manager that lives in your terminal.

Leaf is a **Terminal User Interface (TUI) application** built with Go that lets you create, edit, search, and organize markdown notes without leaving your terminal.

## ✨ Features

- **📝 Markdown-Native**: Write notes in pure markdown with real-time rendering
- **⚡ Blazingly Fast**: Lightweight, responsive, built with Go
- **🎨 Beautiful UI**: Intuitive terminal interface with syntax highlighting
- **🔍 Smart Search**: Find notes by content or filename instantly
- **📂 Organized Storage**: Simple file-based structure, easy to version control
- **⌨️ Keyboard-Driven**: Full keyboard navigation—no mouse required

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**
- **gotestsum** (for enhanced test output)
- Windows, macOS, or Linux

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/leaf.git
cd leaf

# Install dependencies
go mod download

# Install gotestsum (if not already installed)
go install gotest.tools/gotestsum@latest

# Run the application
go run ./cmd/leaf
```

## 🧪 Testing

Leaf uses `gotestsum` for enhanced test output:

```powershell
# Run all tests
gotestsum ./...

# Verbose mode
gotestsum --format=verbose ./...

# Storage package tests only
gotestsum ./tests/storage

# Watch mode (auto re-run on changes)
gotestsum --watch ./...

# Generate coverage report
gotestsum -- -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

**Expected output:**

```
✓ tests/storage
∅ internal/app
∅ internal/storage
∅ internal/ui

DONE 5 tests in 0.032s
```

## 🛠️ Development

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Build the executable
go build -o leaf.exe ./cmd/leaf

# Run the application
go run ./cmd/leaf
```

## 📦 Tech Stack

- **Go 1.21+** — Core language
- **Bubbletea** — TUI framework
- **Lipgloss** — Terminal styling
- **Goldmark** — Markdown parsing
- **Chroma** — Syntax highlighting

## 📝 License

MIT License - See [LICENSE](LICENSE) file for details.
