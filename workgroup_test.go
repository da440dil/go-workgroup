package workgroup

import (
	"errors"
	"testing"
)

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
	if want != got {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
