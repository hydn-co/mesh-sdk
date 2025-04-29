package events

import (
	"github.com/fgrzl/es"
	"github.com/fgrzl/messaging"
)

var (
	_ es.DomainEvent = &SlugReserved{}
	_ es.DomainEvent = &SlugSurrendered{}
)

func init() {
	es.Register(func() *SlugReserved { return &SlugReserved{} })
	es.Register(func() *SlugSurrendered { return &SlugSurrendered{} })
}

// SlugReserved event
type SlugReserved struct {
	es.DomainEventBase
	Parent es.Entity `json:"parent"`
	Owner  es.Entity `json:"owner"`
	Slug   string    `json:"slug"`
}

func (e *SlugReserved) GetDiscriminator() string {
	return "hydn://domain/events/slugs/reserved"
}

func (e *SlugReserved) GetRoute() messaging.Route {
	return messaging.NewInternalRoute("slugs", "reserved")
}

// SlugSurrendered event
type SlugSurrendered struct {
	es.DomainEventBase
	Parent es.Entity `json:"parent"`
	Owner  es.Entity `json:"owner"`
	Slug   string    `json:"slug"`
}

func (e *SlugSurrendered) GetDiscriminator() string {
	return "hydn://domain/events/slugs/surrendered"
}

func (e *SlugSurrendered) GetRoute() messaging.Route {
	return messaging.NewInternalRoute("slugs", "surrendered")
}
