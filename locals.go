package httpmux

import (
	"context"
	"net/http"
)

// SetLocal sets a variable in a request's context
// for retrieving it in another middleware.
func SetLocal(r *http.Request, key, val interface{}) {
	ctx := context.WithValue(r.Context(), key, val)
	*r = *r.WithContext(ctx)
}
