// Package signal provides functions for cancelling execution using notifications on os signals.
package signal

import (
	"os"
	"os/signal"
)

// New creates function for cancelling execution.
func New(sig ...os.Signal) func(<-chan struct{}) error {
	return func(stop <-chan struct{}) error {
		if len(sig) == 0 {
			sig = append(sig, os.Interrupt)
		}
		done := make(chan os.Signal, len(sig))
		signal.Notify(done, sig...)
		select {
		case <-stop:
		case <-done:
		}
		signal.Stop(done)
		close(done)
		return nil
	}
}
