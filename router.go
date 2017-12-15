package httpmux

import (
	"context"
	"net/http"

	"github.com/gbrlsnchs/radix"
)

type Router struct {
	prefix  string
	methods map[string]*radix.Tree
}

func NewRouter() *Router {
	return &Router{methods: make(map[string]*radix.Tree)}
}

func (rt *Router) Handle(m, p string, h http.Handler) {
	rt.add(m, rt.prefix+resolvedPath(p), []interface{}{h})
}

func (rt *Router) HandleFunc(m, p string, hfunc http.HandlerFunc) {
	rt.add(m, rt.prefix+resolvedPath(p), []interface{}{hfunc})
}

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

func (rt *Router) Use(sub *Subrouter) {
	for endp, mids := range sub.endps {
		k := rt.prefix + endp

		for m, v := range mids {
			rt.add(m, k, v)
		}
	}
}

func (rt *Router) WithPrefix(p string) *Router {
	if p == "/" {
		return rt
	}

	rt.prefix = p

	return rt
}

func (rt *Router) add(m, p string, mids []interface{}) {
	if len(mids) == 0 || p == "" {
		return
	}

	if rt.methods[m] == nil {
		rt.methods[m] = radix.New(m)
	}

	rt.methods[m].Add(p, mids)
}
