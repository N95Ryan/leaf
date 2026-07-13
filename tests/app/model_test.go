package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/N95Ryan/leaf/internal/app"
	"github.com/N95Ryan/leaf/internal/storage"
	"github.com/N95Ryan/leaf/tests/testutil"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	t.Run("should initialize model successfully", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModelWithStorage(&mockFileSystem{}, "")

		assert.Equal(app.ModeList, model.Mode(), "mode should be ModeList by default")
		assert.NotNil(model.Notes(), "notes should be initialized")
		assert.Empty(model.Notes(), "notes should be empty initially")
		assert.NotNil(model.Storage(), "storage should be initialized")
		assert.Empty(model.LastError(), "should have no error on successful initialization")
	})
}

func TestNewModel_StorageInitialization(t *testing.T) {
	t.Run("should have storage initialized or error set", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModel()

		if model.Storage() == nil {
			assert.NotEmpty(model.LastError(),
				"if storage initialization fails, lastError should contain error message")
		} else {
			assert.Empty(model.LastError(),
				"if storage initialization succeeds, lastError should be empty")
		}
	})
}

type mockFileSystem struct {
	notes           []*storage.Note
	listErr         error
	saveErr         error
	deleteErr       error
	searchErr       error
	searchResults   []*storage.Note
	listCalled      int
	saveCalled      int
	deleteCalled    int
	searchCalled    int
	lastSearchQuery string
	lastSavedNote   *storage.Note
	lastDeletedID   string
}

func (m *mockFileSystem) ListNotes(ctx context.Context) ([]*storage.Note, error) {
	m.listCalled++
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.notes, nil
}

func (m *mockFileSystem) GetNote(ctx context.Context, id string) (*storage.Note, error) {
	return nil, nil
}

func (m *mockFileSystem) SaveNote(ctx context.Context, note *storage.Note) error {
	m.saveCalled++
	m.lastSavedNote = note
	return m.saveErr
}

func (m *mockFileSystem) DeleteNote(ctx context.Context, id string) error {
	m.deleteCalled++
	m.lastDeletedID = id
	return m.deleteErr
}

func (m *mockFileSystem) SearchNotes(ctx context.Context, query string) ([]*storage.Note, error) {
	m.searchCalled++
	m.lastSearchQuery = query
	if m.searchErr != nil {
		return nil, m.searchErr
	}
	if m.searchResults != nil {
		return m.searchResults, nil
	}
	return nil, nil
}

func newTestModel(mock *mockFileSystem) app.Model {
	notes := []*storage.Note{
		{ID: "1", Title: "Alpha", Content: "first", UpdatedAt: time.Now()},
		{ID: "2", Title: "Beta", Content: "second", UpdatedAt: time.Now().Add(-time.Hour)},
	}
	mock.notes = notes
	model := app.NewModelWithStorage(mock, "")
	updated, _ := model.Update(app.NoteLoadedMsg{Notes: notes})
	return updated.(app.Model)
}

func TestInit(t *testing.T) {
	t.Run("should return nil when storage is nil", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModelWithStorage(nil, "storage error")

		cmd := model.Init()
		assert.Nil(cmd, "Init should return nil when storage is nil")
	})

	t.Run("should return a valid command when storage is initialized", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModelWithStorage(&mockFileSystem{}, "")

		cmd := model.Init()
		assert.NotNil(cmd, "Init should return a command to load notes")
	})
}

func TestUpdate_NotesLoadedMsg(t *testing.T) {
	t.Run("should store notes when loaded successfully", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModelWithStorage(&mockFileSystem{}, "")

		testNotes := []*storage.Note{
			{ID: "1", Title: "Note 1", Content: "Content 1"},
			{ID: "2", Title: "Note 2", Content: "Content 2"},
		}

		updatedModel, _ := model.Update(app.NoteLoadedMsg{Notes: testNotes})
		m := updatedModel.(app.Model)

		assert.Len(m.Notes(), 2, "should have 2 notes loaded")
		assert.Equal("Note 1", m.Notes()[0].Title, "first note should have correct title")
		assert.Empty(m.LastError(), "lastError should be cleared on successful load")
	})

	t.Run("should store error when loading fails", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModelWithStorage(&mockFileSystem{}, "")

		updatedModel, _ := model.Update(app.NoteLoadedMsg{
			Notes: nil,
			Err:   errors.New("failed to load notes"),
		})
		m := updatedModel.(app.Model)

		assert.Equal("failed to load notes", m.LastError(), "lastError should match the error message")
		assert.Len(m.Notes(), 0, "notes should remain empty when load fails")
	})
}

