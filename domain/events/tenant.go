package events

import (
	"github.com/fgrzl/es"
	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
)

var (
	_ es.DomainEvent = &TenantRequested{}
)

func init() {
	es.Register(func() *TenantRequested { return &TenantRequested{} })
	es.Register(func() *TenantRenamed { return &TenantRenamed{} })
	es.Register(func() *TenantSlugChanged { return &TenantSlugChanged{} })
}

type TenantRequested struct {
	es.DomainEventBase
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (e *TenantRequested) GetDiscriminator() string { return "event://tenants/requested" }

func (e *TenantRequested) GetRoute() messaging.Route {
	var orgID *uuid.UUID
	if e != nil {
		x := e.GetAggregateID()
		orgID = &x
	}
	return messaging.NewTenantRoute("tenants", "tenant_requested", orgID)
}

type TenantRenamed struct {
	es.DomainEventBase
	Name string `json:"name"`
}

func (e *TenantRenamed) GetDiscriminator() string { return "event://tenants/renamed" }

type TenantSlugChanged struct {
	es.DomainEventBase
	Slug    string `json:"slug"`
	OldSlug string `json:"old_slug"`
}

func (e *TenantSlugChanged) GetDiscriminator() string {
	return "event://tenants/slug_changed"
}

func (e *TenantSlugChanged) GetRoute() messaging.Route {
	var orgID *uuid.UUID
	if e != nil {
		x := e.GetAggregateID()
		orgID = &x
	}
	return messaging.NewTenantRoute("tenants", "tenant_slug_changed", orgID)
}
