package tenant

import (
	"context"

	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/meshctx"
)

// Backwards-compatibility shim: delegate tenant API to the new meshctx package.
var ErrTenantIDMissing = meshctx.ErrTenantIDMissing

func WithTenantID(ctx context.Context, tenantID uuid.UUID) context.Context {
	return meshctx.WithTenantID(ctx, tenantID)
}

func TenantIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return meshctx.TenantIDFromContext(ctx)
}

func WithStreamPath(ctx context.Context, path string) context.Context {
	return meshctx.WithStreamPath(ctx, path)
}

func StreamPathFromContext(ctx context.Context) (string, bool) {
	return meshctx.StreamPathFromContext(ctx)
}
