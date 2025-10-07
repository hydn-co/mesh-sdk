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

// Lease represents a distributed lease
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

// Acquire represents a request to acquire a lease
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

func (r *Acquire) GetRoute() messaging.Route {
	return GetAcquireRoute(&r.TenantID)
}

// Renew represents a request to renew an existing lease
type Renew struct {
	ID       uuid.UUID     `json:"id"`
	TenantID uuid.UUID     `json:"tenant_id"`
	Key      string        `json:"key"`
	TTL      time.Duration `json:"ttl"`
}

func (*Renew) GetDiscriminator() string {
	return "mesh://leasekit/renew"
}

func (r *Renew) GetRoute() messaging.Route {
	return GetRenewRoute(&r.TenantID)
}

// Release represents a request to release a lease
type Release struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Key      string    `json:"key"`
}

func (*Release) GetDiscriminator() string {
	return "mesh://leasekit/release"
}

func (r *Release) GetRoute() messaging.Route {
	return GetReleaseRoute(&r.TenantID)
}

// Route helper functions
func GetAcquireRoute(tenantID *uuid.UUID) messaging.Route {
	return messaging.NewTenantRoute("leases", "acquire", tenantID)
}

func GetRenewRoute(tenantID *uuid.UUID) messaging.Route {
	return messaging.NewTenantRoute("leases", "renew", tenantID)
}

func GetReleaseRoute(tenantID *uuid.UUID) messaging.Route {
	return messaging.NewTenantRoute("leases", "release", tenantID)
}
