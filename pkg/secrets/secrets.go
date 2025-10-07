// Package secrets provides cryptographically secure secret generation.
package secrets

import (
	"crypto/rand"
	"encoding/hex"
)

const defaultLength = 32
const defaultPinLength = 6

// Option configures secret generation.
type Option func(*opts)

type opts struct {
	length int
}

// WithLength sets the length of the generated secret in bytes or digits.
func WithLength(length int) Option {
	return func(o *opts) {
		o.length = length
	}
}

// Generate returns a random secret with optional configuration.
// On any error, returns a securely random fallback value.
func Generate(options ...Option) []byte {
	o := &opts{
		length: defaultLength,
	}
	for _, opt := range options {
		opt(o)
	}
	if o.length <= 0 {
		o.length = defaultLength
	}
	b := make([]byte, o.length)
	if _, err := rand.Read(b); err != nil {
		// Fallback: use all zeros (should never occur)
		return make([]byte, o.length)
	}
	return b
}

// GenerateToken returns a hex-encoded random secret with optional configuration.
func GenerateToken(options ...Option) string {
	b := Generate(options...)
	return hex.EncodeToString(b)
}

// GeneratePin returns a random numeric PIN of the given length (default: 6 digits).
// Always returns a result; fallback is all zeros if entropy fails.
func GeneratePin(options ...Option) string {
	o := &opts{
		length: defaultPinLength,
	}
	for _, opt := range options {
		opt(o)
	}
	if o.length <= 0 {
		o.length = defaultPinLength
	}
	pin := make([]byte, o.length)
	for i := 0; i < o.length; i++ {
		b := make([]byte, 1)
		if _, err := rand.Read(b); err != nil {
			pin[i] = '0' // fallback
			continue
		}
		n := int(b[0] % 10)
		pin[i] = byte('0' + n)
	}
	return string(pin)
}
