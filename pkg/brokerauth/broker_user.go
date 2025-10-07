package auth

import "github.com/google/uuid"

type AuthBrokerUser struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	ClientID     uuid.UUID `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
}

type BrokerUser struct {
	Token []byte `json:"token"`
	Seed  []byte `json:"seed"`
}
