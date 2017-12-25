package httpmux_test

import (
	"net/http"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
)

var responseNotFound = []byte("404 page not found\n")

func subrouterHelperWithHandler(status int, response string, hasHandler bool) *Subrouter {
	subr := NewSubrouter()
	h := &handlerMockup{status: status, response: response}

	if hasHandler {
		subr.Handle(http.MethodGet, "/", h)

		return subr
	}

	subr.HandleFunc(http.MethodGet, "/", h.ServeHTTP)

	return subr
}

func TestEmptySubrouter(t *testing.T) {
	w, r := responseRequest()

	routerHelper(NewSubrouter()).ServeHTTP(w, r)

	body := w.Body.Bytes()

	testHTTPResponse(t, http.StatusNotFound, w.Code, responseNotFound, body)
}

func TestSubrouterHandle(t *testing.T) {
	testSubrouter(t, true)
}

func TestSubrouterHandleFunc(t *testing.T) {
	testSubrouter(t, false)
}
