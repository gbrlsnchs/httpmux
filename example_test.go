package httpmux_test

import (
	"net/http"

	"github.com/gbrlsnchs/httphandler"
	"github.com/gbrlsnchs/httpmux"
)

func Example() {
	// submuxes
	authFunc := func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		id := r.Context().Value("id")

		return id, nil
	}
	authMux := httpmux.New("/auth").
		SetHandler(httphandler.New(http.StatusOK, authFunc), http.MethodPost).
		SetSubmux(httpmux.New("{id:[0-9]+}"))

	userFunc := func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		id := r.Context().Value("id")

		return id, nil
	}
	userMux := httpmux.New("/user").
		SetSubmux(httpmux.New("/{id}"))

	userMux.Submux("/{id}").
		SetHandler(httphandler.New(http.StatusAccepted, userFunc), http.MethodPost, http.MethodPut)

	testFunc := func(w http.ResponseWriter, r *http.Request) (interface{}, error) { return nil, nil }
	testMux := httpmux.New("/test").
		SetHandler(httphandler.New(http.StatusCreated, testFunc))

	// main mux
	m := httpmux.New("/api").
		SetSubmux(authMux).
		SetSubmux(userMux).
		SetSubmux(testMux)

	http.Handle("/", m)
}
