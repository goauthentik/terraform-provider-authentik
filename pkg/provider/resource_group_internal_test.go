package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// TestResourceGroupReadRolesPreserveConfiguredOrder is the executable spec for list
// ordering that MergeStringList must preserve. Unlike the SDKv2 version, this needs no
// httptest server: fromAPI is a plain function from (prior model, API response) to
// (model, diagnostics), so it can be called directly with a literal *api.Group.
func TestResourceGroupReadRolesPreserveConfiguredOrder(t *testing.T) {
	ctx := context.Background()
	r := &groupResource{}

	priorRoles, diags := types.ListValueFrom(ctx, types.StringType, []string{"role-a", "role-b", "role-c"})
	require.False(t, diags.HasError())

	data := &groupModel{
		Name:        types.StringValue("infrastructure"),
		IsSuperuser: types.BoolValue(false),
		Parents:     types.ListNull(types.StringType),
		Users:       types.ListNull(types.Int32Type),
		Roles:       priorRoles,
		Attributes:  jsontypes.NewNormalizedValue("{}"),
	}

	res := &api.Group{
		Pk:          "group-1",
		NumPk:       1,
		Name:        "infrastructure",
		IsSuperuser: new(false),
		Parents:     []string{},
		Users:       []int32{},
		Attributes:  map[string]any{},
		Roles:       []string{"role-c", "role-a", "role-b", "role-d"},
	}

	diags = r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	var roles []string
	require.False(t, data.Roles.ElementsAs(ctx, &roles, false).HasError())
	assert.Equal(t, []string{
		"role-a",
		"role-b",
		"role-c",
		"role-d",
	}, roles)
}

// TestResourceGroupToRequest_ListEncodings covers the write direction for this resource's
// three lists, which is where two distinct defects lived before helpers.SliceOrEmpty.
//
// users is Optional+Computed, so config omitting it makes the plan value *unknown* rather
// than null - and ElementsAs cannot write an unknown into a []int32, so a plain port
// failed Create outright with "target type cannot handle unknown values". parents and
// roles are plain Optional, so they arrive null, and a null list converted with
// ElementsAs yields a nil slice that GroupRequest.ToMap gates out with !IsNil() - which
// means clearing either one would leave the API's copy in place and then fail the apply
// on the plan/state mismatch. All three must serialise as a present, empty array, exactly
// as SDKv2's CastSlice did.
func TestResourceGroupToRequest_ListEncodings(t *testing.T) {
	ctx := context.Background()
	r := &groupResource{}

	data := &groupModel{
		Name:        types.StringValue("infrastructure"),
		IsSuperuser: types.BoolValue(false),
		Parents:     types.ListNull(types.StringType),
		Users:       types.ListUnknown(types.Int32Type),
		Roles:       types.ListNull(types.StringType),
		Attributes:  jsontypes.NewNormalizedValue("{}"),
	}

	body, diags := r.toRequest(ctx, data)
	require.False(t, diags.HasError(), "an unknown users list must not produce diagnostics: %v", diags)
	require.NotNil(t, body)

	require.NotNil(t, body.Parents, "a nil parents slice would be omitted, so clearing it would not clear it")
	assert.Empty(t, body.Parents)
	require.NotNil(t, body.Roles, "a nil roles slice would be omitted, so clearing it would not clear it")
	assert.Empty(t, body.Roles)
	require.NotNil(t, body.Users, "unknown users must serialise as [], matching SDKv2's CastSlice")
	assert.Empty(t, body.Users)

	// ToMap is what actually decides whether a field reaches the wire.
	serialised, err := body.ToMap()
	require.NoError(t, err)
	for _, key := range []string{"parents", "users", "roles"} {
		assert.Containsf(t, serialised, key, "%s must be present in the request body", key)
	}
}
