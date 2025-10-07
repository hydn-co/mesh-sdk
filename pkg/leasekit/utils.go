package leasekit

import (
	"context"
	"log/slog"
	"time"

	"github.com/hydn-co/mesh-sdk/pkg/tenantkit"
)

// WithLease acquires a lease and runs fn with a context that is canceled when fn completes or the lease is lost.
func WithLease(ctx context.Context, client Client, key string, ttl time.Duration, maxAttempts int, fn func(context.Context) error) error {

	tenantID, err := tenantkit.TenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lease, err := client.Acquire(ctx, tenantID, key, ttl, maxAttempts)
	if err != nil {
		return err
	}

	defer func() {
		if err := client.Release(context.Background(), lease); err != nil {
			slog.Error("failed to release lease", "key", key, "error", err)
		}
	}()

	done := make(chan error, 1)
	go func() {
		done <- fn(ctx)
	}()

	renewTicker := time.NewTicker(ttl / 3) // Renew every 1/3 of TTL
	defer renewTicker.Stop()

	for {
		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return ctx.Err()
		case <-renewTicker.C:
			if err := client.Renew(ctx, lease); err != nil {
				slog.Error("failed to renew lease", "key", key, "error", err)
				cancel() // Cancel the context to signal that the lease is lost
				return err
			}
		}
	}
}
