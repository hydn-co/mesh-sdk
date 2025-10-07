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
	// Arrange
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

	// Act
	lease, err := client.Acquire(context.Background(), tenantID, key, ttl, maxAttempts)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, lease)
	assert.Equal(t, key, lease.Key)
}

func TestLeaseClient_Renew(t *testing.T) {
	// Arrange
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

	// Act
	err := client.Renew(context.Background(), lease)

	// Assert
	assert.NoError(t, err)
}

func TestLeaseClient_Release(t *testing.T) {
	// Arrange
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

	// Act
	err := client.Release(context.Background(), lease)

	// Assert
	assert.NoError(t, err)
}
