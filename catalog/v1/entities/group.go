package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *Group { return &Group{} })
}

type Group struct {
	ConnectorID    uuid.UUID `json:"connector_id"`
	GroupReference string    `json:"group_reference"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
}

func (e *Group) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *Group) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *Group) GetReference() string {
	return e.GroupReference
}

func (e *Group) GetSpace() string {
	return "groups"
}

func (e *Group) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/group"
}
