package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

func redirectURIList(t *testing.T, entries ...map[string]string) types.List {
	t.Helper()

	elems := make([]attr.Value, 0, len(entries))
	for _, e := range entries {
		attrs := map[string]attr.Value{}
		for k, v := range e {
			attrs[k] = types.StringValue(v)
		}
		elems = append(elems, types.MapValueMust(types.StringType, attrs))
	}
	return types.ListValueMust(redirectURIElementType, elems)
}

// TestProviderOAuth2RedirectURIs_RoundTripAndOptionalType covers allowed_redirect_uris, the
// attribute migration-plan.md singles out as the batch's real risk. It is a TypeList of bare
// TypeMap in SDKv2, which CoreConfigSchema expands to list(map(string)), and the element maps
// carry two keys or three depending on whether redirect_uri_type is set. Getting the omission
// wrong would change the state layout for every existing resource.
func TestProviderOAuth2RedirectURIs_RoundTripAndOptionalType(t *testing.T) {
	ctx := t.Context()

	t.Run("three-key and two-key entries both round-trip", func(t *testing.T) {
		in := redirectURIList(t,
			map[string]string{"matching_mode": "strict", "url": "https://a.invalid", "redirect_uri_type": "url"},
			map[string]string{"matching_mode": "regex", "url": "https://b.invalid/.*"},
		)

		parsed, diags := redirectURIsFromList(ctx, in)
		require.False(t, diags.HasError(), diags)
		require.Len(t, parsed, 2)
		assert.Equal(t, api.RedirectURITypeEnum("url"), parsed[0].RedirectURIType)
		assert.Equal(t, api.RedirectURITypeEnum(""), parsed[1].RedirectURIType,
			"an absent redirect_uri_type key must parse as the empty enum, not error")

		out, diags := redirectURIsToList(parsed)
		require.False(t, diags.HasError(), diags)

		var back []map[string]string
		require.False(t, out.ElementsAs(ctx, &back, false).HasError())
		require.Len(t, back, 2)
		assert.Len(t, back[0], 3, "an entry with a type keeps all three keys")
		assert.Len(t, back[1], 2, "an entry without a type must omit the key, not emit an empty string")
		assert.Equal(t, "https://b.invalid/.*", back[1]["url"])
	})

	t.Run("configured order survives an API reordering", func(t *testing.T) {
		r := &providerOAuth2Resource{}

		data := &providerOAuth2Model{
			AllowedRedirectURIs: redirectURIList(t,
				map[string]string{"matching_mode": "strict", "url": "https://a.invalid"},
				map[string]string{"matching_mode": "strict", "url": "https://b.invalid"},
			),
			PropertyMappings:       types.ListNull(types.StringType),
			GrantTypes:             types.ListNull(types.StringType),
			JWTFederationSources:   types.ListNull(types.StringType),
			JWTFederationProviders: types.ListNull(types.Int32Type),
		}

		// API returns the same pair in the opposite order.
		diags := r.fromAPI(ctx, data, &api.OAuth2Provider{
			Pk:                12,
			Name:              "oauth2",
			AuthorizationFlow: "flow",
			InvalidationFlow:  "flow",
			RedirectUris: []api.RedirectURI{
				{MatchingMode: "strict", Url: "https://b.invalid"},
				{MatchingMode: "strict", Url: "https://a.invalid"},
			},
			AuthenticationFlow: *api.NewNullableString(nil),
			SigningKey:         *api.NewNullableString(nil),
			EncryptionKey:      *api.NewNullableString(nil),
		})
		require.False(t, diags.HasError(), diags)

		var back []map[string]string
		require.False(t, data.AllowedRedirectURIs.ElementsAs(ctx, &back, false).HasError())
		require.Len(t, back, 2)
		assert.Equal(t, "https://a.invalid", back[0]["url"], "configured order must win over the API's")
		assert.Equal(t, "https://b.invalid", back[1]["url"])
	})

	t.Run("never-configured list stays null", func(t *testing.T) {
		r := &providerOAuth2Resource{}

		data := &providerOAuth2Model{
			AllowedRedirectURIs:    types.ListNull(redirectURIElementType),
			PropertyMappings:       types.ListNull(types.StringType),
			GrantTypes:             types.ListNull(types.StringType),
			JWTFederationSources:   types.ListNull(types.StringType),
			JWTFederationProviders: types.ListNull(types.Int32Type),
		}

		diags := r.fromAPI(ctx, data, &api.OAuth2Provider{
			Pk: 12, Name: "oauth2", AuthorizationFlow: "flow", InvalidationFlow: "flow",
			AuthenticationFlow: *api.NewNullableString(nil),
			SigningKey:         *api.NewNullableString(nil),
			EncryptionKey:      *api.NewNullableString(nil),
		})
		require.False(t, diags.HasError(), diags)
		assert.True(t, data.AllowedRedirectURIs.IsNull(), "H1: null in, nothing from the API, so null out")
	})
}

