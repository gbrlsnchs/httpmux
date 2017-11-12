package httpmux

import "strings"

func resolvePath(p string) string {
	if p[0] != '/' {
		p = "/" + p
	}

	return strings.ToLower(p)
}
