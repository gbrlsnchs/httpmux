package httpmux_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
)

func responseRequest() (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
}

func testCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("%d != %d\n", expected, actual)
	}
}

func testHTTPResponse(t *testing.T, expectedCode, actualCode int, expectedResponse, actualResponse []byte) {
	testCode(t, expectedCode, actualCode)
	testResponse(t, expectedResponse, actualResponse)
}

func testResponse(t *testing.T, expected []byte, actual []byte) {
	if !bytes.Equal(expected, actual) {
		t.Errorf("%s != %s\n",
			strings.TrimSuffix(string(expected), "\n"),
			strings.TrimSuffix(string(actual), "\n"),
		)
	}
}

func testRouter(t *testing.T, hasHandler bool) {
	testTable := []struct {
		obj              *Router
		expectedStatus   int
		expectedResponse []byte
	}{
		{
			routerHelperWithHandler(http.StatusOK, "foobar", hasHandler),
			http.StatusOK,
			[]byte("foobar"),
		},
		{
			routerHelperWithHandler(http.StatusBadRequest, "bazqux", hasHandler),
			http.StatusBadRequest,
			[]byte("bazqux"),
		},
	}

	for _, tt := range testTable {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		tt.obj.ServeHTTP(w, r)

		body := w.Body.Bytes()

		testHTTPResponse(t, tt.expectedStatus, w.Code, tt.expectedResponse, body)
	}
}

func testSubrouter(t *testing.T, hasHandler bool) {
	testTable := []struct {
		obj              *Subrouter
		expectedStatus   int
		expectedResponse []byte
	}{
		{
			subrouterHelperWithHandler(http.StatusOK, "foobar", hasHandler),
			http.StatusOK,
			[]byte("foobar"),
		},
		{
			subrouterHelperWithHandler(http.StatusBadRequest, "bazqux", hasHandler),
			http.StatusBadRequest,
			[]byte("bazqux"),
		},
	}

	for _, tt := range testTable {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		routerHelper(tt.obj).ServeHTTP(w, r)

		body := w.Body.Bytes()

		testHTTPResponse(t, tt.expectedStatus, w.Code, tt.expectedResponse, body)
	}
}
