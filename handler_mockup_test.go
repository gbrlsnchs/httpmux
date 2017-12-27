package httpmux_test

import (
	"net/http"

	. "github.com/gbrlsnchs/httpmux"
)

type handlerMockup struct {
	status   int
	response []byte
	canceled bool
	finished bool
	err      error
}

func (hm *handlerMockup) Cancel(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	defer Cancel(r)

	hm.canceled = true
}

func (hm *handlerMockup) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(hm.status)

	_, hm.err = w.Write(hm.response)
	hm.finished = true
}
