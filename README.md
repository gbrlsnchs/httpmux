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

## Contribution
### How to help:
- Pull Requests
- Issues
- Opinions

[Go]: https://golang.org
[here]: https://godoc.org/github.com/gbrlsnchs/httpmux
