package spawn

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type ExeSpawner struct{}

func NewExeSpawner() (Spawner, error) {
	return &ExeSpawner{}, nil
}

func (s *ExeSpawner) Run(ctx context.Context, config SpawnConfig) error {
	cmd := exec.CommandContext(ctx, config.Image, config.Args...)
	cmd.Env = os.Environ()
	for k, v := range config.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return cmd.Run()
}
