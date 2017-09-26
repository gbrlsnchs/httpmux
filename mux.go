package httpmux

import (
	"context"
	"errors"
	"net/http"
	"path"
	"regexp"
	"strings"
)

// Multiplexer represents an HTTP request mux.
type Multiplexer interface {
	http.Handler
	// Loopkup finds the match for the request.
	Lookup(*http.Request) (Multiplexer, error)
	// NewSubmux sets a new submux and returns it.
	NewSubmux(string) Multiplexer
	// Path returns the mux's path without the trailing slash and regexps.
	Path() string
	// Regexp returns the mux's regexp contained in the path used to build the mux.
	Regexp() string
	// SetHandler sets a new handler for the mux.
	SetHandler(http.Handler, ...string) Multiplexer
	// SetHeader sets a request header filter.
	SetHeader(string, string) Multiplexer
	// SetSubmux sets a new submux.
	SetSubmux(Multiplexer) Multiplexer
	// Submux returns a direct child of the caller mux.
	// If none is found, it returns nil.
	Submux(string) Multiplexer
}

// mux is a Radix Tree node.
type mux struct {
	// path works as a Radix Tree node label.
	path string
	// regexp holds a pattern to be matched in a dynamic route.
	regexp string
	// header holds a map for matching with the request's header.
	header http.Header
	// handlerMap holds all handlers to be executed depending on the HTTP method.
	handlerMap map[string]http.Handler
	// nodes are all child nodes the mux (node) holds.
	nodes map[string]Multiplexer
	// dynNode represents a special child node that holds a dynamic path.
	// Each mux (or node) can only hold one dynamic path.
	dynNode Multiplexer
}

// New creates a new Multiplexer.
//
// If the path passed as parameter has more than one path, i.e.:
//  httpmux.New("/more/than/one/path")
// it will automatically be recursively broken into submuxes.
func New(p string) Multiplexer {
	spath := strings.Split(strings.TrimPrefix(p, "/"), "/")

	m := &mux{}

	if r, p := extRegexp(spath[0]); r != "" {
		m.regexp = r
		spath[0] = p
	}

	m.path = spath[0]

	if len(spath) > 1 {
		if dynPath(spath[1]) {
			m.dynNode = New(path.Join(spath[1:]...))

			return m
		}

		if m.nodes == nil {
			m.nodes = make(map[string]Multiplexer)
		}

		m.nodes[spath[1]] = New(path.Join(spath[1:]...))
	}

	return m
}

// Lookup searches for the root mux and all submuxes using a Radix Tree algorithm.
func (m *mux) Lookup(r *http.Request) (Multiplexer, error) {
	n := m
	parentHeader := n.header
	found := 0
	p := r.URL.Path
	spath := strings.Split(strings.TrimPrefix(p, "/"), "/")
	key := spath[found]

	// before entering the loop it checks the parent mux,
	// which is the root of the tree
	if n.path == key {
		found = found + 1
	}

	if dynPath(n.path) && found == 0 {
		*r = *r.WithContext(n.ctx(r.Context(), key))

		if n.regexp != "" {
			match, err := regexp.MatchString(n.regexp, key)

			if err != nil {
				return nil, err
			}

			if !match {
				return nil, nil
			}
		}

		found = found + 1
	}

	for found < len(spath) {
		key = spath[found]

		if n.dynNode != nil {
			n = n.dynNode.(*mux)
			*r = *r.WithContext(n.ctx(r.Context(), key))

			if len(n.header) > 0 {
				parentHeader = n.header
			}

			if n.regexp != "" {
				match, err := regexp.MatchString(n.regexp, key)

				if err != nil {
					return nil, err
				}

				if !match {
					return nil, nil
				}
			}

			found = found + 1

			continue
		}

		if n.nodes[key] != nil {
			n = n.nodes[key].(*mux)

			if len(n.header) > 0 {
				parentHeader = n.header
			}

			found = found + 1

			continue
		}

		return nil, nil
	}

	for k := range parentHeader {
		if parentHeader.Get(k) != r.Header.Get(k) {
			return nil, nil
		}
	}

	return n, nil
}

// NewSubmux sets a new submux and returns it.
func (m *mux) NewSubmux(p string) Multiplexer {
	sm := New(p)

	m.SetSubmux(sm)

	return sm
}

// Path returns the mux's path without the trailing slash and regexps.
func (m *mux) Path() string {
	return m.path
}

// Regexp returns the mux's regexp contained in the path used to build the mux.
func (m *mux) Regexp() string {
	return m.regexp
}

// ServeHTTP calls a handler's ServeHTTP function if:
// 1) The Lookup function returns a not nil node;
// 2) The handler function is not nil for the requested HTTP method.
func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n, err := m.Lookup(r)

	if err != nil {
		PanicFunc(err)

		return
	}

	if n != nil {
		mx, ok := n.(*mux)

		if !ok {
			PanicFunc(errors.New("httpmux: Multiplexer conversion to *mux"))
		}

		if mx.handlerMap[All] != nil {
			mx.handlerMap[All].ServeHTTP(w, r)

			return
		}

		if mx.handlerMap[r.Method] != nil {
			mx.handlerMap[r.Method].ServeHTTP(w, r)
		}
	}
}

// SetHandler sets a new handler for the mux.
func (m *mux) SetHandler(h http.Handler, methods ...string) Multiplexer {
	if m.handlerMap == nil {
		m.handlerMap = make(map[string]http.Handler)
	}

	if len(methods) == 0 {
		m.handlerMap[All] = h

		return m
	}

	for _, method := range methods {
		m.handlerMap[method] = h
	}

	return m
}

// SetHeader sets a request header filter.
func (m *mux) SetHeader(k string, v string) Multiplexer {
	if m.header == nil {
		m.header = make(http.Header)
	}

	m.header.Set(k, v)

	return m
}

// SetSubmux sets a new submux.
func (m *mux) SetSubmux(sm Multiplexer) Multiplexer {
	p := sm.Path()

	if dynPath(p) {
		m.dynNode = sm

		return m
	}

	if m.nodes == nil {
		m.nodes = make(map[string]Multiplexer)
	}

	m.nodes[p] = sm

	return m
}

// Submux returns a direct child of the caller mux.
// If none is found, it returns nil.
func (m *mux) Submux(p string) Multiplexer {
	p = strings.TrimPrefix(p, "/")

	if r, extp := extRegexp(p); r != "" {
		p = extp
	}

	if dynPath(p) {
		return m.dynNode
	}

	return m.nodes[p]
}

// ctx builds a context based on an already existent context
// and resolves its name based on the mux's path.
func (m *mux) ctx(ctx context.Context, val string) context.Context {
	ctxName := m.path[1 : len(m.path)-1]

	if i := strings.IndexByte(ctxName, ':'); i >= 0 {
		ctxName = ctxName[:i]
	}

	ctx = context.WithValue(
		ctx,
		ctxName,
		val,
	)

	return ctx
}
