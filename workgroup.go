// Package workgroup provides synchronization for groups of related goroutines.
package workgroup

import (
	"context"
)

// Option is function returned by functions for setting options.
type Option func(g *Group)

// WithContext is helper function which adds a function to the Group
// for cancelling execution using context.
func WithContext(ctx context.Context) Option {
	return func(g *Group) {
		g.Add(func(stop <-chan struct{}) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-stop:
				return nil
			}
		})
	}
}

// Group is a group of related goroutines.
// The zero value for a Group is fully usable without initialization.
type Group struct {
	fns []Run
}

// NewGroup creates new Group.
func NewGroup(options ...Option) *Group {
	g := &Group{}
	for _, fn := range options {
		fn(g)
	}
	return g
}

// Run is a function to execute with other related functions in its own goroutine.
// The closure of the channel passed to Run should trigger return.
type Run func(<-chan struct{}) error

// Add adds a function to the Group.
// The function will be exectuted in its own goroutine when Run is called.
// Add must be called before Run.
func (g *Group) Add(fn Run) {
	g.fns = append(g.fns, fn)
}

// Run exectues each function registered via Add in its own goroutine.
// Run blocks until all functions have returned.
// The first function to return will trigger the closure of the channel passed to each function,
// which should in turn, return.
// The return value from the first function to exit will be returned to the caller of Run.
func (g *Group) Run() error {
	if len(g.fns) == 0 {
		return nil
	}

	stop := make(chan struct{})
	done := make(chan error, len(g.fns))
	for _, fn := range g.fns {
		go func(fn Run) {
			done <- fn(stop)
		}(fn)
	}

	var err error
	for i := 0; i < cap(done); i++ {
		if i == 0 {
			err = <-done
			close(stop)
		} else {
			<-done
		}
	}
	close(done)
	return err
}
