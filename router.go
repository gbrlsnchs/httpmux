package httpmux

import (
	"context"
	"net/http"

	"github.com/gbrlsnchs/radix"
)

// Router is a router that implements a http.Handler.
type Router struct {
	prefix  string
	methods map[string]*radix.Tree
}

// NewRouter creates and initializes a new Router.
func NewRouter() *Router {
	return &Router{methods: make(map[string]*radix.Tree)}
}

// Handle registers an http.Handler for a given method and path.
func (rt *Router) Handle(m, p string, h http.Handler) {
	rt.add(m, rt.prefix+resolvedPath(p), []interface{}{h})
}

// HandleFunc registers an http.HandlerFunc for a given method and path.
func (rt *Router) HandleFunc(m, p string, hfunc http.HandlerFunc) {
	rt.add(m, rt.prefix+resolvedPath(p), []interface{}{hfunc})
}

// HandleMiddlewares registers a stack of middlewares for a given method and path.
func (rt *Router) HandleMiddlewares(m, p string, mids ...interface{}) {
	midsToAdd := make([]interface{}, 0)

	for _, mid := range mids {
		if h, ok := mid.(http.Handler); ok {
			midsToAdd = append(midsToAdd, h)

			continue
		}

		if hfunc, ok := mid.(func(http.ResponseWriter, *http.Request)); ok {
			midsToAdd = append(midsToAdd, http.HandlerFunc(hfunc))
		}
	}

	rt.add(m, rt.prefix+resolvedPath(p), midsToAdd)
}

// ServeHTTP uses a radix tree to find a given URL path from a request.
//
// It is able to pass route parameters using the request's context.
// It can also break out of a middleware stack using the request's context cancelation.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rt.methods == nil || rt.methods[r.Method] == nil {
		http.NotFound(w, r)

		return
	}

	n, p := rt.methods[r.Method].GetByRune(r.URL.Path, ':', '/')

	if n != nil {
		ctx := r.Context()
		mids := n.Value.([]interface{})

		if len(p) > 0 {
			r = r.WithContext(context.WithValue(ctx, Params, p))
		}

		for _, m := range mids {
			switch v := m.(type) {
			case http.Handler:
				v.ServeHTTP(w, r)

			case http.HandlerFunc:
				v(w, r)
			}

			select {
			case <-ctx.Done():
				return

			default:
			}
		}

		return
	}

	http.NotFound(w, r)
}

// Use registers a Subrouter to be used by the Router.
func (rt *Router) Use(sub *Subrouter) {
	for endp, mids := range sub.endps {
		k := rt.prefix + endp

		for m, v := range mids {
			rt.add(m, k, v)
		}
	}
}

// WithPrefix sets a prefix for the Router, what makes
// all registered handlers use the prefix set.
func (rt *Router) WithPrefix(p string) *Router {
	if p == "/" {
		return rt
	}

	rt.prefix = p

	return rt
}

// add adds middlewares to the radix tree.
func (rt *Router) add(m, p string, mids []interface{}) {
	if p == "" {
		p = "/"
	}

	if len(mids) == 0 {
		return
	}

	if rt.methods[m] == nil {
		rt.methods[m] = radix.New(m)
	}

	rt.methods[m].Add(p, mids)
}
