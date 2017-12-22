package httpmux_test

import "net/http"

type handlerMock struct {
	status   int
	response string
}

func (hm *handlerMock) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(hm.status)
	w.Write([]byte(hm.response))
}
