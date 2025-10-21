package localstore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

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

func TestGetBaseAndConfigPaths(t *testing.T) {
	withTempConfigDir(t, func(t *testing.T, tempDir string) {
		base, err := GetBasePath()
		require.NoError(t, err)
		// base should be under tempDir/hyddenlabs/mesh
		require.Contains(t, base, filepath.Join(tempDir, "hyddenlabs", "mesh"))

		conf, err := GetConfigPath()
		require.NoError(t, err)
		require.Contains(t, conf, filepath.Join(tempDir, "hyddenlabs", "mesh", ".config"))
	})
}

func TestTenantPathsAndCreds(t *testing.T) {
	withTempConfigDir(t, func(t *testing.T, tempDir string) {
		tenant := uuid.New()
		path, err := GetTenantPath(tenant)
		require.NoError(t, err)
		require.Contains(t, path, tenant.String())

		dataPath, err := GetDataPath(tenant)
		require.NoError(t, err)
		require.Contains(t, dataPath, filepath.Join(tenant.String(), "data"))

		credsPath, err := GetCredsPath(tenant)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(path, ".creds"), credsPath)
	})
}

func TestProviderPaths(t *testing.T) {
	withTempConfigDir(t, func(t *testing.T, tempDir string) {
		provider := uuid.New()
		base, err := GetProvidersBasePath()
		require.NoError(t, err)
		require.Contains(t, base, filepath.Join(tempDir, "hyddenlabs", "mesh", "providers"))

		pp, err := GetProviderBasePath(provider)
		require.NoError(t, err)
		require.Contains(t, pp, provider.String())
	})
}
