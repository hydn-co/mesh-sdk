package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *Permission { return &Permission{} })
}

type Permission struct {
	ConnectorID         uuid.UUID `json:"connector_id"`
	PermissionReference string    `json:"permission_reference"`
	Name                string    `json:"name"`
	Description         string    `json:"description,omitempty"`
}

func (e *Permission) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *Permission) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *Permission) GetReference() string {
	return e.PermissionReference
}

func (e *Permission) GetSpace() string {
	return "permissions"
}

func (e *Permission) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/permission"
}
