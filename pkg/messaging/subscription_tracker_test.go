package messaging

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSubscription implements messaging.Subscription for testing
type mockSubscription struct {
	id                uuid.UUID
	unsubscribeCalled bool
	unsubscribeErr    error
}

func (m *mockSubscription) GetID() uuid.UUID {
	return m.id
}

func (m *mockSubscription) Unsubscribe() error {
	m.unsubscribeCalled = true
	return m.unsubscribeErr
}

func TestShouldTrackSubscription(t *testing.T) {
	// Arrange
	tracker := NewSubscriptionTracker()
	sub := &mockSubscription{id: uuid.New()}

	// Act
	tracker.Track(sub)

	// Assert - verify by untracking (which calls Unsubscribe)
	tracker.Untrack(sub)
	assert.True(t, sub.unsubscribeCalled)
}

func TestShouldUntrackSubscriptionAndUnsubscribe(t *testing.T) {
	// Arrange
	tracker := NewSubscriptionTracker()
	sub := &mockSubscription{id: uuid.New()}
	tracker.Track(sub)

	// Act
	tracker.Untrack(sub)

	// Assert
	assert.True(t, sub.unsubscribeCalled)
}

func TestShouldShutdownAllSubscriptions(t *testing.T) {
	// Arrange
	tracker := NewSubscriptionTracker()
	sub1 := &mockSubscription{id: uuid.New()}
	sub2 := &mockSubscription{id: uuid.New()}
	tracker.Track(sub1)
	tracker.Track(sub2)
	ctx := context.Background()

	// Act
	err := tracker.Shutdown(ctx)

	// Assert
	require.NoError(t, err)
	assert.True(t, sub1.unsubscribeCalled)
	assert.True(t, sub2.unsubscribeCalled)
}

func TestShouldReturnErrorWhenShutdownFails(t *testing.T) {
	// Arrange
	tracker := NewSubscriptionTracker()
	expectedErr := errors.New("unsubscribe failed")
	sub := &mockSubscription{
		id:             uuid.New(),
		unsubscribeErr: expectedErr,
	}
	tracker.Track(sub)
	ctx := context.Background()

	// Act
	err := tracker.Shutdown(ctx)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unsubscribe")
	assert.True(t, sub.unsubscribeCalled)
}

func TestShouldHandleEmptyTrackerShutdown(t *testing.T) {
	// Arrange
	tracker := NewSubscriptionTracker()
	ctx := context.Background()

	// Act
	err := tracker.Shutdown(ctx)

	// Assert
	assert.NoError(t, err)
}

func TestShouldOverwriteTrackedSubscriptionWithSameID(t *testing.T) {
	// Arrange
	tracker := NewSubscriptionTracker()
	sharedID := uuid.New()
	sub1 := &mockSubscription{id: sharedID}
	sub2 := &mockSubscription{id: sharedID}
	tracker.Track(sub1)

	// Act
	tracker.Track(sub2)
	tracker.Untrack(sub2)

	// Assert
	assert.False(t, sub1.unsubscribeCalled, "first subscription should not be unsubscribed")
	assert.True(t, sub2.unsubscribeCalled, "second subscription should be unsubscribed")
}
