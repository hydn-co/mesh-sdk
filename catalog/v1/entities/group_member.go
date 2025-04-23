package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *GroupMember { return &GroupMember{} })
}

type GroupMember struct {
	ConnectorID      uuid.UUID `json:"connector_id"`
	GroupReference   string    `json:"group_reference"`
	AccountReference string    `json:"account_reference"`
	RoleReference    string    `json:"role_reference,omitempty"`
}

func (e *GroupMember) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *GroupMember) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *GroupMember) GetReference() string {
	return CompoundRef(e.GroupReference, e.AccountReference)
}

func (e *GroupMember) GetSpace() string {
	return "group-members"
}

func (e *GroupMember) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/group-member"
}
