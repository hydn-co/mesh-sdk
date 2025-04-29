//go:build ecs

package spawn

import "context"

type ECSSpawner struct{}

func newECSSpawner() (Spawner, error) {
	return &ECSSpawner{}, nil
}

func (s *ECSSpawner) Run(ctx context.Context, config SpawnConfig) error {
	// TODO: Implement ECS logic using Azure SDK
	return nil
}
