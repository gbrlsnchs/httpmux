package httpmux

import (
	"context"
	"net/http"
)

// Cancel sinalizes to not run the next
// handler in a middleware stack.
func Cancel(r *http.Request) {
	_, cancel := context.WithCancel(r.Context())

	cancel()
}
