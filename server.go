package workgroup

import "errors"

// Server creates function for canceling execution.
// First function passed in should block.
// Second function passed in should unblock first function.
func Server(serve func() error, shutdown func() error) RunFunc {
	return func(stop <-chan struct{}) error {
		done := make(chan error)
		defer close(done)

		go func() {
			done <- serve()
		}()

		select {
		case err := <-done:
			return err
		case <-stop:
			err := shutdown()
			err = errors.Join(err, <-done)
			return err
		}
	}
}
