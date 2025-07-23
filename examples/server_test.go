package workgroup_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/da440dil/go-workgroup"
)

func ExampleServer() {
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
}
