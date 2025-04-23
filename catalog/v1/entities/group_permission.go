package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *GroupPermission { return &GroupPermission{} })
}

type GroupPermission struct {
	ConnectorID         uuid.UUID `json:"connector_id"`
	GroupReference      string    `json:"group_reference"`
	PermissionReference string    `json:"permission_reference"`
}

func (e *GroupPermission) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *GroupPermission) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *GroupPermission) GetReference() string {
	return CompoundRef(e.GroupReference, e.PermissionReference)
}

func (e *GroupPermission) GetSpace() string {
	return "group-permissions"
}

func (e *GroupPermission) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/group-permission"
}
