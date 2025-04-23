package entities

import (
	"testing"

	"github.com/fgrzl/json/polymorphic"
	"github.com/stretchr/testify/require"
)

func TestShouldRegisterDiscriminator(t *testing.T) {

	// Arrange
	var expectedEntities = map[string]Entity{
		"hydn://catalog/v1/entities/account-permission": &AccountPermission{},
		"hydn://catalog/v1/entities/account-role":       &AccountRole{},
		"hydn://catalog/v1/entities/account":            &Account{},
		"hydn://catalog/v1/entities/application":        &Application{},
		"hydn://catalog/v1/entities/group-member":       &GroupMember{},
		"hydn://catalog/v1/entities/group-permission":   &GroupPermission{},
		"hydn://catalog/v1/entities/group-role":         &GroupRole{},
		"hydn://catalog/v1/entities/group":              &Group{},
		"hydn://catalog/v1/entities/multi-factor":       &MultiFactor{},
		"hydn://catalog/v1/entities/network-object":     &NetworkObject{},
		"hydn://catalog/v1/entities/permission":         &Permission{},
		"hydn://catalog/v1/entities/role":               &Role{},
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

func TestShouldHaveSpace(t *testing.T) {

	// Arrange
	var expectedEntities = map[string]Entity{
		"account-permissions": &AccountPermission{},
		"account-roles":       &AccountRole{},
		"accounts":            &Account{},
		"applications":        &Application{},
		"group-members":       &GroupMember{},
		"group-permissions":   &GroupPermission{},
		"group-roles":         &GroupRole{},
		"groups":              &Group{},
		"multi-factors":       &MultiFactor{},
		"network-objects":     &NetworkObject{},
		"permissions":         &Permission{},
		"roles":               &Role{},
	}

	for expected, instance := range expectedEntities {
		t.Run(expected,
			func(t *testing.T) {
				// Act
				actual := instance.GetSpace()

				// Assert
				require.Equal(t, expected, actual)
			})
	}
}
