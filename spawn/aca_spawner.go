//go:build acs

package spawn

import "context"

type ACASpawner struct{}

func newACSSpawner() (Spawner, error) {
	return &ACASpawner{}, nil
}

func (s *ACASpawner) Run(ctx context.Context, config SpawnConfig) error {
	// TODO: Implement ACS logic using Azure SDK
	return nil
}
