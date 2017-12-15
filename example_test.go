package httpmux_test

import (
	"log"
	"net/http"

	"github.com/gbrlsnchs/httpmux"
)

func Example() {
	rt := httpmux.NewRouter()

	rt.HandleMiddlewares(http.MethodGet, "/example",
		// Logger.
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("r.URL.Path = %s\n", r.URL.Path)
		},
		// Guard.
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/forbidden" {
				w.WriteHeader(http.StatusForbidden)
				httpmux.Cancel(r)
			}
		},
		// Handler.
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	http.ListenAndServe("/", rt)
}
