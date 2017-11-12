package httpmux_test

import (
	"net/http"

	"github.com/gbrlsnchs/httpmux"
)

func Example() {
	example := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// "/api/auth"
	authMux := httpmux.NewSubmux("/auth")

	authMux.HandleFunc(http.MethodPost, example)

	// "/api/user/:id"
	userIDMux := httpmux.NewSubmux("/:id")

	userIDMux.HandleFunc(http.MethodGet, example)
	userIDMux.HandleFunc(http.MethodPost, example)

	// "/api/user"
	userMux := httpmux.NewSubmux("/user")

	userMux.HandleFunc(http.MethodPost, example)
	userMux.Add(userIDMux)

	// "/api"
	parentMux := httpmux.New("/api")

	parentMux.Add(authMux)
	parentMux.Add(userMux)

	http.ListenAndServe("/", parentMux)
}
