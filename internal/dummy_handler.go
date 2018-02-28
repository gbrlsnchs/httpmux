package internal

import (
	"net/http"
	"strings"

	"github.com/gbrlsnchs/httpmux"
)

type DummyHandler struct {
	Status       int
	Response     []byte
	Canceled     bool
	Finished     bool
	Err          error
	ReturnParams bool
}

func (d *DummyHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	httpmux.Cancel(r)

	d.Canceled = true
}

func (d *DummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(d.Status)

	content := d.Response

	if d.ReturnParams {
		p, _ := r.Context().Value(ParamsKey).(map[string]string)
		keys := []string{}

		for k, v := range p {
			keys = append(keys, k, v)
		}

		content = []byte(strings.Join(keys, "/"))
	}

	_, d.Err = w.Write(content)
	d.Finished = true
}
