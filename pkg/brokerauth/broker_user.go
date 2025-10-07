package auth

import "github.com/google/uuid"

// AuthBrokerUser is the payload sent to the broker auth endpoint when
// requesting a token for a tenant-scoped client.
type AuthBrokerUser struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	ClientID     uuid.UUID `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
}

// BrokerUser represents the broker-side response containing a token and the
// nkeys seed for the created user. Both fields are binary blobs and may be
// encoded (for example base64) by the transport layer.
type BrokerUser struct {
	Token []byte `json:"token"`
	Seed  []byte `json:"seed"`
}
