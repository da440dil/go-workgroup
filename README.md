# go-workgroup

[![Build Status](https://travis-ci.com/da440dil/go-workgroup.svg?branch=master)](https://travis-ci.com/da440dil/go-workgroup)
[![Coverage Status](https://coveralls.io/repos/github/da440dil/go-workgroup/badge.svg?branch=master)](https://coveralls.io/github/da440dil/go-workgroup?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/da440dil/go-workgroup.svg)](https://pkg.go.dev/github.com/da440dil/go-workgroup)
[![GoDoc](https://godoc.org/github.com/da440dil/go-workgroup?status.svg)](https://godoc.org/github.com/da440dil/go-workgroup)
[![Go Report Card](https://goreportcard.com/badge/github.com/da440dil/go-workgroup)](https://goreportcard.com/report/github.com/da440dil/go-workgroup)

Synchronization for groups of related goroutines.

## Example HTTP server

```go
// Create workgroup
var wg workgroup.Group
// Add function to cancel execution using os signal
wg.Add(workgroup.Signal())
// Create http server
srv := http.Server{Addr: "127.0.0.1:8080"}
// Add function to start and stop http server
wg.Add(workgroup.Server(
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
))
// Create context to cancel execution after 5 seconds
ctx, cancel := context.WithCancel(context.Background())
go func() {
	time.Sleep(time.Second * 5)
	fmt.Println("Context canceled")
	cancel()
}()
// Add function to cancel execution using context
wg.Add(workgroup.Context(ctx))
// Execute each function
err := wg.Run()
fmt.Printf("Workgroup run stopped with error: %v\n", err)
// Output in case execution has been canceled with context:
// Server listen at 127.0.0.1:8080
// Context canceled
// Server is about to shutdown
// Server stopped listening with error: http: Server closed
// Server shutdown with error: <nil>
// Workgroup run stopped with error: context canceled
//
// Output in case execution has been interrupted with os signal:
// Server listen at 127.0.0.1:8080
// Server is about to shutdown
// Server stopped listening with error: http: Server closed
// Server shutdown with error: <nil>
// Workgroup run stopped with error: <nil>
```

## Example gRPC server

```go
// Create workgroup
var wg workgroup.Group
// Add function to cancel execution using os signal
wg.Add(workgroup.Signal())
// Create grpc server
srv := grpc.NewServer()
// Add function to start and stop grpc server
wg.Add(workgroup.Server(
	func() error {
		lis, err := net.Listen("tcp", "127.0.0.1:50051")
		if err != nil {
			return err
		}
		fmt.Printf("Server listen at %v\n", lis.Addr())
		err = srv.Serve(lis)
		fmt.Printf("Server stopped listening with error: %v\n", err)
		return err
	},
	func() error {
		fmt.Println("Server is about to shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		var wg workgroup.Group
		wg.Add(workgroup.Context(ctx))
		wg.Add(func(stop <-chan struct{}) error {
			srv.GracefulStop()
			return nil
		})
		err := wg.Run()
		fmt.Printf("Server shutdown with error: %v\n", err)
		return err
	},
))
// Create context to cancel execution after 5 seconds
ctx, cancel := context.WithCancel(context.Background())
go func() {
	time.Sleep(time.Second * 5)
	fmt.Println("Context canceled")
	cancel()
}()
// Add function to cancel execution using context
wg.Add(workgroup.Context(ctx))
// Execute each function
err := wg.Run()
fmt.Printf("Workgroup run stopped with error: %v\n", err)
// Output in case execution has been canceled with context:
// Server listen at 127.0.0.1:50051
// Context canceled
// Server is about to shutdown
// Server shutdown with error: <nil>
// Server stopped listening with error: <nil>
// Workgroup run stopped with error: context canceled
//
// Output in case execution has been interrupted with os signal:
// Server listen at 127.0.0.1:50051
// Server is about to shutdown
// Server shutdown with error: <nil>
// Server stopped listening with error: <nil>
// Workgroup run stopped with error: <nil>
```
