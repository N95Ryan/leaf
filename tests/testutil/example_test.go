package testutil_test

import (
	"testing"

	"github.com/N95Ryan/leaf/tests/testutil"
)

// ExampleAssert demonstrates the visual output of various assertions
func TestExampleAssertions(t *testing.T) {
	// This test is designed to show how assertions look when they pass
	// Run with: go test ./tests/testutil/... -v -run TestExampleAssertions
	
	t.Run("✓ successful assertions", func(t *testing.T) {
		assert := testutil.New(t)
		
		// These all pass
		assert.Equal(42, 42, "numbers are equal")
		assert.NotNil("not nil", "string is not nil")
		assert.True(true, "boolean is true")
		assert.Empty("", "string is empty")
		assert.Len([]int{1, 2, 3}, 3, "slice has 3 elements")
		assert.NoError(nil, "no error occurred")
		assert.Contains("hello world", "world", "string contains substring")
	})
	
	// Uncomment the tests below to see how failures look with emojis
	// These are commented out so the test suite passes by default
	
	/*
	t.Run("❌ failed Equal", func(t *testing.T) {
		assert := testutil.New(t)
		assert.Equal(42, 43, "these numbers don't match")
	})
	
	t.Run("❌ failed NotNil", func(t *testing.T) {
		assert := testutil.New(t)
		var value *string
		assert.NotNil(value, "this value is nil")
	})
	
	t.Run("❌ failed True", func(t *testing.T) {
		assert := testutil.New(t)
		assert.True(false, "this is false, not true")
	})
	
	t.Run("❌ failed Empty", func(t *testing.T) {
		assert := testutil.New(t)
		assert.Empty("not empty", "this string is not empty")
	})
	
	t.Run("❌ failed Len", func(t *testing.T) {
		assert := testutil.New(t)
		assert.Len([]int{1, 2}, 5, "slice has wrong length")
	})
	
	t.Run("❌ failed NoError", func(t *testing.T) {
		assert := testutil.New(t)
		err := errors.New("something went wrong")
		assert.NoError(err, "expected no error but got one")
	})
	
	t.Run("❌ failed Contains", func(t *testing.T) {
		assert := testutil.New(t)
		assert.Contains("hello", "goodbye", "substring not found")
	})
	*/
}

// This example shows best practices for using testutil
func TestBestPractices(t *testing.T) {
	t.Run("use descriptive messages", func(t *testing.T) {
		assert := testutil.New(t)
		
		user := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		
		// Good: descriptive message
		assert.NotNil(user, "user should be created from database")
		assert.Equal("John", user["name"], "user name should match database value")
		assert.Equal(30, user["age"], "user age should be 30")
	})
	
	t.Run("use formatted messages", func(t *testing.T) {
		assert := testutil.New(t)
		
		for i := 0; i < 3; i++ {
			value := i * 2
			expected := i * 2
			
			// Using formatted messages with context
			assert.Equal(expected, value, "iteration %d: value should be %d", i, expected)
		}
	})
	
	t.Run("group related assertions", func(t *testing.T) {
		assert := testutil.New(t)
		
		// Test related properties together
		slice := []int{1, 2, 3}
		
		assert.NotNil(slice, "slice should be initialized")
		assert.Len(slice, 3, "slice should have 3 elements")
		assert.NotEmpty(slice, "slice should not be empty")
	})
}
