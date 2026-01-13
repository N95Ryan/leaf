package testutil

import (
	"errors"
	"testing"
)

// mockT is a mock implementation of testing.T for testing our assertions
type mockT struct {
	failed  bool
	message string
}

func (m *mockT) Helper()                              {}
func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.message = format
}

func TestAssert_Equal(t *testing.T) {
	t.Run("should pass when values are equal", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		assert.Equal(42, 42)
		
		if mock.failed {
			t.Error("assertion should have passed")
		}
	})
	
	t.Run("should fail when values are not equal", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		assert.Equal(42, 43)
		
		if !mock.failed {
			t.Error("assertion should have failed")
		}
	})
}

func TestAssert_NotNil(t *testing.T) {
	t.Run("should pass when value is not nil", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		value := "not nil"
		assert.NotNil(value)
		
		if mock.failed {
			t.Error("assertion should have passed")
		}
	})
	
	t.Run("should fail when value is nil", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		var value *string
		assert.NotNil(value)
		
		if !mock.failed {
			t.Error("assertion should have failed")
		}
	})
}

func TestAssert_Empty(t *testing.T) {
	t.Run("should pass for empty string", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		assert.Empty("")
		
		if mock.failed {
			t.Error("assertion should have passed")
		}
	})
	
	t.Run("should pass for empty slice", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		var slice []string
		assert.Empty(slice)
		
		if mock.failed {
			t.Error("assertion should have passed")
		}
	})
	
	t.Run("should fail for non-empty string", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		assert.Empty("not empty")
		
		if !mock.failed {
			t.Error("assertion should have failed")
		}
	})
}

func TestAssert_Len(t *testing.T) {
	t.Run("should pass when length matches", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		slice := []int{1, 2, 3}
		assert.Len(slice, 3)
		
		if mock.failed {
			t.Error("assertion should have passed")
		}
	})
	
	t.Run("should fail when length doesn't match", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		slice := []int{1, 2, 3}
		assert.Len(slice, 5)
		
		if !mock.failed {
			t.Error("assertion should have failed")
		}
	})
}

func TestAssert_NoError(t *testing.T) {
	t.Run("should pass when error is nil", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		assert.NoError(nil)
		
		if mock.failed {
			t.Error("assertion should have passed")
		}
	})
	
	t.Run("should fail when error is not nil", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		err := errors.New("test error")
		assert.NoError(err)
		
		if !mock.failed {
			t.Error("assertion should have failed")
		}
	})
}

func TestAssert_True(t *testing.T) {
	t.Run("should pass when value is true", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		assert.True(true)
		
		if mock.failed {
			t.Error("assertion should have passed")
		}
	})
	
	t.Run("should fail when value is false", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		assert.True(false)
		
		if !mock.failed {
			t.Error("assertion should have failed")
		}
	})
}

func TestAssert_Contains(t *testing.T) {
	t.Run("should pass when string contains substring", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		assert.Contains("hello world", "world")
		
		if mock.failed {
			t.Error("assertion should have passed")
		}
	})
	
	t.Run("should fail when string doesn't contain substring", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)
		
		assert.Contains("hello world", "goodbye")
		
		if !mock.failed {
			t.Error("assertion should have failed")
		}
	})
}
