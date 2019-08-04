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
	"os/signal"
	"time"

	"github.com/da440dil/go-workgroup"
)

func main() {
	// Create context to cancel execution after 5 seconds
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("Context cancel")
		cancel()
	}()
	// Create workgroup
	wg := workgroup.NewGroup(workgroup.WithContext(ctx))
	workgroup.WithContext(ctx)(wg)
	// Add function to cancel execution using os signal
	wg.Add(func(stop <-chan struct{}) error {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt)
		select {
		case <-stop:
		case <-done:
		}
		signal.Stop(done)
		close(done)
		return nil
	})
	// Add function to start http server
	wg.Add(func(stop <-chan struct{}) error {
		// Create http server
		srv := http.Server{Addr: "127.0.0.1:8080"}

		done := make(chan error, 1)
		go func() {
			fmt.Printf("Server is about to listen at at %v\n", srv.Addr)
			done <- srv.ListenAndServe()
		}()

		select {
		case err := <-done:
			close(done)
			fmt.Printf("Server stops listening with error: %v\n", err)
			return err
		case <-stop:
			fmt.Println("Server is about to shutdown")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()

			err := srv.Shutdown(ctx)
			fmt.Printf("Server shutdown with error: %v\n", err)
			return err
		}
	})
	// Execute each function
	err := wg.Run()
	fmt.Printf("Workgroup quit with error: %v\n", err)

	// Server is about to listen at at 127.0.0.1:8080
	// Context cancel
	// Server is about to shutdown
	// Server shutdown with error: <nil>
	// Workgroup quit with error: context canceled
}
```

Inspired by [workgroup](https://github.com/heptio/workgroup) package.