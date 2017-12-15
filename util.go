package httpmux

func resolvedPath(p string) string {
	if p == "/" {
		return ""
	}

	return p
}
