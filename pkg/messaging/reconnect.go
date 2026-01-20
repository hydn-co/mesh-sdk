package messaging

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
)

// RequestHandlerMonitor periodically tests NATS connectivity and recreates subscriptions if needed.
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

// NewRequestHandlerMonitor creates a monitor that watches and reconnects request handlers.
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

// monitor periodically checks handler subscriptions and reconnects if needed.
func (m *RequestHandlerMonitor) monitor() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	consecutiveFailures := 0
	maxFailures := 5

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.mu.Lock()
			for key, info := range m.handlers {
				// Test if we can still communicate with NATS by trying to create a test subscription
				testRoute := messaging.NewInternalRoute("_health", uuid.NewString())
				testSub, err := m.bus.SubscribeRequest(testRoute, func(ctx context.Context, req messaging.Request) (messaging.Response, error) {
					return nil, nil
				})

				if err != nil {
					consecutiveFailures++
					slog.WarnContext(m.ctx, "NATS health check failed, connection may be down",
						"consecutive_failures", consecutiveFailures,
						"error", err)

					if consecutiveFailures >= maxFailures {
						slog.ErrorContext(m.ctx, "attempting to reconnect request handlers after sustained failures")
						m.reconnectHandler(key, info)
						consecutiveFailures = 0
					}
				} else {
					// Health check passed, clean up test subscription
					_ = testSub.Unsubscribe()
					if consecutiveFailures > 0 {
						consecutiveFailures = 0
						slog.InfoContext(m.ctx, "NATS health check recovered")
					}
				}
			}
			m.mu.Unlock()
		}
	}
}

// reconnectHandler attempts to recreate a failed subscription.
func (m *RequestHandlerMonitor) reconnectHandler(key string, info handlerInfo) {
	// Unsubscribe the old one if possible
	if info.sub != nil {
		_ = info.sub.Unsubscribe()
	}

	// Try to create new subscription
	newSub, err := m.bus.SubscribeRequest(info.route, info.handler)
	if err != nil {
		slog.ErrorContext(m.ctx, "failed to reconnect request handler",
			"route", info.route.String(),
			"error", err)
		// Keep trying on next monitor cycle
		return
	}

	// Update with new subscription
	m.handlers[key] = handlerInfo{
		route:   info.route,
		handler: info.handler,
		sub:     newSub,
	}

	slog.InfoContext(m.ctx, "successfully reconnected request handler", "route", info.route.String())
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
