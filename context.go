package workgroup

import "context"

func Context(ctx context.Context) RunFunc {
	return func(stop <-chan struct{}) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-stop:
			return nil
		}
	}
}
