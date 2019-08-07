package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	g "github.com/da440dil/go-workgroup"
	gc "github.com/da440dil/go-workgroup/template/context"
	gsh "github.com/da440dil/go-workgroup/template/shutdown"
	gsi "github.com/da440dil/go-workgroup/template/signal"
	"google.golang.org/grpc"
)

func main() {
	// Create workgroup
	var wg g.Group
	// Add function to cancel execution using os signal
	wg.Add(gsi.New(os.Interrupt))
	// Create grpc server
	srv := grpc.NewServer()
	// Add function to start and stop grpc server
	wg.Add(gsh.New(
		func() error {
			lis, err := net.Listen("tcp", "127.0.0.1:50051")
			if err != nil {
				return err
			}
			fmt.Printf("Server is about to listen at %v\n", lis.Addr())
			return srv.Serve(lis)
		},
		func() {
			fmt.Println("Server is about to shutdown")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()

			var wg g.Group
			wg.Add(gc.New(ctx))
			wg.Add(func(stop <-chan struct{}) error {
				srv.GracefulStop()
				return nil
			})
			err := wg.Run()
			fmt.Printf("Server shutdown with error: %v\n", err)
		},
	))
	// Create context to cancel execution after 5 seconds
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("Context cancel")
		cancel()
	}()
	// Add function to cancel execution using context
	wg.Add(gc.New(ctx))
	// Execute each function
	err := wg.Run()
	fmt.Printf("Workgroup quit with error: %v\n", err)
	// Output:
	// Server is about to listen at 127.0.0.1:50051
	// Context cancel
	// Server is about to shutdown
	// Server shutdown with error: <nil>
	// Workgroup quit with error: context canceled
}
