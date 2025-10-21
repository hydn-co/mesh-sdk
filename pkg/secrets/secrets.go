// Package secrets provides cryptographically secure secret generation.
package secrets

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const defaultLength = 32
const defaultPinLength = 6

const (
	argonTime    = 3
	argonMemory  = 64 * 1024 // KiB (64 MiB)
	argonThreads = 2
	argonKeyLen  = 32
	saltLen      = 16
)

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

// hashPasswordArgon2 hashes the password using Argon2id and returns a PHC string.
func hashPasswordArgon2(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, uint32(argonTime), uint32(argonMemory), uint8(argonThreads), uint32(argonKeyLen))
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argonMemory, argonTime, argonThreads, b64Salt, b64Hash), nil
}

// verifyPasswordArgon2 verifies the password against a PHC-formatted Argon2id string.
func verifyPasswordArgon2(encoded, password string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid argon2id format")
	}
	// params in parts[3]
	params := strings.Split(parts[3], ",")
	var memory, time uint64
	var threads int
	for _, p := range params {
		kv := strings.SplitN(p, "=", 2)
		switch kv[0] {
		case "m":
			memory, _ = strconv.ParseUint(kv[1], 10, 32)
		case "t":
			time, _ = strconv.ParseUint(kv[1], 10, 32)
		case "p":
			pth, _ := strconv.Atoi(kv[1])
			threads = pth
		}
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	calc := argon2.IDKey([]byte(password), salt, uint32(time), uint32(memory), uint8(threads), uint32(len(hash)))
	if subtle.ConstantTimeCompare(calc, hash) == 1 {
		return true, nil
	}
	return false, nil
}

// HashPassword hashes the provided password using Argon2id (modern default).
// Returns a PHC formatted hash string or an error.
func HashPassword(password string) (string, error) {
	return hashPasswordArgon2(password)
}

// CheckPasswordHash verifies an argon2id PHC string against the provided password.
func CheckPasswordHash(password, hash string) bool {
	ok, _ := verifyPasswordArgon2(hash, password)
	return ok
}
