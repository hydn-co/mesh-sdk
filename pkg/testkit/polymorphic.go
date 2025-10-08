package testkit

import (
	"testing"

	"github.com/fgrzl/es"
	"github.com/fgrzl/json/polymorphic"
	"github.com/stretchr/testify/require"
)

// TestPolymorphicRegistrations tests that all provided types are properly registered
// with the polymorphic system and can be created from their discriminators.
func TestPolymorphicRegistrations(t *testing.T, types map[string]any) {
	t.Helper()

	for discriminator, instance := range types {
		t.Run(discriminator, func(t *testing.T) {
			// Polymorphic check
			created, err := polymorphic.CreateInstance(discriminator)
			require.NoError(t, err, "could not create instance for %s", discriminator)
			require.NotNil(t, created, "nil returned for %s", discriminator)
			require.IsType(t, instance, created, "instance type mismatch for %s", discriminator)
		})
	}
}

// TestDomainEventRegistrations tests that all provided domain events are properly registered
// with the polymorphic system, implement the DomainEvent interface, and can be created from their discriminators.
func TestDomainEventRegistrations(t *testing.T, events map[string]any) {
	t.Helper()

	for discriminator, instance := range events {
		t.Run(discriminator, func(t *testing.T) {
			// Interface check
			_, ok := instance.(es.DomainEvent)
			require.True(t, ok, "%T does not implement DomainEvent", instance)

			// Polymorphic check
			created, err := polymorphic.CreateInstance(discriminator)
			require.NoError(t, err, "could not create instance for %s", discriminator)
			require.NotNil(t, created, "nil returned for %s", discriminator)
			require.IsType(t, instance, created, "instance type mismatch for %s", discriminator)
		})
	}
}
