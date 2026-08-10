package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// This file is the executable spec for H4 across the four secret-carrying stage resources
// migrated in Phase 3 batch 1 part 6.
//
// SDKv2's d.Set was incremental: a Read that never called SetWrapper for an attribute left
// whatever was already in the ResourceData. resp.State.Set is wholesale - anything the
// model does not carry becomes null - so the equivalent is that Create/Update seed the
// model from req.Plan and Read seeds it from req.State, and fromAPI then simply does not
// assign the attributes the API never returns. These tests assert that the seeded value
// survives, which is the only thing standing between a write-only secret and being wiped
// from state on the next refresh.
//
// For client_secret and private_key, both Required, being wiped would not merely lose the
// value: it would drop a Required attribute out of state and fail the apply. The matching
// acceptance tests use ImportStateVerifyIgnore for exactly these attributes, since import
// starts from an empty model and so genuinely cannot recover them - the same limitation
// SDKv2 had.

func TestStageCaptchaFromAPI_PreservesPrivateKey(t *testing.T) {
	r := &stageCaptchaResource{}

	data := &stageCaptchaModel{
		ID:         types.StringValue("uuid-1"),
		PrivateKey: types.StringValue("super-secret"),
	}

	// A real CaptchaStage response carries no private_key field at all.
	res := &api.CaptchaStage{
		Pk:        "uuid-1",
		Name:      "captcha",
		PublicKey: "public",
	}

	r.fromAPI(data, res)

	require.False(t, data.PrivateKey.IsNull(), "private_key is Required and write-only; nulling it would fail the apply")
	assert.Equal(t, "super-secret", data.PrivateKey.ValueString())
	assert.Equal(t, "public", data.PublicKey.ValueString())
}

func TestStageAuthenticatorDuoFromAPI_PreservesSecrets(t *testing.T) {
	r := &stageAuthenticatorDuoResource{}

	data := &stageAuthenticatorDuoModel{
		ID:             types.StringValue("uuid-1"),
		ClientSecret:   types.StringValue("client-secret"),
		AdminSecretKey: types.StringValue("admin-secret"),
	}

	res := &api.AuthenticatorDuoStage{
		Pk:            "uuid-1",
		Name:          "duo",
		ClientId:      "client-id",
		ApiHostname:   "https://example.invalid",
		ConfigureFlow: *api.NewNullableString(nil),
	}

	r.fromAPI(data, res)

	require.False(t, data.ClientSecret.IsNull(), "client_secret is Required and write-only")
	assert.Equal(t, "client-secret", data.ClientSecret.ValueString())
	require.False(t, data.AdminSecretKey.IsNull(), "admin_secret_key is write-only")
	assert.Equal(t, "admin-secret", data.AdminSecretKey.ValueString())
	assert.Equal(t, "client-id", data.ClientID.ValueString())
}

func TestStageEmailFromAPI_PreservesPassword(t *testing.T) {
	r := &stageEmailResource{}

	data := &stageEmailModel{
		ID:       types.StringValue("uuid-1"),
		Password: types.StringValue("smtp-password"),
	}

	res := &api.EmailStage{Pk: "uuid-1", Name: "email"}

	r.fromAPI(data, res)

	require.False(t, data.Password.IsNull(), "password is write-only and must survive a Read")
	assert.Equal(t, "smtp-password", data.Password.ValueString())
}

func TestStageAuthenticatorEmailFromAPI_PreservesPassword(t *testing.T) {
	r := &stageAuthenticatorEmailResource{}

	data := &stageAuthenticatorEmailModel{
		ID:       types.StringValue("uuid-1"),
		Password: types.StringValue("smtp-password"),
	}

	res := &api.AuthenticatorEmailStage{
		Pk:            "uuid-1",
		Name:          "email-authenticator",
		ConfigureFlow: *api.NewNullableString(nil),
	}

	r.fromAPI(data, res)

	require.False(t, data.Password.IsNull(), "password is write-only and must survive a Read")
	assert.Equal(t, "smtp-password", data.Password.ValueString())
}

// TestStageEmailFromAPI_DefaultedAttributesSurviveImport is discovery #5 applied to the
// resource with the most defaulted attributes in the batch. Every value below is the
// respective type's zero value, which is precisely when the prior-aware helpers would
// collapse a defaulted attribute to null and make the next plan report a spurious update.
// use_tls/use_ssl/username are the only three attributes here without a Default, so they
// are the only three that must stay null.
func TestStageEmailFromAPI_DefaultedAttributesSurviveImport(t *testing.T) {
	r := &stageEmailResource{}

	data := &stageEmailModel{ID: types.StringValue("uuid-1")}

	res := &api.EmailStage{
		Pk:                    "uuid-1",
		Name:                  "email",
		UseGlobalSettings:     new(false),
		Host:                  new(""),
		Port:                  new(int32(0)),
		Timeout:               new(int32(0)),
		FromAddress:           new(""),
		TokenExpiry:           new(""),
		Subject:               new(""),
		Template:              new(""),
		ActivateUserOnSuccess: new(false),
		RecoveryMaxAttempts:   new(int32(0)),
		RecoveryCacheTimeout:  new(""),
	}

	r.fromAPI(data, res)

	for name, v := range map[string]attrNullable{
		"use_global_settings":      data.UseGlobalSettings,
		"host":                     data.Host,
		"port":                     data.Port,
		"timeout":                  data.Timeout,
		"from_address":             data.FromAddress,
		"token_expiry":             data.TokenExpiry,
		"subject":                  data.Subject,
		"template":                 data.Template,
		"activate_user_on_success": data.ActivateUserOnSuccess,
		"recovery_max_attempts":    data.RecoveryMaxAttempts,
		"recovery_cache_timeout":   data.RecoveryCacheTimeout,
	} {
		assert.Falsef(t, v.IsNull(), "%s has a Default, so it must never import as null", name)
	}

	assert.True(t, data.UseTLS.IsNull(), "use_tls has no Default, so null + false stays null")
	assert.True(t, data.UseSSL.IsNull(), "use_ssl has no Default, so null + false stays null")
	assert.True(t, data.Username.IsNull(), "username has no Default, so null + \"\" stays null")
}

// attrNullable is the shared slice of the types.* value API these tests need, so the loop
// above can hold Bool, String and Int32 values in one map.
type attrNullable interface {
	IsNull() bool
}
