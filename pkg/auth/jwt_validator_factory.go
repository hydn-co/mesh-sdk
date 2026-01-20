package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/fgrzl/claims/jwtkit"

	"github.com/hydn-co/mesh-sdk/pkg/env"
)

var (
	cachedPublicKey *rsa.PublicKey
	keyMu           sync.RWMutex
)

// GetValidator returns a jwtkit.Validator backed by the auth-server's JWKS.
// The function caches the first-found RSA public key in memory. For production
// use-cases a more complete JWKS client (kid-aware, refresh-on-expiry) is
// recommended; this function implements a minimal, pragmatic approach for the
// SDK.
func GetValidator() (jwtkit.Validator, error) {
	return GetValidatorWithClient(nil)
}

// GetValidatorWithClient returns a jwtkit.Validator backed by the auth-server's JWKS.
// It accepts an optional HTTP client for custom transport configuration (e.g., TLS settings).
// If client is nil, http.DefaultClient will be used.
func GetValidatorWithClient(client *http.Client) (jwtkit.Validator, error) {
	if pk := getCachedPublicKey(); pk != nil {
		return &jwtkit.RSAValidator{PublicKey: pk}, nil
	}

	// Fetch JWKS from auth server well-known endpoint.
	base := env.GetEnvOrDefaultStr(env.MeshAuthBaseURL, "http://localhost:8080")
	jwksURL := strings.TrimRight(base, "/") + "/.well-known/jwks.json"
	pk, err := fetchFirstRSAFromJWKS(jwksURL, client)
	if err != nil {
		return nil, fmt.Errorf("fetch jwks from %s: %w", jwksURL, err)
	}
	if pk == nil {
		return nil, errors.New("no RSA key available from jwks")
	}
	setCachedPublicKey(pk)
	return &jwtkit.RSAValidator{PublicKey: pk}, nil
}

func getCachedPublicKey() *rsa.PublicKey {
	keyMu.RLock()
	pk := cachedPublicKey
	keyMu.RUnlock()
	return pk
}

func setCachedPublicKey(pk *rsa.PublicKey) {
	keyMu.Lock()
	cachedPublicKey = pk
	keyMu.Unlock()
}

// fetchFirstRSAFromJWKS fetches a JWKS from the given URL and returns the
// first RSA public key it contains. This helper is intentionally simple and
// returns the first RSA key found; callers that require key rotation or
// `kid` selection should replace this with a full JWKS client.
// If client is nil, http.DefaultClient will be used.
func fetchFirstRSAFromJWKS(jwksURL string, client *http.Client) (*rsa.PublicKey, error) {
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch jwks: non-200 status")
	}

	var doc struct {
		Keys []struct {
			Kty string `json:"kty"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, err
	}

	for _, k := range doc.Keys {
		if strings.ToUpper(k.Kty) != "RSA" {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			// try standard base64 as a fallback
			nBytes, err = base64.StdEncoding.DecodeString(k.N)
			if err != nil {
				continue
			}
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil {
			eBytes, err = base64.StdEncoding.DecodeString(k.E)
			if err != nil {
				continue
			}
		}
		// convert eBytes to int
		e := 0
		for _, b := range eBytes {
			e = e<<8 + int(b)
		}
		pub := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: e,
		}
		return pub, nil
	}
	return nil, errors.New("no RSA key found in jwks")
}
