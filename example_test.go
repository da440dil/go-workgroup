package workgroup_test

import (
	"errors"
	"fmt"

	"github.com/da440dil/go-workgroup"
)

func ExampleGroup() {
	var g workgroup.Group

	done := make(chan struct{})

	g.Add(func(<-chan struct{}) error {
		<-done
		fmt.Println("one")
		return errors.New("three")
	})
	g.Add(func(stop <-chan struct{}) error {
		<-stop
		fmt.Println("two")
		return errors.New("four")
	})

	result := make(chan error)
	go func() {
		result <- g.Run()
	}()
	close(done)
	fmt.Printf("%v\n", <-result)
	// Output:
	// one
	// two
	// three
}
