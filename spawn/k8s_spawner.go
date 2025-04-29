//go:build k8s

package spawn

import "context"

type K8sSpawner struct{}

func newK8sSpawner() (Spawner, error) {
	return &K8sSpawner{}, nil
}

func (s *K8sSpawner) Run(ctx context.Context, config SpawnConfig) error {
	// TODO: Implement K8s logic using Azure SDK
	return nil
}
