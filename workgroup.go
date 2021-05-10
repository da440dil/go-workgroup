// Package workgroup provides synchronization for groups of related goroutines.
package workgroup

// RunFunc is a function to execute with other related functions in its own goroutine.
// The closure of the channel passed to RunFunc should trigger return.
type RunFunc func(<-chan struct{}) error

// Group is a group of related goroutines.
// The zero value for a Group is fully usable without initialization.
type Group struct {
	fns []RunFunc
}

// Add adds a function to the Group.
// The function will be exectuted in its own goroutine when Run is called.
// Add must be called before Run.
func (g *Group) Add(fn RunFunc) {
	g.fns = append(g.fns, fn)
}

// Run executes each function registered via Add in its own goroutine.
// Run blocks until all functions have returned, then returns the first non-nil error (if any) from them.
// The first function to return will trigger the closure of the channel passed to each function, which should in turn, return.
func (g *Group) Run() error {
	if len(g.fns) == 0 {
		return nil
	}

	stop := make(chan struct{})
	done := make(chan error, len(g.fns))
	defer close(done)

	for _, fn := range g.fns {
		go func(fn RunFunc) {
			done <- fn(stop)
		}(fn)
	}

	var err error
	for i := 0; i < cap(done); i++ {
		if err == nil {
			err = <-done
		} else {
			<-done
		}
		if i == 0 {
			close(stop)
		}
	}
	return err
}
