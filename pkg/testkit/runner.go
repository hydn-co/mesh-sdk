package testkit

import (
	"context"
)

// Manifest is a very small, local representation used by tests.
// The real project may have a richer structure; this shim only implements
// what's needed by `internal/testkit/invoke.go`.
type Manifest struct {
	Capabilities map[string]any
}

// RunWithArgs is a no-op shim used by tests. It accepts args similar to
// the real runner but performs no real work. Returning nil allows tests
// that just assert successful execution to continue.
func RunWithArgs(ctx context.Context, manifest *Manifest, args ...string) error {
	// no-op for tests
	return nil
}
