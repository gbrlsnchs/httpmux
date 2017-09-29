package httpmux_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	handler func(w http.ResponseWriter, r *http.Request)
}

func (th *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.handler(w, r)
}

func TestMuxLookup(t *testing.T) {
	a := assert.New(t)

	tests := []*struct {
		path               string
		mux                *Mux
		expected           bool
		expectedParams     []string
		expectedParamsVals []string
	}{
		{
			path:     "/test",
			mux:      New("/test"),
			expected: true,
		},
		{
			path:     "/test",
			mux:      New("/:testing"),
			expected: true,
		},
		{
			path:     "/testing",
			mux:      New("/test"),
			expected: false,
		},
		{
			path:               "/123/test",
			mux:                New("/:testing/test"),
			expected:           true,
			expectedParams:     []string{"testing"},
			expectedParamsVals: []string{"123"},
		},
		{
			path:     "/testing/123/nope",
			mux:      New("/testing/:anything"),
			expected: false,
		},
	}

	for i, test := range tests {
		index := strconv.Itoa(i)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, test.path, nil)

		flag := false

		test.mux.SetHandler(
			http.MethodGet,
			&testHandler{
				func(w http.ResponseWriter, r *http.Request) {
					flag = true

					for i := range test.expectedParams {
						val := r.Context().Value(test.expectedParams[i])

						a.Exactly(test.expectedParamsVals[i], val, index)
					}
				},
			},
		)

		test.mux.Root().ServeHTTP(w, r)

		a.Exactly(test.expected, flag, index)
	}
}
