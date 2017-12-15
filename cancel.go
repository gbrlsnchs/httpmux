package httpmux

import (
	"context"
	"net/http"
)

func Cancel(r *http.Request) {
	_, cancel := context.WithCancel(r.Context())

	cancel()
}
