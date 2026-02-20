package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goauthentik.io/api/v3"
)

type TestResource map[string]any

func (tr TestResource) GetOk(key string) (any, bool) {
	v, ok := tr[key]
	return v, ok
}

func (tr TestResource) Set(k string, v any) error {
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

func Test_GetIntP_CastsInt(t *testing.T) {
	val := GetIntP(TestResource{
		"foo": 123,
	}, "foo")

	assert.NotNil(t, val)
	assert.Equal(t, int32(123), *val)
}

func Test_GetInt64P_CastsInt(t *testing.T) {
	val := GetInt64P(TestResource{
		"foo": 123,
	}, "foo")

	assert.NotNil(t, val)
	assert.Equal(t, int64(123), *val)
}

func Test_CastSlice_New(t *testing.T) {
	v := CastSlice[string](TestResource{
		"foo": []any{
			"foo",
			"bar",
		},
	}, "foo")
	assert.NotNil(t, v)
	assert.Equal(t, []string{"foo", "bar"}, v)
}
