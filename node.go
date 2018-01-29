package httpmux

import "net/http"

type node struct {
	len          int
	handlerFuncs []http.HandlerFunc
	handlers     []http.Handler
}

func newNode(len int) *node {
	return &node{
		len:          len,
		handlerFuncs: make([]http.HandlerFunc, len),
		handlers:     make([]http.Handler, len),
	}
}

func (n *node) add(pos int, h interface{}) {
	if h, ok := h.(http.HandlerFunc); ok {
		n.handlerFuncs[pos] = h

		return
	}

	if h, ok := h.(http.Handler); ok {
		n.handlers[pos] = h
	}
}
