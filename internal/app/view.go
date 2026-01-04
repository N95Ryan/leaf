package app

import (
	"fmt"
	"strings"
)

// View renders the user interface based on model state (Elm Pattern)
// This function must be pure: no logic, only rendering
func (m Model) View() string {
	switch m.mode {
	case ModeList:
		return m.renderList()
	case ModeEdit:
		return m.renderEdit()
	case ModeSearch:
		return m.renderSearch()
	case ModeCreate:
		return m.renderCreate()
	default:
		return "Unknown mode"
	}
}

// renderList displays the list of notes
func (m Model) renderList() string {
	var b strings.Builder

	b.WriteString("🌱 Leaf - Note Manager\n\n")

	if len(m.notes) == 0 {
		b.WriteString("No notes. Press 'n' to create a note.\n")
	} else {
		for i, note := range m.notes {
			prefix := "  "
			if i == m.selectedIdx {
				prefix = "> "
			}
			b.WriteString(fmt.Sprintf("%s%s\n", prefix, note.Title))
		}
	}

	b.WriteString("\nShortcuts: n (new), e (edit), / (search), q (quit)")

	return b.String()
}

// renderEdit displays the note editor
func (m Model) renderEdit() string {
	if m.currentNote == nil {
		return "No note selected"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Editing: %s\n\n", m.currentNote.Title))
	b.WriteString(m.currentNote.Content)
	b.WriteString("\n\nPress 'esc' to return to list")

	return b.String()
}

// renderSearch displays the search interface
func (m Model) renderSearch() string {
	var b strings.Builder
	b.WriteString("Search: ")
	b.WriteString(m.searchQuery)
	b.WriteString("_\n\n")
	b.WriteString("Type your search and press 'esc' to cancel")

	return b.String()
}

// renderCreate displays the note creation interface
func (m Model) renderCreate() string {
	return "Create a new note\n\n" +
		"Title: _\n\n" +
		"Press 'esc' to cancel"
}

// TODO: Add Lipgloss styling to improve appearance
