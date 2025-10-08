package localstore

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	lsDirOrg       = "hyddenlabs"
	lsDirApp       = "mesh"
	lsDirConfig    = ".config"
	lsDirTenants   = "tenants"
	lsDirData      = "data"
	lsDirProviders = "providers"
	lsCredsFile    = ".creds"
)

// GetBasePath returns the base config directory: <user config dir>/hyddenlabs/mesh
// - macOS:   $HOME/Library/Application Support/hyddenlabs/mesh
// - Linux:   $XDG_CONFIG_HOME/hyddenlabs/mesh or $HOME/.config/hyddenlabs/mesh
// - Windows: %AppData%\hyddenlabs\mesh
func GetBasePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return getOrCreatePath(filepath.Join(configDir, lsDirOrg, lsDirApp))
}

// GetConfigPath returns: <base config dir>/.config
// - macOS:   $HOME/Library/Application Support/hyddenlabs/mesh/.config
// - Linux:   $XDG_CONFIG_HOME/hyddenlabs/mesh/.config or $HOME/.config/hyddenlabs/mesh/.config
// - Windows: %AppData%\hyddenlabs\mesh\.config
func GetConfigPath() (string, error) {
	p, err := joinBasePath(lsDirConfig)
	if err != nil {
		return "", err
	}
	return getOrCreatePath(p)
}

// GetTenantRootPath returns: <base config dir>/tenants
// - macOS:   $HOME/Library/Application Support/hyddenlabs/mesh/tenants
// - Linux:   $XDG_CONFIG_HOME/hyddenlabs/mesh/tenants or $HOME/.config/hyddenlabs/mesh/tenants
// - Windows: %AppData%\hyddenlabs\mesh\tenants
func GetTenantRootPath() (string, error) {
	p, err := joinBasePath(lsDirTenants)
	if err != nil {
		return "", err
	}
	return getOrCreatePath(p)
}

// GetTenantPath returns: <base config dir>/tenants/{tenantID}
// - macOS:   $HOME/Library/Application Support/hyddenlabs/mesh/tenants/{tenantID}
// - Linux:   $XDG_CONFIG_HOME/hyddenlabs/mesh/tenants/{tenantID} or $HOME/.config/hyddenlabs/mesh/tenants/{tenantID}
// - Windows: %AppData%\hyddenlabs\mesh\tenants\{tenantID}
func GetTenantPath(tenantID uuid.UUID) (string, error) {
	base, err := GetTenantRootPath()
	if err != nil {
		return "", err
	}
	return getOrCreatePath(filepath.Join(base, tenantID.String()))
}

// GetDataPath returns: <base config dir>/tenants/{tenantID}/data
// - macOS:   $HOME/Library/Application Support/hyddenlabs/mesh/tenants/{tenantID}/data
// - Linux:   $XDG_CONFIG_HOME/hyddenlabs/mesh/tenants/{tenantID}/data or $HOME/.config/hyddenlabs/mesh/tenants/{tenantID}/data
// - Windows: %AppData%\hyddenlabs\mesh\tenants\{tenantID}\data
func GetDataPath(tenantID uuid.UUID) (string, error) {
	tenantPath, err := GetTenantPath(tenantID)
	if err != nil {
		return "", err
	}
	return getOrCreatePath(filepath.Join(tenantPath, lsDirData))
}

// GetProvidersBasePath returns: <base config dir>/providers
// - macOS:   $HOME/Library/Application Support/hyddenlabs/mesh/providers
// - Linux:   $XDG_CONFIG_HOME/hyddenlabs/mesh/providers or $HOME/.config/hyddenlabs/mesh/providers
// - Windows: %AppData%\hyddenlabs\mesh\providers
func GetProvidersBasePath() (string, error) {
	base, err := GetBasePath()
	if err != nil {
		return "", err
	}
	return getOrCreatePath(filepath.Join(base, lsDirProviders))
}

// GetProviderBasePath returns: <base config dir>/providers/{providerID}
// Useful for provider-scoped storage that is not tenant-specific.
func GetProviderBasePath(providerID uuid.UUID) (string, error) {
	providersBase, err := GetProvidersBasePath()
	if err != nil {
		return "", err
	}
	return getOrCreatePath(filepath.Join(providersBase, providerID.String()))
}

// GetCredsPath returns the path to the .creds file under <base config dir>/tenants/{tenantID}
// - macOS:   $HOME/Library/Application Support/hyddenlabs/mesh/tenants/{tenantID}/.creds
// - Linux:   $XDG_CONFIG_HOME/hyddenlabs/mesh/tenants/{tenantID}/.creds or $HOME/.config/hyddenlabs/mesh/tenants/{tenantID}/.creds
// - Windows: %AppData%\hyddenlabs\mesh\tenants\{tenantID}\.creds
func GetCredsPath(tenantID uuid.UUID) (string, error) {
	tenantPath, err := GetTenantPath(tenantID)
	if err != nil {
		return "", err
	}
	return filepath.Join(tenantPath, lsCredsFile), nil
}

// getOrCreatePath ensures the given path exists and returns it.
func getOrCreatePath(path string) (string, error) {
	if err := os.MkdirAll(path, 0700); err != nil {
		return "", err
	}
	return path, nil
}

// joinBasePath joins the base config dir with the provided relative parts.
// It does not create the path on disk (use getOrCreatePath for that).
func joinBasePath(parts ...string) (string, error) {
	base, err := GetBasePath()
	if err != nil {
		return "", err
	}
	all := append([]string{base}, parts...)
	return filepath.Join(all...), nil
}
