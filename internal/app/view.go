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
	case ModeView:
		return m.renderView()
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

	b.WriteString("\nShortcuts: n (new), r (read), e (edit), t (sort), d (delete), q (quit)")
	b.WriteString(m.renderSortIndicator())
	b.WriteString(m.renderDeleteConfirm())
	b.WriteString(m.renderError())

	return b.String()
}

// renderView displays the note in read-only mode
func (m Model) renderView() string {
	if m.currentNote == nil {
		return "No note selected"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("📖 %s\n\n", m.currentNote.Title))
	b.WriteString(m.currentNote.Content)
	b.WriteString("\n\nShortcuts: i/e (edit), Esc (back to list)")
	b.WriteString(m.renderError())

	return b.String()
}

// renderEdit displays the note editor
func (m Model) renderEdit() string {
	if m.currentNote == nil {
		return "No note selected"
	}

	var b strings.Builder
	b.WriteString("✏️  Editing note\n\n")

	// Show title input
	focusIndicator := " "
	if m.editFocus == "title" {
		focusIndicator = "›"
	}
	b.WriteString(fmt.Sprintf("%s Title:\n", focusIndicator))
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\n")

	// Show content editor
	focusIndicator = " "
	if m.editFocus == "content" {
		focusIndicator = "›"
	}
	b.WriteString(fmt.Sprintf("%s Content:\n", focusIndicator))
	b.WriteString(m.contentEditor.View())

	b.WriteString("\n\nShortcuts: Tab (switch field), Ctrl+S (save), Esc (cancel)")
	b.WriteString(m.renderError())

	return b.String()
}

// renderSearch displays the search interface
func (m Model) renderSearch() string {
	var b strings.Builder
	b.WriteString("Search: ")
	b.WriteString(m.searchQuery)
	b.WriteString("_\n\n")
	b.WriteString("Type your search and press 'esc' to cancel")
	b.WriteString(m.renderError())

	return b.String()
}

// renderCreate displays the note creation interface
func (m Model) renderCreate() string {
	var b strings.Builder

	b.WriteString("🌱 Create a new note\n\n")

	// Show title input or content editor based on editMode
	if m.editMode == "title" {
		b.WriteString("Title:\n")
		b.WriteString(m.titleInput.View())
		b.WriteString("\n\n")
		b.WriteString("Shortcuts: Enter (next), Esc (cancel)")
	} else {
		// Show title as read-only and content editor
		if m.creatingNote != nil {
			b.WriteString(fmt.Sprintf("Title: %s\n\n", m.creatingNote.Title))
		}
		b.WriteString("Content:\n")
		b.WriteString(m.contentEditor.View())
		b.WriteString("\n\n")
		b.WriteString("Shortcuts: Ctrl+S (save), Esc (back to title)")
	}

	b.WriteString(m.renderError())

	return b.String()
}

// renderDeleteConfirm displays the delete confirmation message
func (m Model) renderDeleteConfirm() string {
	if !m.deleteConfirm || m.noteToDelete == nil {
		return ""
	}
	return fmt.Sprintf("\n⚠️  Press 'd' again to confirm deletion of '%s' (Esc to cancel)", m.noteToDelete.Title)
}

// renderSortIndicator displays the current sort mode
func (m Model) renderSortIndicator() string {
	var sortName string
	switch m.sortMode {
	case SortByUpdatedDesc:
		sortName = "Updated ↓"
	case SortByUpdatedAsc:
		sortName = "Updated ↑"
	case SortByCreatedDesc:
		sortName = "Created ↓"
	case SortByCreatedAsc:
		sortName = "Created ↑"
	case SortByTitleAsc:
		sortName = "Title A-Z"
	case SortByTitleDesc:
		sortName = "Title Z-A"
	}
	return fmt.Sprintf("\n[Sort: %s]", sortName)
}

// renderError displays error messages if any
func (m Model) renderError() string {
	if m.lastError == "" {
		return ""
	}
	return fmt.Sprintf("\n❌ Error: %s\n", m.lastError)
}

// TODO: Add Lipgloss styling to improve appearance
