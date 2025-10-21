# Lifecycle conventions

This project follows a simple Start/Stop convention for long-running services:

- `Start(ctx)` must be non-blocking. It should perform only quick initialization and launch background work (goroutines). If startup can fail, return an error. Prefer short, context-aware initialization.

- `Stop(ctx)` should cancel or terminate the work started by `Start`. It may block while waiting for shutdown and should respect the provided context for timeouts.

Note: implementers should make Start non-blocking and implement Stop to cancel/wait for shutdown. There's no built-in helper in this repo; prefer explicit Start/Stop implementations so the behavior is obvious.

Migration checklist

1. Find services whose `Start` blocks (long sleeps, loops, network calls, or `Consume` loops).
2. Move the blocking loop into a `run(ctx)`-style method that returns an error when the loop exits.
3. Make `Start(ctx)` launch `go run(ctx)` (as a goroutine) and return quickly.
4. Implement `Stop(ctx)` to cancel the context used by the run loop and wait for shutdown to complete.

Example

```go
// before: blocking Start
func (s *svc) Start(ctx context.Context) error {
    // long-running loop
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // work
        }
    }
}

// after: non-blocking Start that spawns run
func (s *svc) Start(ctx context.Context) error {
    s.runCtx, s.cancel = context.WithCancel(ctx)
    go s.run(s.runCtx)
    return nil
}

func (s *svc) run(ctx context.Context) error {
    // long-running loop
}

func (s *svc) Stop(ctx context.Context) error {
    s.cancel()
    // wait for run to finish or ctx timeout
    return nil
}
```

If you'd like, I can automatically scan and propose changes for services that look blocking. Start with a repo-wide scan and a small diff for one service? Let me know which service to refactor next.
