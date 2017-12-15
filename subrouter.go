package httpmux

import "net/http"

type Subrouter struct {
	prefix string
	endps  map[string]map[string][]interface{}
}

func NewSubrouter() *Subrouter {
	return &Subrouter{endps: make(map[string]map[string][]interface{})}
}

func (s *Subrouter) Handle(m, p string, h http.Handler) {
	k := s.prefix + resolvedPath(p)

	s.initEndp(k)

	s.endps[k][m] = []interface{}{h}
}

func (s *Subrouter) HandleFunc(m, p string, hfunc http.HandlerFunc) {
	k := s.prefix + resolvedPath(p)

	s.initEndp(k)

	s.endps[k][m] = []interface{}{hfunc}
}

func (s *Subrouter) HandleMiddlewares(m, p string, mids ...interface{}) {
	k := s.prefix + resolvedPath(p)

	s.initEndp(k)

	for _, mid := range mids {
		if h, ok := mid.(http.Handler); ok {
			s.endps[k][m] = append(s.endps[k][m], h)

			continue
		}

		if hfunc, ok := mid.(func(http.ResponseWriter, *http.Request)); ok {
			s.endps[k][m] = append(s.endps[k][m], http.HandlerFunc(hfunc))
		}
	}
}

func (s *Subrouter) Use(sub *Subrouter) {
	for endp, mids := range sub.endps {
		k := s.prefix + endp

		s.initEndp(k)

		for m, v := range mids {
			s.endps[k][m] = append(s.endps[k][m], v)
		}
	}
}

func (s *Subrouter) WithPrefix(p string) *Subrouter {
	if p == "/" {
		return s
	}

	s.prefix = p

	return s
}

func (s *Subrouter) initEndp(e string) {
	if s.endps[e] == nil {
		s.endps[e] = make(map[string][]interface{})
	}
}
