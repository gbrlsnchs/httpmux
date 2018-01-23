package httpmux

import "net/http"

// Params returns a map with parameters extracted from the request.
func Params(r *http.Request) map[string]string {
	if p, ok := r.Context().Value(ParamsKey).(map[string]string); ok {
		return p
	}

	return nil
}
