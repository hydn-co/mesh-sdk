package leasekit

import (
	"context"
	"testing"
	"time"

	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLeaseClient_Acquire(t *testing.T) {
	bus := testkit.NewMockMessageBus()
	bus.Mock.On("RequestWithContext",
		mock.Anything,
		mock.MatchedBy(func(r messaging.Request) bool {
			_, ok := r.(*Acquire)
			return ok
		}),
		mock.Anything,
	).Return(&Lease{ID: uuid.New(), TenantID: uuid.New(), Key: "test-key", TTL: time.Second, ExpireAt: time.Now().Add(time.Second)}, nil).Once()

	busf := testkit.ConfigureBusFactory(bus)
	client := NewLeaseClient(busf)
	tenantID := uuid.New()
	key := "test-key"
	ttl := time.Second
	maxAttempts := 1

	lease, err := client.Acquire(context.Background(), tenantID, key, ttl, maxAttempts)
	assert.NoError(t, err)
	assert.NotNil(t, lease, "Expected lease to be non-nil")
	assert.Equal(t, key, lease.Key, "Expected lease with key %q, got %+v", key, lease)
}

func TestLeaseClient_Renew(t *testing.T) {
	bus := testkit.NewMockMessageBus()
	bus.Mock.On("RequestWithContext",
		mock.Anything,
		mock.MatchedBy(func(r messaging.Request) bool {
			_, ok := r.(*Renew)
			return ok
		}),
		mock.Anything,
	).Return(&messaging.Accepted{}, nil).Once()

	busf := testkit.ConfigureBusFactory(bus)
	client := NewLeaseClient(busf)
	// Set ExpireAt to a future time to ensure lease is not expired
	lease := &Lease{ID: uuid.New(), Key: "test-key", TTL: time.Second, ExpireAt: time.Now().Add(time.Minute)}
	err := client.Renew(context.Background(), lease)
	assert.NoError(t, err, "Renew returned error")
}

func TestLeaseClient_Release(t *testing.T) {
	bus := testkit.NewMockMessageBus()
	bus.Mock.On("RequestWithContext",
		mock.Anything,
		mock.MatchedBy(func(r messaging.Request) bool {
			_, ok := r.(*Release)
			return ok
		}),
		mock.Anything,
	).Return(&messaging.Accepted{}, nil).Once()

	busf := testkit.ConfigureBusFactory(bus)
	client := NewLeaseClient(busf)
	lease := &Lease{ID: uuid.New(), Key: "test-key"}
	err := client.Release(context.Background(), lease)
	assert.NoError(t, err, "Release returned error")
}
