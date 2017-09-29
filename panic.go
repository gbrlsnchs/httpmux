package httpmux

// PanicFunc is the function called when the Lookup method
// returns an error inside the ServeHTTP method.
//
// By default, the built-in panic function is called.
var PanicFunc func(error)

// callPanic calls a custom panic function if any is set.
// Otherwise, it calls the built-in panic function.
func callPanic(err error) {
	if PanicFunc != nil {
		PanicFunc(err)

		return
	}

	panic(err)
}
