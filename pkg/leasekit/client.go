package leasekit

import (
	"context"
	"errors"
	"time"

	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
)

var (
	// Empty is a zero-value Lease used for comparison and tests.
	Empty Lease
	// ErrInvalidTTL indicates an invalid TTL was provided to a lease operation.
	ErrInvalidTTL = errors.New("invalid TTL")
	// ErrNilLease indicates a nil lease pointer was passed to a client method.
	ErrNilLease = errors.New("invalid lease: nil")
	// ErrLeaseExpired indicates the target lease has already expired and
	// cannot be renewed.
	ErrLeaseExpired = errors.New("cannot renew: lease already expired")
)

// Client is a convenience wrapper for interacting with the lease Manager
// over the messaging bus. It provides synchronous helper methods for
// acquiring, renewing and releasing leases.
type Client interface {
	Acquire(ctx context.Context, tenantID uuid.UUID, key string, ttl time.Duration, maxAttempts int) (*Lease, error)
	Renew(ctx context.Context, lease *Lease) error
	Release(ctx context.Context, lease *Lease) error
}

// NewLeaseClient creates a new Client that communicates with the lease
// Manager using the provided messaging.MessageBusFactory.
func NewLeaseClient(busFactory messaging.MessageBusFactory) Client {
	return &client{busFactory: busFactory}
}

type client struct {
	busFactory messaging.MessageBusFactory
}

func (l *client) Acquire(
	ctx context.Context,
	tenantID uuid.UUID,
	key string,
	ttl time.Duration,
	maxAttempts int,
) (*Lease, error) {
	if key == "" || ttl <= 0 {
		return nil, ErrInvalidTTL
	}

	bus, err := l.busFactory.Get(ctx)
	if err != nil {
		return nil, err
	}

	req := &Acquire{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Key:         key,
		TTL:         ttl,
		MaxAttempts: maxAttempts,
	}

	return messaging.SendRequestWithContext[*Acquire, *Lease](ctx, bus, req, 5*time.Second)
}

func (l *client) Renew(ctx context.Context, lease *Lease) error {
	if lease == nil {
		return ErrNilLease
	}
	if lease.ExpireAt.Before(time.Now()) {
		return ErrLeaseExpired
	}
	if lease.TTL <= 0 {
		return ErrInvalidTTL
	}

	bus, err := l.busFactory.Get(ctx)
	if err != nil {
		return err
	}

	req := &Renew{
		ID:       lease.ID,
		TenantID: lease.TenantID,
		Key:      lease.Key,
		TTL:      lease.TTL,
	}

	if _, err := messaging.SendRequestWithContext[*Renew, *messaging.Accepted](ctx, bus, req, 5*time.Second); err != nil {
		return err
	}

	lease.ExpireAt = time.Now().Add(lease.TTL)
	return nil
}

func (l *client) Release(ctx context.Context, lease *Lease) error {
	if lease == nil {
		return ErrNilLease
	}

	bus, err := l.busFactory.Get(ctx)
	if err != nil {
		return err
	}

	req := &Release{
		ID:       lease.ID,
		TenantID: lease.TenantID,
		Key:      lease.Key,
	}

	_, err = messaging.SendRequestWithContext[*Release, *messaging.Accepted](ctx, bus, req, 5*time.Second)
	return err
}
