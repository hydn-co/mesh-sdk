package lifecycle

import (
	"context"
	"sync"
	"time"
)

// FailureAction controls what the supervisor does when the run function returns an error.
type FailureAction int

const (
	FailHost FailureAction = iota
	Restart
	Ignore
)

// Policy controls restart/backoff behavior.
type Policy struct {
	Action      FailureAction
	MaxRestarts int // -1 for unlimited
	Backoff     func(attempt int) time.Duration
}

// Supervisor runs a blocking runFunc(ctx) and optionally restarts it according to Policy.
// Supervisor implements Service so it can be passed to Host.
type Supervisor struct {
	runFunc func(ctx context.Context) error
	policy  Policy

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
	// last error from runFunc
	lastErr error
}

// NewSupervisor creates a supervisor for the provided run function and policy.
// runFunc should block until finished or ctx is canceled. It returns an error on failure.
func NewSupervisor(runFunc func(ctx context.Context) error, policy Policy) *Supervisor {
	return &Supervisor{runFunc: runFunc, policy: policy}
}

// Start launches the supervisor loop which runs runFunc(ctx) and applies the policy.
func (s *Supervisor) Start(parentCtx context.Context) error {
	s.mu.Lock()
	if s.done != nil {
		s.mu.Unlock()
		return nil
	}
	ctx, cancel := context.WithCancel(parentCtx)
	s.cancel = cancel
	s.done = make(chan struct{})
	s.mu.Unlock()

	go func() {
		defer func() {
			s.mu.Lock()
			if s.done != nil {
				close(s.done)
			}
			s.done = nil
			s.cancel = nil
			s.mu.Unlock()
		}()

		attempts := 0
		for {
			// respect parent cancelation
			if ctx.Err() != nil {
				return
			}

			err := s.runFunc(ctx)
			if err == nil {
				// clean exit
				return
			}

			s.mu.Lock()
			s.lastErr = err
			s.mu.Unlock()

			switch s.policy.Action {
			case Ignore:
				return
			case FailHost:
				// notify by setting lastErr and returning; Host may observe this or runner can use OnFatal
				return
			case Restart:
				attempts++
				if s.policy.MaxRestarts >= 0 && attempts > s.policy.MaxRestarts {
					return
				}
				wait := time.Second
				if s.policy.Backoff != nil {
					wait = s.policy.Backoff(attempts)
				} else {
					// default exponential backoff capped
					wait = time.Duration(1<<uint(attempts-1)) * 500 * time.Millisecond
					if wait > 30*time.Second {
						wait = 30 * time.Second
					}
				}
				select {
				case <-time.After(wait):
					continue
				case <-ctx.Done():
					return
				}
			default:
				return
			}
		}
	}()
	return nil
}

// Stop cancels the running supervisor and waits for the worker to exit.
func (s *Supervisor) Stop(ctx context.Context) error {
	s.mu.Lock()
	cancel := s.cancel
	done := s.done
	s.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if done == nil {
		return nil
	}

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// LastError returns the last error reported by the run function (if any).
func (s *Supervisor) LastError() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastErr
}

// Helper to create a Policy with restart semantics.
func RestartPolicy(maxRestarts int, backoff func(attempt int) time.Duration) Policy {
	return Policy{Action: Restart, MaxRestarts: maxRestarts, Backoff: backoff}
}

// Helper to create a FailHost policy.
func FailPolicy() Policy { return Policy{Action: FailHost} }

// Helper to create an Ignore policy.
func IgnorePolicy() Policy { return Policy{Action: Ignore} }
