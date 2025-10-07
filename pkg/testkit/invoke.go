package testkit

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-sdk/pkg/tenantkit"
	"github.com/stretchr/testify/assert"
)

func InvokeDescribe(t *testing.T, manifest *runner.Manifest) {
	t.Helper()
	ctx := tenantkit.WithStreamPath(t.Context(), t.TempDir())

	err := runner.RunWithArgs(ctx, manifest, "-describe")
	assert.NoError(t, err)
}

func InvokeList(t *testing.T, manifest *runner.Manifest) {
	t.Helper()
	ctx := tenantkit.WithStreamPath(t.Context(), t.TempDir())

	// derive expected list from manifest keys
	var expected []string
	for name := range manifest.Capabilities {
		expected = append(expected, name)
	}

	output, err := CaptureOutput(func() error {
		return runner.RunWithArgs(ctx, manifest, "-list")
	})
	assert.NoError(t, err)

	for _, capID := range expected {
		assert.Contains(t, output, capID, "expected capability ID %q in output", capID)
	}
}

func InvokeGenerate(t *testing.T, manifest *runner.Manifest, capability, connectorID, secretID string) {
	t.Helper()
	ctx := tenantkit.WithStreamPath(t.Context(), t.TempDir())

	err := runner.RunWithArgs(ctx, manifest,
		"-generate",
		"-capability", capability,
		"-connector-id", connectorID,
		"-secret-id", secretID,
	)
	assert.NoError(t, err)
}

func InvokeRun(t *testing.T, manifest *runner.Manifest, capability, connectorID string) {
	t.Helper()
	ctx := tenantkit.WithStreamPath(t.Context(), t.TempDir())
	err := runner.RunWithArgs(ctx, manifest,
		"-run",
		"-capability", capability,
		"-connector-id", connectorID,
		"-tenant-id", uuid.Nil.String(),
		"-mock",
	)
	assert.NoError(t, err)
}
