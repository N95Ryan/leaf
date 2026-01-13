package app_test

import (
	"testing"

	"github.com/N95Ryan/leaf/internal/app"
	"github.com/N95Ryan/leaf/tests/testutil"
)

// TestDemo_FailureExamples demonstrates how emoji assertions look when they fail
// This file is for demonstration purposes only
// Rename this file to model_test_demo.go.disabled to disable these tests
func TestDemo_FailureExamples(t *testing.T) {
	t.Skip("⚠️  DEMO: Uncomment this line to see emoji failures in action")
	
	assert := testutil.New(t)
	model := app.NewModel()
	
	// These will fail to show the emoji output
	t.Run("Equal failure example", func(t *testing.T) {
		assert := testutil.New(t)
		assert.Equal(app.ModeEdit, model.Mode(), "this will fail to show ❌ Equal")
	})
	
	t.Run("Len failure example", func(t *testing.T) {
		assert := testutil.New(t)
		assert.Len(model.Notes(), 5, "this will fail to show ❌ Len with 📏")
	})
	
	t.Run("True failure example", func(t *testing.T) {
		assert := testutil.New(t)
		assert.True(false, "this will fail to show ❌ True with ✓/✗")
	})
	
	t.Run("Empty failure example", func(t *testing.T) {
		assert := testutil.New(t)
		assert.Empty("not empty", "this will fail to show ❌ Empty with 📭")
	})
}
