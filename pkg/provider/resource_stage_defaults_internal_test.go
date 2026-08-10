package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// The tests below are the executable spec for the rule that an attribute carrying a
// schema Default must be mapped from the API verbatim, while an attribute that can
// actually be null keeps the prior-aware H1 helpers.
//
// They all simulate `terraform import`, which is the only code path where the
// distinction is observable: ImportState writes just the id, so every other attribute
// is null when Read seeds its model from prior state. Feeding that null prior into
// BoolOrNull/Int32OrNull/StringOrNull alongside an API zero value returns null - and
// since an attribute with a Default can never legitimately be null, the next plan
// applies the default and reports a spurious "null -> false" update. SDKv2 wrote the
// concrete value here and planned clean.
//
// Note these cases are exactly the ones a docs diff cannot catch (docs render
// identically either way) and that only bite when the API value equals the type's zero
// value, which is why they need pinning explicitly rather than by inspection.

func TestStagePasswordFromAPI_DefaultedAttributesSurviveImport(t *testing.T) {
	ctx := context.Background()
	r := &stagePasswordResource{}

	// Everything null, as after ImportState.
	data := &stagePasswordModel{
		ID:       types.StringValue("uuid-1"),
		Backends: types.ListNull(types.StringType),
	}

	res := &api.PasswordStage{
		Pk:                         "uuid-1",
		Name:                       "stage",
		Backends:                   []api.BackendsEnum{api.BACKENDSENUM_AUTHENTIK_CORE_AUTH_INBUILT_BACKEND},
		AllowShowPassword:          new(false),
		FailedAttemptsBeforeCancel: new(int32(5)),
		ConfigureFlow:              *api.NewNullableString(nil),
	}

	diags := r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	require.False(t, data.AllowShowPassword.IsNull(), "allow_show_password has a Default, so it must never import as null")
	assert.False(t, data.AllowShowPassword.ValueBool())
	require.False(t, data.FailedAttemptsBeforeCancel.IsNull(), "failed_attempts_before_cancel has a Default, so it must never import as null")
	assert.Equal(t, int32(5), data.FailedAttemptsBeforeCancel.ValueInt32())

	// configure_flow is Optional-only and genuinely nullable, so it stays null.
	assert.True(t, data.ConfigureFlow.IsNull(), "configure_flow has no Default and the API returned null")
}

func TestStageAccountLockdownFromAPI_DefaultedAttributesSurviveImport(t *testing.T) {
	r := &stageAccountLockdownResource{}

	data := &stageAccountLockdownModel{ID: types.StringValue("uuid-1")}

	// All four bools default to true in the schema, so false is the interesting case:
	// it is both a value a user can legitimately set and the type's zero value.
	res := &api.AccountLockdownStage{
		Pk:                        "uuid-1",
		Name:                      "stage",
		DeactivateUser:            new(false),
		SetUnusablePassword:       new(false),
		DeleteSessions:            new(false),
		RevokeTokens:              new(false),
		SelfServiceCompletionFlow: *api.NewNullableString(nil),
	}

	r.fromAPI(data, res)

	for name, v := range map[string]types.Bool{
		"deactivate_user":       data.DeactivateUser,
		"set_unusable_password": data.SetUnusablePassword,
		"delete_sessions":       data.DeleteSessions,
		"revoke_tokens":         data.RevokeTokens,
	} {
		require.Falsef(t, v.IsNull(), "%s has a Default, so it must never import as null", name)
		assert.Falsef(t, v.ValueBool(), "%s should round-trip the API's false", name)
	}

	assert.True(t, data.SelfServiceCompletionFlow.IsNull(), "self_service_completion_flow has no Default and the API returned null")
}

// TestStageRedirectFromAPI_DefaultAndNullableCoexist pins both halves of the rule in one
// resource: keep_context has a Default and must come back concrete, while target_static
// is Optional-only and must stay null even though the API returns "" for it (H1).
func TestStageRedirectFromAPI_DefaultAndNullableCoexist(t *testing.T) {
	r := &stageRedirectResource{}

	data := &stageRedirectModel{ID: types.StringValue("uuid-1")}

	empty := ""
	res := &api.RedirectStage{
		Pk:           "uuid-1",
		Name:         "stage",
		Mode:         api.REDIRECTSTAGEMODEENUM_STATIC,
		KeepContext:  new(false),
		TargetStatic: &empty,
		TargetFlow:   *api.NewNullableString(nil),
	}

	r.fromAPI(data, res)

	require.False(t, data.KeepContext.IsNull(), "keep_context has a Default, so it must never import as null")
	assert.False(t, data.KeepContext.ValueBool())

	assert.True(t, data.TargetStatic.IsNull(), "target_static has no Default and the API returned \"\", so it stays null")
	assert.True(t, data.TargetFlow.IsNull())
}
