package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// TestStageIdentificationToRequest_UnsetEncodings is the executable spec for this
// resource's three different "unset" encodings, which SDKv2 expressed through three
// different helpers and which collapse into one another if ported carelessly:
//
//	new(d.Get(...).(string))  -> pointer to ""  -> clear the field  (password/captcha/webauthn stage)
//	helpers.GetP[string]      -> nil            -> leave unchanged  (the three flow refs)
//	helpers.CastSlice         -> non-nil []T{}  -> send empty array (user_fields, sources)
//
// The list case is the one with teeth: the request model gates list fields on
// `if !IsNil(o.X)`, so a nil slice is omitted, the server keeps the old contents, and
// because the plan said null the framework then raises "Provider produced inconsistent
// result after apply".
func TestStageIdentificationToRequest_UnsetEncodings(t *testing.T) {
	ctx := context.Background()
	r := &stageIdentificationResource{}

	data := &stageIdentificationModel{
		Name:                    types.StringValue("ident"),
		UserFields:              types.ListNull(types.StringType),
		Sources:                 types.ListNull(types.StringType),
		PasswordStage:           types.StringNull(),
		CaptchaStage:            types.StringNull(),
		WebauthnStage:           types.StringNull(),
		EnrollmentFlow:          types.StringNull(),
		RecoveryFlow:            types.StringNull(),
		PasswordlessFlow:        types.StringNull(),
		CaseInsensitiveMatching: types.BoolNull(),
		ShowMatchedUser:         types.BoolValue(true),
		PretendUserExists:       types.BoolValue(true),
		ShowSourceLabels:        types.BoolValue(false),
		EnableRememberMe:        types.BoolValue(false),
	}

	body, diags := r.toRequest(ctx, data)
	require.False(t, diags.HasError(), diags)

	// Lists must be present-but-empty so clearing them clears server-side.
	require.NotNil(t, body.UserFields, "user_fields must be sent as [] to clear it, not omitted")
	assert.Empty(t, body.UserFields)
	require.NotNil(t, body.Sources, "sources must be sent as [] to clear it, not omitted")
	assert.Empty(t, body.Sources)

	// The three stage references clear the field: set, with an empty-string value.
	for name, v := range map[string]api.NullableString{
		"password_stage": body.PasswordStage,
		"captcha_stage":  body.CaptchaStage,
		"webauthn_stage": body.WebauthnStage,
	} {
		require.Truef(t, v.IsSet(), "%s must be sent to clear it server-side", name)
		require.NotNilf(t, v.Get(), "%s must be \"\", not JSON null", name)
		assert.Equalf(t, "", *v.Get(), "%s should clear via an empty string", name)
	}

	// The three flow references leave the field unchanged: omitted entirely.
	for name, v := range map[string]api.NullableString{
		"enrollment_flow":   body.EnrollmentFlow,
		"recovery_flow":     body.RecoveryFlow,
		"passwordless_flow": body.PasswordlessFlow,
	} {
		require.Truef(t, v.IsSet(), "%s is built via NewNullableString, so the wrapper is always set", name)
		assert.Nilf(t, v.Get(), "%s must serialise as JSON null (GetP semantics), not \"\"", name)
	}

	// case_insensitive_matching has no Default but SDKv2 always sent it as false.
	require.NotNil(t, body.CaseInsensitiveMatching)
	assert.False(t, *body.CaseInsensitiveMatching)
}

// TestStageIdentificationFromAPI_PreservesNullAndOrder covers the read direction: a
// never-configured optional attribute stays null even though the API reports it as
// absent, the four defaulted bools come back concrete so import doesn't plan a spurious
// update, and user_fields keeps the configured order rather than the API's.
func TestStageIdentificationFromAPI_PreservesNullAndOrder(t *testing.T) {
	ctx := context.Background()
	r := &stageIdentificationResource{}

	configured, d := types.ListValueFrom(ctx, types.StringType, []string{"email", "username"})
	require.False(t, d.HasError(), d)

	data := &stageIdentificationModel{
		UserFields: configured,
		Sources:    types.ListNull(types.StringType),
	}

	res := &api.IdentificationStage{
		Pk:   "uuid-1",
		Name: "ident",
		// API reports the same set in the opposite order.
		UserFields:        []api.UserFieldsEnum{api.USERFIELDSENUM_USERNAME, api.USERFIELDSENUM_EMAIL},
		PasswordStage:     *api.NewNullableString(nil),
		CaptchaStage:      *api.NewNullableString(nil),
		WebauthnStage:     *api.NewNullableString(nil),
		EnrollmentFlow:    *api.NewNullableString(nil),
		RecoveryFlow:      *api.NewNullableString(nil),
		PasswordlessFlow:  *api.NewNullableString(nil),
		ShowMatchedUser:   new(false),
		PretendUserExists: new(false),
		ShowSourceLabels:  new(false),
		EnableRememberMe:  new(false),
	}

	diags := r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	var gotUserFields []string
	require.False(t, data.UserFields.ElementsAs(ctx, &gotUserFields, false).HasError())
	assert.Equal(t, []string{"email", "username"}, gotUserFields, "configured order must win over the API's")

	assert.True(t, data.Sources.IsNull(), "sources was null and the API returned nothing, so it stays null")
	assert.True(t, data.PasswordStage.IsNull())
	assert.True(t, data.EnrollmentFlow.IsNull())

	// All four have Defaults, so false must survive as false rather than becoming null.
	require.False(t, data.ShowMatchedUser.IsNull())
	assert.False(t, data.ShowMatchedUser.ValueBool())
	require.False(t, data.PretendUserExists.IsNull())
	require.False(t, data.ShowSourceLabels.IsNull())
	require.False(t, data.EnableRememberMe.IsNull())

	// case_insensitive_matching has no Default, so null + zero value stays null.
	assert.True(t, data.CaseInsensitiveMatching.IsNull())
}
