package ui

import (
	"fmt"
	"strings"

	"github.com/N95Ryan/leaf/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
)

// RenderTitle renders the application title.
func RenderTitle(text string) string {
	return TitleStyle.Render(text)
}

// RenderAppTitle renders the main application header.
func RenderAppTitle(text string) string {
	return AppTitleStyle.Render(text)
}

// RenderNoteList renders a selectable list of notes.
func RenderNoteList(notes []*storage.Note, selectedIdx int) string {
	if len(notes) == 0 {
		return ListItemStyle.Render("No notes found.")
	}

	var b strings.Builder
	for i, note := range notes {
		line := fmt.Sprintf("%s", note.Title)
		if i == selectedIdx {
			b.WriteString(SelectedItemStyle.Render("> " + line))
		} else {
			b.WriteString(ListItemStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// RenderSearchBar renders the search input bar.
func RenderSearchBar(input textinput.Model) string {
	return TitleStyle.Render("Search") + "\n" + input.View()
}

// Shortcut represents a keyboard shortcut hint.
type Shortcut struct {
	Key   string
	Label string
}

// RenderShortcuts renders styled keyboard shortcut hints.
func RenderShortcuts(items []Shortcut) string {
	var b strings.Builder
	for i, item := range items {
		if i > 0 {
			b.WriteString(ShortcutSepStyle.Render(" │ "))
		}
		b.WriteString(ShortcutKeyStyle.Render(item.Key))
		b.WriteString(ShortcutLabelStyle.Render(" " + item.Label))
	}
	return StatusBarStyle.Render(b.String())
}

// RenderStatusBar renders keyboard shortcut hints from a plain string.
// Prefer RenderShortcuts for styled output.
func RenderStatusBar(shortcuts string) string {
	return StatusBarStyle.Render(shortcuts)
}

// RenderError renders an error message.
func RenderError(msg string) string {
	if msg == "" {
		return ""
	}
	return ErrorStyle.Render("Error: " + msg)
}

// RenderSectionTitle renders a section heading.
func RenderSectionTitle(text string) string {
	return TitleStyle.Render(text)
}

// RenderEditorSection wraps editor content with a border.
func RenderEditorSection(content string) string {
	return EditorStyle.Render(content)
}

// RenderSortIndicator renders the current sort mode label.
func RenderSortIndicator(sortName string) string {
	return ListItemStyle.Render(fmt.Sprintf("[Sort: %s]", sortName))
}

// RenderDeleteConfirm renders the delete confirmation prompt.
func RenderDeleteConfirm(title string) string {
	return ErrorStyle.Render(fmt.Sprintf("Press 'd' again to confirm deletion of '%s' (Esc to cancel)", title))
}

// RenderViewHeader renders a note title in read mode.
func RenderViewHeader(title string) string {
	return TitleStyle.Render("Reading: " + title)
}
