package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/N95Ryan/leaf/internal/app"
	"github.com/N95Ryan/leaf/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	notesDir := flag.String("notes-dir", "", "Directory to store notes (default: ~/.leaf/notes)")
	flag.Parse()

	fs, err := initStorage(*notesDir)
	var lastErr string
	if err != nil {
		lastErr = err.Error()
		fs = nil
	}

	m := app.NewModelWithStorage(fs, lastErr)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func initStorage(notesDir string) (storage.FileSystem, error) {
	if notesDir != "" {
		return storage.NewLocalFileSystemWithDir(notesDir)
	}
	return storage.NewLocalFileSystem()
}
