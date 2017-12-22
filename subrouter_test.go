package httpmux_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
)

var responseNotFound = []byte("404 page not found\n")

func subrouterHelperWithHandler(status int, response string) *Subrouter {
	subr := NewSubrouter()

	subr.Handle(http.MethodGet, "/", &handlerMock{status, response})

	return subr
}

func TestSubrouterHandle(t *testing.T) {
	testTable := []struct {
		obj              *Subrouter
		expectedStatus   int
		expectedResponse []byte
	}{
		{NewSubrouter(), http.StatusNotFound, responseNotFound},
		{subrouterHelperWithHandler(http.StatusOK, "foobar"), http.StatusOK, []byte("foobar")},
	}

	for _, tt := range testTable {
		rt := routerHelper(tt.obj)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		rt.ServeHTTP(w, r)

		if w.Code != tt.expectedStatus {
			t.Errorf("%d != %d\n", w.Code, tt.expectedStatus)
		}

		body := w.Body.Bytes()

		if !bytes.Equal(body, tt.expectedResponse) {
			t.Errorf("%s != %s\n",
				strings.TrimSuffix(string(body), "\n"),
				strings.TrimSuffix(string(tt.expectedResponse), "\n"),
			)
		}
	}
}
