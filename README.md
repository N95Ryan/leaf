# Leaf 🌴

> An elegant, blazingly fast markdown note manager that lives in your terminal.

Leaf is a **Terminal User Interface (TUI) application** built with Go that lets you create, edit, search, and organize markdown notes without leaving your terminal.

## Features

- **Markdown-Native**: Write notes in pure markdown; rendered in read mode with syntax highlighting
- **Blazingly Fast**: Lightweight, responsive, built with Go
- **Beautiful UI**: Styled terminal interface with Lipgloss
- **Smart Search**: Find notes by title or content
- **Organized Storage**: File-based structure in `~/.leaf/notes/`, easy to version control
- **Keyboard-Driven**: Full keyboard navigation — no mouse required
- **Cross-Platform**: Linux, macOS, and Windows (including Git Bash)

## Prerequisites

- **Go 1.25+** ([download](https://go.dev/dl/)) — required to build from source
- A **modern terminal** with ANSI color support:
  - **Windows**: [Windows Terminal](https://aka.ms/terminal) + Git Bash (recommended), or Git Bash (Mintty) alone
  - **macOS**: Terminal.app, iTerm2, or Warp
  - **Linux**: any modern terminal emulator

Go 1.21+ can auto-download the required toolchain when you run `go` commands.

## Quick Start

### Option 1 — Clone and run (Git Bash / Linux / macOS)

```bash
git clone https://github.com/N95Ryan/leaf.git
cd leaf
go mod download
go run ./cmd/leaf
```

### Option 2 — Install script (Git Bash / Linux / macOS)

```bash
git clone https://github.com/N95Ryan/leaf.git
cd leaf
bash scripts/install.sh
leaf
```

### Option 3 — go install (any platform with Go)

```bash
go install github.com/N95Ryan/leaf/cmd/leaf@latest
leaf
```

### Option 4 — Pre-built binary (no Go required)

Download the latest release for your platform from [GitHub Releases](https://github.com/N95Ryan/leaf/releases):

| Platform              | Binary                   |
| --------------------- | ------------------------ |
| Linux (x86_64)        | `leaf-linux-amd64`       |
| Linux (ARM64)         | `leaf-linux-arm64`       |
| Windows (x86_64)      | `leaf-windows-amd64.exe` |
| macOS (Intel)         | `leaf-darwin-amd64`      |
| macOS (Apple Silicon) | `leaf-darwin-arm64`      |

**Git Bash (Windows):**

```bash
chmod +x leaf-windows-amd64.exe
./leaf-windows-amd64.exe
```

**Linux / macOS:**

```bash
chmod +x leaf-linux-amd64   # or leaf-darwin-arm64, etc.
./leaf-linux-amd64
```

Releases are published automatically when a version tag is pushed (e.g. `v0.1.0`).

## Build from source

```bash
go build -o leaf ./cmd/leaf
```

On Windows, Go produces `leaf.exe`. Run with `./leaf` (Git Bash) or `.\leaf.exe` (PowerShell).

## Configuration

### Custom notes directory

By default, notes are stored in `~/.leaf/notes/`. Override with a flag:

```bash
go run ./cmd/leaf --notes-dir ./my-notes
leaf --notes-dir "$HOME/Documents/leaf-notes"
```

Useful for testing, portable setups, or syncing a specific folder.

## Keyboard Shortcuts

| Key       | Mode          | Action                          |
| --------- | ------------- | ------------------------------- |
| `n`       | List          | Create a new note               |
| `r`       | List / Search | Read selected note              |
| `e`       | List          | Edit selected note              |
| `/`       | List          | Open search                     |
| `Enter`   | Search        | Run search                      |
| `j` / `k` | List / Search | Navigate up/down                |
| `t`       | List          | Cycle sort mode                 |
| `d`       | List          | Delete (press twice to confirm) |
| `Tab`     | Edit          | Switch title/content focus      |
| `Ctrl+S`  | Create / Edit | Save note                       |
| `Esc`     | Any           | Cancel / back to list           |
| `q`       | Any           | Quit                            |

## Troubleshooting

### `go: command not found`

Install Go 1.25+ from [go.dev/dl](https://go.dev/dl/) and restart your terminal.

### TUI looks broken (no colors, garbled layout)

Use a modern terminal. On Windows, prefer **Windows Terminal** with the Git Bash profile instead of legacy CMD.

### `Ctrl+S` does not save (Git Bash)

Some terminals intercept `Ctrl+S` as flow control (XOFF). If the terminal freezes, press `Ctrl+Q` to resume. Prefer **Windows Terminal** for reliable `Ctrl+S` handling, or save from a different terminal.

### Emoji / special characters display as boxes

Change your terminal font to one with emoji support (e.g. Cascadia Code, JetBrains Mono, Nerd Fonts).

### Notes location

Notes are plain Markdown files at `~/.leaf/notes/{uuid}.md`. You can edit them with any text editor; restart Leaf or re-open the app to refresh the list.

## Testing

```bash
go test ./...

# With gotestsum (optional)
gotestsum ./...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

CI runs on **Linux, Windows, and macOS** on every push and pull request to `main`.

## Development

```bash
go fmt ./...
go vet ./...
go build -o leaf ./cmd/leaf
go run ./cmd/leaf
```

## Tech Stack

- **Go 1.25+** — Core language
- **Bubbletea** — TUI framework
- **Lipgloss** — Terminal styling
- **Goldmark** — Markdown parsing
- **Chroma** — Syntax highlighting

## License

MIT License — See [LICENSE](LICENSE) file for details.
