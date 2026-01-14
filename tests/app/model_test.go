package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/N95Ryan/leaf/internal/app"
	"github.com/N95Ryan/leaf/internal/storage"
	"github.com/N95Ryan/leaf/tests/testutil"
)

func TestNewModel(t *testing.T) {
	t.Run("should initialize model successfully", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModel()

		// Check that the model is initialized with correct defaults
		assert.Equal(app.ModeList, model.Mode(), "mode should be ModeList by default")

		// Check that notes are initialized as empty slice
		assert.NotNil(model.Notes(), "notes should be initialized")
		assert.Empty(model.Notes(), "notes should be empty initially")

		// Check that storage is initialized
		assert.NotNil(model.Storage(), "storage should be initialized")

		// Check that there is no error on successful initialization
		assert.Empty(model.LastError(), "should have no error on successful initialization")
	})
}

func TestNewModel_StorageInitialization(t *testing.T) {
	t.Run("should have storage initialized or error set", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModel()

		// Either storage should be initialized OR lastError should be set
		// (both being nil/empty would indicate a problem)
		if model.Storage() == nil {
			assert.NotEmpty(model.LastError(),
				"if storage initialization fails, lastError should contain error message")
		} else {
			assert.Empty(model.LastError(),
				"if storage initialization succeeds, lastError should be empty")
		}
	})
}

func TestNewModel_DefaultValues(t *testing.T) {
	t.Run("should have correct default values", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModel()

		// Verify all expected default values
		assert.NotNil(model.Notes(), "notes should not be nil")
		assert.Len(model.Notes(), 0, "notes should start empty")
		assert.Equal(app.ModeList, model.Mode(), "should start in ModeList")
	})
}

// Mock FileSystem for testing
type mockFileSystem struct {
	notes []*storage.Note
	err   error
}

func (m *mockFileSystem) ListNotes(ctx context.Context) ([]*storage.Note, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.notes, nil
}

func (m *mockFileSystem) GetNote(ctx context.Context, id string) (*storage.Note, error) {
	return nil, nil
}

func (m *mockFileSystem) SaveNote(ctx context.Context, note *storage.Note) error {
	return nil
}

func (m *mockFileSystem) DeleteNote(ctx context.Context, id string) error {
	return nil
}

func (m *mockFileSystem) SearchNotes(ctx context.Context, query string) ([]*storage.Note, error) {
	return nil, nil
}

func TestInit(t *testing.T) {
	t.Run("should return nil when storage is nil", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModel()

		// Force storage to nil to test edge case
		// We can't directly set it, but if initialization failed, storage would be nil
		// We test the Init behavior indirectly
		cmd := model.Init()

		// If storage is nil, Init should return nil
		// Otherwise, it should return a command
		if model.Storage() == nil {
			assert.Nil(cmd, "Init should return nil when storage is nil")
		} else {
			assert.NotNil(cmd, "Init should return a command when storage is initialized")
		}
	})

	t.Run("should return a valid command when storage is initialized", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModel()

		// Only test if storage was successfully initialized
		if model.Storage() != nil {
			cmd := model.Init()
			assert.NotNil(cmd, "Init should return a command to load notes")
		}
	})
}

func TestUpdate_NotesLoadedMsg(t *testing.T) {
	t.Run("should store notes when loaded successfully", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModel()

		// Create test notes
		testNotes := []*storage.Note{
			{ID: "1", Title: "Note 1", Content: "Content 1"},
			{ID: "2", Title: "Note 2", Content: "Content 2"},
		}

		// Simulate receiving NoteLoadedMsg with success
		msg := app.NoteLoadedMsg{
			Notes: testNotes,
			Err:   nil,
		}

		updatedModel, _ := model.Update(msg)
		m := updatedModel.(app.Model)

		// Verify notes are stored
		assert.Len(m.Notes(), 2, "should have 2 notes loaded")
		assert.Equal("Note 1", m.Notes()[0].Title, "first note should have correct title")
		assert.Equal("Note 2", m.Notes()[1].Title, "second note should have correct title")

		// Verify error is cleared
		assert.Empty(m.LastError(), "lastError should be cleared on successful load")
	})

	t.Run("should store error when loading fails", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModel()

		// Simulate receiving NoteLoadedMsg with error
		testErr := errors.New("failed to load notes")
		msg := app.NoteLoadedMsg{
			Notes: nil,
			Err:   testErr,
		}

		updatedModel, _ := model.Update(msg)
		m := updatedModel.(app.Model)

		// Verify error is stored
		assert.NotEmpty(m.LastError(), "lastError should contain error message")
		assert.Equal("failed to load notes", m.LastError(), "lastError should match the error message")

		// Verify notes remain unchanged (should still be empty)
		assert.Len(m.Notes(), 0, "notes should remain empty when load fails")
	})

	t.Run("should clear previous error on successful load", func(t *testing.T) {
		assert := testutil.New(t)
		model := app.NewModel()

		// First, simulate an error
		errMsg := app.NoteLoadedMsg{
			Notes: nil,
			Err:   errors.New("previous error"),
		}
		updatedModel, _ := model.Update(errMsg)
		model = updatedModel.(app.Model)

		// Verify error is set
		assert.NotEmpty(model.LastError(), "lastError should be set")

		// Now simulate a successful load
		testNotes := []*storage.Note{
			{ID: "1", Title: "Note 1", Content: "Content 1"},
		}
		successMsg := app.NoteLoadedMsg{
			Notes: testNotes,
			Err:   nil,
		}
		updatedModel, _ = model.Update(successMsg)
		model = updatedModel.(app.Model)

		// Verify error is cleared
		assert.Empty(model.LastError(), "lastError should be cleared on successful load")
		assert.Len(model.Notes(), 1, "notes should be loaded")
	})
}
