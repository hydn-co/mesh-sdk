package events

import (
	"testing"

	"github.com/fgrzl/es"
	"github.com/fgrzl/json/polymorphic"
	"github.com/stretchr/testify/require"
)

// make sure we register all domain events and use consistent discriminators
// keep this list up to date with the events in the domain package
// and lets keep it sorted by the discriminator
var domainEvents = map[string]any{
	"hydn://domain/events/clients/disabled":   &ClientDisabled{},
	"hydn://domain/events/clients/enabled":    &ClientEnabled{},
	"hydn://domain/events/clients/registered": &ClientRegistered{},
	"hydn://domain/events/clients/removed":    &ClientRemoved{},
	"hydn://domain/events/owners/assigned":    &OwnerAssigned{},
	"hydn://domain/events/owners/unassigned":  &OwnerUnassigned{},
	"hydn://domain/events/slugs/reserved":     &SlugReserved{},
	"hydn://domain/events/slugs/surrendered":  &SlugSurrendered{},

	// Add more here
}

func TestShouldRegisterPolymorphicDomainEvents(t *testing.T) {

	for discriminator, instance := range domainEvents {
		t.Run(discriminator, func(t *testing.T) {
			// Interface check
			_, ok := instance.(es.DomainEvent)
			require.True(t, ok, "%T does not implement DomainEvent", instance)

			// Polymorphic check
			created, err := polymorphic.CreateInstance(discriminator)
			require.NoError(t, err, "could not create instance for %s", discriminator)
			require.NotNil(t, created, "nil returned for %s", discriminator)
		})
	}
}
