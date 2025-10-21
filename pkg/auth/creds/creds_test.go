package creds

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/localstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to set temp config dir for tests
func withTempConfigDir(t *testing.T, fn func(t *testing.T, tempDir string)) {
	t.Helper()
	tempDir := t.TempDir()
	// Set both XDG_CONFIG_HOME and APPDATA for cross-platform tests
	origXDG := os.Getenv("XDG_CONFIG_HOME")
	origAPP := os.Getenv("APPDATA")
	defer func() {
		os.Setenv("XDG_CONFIG_HOME", origXDG)
		os.Setenv("APPDATA", origAPP)
	}()
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	os.Setenv("APPDATA", tempDir)

	fn(t, tempDir)
}

func TestShouldUseClientIDAndSecretFromEnv(t *testing.T) {
	withTempConfigDir(t, func(t *testing.T, _ string) {
		// Arrange
		tenant := uuid.New()
		clientID := uuid.New()
		clientSecret := "env-secret"
		seed := "SUAxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

		require.NoError(t, os.Setenv("MESH_CLIENT_ID", clientID.String()))
		require.NoError(t, os.Setenv("MESH_CLIENT_SECRET", clientSecret))
		require.NoError(t, os.Setenv("MESH_CLIENT_SEED", seed))
		defer func() {
			os.Unsetenv("MESH_CLIENT_ID")
			os.Unsetenv("MESH_CLIENT_SECRET")
			os.Unsetenv("MESH_CLIENT_SEED")
		}()

		// Act
		creds, err := LoadOrCreateCreds(tenant)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, clientID, creds.ClientID)
		assert.Equal(t, clientSecret, creds.ClientSecret)
		pub, err := creds.User.PublicKey()
		require.NoError(t, err)
		assert.NotEmpty(t, pub)
	})
}

func TestShouldCreateNewNkeyIfSeedNotAvailable(t *testing.T) {
	withTempConfigDir(t, func(t *testing.T, _ string) {
		// Arrange
		tenant := uuid.New()
		clientID := uuid.New()
		clientSecret := "env-secret-2"

		require.NoError(t, os.Setenv("MESH_CLIENT_ID", clientID.String()))
		require.NoError(t, os.Setenv("MESH_CLIENT_SECRET", clientSecret))
		defer func() {
			os.Unsetenv("MESH_CLIENT_ID")
			os.Unsetenv("MESH_CLIENT_SECRET")
		}()

		// Act
		creds, err := LoadOrCreateCreds(tenant)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, clientID, creds.ClientID)
		assert.Equal(t, clientSecret, creds.ClientSecret)
		pub, err := creds.User.PublicKey()
		require.NoError(t, err)
		assert.NotEmpty(t, pub)
	})
}

func TestShouldPersistFileIfNotExists(t *testing.T) {
	withTempConfigDir(t, func(t *testing.T, tempDir string) {
		// Arrange
		tenant := uuid.New()
		// Ensure no env vars
		os.Unsetenv("MESH_CLIENT_ID")
		os.Unsetenv("MESH_CLIENT_SECRET")
		os.Unsetenv("MESH_CLIENT_SEED")

		// Act
		creds, err := LoadOrCreateCreds(tenant)

		// Assert
		require.NoError(t, err)
		// Check file exists
		path, err := localstore.GetCredsPath(tenant)
		require.NoError(t, err)
		_, err = os.Stat(path)
		require.NoError(t, err)

		// Validate file content parses
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		var cf credsFile
		require.NoError(t, json.Unmarshal(data, &cf))
		assert.Equal(t, cf.ClientID, creds.ClientID)
		assert.Equal(t, cf.ClientSecret, creds.ClientSecret)
	})
}

func TestEnvTakesPrecedenceOverFile(t *testing.T) {
	withTempConfigDir(t, func(t *testing.T, tempDir string) {
		// Arrange - create a creds file with one set of values
		tenant := uuid.New()
		fileClientID := uuid.New()
		fileSecret := "file-secret"
		cf := credsFile{ClientID: fileClientID, ClientSecret: fileSecret, ClientSeed: []byte("seed-bytes")}
		path, err := localstore.GetCredsPath(tenant)
		require.NoError(t, err)
		// Ensure parent dirs exist
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0700))
		data, err := json.MarshalIndent(cf, "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(path, data, 0600))

		// Now set env to override
		envClientID := uuid.New()
		envClientSecret := "env-override-secret"
		require.NoError(t, os.Setenv("MESH_CLIENT_ID", envClientID.String()))
		require.NoError(t, os.Setenv("MESH_CLIENT_SECRET", envClientSecret))
		defer func() {
			os.Unsetenv("MESH_CLIENT_ID")
			timeout := os.Unsetenv("MESH_CLIENT_SECRET")
			_ = timeout
		}()

		// Act
		creds, err := LoadOrCreateCreds(tenant)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, envClientID, creds.ClientID)
		assert.Equal(t, envClientSecret, creds.ClientSecret)
	})
}
