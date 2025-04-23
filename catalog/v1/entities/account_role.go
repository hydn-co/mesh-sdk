package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *AccountRole { return &AccountRole{} })
}

type AccountRole struct {
	ConnectorID      uuid.UUID `json:"connector_id"`
	AccountReference string    `json:"account_reference"`
	RoleReference    string    `json:"role_reference"`
}

func (e *AccountRole) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *AccountRole) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *AccountRole) GetReference() string {
	return CompoundRef(e.AccountReference, e.RoleReference)
}

func (e *AccountRole) GetSpace() string {
	return "account-roles"
}

func (e *AccountRole) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/account-role"
}
