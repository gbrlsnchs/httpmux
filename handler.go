package httpmux

import (
	"net/http"
)

type Handler struct {
	h      http.Handler
	method string
}
