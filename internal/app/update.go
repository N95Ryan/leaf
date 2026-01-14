package app

import (
	"context"

	"github.com/N95Ryan/leaf/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// Message represents the different types of messages in the Elm pattern
type Message interface{}

// KeyMsg represents a key press
type KeyMsg struct {
	Key string
}

// NoteLoadedMsg is sent when notes are loaded
type NoteLoadedMsg struct {
	Notes []*storage.Note
	Err   error
}

// NoteSavedMsg is sent when a note is saved
type NoteSavedMsg struct {
	Note *storage.Note
	Err  error
}

// Update handles messages and returns a new model (Elm Pattern)
// This function must be pure: no side effects, only state mutation
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case NoteLoadedMsg:
		if msg.Err != nil {
			// Store error message to display in view
			m.lastError = msg.Err.Error()
			return m, nil
		}
		// Clear any previous error and store loaded notes
		m.lastError = ""
		m.notes = msg.Notes
		return m, nil

	case NoteSavedMsg:
		if msg.Err != nil {
			// Store error message to display in view
			m.lastError = msg.Err.Error()
			return m, nil
		}
		// Clear any previous error and reload notes to show the new one
		m.lastError = ""
		m.creatingNote = nil
		return m, loadNotesCmd(m.storage)

	default:
		return m, nil
	}
}

// handleKeyPress handles key presses based on current mode
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Special handling for ModeCreate: delegate to textinput
	if m.mode == ModeCreate {
		return m.handleCreateMode(msg)
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "n":
		// Create a new note
		if m.mode == ModeList {
			m.mode = ModeCreate
			m.titleInput.SetValue("") // Reset input
			m.titleInput.Focus()      // Ensure it has focus
			m.creatingNote = nil      // Clear any previous note
			return m, nil
		}

	case "e":
		// Edit selected note
		if m.mode == ModeList && len(m.notes) > 0 {
			m.mode = ModeEdit
			m.currentNote = m.notes[m.selectedIdx]
			return m, nil
		}

	case "/":
		// Activate search
		if m.mode == ModeList {
			m.mode = ModeSearch
			m.searchQuery = ""
			return m, nil
		}

	case "esc":
		// Return to list
		if m.mode == ModeEdit || m.mode == ModeSearch || m.mode == ModeCreate {
			m.mode = ModeList
			m.currentNote = nil
			m.searchQuery = ""
			return m, nil
		}

	case "j", "down":
		// Navigate down in list
		if m.mode == ModeList && m.selectedIdx < len(m.notes)-1 {
			m.selectedIdx++
			return m, nil
		}

	case "k", "up":
		// Navigate up in list
		if m.mode == ModeList && m.selectedIdx > 0 {
			m.selectedIdx--
			return m, nil
		}
	}

	return m, nil
}

// handleCreateMode handles key presses in ModeCreate
func (m Model) handleCreateMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel creation and return to list
		m.mode = ModeList
		m.titleInput.SetValue("") // Reset input
		m.creatingNote = nil
		return m, nil

	case "enter":
		// Confirm title and prepare note creation
		title := m.titleInput.Value()
		if title == "" {
			// Don't create note with empty title
			return m, nil
		}

		// Create the note with title
		note := storage.NewNote(title, "")
		m.creatingNote = note

		// For now, save immediately and return to list
		// Phase 3.3 will add content editing before saving
		m.mode = ModeList
		m.titleInput.SetValue("") // Reset for next time

		// Save the note asynchronously
		return m, saveNoteCmd(m.storage, note)

	default:
		// Delegate all other keys to textinput
		var cmd tea.Cmd
		m.titleInput, cmd = m.titleInput.Update(msg)
		return m, cmd
	}
}

// loadNotesCmd is a command that loads all notes from storage
// It runs asynchronously and returns a NoteLoadedMsg
func loadNotesCmd(fs storage.FileSystem) tea.Cmd {
	return func() tea.Msg {
		// Use background context for loading notes
		notes, err := fs.ListNotes(context.Background())
		return NoteLoadedMsg{
			Notes: notes,
			Err:   err,
		}
	}
}

// saveNoteCmd is a command that saves a note to storage
// It runs asynchronously and returns a NoteSavedMsg
func saveNoteCmd(fs storage.FileSystem, note *storage.Note) tea.Cmd {
	return func() tea.Msg {
		// Use background context for saving note
		err := fs.SaveNote(context.Background(), note)
		return NoteSavedMsg{
			Note: note,
			Err:  err,
		}
	}
}
