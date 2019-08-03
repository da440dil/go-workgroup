# go-workgroup

[![Build Status](https://travis-ci.com/da440dil/go-workgroup.svg?branch=master)](https://travis-ci.com/da440dil/go-workgroup)
[![Coverage Status](https://coveralls.io/repos/github/da440dil/go-workgroup/badge.svg?branch=master)](https://coveralls.io/github/da440dil/go-workgroup?branch=master)
[![GoDoc](https://godoc.org/github.com/da440dil/go-workgroup?status.svg)](https://godoc.org/github.com/da440dil/go-workgroup)
[![Go Report Card](https://goreportcard.com/badge/github.com/da440dil/go-workgroup)](https://goreportcard.com/report/github.com/da440dil/go-workgroup)

Synchronization for groups of related goroutines.

## Example

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/da440dil/go-workgroup"
)

func main() {
	// Create context to cancel execution after 5 seconds
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("Cancel context")
		cancel()
	}()

	// Create workgroup
	wg := workgroup.NewGroup(
		// cancel execution using context
		workgroup.WithContext(ctx),
		// cancel execution using os signal
		workgroup.WithSignal(os.Interrupt),
	)

	// Add function to start http server
	wg.Add(func(stop <-chan struct{}) error {
		srv := http.Server{Addr: "127.0.0.1:8080"}

		done := make(chan error, 2)
		go func() {
			<-stop
			fmt.Println("Server is about to stop")
			done <- srv.Shutdown(context.Background())
		}()

		go func() {
			fmt.Println("Server starts listening")
			done <- srv.ListenAndServe()
		}()

		for i := 0; i < cap(done); i++ {
			if err := <-done; err != nil && err != http.ErrServerClosed {
				return err
			}
		}
		fmt.Println("Server stopped")
		return nil
	})

	// Execute each function
	if err := wg.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Server starts listening
	// Cancel context
	// Server is about to stop
	// Server stopped
	// Error: context canceled
}
```

Inspired by [workgroup](https://github.com/heptio/workgroup) package.