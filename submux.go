package httpmux

import (
	"bytes"
	"net/http"
)

// Submux is a child of a Mux.
type Submux struct {
	path         string
	handlers     map[string]http.Handler
	handlerFuncs map[string]http.HandlerFunc
	submuxes     []*Submux
	parent       *Submux
}

// NewSubmux creates a new submultiplexer,
// which can be attached to a Mux.
func NewSubmux(p string) *Submux {
	return &Submux{path: resolvePath(p)}
}

// Add adds a new submultiplexer as a child of the submux.
func (smux *Submux) Add(c *Submux) {
	if smux.submuxes == nil {
		smux.submuxes = make([]*Submux, 0)
	}

	c.parent = smux
	smux.submuxes = append(smux.submuxes, c)
}

// Handle adds an http.Handler to the submultiplexer's handlers map.
func (smux *Submux) Handle(m string, h http.Handler) {
	if smux.handlers == nil {
		smux.handlers = make(map[string]http.Handler)
	}

	if smux.handlerFuncs != nil {
		smux.handlerFuncs[m] = nil
	}

	smux.handlers[m] = h
}

// Handle adds an http.HandlerFunc to the submultiplexer's handlerFuncs map.
func (smux *Submux) HandleFunc(m string, hfunc http.HandlerFunc) {
	if smux.handlerFuncs == nil {
		smux.handlerFuncs = make(map[string]http.HandlerFunc)
	}

	if smux.handlers != nil {
		smux.handlers[m] = nil
	}

	smux.handlerFuncs[m] = hfunc
}

// Path resolves the whole path a submultiplexer
// has considering its children.
func (smux *Submux) Path() string {
	var buf bytes.Buffer

	for smux != nil {
		p := buf.String()

		buf.Reset()
		buf.WriteString(smux.path)
		buf.WriteString(p)

		smux = smux.parent
	}

	return buf.String()
}
