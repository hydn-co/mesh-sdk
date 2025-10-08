package factories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/fgrzl/streamkit"
	"github.com/fgrzl/streamkit/pkg/transport/wskit"
	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/env"
	"github.com/hydn-co/mesh-sdk/pkg/localstore"
)

// StreamkitClientOptions holds configuration for the streamkit client factory.
//
// Fields left zero-valued will be replaced with sensible defaults by the factory
// (for example, AuthURL, StreamURL).
type StreamkitClientOptions struct {
	TenantID          uuid.UUID
	ClientCredentials localstore.ClientCredentials
	AuthURL           string
	StreamURL         string
	HTTPClient        *http.Client // Optional HTTP client for auth requests
}

// NewStreamkitClientFactory builds a streamkit.ClientFactory configured to
// authenticate against the mesh portal using the provided options struct.
// The returned factory uses environment variables to discover the portal and
// stream endpoints if not set in options, and will perform an HTTP call to obtain a JWT when
// establishing connections.
func NewStreamkitClientFactory(opts StreamkitClientOptions) streamkit.ClientFactory {
	authURL := opts.AuthURL
	if authURL == "" {
		authBase := env.GetEnvOrDefaultStr(env.MeshPortalBaseURL, "http://localhost:8080")
		authURL = strings.TrimRight(authBase, "/") + "/auth/stream/user"
	}
	streamURL := opts.StreamURL
	if streamURL == "" {
		streamURL = env.GetEnvOrDefaultStr(env.MeshStreamBaseURL, "ws://localhost:9444")
	}
	return &DefaultStreamkitClientFactory{
		tenantID:     opts.TenantID,
		clientID:     opts.ClientCredentials.ClientID,
		clientSecret: opts.ClientCredentials.ClientSecret,
		authURL:      authURL,
		streamURL:    streamURL,
		httpClient:   opts.HTTPClient,
	}
}

// DefaultStreamkitClientFactory is the default implementation of
// streamkit.ClientFactory used by the SDK. It obtains a JWT from the portal
// auth endpoint and uses it when creating websocket stream connections.
type DefaultStreamkitClientFactory struct {
	tenantID     uuid.UUID
	clientID     uuid.UUID
	clientSecret string
	authURL      string
	streamURL    string
	httpClient   *http.Client
}

// Get implements streamkit.ClientFactory.
func (f *DefaultStreamkitClientFactory) Get(ctx context.Context) (streamkit.Client, error) {

	provider := wskit.NewBidiStreamProvider(f.streamURL, f.fetchJWT)
	client := streamkit.NewClient(provider)
	return client, nil
}

func (f *DefaultStreamkitClientFactory) fetchJWT() (string, error) {
	type authenticateUser struct {
		TenantID     uuid.UUID `json:"tenant_id"`
		ClientID     uuid.UUID `json:"client_id"`
		ClientSecret string    `json:"client_secret"`
		Scopes       []string  `json:"scopes"`
	}

	payload := &authenticateUser{
		TenantID:     f.tenantID,
		ClientID:     f.clientID,
		ClientSecret: f.clientSecret,
		Scopes:       []string{"streamkit"},
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return "", err
	}

	client := f.httpClient
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequest(http.MethodPost, f.authURL, buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("auth failed: " + string(body))
	}

	jwtBytes, _ := io.ReadAll(resp.Body)
	return string(jwtBytes), nil
}
