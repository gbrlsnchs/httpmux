package httpmux

import (
	"context"
	"net/http"

	"github.com/gbrlsnchs/radix"
)

// ParamsKey is the value for used for
// retrieving the parameters map in a Context.
var ParamsKey interface{}

// Router is a router that implements a http.Handler.
type Router struct {
	prefix  string
	methods map[string]*radix.Tree
	common  []interface{}
}

// NewRouter creates and initializes a new Router.
func NewRouter() *Router {
	return &Router{methods: make(map[string]*radix.Tree)}
}

// Debug prints the router's radix tree structure
// for each HTTP method registered.
func (rt *Router) Debug() error {
	for _, m := range rt.methods {
		if err := m.Debug(); err != nil {
			return err
		}
	}

	return nil
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
		tnode := n.Value.(*node)

		if len(p) > 0 {
			r = r.WithContext(context.WithValue(r.Context(), ParamsKey, p))
		}

		for i := 0; i < tnode.len; i++ {
			if tnode.handlers[i] != nil {
				tnode.handlers[i].ServeHTTP(w, r)

				goto healthCheck
			}

			if tnode.handlerFuncs[i] != nil {
				tnode.handlerFuncs[i](w, r)

				goto healthCheck
			}

			http.NotFound(w, r)

			return

		healthCheck:
			select {
			case <-r.Context().Done():
				return

			default:
			}
		}

		return
	}

	http.NotFound(w, r)
}

// SetCommon sets middlewares that run for all endpoints.
func (rt *Router) SetCommon(mids ...interface{}) {
	rt.common = mids
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

	mids = append(rt.common, mids...)

	if rt.methods[m] == nil {
		rt.methods[m] = radix.New(m)
	}

	// Filter middlewares before creating a node.
	for i := 0; i < len(mids); i++ {
		if _, ok := mids[i].(http.HandlerFunc); ok {
			continue
		}

		if _, ok := mids[i].(http.Handler); ok {
			continue
		}

		mids = append(mids[:i], mids[i+1:]...)
		i--
	}

	n := newNode(len(mids))

	for i, m := range mids {
		n.add(i, m)
	}

	rt.methods[m].Add(p, n)
}
