package events

import (
	"github.com/fgrzl/es"
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *ClientRegistered { return &ClientRegistered{} })
	polymorphic.Register(func() *ClientEnabled { return &ClientEnabled{} })
	polymorphic.Register(func() *ClientDisabled { return &ClientDisabled{} })
	polymorphic.Register(func() *ClientRemoved { return &ClientRemoved{} })
}

type ClientRegistered struct {
	es.DomainEventBase
	TenantID     uuid.UUID `json:"tenant_id"`
	ClientID     uuid.UUID `json:"client_id"`
	ClientSecret []byte    `json:"client_secret"`
}

func (e *ClientRegistered) GetDiscriminator() string {
	return "hydn://domain/events/clients/registered"
}

type ClientEnabled struct {
	es.DomainEventBase
	TenantID uuid.UUID `json:"tenant_id"`
	ClientID uuid.UUID `json:"client_id"`
}

func (e *ClientEnabled) GetDiscriminator() string {
	return "hydn://domain/events/clients/enabled"
}

type ClientDisabled struct {
	es.DomainEventBase
	TenantID uuid.UUID `json:"tenant_id"`
	ClientID uuid.UUID `json:"client_id"`
}

func (e *ClientDisabled) GetDiscriminator() string {
	return "hydn://domain/events/clients/disabled"
}

type ClientRemoved struct {
	es.DomainEventBase
	TenantID uuid.UUID `json:"tenant_id"`
	ClientID uuid.UUID `json:"client_id"`
}

func (e *ClientRemoved) GetDiscriminator() string {
	return "hydn://domain/events/clients/removed"
}
