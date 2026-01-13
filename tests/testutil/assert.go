package testutil

import (
	"fmt"
	"reflect"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Helper()
	Errorf(format string, args ...interface{})
}

// Assert provides custom assertion helpers for tests
type Assert struct {
	t TestingT
}

// New creates a new Assert instance
func New(t TestingT) *Assert {
	return &Assert{t: t}
}

// Equal checks if two values are equal
func (a *Assert) Equal(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ Equal failed\n  Expected: %v\n  Actual:   %v\n  %s", expected, actual, msg)
	}
}

// NotEqual checks if two values are not equal
func (a *Assert) NotEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if reflect.DeepEqual(expected, actual) {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ NotEqual failed\n  Expected values to be different\n  Got: %v\n  %s", actual, msg)
	}
}

// Nil checks if a value is nil
func (a *Assert) Nil(value interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if !isNil(value) {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ Nil failed\n  Expected: nil\n  Got:      %v\n  %s", value, msg)
	}
}

// NotNil checks if a value is not nil
func (a *Assert) NotNil(value interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if isNil(value) {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ NotNil failed\n  Expected: non-nil value\n  Got:      nil\n  %s", msg)
	}
}

// True checks if a value is true
func (a *Assert) True(value bool, msgAndArgs ...interface{}) {
	a.t.Helper()
	if !value {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ True failed\n  Expected: true ✓\n  Got:      false ✗\n  %s", msg)
	}
}

// False checks if a value is false
func (a *Assert) False(value bool, msgAndArgs ...interface{}) {
	a.t.Helper()
	if value {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ False failed\n  Expected: false ✗\n  Got:      true ✓\n  %s", msg)
	}
}

// Empty checks if a value is empty (nil, empty string, empty slice, etc.)
func (a *Assert) Empty(value interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if !isEmpty(value) {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ Empty failed\n  Expected: empty 📭\n  Got:      %v\n  %s", value, msg)
	}
}

// NotEmpty checks if a value is not empty
func (a *Assert) NotEmpty(value interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if isEmpty(value) {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ NotEmpty failed\n  Expected: non-empty value 📬\n  Got:      empty\n  %s", msg)
	}
}

// Len checks if a slice/map/array has the expected length
func (a *Assert) Len(value interface{}, expectedLen int, msgAndArgs ...interface{}) {
	a.t.Helper()
	actualLen := getLen(value)
	if actualLen != expectedLen {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ Len failed\n  Expected length: %d 📏\n  Actual length:   %d\n  %s", expectedLen, actualLen, msg)
	}
}

// NoError checks if an error is nil
func (a *Assert) NoError(err error, msgAndArgs ...interface{}) {
	a.t.Helper()
	if err != nil {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ NoError failed\n  Expected: no error ✓\n  Got:      %v ⚠️\n  %s", err, msg)
	}
}

// Error checks if an error is not nil
func (a *Assert) Error(err error, msgAndArgs ...interface{}) {
	a.t.Helper()
	if err == nil {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ Error failed\n  Expected: an error ⚠️\n  Got:      nil (no error)\n  %s", msg)
	}
}

// Contains checks if a string contains a substring
func (a *Assert) Contains(haystack, needle string, msgAndArgs ...interface{}) {
	a.t.Helper()
	if !contains(haystack, needle) {
		msg := formatMessage(msgAndArgs...)
		a.t.Errorf("❌ Contains failed\n  Expected: %q to contain %q 🔍\n  %s", haystack, needle, msg)
	}
}

// Helper functions

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	kind := v.Kind()
	
	// Check if the type can be nil
	switch kind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	
	return false
}

func isEmpty(value interface{}) bool {
	if isNil(value) {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return isEmpty(v.Elem().Interface())
	default:
		return false
	}
}

func getLen(value interface{}) int {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len()
	default:
		return 0
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && 
		(haystack == needle || len(needle) == 0 || indexOf(haystack, needle) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func formatMessage(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 {
		return ""
	}
	
	if len(msgAndArgs) == 1 {
		return fmt.Sprintf("%v", msgAndArgs[0])
	}
	
	format, ok := msgAndArgs[0].(string)
	if !ok {
		return fmt.Sprint(msgAndArgs...)
	}
	
	return fmt.Sprintf(format, msgAndArgs[1:]...)
}
