package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgon2HashAndVerify(t *testing.T) {
	password := "correcthorsebatterystaple"
	enc, err := hashPasswordArgon2(password)
	require.NoError(t, err)
	assert.True(t, len(enc) > 0)
	ok, err := verifyPasswordArgon2(enc, password)
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestCheckPasswordHashArgon2(t *testing.T) {
	password := "hunter2"
	enc, err := HashPassword(password)
	require.NoError(t, err)
	assert.True(t, CheckPasswordHash(password, enc))
	assert.False(t, CheckPasswordHash("wrong", enc))
}
