package httpmux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
)

const longPath = "/test/test/test/test/test/test"

var handlerFunc1 = func(w http.ResponseWriter, r *http.Request) {}
var handlerFunc2 = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func BenchmarkHTTPMux(b *testing.B) {
	m := New(longPath+"/:test").SetHandler(http.MethodGet, &testHandler{handlerFunc1})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, longPath+"/123", nil)

	for i := 0; i < b.N; i++ {
		m.Root().ServeHTTP(w, r)
	}
}

func BenchmarkHTTPRouter(b *testing.B) {
	hr := httprouter.New()
	hr.GET(longPath+"/:test", handlerFunc2)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, longPath+"/123", nil)

	for i := 0; i < b.N; i++ {
		hr.ServeHTTP(w, r)
	}
}

func BenchmarkGorillaMux(b *testing.B) {
	gm := mux.NewRouter()

	gm.HandleFunc(longPath+"/{test}", handlerFunc1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, longPath+"/123", nil)

	for i := 0; i < b.N; i++ {
		gm.ServeHTTP(w, r)
	}
}
