package lifecycle

import (
	"context"
)

// Service represents a long-running component managed by a Host. Implementers
// should ensure Start begins background activity and Stop gracefully terminates it.
type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}
