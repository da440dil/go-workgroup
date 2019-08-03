package workgroup

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var g Group

	wait := make(chan struct{})
	err1 := errors.New("first")
	err2 := errors.New("second")

	g.Add(func(<-chan struct{}) error {
		<-wait
		return err1
	})
	g.Add(func(stop <-chan struct{}) error {
		<-stop
		return err2
	})

	result := make(chan error)
	go func() {
		result <- g.Run()
	}()
	close(wait)
	assert.Equal(t, err1, <-result)
}

func TestWithContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	g := NewGroup(WithContext(ctx))

	err1 := errors.New("first")
	err2 := errors.New("second")

	g.Add(func(stop <-chan struct{}) error {
		<-stop
		return err1
	})
	g.Add(func(stop <-chan struct{}) error {
		<-stop
		return err2
	})

	result := make(chan error)
	go func() {
		result <- g.Run()
	}()
	cancel()
	assert.Equal(t, context.Canceled, <-result)
}

func TestWithContextStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	g := NewGroup(WithContext(ctx))

	wait := make(chan struct{})
	err1 := errors.New("first")
	err2 := errors.New("second")

	g.Add(func(stop <-chan struct{}) error {
		<-stop
		return err1
	})
	g.Add(func(stop <-chan struct{}) error {
		<-wait
		return err2
	})

	result := make(chan error)
	go func() {
		result <- g.Run()
	}()
	close(wait)
	assert.Equal(t, err2, <-result)
	cancel()
}

func TestWithSignal(t *testing.T) {
	g := NewGroup(WithSignal(os.Interrupt))

	wait := make(chan struct{})
	err1 := errors.New("first")
	err2 := errors.New("second")

	g.Add(func(<-chan struct{}) error {
		<-wait
		return err1
	})
	g.Add(func(stop <-chan struct{}) error {
		<-stop
		return err2
	})

	result := make(chan error)
	go func() {
		result <- g.Run()
	}()
	close(wait)
	assert.Equal(t, err1, <-result)
}

func TestWithSignalEmpty(t *testing.T) {
	g := NewGroup(WithSignal())
	assert.Equal(t, len(g.fns), 0)
}

func TestZeroValue(t *testing.T) {
	var g Group

	err := g.Run()
	assert.NoError(t, err)
}
