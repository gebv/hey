package utils

import (
	"fmt"
	"regexp"
)

const NameChars = `[a-zA-Z0-9][a-zA-Z0-9_-]`

var NamePattern = regexp.MustCompile(`^/?` + NameChars + `+$`)

func ValidName(name string) error {
	if !NamePattern.MatchString(name) {
		return fmt.Errorf("invalid name")
	}

	return nil
}

var NameSeparator = `.`
var SpecialNamePattern = regexp.
	MustCompile(
		"^(?P<channel>" +
			NameChars +
			"+)?" +
			NameSeparator + "?" +
			"(?P<thread>" +
			NameChars +
			"+)?$",
	)

type SpecialName string

func (s SpecialName) Valid() bool {
	return SpecialNamePattern.MatchString(string(s))
}

func (s SpecialName) Channel() string {
	if !s.Valid() {
		return ""
	}
	arr := SpecialNamePattern.FindStringSubmatch(string(s))

	return arr[1]
}

func (s SpecialName) Thread() string {
	if !s.Valid() {
		return ""
	}
	arr := SpecialNamePattern.FindStringSubmatch(string(s))

	return arr[2]
}
