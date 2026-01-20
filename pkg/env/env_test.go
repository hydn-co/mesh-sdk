package env

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEnvCleanup(t *testing.T, key, val string) {
	t.Helper()
	old, ok := os.LookupEnv(key)
	require.NoError(t, os.Setenv(key, val))
	t.Cleanup(func() {
		if ok {
			_ = os.Setenv(key, old)
		} else {
			_ = os.Unsetenv(key)
		}
	})
}

func unsetEnvCleanup(t *testing.T, key string) {
	t.Helper()
	old, ok := os.LookupEnv(key)
	_ = os.Unsetenv(key)
	t.Cleanup(func() {
		if ok {
			_ = os.Setenv(key, old)
		}
	})
}

func TestShouldReturnStringWhenEnvVarPresent(t *testing.T) {
	// Arrange
	setEnvCleanup(t, "TEST_ENV_STR", "value")

	// Act
	v, ok := TryGetEnvStr("TEST_ENV_STR")
	res := GetEnvOrDefaultStr("TEST_ENV_STR", "default")

	// Assert
	require.True(t, ok)
	assert.Equal(t, "value", v)
	assert.Equal(t, "value", res)
}

func TestShouldReturnDefaultWhenEnvVarMissing(t *testing.T) {
	// Arrange
	unsetEnvCleanup(t, "TEST_ENV_MISSING")

	// Act
	res := GetEnvOrDefaultStr("TEST_ENV_MISSING", "def")
	v, ok := TryGetEnvStr("TEST_ENV_MISSING")

	// Assert
	assert.False(t, ok)
	assert.Equal(t, "", v)
	assert.Equal(t, "def", res)
}

func TestShouldParseIntAndHandleInvalid(t *testing.T) {
	// valid int
	setEnvCleanup(t, "TEST_ENV_INT", "42")
	val, ok := TryGetEnvInt("TEST_ENV_INT")
	assert.True(t, ok)
	assert.Equal(t, 42, val)
	assert.Equal(t, 42, GetEnvOrDefaultInt("TEST_ENV_INT", 7))

	// invalid int
	setEnvCleanup(t, "TEST_ENV_INT", "notint")
	v2, ok2 := TryGetEnvInt("TEST_ENV_INT")
	assert.False(t, ok2)
	assert.Equal(t, 0, v2)
	assert.Equal(t, 7, GetEnvOrDefaultInt("TEST_ENV_INT", 7))
}

func TestShouldParseBool(t *testing.T) {
	setEnvCleanup(t, "TEST_ENV_BOOL", "true")
	b, ok := TryGetEnvBool("TEST_ENV_BOOL")
	assert.True(t, ok)
	assert.True(t, b)

	setEnvCleanup(t, "TEST_ENV_BOOL", "false")
	b2, ok2 := TryGetEnvBool("TEST_ENV_BOOL")
	assert.True(t, ok2)
	assert.False(t, b2)
}

func TestShouldParseUUIDAndMustGet(t *testing.T) {
	id := uuid.New()
	setEnvCleanup(t, "TEST_ENV_UUID", id.String())
	got, ok := TryGetEnvUUID("TEST_ENV_UUID")
	assert.True(t, ok)
	assert.Equal(t, id, got)

	// MustGetEnvUUID returns value when present
	req := MustGetEnvUUID("TEST_ENV_UUID")
	assert.Equal(t, id, req)

	// invalid uuid causes TryGetEnvUUID to return false
	setEnvCleanup(t, "TEST_ENV_UUID", "not-a-uuid")
	got2, ok3 := TryGetEnvUUID("TEST_ENV_UUID")
	assert.False(t, ok3)
	assert.Equal(t, uuid.Nil, got2)

	// MustGetEnvUUID should panic on invalid
	setEnvCleanup(t, "TEST_ENV_UUID", "not-a-uuid")
	require.Panics(t, func() { MustGetEnvUUID("TEST_ENV_UUID") })
}

func TestShouldReturnStringSlice(t *testing.T) {
	setEnvCleanup(t, "TEST_ENV_SLICE", "a,b,c")
	arr, ok := TryGetEnvStrSlice("TEST_ENV_SLICE")
	assert.True(t, ok)
	assert.Equal(t, []string{"a", "b", "c"}, arr)
}

func TestShouldReturnIntSlice(t *testing.T) {
	setEnvCleanup(t, "TEST_ENV_INT_SLICE", "1,2,5")
	arr, ok := TryGetEnvIntSlice("TEST_ENV_INT_SLICE")
	require.True(t, ok)
	assert.Equal(t, []int{1, 2, 5}, arr)
}

func TestMustGetPanicsOnMissing(t *testing.T) {
	unsetEnvCleanup(t, "TEST_MUST_INT")
	require.Panics(t, func() { MustGetEnvInt("TEST_MUST_INT") })
	unsetEnvCleanup(t, "TEST_MUST_STR")
	require.Panics(t, func() { MustGetEnvStr("TEST_MUST_STR") })
}
