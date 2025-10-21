package localstore

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcquireFileLockCreatesLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.lock")
	// Ensure no pre-existing lock
	_, err := os.Stat(path)
	assert.True(t, os.IsNotExist(err) || err == nil)

	// Act
	require.NoError(t, AcquireFileLock(path))

	// Assert the lock file exists
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.False(t, info.IsDir())
}

func TestAcquireFileLockRemovesStaleLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stale.lock")
	// Create a lock file with old mod time
	require.NoError(t, os.WriteFile(path, []byte("lock"), 0600))
	old := time.Now().Add(-10 * time.Second)
	require.NoError(t, os.Chtimes(path, old, old))

	// Act: should remove the stale lock and acquire
	require.NoError(t, AcquireFileLock(path))

	// Assert the lock is present
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.False(t, info.IsDir())
}

func TestAcquireFileLockTimeoutWhenHeld(t *testing.T) {
	path := filepath.Join(t.TempDir(), "held.lock")
	// Create a lock file and keep it fresh (simulate another process holding it)
	require.NoError(t, os.WriteFile(path, []byte("lock"), 0600))
	fresh := time.Now()
	require.NoError(t, os.Chtimes(path, fresh, fresh))

	// Make AcquireFileLock relatively quick by overriding loop? Not possible
	// But we can test by invoking and expect timeout after 20 attempts ~2s.
	// Act
	err := AcquireFileLock(path)

	// Assert: should return timeout error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out waiting for lock")
}
