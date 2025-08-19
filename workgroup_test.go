package workgroup

import (
	"errors"
	"testing"
)

func TestStopped(t *testing.T) {
	wait := make(chan struct{})
	err := errors.New("err")

	f1 := func(<-chan struct{}) error {
		<-wait
		return err
	}
	f2 := func(stop <-chan struct{}) error {
		<-stop
		return nil
	}

	close(wait)
	assert(t, err, Run(f1, f2))
}

func TestError(t *testing.T) {
	wait := make(chan struct{})
	err := errors.New("err")

	f1 := func(<-chan struct{}) error {
		<-wait
		return nil
	}
	f2 := func(stop <-chan struct{}) error {
		<-stop
		return err
	}

	close(wait)
	assert(t, err, Run(f1, f2))
}

func TestGroupZeroValue(t *testing.T) {
	var g Group
	assert(t, nil, g.Run())
}

func TestGroupStopped(t *testing.T) {
	var g Group

	wait := make(chan struct{})
	err := errors.New("err")

	g.Add(func(<-chan struct{}) error {
		<-wait
		return err
	})
	g.Add(func(stop <-chan struct{}) error {
		<-stop
		return nil
	})

	close(wait)
	assert(t, err, g.Run())
}

func TestGroupError(t *testing.T) {
	var g Group

	wait := make(chan struct{})
	err := errors.New("err")

	g.Add(func(<-chan struct{}) error {
		<-wait
		return nil
	})
	g.Add(func(stop <-chan struct{}) error {
		<-stop
		return err
	})

	close(wait)
	assert(t, err, g.Run())
}

func assert(t *testing.T, want, got error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}