// TestProviderOAuth2ToRequest_GrantTypesOmittedWhenEmpty is the one place discovery #6's
// "lists go out as present-but-empty" rule is deliberately *not* applied: the API rejects an
// empty grant_types and otherwise derives the value from the rest of the configuration.
func TestProviderOAuth2ToRequest_GrantTypesOmittedWhenEmpty(t *testing.T) {
	r := &providerOAuth2Resource{}

	base := func() *providerOAuth2Model {
		return &providerOAuth2Model{
			Name:                   types.StringValue("oauth2"),
			AuthorizationFlow:      types.StringValue("flow"),
			InvalidationFlow:       types.StringValue("flow"),
			ClientID:               types.StringValue("client"),
			ClientType:             types.StringValue(string(api.CLIENTTYPEENUM_CONFIDENTIAL)),
			AccessCodeValidity:     types.StringValue("minutes=1"),
			AccessTokenValidity:    types.StringValue("minutes=10"),
			RefreshTokenValidity:   types.StringValue("days=30"),
			RefreshTokenThreshold:  types.StringValue("seconds=0"),
			IncludeClaimsInIDToken: types.BoolValue(true),
			LogoutMethod:           types.StringValue(string(api.OAUTH2PROVIDERLOGOUTMETHODENUM_BACKCHANNEL)),
			SubMode:                types.StringValue(string(api.SUBMODEENUM_HASHED_USER_ID)),
			IssuerMode:             types.StringValue(string(api.ISSUERMODEENUM_PER_PROVIDER)),
			PropertyMappings:       types.ListNull(types.StringType),
			JWTFederationSources:   types.ListNull(types.StringType),
			JWTFederationProviders: types.ListNull(types.Int32Type),
			AllowedRedirectURIs:    types.ListNull(redirectURIElementType),
			GrantTypes:             types.ListNull(types.StringType),
		}
	}

	t.Run("null grant_types is omitted entirely", func(t *testing.T) {
		body, diags := r.toRequest(t.Context(), base())
		require.False(t, diags.HasError(), diags)
		assert.Nil(t, body.GrantTypes, "the API rejects an empty grant_types, so it must be omitted")

		// The other lists still follow discovery #6.
		require.NotNil(t, body.PropertyMappings)
		assert.Empty(t, body.PropertyMappings)
		require.NotNil(t, body.RedirectUris)
		assert.Empty(t, body.RedirectUris)
	})

	t.Run("explicit grant_types is sent", func(t *testing.T) {
		data := base()
		data.GrantTypes = types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(string(api.GRANTTYPESENUM_AUTHORIZATION_CODE)),
		})

		body, diags := r.toRequest(t.Context(), data)
		require.False(t, diags.HasError(), diags)
		assert.Equal(t, []api.GrantTypesEnum{api.GRANTTYPESENUM_AUTHORIZATION_CODE}, body.GrantTypes)
	})
}

// TestProviderOAuth2FromAPI_PreservesJWKSSources pins the inert H4 attribute: SDKv2 neither
// sent jwks_sources nor read it back, so it holds only what config put there.
func TestProviderOAuth2FromAPI_PreservesJWKSSources(t *testing.T) {
	r := &providerOAuth2Resource{}
	ctx := t.Context()

	configured := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("source-slug")})
	data := &providerOAuth2Model{
		JWKSSources:            configured,
		AllowedRedirectURIs:    types.ListNull(redirectURIElementType),
		PropertyMappings:       types.ListNull(types.StringType),
		GrantTypes:             types.ListNull(types.StringType),
		JWTFederationSources:   types.ListNull(types.StringType),
		JWTFederationProviders: types.ListNull(types.Int32Type),
	}

	diags := r.fromAPI(ctx, data, &api.OAuth2Provider{
		Pk: 12, Name: "oauth2", AuthorizationFlow: "flow", InvalidationFlow: "flow",
		AuthenticationFlow: *api.NewNullableString(nil),
		SigningKey:         *api.NewNullableString(nil),
		EncryptionKey:      *api.NewNullableString(nil),
	})
	require.False(t, diags.HasError(), diags)

	var back []string
	require.False(t, data.JWKSSources.ElementsAs(ctx, &back, false).HasError())
	assert.Equal(t, []string{"source-slug"}, back, "jwks_sources is never refreshed from the server")
}
