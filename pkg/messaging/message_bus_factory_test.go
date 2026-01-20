package messaging

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nkeys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testRoundTripper simulates a sequence of responses: first N errors, then a success
type testRoundTripper struct {
	failsBefore int
	called      int
	respBody    string
	respStatus  int
}

func (t *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t.called++
	if t.called <= t.failsBefore {
		return nil, errors.New("simulated transport error")
	}
	// success response
	r := &http.Response{
		StatusCode: t.respStatus,
		Body:       io.NopCloser(strings.NewReader(t.respBody)),
		Header:     make(http.Header),
	}
	return r, nil
}

func buildFactoryForTest(t *testing.T) *DefaultMessageBusFactory {
	// Arrange: create a real nkeys user keypair so it implements the full interface
	user, err := nkeys.CreateUser()
	require.NoError(t, err)

	// Act: build factory
	return &DefaultMessageBusFactory{
		tenantID:     uuid.Nil,
		clientID:     uuid.Nil,
		clientSecret: "s",
		user:         user,
		authURL:      "http://example.local/auth/broker/user",
		brokerURL:    "ws://localhost:9222",
	}
}

func TestFetchJWTSucceedsAfterRetries(t *testing.T) {
	f := buildFactoryForTest(t)

	// Arrange: replace default client's transport to simulate two failures then success
	tr := &testRoundTripper{failsBefore: 2, respBody: "MYTOKEN", respStatus: 200}
	oldTransport := http.DefaultTransport
	http.DefaultTransport = tr
	defer func() { http.DefaultTransport = oldTransport }()

	// Act: obtain fetch function and call it
	fn := f.fetchJWT(context.Background())

	// give generous timeout (not strictly necessary for this unit test but kept for parity)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tok, err := fn()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "MYTOKEN", tok)
	assert.Equal(t, 3, tr.called)
	_ = ctx
}

func TestFetchJWTFailsAfterRetries(t *testing.T) {
	f := buildFactoryForTest(t)

	// Arrange: transport that always fails
	tr := &testRoundTripper{failsBefore: 10}
	oldTransport := http.DefaultTransport
	http.DefaultTransport = tr
	defer func() { http.DefaultTransport = oldTransport }()

	// Act
	fn := f.fetchJWT(context.Background())
	_, err := fn()

	// Assert
	require.Error(t, err)
}

func TestFetchJWTNon200ReturnsError(t *testing.T) {
	// Arrange
	f := buildFactoryForTest(t)
	tr := &testRoundTripper{failsBefore: 0, respBody: "server error", respStatus: 500}
	oldTransport := http.DefaultTransport
	http.DefaultTransport = tr
	defer func() { http.DefaultTransport = oldTransport }()

	// Act
	fn := f.fetchJWT(context.Background())
	tok, err := fn()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", tok)
	assert.Equal(t, 1, tr.called)
}

func TestFetchJWTEmptyTokenRetriesThenError(t *testing.T) {
	// Arrange
	f := buildFactoryForTest(t)
	// ensure auth attempts default to 3 for test determinism
	f.authAttempts = 3
	tr := &testRoundTripper{failsBefore: 0, respBody: "", respStatus: 200}
	oldTransport := http.DefaultTransport
	http.DefaultTransport = tr
	defer func() { http.DefaultTransport = oldTransport }()

	// Act
	fn := f.fetchJWT(context.Background())
	tok, err := fn()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", tok)
	assert.Equal(t, 3, tr.called)
}
