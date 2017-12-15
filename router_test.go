package httpmux_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	a := assert.New(t)
	tests := []*struct {
		path           string
		rt             *Router
		method         string
		endp           string
		mids           []interface{}
		expected       int
		expectedParams map[string]string
		sub            *Subrouter
	}{
		// #0
		{
			path:     "/foo",
			rt:       NewRouter(),
			method:   http.MethodGet,
			endp:     "/foo",
			expected: http.StatusOK,
		},
		// #1
		{
			path:     "/foo",
			rt:       NewRouter(),
			method:   http.MethodGet,
			expected: http.StatusNotFound,
		},
		// #2
		{
			path:           "/foo/123",
			rt:             NewRouter().WithPrefix("/foo/:bar"),
			method:         http.MethodGet,
			endp:           "/",
			expected:       http.StatusOK,
			expectedParams: map[string]string{"bar": "123"},
		},
		// #3
		{
			path:           "/foo/123/456",
			rt:             NewRouter(),
			method:         http.MethodGet,
			endp:           "/foo/:bar/:baz",
			expected:       http.StatusOK,
			expectedParams: map[string]string{"bar": "123", "baz": "456"},
		},
		// #4
		{
			path:           "/foo/123/456",
			rt:             NewRouter().WithPrefix("/foo/:bar"),
			method:         http.MethodGet,
			expected:       http.StatusOK,
			expectedParams: map[string]string{"bar": "123", "baz": "456"},
			sub:            NewSubrouter().WithPrefix("/:baz"),
		},
		// #5
		{
			path:           "/foo/123/456",
			rt:             NewRouter().WithPrefix("/foo"),
			method:         http.MethodGet,
			expected:       http.StatusOK,
			expectedParams: map[string]string{"bar": "123", "baz": "456"},
			sub:            NewSubrouter().WithPrefix("/:bar/:baz"),
		},
		// #6
		{
			path:     "/foo/123",
			rt:       NewRouter(),
			endp:     "/foo/:bar/:baz",
			expected: http.StatusNotFound,
		},
		// #7
		{
			path:     "/foo",
			rt:       NewRouter(),
			method:   http.MethodGet,
			endp:     "/foo",
			expected: http.StatusUnauthorized,
			mids: []interface{}{
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					Cancel(r)
				},
			},
		},
		// #8
		{
			path:     "/foo",
			rt:       NewRouter(),
			method:   http.MethodGet,
			endp:     "/foo",
			expected: http.StatusOK,
			mids: []interface{}{
				func(w http.ResponseWriter, r *http.Request) {
					// ignore this middleware
				},
			},
		},
	}

	for i, test := range tests {
		index := strconv.Itoa(i)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, test.path, nil)
		midsToAdd := make([]interface{}, 0)
		midsToAdd = append(midsToAdd, test.mids...)
		midsToAdd = append(midsToAdd, func(w http.ResponseWriter, r *http.Request) {
			p, ok := r.Context().Value(Params).(map[string]string)

			if len(test.expectedParams) > 0 {
				a.NotNil(p, index)
				a.True(ok, index)
			}

			for k, v := range test.expectedParams {
				val := test.expectedParams[k]

				a.Exactly(v, val, index)
			}

			w.WriteHeader(http.StatusOK)
		})

		if test.sub != nil {
			test.sub.HandleMiddlewares(test.method, test.endp, midsToAdd...)
			test.rt.Use(test.sub)
		} else {
			test.rt.HandleMiddlewares(test.method, test.endp, midsToAdd...)
		}

		test.rt.ServeHTTP(w, r)
		a.Exactly(test.expected, w.Code, index)
	}
}
