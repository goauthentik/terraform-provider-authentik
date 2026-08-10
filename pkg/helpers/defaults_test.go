package helpers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringDefault(t *testing.T) {
	d := StringDefault("foo")
	assert.Equal(t, "foo", d.Value())

	// still satisfies defaults.String directly, for schema.StringAttribute.Default
	var _ defaults.String = d
	resp := &defaults.StringResponse{}
	d.DefaultString(context.Background(), defaults.StringRequest{}, resp)
	require.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "foo", resp.PlanValue.ValueString())
}

func TestBoolDefault(t *testing.T) {
	d := BoolDefault(true)
	assert.True(t, d.Value())

	resp := &defaults.BoolResponse{}
	d.DefaultBool(context.Background(), defaults.BoolRequest{}, resp)
	require.False(t, resp.Diagnostics.HasError())
	assert.True(t, resp.PlanValue.ValueBool())
}

func TestInt32Default(t *testing.T) {
	d := Int32Default(5)
	assert.Equal(t, int32(5), d.Value())

	resp := &defaults.Int32Response{}
	d.DefaultInt32(context.Background(), defaults.Int32Request{}, resp)
	require.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, int32(5), resp.PlanValue.ValueInt32())
}

func TestFloat64Default(t *testing.T) {
	d := Float64Default(0.5)
	assert.InDelta(t, 0.5, d.Value(), 0.0001)

	resp := &defaults.Float64Response{}
	d.DefaultFloat64(context.Background(), defaults.Float64Request{}, resp)
	require.False(t, resp.Diagnostics.HasError())
	assert.InDelta(t, 0.5, resp.PlanValue.ValueFloat64(), 0.0001)
}
