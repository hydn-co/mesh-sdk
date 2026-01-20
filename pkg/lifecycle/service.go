package lifecycle

import (
	"context"
)

// Service represents a managed, long-running component.
//
// Conventions:
//   - Start(ctx) should be non-blocking: it should begin necessary background
//     activity (e.g., launch goroutines) and return quickly. Any initialization
//     that may block should be performed with respect to the provided context
//     and kept minimal.
//   - Stop(ctx) is used to cancel or terminate work started by Start. Stop is
//     allowed to block until shutdown completes, and should respect the provided
//     context for timeouts/cancellation.
//
// This convention makes it easy for Hosts to start many services quickly and
// later stop them deterministically.
type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}
