package httpmux_test

import (
	"log"
	"net/http"

	"github.com/gbrlsnchs/httpmux"
	"github.com/gbrlsnchs/httpmux/internal"
)

func Example() {
	rt := httpmux.NewRouter()

	rt.SetParamsKey(internal.ParamsKey)
	rt.HandleMiddlewares(http.MethodGet, "/:path",
		// Logger.
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("r.URL.Path = %s\n", r.URL.Path)
		},
		// Guard.
		func(w http.ResponseWriter, r *http.Request) {
			if params, ok := r.Context().Value(internal.ParamsKey).(map[string]string); ok {
				if params["path"] == "forbidden" {
					w.WriteHeader(http.StatusForbidden)
					httpmux.Cancel(r)
				}

				return
			}

			httpmux.Cancel(r)
		},
		// Handler.
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	http.ListenAndServe("/", rt)
}
