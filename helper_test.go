package httpmux_test

import (
	"bytes"
	"strings"
	"testing"
)

func testCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("%d != %d\n", expected, actual)
	}
}

func testResponse(t *testing.T, expected []byte, actual []byte) {
	if !bytes.Equal(expected, actual) {
		t.Errorf("%s != %s\n",
			strings.TrimSuffix(string(expected), "\n"),
			strings.TrimSuffix(string(actual), "\n"),
		)
	}
}

func testHTTPResponse(t *testing.T, expectedCode, actualCode int, expectedResponse, actualResponse []byte) {
	testCode(t, expectedCode, actualCode)
	testResponse(t, expectedResponse, actualResponse)
}
