package lifecycle

import (
	"context"
	"errors"
)

// Host manages the lifecycle of multiple services. It starts services and
// records them so they can later be stopped as a group.
type Host interface {
	StartServices(ctx context.Context, services ...Service) error
	Stop(ctx context.Context) error
}

type defaultHost struct {
	services []Service
}

// NewHost returns a new, empty lifecycle Host.
func NewHost() Host {
	m := &defaultHost{
		services: make([]Service, 0),
	}
	return m
}

// StartServices starts each provided Service in order. If any service fails to
// start, the error is returned immediately and previously started services are
// left running.
func (h *defaultHost) StartServices(ctx context.Context, services ...Service) error {

	for _, service := range services {

		if err := service.Start(ctx); err != nil {
			return err
		}

		h.services = append(h.services, service)
	}
	return nil
}

// Stop stops all previously started services and aggregates any errors.
func (m *defaultHost) Stop(ctx context.Context) error {
	var errs []error

	for _, srv := range m.services {
		if err := srv.Stop(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
