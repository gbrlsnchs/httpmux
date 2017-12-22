package httpmux_test

import "net/http"

type handlerMockup struct {
	status   int
	response string
	run      bool
	err      error
}

func (hm *handlerMockup) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(hm.status)
	_, hm.err = w.Write([]byte(hm.response))
	hm.run = true
}
