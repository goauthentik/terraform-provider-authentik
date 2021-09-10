package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sliceToString(t *testing.T) {
	foo := []interface{}{"test"}
	bar := sliceToString(foo)
	assert.Equal(t, bar, []string{"test"})
}

func Test_sliceToStringPointer(t *testing.T) {
	foo := []interface{}{"test"}
	bar := sliceToStringPointer(foo)
	assert.Equal(t, bar, []*string{stringToPointer("test")})
}
