package httpmux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
)

var responseNotFound = []byte("404 page not found\n")

func subrouterHelperWithHandler(status int, response string) *Subrouter {
	subr := NewSubrouter()

	subr.Handle(http.MethodGet, "/", &handlerMockup{status: status, response: response})

	return subr
}

func TestEmptySubrouter(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	routerHelper(NewSubrouter()).ServeHTTP(w, r)

	body := w.Body.Bytes()

	testHTTPResponse(t, http.StatusNotFound, w.Code, responseNotFound, body)
}

func TestSubrouterHandle(t *testing.T) {
	testTable := []struct {
		obj              *Subrouter
		expectedStatus   int
		expectedResponse []byte
	}{
		{subrouterHelperWithHandler(http.StatusOK, "foobar"), http.StatusOK, []byte("foobar")},
	}

	for _, tt := range testTable {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		routerHelper(tt.obj).ServeHTTP(w, r)

		body := w.Body.Bytes()

		testHTTPResponse(t, tt.expectedStatus, w.Code, tt.expectedResponse, body)
	}
}
