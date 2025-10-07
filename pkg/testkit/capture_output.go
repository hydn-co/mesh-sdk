package testkit

import (
	"bytes"
	"io"
	"os"
)

// CaptureOutput captures stdout while running the given function.
func CaptureOutput(fn func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var buf bytes.Buffer
	done := make(chan struct{})

	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	err := fn()
	_ = w.Close()
	os.Stdout = orig
	<-done

	return buf.String(), err
}
