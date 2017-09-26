package httpmux

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	handler func(http.ResponseWriter, *http.Request)
}

func (th *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.handler(w, r)
}

func TestNewMultiplexer(t *testing.T) {
	a := assert.New(t)
	tests := []*struct {
		endpoint   string
		method     string
		mux        Multiplexer
		muxMethod  string
		params     []string
		paramsVals []string
		subs       []string
		flag       bool
		rHeader    [2]string
		expected   bool
	}{
		// #0
		{
			endpoint: "/test",
			method:   http.MethodGet,
			mux:      New("/test"),
			expected: true,
		},
		// #1
		{
			endpoint:  "/test",
			method:    http.MethodGet,
			mux:       New("/test"),
			muxMethod: http.MethodGet,
			expected:  true,
		},
		// #2
		{
			endpoint:  "/test",
			method:    http.MethodGet,
			mux:       New("/test"),
			muxMethod: http.MethodPost,
			expected:  false,
		},
		// #3
		{
			endpoint: "/test/test",
			method:   http.MethodGet,
			mux:      New("/test"),
			expected: false,
		},
		// #4
		{
			endpoint:  "/test/test",
			method:    http.MethodGet,
			mux:       New("/test"),
			muxMethod: http.MethodGet,
			expected:  false,
		},
		// #5
		{
			endpoint:  "/test/test",
			method:    http.MethodGet,
			mux:       New("/test"),
			muxMethod: http.MethodPost,
			expected:  false,
		},
		// #6
		{
			endpoint: "/test/test",
			method:   http.MethodGet,
			mux:      New("/test").SetSubmux(New("/test")),
			subs:     []string{"/test"},
			expected: true,
		},
		// #7
		{
			endpoint:  "/test/test",
			method:    http.MethodGet,
			mux:       New("/test").SetSubmux(New("/test")),
			muxMethod: http.MethodGet,
			subs:      []string{"/test"},
			expected:  true,
		},
		// #8
		{
			endpoint:  "/test/test",
			method:    http.MethodGet,
			mux:       New("/test").SetSubmux(New("/test")),
			muxMethod: http.MethodPost,
			subs:      []string{"/test"},
			expected:  false,
		},
		// #9
		{
			endpoint:   "/test/123",
			method:     http.MethodGet,
			mux:        New("/test").SetSubmux(New("/{test}")),
			params:     []string{"test"},
			paramsVals: []string{"123"},
			subs:       []string{"/{test}"},
			expected:   true,
		},
		// #10
		{
			endpoint:   "/test/123",
			method:     http.MethodGet,
			mux:        New("/test").SetSubmux(New("/{test}")),
			muxMethod:  http.MethodGet,
			params:     []string{"test"},
			paramsVals: []string{"123"},
			subs:       []string{"/{test}"},
			expected:   true,
		},
		// #11
		{
			endpoint:   "/test/123",
			method:     http.MethodGet,
			mux:        New("/test").SetSubmux(New("/{test}")),
			muxMethod:  http.MethodPost,
			params:     []string{"test"},
			paramsVals: []string{"123"},
			subs:       []string{"/{test}"},
			expected:   false,
		},
		// #12
		{
			endpoint:   "/test/abc",
			method:     http.MethodGet,
			mux:        New("/test").SetSubmux(New("/{test:[0-9]+}")),
			params:     []string{"test"},
			paramsVals: []string{"abc"},
			subs:       []string{"/{test:[0-9]+}"},
			expected:   false,
		},
		// #13
		{
			endpoint:   "/test/abc",
			method:     http.MethodGet,
			mux:        New("/test").SetSubmux(New("/{test:[0-9]+}")),
			muxMethod:  http.MethodGet,
			params:     []string{"test"},
			paramsVals: []string{"abc"},
			subs:       []string{"/{test:[0-9]+}"},
			expected:   false,
		},
		// #14
		{
			endpoint:   "/test/abc",
			method:     http.MethodGet,
			mux:        New("/test").SetSubmux(New("/{test:[0-9]+}")),
			muxMethod:  http.MethodPost,
			params:     []string{"test"},
			paramsVals: []string{"abc"},
			subs:       []string{"/{test:[0-9]+}"},
			expected:   false,
		},
		// #15
		{
			endpoint:   "/test/golang/123",
			method:     http.MethodGet,
			mux:        New("/test").SetSubmux(New("/{name}").SetSubmux(New("/{value}"))),
			params:     []string{"name", "value"},
			paramsVals: []string{"golang", "123"},
			subs:       []string{"/{name}", "/{value}"},
			expected:   true,
		},
		// #16
		{
			endpoint:   "/test/golang/123",
			method:     http.MethodGet,
			mux:        New("/test").SetSubmux(New("/{name}").SetSubmux(New("/{value}"))),
			muxMethod:  http.MethodGet,
			params:     []string{"name", "value"},
			paramsVals: []string{"golang", "123"},
			subs:       []string{"/{name}", "/{value}"},
			expected:   true,
		},
		// #17
		{
			endpoint:   "/test/golang/123",
			method:     http.MethodGet,
			mux:        New("/test").SetSubmux(New("/{name}").SetSubmux(New("/{value}"))),
			muxMethod:  http.MethodPost,
			params:     []string{"name", "value"},
			paramsVals: []string{"golang", "123"},
			subs:       []string{"/{name}", "/{value}"},
			expected:   false,
		},
		// #18
		{
			endpoint: "/test/recursive/new",
			method:   http.MethodGet,
			mux:      New("/test/recursive/new"),
			subs:     []string{"/recursive", "/new"},
			expected: true,
		},
		// #19
		{
			endpoint:  "/test/recursive/new",
			method:    http.MethodGet,
			mux:       New("/test/recursive/new"),
			muxMethod: http.MethodGet,
			subs:      []string{"/recursive", "/new"},
			expected:  true,
		},
		// #20
		{
			endpoint:  "/test/recursive/new",
			method:    http.MethodGet,
			mux:       New("/test/recursive/new"),
			muxMethod: http.MethodPost,
			subs:      []string{"/recursive", "/new"},
			expected:  false,
		},
		// #21
		{
			endpoint: "/test/recursive/new",
			method:   http.MethodGet,
			mux:      New("/test/recursive/old"),
			subs:     []string{"/recursive", "/old"},
			expected: false,
		},
		// #22
		{
			endpoint:  "/test/recursive/new",
			method:    http.MethodGet,
			mux:       New("/test/recursive/old"),
			muxMethod: http.MethodGet,
			subs:      []string{"/recursive", "/old"},
			expected:  false,
		},
		// #23
		{
			endpoint:  "/test/recursive/new",
			method:    http.MethodGet,
			mux:       New("/test/recursive/old"),
			muxMethod: http.MethodPost,
			subs:      []string{"/recursive", "/old"},
			expected:  false,
		},
		// #24
		{
			endpoint:   "/test/123/456",
			method:     http.MethodGet,
			mux:        New("/test/{val1:[0-9]+}/{val2}"),
			params:     []string{"val1", "val2"},
			paramsVals: []string{"123", "456"},
			subs:       []string{"/{val1:[0-9]+}", "/{val2}"},
			expected:   true,
		},
		// #25
		{
			endpoint:   "/test/123/456",
			method:     http.MethodGet,
			mux:        New("/test/{val1:[0-9]+}/{val2}"),
			muxMethod:  http.MethodGet,
			params:     []string{"val1", "val2"},
			paramsVals: []string{"123", "456"},
			subs:       []string{"/{val1:[0-9]+}", "/{val2}"},
			expected:   true,
		},
		// #26
		{
			endpoint:   "/test/123/456",
			method:     http.MethodGet,
			mux:        New("/test/{val1:[0-9]+}/{val2}"),
			muxMethod:  http.MethodPost,
			params:     []string{"val1", "val2"},
			paramsVals: []string{"123", "456"},
			subs:       []string{"/{val1:[0-9]+}", "/{val2}"},
			expected:   false,
		},
		// #27
		{
			endpoint:   "/test/123/456",
			method:     http.MethodGet,
			mux:        New("/test/{val1:[a-zA-Z]+}/{val2}"),
			params:     []string{"val1", "val2"},
			paramsVals: []string{"123", "456"},
			subs:       []string{"/{val1:[0-9]+}", "/{val2}"},
			expected:   false,
		},
		// #28
		{
			endpoint:   "/test/123/456",
			method:     http.MethodGet,
			mux:        New("/test/{val1:[a-zA-Z]+}/{val2}"),
			muxMethod:  http.MethodGet,
			params:     []string{"val1", "val2"},
			paramsVals: []string{"123", "456"},
			subs:       []string{"/{val1:[0-9]+}", "/{val2}"},
			expected:   false,
		},
		// #29
		{
			endpoint:   "/test/123/456",
			method:     http.MethodGet,
			mux:        New("/test/{val1:[a-zA-Z]+}/{val2}"),
			muxMethod:  http.MethodPost,
			params:     []string{"val1", "val2"},
			paramsVals: []string{"123", "456"},
			subs:       []string{"/{val1:[0-9]+}", "/{val2}"},
			expected:   false,
		},
		// #30
		{
			endpoint: "/test/header",
			method:   http.MethodGet,
			mux:      New("/test/header").SetHeader("Content-Type", "application/json"),
			subs:     []string{"/header"},
			rHeader:  [2]string{"Content-Type", "application/json"},
			expected: true,
		},
		// #31
		{
			endpoint: "/test/header",
			method:   http.MethodGet,
			mux:      New("/test/header").SetHeader("Content-Type", "application/json"),
			subs:     []string{"/header"},
			expected: false,
		},
		// #32
		{
			endpoint: "/test/header",
			method:   http.MethodGet,
			mux:      New("/test/header").SetHeader("Content-Type", "application/json"),
			subs:     []string{"/header"},
			rHeader:  [2]string{"Content-Type", "application/xml"},
			expected: false,
		},
		// #33
		{
			endpoint: "/test/header",
			method:   http.MethodGet,
			mux:      New("/test").SetSubmux(New("/header").SetHeader("Content-Type", "application/json")),
			subs:     []string{"/header"},
			rHeader:  [2]string{"Content-Type", "application/json"},
			expected: true,
		},
		// #34
		{
			endpoint: "/test/header",
			method:   http.MethodGet,
			mux:      New("/test").SetSubmux(New("/header").SetHeader("Content-Type", "application/json")),
			subs:     []string{"/header"},
			expected: false,
		},
		// #35
		{
			endpoint: "/test/header",
			method:   http.MethodGet,
			mux:      New("/test").SetSubmux(New("/header").SetHeader("Content-Type", "application/json")),
			subs:     []string{"/header"},
			rHeader:  [2]string{"Content-Type", "application/xml"},
			expected: false,
		},
		// #36
		{
			endpoint: "/test/header",
			method:   http.MethodGet,
			mux: New("/test").SetHeader("Content-Type", "application/xml").
				SetSubmux(New("/header").SetHeader("Content-Type", "application/json")),
			subs:     []string{"/header"},
			rHeader:  [2]string{"Content-Type", "application/json"},
			expected: true,
		},
		// #37
		{
			endpoint: "/test/header",
			method:   http.MethodGet,
			mux: New("/test").SetHeader("Content-Type", "application/xml").
				SetSubmux(New("/header").SetHeader("Content-Type", "application/json")),
			subs:     []string{"/header"},
			rHeader:  [2]string{"Content-Type", "application/xml"},
			expected: false,
		},
	}

	for i, test := range tests {
		index := strconv.Itoa(i)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(test.method, test.endpoint, nil)

		if test.rHeader[0] != "" && test.rHeader[1] != "" {
			r.Header.Set(test.rHeader[0], test.rHeader[1])
		}

		th := &testHandler{handler: func(w http.ResponseWriter, r *http.Request) {
			for i, p := range test.params {
				v := r.Context().Value(p)

				a.Exactly(test.paramsVals[i], v, index)
			}

			test.flag = true
		}}

		method := test.muxMethod

		if method == "" {
			method = All
		}

		if len(test.subs) > 0 {
			subm := test.mux

			for _, s := range test.subs {
				subm = subm.Submux(s)
			}

			subm.SetHandler(th, method)
		} else {
			test.mux.SetHandler(th, method)
		}

		a.NotNil(test.mux, index)
		test.mux.ServeHTTP(w, r)
		a.Exactly(test.expected, test.flag, index)
	}
}
