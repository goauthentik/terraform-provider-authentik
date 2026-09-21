package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffSuppressExpression(t *testing.T) {
	cases := []struct {
		name   string
		old    string
		new    string
		expect bool
	}{
		{"identical", "a {}", "a {}", true},
		{"new adds trailing newline", "a {}", "a {}\n", true},
		{"new adds two trailing newlines", "a {}", "a {}\n\n", false},
		{"changed content", "a {}", "b {}\n", false},
		{"removed trailing newline", "a {}\n", "a {}", false},
		{"empty vs newline", "", "\n", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, DiffSuppressExpression("", tc.old, tc.new, nil))
		})
	}
}
