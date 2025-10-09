package processors

import (
	"context"
	"fmt"

	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
)

// SubscriptionTracker defines lifecycle methods for managing messaging.Subscriptions.
type SubscriptionTracker interface {
	Track(messaging.Subscription)
	Untrack(messaging.Subscription)
	Shutdown(context.Context) error
}

type subscriptionTracker struct {
	subs map[uuid.UUID]messaging.Subscription
}

func NewSubscriptionTracker() SubscriptionTracker {
	return &subscriptionTracker{
		subs: make(map[uuid.UUID]messaging.Subscription),
	}
}

func (s *subscriptionTracker) Track(sub messaging.Subscription) {
	s.subs[sub.GetID()] = sub
}

func (s *subscriptionTracker) Untrack(sub messaging.Subscription) {
	delete(s.subs, sub.GetID())
	_ = sub.Unsubscribe()
}

func (s *subscriptionTracker) Shutdown(ctx context.Context) error {
	for id, sub := range s.subs {
		if err := sub.Unsubscribe(); err != nil {
			return fmt.Errorf("failed to unsubscribe from %s subscription %s", id, err)
		}
		delete(s.subs, id)
	}
	return nil
}
