package requests

import (
	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
)

type GetTenant struct {
	TenantID uuid.UUID `json:"tenant_id"`
}

func (g *GetTenant) GetDiscriminator() string {
	return "query://tenants/get"
}

func (g *GetTenant) GetRoute() messaging.Route {
	return messaging.NewInternalRoute("slug", "reserve")
}

type RegisterTenant struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	RequestedBy uuid.UUID `json:"requested_by"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
}

func (c *RegisterTenant) GetDiscriminator() string {
	return "command://tenants/register"
}

func (c *RegisterTenant) GetRoute() messaging.Route {
	var tenantID *uuid.UUID
	if c != nil {
		tenantID = &c.TenantID
	}
	return messaging.NewTenantRoute("org", "get_org", tenantID)
}

type RenameTenant struct {
	TenantID uuid.UUID `json:"tenant_id"`
	Name     string    `json:"slug"`
}

func (c *RenameTenant) GetDiscriminator() string {
	return "command://tenants/rename"
}

func (c *RenameTenant) GetRoute() messaging.Route {
	var orgID *uuid.UUID
	if c != nil {
		orgID = &c.TenantID
	}
	return messaging.NewTenantRoute("org", "rename", orgID)
}

type ChangeTenantSlug struct {
	TenantID uuid.UUID `json:"tenant_id"`
	Slug     string    `json:"slug"`
}

func (c *ChangeTenantSlug) GetDiscriminator() string {
	return "command://tenants/change_slug"
}

func (c *ChangeTenantSlug) GetRoute() messaging.Route {
	return messaging.NewInternalRoute("slug", "reserve")
}

type GetTenants struct{}

func (c *GetTenants) GetDiscriminator() string {
	return "hydn://tentants/get"
}

func (c *GetTenants) GetRoute() messaging.Route {
	return messaging.NewInternalRoute("slug", "reserve")
}
