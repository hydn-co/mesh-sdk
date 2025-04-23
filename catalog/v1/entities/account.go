package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *Account { return &Account{} })
}

type AccountType string

type Account struct {
	ConnectorID      uuid.UUID   `json:"connector_id"`
	AccountReference string      `json:"account_reference"`
	AccountType      AccountType `json:"account_type"`
	Name             string      `json:"name,omitempty"`
	Description      string      `json:"description,omitempty"`
	DisplayName      string      `json:"display_name,omitempty"`
	FirstName        string      `json:"first_name,omitempty"`
	MiddleName       string      `json:"middle_name,omitempty"`
	LastName         string      `json:"last_name,omitempty"`
	PrimaryEmail     *Email      `json:"primary_email"`
	AlternateEmails  []*Email    `json:"alternate_emails,omitempty"`
	PrimaryPhone     *Phone      `json:"primary_phone"`
	AlternatePhones  []*Phone    `json:"alternate_phones,omitempty"`
}

func (e *Account) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *Account) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *Account) GetReference() string {
	return e.AccountReference
}

func (e *Account) GetSpace() string {
	return "accounts"
}

func (e *Account) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/account"
}

type Email struct {
	Address string `json:"address"`
}

type Phone struct {
	Number string `json:"number"`
}
