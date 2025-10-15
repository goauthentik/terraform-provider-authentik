package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goauthentik.io/api/v3"
)

func Test_CastSlice(t *testing.T) {
	v := CastSlice[string](TestResource{
		"foo": []interface{}{"bar", "baz"},
	}, "foo")
	assert.NotNil(t, v)
	assert.Equal(t, []string{"bar", "baz"}, v)
}

func Test_CastSliceString(t *testing.T) {
	v := CastSliceString[api.IntentEnum]([]string{
		string(api.INTENTENUM_API),
	})
	assert.Equal(t, []api.IntentEnum{api.INTENTENUM_API}, v)
}
