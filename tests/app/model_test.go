package app_test

import (
	"testing"

	"github.com/N95Ryan/leaf/internal/app"
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
