# go-workgroup

[![CI](https://github.com/da440dil/go-workgroup/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/da440dil/go-workgroup/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/da440dil/go-workgroup/badge.svg?branch=master)](https://coveralls.io/github/da440dil/go-workgroup?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/da440dil/go-workgroup.svg)](https://pkg.go.dev/github.com/da440dil/go-workgroup)
[![Go Report Card](https://goreportcard.com/badge/github.com/da440dil/go-workgroup)](https://goreportcard.com/report/github.com/da440dil/go-workgroup)

Synchronization for groups of related goroutines.

## [Example](./examples/server_test.go) HTTP server

```go
srv := http.Server{Addr: "127.0.0.1:8080"}
// Create context to cancel execution after 100 milliseconds,
// increase timeout to manually interrupt execution with SIGINT
ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
defer cancel()

// Execute each function in its own goroutine
err := workgroup.Run(
	workgroup.Server(
		func() error {
			fmt.Printf("Server listen at %v\n", srv.Addr)
			err := srv.ListenAndServe()
			fmt.Printf("Server stopped listening with error: %v\n", err)
			if err != http.ErrServerClosed {
				return err
			}
			return nil
		},
		func() error {
			fmt.Println("Server is about to shutdown")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
			defer cancel()

			err := srv.Shutdown(ctx)
			fmt.Printf("Server shutdown with error: %v\n", err)
			return err
		},
	),
	workgroup.Context(ctx),
	workgroup.Signal(),
)
fmt.Printf("Workgroup run stopped with error: %v\n", err)
// Output:
// Server listen at 127.0.0.1:8080
// Server is about to shutdown
// Server stopped listening with error: http: Server closed
// Server shutdown with error: <nil>
// Workgroup run stopped with error: context deadline exceeded
```
