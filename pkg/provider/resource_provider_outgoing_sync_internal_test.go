package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// TestProviderOutgoingSyncFromAPI_FilterGroupUnwrapped is the regression test for a provider
// *panic* found while migrating these two resources.
//
// SDKv2's Read did, in both files:
//
//	helpers.SetWrapper(d, "filter_group", res.FilterGroup)
//
// FilterGroup is an api.NullableString, i.e. a struct - the code passed the wrapper itself
// rather than calling .Get() on it. d.Set rejects that with
//
//	filter_group: '' expected type 'string', got unconvertible type 'api.NullableString'
//
// and helpers.SetWrapper panics on any error from d.Set. So *every* Read of an
// authentik_provider_google_workspace or authentik_provider_microsoft_entra resource crashed
// the provider - on refresh, on plan, on import. Neither resource has an acceptance test,
// which is why a guaranteed crash went unnoticed.
//
// The framework port unwraps properly via helpers.StringPtrOrNull, so the bug is fixed by
// construction; this test exists so nobody reintroduces the wrapper-passing shape.
func TestProviderOutgoingSyncFromAPI_FilterGroupUnwrapped(t *testing.T) {
	t.Run("google_workspace", func(t *testing.T) {
		r := &providerGoogleWorkspaceResource{}
		data := &providerGoogleWorkspaceModel{
			PropertyMappings:      types.ListNull(types.StringType),
			PropertyMappingsGroup: types.ListNull(types.StringType),
		}

		diags := r.fromAPI(t.Context(), data, &api.GoogleWorkspaceProvider{
			Pk:                      3,
			Name:                    "gws",
			DelegatedSubject:        "admin@example.com",
			DefaultGroupEmailDomain: "example.com",
			FilterGroup:             *api.NewNullableString(api.PtrString("group-uuid")),
		})
		require.False(t, diags.HasError(), diags)

		assert.Equal(t, "group-uuid", data.FilterGroup.ValueString(),
			"filter_group must be unwrapped with .Get(), not passed as the NullableString struct")
		assert.Equal(t, "3", data.ID.ValueString())
	})

	t.Run("google_workspace null filter_group stays null", func(t *testing.T) {
		r := &providerGoogleWorkspaceResource{}
		data := &providerGoogleWorkspaceModel{
			PropertyMappings:      types.ListNull(types.StringType),
			PropertyMappingsGroup: types.ListNull(types.StringType),
		}

		diags := r.fromAPI(t.Context(), data, &api.GoogleWorkspaceProvider{
			Pk:          3,
			Name:        "gws",
			FilterGroup: *api.NewNullableString(nil),
		})
		require.False(t, diags.HasError(), diags)
		assert.True(t, data.FilterGroup.IsNull())
	})

	t.Run("microsoft_entra", func(t *testing.T) {
		r := &providerMicrosoftEntraResource{}
		data := &providerMicrosoftEntraModel{
			PropertyMappings:      types.ListNull(types.StringType),
			PropertyMappingsGroup: types.ListNull(types.StringType),
		}

		diags := r.fromAPI(t.Context(), data, &api.MicrosoftEntraProvider{
			Pk:           4,
			Name:         "entra",
			ClientId:     "client-id",
			ClientSecret: "client-secret",
			TenantId:     "tenant-id",
			FilterGroup:  *api.NewNullableString(api.PtrString("group-uuid")),
		})
		require.False(t, diags.HasError(), diags)

		assert.Equal(t, "group-uuid", data.FilterGroup.ValueString())
		// client_secret is Sensitive but *is* returned, so unlike the H4 secrets it is
		// refreshed from the server rather than preserved from the plan.
		assert.Equal(t, "client-secret", data.ClientSecret.ValueString())
		assert.Equal(t, "4", data.ID.ValueString())
	})
}

// TestProviderWSFederationToRequest_SignLogoutRequestCanBeTurnedOff is the sixth instance of
// the GetP[bool] defect. Asserted on the pointer and the value; the marshalled-body form of
// this check lives in the geoip and telegram tests.
func TestProviderWSFederationToRequest_SignLogoutRequestCanBeTurnedOff(t *testing.T) {
	r := &providerWSFederationResource{}

	body, diags := r.toRequest(t.Context(), &providerWSFederationModel{
		Name:                       types.StringValue("wsfed"),
		AuthorizationFlow:          types.StringValue("flow"),
		InvalidationFlow:           types.StringValue("flow"),
		ReplyURL:                   types.StringValue("https://example.invalid"),
		Wtrealm:                    types.StringValue("realm"),
		PropertyMappings:           types.ListNull(types.StringType),
		AssertionValidNotBefore:    types.StringValue("minutes=-5"),
		AssertionValidNotOnOrAfter: types.StringValue("minutes=5"),
		SessionValidNotOnOrAfter:   types.StringValue("minutes=86400"),
		DigestAlgorithm:            types.StringValue(string(api.DIGESTALGORITHMENUM_HTTP___WWW_W3_ORG_2001_04_XMLENCSHA256)),
		SignatureAlgorithm:         types.StringValue(string(api.SIGNATUREALGORITHMENUM_HTTP___WWW_W3_ORG_2001_04_XMLDSIG_MORERSA_SHA256)),
		SignAssertion:              types.BoolValue(true),
		SignLogoutRequest:          types.BoolValue(false),
	})
	require.False(t, diags.HasError(), diags)

	require.NotNil(t, body.SignLogoutRequest,
		"an explicit false must reach the API or sign_logout_request can never be disabled")
	assert.False(t, *body.SignLogoutRequest)

	// discovery #6: a null list serialises as present-but-empty.
	require.NotNil(t, body.PropertyMappings)
	assert.Empty(t, body.PropertyMappings)
}
