package httpmux

// resolvedPath returns a resolved path
// to be used by Routers and Subrouters.
func resolvedPath(p string) string {
	if p == "/" {
		return ""
	}

	return p
}
