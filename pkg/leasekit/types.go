package leasekit

import (
	"time"

	"github.com/fgrzl/json/polymorphic"
	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
)

func init() {
	polymorphic.RegisterType[Lease]()
	polymorphic.RegisterType[Acquire]()
	polymorphic.RegisterType[Renew]()
	polymorphic.RegisterType[Release]()
}

// Lease represents a distributed lease.
//
// A Lease contains identifying information and an expiry timestamp. The
// TTL field indicates the lease duration used when the lease was issued.
type Lease struct {
	ID       uuid.UUID     `json:"id"`
	TenantID uuid.UUID     `json:"tenant_id"`
	Key      string        `json:"key"`
	TTL      time.Duration `json:"ttl"`       // duration for which lease is valid
	ExpireAt time.Time     `json:"expire_at"` // when the lease will expire
}

func (*Lease) GetDiscriminator() string {
	return "mesh://leasekit/lease"
}

// Acquire represents a request to acquire a Lease.
//
// MaxAttempts allows callers to indicate how many attempts the manager
// should try before failing the request.
type Acquire struct {
	ID          uuid.UUID     `json:"id"`
	TenantID    uuid.UUID     `json:"tenant_id"`
	Key         string        `json:"key"`
	TTL         time.Duration `json:"ttl"`
	MaxAttempts int           `json:"max_attempts"`
}

func (*Acquire) GetDiscriminator() string {
	return "mesh://leasekit/acquire"
}

// GetRoute returns the messaging route for an Acquire request.
func (r *Acquire) GetRoute() messaging.Route {
	return GetAcquireRoute(&r.TenantID)
}

// Renew represents a request to renew an existing Lease.
type Renew struct {
	ID       uuid.UUID     `json:"id"`
	TenantID uuid.UUID     `json:"tenant_id"`
	Key      string        `json:"key"`
	TTL      time.Duration `json:"ttl"`
}

func (*Renew) GetDiscriminator() string {
	return "mesh://leasekit/renew"
}

// GetRoute returns the messaging route for a Renew request.
func (r *Renew) GetRoute() messaging.Route {
	return GetRenewRoute(&r.TenantID)
}

// Release represents a request to release a Lease.
type Release struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Key      string    `json:"key"`
}

func (*Release) GetDiscriminator() string {
	return "mesh://leasekit/release"
}

// GetRoute returns the messaging route for a Release request.
func (r *Release) GetRoute() messaging.Route {
	return GetReleaseRoute(&r.TenantID)
}

// Route helper functions
// GetAcquireRoute returns the tenant-scoped messaging route for acquire requests.
func GetAcquireRoute(tenantID *uuid.UUID) messaging.Route {
	return messaging.NewTenantRoute("leases", "acquire", tenantID)
}

// GetRenewRoute returns the tenant-scoped messaging route for renew requests.
func GetRenewRoute(tenantID *uuid.UUID) messaging.Route {
	return messaging.NewTenantRoute("leases", "renew", tenantID)
}

// GetReleaseRoute returns the tenant-scoped messaging route for release requests.
func GetReleaseRoute(tenantID *uuid.UUID) messaging.Route {
	return messaging.NewTenantRoute("leases", "release", tenantID)
}
