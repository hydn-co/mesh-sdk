package spawn

import (
	"context"
	"fmt"
)

type SpawnMode string

const (
	SpawnModeExe SpawnMode = "exe"
	SpawnModeACS SpawnMode = "acs"
	SpawnModeK8s SpawnMode = "k8s"
	SpawnModeECS SpawnMode = "ecs"
)

func (m SpawnMode) IsValid() bool {
	return m == SpawnModeExe || m == SpawnModeACS || m == SpawnModeK8s || m == SpawnModeECS
}

type SpawnConfig struct {
	JobID       string
	Image       string // image name or binary path
	Args        []string
	Env         map[string]string
	Integration string
	Version     string
}

type Spawner interface {
	Run(ctx context.Context, config SpawnConfig) error
}

func GetSpawner(mode SpawnMode) (Spawner, error) {
	switch mode {
	case SpawnModeExe:
		return NewExeSpawner()
	// case SpawnModeACS:
	// 	return &ACASpawner{}, nil
	// case SpawnModeK8s:
	// 	return &K8sSpawner{}, nil
	// case SpawnModeECS:
	// 	return &ECSSpawner{}, nil
	default:
		return nil, fmt.Errorf("unsupported SPAWN_MODE: %s", mode)
	}
}
