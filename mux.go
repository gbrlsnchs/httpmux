package httpmux

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type Mux struct {
	path     string
	regexp   string
	header   http.Header
	handlers []*Handler
	children []*Mux
	dynChild *Mux
	parent   *Mux
}

func New(path string) *Mux {
	startIndex := strings.IndexByte(path, '/')
	dynIndex := strings.IndexByte(path, ':')
	endIndex := strings.IndexByte(path[startIndex+1:], '/')

	if endIndex >= 0 {
		endIndex = endIndex + 1

	} else {
		endIndex = len(path)
	}

	regexpIndex := strings.IndexByte(path[dynIndex+1:], ':')

	if regexpIndex >= 0 {
		regexpIndex = regexpIndex + 2

		if regexpIndex > endIndex {
			regexpIndex = regexpIndex - endIndex
		}
	} else {
		regexpIndex = endIndex
	}

	m := &Mux{path: path[startIndex+1 : regexpIndex]}

	if regexpIndex != endIndex {
		m.regexp = path[regexpIndex+1 : endIndex]
	}

	if endIndex != len(path) {
		child := New(path[endIndex:])

		if strings.IndexByte(path[endIndex:], ':') >= 0 {
			m.dynChild = child
		} else {
			m.children = []*Mux{child}
		}

		child.parent = m

		return child
	}

	return m
}

func (m *Mux) Lookup(r *http.Request) (*Mux, error) {
	var (
		n      *Mux
		header http.Header
	)
	startIndex := strings.IndexByte(r.URL.Path, '/')
	endIndex := strings.IndexByte(r.URL.Path[startIndex+1:], '/')

	if endIndex < 0 {
		endIndex = len(r.URL.Path)
	} else {
		endIndex = endIndex + 1
	}

	if len(m.header) > 0 {
		header = n.header
	}

	if m.path == r.URL.Path[startIndex+1:endIndex] {
		n = m

		goto walk
	}

	if i := strings.IndexByte(m.path, ':'); i >= 0 {
		if m.regexp != "" {
			match, err := regexp.MatchString(m.regexp, r.URL.Path)

			if err != nil {
				return nil, err
			}

			if !match {
				return nil, nil
			}
		}

		n = m
		*r = *r.WithContext(n.ctx(r.Context(), r.URL.Path[startIndex+1:endIndex]))
	}

walk:
	startIndex = strings.IndexByte(r.URL.Path[endIndex:], '/') + endIndex
	endIndex = strings.IndexByte(r.URL.Path[startIndex+1:], '/')

	if endIndex < 0 {
		endIndex = len(r.URL.Path)
	} else {
		endIndex = endIndex + startIndex + 1
	}

	if n != nil && len(r.URL.Path[startIndex+1:endIndex]) > 0 {
		for _, c := range n.children {
			if c.path == r.URL.Path[startIndex+1:endIndex] {
				n = c

				if len(n.header) > 0 {
					header = n.header
				}

				goto walk
			}
		}

		if n.dynChild != nil {
			if n.dynChild.regexp != "" {
				match, err := regexp.MatchString(n.dynChild.regexp, r.URL.Path[startIndex+1:endIndex])

				if err != nil {
					return nil, err
				}

				if !match {
					return nil, nil
				}
			}

			n = n.dynChild

			if len(n.header) > 0 {
				header = n.header
			}

			*r = *r.WithContext(n.ctx(r.Context(), r.URL.Path[startIndex+1:endIndex]))

			goto walk
		}

		return nil, nil
	}

	for k := range header {
		if header.Get(k) != r.Header.Get(k) {
			return nil, nil
		}
	}

	return n, nil
}

func (m *Mux) Parent() *Mux {
	return m.parent
}

func (m *Mux) Root() *Mux {
	for m.parent != nil {
		return m.parent.Root()
	}

	return m
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n, err := m.Lookup(r)

	if err != nil {
		PanicFunc(err)

		return
	}

	if n != nil {
		for _, hndr := range n.handlers {
			if hndr.method == r.Method {
				hndr.h.ServeHTTP(w, r)

				return
			}
		}
	}
}

func (m *Mux) SetHandler(method string, h http.Handler) *Mux {
	m.handlers = append(m.handlers, &Handler{h: h, method: method})

	return m
}

func (m *Mux) SetHeader(k string, v string) *Mux {
	if m.header == nil {
		m.header = make(http.Header)
	}

	m.header.Set(k, v)

	return m
}

func (m *Mux) SetSubmux(sm *Mux) *Mux {
	if i := strings.IndexByte(sm.path, ':'); i >= 0 {
		m.dynChild = sm

		return m
	}

	m.children = append(m.children, sm)

	return sm
}

func (m *Mux) ctx(ctx context.Context, val string) context.Context {
	ctx = context.WithValue(
		ctx,
		m.path[1:],
		val,
	)

	return ctx
}
