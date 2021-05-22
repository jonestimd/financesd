package domain

import (
	"errors"
	"strings"
	"unicode"
)

var validateName = func(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return errors.New("name must not be empty")
	}
	if unicode.IsSpace(rune(name[0])) || unicode.IsSpace(rune(name[len(name)-1])) {
		return errors.New("name must not contain leading or trailing white space")
	}
	return nil
}
