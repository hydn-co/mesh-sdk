package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *MultiFactor { return &MultiFactor{} })
}

func NewMultiFactor() *MultiFactor { return &MultiFactor{} }

// MultiFactorStatus represents the status of multi-factor authentication.
type MultiFactorStatus string

const (
	MFAStatusNotApplicable MultiFactorStatus = "NOT_APPLICABLE"
	MFAStatusPending       MultiFactorStatus = "PENDING"
	MFAStatusEnabled       MultiFactorStatus = "ENABLED"
)

// MultiFactorKind represents the kind of multi-factor authentication.
type MultiFactorKind string

const (
	MFAKindUnknown       MultiFactorKind = "UNKNOWN"
	MFAKindSMS           MultiFactorKind = "SMS"
	MFAKindEmail         MultiFactorKind = "EMAIL"
	MFAKindAuthenticator MultiFactorKind = "AUTHENTICATOR"
	MFAKindOther         MultiFactorKind = "OTHER"
)

type MultiFactor struct {
	ConnectorID      uuid.UUID         `json:"connector_id"`
	AccountReference string            `json:"account_reference"`
	Status           MultiFactorStatus `json:"status" enum:"NOT_APPLICABLE,PENDING,ENABLED"`
	Kind             MultiFactorKind   `json:"kind" enum:"UNKNOWN,SMS,EMAIL,AUTHENTICATOR,OTHER"`
}

func (e *MultiFactor) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *MultiFactor) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *MultiFactor) GetReference() string {
	return e.AccountReference
}

func (e *MultiFactor) GetSpace() string {
	return "multi-factors"
}

func (e *MultiFactor) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/multi-factor"
}
