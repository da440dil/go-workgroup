// Package workgroup provides synchronization for groups of related goroutines.
package workgroup

import "context"

// WithContext allocates and returns new Group.
// Context allows cancelling execution.
func WithContext(ctx context.Context) *Group {
	return &Group{ctx: ctx}
}

// Group is a group of related goroutines.
// The zero value for a Group is fully usable without initalisation.
type Group struct {
	fns []Func
	ctx context.Context
}

// Add adds a function to the Group.
// The function will be exectuted in its own goroutine when Wait is called.
// Add must be called before Wait.
func (g *Group) Add(fn Func) {
	g.fns = append(g.fns, fn)
}

// Wait exectues each function registered via Add in its own goroutine.
// Wait blocks until all functions have returned.
// The first function to return will trigger the closure of the channel passed to each function,
// which should in turn, return.
// The return value from the first function to exit will be returned to the caller of Wait.
func (g *Group) Wait() error {
	if len(g.fns) < 1 {
		return nil
	}

	fns := g.fns
	if g.ctx != nil {
		fns = append(g.fns, func(stop <-chan struct{}) error {
			select {
			case <-g.ctx.Done():
				return g.ctx.Err()
			case <-stop:
				return nil
			}
		})
	}
	g.fns = nil

	stop := make(chan struct{})
	done := make(chan error, len(fns))
	for _, fn := range fns {
		go func(fn Func) {
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

// Func is a function to execute with other related functions in its own goroutine.
// The closure of the channel passed to Func should trigger return.
type Func func(<-chan struct{}) error
