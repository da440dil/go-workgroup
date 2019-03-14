package workgroup_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/da440dil/go-workgroup"
)

func Example() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("cancel context")
		cancel()
	}()

	wg := workgroup.WithContext(ctx)
	wg.Add(serveHTTP("127.0.0.1:8081", newServeMux("I'm ready")))
	wg.Add(serveHTTP("127.0.0.1:8082", newServeMux("I'm live")))
	wg.Add(listenSignal(os.Interrupt))
	if err := wg.Wait(); err != nil {
		fmt.Println("error:", err)
	}

	// server listen 127.0.0.1:8081
	// server listen 127.0.0.1:8082
	// listen os signal
	// cancel context
	// stop listening os signal
	// stop listening 127.0.0.1:8082
	// stop listening 127.0.0.1:8081
	// error: context canceled
}

func serveHTTP(addr string, handler http.Handler) workgroup.Func {
	return func(stop <-chan struct{}) error {
		srv := http.Server{
			Addr:    addr,
			Handler: handler,
		}

		go func() {
			<-stop
			srv.Shutdown(context.Background())
		}()

		defer fmt.Println("stop listening", addr)

		fmt.Println("server listen", addr)
		return srv.ListenAndServe()
	}
}

func newServeMux(v string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, v)
	})
	return mux
}

func listenSignal(v ...os.Signal) workgroup.Func {
	return func(stop <-chan struct{}) error {
		sig := make(chan os.Signal)
		signal.Notify(sig, v...)

		go func() {
			<-stop
			close(sig)
		}()

		fmt.Println("listen os signal")
		<-sig
		signal.Stop(sig)
		fmt.Println("stop listening os signal")

		return nil
	}
}
