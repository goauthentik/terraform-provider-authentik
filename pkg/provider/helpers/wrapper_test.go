package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goauthentik.io/api/v3"
)

func Test_CastSlice(t *testing.T) {
	foo := []interface{}{"test"}
	bar := CastSlice[string](foo)
	assert.Equal(t, bar, []string{"test"})
}

type TestResource map[string]interface{}

func (tr TestResource) GetOk(key string) (interface{}, bool) {
	v, ok := tr[key]
	return v, ok
}

func (tr TestResource) Set(k string, v interface{}) error {
	tr[k] = v
	return nil
}

func Test_GetP_Enum(t *testing.T) {
	v := CastString[api.NotConfiguredActionEnum](GetP[string](TestResource{
		"foo": "skip",
	}, "foo"))
	assert.NotNil(t, v)
	assert.Equal(t, api.NOTCONFIGUREDACTIONENUM_SKIP, *v)
}
