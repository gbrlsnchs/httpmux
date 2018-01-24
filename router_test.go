package httpmux_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
)

func TestEmptyRouter(t *testing.T) {
	expectedStatus := http.StatusNotFound
	expectedResponse := []byte("404 page not found\n")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rt := NewRouter()

	rt.ServeHTTP(w, r)

	body := w.Body.Bytes()

	if expectedStatus != w.Code {
		t.Errorf("%d is not expected status (%d)\n", w.Code, expectedStatus)
	}

	if !bytes.Equal(expectedResponse, body) {
		t.Errorf("%x is not expected response (%x)\n", body, expectedResponse)
	}
}

func TestRouterHandle(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedResponse := []byte("foobar")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rt := NewRouter()
	h := &handlerMockup{status: expectedStatus, response: expectedResponse}

	rt.Handle(http.MethodGet, "/", h)
	rt.ServeHTTP(w, r)

	body := w.Body.Bytes()

	if expectedStatus != w.Code {
		t.Errorf("%d is not expected status (%d)\n", w.Code, expectedStatus)
	}

	if !bytes.Equal(expectedResponse, body) {
		t.Errorf("%x is not expected response (%x)\n", body, expectedResponse)
	}

	if !h.finished {
		t.Error("http.Handler has not run")
	}
}

func TestRouterHandleFunc(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedResponse := []byte("foobar")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rt := NewRouter()
	h := &handlerMockup{status: expectedStatus, response: expectedResponse}

	rt.HandleFunc(http.MethodGet, "/", h.ServeHTTP)
	rt.ServeHTTP(w, r)

	body := w.Body.Bytes()

	if expectedStatus != w.Code {
		t.Errorf("%d is not expected status (%d)\n", w.Code, expectedStatus)
	}

	if !bytes.Equal(expectedResponse, body) {
		t.Errorf("%x is not expected response (%x)\n", body, expectedResponse)
	}

	if !h.finished {
		t.Error("http.HandlerFunc has not run")
	}
}

func TestRouterWithCancel(t *testing.T) {
	expectedStatus := http.StatusBadRequest
	expectedResponse := []byte{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rt := NewRouter()
	h := &handlerMockup{status: expectedStatus, response: expectedResponse}

	rt.HandleMiddlewares(http.MethodGet, "/", h.Cancel, h.ServeHTTP)
	rt.ServeHTTP(w, r)

	body := w.Body.Bytes()

	if expectedStatus != w.Code {
		t.Errorf("%d is not expected status (%d)\n", w.Code, expectedStatus)
	}

	if !bytes.Equal(expectedResponse, body) {
		t.Errorf("%x is not expected response (%x)\n", body, expectedResponse)
	}

	if !h.canceled {
		t.Error("http.HandlerFunc has not been canceled")
	}

	if h.finished {
		t.Error("http.HandlerFunc has run")
	}
}

func TestRouterHandleWithParams(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedResponse := map[string]string{"foo": "123", "bar": "456"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/123/456", nil)
	rt := NewRouter()
	h := &handlerMockup{status: expectedStatus, returnParams: true}

	rt.Handle(http.MethodGet, "/:foo/:bar", h)
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

	if !h.finished {
		t.Error("http.Handler has not run")
	}
}
