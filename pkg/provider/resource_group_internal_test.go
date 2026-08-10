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
