package httpmux_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
	"github.com/stretchr/testify/assert"
)

func TestMux(t *testing.T) {
	a := assert.New(t)
	tests := []*struct {
		path           string
		mux            *Mux
		submuxes       []string
		expected       int
		expectedParams map[string]string
	}{
		// #0
		{
			path:     "/foo",
			mux:      New("foo"),
			expected: http.StatusOK,
		},
		// #1
		{
			path:     "/foo",
			mux:      New("/bar"),
			expected: http.StatusNotFound,
		},
		// #2
		{
			path:           "/foo/123",
			mux:            New("/foo/:bar"),
			expected:       http.StatusOK,
			expectedParams: map[string]string{"bar": "123"},
		},
		// #3
		{
			path:           "/foo/123/456",
			mux:            New("/foo/:bar/:baz"),
			expected:       http.StatusOK,
			expectedParams: map[string]string{"bar": "123", "baz": "456"},
		},
		// #4
		{
			path:           "/foo/123/456",
			mux:            New("/foo/:bar"),
			submuxes:       []string{"/:baz"},
			expected:       http.StatusOK,
			expectedParams: map[string]string{"bar": "123", "baz": "456"},
		},
		// #5
		{
			path:           "/foo/123/456",
			mux:            New("/foo"),
			submuxes:       []string{"/:bar", "/:baz"},
			expected:       http.StatusOK,
			expectedParams: map[string]string{"bar": "123", "baz": "456"},
		},
		// #6
		{
			path:     "/foo/123",
			mux:      New("/foo/:bar/:baz"),
			expected: http.StatusNotFound,
		},
	}

	for i, test := range tests {
		index := strconv.Itoa(i)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, test.path, nil)
		hfunc := func(w http.ResponseWriter, r *http.Request) {
			params, ok := r.Context().Value(Params).(map[string]string)

			if len(test.expectedParams) > 0 {
				a.NotNil(params, index)
				a.True(ok, index)
			}

			for k, v := range test.expectedParams {
				val := params[k]

				a.Exactly(v, val, index)
			}

			w.WriteHeader(http.StatusOK)
		}

		var smux *Submux

		for i := len(test.submuxes) - 1; i >= 0; i-- {
			tmp := smux
			smux = NewSubmux(test.submuxes[i])

			smux.HandleFunc(http.MethodGet, hfunc)

			if tmp == nil {
				continue
			}

			smux.Add(tmp)
		}

		if smux != nil {
			test.mux.Add(smux)
		} else {
			test.mux.HandleFunc(http.MethodGet, hfunc)
		}

		err := test.mux.Debug()

		a.Nil(err, index)

		test.mux.ServeHTTP(w, r)
		a.Exactly(test.expected, w.Code, index)
	}
}
