package sdkprovider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

func TestResourceStageAuthenticatorValidateSchemaToProvider(t *testing.T) {
	for _, tc := range []struct {
		name     string
		action   string
		expected api.NotConfiguredActionEnum
	}{
		{"skip", "skip", api.NOTCONFIGUREDACTIONENUM_SKIP},
		{"deny", "deny", api.NOTCONFIGUREDACTIONENUM_DENY},
		{"configure", "configure", api.NOTCONFIGUREDACTIONENUM_CONFIGURE},
		// `not_configured_action` is Required, so it must always end up in the request
		// body. Sourcing it via `helpers.GetP` (`d.GetOk`) yielded a nil pointer for the
		// zero value, and the model's `omitempty` then dropped the field from the UPDATE
		// PUT entirely, leaving authentik on its own default.
		{"zero value is still serialized", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceStageAuthenticatorValidate().Schema, map[string]any{
				"name":                  "test",
				"not_configured_action": tc.action,
			})
			r := resourceStageAuthenticatorValidateSchemaToProvider(d)

			require.NotNil(t, r.NotConfiguredAction, "not_configured_action must never be nil")
			assert.Equal(t, tc.expected, *r.NotConfiguredAction)

			b, err := json.Marshal(r)
			assert.NoError(t, err)
			assert.Contains(t, string(b), `"not_configured_action"`, "not_configured_action must never be omitted from the request body")
		})
	}
}
