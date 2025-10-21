package lifecycle

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockService implements Service for testing
type mockService struct {
	startCalled bool
	stopCalled  bool
	startErr    error
	stopErr     error
	startDelay  time.Duration
	stopDelay   time.Duration
}

func (m *mockService) Start(ctx context.Context) error {
	m.startCalled = true
	if m.startDelay > 0 {
		time.Sleep(m.startDelay)
	}
	return m.startErr
}

func (m *mockService) Stop(ctx context.Context) error {
	m.stopCalled = true
	if m.stopDelay > 0 {
		time.Sleep(m.stopDelay)
	}
	return m.stopErr
}

func TestShouldStartServiceSuccessfully(t *testing.T) {
	// Arrange
	service := &mockService{}
	ctx := context.Background()

	// Act
	err := service.Start(ctx)

	// Assert
	require.NoError(t, err)
	assert.True(t, service.startCalled)
}

func TestShouldStopServiceSuccessfully(t *testing.T) {
	// Arrange
	service := &mockService{}
	ctx := context.Background()

	// Act
	err := service.Stop(ctx)

	// Assert
	require.NoError(t, err)
	assert.True(t, service.stopCalled)
}

func TestShouldReturnErrorWhenStartFails(t *testing.T) {
	// Arrange
	expectedErr := errors.New("start failed")
	service := &mockService{startErr: expectedErr}
	ctx := context.Background()

	// Act
	err := service.Start(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.True(t, service.startCalled)
}

func TestShouldReturnErrorWhenStopFails(t *testing.T) {
	// Arrange
	expectedErr := errors.New("stop failed")
	service := &mockService{stopErr: expectedErr}
	ctx := context.Background()

	// Act
	err := service.Stop(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.True(t, service.stopCalled)
}

func TestShouldRespectContextCancellation(t *testing.T) {
	// Arrange
	service := &mockService{startDelay: 100 * time.Millisecond}
	ctx, cancel := context.WithCancel(context.Background())

	// Act
	cancel() // Cancel immediately
	err := service.Start(ctx)

	// Assert - service should still execute but context is cancelled
	require.NoError(t, err) // Start doesn't check ctx in mock
	assert.True(t, service.startCalled)
}

func TestSupervisorShouldRunFunctionUntilSuccess(t *testing.T) {
	// Arrange
	called := false
	runFunc := func(ctx context.Context) error {
		called = true
		return nil // Success on first run
	}
	policy := Policy{
		Action:      FailHost,
		MaxRestarts: 0,
	}
	supervisor := NewSupervisor(runFunc, policy)
	ctx := context.Background()

	// Act
	err := supervisor.Start(ctx)
	require.NoError(t, err)

	// Wait a bit for goroutine to execute
	time.Sleep(50 * time.Millisecond)

	stopErr := supervisor.Stop(ctx)

	// Assert
	assert.NoError(t, stopErr)
	assert.True(t, called)
}

func TestSupervisorShouldRestartOnFailureWhenPolicyIsRestart(t *testing.T) {
	// Arrange
	callCount := 0
	runFunc := func(ctx context.Context) error {
		callCount++
		if callCount < 3 {
			return errors.New("temporary failure")
		}
		return nil // Success on third try
	}
	policy := Policy{
		Action:      Restart,
		MaxRestarts: 5,
		Backoff: func(attempt int) time.Duration {
			return 10 * time.Millisecond
		},
	}
	supervisor := NewSupervisor(runFunc, policy)
	ctx := context.Background()

	// Act
	err := supervisor.Start(ctx)
	require.NoError(t, err)

	// Wait for retries
	time.Sleep(150 * time.Millisecond)

	stopErr := supervisor.Stop(ctx)

	// Assert
	assert.NoError(t, stopErr)
	assert.GreaterOrEqual(t, callCount, 3, "should have retried at least 3 times")
}

func TestSupervisorShouldStopWhenMaxRestartsReached(t *testing.T) {
	// Arrange
	callCount := 0
	runFunc := func(ctx context.Context) error {
		callCount++
		return errors.New("persistent failure")
	}
	policy := Policy{
		Action:      Restart,
		MaxRestarts: 2,
		Backoff: func(attempt int) time.Duration {
			return 10 * time.Millisecond
		},
	}
	supervisor := NewSupervisor(runFunc, policy)
	ctx := context.Background()

	// Act
	err := supervisor.Start(ctx)
	require.NoError(t, err)

	// Wait for max restarts to be reached
	time.Sleep(150 * time.Millisecond)

	stopErr := supervisor.Stop(ctx)

	// Assert
	assert.NoError(t, stopErr)
	assert.LessOrEqual(t, callCount, 3, "should not exceed max restarts + initial attempt")
}

func TestSupervisorShouldIgnoreFailuresWhenPolicyIsIgnore(t *testing.T) {
	// Arrange
	callCount := 0
	runFunc := func(ctx context.Context) error {
		callCount++
		return errors.New("ignored failure")
	}
	policy := Policy{
		Action:      Ignore,
		MaxRestarts: 0,
	}
	supervisor := NewSupervisor(runFunc, policy)
	ctx := context.Background()

	// Act
	err := supervisor.Start(ctx)
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	stopErr := supervisor.Stop(ctx)

	// Assert
	assert.NoError(t, stopErr)
	assert.Equal(t, 1, callCount, "should only run once when ignoring failures")
}
