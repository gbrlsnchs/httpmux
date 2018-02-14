package httpmux

import "net/http"

var paramsKey interface{}

// Params returns a map with parameters extracted from the request.
func Params(r *http.Request) map[string]string {
	if paramsKey == nil {
		return nil
	}

	if p, ok := r.Context().Value(paramsKey).(map[string]string); ok {
		return p
	}

	return nil
}

// SetParamsKey sets a key for retrieving
// URL parameters only once, for avoiding
// unwanted manipulation.
func SetParamsKey(v interface{}) {
	if paramsKey == nil {
		paramsKey = v
	}
}
