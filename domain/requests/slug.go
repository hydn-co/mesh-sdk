package requests

import (
	"github.com/fgrzl/es"
	"github.com/fgrzl/messaging"
)

type GetSlug struct {
	Parent es.Entity `json:"parent"`
	Slug   string    `json:"slug"`
}

// GetDiscriminator implements messaging.Request.
func (r *GetSlug) GetDiscriminator() string {
	return "query://slug/get"
}

type CheckSlugAvailability struct {
	Parent es.Entity `json:"parent"`
	Slug   string    `json:"slug"`
}

func (c *CheckSlugAvailability) GetDiscriminator() string {
	return "query://slug/availability"
}

type ReserveSlug struct {
	Parent es.Entity `json:"parent"`
	Owner  es.Entity `json:"owner"`
	Slug   string    `json:"slug"`
}

// GetDiscriminator implements messaging.Request.
func (r *ReserveSlug) GetDiscriminator() string {
	return "command://slug/reserve"
}

// GetRoute implements messaging.Request.
func (r *ReserveSlug) GetRoute() messaging.Route {
	return messaging.NewInternalRoute("slug", "reserve")
}

type SurrenderSlug struct {
	Parent es.Entity `json:"parent"`
	Owner  es.Entity `json:"owner"`
	Slug   string    `json:"slug"`
}

// GetDiscriminator implements messaging.Request.
func (s *SurrenderSlug) GetDiscriminator() string {
	return "command://slug/surrender"
}

// GetRoute implements messaging.Request.
func (s *SurrenderSlug) GetRoute() messaging.Route {
	return messaging.NewInternalRoute("slug", "surrender")
}
