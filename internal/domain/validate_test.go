package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateName(t *testing.T) {
	tests := []struct {
		name  string
		value string
		err   error
	}{
		{"panics for empty name", "", errors.New("name must not be empty")},
		{"panics for leading whitespace", " x", errors.New("name must not contain leading or trailing white space")},
		{"panics for trailing whitespace", "x ", errors.New("name must not contain leading or trailing white space")},
		{"returns for valid name", "x", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, test.err, err)
				} else if test.err != nil {
					assert.Fail(t, "expected an error")
				}
			}()

			validateName(test.value)
		})
	}
}
