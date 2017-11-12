package httpmux

import (
	"context"
	"net/http"

	"github.com/gbrlsnchs/patricia"
)

// Mux is an HTTP multiplexer that implements
// an http.Handler and serves as the parent of
// other multiplexers.
type Mux struct {
	path  string
	trees map[string]*patricia.Tree
}

// New creates a new Mux. Its path is then
// set in lowercase and, if needed, prepended a leading slash.
func New(p string) *Mux {
	return &Mux{path: resolvePath(p)}
}

// Add adds a submultiplexer's handler to
// the multiplexer's handlers tree.
func (mux *Mux) Add(smux *Submux) {
	for k, v := range smux.handlers {
		mux.addValue(k, smux.Path(), v)
	}

	for k, v := range smux.handlerFuncs {
		mux.addValue(k, smux.Path(), v)
	}

	for _, c := range smux.submuxes {
		mux.Add(c)
	}
}

// Debug prints the multiplexer's handlers tree in debug mode.
func (mux *Mux) Debug() error {
	for _, t := range mux.trees {
		err := t.Debug()

		if err != nil {
			return err
		}
	}

	return nil
}

// Handle adds an http.Handler to the multiplexer's handlers tree.
func (mux *Mux) Handle(m string, h http.Handler) {
	mux.addValue(m, "", h)
}

// HandleFunc adds an http.HandlerFunc to the multiplexer's handlers tree.
func (mux *Mux) HandleFunc(m string, hfunc http.HandlerFunc) {
	mux.addValue(m, "", hfunc)
}

// Print prints the multiplexer's handlers tree.
func (mux *Mux) Print() error {
	for _, t := range mux.trees {
		err := t.Print()

		if err != nil {
			return err
		}
	}

	return nil
}

// ServeHTTP looks for a match inside the tree, according to the request method,
// and if a node is found, the node's value is run if is either an http.Handler
// an http.HandlerFunc.
func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mux.trees == nil || mux.trees[r.Method] == nil {
		http.NotFound(w, r)

		return
	}

	n, p := mux.trees[r.Method].
		GetByRune(r.URL.Path, ':', '/')

	if n != nil && n.Value != nil {
		if len(p) > 0 {
			r = r.WithContext(context.WithValue(r.Context(), Params, p))
		}

		switch v := n.Value.(type) {
		case http.Handler:
			v.ServeHTTP(w, r)

			return

		case http.HandlerFunc:
			v(w, r)

			return
		}
	}

	http.NotFound(w, r)
}

// addValue resolves a path and merges it with the multiplexer's path
// and then adds a value to the handlers tree.
func (mux *Mux) addValue(m string, p string, v interface{}) *Mux {
	if v == nil {
		return mux
	}

	if mux.trees == nil {
		mux.trees = make(map[string]*patricia.Tree)
	}

	if mux.trees[m] == nil {
		mux.trees[m] = patricia.New(m)
	}

	mux.trees[m].Add(mux.path+p, v)

	return mux
}
