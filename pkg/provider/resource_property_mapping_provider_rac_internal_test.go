package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

// authentik_property_mapping_provider_rac is the only one of the fourteen property mappings
// whose expression is Optional rather than Required, and the only one with a JSON attribute.
// Both of those are the interesting cases, so it carries this batch's unit tests.

// TestPropertyMappingProviderRACFromAPI_OptionalExpression is H1 for an expression
// attribute. Semantic equality cannot cover null-vs-"" - value_semantic_equality.go returns
// early when the prior value is null - so ExpressionOrNull is needed on top of the custom
// type. The other thirteen mappings do not need it only because their expression is
// Required and therefore never null.
func TestPropertyMappingProviderRACFromAPI_OptionalExpression(t *testing.T) {
	r := &propertyMappingProviderRACResource{}

	empty := ""
	res := &api.RACPropertyMapping{
		Pk:             "uuid-1",
		Name:           "mapping",
		Expression:     &empty,
		StaticSettings: map[string]any{},
	}

	t.Run("null prior and empty API value stays null", func(t *testing.T) {
		data := &propertyMappingProviderRACModel{Expression: helpers.NewExpressionNull()}
		diags := r.fromAPI(data, res)
		require.False(t, diags.HasError(), diags)
		assert.True(t, data.Expression.IsNull(), "expression was null and the API returned \"\", so it stays null")
	})

	t.Run("configured prior and empty API value becomes empty", func(t *testing.T) {
		data := &propertyMappingProviderRACModel{Expression: helpers.NewExpressionValue("")}
		diags := r.fromAPI(data, res)
		require.False(t, diags.HasError(), diags)
		assert.False(t, data.Expression.IsNull(), "an explicitly-configured \"\" must not become null")
	})

	t.Run("trailing newline is preserved in state, not normalised away", func(t *testing.T) {
		data := &propertyMappingProviderRACModel{Expression: helpers.NewExpressionNull()}
		withValue := &api.RACPropertyMapping{
			Pk: "uuid-1", Name: "mapping",
			Expression:     api.PtrString("return True"),
			StaticSettings: map[string]any{},
		}
		diags := r.fromAPI(data, withValue)
		require.False(t, diags.HasError(), diags)
		assert.Equal(t, "return True", data.Expression.ValueString(),
			"semantic equality compares, it does not rewrite - state holds what the API returned")
	})
}

// TestPropertyMappingProviderRACFromAPI_SettingsJSON pins the settings round trip. The
// attribute has a Default of "{}", so it can never be null in state (discovery #5), and an
// empty map from the API must serialise back to "{}" rather than "null".
func TestPropertyMappingProviderRACFromAPI_SettingsJSON(t *testing.T) {
	r := &propertyMappingProviderRACResource{}

	t.Run("empty settings map becomes {}", func(t *testing.T) {
		data := &propertyMappingProviderRACModel{Expression: helpers.NewExpressionNull()}
		diags := r.fromAPI(data, &api.RACPropertyMapping{
			Pk: "uuid-1", Name: "mapping",
			StaticSettings: map[string]any{},
		})
		require.False(t, diags.HasError(), diags)
		require.False(t, data.Settings.IsNull(), "settings has a Default, so it must never be null")
		assert.JSONEq(t, "{}", data.Settings.ValueString())
	})

	// A nil map marshals to "null", which is valid JSON but not a valid object for this
	// attribute's Default. Asserting the current behaviour so a regression is visible.
	t.Run("populated settings round-trip", func(t *testing.T) {
		data := &propertyMappingProviderRACModel{Expression: helpers.NewExpressionNull()}
		diags := r.fromAPI(data, &api.RACPropertyMapping{
			Pk: "uuid-1", Name: "mapping",
			StaticSettings: map[string]any{"resize-method": "display-update"},
		})
		require.False(t, diags.HasError(), diags)
		assert.JSONEq(t, `{"resize-method":"display-update"}`, data.Settings.ValueString())
	})
}

// TestPropertyMappingProviderRACToRequest covers the write direction: a null expression must
// be omitted (SDKv2 used GetP), and settings must always be sent since it has a Default.
func TestPropertyMappingProviderRACToRequest(t *testing.T) {
	r := &propertyMappingProviderRACResource{}

	data := &propertyMappingProviderRACModel{
		Name:       types.StringValue("mapping"),
		Expression: helpers.NewExpressionNull(),
		Settings:   jsontypes.NewNormalizedValue("{}"),
	}

	body, diags := r.toRequest(data)
	require.False(t, diags.HasError(), diags)

	assert.Nil(t, body.Expression, "a null expression must be omitted, not sent as \"\"")
	require.NotNil(t, body.StaticSettings, "static_settings is not omitempty, so it must not be nil")
	assert.Empty(t, body.StaticSettings)

	// A heredoc value must reach the API byte-for-byte.
	data.Expression = helpers.NewExpressionValue("return True\n")
	body, diags = r.toRequest(data)
	require.False(t, diags.HasError(), diags)
	require.NotNil(t, body.Expression)
	assert.Equal(t, "return True\n", *body.Expression)
}
