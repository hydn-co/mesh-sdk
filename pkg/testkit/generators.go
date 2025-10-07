package testkit

import (
	"github.com/google/uuid"
)

func NewConnectorID(name string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("connector"+name)).String()
}

func NewSecretID(name string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("secret"+name)).String()
}

// GetConnectorArgs converts a collector name to capability and connector ID.
// E.g. "IamUserCollector" => "iam-user-collector", "iam-user-collector-123"
func GetConnectorArgs(name string) (capability, connectorID string) {
	if name == "" {
		panic("GetConnectorArgs: name must not be empty")
	}

	capability = name
	connectorID = NewConnectorID(capability)
	return capability, connectorID
}
