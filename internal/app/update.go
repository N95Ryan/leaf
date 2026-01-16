package app

import (
	"context"
	"sort"
	"strings"

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

// NoteDeletedMsg is sent when a note is deleted
type NoteDeletedMsg struct {
	NoteID string
	Err    error
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
		m.sortNotes() // Apply current sort mode
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
		m.sortNotes() // Apply current sort mode after saving
		return m, loadNotesCmd(m.storage)

	case NoteDeletedMsg:
		if msg.Err != nil {
			// Store error message to display in view
			m.lastError = msg.Err.Error()
			m.deleteConfirm = false
			m.noteToDelete = nil
			return m, nil
		}
		// Clear any previous error and reload notes
		m.lastError = ""
		m.deleteConfirm = false
		m.noteToDelete = nil
		return m, loadNotesCmd(m.storage)

	default:
		return m, nil
	}
}

// handleKeyPress handles key presses based on current mode
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Special handling for ModeCreate: delegate based on editMode
	if m.mode == ModeCreate {
		return m.handleCreateMode(msg)
	}

	// Special handling for ModeEdit: delegate to textarea
	if m.mode == ModeEdit {
		return m.handleEditMode(msg)
	}

	// Special handling for ModeView: read-only mode
	if m.mode == ModeView {
		return m.handleViewMode(msg)
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "n":
		// Create a new note
		if m.mode == ModeList {
			m.mode = ModeCreate
			m.editMode = "title"         // Start with title
			m.titleInput.SetValue("")    // Reset input
			m.titleInput.Focus()         // Ensure it has focus
			m.contentEditor.SetValue("") // Reset content
			m.contentEditor.Blur()       // Blur content editor
			m.creatingNote = nil         // Clear any previous note
			return m, nil
		}

	case "r":
		// View selected note (read-only)
		if m.mode == ModeList && len(m.notes) > 0 {
			m.mode = ModeView
			m.currentNote = m.notes[m.selectedIdx]
			return m, nil
		}

	case "e":
		// Edit selected note
		if m.mode == ModeList && len(m.notes) > 0 {
			m.mode = ModeEdit
			m.currentNote = m.notes[m.selectedIdx]
			// Load note title and content into editors
			m.titleInput.SetValue(m.currentNote.Title)
			m.contentEditor.SetValue(m.currentNote.Content)
			// Start with content focused
			m.editFocus = "content"
			m.titleInput.Blur()
			m.contentEditor.Focus()
			return m, nil
		}

	case "/":
		// Activate search
		if m.mode == ModeList {
			m.mode = ModeSearch
			m.searchQuery = ""
			return m, nil
		}

	case "t":
		// Cycle through sort modes
		if m.mode == ModeList {
			m.sortMode = (m.sortMode + 1) % 6 // Cycle through 6 sort modes
			m.sortNotes()
			m.deleteConfirm = false // Cancel delete confirmation
			m.noteToDelete = nil
			return m, nil
		}

	case "d":
		// Delete selected note (with confirmation)
		if m.mode == ModeList && len(m.notes) > 0 {
			if !m.deleteConfirm {
				// First press: ask for confirmation
				m.deleteConfirm = true
				m.noteToDelete = m.notes[m.selectedIdx]
				return m, nil
			} else {
				// Second press: confirm deletion
				if m.noteToDelete != nil {
					note := m.noteToDelete
					m.deleteConfirm = false
					m.noteToDelete = nil
					return m, deleteNoteCmd(m.storage, note.ID)
				}
			}
		}

	case "esc":
		// Cancel delete confirmation if active
		if m.deleteConfirm {
			m.deleteConfirm = false
			m.noteToDelete = nil
			return m, nil
		}
		// Return to list
		if m.mode == ModeView || m.mode == ModeEdit || m.mode == ModeSearch || m.mode == ModeCreate {
			m.mode = ModeList
			m.currentNote = nil
			m.searchQuery = ""
			return m, nil
		}

	case "j", "down":
		// Navigate down in list
		if m.mode == ModeList && m.selectedIdx < len(m.notes)-1 {
			m.selectedIdx++
			m.deleteConfirm = false // Cancel delete confirmation on navigation
			m.noteToDelete = nil
			return m, nil
		}

	case "k", "up":
		// Navigate up in list
		if m.mode == ModeList && m.selectedIdx > 0 {
			m.selectedIdx--
			m.deleteConfirm = false // Cancel delete confirmation on navigation
			m.noteToDelete = nil
			return m, nil
		}
	}

	return m, nil
}

