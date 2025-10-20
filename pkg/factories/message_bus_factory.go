package factories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fgrzl/messaging"
	"github.com/fgrzl/messaging/pkg/natsbus"
	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/localstore"
	"github.com/nats-io/nkeys"
)

// MessageBusOptions holds all the necessary configuration for message bus factory.
//
// Fields left zero-valued will be replaced with sensible defaults by the
// factory (for example AuthAttempts, AuthBaseDelay, MaxAuthRespBytes).
type MessageBusOptions struct {
	TenantID          uuid.UUID
	ClientCredentials localstore.ClientCredentials
	MessageBusAuthURL string
	BrokerURL         string
	// Optional HTTP client to use for auth requests. If nil, http.DefaultClient
	// will be used.
	HTTPClient *http.Client
	// Optional max attempts for auth (default 3)
	AuthAttempts int
	// Optional base delay for auth backoff (default 300ms)
	AuthBaseDelay time.Duration
	// Maximum number of bytes to read from auth response (default 64KB)
	MaxAuthRespBytes int64
}

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// NewMessageBusFactory initializes a new factory that creates or reuses a
// MessageBus instance. The returned factory is safe for concurrent use and
// will lazily establish the underlying connection when Get is called.
func NewMessageBusFactory(opts MessageBusOptions) messaging.MessageBusFactory {
	return &DefaultMessageBusFactory{
		tenantID:         opts.TenantID,
		clientID:         opts.ClientCredentials.ClientID,
		clientSecret:     opts.ClientCredentials.ClientSecret,
		user:             opts.ClientCredentials.User,
		authURL:          opts.MessageBusAuthURL,
		brokerURL:        opts.BrokerURL,
		httpClient:       opts.HTTPClient,
		authAttempts:     opts.AuthAttempts,
		authBaseDelay:    opts.AuthBaseDelay,
		maxAuthRespBytes: opts.MaxAuthRespBytes,
	}
}

type DefaultMessageBusFactory struct {
	tenantID         uuid.UUID
	clientID         uuid.UUID
	clientSecret     string
	user             nkeys.KeyPair
	authURL          string
	brokerURL        string
	httpClient       *http.Client
	authAttempts     int
	authBaseDelay    time.Duration
	maxAuthRespBytes int64

	mu  sync.Mutex
	bus messaging.MessageBus
}

// Get returns a cached or newly established message bus connection.
func (f *DefaultMessageBusFactory) Get(ctx context.Context) (messaging.MessageBus, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.bus != nil {
		return f.bus, nil
	}

	// Diagnostic: log attempt to create message bus (do not log secrets)
	slog.Debug("creating message bus", "broker_url", f.brokerURL, "auth_url", f.authURL, "tenant_id", f.tenantID, "client_id", f.clientID)

	// Retry loop with exponential backoff + jitter
	maxAttempts := 5
	baseDelay := 500 * time.Millisecond
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		slog.Debug("attempting to create message bus", "attempt", attempt)
		bus, err := natsbus.NewBus(f.brokerURL, f.fetchJWT(ctx), f.signNonce)
		if err == nil {
			f.bus = bus
			return bus, nil
		}
		lastErr = err
		slog.Warn("message bus creation attempt failed", "attempt", attempt, "err", err)
		if attempt == maxAttempts {
			slog.Error("failed to create message bus after retries", "broker_url", f.brokerURL, "err", err)
			return nil, fmt.Errorf("create message bus: %w", err)
		}
		backoff := baseDelay * time.Duration(1<<(attempt-1))
		jitter := time.Duration(rng.Int63n(int64(backoff / 2)))
		sleep := backoff + jitter
		slog.Debug("waiting before retrying message bus creation", "sleep_ms", sleep.Milliseconds())
		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if lastErr != nil {
		return nil, fmt.Errorf("create message bus: %w", lastErr)
	}
	return nil, fmt.Errorf("create message bus: unknown error")
}

func (f *DefaultMessageBusFactory) fetchJWT(ctx context.Context) func() (string, error) {
	return func() (string, error) {
		type authRequest struct {
			TenantID     uuid.UUID `json:"tenant_id"`
			ClientID     uuid.UUID `json:"client_id"`
			ClientSecret string    `json:"client_secret"`
			UserPub      string    `json:"user_public_key"`
		}

		userPub, err := f.user.PublicKey()
		if err != nil {
			return "", fmt.Errorf("get user public key: %w", err)
		}

		payload := authRequest{
			TenantID:     f.tenantID,
			ClientID:     f.clientID,
			ClientSecret: f.clientSecret,
			UserPub:      userPub,
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(payload); err != nil {
			return "", fmt.Errorf("encode auth payload: %w", err)
		}

		payloadBytes := buf.Bytes()

		// Diagnostic: log auth request metadata (no secrets)
		slog.Debug("performing auth request for message bus", "auth_url", f.authURL, "tenant_id", f.tenantID, "client_id", f.clientID)

		// Prepare retry parameters
		maxAuthAttempts := f.authAttempts
		if maxAuthAttempts <= 0 {
			maxAuthAttempts = 3
		}
		authBaseDelay := f.authBaseDelay
		if authBaseDelay <= 0 {
			authBaseDelay = 300 * time.Millisecond
		}
		maxResp := f.maxAuthRespBytes
		if maxResp <= 0 {
			maxResp = 64 * 1024 // 64 KiB
		}

		client := f.httpClient
		if client == nil {
			client = http.DefaultClient
		}

		for a := 1; a <= maxAuthAttempts; a++ {
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			// Build a fresh request each attempt
			reqReader := bytes.NewReader(payloadBytes)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.authURL, reqReader)
			if err != nil {
				return "", fmt.Errorf("create auth request: %w", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				slog.Warn("auth request attempt failed", "attempt", a, "err", err)
				if a == maxAuthAttempts {
					slog.Error("auth request failed after retries", "auth_url", f.authURL, "err", err)
					return "", fmt.Errorf("auth request failed: %w", err)
				}
				backoff := authBaseDelay * time.Duration(1<<(a-1))
				jitter := time.Duration(rng.Int63n(int64(backoff / 2)))
				sleep := backoff + jitter
				select {
				case <-time.After(sleep):
				case <-ctx.Done():
					return "", ctx.Err()
				}
				continue
			}

			limited := io.LimitReader(resp.Body, maxResp)
			bodyBytes, _ := io.ReadAll(limited)
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				trimmed := string(bodyBytes)
				if len(trimmed) > 200 {
					trimmed = trimmed[:200] + "..."
				}
				slog.Error("auth endpoint returned non-200", "status", resp.StatusCode, "body", trimmed)
				return "", fmt.Errorf("auth failed: %s", trimmed)
			}

			token := strings.TrimSpace(string(bodyBytes))
			slog.Debug("received auth token for message bus", "len", len(token))
			if token == "" {
				// treat empty token as transient error so we can retry
				if a == maxAuthAttempts {
					return "", errors.New("empty token from auth endpoint")
				}
				backoff := authBaseDelay * time.Duration(1<<(a-1))
				jitter := time.Duration(rng.Int63n(int64(backoff / 2)))
				sleep := backoff + jitter
				select {
				case <-time.After(sleep):
				case <-ctx.Done():
					return "", ctx.Err()
				}
				continue
			}
			return token, nil
		}

		return "", fmt.Errorf("auth failed after retries")
	}
}

func (f *DefaultMessageBusFactory) signNonce(nonce []byte) ([]byte, error) {
	sig, err := f.user.Sign(nonce)
	if err != nil {
		return nil, fmt.Errorf("sign nonce: %w", err)
	}
	return sig, nil
}
