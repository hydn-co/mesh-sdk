package testkit

import (
	"github.com/google/uuid"
)

// NewConnectorID returns a deterministic connector UUID for a given name.
// It is used in tests to produce stable IDs across runs.
func NewConnectorID(name string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("connector"+name)).String()
}

// NewSecretID returns a deterministic secret UUID for a given name.
// Useful for test fixtures where reproducible identifiers are desired.
func NewSecretID(name string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("secret"+name)).String()
}

// GetConnectorArgs converts a capability name into a capability identifier and
// a deterministic connector ID. The function panics if name is empty and is
// intended for use in test code where a missing name indicates a programmer
// error.
func GetConnectorArgs(name string) (capability, connectorID string) {
	if name == "" {
		panic("GetConnectorArgs: name must not be empty")
	}

	capability = name
	connectorID = NewConnectorID(capability)
	return capability, connectorID
}
