package factories

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
	// create a real nkeys user keypair so it implements the full interface
	user, err := nkeys.CreateUser()
	if err != nil {
		t.Fatalf("failed to create test nkeys user: %v", err)
	}

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

	// replace default client's transport
	tr := &testRoundTripper{failsBefore: 2, respBody: "MYTOKEN", respStatus: 200}
	oldTransport := http.DefaultTransport
	http.DefaultTransport = tr
	defer func() { http.DefaultTransport = oldTransport }()

	fn := f.fetchJWT(context.Background())

	// give generous timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tok, err := fn()
	if err != nil {
		t.Fatalf("expected success after retries, got err: %v", err)
	}
	if tok != "MYTOKEN" {
		t.Fatalf("unexpected token: %s", tok)
	}
	if tr.called != 3 {
		t.Fatalf("expected 3 calls, got %d", tr.called)
	}
	_ = ctx
}

func TestFetchJWTFailsAfterRetries(t *testing.T) {
	f := buildFactoryForTest(t)

	tr := &testRoundTripper{failsBefore: 10}
	oldTransport := http.DefaultTransport
	http.DefaultTransport = tr
	defer func() { http.DefaultTransport = oldTransport }()

	fn := f.fetchJWT(context.Background())
	_, err := fn()
	if err == nil {
		t.Fatalf("expected error after retries, got nil")
	}
}
