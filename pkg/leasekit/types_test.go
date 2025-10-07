package leasekit

import (
	"testing"

	"github.com/hydn-co/mesh-sdk/pkg/testkit"
)

// make sure we register all leasekit types and use consistent discriminators
// keep this list up to date with the types in the leasekit package
// and lets keep it sorted by the discriminator
var leasekitTypes = map[string]any{
	"mesh://leasekit/acquire": &Acquire{},
	"mesh://leasekit/lease":   &Lease{},
	"mesh://leasekit/release": &Release{},
	"mesh://leasekit/renew":   &Renew{},
}

func TestShouldRegisterPolymorphicLeasekitTypes(t *testing.T) {
	testkit.TestPolymorphicRegistrations(t, leasekitTypes)
}
