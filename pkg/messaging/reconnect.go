package messaging

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/fgrzl/messaging"
)

// RequestHandlerMonitor periodically refreshes NATS subscriptions to ensure they remain valid.
// This uses a proactive approach - recreating all subscriptions every cycle rather than
// relying on health checks, since the Subscription interface doesn't expose connection state.
type RequestHandlerMonitor struct {
	ctx      context.Context
	cancel   context.CancelFunc
	bus      messaging.MessageBus
	mu       sync.RWMutex
	handlers map[string]handlerInfo
}

type handlerInfo struct {
	route   messaging.Route
	handler messaging.RequestHandler
	sub     messaging.Subscription
}

// NewRequestHandlerMonitor creates a monitor that proactively refreshes request handlers.
func NewRequestHandlerMonitor(ctx context.Context, bus messaging.MessageBus) *RequestHandlerMonitor {
	ctx, cancel := context.WithCancel(ctx)
	m := &RequestHandlerMonitor{
		ctx:      ctx,
		cancel:   cancel,
		bus:      bus,
		handlers: make(map[string]handlerInfo),
	}
	go m.monitor()
	return m
}

// Register adds a request handler to be monitored.
func (m *RequestHandlerMonitor) Register(route messaging.Route, handler messaging.RequestHandler) (messaging.Subscription, error) {
	sub, err := m.bus.SubscribeRequest(route, handler)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.handlers[route.String()] = handlerInfo{
		route:   route,
		handler: handler,
		sub:     sub,
	}
	m.mu.Unlock()

	slog.InfoContext(m.ctx, "registered monitored request handler", "route", route.String())
	return sub, nil
}

// monitor proactively refreshes all handler subscriptions every cycle.
// This approach handles cases where subscriptions die silently without any way to detect it.
func (m *RequestHandlerMonitor) monitor() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.mu.Lock()
			// Proactively recreate all subscriptions to ensure they're valid
			// This handles cases where the underlying connection died silently
			for key, info := range m.handlers {
				// Unsubscribe the old one (may fail if already dead, that's ok)
				if info.sub != nil {
					_ = info.sub.Unsubscribe()
				}

				// Create fresh subscription
				newSub, err := m.bus.SubscribeRequest(info.route, info.handler)
				if err != nil {
					slog.ErrorContext(m.ctx, "failed to refresh request handler subscription",
						"route", info.route.String(),
						"error", err)
					continue
				}

				// Update with new subscription
				m.handlers[key] = handlerInfo{
					route:   info.route,
					handler: info.handler,
					sub:     newSub,
				}
				slog.DebugContext(m.ctx, "refreshed request handler subscription", "route", info.route.String())
			}
			m.mu.Unlock()
		}
	}
}

// Stop stops the monitor and unsubscribes all handlers.
func (m *RequestHandlerMonitor) Stop() {
	m.cancel()

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, info := range m.handlers {
		if info.sub != nil {
			_ = info.sub.Unsubscribe()
		}
	}
	m.handlers = make(map[string]handlerInfo)
}
