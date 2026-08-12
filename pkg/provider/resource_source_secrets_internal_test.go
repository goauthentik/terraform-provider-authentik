package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// TestSourceKerberosFromAPI_GroupMatchingModeUsesGroupField pins a real bug found while
// migrating this resource. SDKv2's Read did:
//
//	helpers.SetWrapper(d, "group_matching_mode", res.UserMatchingMode)
//
// - a copy-paste slip that no other source has (saml and oauth both read GroupMatchingMode
// correctly). Under SDKv2 it silently wrote the wrong value into state whenever the two
// modes differed, producing a permanent diff. Under the framework the same slip would be a
// hard "Provider produced inconsistent result after apply" error, because the plan carries
// the configured group_matching_mode while state would get the user one.
//
// The two modes are set to *different* values here; with the old behaviour
// group_matching_mode would come back as "email_link".
func TestSourceKerberosFromAPI_GroupMatchingModeUsesGroupField(t *testing.T) {
	r := &sourceKerberosResource{}

	data := &sourceKerberosModel{}

	res := &api.KerberosSource{
		Pk:                 "uuid-1",
		Name:               "kerberos",
		Slug:               "kerberos",
		Realm:              "EXAMPLE.COM",
		UserMatchingMode:   api.USERMATCHINGMODEENUM_EMAIL_LINK.Ptr(),
		GroupMatchingMode:  api.GROUPMATCHINGMODEENUM_NAME_LINK.Ptr(),
		AuthenticationFlow: *api.NewNullableString(nil),
		EnrollmentFlow:     *api.NewNullableString(nil),
	}

	r.fromAPI(data, res)

	assert.Equal(t, string(api.USERMATCHINGMODEENUM_EMAIL_LINK), data.UserMatchingMode.ValueString())
	assert.Equal(t, string(api.GROUPMATCHINGMODEENUM_NAME_LINK), data.GroupMatchingMode.ValueString(),
		"group_matching_mode must come from res.GroupMatchingMode, not res.UserMatchingMode")
}

// TestSourceKerberosFromAPI_PreservesSecrets covers three of the 18 H4 attributes at once.
func TestSourceKerberosFromAPI_PreservesSecrets(t *testing.T) {
	r := &sourceKerberosResource{}

	data := &sourceKerberosModel{
		SyncPassword: types.StringValue("sync-secret"),
		SyncKeytab:   types.StringValue("sync-keytab"),
		SpnegoKeytab: types.StringValue("spnego-keytab"),
	}

	res := &api.KerberosSource{
		Pk:                 "uuid-1",
		Name:               "kerberos",
		Slug:               "kerberos",
		Realm:              "EXAMPLE.COM",
		AuthenticationFlow: *api.NewNullableString(nil),
		EnrollmentFlow:     *api.NewNullableString(nil),
	}

	r.fromAPI(data, res)

	assert.Equal(t, "sync-secret", data.SyncPassword.ValueString())
	assert.Equal(t, "sync-keytab", data.SyncKeytab.ValueString())
	assert.Equal(t, "spnego-keytab", data.SpnegoKeytab.ValueString())
}

func TestSourceLDAPFromAPI_PreservesBindPassword(t *testing.T) {
	ctx := context.Background()
	r := &sourceLDAPResource{}

	data := &sourceLDAPModel{
		BindPassword:          types.StringValue("bind-secret"),
		PropertyMappings:      types.ListNull(types.StringType),
		PropertyMappingsGroup: types.ListNull(types.StringType),
	}

	res := &api.LDAPSource{
		Pk:              "uuid-1",
		Name:            "ldap",
		Slug:            "ldap",
		BaseDn:          "dc=example,dc=com",
		ServerUri:       "ldap://example.invalid",
		SyncParentGroup: *api.NewNullableString(nil),
	}

	diags := r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	require.False(t, data.BindPassword.IsNull(), "bind_password is Required and write-only")
	assert.Equal(t, "bind-secret", data.BindPassword.ValueString())

	// additional_user_dn/additional_group_dn default to "" - exactly what the API returns
	// when unset - so discovery #5 says they must come back concrete, not null.
	require.False(t, data.AdditionalUserDN.IsNull(), "additional_user_dn has a Default of \"\"")
	assert.Equal(t, "", data.AdditionalUserDN.ValueString())
	require.False(t, data.AdditionalGroupDN.IsNull())
}

func TestSourceOAuthFromAPI_PreservesConsumerSecret(t *testing.T) {
	ctx := context.Background()
	r := &sourceOAuthResource{}

	data := &sourceOAuthModel{
		ConsumerSecret:        types.StringValue("consumer-secret"),
		PropertyMappings:      types.ListNull(types.StringType),
		PropertyMappingsGroup: types.ListNull(types.StringType),
	}

	res := &api.OAuthSource{
		Pk:                 "uuid-1",
		Name:               "oauth",
		Slug:               "oauth",
		ProviderType:       api.PROVIDERTYPEENUM_GITHUB,
		ConsumerKey:        "consumer-key",
		CallbackUrl:        "https://authentik.invalid/source/oauth/callback/oauth/",
		AuthenticationFlow: *api.NewNullableString(nil),
		EnrollmentFlow:     *api.NewNullableString(nil),
		RequestTokenUrl:    *api.NewNullableString(nil),
		AuthorizationUrl:   *api.NewNullableString(nil),
		AccessTokenUrl:     *api.NewNullableString(nil),
		ProfileUrl:         *api.NewNullableString(nil),
	}

	diags := r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	require.False(t, data.ConsumerSecret.IsNull(), "consumer_secret is Required and write-only")
	assert.Equal(t, "consumer-secret", data.ConsumerSecret.ValueString())
	assert.Equal(t, "consumer-key", data.ConsumerKey.ValueString())

	// Nullable URL fields stay null rather than becoming "".
	assert.True(t, data.RequestTokenURL.IsNull())
	assert.True(t, data.AuthorizationURL.IsNull())

	// oidc_jwks is Optional+Computed with no Default; an absent map must still be valid JSON.
	require.False(t, data.OIDCJWKS.IsNull())
	assert.JSONEq(t, "null", data.OIDCJWKS.ValueString())
}
