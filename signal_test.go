package workgroup

import (
	"errors"
	"testing"
)

func TestSignalStopped(t *testing.T) {
	var g Group

	g.Add(Signal())

	wait := make(chan struct{})
	err := errors.New("err")
	g.Add(func(<-chan struct{}) error {
		<-wait
		return err
	})

	close(wait)
	assert(t, err, g.Run())
}
