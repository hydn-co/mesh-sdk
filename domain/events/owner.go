package events

import (
	"github.com/fgrzl/es"
	"github.com/google/uuid"
)

var (
	_ es.DomainEvent = &OwnerAssigned{}
	_ es.DomainEvent = &OwnerUnassigned{}
)

func init() {
	es.Register(func() *OwnerAssigned { return &OwnerAssigned{} })
	es.Register(func() *OwnerUnassigned { return &OwnerUnassigned{} })
}

// OwnerAssigned event
type OwnerAssigned struct {
	es.DomainEventBase
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"user_id"`
}

func (e *OwnerAssigned) GetDiscriminator() string {
	return "hydn://domain/events/owners/assigned"
}

// OwnerUnassigned event
type OwnerUnassigned struct {
	es.DomainEventBase
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"user_id"`
}

func (e *OwnerUnassigned) GetDiscriminator() string {
	return "hydn://domain/events/owners/unassigned"
}
