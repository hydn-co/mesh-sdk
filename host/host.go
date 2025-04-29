package host

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ShutdownFunc func(context.Context) error

func WaitForShutdown(fns ...ShutdownFunc) {
	WaitForShutdownWithTimeout(5*time.Second, fns...)
}

func WaitForShutdownWithTimeout(timeout time.Duration, fns ...ShutdownFunc) {
	// Listen for shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // wait for signal

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, fn := range fns {
		wg.Add(1)
		go func(fn ShutdownFunc) {
			defer wg.Done()
			_ = fn(ctx)
		}(fn)
	}

	wg.Wait()
}
