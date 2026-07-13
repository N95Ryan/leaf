package app

import (
	"fmt"
	"strings"

	"github.com/N95Ryan/leaf/internal/markdown"
	"github.com/N95Ryan/leaf/internal/ui"
)

// View renders the user interface based on model state (Elm Pattern)
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

func (m Model) renderList() string {
	var b strings.Builder
	b.WriteString(ui.RenderAppTitle("Leaf 🌴 - Note Manager"))
	b.WriteString("\n\n")

	if len(m.notes) == 0 {
		b.WriteString(ui.ListItemStyle.Render("No notes. Press 'n' to create a note."))
		b.WriteString("\n")
	} else {
		b.WriteString(ui.RenderNoteList(m.notes, m.selectedIdx))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(ui.RenderShortcuts(ui.ListShortcuts))
	b.WriteString("\n")
	b.WriteString(ui.RenderSortIndicator(m.sortName()))
	if m.deleteConfirm && m.noteToDelete != nil {
		b.WriteString("\n")
		b.WriteString(ui.RenderDeleteConfirm(m.noteToDelete.Title))
	}
	if err := ui.RenderError(m.lastError); err != "" {
		b.WriteString("\n")
		b.WriteString(err)
	}

	return b.String()
}

func (m Model) renderView() string {
	if m.currentNote == nil {
		return ui.RenderError("No note selected")
	}

	var b strings.Builder
	b.WriteString(ui.RenderViewHeader(m.currentNote.Title))
	b.WriteString("\n\n")

	content := m.currentNote.Content
	if m.renderMarkdown {
		rendered, err := markdown.Render(content, m.contentWidth())
		if err != nil {
			b.WriteString(content)
			b.WriteString("\n\n")
			b.WriteString(ui.RenderError(err.Error()))
		} else {
			b.WriteString(rendered)
		}
	} else {
		b.WriteString(content)
	}

	b.WriteString("\n\n")
	b.WriteString(ui.RenderShortcuts(ui.ViewShortcuts))
	if err := ui.RenderError(m.lastError); err != "" {
		b.WriteString("\n")
		b.WriteString(err)
	}

	return b.String()
}

func (m Model) renderEdit() string {
	if m.currentNote == nil {
		return ui.RenderError("No note selected")
	}

	var b strings.Builder
	b.WriteString(ui.RenderSectionTitle("Editing note"))
	b.WriteString("\n\n")

	focusIndicator := " "
	if m.editFocus == "title" {
		focusIndicator = ">"
	}
	b.WriteString(fmt.Sprintf("%s Title:\n", focusIndicator))
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\n")

	focusIndicator = " "
	if m.editFocus == "content" {
		focusIndicator = ">"
	}
	b.WriteString(fmt.Sprintf("%s Content:\n", focusIndicator))
	b.WriteString(ui.RenderEditorSection(m.contentEditor.View()))

	b.WriteString("\n\n")
	b.WriteString(ui.RenderShortcuts(ui.EditShortcuts))
	if err := ui.RenderError(m.lastError); err != "" {
		b.WriteString("\n")
		b.WriteString(err)
	}

	return b.String()
}

func (m Model) renderSearch() string {
	var b strings.Builder
	b.WriteString(ui.RenderSearchBar(m.searchInput))
	b.WriteString("\n\n")

	if m.searchSubmitted {
		if len(m.notes) == 0 {
			b.WriteString(ui.ListItemStyle.Render("No matching notes."))
		} else {
			b.WriteString(ui.RenderNoteList(m.notes, m.selectedIdx))
		}
		b.WriteString("\n\n")
		b.WriteString(ui.RenderShortcuts(ui.SearchResultsShortcuts))
	} else {
		b.WriteString(ui.RenderShortcuts(ui.SearchInputShortcuts))
	}

	if err := ui.RenderError(m.lastError); err != "" {
		b.WriteString("\n")
		b.WriteString(err)
	}

	return b.String()
}

func (m Model) renderCreate() string {
	var b strings.Builder
	b.WriteString(ui.RenderTitle("Create a new note"))
	b.WriteString("\n\n")

	if m.editMode == "title" {
		b.WriteString("Title:\n")
		b.WriteString(m.titleInput.View())
		b.WriteString("\n\n")
		b.WriteString(ui.RenderShortcuts(ui.CreateTitleShortcuts))
	} else {
		if m.creatingNote != nil {
			b.WriteString(fmt.Sprintf("Title: %s\n\n", m.creatingNote.Title))
		}
		b.WriteString("Content:\n")
		b.WriteString(ui.RenderEditorSection(m.contentEditor.View()))
		b.WriteString("\n\n")
		b.WriteString(ui.RenderShortcuts(ui.CreateContentShortcuts))
	}

	if err := ui.RenderError(m.lastError); err != "" {
		b.WriteString("\n")
		b.WriteString(err)
	}

	return b.String()
}

func (m Model) sortName() string {
	switch m.sortMode {
	case SortByUpdatedDesc:
		return "Updated (newest)"
	case SortByUpdatedAsc:
		return "Updated (oldest)"
	case SortByCreatedDesc:
		return "Created (newest)"
	case SortByCreatedAsc:
		return "Created (oldest)"
	case SortByTitleAsc:
		return "Title A-Z"
	case SortByTitleDesc:
		return "Title Z-A"
	default:
		return "Unknown"
	}
}

func (m Model) contentWidth() int {
	width := m.width - 4
	if width < 20 {
		return 80
	}
	return width
}
