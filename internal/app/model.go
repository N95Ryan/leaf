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

	// UI
	width  int
	height int
}

// NewModel creates a new model with initial state
func NewModel() Model {
	return Model{
		mode:        ModeList,
		notes:       []*storage.Note{},
		selectedIdx: 0,
		storage:     nil, // TODO: Initialize concrete FileSystem
	}
}

// Init is called when the program starts (Bubbletea)
func (m Model) Init() tea.Cmd {
	// Return a command to load notes at startup
	return nil // TODO: Return a command to load notes
}
