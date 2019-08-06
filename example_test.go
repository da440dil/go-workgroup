package workgroup_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	g "github.com/da440dil/go-workgroup"
	gc "github.com/da440dil/go-workgroup/context"
	gq "github.com/da440dil/go-workgroup/quit"
	gs "github.com/da440dil/go-workgroup/signal"
)

func ExampleRun_hTTPServer() {
	// Create workgroup
	var wg g.Group
	// Add function to cancel execution using os signal
	wg.Add(gs.New(os.Interrupt))
	// Create http server
	srv := http.Server{Addr: "127.0.0.1:8080"}
	// Add function to start and stop http server
	wg.Add(gq.New(
		func() error {
			fmt.Printf("Server is about to listen at %v\n", srv.Addr)
			return srv.ListenAndServe()
		},
		func() {
			fmt.Println("Server is about to shutdown")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()

			err := srv.Shutdown(ctx)
			fmt.Printf("Server shutdown with error: %v\n", err)
		},
	))
	// Create context to cancel execution after 1 second
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second)
		fmt.Println("Context cancel")
		cancel()
	}()
	// Add function to cancel execution using context
	wg.Add(gc.New(ctx))
	// Execute each function
	err := wg.Run()
	fmt.Printf("Workgroup quit with error: %v\n", err)
	// Output:
	// Server is about to listen at 127.0.0.1:8080
	// Context cancel
	// Server is about to shutdown
	// Server shutdown with error: <nil>
	// Workgroup quit with error: context canceled
}
