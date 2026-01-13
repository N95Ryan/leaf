# 🧪 TestUtil - Custom Testing Library for Leaf

> Custom assertion library based on Go's standard `testing` package

---

## 📋 Overview

This library provides simple and expressive assertions for Go tests, without external dependencies. It is specifically designed for the Leaf project.

**✨ New:** All error messages include emojis for quick reading! See [EMOJIS.md](./EMOJIS.md) for the complete guide.

## 🚀 Usage

### Import

```go
import (
    "testing"
    "github.com/N95Ryan/leaf/tests/testutil"
)
```

### Basic Example

```go
func TestExample(t *testing.T) {
    assert := testutil.New(t)

    // Equality test
    assert.Equal(expected, actual, "values should be equal")

    // Nil test
    assert.NotNil(value, "value should not be nil")

    // Empty slice test
    assert.Empty(mySlice, "slice should be empty")
}
```

## 📚 API

### Comparison Assertions

#### `Equal(expected, actual interface{}, msgAndArgs ...interface{})`

Checks that two values are equal (uses `reflect.DeepEqual`).

```go
assert.Equal(42, actualValue, "should be 42")
assert.Equal("hello", actualString, "message: %s", "test")
```

#### `NotEqual(expected, actual interface{}, msgAndArgs ...interface{})`

Checks that two values are different.

```go
assert.NotEqual(0, count, "count should not be zero")
```

### Nil Assertions

#### `Nil(value interface{}, msgAndArgs ...interface{})`

Checks that a value is nil.

```go
assert.Nil(err, "error should be nil")
```

#### `NotNil(value interface{}, msgAndArgs ...interface{})`

Checks that a value is not nil.

```go
assert.NotNil(result, "result should not be nil")
```

### Boolean Assertions

#### `True(value bool, msgAndArgs ...interface{})`

Checks that a value is true.

```go
assert.True(isValid, "should be valid")
```

#### `False(value bool, msgAndArgs ...interface{})`

Checks that a value is false.

```go
assert.False(hasError, "should not have error")
```

### Empty Assertions

#### `Empty(value interface{}, msgAndArgs ...interface{})`

Checks that a value is empty (nil, "", empty slice, empty map, etc.).

```go
assert.Empty(mySlice, "slice should be empty")
assert.Empty("", "string should be empty")
```

#### `NotEmpty(value interface{}, msgAndArgs ...interface{})`

Checks that a value is not empty.

```go
assert.NotEmpty(notes, "should have notes")
```

### Length Assertions

#### `Len(value interface{}, expectedLen int, msgAndArgs ...interface{})`

Checks the length of a slice, map, array or string.

```go
assert.Len(mySlice, 5, "should have 5 elements")
assert.Len("hello", 5, "string should have 5 characters")
```

### Error Assertions

#### `NoError(err error, msgAndArgs ...interface{})`

Checks that an error is nil.

```go
assert.NoError(err, "should not return error")
```

#### `Error(err error, msgAndArgs ...interface{})`

Checks that an error is not nil.

```go
assert.Error(err, "should return error")
```

### String Assertions

#### `Contains(haystack, needle string, msgAndArgs ...interface{})`

Checks that a string contains a substring.

```go
assert.Contains(message, "error", "message should contain 'error'")
```

## 🎯 Message Formatting

All messages can be formatted with `fmt.Sprintf`:

```go
// Simple message
assert.Equal(expected, actual, "should match")

// Formatted message
assert.Equal(expected, actual, "expected %d, got %d", expected, actual)

// Multiple arguments
assert.NotNil(value, "value should be set for user %s", userName)
```

## ✨ Advantages

- ✅ **No external dependencies**: uses only the standard `testing` package
- ✅ **Full control**: you can modify and extend the library as needed
- ✅ **Clear messages**: errors are formatted in a readable way with emojis 🎨
- ✅ **Helper marks**: uses `t.Helper()` to display the correct error line
- ✅ **Flexible**: supports formatted messages with `fmt.Sprintf`
- ✅ **Quick reading**: instant visual identification with emojis (❌, ✓, ⚠️, 📏, etc.)

## 🔧 Extension

To add new assertions, simply modify `tests/testutil/assert.go`:

```go
// Example of custom assertion
func (a *Assert) Between(value, min, max int, msgAndArgs ...interface{}) {
    a.t.Helper()
    if value < min || value > max {
        msg := formatMessage(msgAndArgs...)
        a.t.Errorf("Expected %d to be between %d and %d\n%s", value, min, max, msg)
    }
}
```

## 📝 Complete Examples

### Simple Test

```go
func TestUser(t *testing.T) {
    assert := testutil.New(t)

    user := NewUser("John", "Doe")

    assert.NotNil(user, "user should be created")
    assert.Equal("John", user.FirstName, "first name should match")
    assert.Equal("Doe", user.LastName, "last name should match")
}
```

### Test with Subtests

```go
func TestCalculator(t *testing.T) {
    t.Run("addition", func(t *testing.T) {
        assert := testutil.New(t)
        result := Add(2, 3)
        assert.Equal(5, result, "2 + 3 should equal 5")
    })

    t.Run("division", func(t *testing.T) {
        assert := testutil.New(t)
        result, err := Divide(10, 2)
        assert.NoError(err, "should not error on valid division")
        assert.Equal(5, result, "10 / 2 should equal 5")
    })
}
```

---

**Created for the Leaf project** 🌱
