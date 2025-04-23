package entities

import (
	"strings"

	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

type Entity interface {
	polymorphic.Polymorphic
	GetSpace() string
	GetReference() string
	GetConnectorID() uuid.UUID
	GetDistinctID() uuid.UUID
}

func CompoundRef(parts ...string) string {
	return strings.Join(parts, ":")
}
