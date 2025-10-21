package broker

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateAuthBrokerUserWithValidFields(t *testing.T) {
	// Arrange
	tenantID := uuid.New()
	clientID := uuid.New()
	secret := "test-secret"

	// Act
	user := AuthBrokerUser{
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: secret,
	}

	// Assert
	assert.Equal(t, tenantID, user.TenantID)
	assert.Equal(t, clientID, user.ClientID)
	assert.Equal(t, secret, user.ClientSecret)
}

func TestShouldCreateBrokerUserWithTokenAndSeed(t *testing.T) {
	// Arrange
	token := []byte("test-token")
	seed := []byte("test-seed")

	// Act
	user := BrokerUser{
		Token: token,
		Seed:  seed,
	}

	// Assert
	assert.Equal(t, token, user.Token)
	assert.Equal(t, seed, user.Seed)
}

func TestShouldAllowEmptyBrokerUserFields(t *testing.T) {
	// Arrange & Act
	user := BrokerUser{}

	// Assert
	assert.Nil(t, user.Token)
	assert.Nil(t, user.Seed)
}
