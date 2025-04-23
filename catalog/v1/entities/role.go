package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *Role { return &Role{} })
}

type Role struct {
	ConnectorID   uuid.UUID `json:"connector_id"`
	RoleReference string    `json:"role_reference"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
}

func (e *Role) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *Role) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *Role) GetReference() string {
	return e.RoleReference
}

func (e *Role) GetSpace() string {
	return "roles"
}

func (e *Role) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/role"
}
