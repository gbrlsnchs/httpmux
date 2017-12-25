package httpmux_test

import (
	"net/http"
	"testing"

	. "github.com/gbrlsnchs/httpmux"
)

func routerHelper(subrs ...*Subrouter) *Router {
	r := NewRouter()

	for _, subr := range subrs {
		r.Use(subr)
	}

	return r
}

func routerHelperWithHandler(status int, response string, hasHandler bool) *Router {
	rt := NewRouter()
	h := &handlerMockup{status: status, response: response}

	if hasHandler {
		rt.Handle(http.MethodGet, "/", h)

		return rt
	}

	rt.HandleFunc(http.MethodGet, "/", h.ServeHTTP)

	return rt
}

func TestRouterHandle(t *testing.T) {
	testRouter(t, true)
}

func TestRouterHandleFunc(t *testing.T) {
	testRouter(t, false)
}
