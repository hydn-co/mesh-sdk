package testkit

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-sdk/pkg/tenantkit"
	"github.com/stretchr/testify/assert"
)

// InvokeDescribe runs the runner with -describe and fails the test if the
// runner returns an error. It sets up a temporary stream path in the test
// context before invoking the runner.
func InvokeDescribe(t *testing.T, manifest *runner.Manifest) {
	t.Helper()
	ctx := tenantkit.WithStreamPath(t.Context(), t.TempDir())

	err := runner.RunWithArgs(ctx, manifest, "-describe")
	assert.NoError(t, err)
}

// InvokeList runs the runner with -list, captures stdout, and asserts that the
// manifest capabilities are present in the output. This helper is intended for
// test assertions and will fail the test on error.
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

// InvokeGenerate runs the runner to generate artifacts for the given capability
// using the provided connector and secret IDs. It fails the test on error.
// This is a test helper used by integration-style tests.
//
// Example usage:
//
//	testkit.InvokeGenerate(t, manifest, "capability-id", connectorID, secretID)
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

// InvokeRun runs the runner for the given capability and connector using a mock tenant.
// It fails the test on error. This helper is useful for smoke tests of runtime behavior.
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
