# go-workgroup

[![Build Status](https://travis-ci.com/da440dil/go-workgroup.svg?branch=master)](https://travis-ci.com/da440dil/go-workgroup)
[![Coverage Status](https://coveralls.io/repos/github/da440dil/go-workgroup/badge.svg?branch=master)](https://coveralls.io/github/da440dil/go-workgroup?branch=master)
[![GoDoc](https://godoc.org/github.com/da440dil/go-workgroup?status.svg)](https://godoc.org/github.com/da440dil/go-workgroup)
[![Go Report Card](https://goreportcard.com/badge/github.com/da440dil/go-workgroup)](https://goreportcard.com/report/github.com/da440dil/go-workgroup)

Synchronization for groups of related goroutines.

## Basic usage

```go
// Create workgroup
var wg workgroup.Group
// Create http server
srv := http.Server{Addr: "127.0.0.1:8080"}
wg.Add(func(stop <-chan struct{}) error {
	go func() {
		<-stop
		// Stop http server
		srv.Close()
	}()
	// Start http server
	return srv.ListenAndServe()
})
```

## Example usage

- [example](./examples/workgroup-http-server/main.go) [http](https://golang.org/pkg/net/http/) server
- [example](./examples/workgroup-grpc-server/main.go) [grpc](https://github.com/grpc/grpc-go) server

Inspired by [workgroup](https://github.com/heptio/workgroup) package.