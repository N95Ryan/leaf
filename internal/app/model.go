package app

import (
	"github.com/N95Ryan/leaf/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// Mode represents the different modes of the application
type Mode int

const (
	ModeList Mode = iota
	ModeEdit
	ModeSearch
	ModeCreate
)

// Model is the main application model (Elm Pattern)
type Model struct {
	// Application state
	mode Mode

	// Notes
	notes       []*storage.Note
	selectedIdx int
	currentNote *storage.Note

	// Search
	searchQuery string

	// Storage
	storage storage.FileSystem

	// Error handling
	lastError string

	// UI
	width  int
	height int
}

// NewModel creates a new model with initial state
func NewModel() Model {
	// Initialize the local filesystem storage
	fs, err := storage.NewLocalFileSystem()

	var lastErr string
	if err != nil {
		// Store the error to display in the UI
		lastErr = err.Error()
		fs = nil // Ensure storage is nil on error
	}

	return Model{
		mode:        ModeList,
		notes:       []*storage.Note{},
		selectedIdx: 0,
		storage:     fs,
		lastError:   lastErr,
	}
}

// Init is called when the program starts (Bubbletea)
func (m Model) Init() tea.Cmd {
	// Return a command to load notes at startup
	return nil // TODO: Return a command to load notes
}

// Getters for testing and external access

// Mode returns the current mode
func (m Model) Mode() Mode {
	return m.mode
}

// Notes returns the current notes list
func (m Model) Notes() []*storage.Note {
	return m.notes
}

// Storage returns the storage instance
func (m Model) Storage() storage.FileSystem {
	return m.storage
}

// LastError returns the last error message
func (m Model) LastError() string {
	return m.lastError
}