// handleCreateMode handles key presses in ModeCreate
func (m Model) handleCreateMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If we're editing title
	if m.editMode == "title" {
		switch msg.String() {
		case "esc":
			// Cancel creation and return to list
			m.mode = ModeList
			m.editMode = "title"
			m.titleInput.SetValue("")
			m.contentEditor.SetValue("")
			m.creatingNote = nil
			return m, nil

		case "enter":
			// Confirm title and move to content editing
			title := m.titleInput.Value()
			if title == "" {
				// Don't proceed with empty title
				return m, nil
			}

			// Switch to content editing mode
			m.editMode = "content"
			m.titleInput.Blur()
			m.contentEditor.Focus()
			m.creatingNote = storage.NewNote(title, "")
			return m, nil

		default:
			// Delegate to textinput
			var cmd tea.Cmd
			m.titleInput, cmd = m.titleInput.Update(msg)
			return m, cmd
		}
	}

	// If we're editing content
	if m.editMode == "content" {
		switch msg.String() {
		case "esc":
			// Go back to title editing
			m.editMode = "title"
			m.contentEditor.Blur()
			m.titleInput.Focus()
			return m, nil

		case "ctrl+s":
			// Save the note with content
			if m.creatingNote == nil {
				return m, nil
			}

			// Update note content
			m.creatingNote.Content = m.contentEditor.Value()

			// Reset and return to list
			m.mode = ModeList
			m.editMode = "title"
			m.titleInput.SetValue("")
			m.contentEditor.SetValue("")

			// Save asynchronously
			return m, saveNoteCmd(m.storage, m.creatingNote)

		default:
			// Delegate to textarea
			var cmd tea.Cmd
			m.contentEditor, cmd = m.contentEditor.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// handleViewMode handles key presses in ModeView (read-only)
func (m Model) handleViewMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Return to list
		m.mode = ModeList
		m.currentNote = nil
		return m, nil

	case "i", "e":
		// Switch to edit mode
		if m.currentNote == nil {
			return m, nil
		}
		m.mode = ModeEdit
		// Load note title and content into editors
		m.titleInput.SetValue(m.currentNote.Title)
		m.contentEditor.SetValue(m.currentNote.Content)
		// Start with content focused
		m.editFocus = "content"
		m.titleInput.Blur()
		m.contentEditor.Focus()
		return m, nil
	}

	return m, nil
}

// handleEditMode handles key presses in ModeEdit
func (m Model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel editing and return to list
		m.mode = ModeList
		m.currentNote = nil
		m.titleInput.SetValue("")
		m.titleInput.Blur()
		m.contentEditor.SetValue("")
		m.contentEditor.Blur()
		m.editFocus = "content"
		return m, nil

	case "tab":
		// Toggle focus between title and content
		if m.editFocus == "title" {
			m.editFocus = "content"
			m.titleInput.Blur()
			m.contentEditor.Focus()
		} else {
			m.editFocus = "title"
			m.contentEditor.Blur()
			m.titleInput.Focus()
		}
		return m, nil

	case "ctrl+s":
		// Save the edited note
		if m.currentNote == nil {
			return m, nil
		}

		// Update note title and content
		newTitle := m.titleInput.Value()
		if newTitle == "" {
			// Don't save with empty title
			m.lastError = "Title cannot be empty"
			return m, nil
		}

		m.currentNote.Title = newTitle
		m.currentNote.Content = m.contentEditor.Value()

		// Return to list
		m.mode = ModeList
		note := m.currentNote
		m.currentNote = nil
		m.titleInput.SetValue("")
		m.titleInput.Blur()
		m.contentEditor.SetValue("")
		m.contentEditor.Blur()
		m.editFocus = "content"

		// Save asynchronously
		return m, saveNoteCmd(m.storage, note)

	default:
		// Delegate to the focused component
		var cmd tea.Cmd
		if m.editFocus == "title" {
			m.titleInput, cmd = m.titleInput.Update(msg)
		} else {
			m.contentEditor, cmd = m.contentEditor.Update(msg)
		}
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

// deleteNoteCmd is a command that deletes a note from storage
// It runs asynchronously and returns a NoteDeletedMsg
func deleteNoteCmd(fs storage.FileSystem, noteID string) tea.Cmd {
	return func() tea.Msg {
		// Use background context for deleting note
		err := fs.DeleteNote(context.Background(), noteID)
		return NoteDeletedMsg{
			NoteID: noteID,
			Err:    err,
		}
	}
}

// sortNotes sorts the notes list according to the current sort mode
func (m *Model) sortNotes() {
	switch m.sortMode {
	case SortByUpdatedDesc:
		sort.Slice(m.notes, func(i, j int) bool {
			return m.notes[i].UpdatedAt.After(m.notes[j].UpdatedAt)
		})
	case SortByUpdatedAsc:
		sort.Slice(m.notes, func(i, j int) bool {
			return m.notes[i].UpdatedAt.Before(m.notes[j].UpdatedAt)
		})
	case SortByCreatedDesc:
		sort.Slice(m.notes, func(i, j int) bool {
			return m.notes[i].CreatedAt.After(m.notes[j].CreatedAt)
		})
	case SortByCreatedAsc:
		sort.Slice(m.notes, func(i, j int) bool {
			return m.notes[i].CreatedAt.Before(m.notes[j].CreatedAt)
		})
	case SortByTitleAsc:
		sort.Slice(m.notes, func(i, j int) bool {
			return strings.ToLower(m.notes[i].Title) < strings.ToLower(m.notes[j].Title)
		})
	case SortByTitleDesc:
		sort.Slice(m.notes, func(i, j int) bool {
			return strings.ToLower(m.notes[i].Title) > strings.ToLower(m.notes[j].Title)
		})
	}
}
