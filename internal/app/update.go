package app

import (
	"context"
	"sort"
	"strings"

	"github.com/N95Ryan/leaf/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

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

// SearchResultsMsg is sent when search completes
type SearchResultsMsg struct {
	Notes []*storage.Note
	Query string
	Err   error
}

// Update handles messages and returns a new model (Elm Pattern)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.applyWindowSize()
		return m, nil

	case NoteLoadedMsg:
		if msg.Err != nil {
			m.lastError = msg.Err.Error()
			return m, nil
		}
		m.lastError = ""
		m.notes = msg.Notes
		m.sortNotes()
		if m.selectedIdx >= len(m.notes) {
			m.selectedIdx = max(0, len(m.notes)-1)
		}
		return m, nil

	case NoteSavedMsg:
		if msg.Err != nil {
			m.lastError = msg.Err.Error()
			return m, nil
		}
		m.lastError = ""
		m.creatingNote = nil
		return m, loadNotesCmd(m.storage)

	case NoteDeletedMsg:
		if msg.Err != nil {
			m.lastError = msg.Err.Error()
			m.deleteConfirm = false
			m.noteToDelete = nil
			return m, nil
		}
		m.lastError = ""
		m.deleteConfirm = false
		m.noteToDelete = nil
		return m, loadNotesCmd(m.storage)

	case SearchResultsMsg:
		if msg.Err != nil {
			m.lastError = msg.Err.Error()
			return m, nil
		}
		m.lastError = ""
		m.notes = msg.Notes
		m.selectedIdx = 0
		m.searchSubmitted = true
		m.sortNotes()
		return m, nil

	default:
		return m, nil
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.mode == ModeCreate {
		return m.handleCreateMode(msg)
	}
	if m.mode == ModeEdit {
		return m.handleEditMode(msg)
	}
	if m.mode == ModeView {
		return m.handleViewMode(msg)
	}
	if m.mode == ModeSearch {
		return m.handleSearchMode(msg)
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "n":
		if m.mode == ModeList {
			m.mode = ModeCreate
			m.editMode = "title"
			m.titleInput.SetValue("")
			m.titleInput.Focus()
			m.contentEditor.SetValue("")
			m.contentEditor.Blur()
			m.creatingNote = nil
			return m, nil
		}

	case "r":
		if m.mode == ModeList && len(m.notes) > 0 {
			m.mode = ModeView
			m.currentNote = m.notes[m.selectedIdx]
			return m, nil
		}

	case "e":
		if m.mode == ModeList && len(m.notes) > 0 {
			m.mode = ModeEdit
			m.currentNote = m.notes[m.selectedIdx]
			m.titleInput.SetValue(m.currentNote.Title)
			m.contentEditor.SetValue(m.currentNote.Content)
			m.editFocus = "content"
			m.titleInput.Blur()
			m.contentEditor.Focus()
			return m, nil
		}

	case "/":
		if m.mode == ModeList {
			m.mode = ModeSearch
			m.searchInput.SetValue("")
			m.searchInput.Focus()
			m.searchSubmitted = false
			return m, nil
		}

	case "t":
		if m.mode == ModeList {
			m.sortMode = (m.sortMode + 1) % 6
			m.sortNotes()
			m.deleteConfirm = false
			m.noteToDelete = nil
			return m, nil
		}

	case "d":
		if m.mode == ModeList && len(m.notes) > 0 {
			if !m.deleteConfirm {
				m.deleteConfirm = true
				m.noteToDelete = m.notes[m.selectedIdx]
				return m, nil
			}
			if m.noteToDelete != nil {
				note := m.noteToDelete
				m.deleteConfirm = false
				m.noteToDelete = nil
				return m, deleteNoteCmd(m.storage, note.ID)
			}
		}

	case "esc":
		if m.deleteConfirm {
			m.deleteConfirm = false
			m.noteToDelete = nil
			return m, nil
		}
		if m.mode == ModeView || m.mode == ModeEdit || m.mode == ModeSearch || m.mode == ModeCreate {
			return m.exitToList(false)
		}

	case "j", "down":
		if m.mode == ModeList && m.selectedIdx < len(m.notes)-1 {
			m.selectedIdx++
			m.deleteConfirm = false
			m.noteToDelete = nil
			return m, nil
		}

	case "k", "up":
		if m.mode == ModeList && m.selectedIdx > 0 {
			m.selectedIdx--
			m.deleteConfirm = false
			m.noteToDelete = nil
			return m, nil
		}
	}

	return m, nil
}

func (m Model) exitToList(reload bool) (Model, tea.Cmd) {
	m.mode = ModeList
	m.currentNote = nil
	m.searchInput.SetValue("")
	m.searchInput.Blur()
	m.searchSubmitted = false
	if reload && m.storage != nil {
		return m, loadNotesCmd(m.storage)
	}
	return m, nil
}

