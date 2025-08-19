package workgroup_test

import (
	"errors"
	"fmt"

	"github.com/da440dil/go-workgroup"
)

func ExampleGroup() {
	var g workgroup.Group

	wait := make(chan struct{})

	g.Add(func(<-chan struct{}) error {
		<-wait
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

	close(wait)
	fmt.Printf("%v\n", <-result)
	// Output:
	// one
	// two
	// three
	// four
}
