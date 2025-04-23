package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *Application { return &Application{} })
}

type Application struct {
	ConnectorID          uuid.UUID `json:"connector_id"`
	ApplicationReference string    `json:"application_reference"`
	Name                 string    `json:"name"`
	Description          string    `json:"description,omitempty"`
}

func (e *Application) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *Application) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *Application) GetReference() string {
	return e.ApplicationReference
}

func (e *Application) GetSpace() string {
	return "applications"
}

func (e *Application) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/application"
}
