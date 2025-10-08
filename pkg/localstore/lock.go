package localstore

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// AcquireFileLock creates a simple file-based lock at the provided path.
// If the lock already exists, it will treat locks older than a small
// threshold as stale and remove them. The function returns an error if it
// cannot acquire the lock after multiple attempts.
func AcquireFileLock(path string) error {
	const maxAge = 5 * time.Second

	// try up to 20 attempts with small backoff
	for i := 0; i < 20; i++ {
		info, err := os.Stat(path)
		if err == nil {
			if time.Since(info.ModTime()) > maxAge {
				_ = os.Remove(path) // stale lock
				continue
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0600)
		if err == nil {
			f.Close()
			return nil
		}
		if !errors.Is(err, os.ErrExist) {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("timed out waiting for lock: " + filepath.Base(path))
}
