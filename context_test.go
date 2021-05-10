package workgroup

import (
	"context"
	"errors"
	"testing"
)

func TestContextStopped(t *testing.T) {
	var g Group

	g.Add(Context(context.TODO()))

	wait := make(chan struct{})
	err := errors.New("err")
	g.Add(func(<-chan struct{}) error {
		<-wait
		return err
	})

	close(wait)
	assert(t, err, g.Run())
}

func TestContextCanceled(t *testing.T) {
	var g Group

	ctx, cancel := context.WithCancel(context.Background())
	g.Add(Context(ctx))

	g.Add(func(stop <-chan struct{}) error {
		<-stop
		return nil
	})

	cancel()
	assert(t, context.Canceled, g.Run())
}
