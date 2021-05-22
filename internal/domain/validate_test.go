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
		{"returns error for empty name", "", errors.New("name must not be empty")},
		{"returns error leading whitespace", " x", errors.New("name must not contain leading or trailing white space")},
		{"returns error trailing whitespace", "x ", errors.New("name must not contain leading or trailing white space")},
		{"returns nil for valid name", "x", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.err, validateName(test.value))
		})
	}
}
