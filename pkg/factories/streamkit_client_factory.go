package factories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/fgrzl/streamkit"
	"github.com/fgrzl/streamkit/pkg/transport/wskit"
	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/env"
	"github.com/hydn-co/mesh-sdk/pkg/localstore"
)

func NewStreamkitClientFactory(tenantID uuid.UUID, creds localstore.ClientCredentials) streamkit.ClientFactory {
	return &DefaultStreamkitClientFactory{
		tenantID:     tenantID,
		clientID:     creds.ClientID,
		clientSecret: creds.ClientSecret,
		authURL:      env.GetEnvOrDefaultStr(env.MeshPortalURL, "http://localhost:8080") + "/auth/stream/user",
		streamURL:    env.GetEnvOrDefaultStr(env.MeshStreamURL, "ws://localhost:9444"),
	}
}

type DefaultStreamkitClientFactory struct {
	tenantID     uuid.UUID
	clientID     uuid.UUID
	clientSecret string
	authURL      string
	streamURL    string
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

	resp, err := http.Post(f.authURL, "application/json", buf)
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