func (m Model) handleSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return m.exitToList(true)

	case "enter":
		if m.storage == nil {
			return m, nil
		}
		query := strings.TrimSpace(m.searchInput.Value())
		if query == "" {
			return m, nil
		}
		return m, searchNotesCmd(m.storage, query)

	case "r":
		if m.searchSubmitted && len(m.notes) > 0 {
			m.mode = ModeView
			m.currentNote = m.notes[m.selectedIdx]
			return m, nil
		}

	case "j", "down":
		if m.searchSubmitted && m.selectedIdx < len(m.notes)-1 {
			m.selectedIdx++
			return m, nil
		}

	case "k", "up":
		if m.searchSubmitted && m.selectedIdx > 0 {
			m.selectedIdx--
			return m, nil
		}

	default:
		if !m.searchSubmitted {
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) handleCreateMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editMode == "title" {
		switch msg.String() {
		case "esc":
			updated, cmd := m.exitToList(false)
			updated.editMode = "title"
			updated.titleInput.SetValue("")
			updated.contentEditor.SetValue("")
			updated.creatingNote = nil
			return updated, cmd

		case "enter":
			title := m.titleInput.Value()
			if title == "" {
				return m, nil
			}
			m.editMode = "content"
			m.titleInput.Blur()
			m.contentEditor.Focus()
			m.creatingNote = storage.NewNote(title, "")
			return m, nil

		default:
			var cmd tea.Cmd
			m.titleInput, cmd = m.titleInput.Update(msg)
			return m, cmd
		}
	}

	if m.editMode == "content" {
		switch msg.String() {
		case "esc":
			m.editMode = "title"
			m.contentEditor.Blur()
			m.titleInput.Focus()
			return m, nil

		case "ctrl+s":
			if m.creatingNote == nil {
				return m, nil
			}
			m.creatingNote.Content = m.contentEditor.Value()
			note := m.creatingNote
			m.mode = ModeList
			m.editMode = "title"
			m.titleInput.SetValue("")
			m.contentEditor.SetValue("")
			return m, saveNoteCmd(m.storage, note)

		default:
			var cmd tea.Cmd
			m.contentEditor, cmd = m.contentEditor.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) handleViewMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeList
		m.currentNote = nil
		return m, nil

	case "i", "e":
		if m.currentNote == nil {
			return m, nil
		}
		m.mode = ModeEdit
		m.titleInput.SetValue(m.currentNote.Title)
		m.contentEditor.SetValue(m.currentNote.Content)
		m.editFocus = "content"
		m.titleInput.Blur()
		m.contentEditor.Focus()
		return m, nil
	}

	return m, nil
}

func (m Model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeList
		m.currentNote = nil
		m.titleInput.SetValue("")
		m.titleInput.Blur()
		m.contentEditor.SetValue("")
		m.contentEditor.Blur()
		m.editFocus = "content"
		return m, nil

	case "tab":
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
		if m.currentNote == nil {
			return m, nil
		}
		newTitle := m.titleInput.Value()
		if newTitle == "" {
			m.lastError = "Title cannot be empty"
			return m, nil
		}
		m.currentNote.Title = newTitle
		m.currentNote.Content = m.contentEditor.Value()
		note := m.currentNote
		m.mode = ModeList
		m.currentNote = nil
		m.titleInput.SetValue("")
		m.titleInput.Blur()
		m.contentEditor.SetValue("")
		m.contentEditor.Blur()
		m.editFocus = "content"
		return m, saveNoteCmd(m.storage, note)

	default:
		var cmd tea.Cmd
		if m.editFocus == "title" {
			m.titleInput, cmd = m.titleInput.Update(msg)
		} else {
			m.contentEditor, cmd = m.contentEditor.Update(msg)
		}
		return m, cmd
	}
}

func loadNotesCmd(fs storage.FileSystem) tea.Cmd {
	return func() tea.Msg {
		notes, err := fs.ListNotes(context.Background())
		return NoteLoadedMsg{Notes: notes, Err: err}
	}
}

func saveNoteCmd(fs storage.FileSystem, note *storage.Note) tea.Cmd {
	return func() tea.Msg {
		err := fs.SaveNote(context.Background(), note)
		return NoteSavedMsg{Note: note, Err: err}
	}
}

func deleteNoteCmd(fs storage.FileSystem, noteID string) tea.Cmd {
	return func() tea.Msg {
		err := fs.DeleteNote(context.Background(), noteID)
		return NoteDeletedMsg{NoteID: noteID, Err: err}
	}
}

func searchNotesCmd(fs storage.FileSystem, query string) tea.Cmd {
	return func() tea.Msg {
		notes, err := fs.SearchNotes(context.Background(), query)
		return SearchResultsMsg{Notes: notes, Query: query, Err: err}
	}
}

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
