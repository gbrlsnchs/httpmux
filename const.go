package httpmux

// CtxKey is the custom type for the key value
// inside a request's context.
type CtxKey uint

const (
	// ParamsKey is the value a request's context
	// holds if a route with named parameters is match.
	ParamsKey CtxKey = iota
)
