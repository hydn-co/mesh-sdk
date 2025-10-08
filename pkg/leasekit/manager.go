package leasekit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/lifecycle"
)

// Manager is a lifecycle.Service that manages distributed leases. It
// registers request handlers on the message bus to respond to Acquire,
// Renew and Release requests.
type Manager struct {
	messaging.Processor
	mu     sync.Mutex
	busf   messaging.MessageBusFactory
	leases map[uuid.UUID]map[string]Lease // tenantID -> key -> lease
}

// NewManager constructs a lease Manager that will use the provided
// messaging.MessageBusFactory to register handlers when started.
func NewManager(busf messaging.MessageBusFactory) lifecycle.Service {
	return &Manager{
		busf:   busf,
		leases: make(map[uuid.UUID]map[string]Lease),
	}
}

func (m *Manager) Start(ctx context.Context) error {

	bus, err := m.busf.Get(ctx)
	if err != nil {
		return err
	}
	m.Processor = messaging.NewProcessor(bus)

	if err := messaging.RegisterRequestHandler(m, GetAcquireRoute(nil), m.acquire); err != nil {
		return err
	}

	if err := messaging.RegisterRequestHandler(m, GetRenewRoute(nil), m.renew); err != nil {
		return err
	}

	if err := messaging.RegisterRequestHandler(m, GetReleaseRoute(nil), m.release); err != nil {
		return err
	}

	return nil
}

// Acquire attempts to acquire a lease, retrying with exponential backoff if needed.
func (m *Manager) acquire(ctx context.Context, msg *Acquire) (*Lease, error) {
	if msg.MaxAttempts <= 0 {
		msg.MaxAttempts = 1
	}

	var attempt int
	for {
		m.mu.Lock()
		now := time.Now()

		if _, ok := m.leases[msg.TenantID]; !ok {
			m.leases[msg.TenantID] = make(map[string]Lease)
		}

		existing, exists := m.leases[msg.TenantID][msg.Key]
		if !exists || existing.ExpireAt.Before(now) {
			lease := Lease{
				ID:       msg.ID,
				TenantID: msg.TenantID,
				Key:      msg.Key,
				TTL:      msg.TTL,
				ExpireAt: now.Add(msg.TTL),
			}
			m.leases[msg.TenantID][msg.Key] = lease
			m.mu.Unlock()
			return &lease, nil
		}
		m.mu.Unlock()

		attempt++
		if attempt >= msg.MaxAttempts {
			return nil, fmt.Errorf("lease already held for key: %s (after %d attempts)", msg.Key, attempt)
		}

		backoff := time.Duration(100*(1<<uint(attempt-1))) * time.Millisecond
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
			continue
		}
	}
}

func (m *Manager) renew(ctx context.Context, msg *Renew) (*messaging.Accepted, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tenantLeases, ok := m.leases[msg.TenantID]
	if !ok {
		return nil, fmt.Errorf("no leases found for tenant %s", msg.TenantID)
	}

	lease, ok := tenantLeases[msg.Key]
	if !ok {
		return nil, fmt.Errorf("no lease for key: %s", msg.Key)
	}

	if lease.ID != msg.ID {
		return nil, fmt.Errorf("lease ID mismatch")
	}

	if lease.ExpireAt.Before(time.Now()) {
		return nil, fmt.Errorf("cannot renew expired lease")
	}

	lease.ExpireAt = time.Now().Add(msg.TTL)
	tenantLeases[msg.Key] = lease
	return &messaging.Accepted{}, nil
}

func (m *Manager) release(ctx context.Context, msg *Release) (*messaging.Accepted, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tenantLeases, ok := m.leases[msg.TenantID]
	if !ok {
		return &messaging.Accepted{}, nil
	}

	lease, ok := tenantLeases[msg.Key]
	if !ok {
		return &messaging.Accepted{}, nil
	}

	if lease.ID == msg.ID {
		delete(tenantLeases, msg.Key)
	}
	return &messaging.Accepted{}, nil
}
