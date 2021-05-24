package domain

import (
	"errors"
	"strings"
	"unicode"
)

var validateName = func(name string) {
	if len(strings.TrimSpace(name)) == 0 {
		panic(errors.New("name must not be empty"))
	}
	if unicode.IsSpace(rune(name[0])) || unicode.IsSpace(rune(name[len(name)-1])) {
		panic(errors.New("name must not contain leading or trailing white space"))
	}
}
