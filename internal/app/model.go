package app

import (
	"github.com/N95Ryan/leaf/internal/storage"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Mode represents the different modes of the application
type Mode int

const (
	ModeList Mode = iota
	ModeView
	ModeEdit
	ModeSearch
	ModeCreate
)

// SortMode represents the different ways to sort notes
type SortMode int

const (
	SortByUpdatedDesc SortMode = iota // Most recently updated first (default)
	SortByUpdatedAsc                  // Oldest updated first
	SortByCreatedDesc                 // Most recently created first
	SortByCreatedAsc                  // Oldest created first
	SortByTitleAsc                    // A-Z
	SortByTitleDesc                   // Z-A
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

	// Input components
	titleInput    textinput.Model
	contentEditor textarea.Model
	creatingNote  *storage.Note

	// Edit mode: "title" or "content"
	editMode string

	// Edit focus: which component has focus in ModeEdit ("title" or "content")
	editFocus string

	// Sort mode for notes list
	sortMode SortMode

	// Delete confirmation
	deleteConfirm bool
	noteToDelete  *storage.Note
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
		mode:          ModeList,
		notes:         []*storage.Note{},
		selectedIdx:   0,
		storage:       fs,
		lastError:     lastErr,
		titleInput:    newTitleInput(),
		contentEditor: newContentEditor(),
		creatingNote:  nil,
		editMode:      "title",
		editFocus:     "content",
		sortMode:      SortByUpdatedDesc,
		deleteConfirm: false,
		noteToDelete:  nil,
	}
}

// newTitleInput creates a new title input component
func newTitleInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your note title"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50
	return ti
}

// newContentEditor creates a new content editor component
func newContentEditor() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Write your note content here..."
	ta.CharLimit = 10000
	ta.SetWidth(80)
	ta.SetHeight(10) // Reduced from 15 to 10 to make room for title
	return ta
}

// Init is called when the program starts (Bubbletea)
func (m Model) Init() tea.Cmd {
	// If storage failed to initialize, don't try to load notes
	if m.storage == nil {
		return nil
	}

	// Return a command to load notes at startup
	return loadNotesCmd(m.storage)
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
