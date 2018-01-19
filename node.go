package httpmux

import "net/http"

type node struct {
	count      uint
	handler    []http.Handler
	handleFunc []http.HandlerFunc
}
