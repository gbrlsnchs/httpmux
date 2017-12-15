package httpmux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
)

func BenchmarkStatic(b *testing.B) {
	b.ReportAllocs()

	rt := NewRouter().WithPrefix("/foo/bar/baz/qux")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo/bar/baz/qux", nil)

	rt.HandleFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < b.N; i++ {
		rt.ServeHTTP(w, r)
	}
}

func BenchmarkDynamic(b *testing.B) {
	b.ReportAllocs()

	rt := NewRouter().WithPrefix("/foo/:bar/:baz/:qux")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo/123/456/789", nil)

	rt.HandleFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < b.N; i++ {
		rt.ServeHTTP(w, r)
	}
}

func BenchmarkStaticWithCancel(b *testing.B) {
	b.ReportAllocs()

	rt := NewRouter().WithPrefix("/foo/bar/baz/qux")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo/bar/baz/qux", nil)

	rt.HandleMiddlewares(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		Cancel(r)
	}, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < b.N; i++ {
		rt.ServeHTTP(w, r)
	}
}

func BenchmarkDynamicWithCancel(b *testing.B) {
	b.ReportAllocs()

	rt := NewRouter().WithPrefix("/foo/:bar/:baz/:qux")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo/123/456/789", nil)

	rt.HandleMiddlewares(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		Cancel(r)
	}, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < b.N; i++ {
		rt.ServeHTTP(w, r)
	}
}
