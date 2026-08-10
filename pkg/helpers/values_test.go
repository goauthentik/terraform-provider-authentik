package helpers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringOrNull(t *testing.T) {
	// prior null + API returns "" -> stays null, not ""
	assert.True(t, StringOrNull(types.StringNull(), "").IsNull())
	// prior non-null (e.g. explicit "" in config) + API returns "" -> stays ""
	assert.Equal(t, types.StringValue(""), StringOrNull(types.StringValue(""), ""))
	// API returns a real value regardless of prior
	assert.Equal(t, types.StringValue("foo"), StringOrNull(types.StringNull(), "foo"))
}

func TestStringPtrOrNull(t *testing.T) {
	assert.True(t, StringPtrOrNull(nil).IsNull())
	v := "foo"
	assert.Equal(t, types.StringValue("foo"), StringPtrOrNull(&v))
	empty := ""
	assert.Equal(t, types.StringValue(""), StringPtrOrNull(&empty))
}

func TestInt32PtrOrNull(t *testing.T) {
	assert.True(t, Int32PtrOrNull(nil).IsNull())
	v := int32(5)
	assert.Equal(t, types.Int32Value(5), Int32PtrOrNull(&v))
	zero := int32(0)
	assert.Equal(t, types.Int32Value(0), Int32PtrOrNull(&zero))
}

func TestInt32Ptr(t *testing.T) {
	assert.Nil(t, Int32Ptr(types.Int32Null()))
	v := Int32Ptr(types.Int32Value(5))
	require.NotNil(t, v)
	assert.Equal(t, int32(5), *v)
}

func TestInt32OrNull(t *testing.T) {
	assert.True(t, Int32OrNull(types.Int32Null(), 0).IsNull())
	assert.Equal(t, types.Int32Value(0), Int32OrNull(types.Int32Value(0), 0))
	assert.Equal(t, types.Int32Value(5), Int32OrNull(types.Int32Null(), 5))
}

func TestBoolOrNull(t *testing.T) {
	assert.True(t, BoolOrNull(types.BoolNull(), false).IsNull())
	assert.Equal(t, types.BoolValue(false), BoolOrNull(types.BoolValue(false), false))
	assert.Equal(t, types.BoolValue(true), BoolOrNull(types.BoolNull(), true))
}

func TestListOrNull(t *testing.T) {
	ctx := context.Background()

	l, diags := ListOrNull(ctx, types.ListNull(types.StringType), types.StringType, []string{})
	require.False(t, diags.HasError())
	assert.True(t, l.IsNull())

	prior, diags := types.ListValueFrom(ctx, types.StringType, []string{})
	require.False(t, diags.HasError())
	l, diags = ListOrNull(ctx, prior, types.StringType, []string{})
	require.False(t, diags.HasError())
	assert.False(t, l.IsNull())

	l, diags = ListOrNull(ctx, types.ListNull(types.StringType), types.StringType, []string{"a"})
	require.False(t, diags.HasError())
	assert.False(t, l.IsNull())
}

func TestStringPtr(t *testing.T) {
	assert.Nil(t, StringPtr(types.StringNull()))
	v := StringPtr(types.StringValue("foo"))
	require.NotNil(t, v)
	assert.Equal(t, "foo", *v)
}

func TestStringPtrEmpty(t *testing.T) {
	v := StringPtrEmpty(types.StringNull())
	require.NotNil(t, v)
	assert.Equal(t, "", *v)

	v = StringPtrEmpty(types.StringValue("foo"))
	require.NotNil(t, v)
	assert.Equal(t, "foo", *v)
}

// Test_BoolPtr pins the three-way distinction SDKv2's GetP[bool] could not express. It was
// d.GetOk underneath, so it returned nil for *any* false value - unset and explicitly-false
// were indistinguishable. Combined with omitempty on the request field and the API leaving
// absent fields unchanged on PUT, that meant a boolean could be switched on but never back
// off. authentik_policy_geoip's check_history_distance/check_impossible_travel are the case
// migration-plan.md calls out as user-visible.
func Test_BoolPtr(t *testing.T) {
	t.Run("null omits the field", func(t *testing.T) {
		assert.Nil(t, BoolPtr(types.BoolNull()))
	})

	t.Run("explicit false is sent, not omitted", func(t *testing.T) {
		v := BoolPtr(types.BoolValue(false))
		require.NotNil(t, v, "an explicit false must reach the API, or the field can never be turned off")
		assert.False(t, *v)
	})

	t.Run("true is sent", func(t *testing.T) {
		v := BoolPtr(types.BoolValue(true))
		require.NotNil(t, v)
		assert.True(t, *v)
	})
}

func Test_Int64Ptr(t *testing.T) {
	assert.Nil(t, Int64Ptr(types.Int64Null()))
	v := Int64Ptr(types.Int64Value(0))
	require.NotNil(t, v, "an explicit 0 must be sent, matching BoolPtr's handling of false")
	assert.Equal(t, int64(0), *v)
}

func Test_Int64OrNull(t *testing.T) {
	assert.True(t, Int64OrNull(types.Int64Null(), 0).IsNull())
	assert.Equal(t, int64(0), Int64OrNull(types.Int64Value(0), 0).ValueInt64())
	assert.Equal(t, int64(100), Int64OrNull(types.Int64Null(), 100).ValueInt64())
}
