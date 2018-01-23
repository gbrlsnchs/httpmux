package httpmux_test

import (
	"net/http"
	"strings"

	. "github.com/gbrlsnchs/httpmux"
)

type handlerMockup struct {
	status       int
	response     []byte
	canceled     bool
	finished     bool
	err          error
	returnParams bool
}

func (hm *handlerMockup) Cancel(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	Cancel(r)

	hm.canceled = true
}

func (hm *handlerMockup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(hm.status)

	content := hm.response

	if hm.returnParams {
		p := Params(r)
		keys := []string{}

		for k, v := range p {
			keys = append(keys, strings.Join([]string{k, v}, "="))
		}

		content = []byte(strings.Join(keys, "/"))
	}

	_, hm.err = w.Write(content)
	hm.finished = true
}
