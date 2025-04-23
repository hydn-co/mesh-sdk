package events

import (
	"context"

	"github.com/fgrzl/json/jsonpatch"
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *EntityPatched { return &EntityPatched{} })
	polymorphic.Register(func() *EntityRemoved { return &EntityRemoved{} })
}

type EntityEvent interface {
	GetCausationID() uuid.UUID
	GetConnectorID() uuid.UUID
	GetCorrelationID() uuid.UUID
	GetEventID() uuid.UUID
	GetID() uuid.UUID
	GetReference() string
	GetSpace() string
	GetTimestamp() int64
}

// Base struct for all Mesh-derived events
type EntityEventBase struct {
	ID            uuid.UUID `json:"id"`
	ConnectorID   uuid.UUID `json:"connector_id"`
	Reference     string    `json:"reference"`
	Space         string    `json:"space"`
	EventID       uuid.UUID `json:"event_id"`
	CorrelationID uuid.UUID `json:"correlation_id,omitempty"`
	CausationID   uuid.UUID `json:"causation_id,omitempty"`
	Timestamp     int64     `json:"timestamp"`
	Sequence      uint64    `json:"sequence,omitempty"`
}

func (e *EntityEventBase) GetCausationID() uuid.UUID   { return e.CausationID }
func (e *EntityEventBase) GetConnectorID() uuid.UUID   { return e.ConnectorID }
func (e *EntityEventBase) GetCorrelationID() uuid.UUID { return e.CorrelationID }
func (e *EntityEventBase) GetEventID() uuid.UUID       { return e.EventID }
func (e *EntityEventBase) GetID() uuid.UUID            { return e.ID }
func (e *EntityEventBase) GetReference() string        { return e.Reference }
func (e *EntityEventBase) GetSpace() string            { return e.Space }
func (e *EntityEventBase) GetTimestamp() int64         { return e.Timestamp }

func NewEntityPatched(ctx context.Context, connectorId uuid.UUID, space, ref string, patch []jsonpatch.Patch) *EntityPatched {

	return &EntityPatched{
		EntityEventBase: EntityEventBase{
			EventID: uuid.New(),
		},
		Patch: patch,
	}
}

type EntityPatched struct {
	EntityEventBase
	Patch []jsonpatch.Patch `json:"patch"`
}

func (e *EntityPatched) GetDiscriminator() string {
	return "hydn://catalog/v1/events/entity-patched"
}

func NewEntityRemoved(ctx context.Context, connectorId uuid.UUID, space, ref string) *EntityRemoved {
	return &EntityRemoved{
		EntityEventBase: EntityEventBase{},
	}
}

type EntityRemoved struct {
	EntityEventBase
}

func (e *EntityRemoved) GetDiscriminator() string {
	return "hydn://catalog/v1/events/entity-removed"
}
