// Package context provides functions for cancelling execution using context.
package context

import "context"

// New creates function for cancelling execution.
func New(ctx context.Context) func(<-chan struct{}) error {
	return func(stop <-chan struct{}) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-stop:
			return nil
		}
	}
}