func TestUpdate_SearchMode(t *testing.T) {
	assert := testutil.New(t)
	mock := &mockFileSystem{
		searchResults: []*storage.Note{
			{ID: "1", Title: "Go note", Content: "content"},
		},
	}
	model := newTestModel(mock)

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m := updated.(app.Model)
	assert.Equal(app.ModeSearch, m.Mode(), "should enter search mode")

	for _, r := range "go" {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = updated.(app.Model)
	}

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.NotNil(cmd, "enter should trigger search command")

	msg := cmd()
	results, ok := msg.(app.SearchResultsMsg)
	if !ok {
		t.Fatalf("expected SearchResultsMsg, got %T", msg)
	}

	updated, _ = m.Update(results)
	m = updated.(app.Model)
	assert.True(m.SearchSubmitted(), "search should be submitted")
	assert.Len(m.Notes(), 1, "should show search results")
}

func TestUpdate_SearchResultsMsg(t *testing.T) {
	assert := testutil.New(t)
	model := app.NewModelWithStorage(&mockFileSystem{}, "")

	results := []*storage.Note{{ID: "1", Title: "Match", Content: "data"}}
	updated, _ := model.Update(app.SearchResultsMsg{Notes: results, Query: "match"})
	m := updated.(app.Model)

	assert.Len(m.Notes(), 1, "should store search results")
	assert.Empty(m.LastError(), "should clear errors")
}

func TestUpdate_SortCycle(t *testing.T) {
	assert := testutil.New(t)
	model := newTestModel(&mockFileSystem{})

	initial := model.SortMode()
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m := updated.(app.Model)
	assert.NotEqual(initial, m.SortMode(), "sort mode should cycle")
}

func TestUpdate_DeleteConfirmation(t *testing.T) {
	assert := testutil.New(t)
	mock := &mockFileSystem{}
	model := newTestModel(mock)

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m := updated.(app.Model)

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	assert.NotNil(cmd, "second delete should trigger delete command")

	msg := cmd()
	deleted, ok := msg.(app.NoteDeletedMsg)
	if !ok {
		t.Fatalf("expected NoteDeletedMsg, got %T", msg)
	}

	updated, cmd = m.Update(deleted)
	assert.NotNil(cmd, "successful delete should reload notes")
	_ = updated
	assert.Equal(1, mock.deleteCalled, "delete should be called once")
}

func TestUpdate_NoteSavedError(t *testing.T) {
	assert := testutil.New(t)
	model := app.NewModelWithStorage(&mockFileSystem{}, "")

	updated, _ := model.Update(app.NoteSavedMsg{
		Note: &storage.Note{ID: "1"},
		Err:  errors.New("save failed"),
	})
	m := updated.(app.Model)

	assert.Equal("save failed", m.LastError(), "should store save error")
}

func TestUpdate_WindowResize(t *testing.T) {
	assert := testutil.New(t)
	model := app.NewModelWithStorage(&mockFileSystem{}, "")

	updated, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m := updated.(app.Model)

	assert.Equal(120, m.Width(), "width should update")
	assert.Equal(40, m.Height(), "height should update")
}

func TestUpdate_SearchEscReloadsNotes(t *testing.T) {
	assert := testutil.New(t)
	mock := &mockFileSystem{}
	model := newTestModel(mock)

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m := updated.(app.Model)

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	assert.NotNil(cmd, "esc from search should reload notes")
	_ = updated.(app.Model)
	assert.Equal(app.ModeList, updated.(app.Model).Mode(), "should return to list mode")
}
