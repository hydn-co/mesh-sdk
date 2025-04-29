package requests

import (
	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
)

type RegisterClient struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	ClientID     uuid.UUID `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
}

func (c *RegisterClient) GetDiscriminator() string {
	return "client_registered"
}

func (e *RegisterClient) GetRoute() messaging.Route {
	return messaging.NewTenantRoute("slugs", "reserved", &e.TenantID)
}

type EnableClient struct {
	TenantID uuid.UUID `json:"tenant_id"`
	ClientID uuid.UUID `json:"client_id"`
}

type DisableClient struct {
	TenantID uuid.UUID `json:"tenant_id"`
	ClientID uuid.UUID `json:"client_id"`
}

type RemoveClient struct {
	TenantID uuid.UUID `json:"tenant_id"`
	ClientID uuid.UUID `json:"client_id"`
}
