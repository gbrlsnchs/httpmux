package httpmux_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
	. "github.com/gbrlsnchs/httpmux/internal"
)

func TestEmptySubrouter(t *testing.T) {
	expectedStatus := http.StatusNotFound
	expectedResponse := []byte("404 page not found\n")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rt := NewRouter()

	rt.Use(NewSubrouter())
	rt.ServeHTTP(w, r)

	body := w.Body.Bytes()

	if expectedStatus != w.Code {
		t.Errorf("%d is not expected status (%d)\n", w.Code, expectedStatus)
	}

	if !bytes.Equal(expectedResponse, body) {
		t.Errorf("%x is not expected response (%x)\n", body, expectedResponse)
	}
}

func TestSubrouterHandle(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedResponse := []byte("foobar")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	subr := NewSubrouter()
	rt := NewRouter()
	h := &DummyHandler{Status: expectedStatus, Response: expectedResponse}

	subr.Handle(http.MethodGet, "/", h)
	rt.Use(subr)
	rt.ServeHTTP(w, r)

	body := w.Body.Bytes()

	if expectedStatus != w.Code {
		t.Errorf("%d is not expected status (%d)\n", w.Code, expectedStatus)
	}

	if !bytes.Equal(expectedResponse, body) {
		t.Errorf("%x is not expected response (%x)\n", body, expectedResponse)
	}

	if !h.Finished {
		t.Error("http.Handler has not run")
	}
}

func TestSubrouterHandleFunc(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedResponse := []byte("foobar")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	subr := NewSubrouter()
	rt := NewRouter()
	h := &DummyHandler{Status: expectedStatus, Response: expectedResponse}

	subr.HandleFunc(http.MethodGet, "/", h.ServeHTTP)
	rt.Use(subr)
	rt.ServeHTTP(w, r)

	body := w.Body.Bytes()

	if expectedStatus != w.Code {
		t.Errorf("%d is not expected status (%d)\n", w.Code, expectedStatus)
	}

	if !bytes.Equal(expectedResponse, body) {
		t.Errorf("%x is not expected response (%x)\n", body, expectedResponse)
	}

	if !h.Finished {
		t.Error("http.HandlerFunc has not run")
	}
}

func TestSubrouterWithCancel(t *testing.T) {
	expectedStatus := http.StatusBadRequest
	expectedResponse := []byte{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	subr := NewSubrouter()
	rt := NewRouter()
	h := &DummyHandler{Status: expectedStatus, Response: expectedResponse}

	subr.HandleMiddlewares(http.MethodGet, "/", h.Cancel, h.ServeHTTP)
	rt.Use(subr)
	rt.ServeHTTP(w, r)

	body := w.Body.Bytes()

	if expectedStatus != w.Code {
		t.Errorf("%d is not expected status (%d)\n", w.Code, expectedStatus)
	}

	if !bytes.Equal(expectedResponse, body) {
		t.Errorf("%x is not expected response (%x)\n", body, expectedResponse)
	}

	if !h.Canceled {
		t.Error("http.HandlerFunc has not been canceled")
	}

	if h.Finished {
		t.Error("http.HandlerFunc has run")
	}
}

func TestSubrouterHandleWithParams(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedResponse := map[string]string{"foo": "123", "bar": "456"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/123/456", nil)
	subr := NewSubrouter()
	rt := NewRouter()
	h := &DummyHandler{Status: expectedStatus, ReturnParams: true}

	subr.Handle(http.MethodGet, "/:foo/:bar", h)
	rt.SetParamsKey(ParamsKey)
	rt.Use(subr)
	rt.ServeHTTP(w, r)

	body := strings.Split(w.Body.String(), "/")
	actualResponse := map[string]string{}

	for i := 0; i < len(body); i += 2 {
		actualResponse[body[i]] = body[i+1]
	}

	if expectedStatus != w.Code {
		t.Errorf("%d is not expected status (%d)\n", w.Code, expectedStatus)
	}

	if !reflect.DeepEqual(expectedResponse, actualResponse) {
		t.Errorf("%v is not expected response (%v)\n", actualResponse, expectedResponse)
	}

	if !h.Finished {
		t.Error("http.Handler has not run")
	}
}
