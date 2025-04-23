package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *AccountPermission { return &AccountPermission{} })
}

type AccountPermission struct {
	ConnectorID         uuid.UUID `json:"connector_id"`
	AccountReference    string    `json:"account_reference"`
	PermissionReference string    `json:"permission_reference"`
}

func (e *AccountPermission) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *AccountPermission) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *AccountPermission) GetReference() string {
	return CompoundRef(e.AccountReference, e.PermissionReference)
}

func (e *AccountPermission) GetSpace() string {
	return "account-permissions"
}

func (e *AccountPermission) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/account-permission"
}
