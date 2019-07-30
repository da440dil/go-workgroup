# go-workgroup

[![Build Status](https://travis-ci.com/da440dil/go-workgroup.svg?branch=master)](https://travis-ci.com/da440dil/go-workgroup)
[![GoDoc](https://godoc.org/github.com/da440dil/go-workgroup?status.svg)](https://godoc.org/github.com/da440dil/go-workgroup)

Synchronization for groups of related goroutines.

## Example

```go
// Create workgroup
var wg workgroup.Group

// Add function to start http server
wg.Add(func(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	addr := "127.0.0.1:8080"
	srv := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	done := make(chan error, 2)
	go func() {
		<-stop
		fmt.Println("Server is about to stop")
		done <- srv.Shutdown(context.Background())
	}()

	go func() {
		fmt.Printf("Server starts listening at %s\n", addr)
		done <- srv.ListenAndServe()
	}()

	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server stopped with error %v\n", err)
			return err
		}
	}
	fmt.Println("Server stopped")
	return nil
})

// Add function to start listening os signal
wg.Add(func(stop <-chan struct{}) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-stop
		close(sig)
	}()

	fmt.Println("Server starts listening os signal")
	<-sig
	fmt.Println("Server stops listening os signal")
	signal.Stop(sig)
	return nil
})

// Create context to cancel execution
ctx, cancel := context.WithCancel(context.Background())
go func() {
	time.Sleep(time.Second * 10)
	fmt.Println("Cancel context")
	cancel()
}()

// Create workgroup with context
wg = workgroup.WithContext(ctx, wg)

// Execute each function
if err := wg.Run(); err != nil {
	fmt.Println("Error:", err)
}

// Server starts listening at 127.0.0.1:8080
// Server starts listening os signal
// Cancel context
// Server stops listening os signal
// Server is about to stop
// Server stopped
// Error: context canceled
```

Inspired by [workgroup](https://github.com/heptio/workgroup) package.