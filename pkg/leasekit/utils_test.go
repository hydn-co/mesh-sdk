package leasekit

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/tenantkit"
)

func TestWithLease_ContextValidation(t *testing.T) {
	// Test that WithLease properly validates context
	ctx := context.Background()

	err := WithLease(ctx, nil, "test-key", time.Minute, 1, nil)
	if err == nil {
		t.Error("Expected error when context missing tenant ID")
	}
}

func TestWithLease_WithValidContext(t *testing.T) {
	// Test WithLease with a valid tenant context
	tenantID := uuid.New()
	ctx := tenantkit.WithTenantID(context.Background(), tenantID)

	// This would need a mock client for full testing
	client := &mockLeaseClient{}

	executed := false
	err := WithLease(ctx, client, "test-key", time.Minute, 1, func(ctx context.Context) error {
		executed = true
		return nil
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !executed {
		t.Error("Expected function to be executed")
	}
}

// Mock lease client for testing
type mockLeaseClient struct{}

func (m *mockLeaseClient) Acquire(ctx context.Context, tenantID uuid.UUID, key string, ttl time.Duration, maxAttempts int) (*Lease, error) {
	return &Lease{
		ID:       uuid.New(),
		TenantID: tenantID,
		Key:      key,
		TTL:      ttl,
		ExpireAt: time.Now().Add(ttl),
	}, nil
}

func (m *mockLeaseClient) Renew(ctx context.Context, lease *Lease) error {
	return nil
}

func (m *mockLeaseClient) Release(ctx context.Context, lease *Lease) error {
	return nil
}
