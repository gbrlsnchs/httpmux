package httpmux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
)

func BenchmarkSingleStatic(b *testing.B) {
	b.ReportAllocs()

	mux := New("/foo/bar/baz/qux")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo/bar/baz/qux", nil)

	mux.HandleFunc(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, r)
	}
}

func BenchmarkSingleDynamic(b *testing.B) {
	b.ReportAllocs()

	mux := New("/foo/:bar/:baz/:qux")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo/123/456/789", nil)

	mux.HandleFunc(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, r)
	}
}

func BenchmarkMultipleStatic(b *testing.B) {
	b.ReportAllocs()

	qux := NewSubmux("/qux")
	baz := NewSubmux("/baz")
	bar := NewSubmux("/bar")
	mux := New("/foo")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo/bar/baz/qux", nil)

	qux.HandleFunc(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	baz.Add(qux)
	bar.Add(baz)
	mux.Add(bar)

	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, r)
	}
}

func BenchmarkMultipleDynamic(b *testing.B) {
	b.ReportAllocs()

	qux := NewSubmux("/:qux")
	baz := NewSubmux("/:baz")
	bar := NewSubmux("/:bar")
	mux := New("/foo")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo/123/456/789", nil)

	qux.HandleFunc(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	baz.Add(qux)
	bar.Add(baz)
	mux.Add(bar)

	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, r)
	}
}
