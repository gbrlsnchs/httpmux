# httpmux (HTTP request multiplexer)
[![Build Status](https://travis-ci.org/gbrlsnchs/httpmux.svg?branch=master)](https://travis-ci.org/gbrlsnchs/httpmux)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/gbrlsnchs/httpmux)

## About
This package is an HTTP request multiplexer that enables router nesting 
and middleware stacking for [Go] (or Golang) HTTP servers.

It uses standard approaches, such as the `context` package for retrieving
route parameters and short-circuiting middlewares, what makes it easy to be used
on both new and old projects, since it doesn't present any new pattern.

## Usage
Full documentation [here].

## Example (from example_test.go)
```go
package httpmux_test

import (
	"log"
	"net/http"

	"github.com/gbrlsnchs/httpmux"
)

func Example() {
	rt := httpmux.NewRouter()

	rt.HandleMiddlewares(http.MethodGet, "/:path",
		// Logger.
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("r.URL.Path = %s\n", r.URL.Path)
		},
		// Guard.
		func(w http.ResponseWriter, r *http.Request) {
			params := httpmux.Params(r)

			if params["path"] == "forbidden" {
				w.WriteHeader(http.StatusForbidden)
				httpmux.Cancel(r)
			}
		},
		// Handler.
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	http.ListenAndServe("/", rt)
}
```

## Contribution
### How to help:
- Pull Requests
- Issues
- Opinions

[Go]: https://golang.org
[here]: https://godoc.org/github.com/gbrlsnchs/httpmux
