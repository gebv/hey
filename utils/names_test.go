package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNames_simple_validName(t *testing.T) {
	type casetest struct {
		Name  string
		Valid bool
	}

	var cases = []casetest{
		{"", false},
		{"_", false},
		{"0", false},
		{"09", true},
		{"0.", false},
		{"0?", false},
		{":qwd", false},
		{"qwd", true},
		{"qwdqwd.qwd", false},
	}

	for _, _case := range cases {
		assert.Equal(t, ValidName(_case.Name) == nil, _case.Valid)
	}
}

func TestNames_simple_specialNames(t *testing.T) {
	type casetest struct {
		Name    string
		Valid   bool
		Channel string
		Thread  string
	}
	var cases = []casetest{
		{"ab.cd", true, "ab", "cd"},
		{"ab", true, "ab", ""},
		{".cd", true, "", "cd"},
		{".b.", false, "", ""},
		{"a.b", false, "", ""},
	}

	for _, _case := range cases {
		sn := SpecialName(_case.Name)
		assert.Equal(t, sn.Valid(), _case.Valid)
		assert.Equal(t, sn.Channel(), _case.Channel)
		assert.Equal(t, sn.Thread(), _case.Thread)
	}
}
