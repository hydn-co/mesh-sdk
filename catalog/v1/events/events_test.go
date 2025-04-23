package events

import (
	"testing"

	"github.com/fgrzl/json/polymorphic"
	"github.com/stretchr/testify/require"
)

func TestShouldRegisterDiscriminator(t *testing.T) {

	// Arrange
	var expectedEntities = map[string]EntityEvent{
		"hydn://catalog/v1/events/entity-patched": &EntityPatched{},
		"hydn://catalog/v1/events/entity-removed": &EntityRemoved{},
	}

	for discriminator, expectedType := range expectedEntities {
		t.Run(discriminator,
			func(t *testing.T) {
				// Act
				instance, err := polymorphic.CreateInstance(discriminator)

				// Assert
				require.NoError(t, err, "Failed to create instance for %s", discriminator)
				require.IsType(t, expectedType, instance, "Unexpected type for %s", discriminator)
			})
	}
}
