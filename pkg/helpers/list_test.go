package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goauthentik.io/api/v3"
)

func Test_CastSlice(t *testing.T) {
	v := CastSlice[string](TestResource{
		"foo": []any{"bar", "baz"},
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

func Test_CastSliceInt32(t *testing.T) {
	v := CastSliceInt32([]any{
		1, 2, 3,
	})
	assert.Equal(t, []int32{1, 2, 3}, v)
}

func Test_Slice32ToInt(t *testing.T) {
	v := Slice32ToInt([]int32{
		1, 2, 3,
	})
	assert.Equal(t, []int{1, 2, 3}, v)
}

func Test_ListConsistentMerge(t *testing.T) {
	v := ListConsistentMerge(
		[]string{"foo", "bar", "baz"},
		[]string{"baz", "foo", "bar", "quox"},
	)
	assert.Equal(t, []string{
		"foo", "bar", "baz", "quox",
	}, v)
}
