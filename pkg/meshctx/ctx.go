package meshctx

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type tenantKeyType string

const (
	tenantKey     tenantKeyType = "tenant_id"
	streamPathKey tenantKeyType = "stream_path"
)

// ErrTenantIDMissing is returned when a tenant ID is absent from context.
var ErrTenantIDMissing = errors.New("tenant ID not found in context")

// WithTenantID adds a tenant UUID to the context.
func WithTenantID(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, tenantKey, tenantID)
}

// TenantIDFromContext retrieves the tenant UUID from context.
// Returns uuid.Nil if the value is not found or invalid.
func TenantIDFromContext(ctx context.Context) (uuid.UUID, error) {
	tenantID, ok := ctx.Value(tenantKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrTenantIDMissing
	}
	return tenantID, nil
}

// WithStreamPath stores a stream path in the context. The stored value can be
// retrieved with StreamPathFromContext.
func WithStreamPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, streamPathKey, path)
}

// StreamPathFromContext returns the stream path previously stored in the
// context. The second return value indicates whether a value was present.
func StreamPathFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(streamPathKey).(string)
	return v, ok
}
