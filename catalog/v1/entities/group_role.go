package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

var _ Entity = (*GroupRole)(nil)

func init() {
	polymorphic.Register(func() *GroupRole { return &GroupRole{} })
}

type GroupRole struct {
	ConnectorID    uuid.UUID `json:"connector_id"`
	GroupReference string    `json:"group_reference"`
	RoleReference  string    `json:"role_reference"`
}

func (e *GroupRole) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *GroupRole) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *GroupRole) GetReference() string {
	return CompoundRef(e.GroupReference, e.RoleReference)
}

func (e *GroupRole) GetSpace() string {
	return "group-roles"
}

func (e *GroupRole) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/group-role"
}
