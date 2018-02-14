package httpmux

import (
	"context"
	"net/http"
)

// SetLocal sets a variable in a request's context
// for retrieving it in another middleware.
func SetLocal(r *http.Request, v, k interface{}) {
	*r = *r.WithContext(context.WithValue(r.Context(), k, v))
}
