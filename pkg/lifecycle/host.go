package lifecycle

import (
	"context"
	"errors"
)

type Host interface {
	StartServices(ctx context.Context, services ...Service) error
	Stop(ctx context.Context) error
}

type defaultHost struct {
	services []Service
}

func NewHost() Host {
	m := &defaultHost{
		services: make([]Service, 0),
	}
	return m
}

// Start will attempt to start the service or return an error
func (h *defaultHost) StartServices(ctx context.Context, services ...Service) error {

	for _, service := range services {

		if err := service.Start(ctx); err != nil {
			return err
		}

		h.services = append(h.services, service)
	}
	return nil
}

func (m *defaultHost) Stop(ctx context.Context) error {
	var errs []error

	for _, srv := range m.services {
		if err := srv.Stop(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
