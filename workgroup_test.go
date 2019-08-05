package workgroup

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
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

func TestGroupZeroValue(t *testing.T) {
	var g Group

	err := g.Run()
	assert.NoError(t, err)
}
