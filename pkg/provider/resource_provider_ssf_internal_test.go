package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// TestProviderSSFFromAPI_PreservesUnrefreshedAttributes is the widest H4 case in the
// provider: SDKv2's Read wrote only name and push_verify_certificates, leaving three of the
// five attributes untouched. fromAPI reproduces that by not assigning them, which only works
// because callers seed the model from req.Plan or req.State first.
//
// event_retention is the interesting one - it has a schema Default, so discovery #5 would
// normally say "map verbatim", but H4 wins where the two collide: the attribute is not
// refreshed at all. The API response below deliberately carries *different* values for all
// three to prove none of them leaks into state.
func TestProviderSSFFromAPI_PreservesUnrefreshedAttributes(t *testing.T) {
	r := &providerSSFResource{}

	configured, diags := types.ListValueFrom(t.Context(), types.Int32Type, []int32{7, 3})
	require.False(t, diags.HasError(), diags)

	data := &providerSSFModel{
		SigningKey:             types.StringValue("configured-key"),
		JWTFederationProviders: configured,
		EventRetention:         types.StringValue("days=30"),
	}

	res := &api.SSFProvider{
		Pk:                     11,
		Name:                   "ssf",
		SigningKey:             "server-key",
		EventRetention:         api.PtrString("days=1"),
		OidcAuthProviders:      []int32{99},
		PushVerifyCertificates: new(false),
	}

	r.fromAPI(data, res)

	assert.Equal(t, "11", data.ID.ValueString(), "the int32 PK is stringified into state")
	assert.Equal(t, "ssf", data.Name.ValueString())
	assert.False(t, data.PushVerifyCertificates.ValueBool(), "push_verify_certificates *is* refreshed")

	assert.Equal(t, "configured-key", data.SigningKey.ValueString(),
		"signing_key is never refreshed from the server (H4)")
	assert.Equal(t, "days=30", data.EventRetention.ValueString(),
		"event_retention is never refreshed either, despite having a Default")

	var got []int32
	require.False(t, data.JWTFederationProviders.ElementsAs(t.Context(), &got, false).HasError())
	assert.Equal(t, []int32{7, 3}, got, "jwt_federation_providers keeps the configured value and order")
}

// TestProviderSSFToRequest_MapsJWTFederationProvidersToOIDCField pins the attribute-to-wire
// rename, which is easy to get wrong because the names differ.
func TestProviderSSFToRequest_MapsJWTFederationProvidersToOIDCField(t *testing.T) {
	r := &providerSSFResource{}

	providers, diags := types.ListValueFrom(t.Context(), types.Int32Type, []int32{4, 1})
	require.False(t, diags.HasError(), diags)

	body, diags := r.toRequest(t.Context(), &providerSSFModel{
		Name:                   types.StringValue("ssf"),
		SigningKey:             types.StringValue("key"),
		JWTFederationProviders: providers,
		EventRetention:         types.StringValue("days=30"),
		PushVerifyCertificates: types.BoolValue(true),
	})
	require.False(t, diags.HasError(), diags)

	assert.Equal(t, []int32{4, 1}, body.OidcAuthProviders,
		"jwt_federation_providers maps to the wire's oidc_auth_providers, in order")

	// discovery #6: a null list must still serialise as present-but-empty.
	body, diags = r.toRequest(t.Context(), &providerSSFModel{
		Name:                   types.StringValue("ssf"),
		SigningKey:             types.StringValue("key"),
		JWTFederationProviders: types.ListNull(types.Int32Type),
		EventRetention:         types.StringValue("days=30"),
		PushVerifyCertificates: types.BoolValue(true),
	})
	require.False(t, diags.HasError(), diags)
	require.NotNil(t, body.OidcAuthProviders)
	assert.Empty(t, body.OidcAuthProviders)
}
