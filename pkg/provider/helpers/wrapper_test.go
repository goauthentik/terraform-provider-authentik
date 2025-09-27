package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_castSlice(t *testing.T) {
	foo := []interface{}{"test"}
	bar := CastSlice[string](foo)
	assert.Equal(t, bar, []string{"test"})
}
