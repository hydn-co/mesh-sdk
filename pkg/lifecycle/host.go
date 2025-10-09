package lifecycle

import (
	"context"
	"errors"
	"time"
)

// Host manages the lifecycle of multiple services.
//
// New behavior:
//   - Start(ctx) starts all managed services in order and must be non-blocking.
//   - Stop(ctx) stops all previously started services (in reverse order) and may block while waiting for shutdown.
//   - Run(ctx) is a blocking convenience that starts services, waits until ctx is canceled,
//     and then attempts a graceful shutdown (Stop) before returning.
type Host interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Run(ctx context.Context) error
}

type defaultHost struct {
	services []Service
}

// NewHost returns a new lifecycle Host that will manage the provided services.
// Services passed to NewHost are the ones the Host will Start/Stop/Run.
func NewHost(services ...Service) Host {
	m := &defaultHost{
		services: make([]Service, 0, len(services)),
	}
	m.services = append(m.services, services...)
	return m
}

// Start starts each managed Service in order. Start must be non-blocking per Service convention.
// If a service fails to start, the error is returned and previously started services are left running.
func (h *defaultHost) Start(ctx context.Context) error {
	for _, service := range h.services {
		if err := service.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Stop stops all previously started services in reverse order and aggregates any errors.
func (h *defaultHost) Stop(ctx context.Context) error {
	var errs []error
	for i := len(h.services) - 1; i >= 0; i-- {
		if err := h.services[i].Stop(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Run starts the services and blocks until the provided context is canceled.
// After cancellation it requests a graceful shutdown via Stop with a 30s timeout.
func (h *defaultHost) Run(ctx context.Context) error {
	// start with a background context for services; Run controls their lifetime
	svcCtx, svcCancel := context.WithCancel(context.Background())
	defer svcCancel()

	if err := h.Start(svcCtx); err != nil {
		return err
	}

	// Block until ctx is canceled
	<-ctx.Done()

	// Attempt graceful shutdown with a timeout
	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return h.Stop(stopCtx)
}
