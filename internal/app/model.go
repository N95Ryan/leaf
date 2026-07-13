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
	SortByUpdatedDesc SortMode = iota
	SortByUpdatedAsc
	SortByCreatedDesc
	SortByCreatedAsc
	SortByTitleAsc
	SortByTitleDesc
)

// Model is the main application model (Elm Pattern)
type Model struct {
	mode Mode

	notes       []*storage.Note
	selectedIdx int
	currentNote *storage.Note

	storage storage.FileSystem

	lastError string

	width  int
	height int

	titleInput    textinput.Model
	contentEditor textarea.Model
	searchInput   textinput.Model
	creatingNote  *storage.Note

	editMode  string
	editFocus string

	sortMode SortMode

	deleteConfirm bool
	noteToDelete  *storage.Note

	searchSubmitted bool
	renderMarkdown  bool
}

// NewModel creates a new model with initial state
func NewModel() Model {
	fs, err := storage.NewLocalFileSystem()
	var lastErr string
	if err != nil {
		lastErr = err.Error()
		fs = nil
	}
	return NewModelWithStorage(fs, lastErr)
}

// NewModelWithStorage creates a model with the given storage backend.
func NewModelWithStorage(fs storage.FileSystem, lastErr string) Model {
	m := Model{
		mode:            ModeList,
		notes:           []*storage.Note{},
		selectedIdx:     0,
		storage:         fs,
		lastError:       lastErr,
		titleInput:      newTitleInput(),
		contentEditor:   newContentEditor(),
		searchInput:     newSearchInput(),
		creatingNote:    nil,
		editMode:        "title",
		editFocus:       "content",
		sortMode:        SortByUpdatedDesc,
		deleteConfirm:   false,
		noteToDelete:    nil,
		searchSubmitted: false,
		renderMarkdown:  true,
	}
	return m
}

func newTitleInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your note title"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50
	return ti
}

func newContentEditor() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Write your note content here..."
	ta.CharLimit = 10000
	ta.SetWidth(80)
	ta.SetHeight(10)
	return ta
}

func newSearchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search notes..."
	ti.CharLimit = 200
	ti.Width = 50
	return ti
}

func (m *Model) applyWindowSize() {
	if m.width <= 0 {
		return
	}
	inputWidth := m.width - 4
	if inputWidth < 20 {
		inputWidth = 20
	}
	m.titleInput.Width = inputWidth
	m.searchInput.Width = inputWidth
	m.contentEditor.SetWidth(inputWidth)

	editorHeight := m.height - 12
	if editorHeight < 5 {
		editorHeight = 5
	}
	m.contentEditor.SetHeight(editorHeight)
}

// Init is called when the program starts (Bubbletea)
func (m Model) Init() tea.Cmd {
	if m.storage == nil {
		return nil
	}
	return loadNotesCmd(m.storage)
}

func (m Model) Mode() Mode {
	return m.mode
}

func (m Model) Notes() []*storage.Note {
	return m.notes
}

func (m Model) Storage() storage.FileSystem {
	return m.storage
}

func (m Model) LastError() string {
	return m.lastError
}

func (m Model) SortMode() SortMode {
	return m.sortMode
}

func (m Model) SelectedIdx() int {
	return m.selectedIdx
}

func (m Model) Width() int {
	return m.width
}

func (m Model) Height() int {
	return m.height
}

func (m Model) SearchSubmitted() bool {
	return m.searchSubmitted
}
