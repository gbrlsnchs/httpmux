package httpmux

import "strings"

// PanicFunc is the function called when the Lookup method
// returns an error inside the ServeHTTP method.
//
// By default, the built-in panic function is called.
var PanicFunc func(interface{})

func callPanic(m interface{}) {
	if PanicFunc != nil {
		PanicFunc(m)

		return
	}

	panic(m)
}

func dynPath(p string) bool {
	return strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}")
}

func extRegexp(p string) (r string, extp string) {
	if !dynPath(p) {
		return "", p
	}

	if i := strings.IndexByte(p, ':'); dynPath(p) && i >= 0 {
		return p[i+1 : len(p)-1], p[:i] + "}"
	}

	return "", p
}
