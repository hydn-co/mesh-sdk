package requests

import (
	"testing"

	"github.com/fgrzl/json/polymorphic"
	"github.com/stretchr/testify/require"
)

// make sure we register all domain events and use consistent discriminators
// keep this list up to date with the events in the domain package
// and lets keep it sorted by the discriminator
var domainRequests = map[string]any{
	"hydn://domain/requests/clients/disabled": &DisableClient{},
	"hydn://domain/requests/clients/enable":   &EnableClient{},

	// Add more here
}

func TestShouldRegisterPolymorphicDomainRequests(t *testing.T) {

	for discriminator, instance := range domainRequests {
		t.Run(discriminator, func(t *testing.T) {
			// Polymorphic check
			created, err := polymorphic.CreateInstance(discriminator)
			require.IsType(t, instance, created, "instance type mismatch for %s", discriminator)
			require.NoError(t, err, "could not create instance for %s", discriminator)
			require.NotNil(t, created, "nil returned for %s", discriminator)
		})
	}
}
