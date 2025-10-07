package leasekit

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/tenantkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithLease_ContextValidation(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	err := WithLease(ctx, nil, "test-key", time.Minute, 1, nil)

	// Assert
	assert.Error(t, err)
}

func TestWithLease_WithValidContext(t *testing.T) {
	// Arrange
	tenantID := uuid.New()
	ctx := tenantkit.WithTenantID(context.Background(), tenantID)
	client := &mockLeaseClient{}

	// Act
	executed := false
	err := WithLease(ctx, client, "test-key", time.Minute, 1, func(ctx context.Context) error {
		executed = true
		return nil
	})

	// Assert
	require.NoError(t, err)
	assert.True(t, executed)
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
