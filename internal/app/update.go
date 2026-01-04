package app

import (
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
			// TODO: Handle error (display in view)
			return m, nil
		}
		m.notes = msg.Notes
		return m, nil

	case NoteSavedMsg:
		if msg.Err != nil {
			// TODO: Handle error
			return m, nil
		}
		// TODO: Update note list
		return m, nil

	default:
		return m, nil
	}
}

// handleKeyPress handles key presses based on current mode
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "n":
		// Create a new note
		if m.mode == ModeList {
			m.mode = ModeCreate
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
