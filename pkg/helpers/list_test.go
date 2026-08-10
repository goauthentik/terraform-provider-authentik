package helpers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestMergeStringList(t *testing.T) {
	ctx := context.Background()

	// prior null + API returns nothing -> stays null
	l, diags := MergeStringList(ctx, types.ListNull(types.StringType), []string{})
	require.False(t, diags.HasError())
	assert.True(t, l.IsNull())

	// prior order is preserved, new entries appended - same guarantee as
	// ListConsistentMerge itself, just through the types.List round trip
	prior, diags := types.ListValueFrom(ctx, types.StringType, []string{"foo", "bar", "baz"})
	require.False(t, diags.HasError())
	l, diags = MergeStringList(ctx, prior, []string{"baz", "foo", "bar", "quox"})
	require.False(t, diags.HasError())
	var got []string
	require.False(t, l.ElementsAs(ctx, &got, false).HasError())
	assert.Equal(t, []string{"foo", "bar", "baz", "quox"}, got)
}

func TestMergeInt32List(t *testing.T) {
	ctx := context.Background()

	l, diags := MergeInt32List(ctx, types.ListNull(types.Int32Type), []int32{})
	require.False(t, diags.HasError())
	assert.True(t, l.IsNull())

	l, diags = MergeInt32List(ctx, types.ListNull(types.Int32Type), []int32{1, 2})
	require.False(t, diags.HasError())
	var got []int32
	require.False(t, l.ElementsAs(ctx, &got, false).HasError())
	assert.Equal(t, []int32{1, 2}, got)
}
