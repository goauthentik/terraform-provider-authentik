package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// TestStageAuthenticatorValidateToRequest_NotConfiguredAction is the executable spec for
// #935, ported from the SDKv2 test of the same name. It gets simpler rather than harder in
// the framework: there is no schema.TestResourceDataRaw and no Terraform machinery at all,
// just a model literal and a direct toRequest call.
//
// not_configured_action is Required, so it must always end up in the request body. It used
// to be sourced through d.GetOk, which yielded a nil pointer for the zero value, and the
// model's omitempty then dropped the field from the UPDATE PUT entirely - leaving authentik
// on its own default. The zero-value row is the regression: it cannot happen through
// normal config (the attribute is Required and validated against the enum), which is
// exactly why it needs pinning here rather than in an acceptance test.
func TestStageAuthenticatorValidateToRequest_NotConfiguredAction(t *testing.T) {
	ctx := context.Background()
	r := &stageAuthenticatorValidateResource{}

	for _, tc := range []struct {
		name     string
		action   string
		expected api.NotConfiguredActionEnum
	}{
		{"skip", "skip", api.NOTCONFIGUREDACTIONENUM_SKIP},
		{"deny", "deny", api.NOTCONFIGUREDACTIONENUM_DENY},
		{"configure", "configure", api.NOTCONFIGUREDACTIONENUM_CONFIGURE},
		{"zero value is still serialized", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := &stageAuthenticatorValidateModel{
				Name:                       types.StringValue("test"),
				NotConfiguredAction:        types.StringValue(tc.action),
				DeviceClasses:              types.ListNull(types.StringType),
				ConfigurationStages:        types.ListNull(types.StringType),
				WebauthnAllowedDeviceTypes: types.ListNull(types.StringType),
				WebauthnHints:              types.ListNull(types.StringType),
				LastAuthThreshold:          types.StringValue("seconds=0"),
				WebauthnUserVerification:   types.StringValue(string(api.USERVERIFICATIONENUM_PREFERRED)),
				EmailOtpThrottlingFactor:   types.Float64Value(1),
				SmsOtpThrottlingFactor:     types.Float64Value(1),
				TotpOtpThrottlingFactor:    types.Float64Value(1),
				StaticOtpThrottlingFactor:  types.Float64Value(1),
			}

			body, diags := r.toRequest(ctx, data)
			require.False(t, diags.HasError(), diags)

			require.NotNil(t, body.NotConfiguredAction, "not_configured_action must never be nil")
			assert.Equal(t, tc.expected, *body.NotConfiguredAction)

			b, err := json.Marshal(body)
			require.NoError(t, err)
			assert.Contains(t, string(b), `"not_configured_action"`, "not_configured_action must never be omitted from the request body")
		})
	}
}

// TestStageAuthenticatorValidateFromAPI_MergesEveryList pins the read direction for all
// four lists. Each is Optional-only, so state must match config exactly after apply or the
// framework raises a data-consistency error - which means the two lists SDKv2 did not merge
// (device_classes, webauthn_hints) have to be merged here too.
func TestStageAuthenticatorValidateFromAPI_MergesEveryList(t *testing.T) {
	ctx := context.Background()
	r := &stageAuthenticatorValidateResource{}

	deviceClasses, d := types.ListValueFrom(ctx, types.StringType, []string{"webauthn", "totp"})
	require.False(t, d.HasError(), d)
	hints, d := types.ListValueFrom(ctx, types.StringType, []string{"client-device", "security-key"})
	require.False(t, d.HasError(), d)

	data := &stageAuthenticatorValidateModel{
		DeviceClasses:              deviceClasses,
		WebauthnHints:              hints,
		ConfigurationStages:        types.ListNull(types.StringType),
		WebauthnAllowedDeviceTypes: types.ListNull(types.StringType),
	}

	res := &api.AuthenticatorValidateStage{
		Pk:   "uuid-1",
		Name: "validate",
		// API reports both in the opposite order.
		DeviceClasses: []api.DeviceClassesEnum{api.DEVICECLASSESENUM_TOTP, api.DEVICECLASSESENUM_WEBAUTHN},
		WebauthnHints: []api.WebAuthnHintEnum{api.WEBAUTHNHINTENUM_SECURITY_KEY, api.WEBAUTHNHINTENUM_CLIENT_DEVICE},
	}

	diags := r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	var gotClasses []string
	require.False(t, data.DeviceClasses.ElementsAs(ctx, &gotClasses, false).HasError())
	assert.Equal(t, []string{"webauthn", "totp"}, gotClasses, "configured order must win over the API's")

	var gotHints []string
	require.False(t, data.WebauthnHints.ElementsAs(ctx, &gotHints, false).HasError())
	assert.Equal(t, []string{"client-device", "security-key"}, gotHints, "configured order must win over the API's")

	assert.True(t, data.ConfigurationStages.IsNull(), "null in, nothing from the API, so null out")
	assert.True(t, data.WebauthnAllowedDeviceTypes.IsNull())
}
