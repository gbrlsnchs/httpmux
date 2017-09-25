package httpmux

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynPath(t *testing.T) {
	a := assert.New(t)
	tests := []*struct {
		txt      string
		expected bool
	}{
		// #1
		{
			txt:      "test",
			expected: false,
		},
		// #2
		{
			txt:      "{test}",
			expected: true,
		},
		// #3
		{
			txt:      "{test",
			expected: false,
		},
		// #4
		{
			txt:      "test}",
			expected: false,
		},
		// #5
		{
			txt:      "[test]",
			expected: false,
		},
	}

	for i, test := range tests {
		index := strconv.Itoa(i)

		a.Exactly(test.expected, dynPath(test.txt), index)
	}
}

func TestExtRegexp(t *testing.T) {
	a := assert.New(t)
	tests := []*struct {
		txt            string
		expectedTxt    string
		expectedRegexp string
	}{
		// #1
		{
			txt:            "test:[0-9]+",
			expectedTxt:    "test:[0-9]+",
			expectedRegexp: "",
		},
		// #2
		{
			txt:            "{test:[0-9]+}",
			expectedTxt:    "{test}",
			expectedRegexp: "[0-9]+",
		},
	}

	for i, test := range tests {
		index := strconv.Itoa(i)
		r, p := extRegexp(test.txt)

		a.Exactly(test.expectedTxt, p, index)
		a.Exactly(test.expectedRegexp, r, index)
	}
}
