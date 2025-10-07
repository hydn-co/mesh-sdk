package factories

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

func GetValidator() (jwtkit.Validator, error) {
	if pk := getCachedPublicKey(); pk != nil {
		return &jwtkit.RSAValidator{PublicKey: pk}, nil
	}

	// Fetch JWKS from auth server well-known endpoint
	portal := env.GetEnvOrDefaultStr(env.MeshPortalURL, "http://localhost:8080")
	jwksURL := strings.TrimRight(portal, "/") + "/.well-known/jwks.json"
	pk, err := fetchFirstRSAFromJWKS(jwksURL)
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
// first RSA public key it contains.
func fetchFirstRSAFromJWKS(jwksURL string) (*rsa.PublicKey, error) {
	resp, err := http.Get(jwksURL)
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
