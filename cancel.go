package httpmux

import (
	"context"
	"net/http"
)

// Cancel sinalizes to not run the next
// handler in a middleware stack.
func Cancel(r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	*r = *r.WithContext(ctx)

	cancel()
}
