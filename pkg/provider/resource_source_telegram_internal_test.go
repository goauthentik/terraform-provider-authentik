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

// authentik_source_telegram is the first resource to hit H3 and H4 at once: its id is
// derived from the mutable slug, and bot_token is write-only. It carries this part's unit
// tests for that reason.

// TestSourceTelegramFromAPI_PreservesBotTokenAndSlugID pins both hazards together.
func TestSourceTelegramFromAPI_PreservesBotTokenAndSlugID(t *testing.T) {
	ctx := context.Background()
	r := &sourceTelegramResource{}

	data := &sourceTelegramModel{
		ID:                    types.StringValue("old-slug"),
		BotToken:              types.StringValue("secret-token"),
		PropertyMappings:      types.ListNull(types.StringType),
		PropertyMappingsGroup: types.ListNull(types.StringType),
	}

	// A real TelegramSource response has no bot_token field at all, and here the slug has
	// been renamed server-side.
	res := &api.TelegramSource{
		Pk:                    "uuid-1",
		Name:                  "telegram",
		Slug:                  "new-slug",
		BotUsername:           "bot",
		PreAuthenticationFlow: "flow-uuid",
		AuthenticationFlow:    *api.NewNullableString(nil),
		EnrollmentFlow:        *api.NewNullableString(nil),
	}

	diags := r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	require.False(t, data.BotToken.IsNull(), "bot_token is Required and write-only; nulling it would fail the apply")
	assert.Equal(t, "secret-token", data.BotToken.ValueString())

	// H3: id tracks the slug, not the Pk - and it must follow a server-side rename.
	assert.Equal(t, "new-slug", data.ID.ValueString(), "id is res.Slug for this resource, not res.Pk")
	assert.Equal(t, "new-slug", data.Slug.ValueString())
	assert.Equal(t, "uuid-1", data.UUID.ValueString(), "uuid is the stable Pk")

	// The two flow references are genuinely nullable and unset here.
	assert.True(t, data.AuthenticationFlow.IsNull())
	assert.True(t, data.EnrollmentFlow.IsNull())
}

// TestSourceTelegramToRequest_RequestMessageAccessCanBeTurnedOff is the GetP[bool] defect
// again, this time on a source. SDKv2 built request_message_access with GetP[bool], so it was
// omitted whenever false and - the API leaving absent fields unchanged on PUT - could not be
// switched back off once enabled. Asserted on the marshalled body because omitempty is what
// decided the old outcome.
func TestSourceTelegramToRequest_RequestMessageAccessCanBeTurnedOff(t *testing.T) {
	ctx := context.Background()
	r := &sourceTelegramResource{}

	data := &sourceTelegramModel{
		Name:                  types.StringValue("telegram"),
		Slug:                  types.StringValue("telegram"),
		Enabled:               types.BoolValue(true),
		UserPathTemplate:      types.StringValue("goauthentik.io/sources/%(slug)s"),
		PolicyEngineMode:      types.StringValue(string(api.POLICYENGINEMODE_ANY)),
		UserMatchingMode:      types.StringValue(string(api.USERMATCHINGMODEENUM_IDENTIFIER)),
		PreAuthenticationFlow: types.StringValue("flow-uuid"),
		BotUsername:           types.StringValue("bot"),
		BotToken:              types.StringValue("secret-token"),
		RequestMessageAccess:  types.BoolValue(false),
		PropertyMappings:      types.ListNull(types.StringType),
		PropertyMappingsGroup: types.ListNull(types.StringType),
	}

	body, diags := r.toRequest(ctx, data)
	require.False(t, diags.HasError(), diags)

	require.NotNil(t, body.RequestMessageAccess)
	assert.False(t, *body.RequestMessageAccess)

	b, err := json.Marshal(body)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"request_message_access"`,
		"an explicit false must reach the API or the setting can never be disabled")

	// discovery #6: both lists must serialise as present-but-empty, not be omitted.
	require.NotNil(t, body.UserPropertyMappings)
	assert.Empty(t, body.UserPropertyMappings)
	require.NotNil(t, body.GroupPropertyMappings)
	assert.Empty(t, body.GroupPropertyMappings)
}

// TestSourceSCIMFromAPI_ComputedURLAndToken covers the two Computed-only attributes SDKv2
// filled from a redundant second retrieve of the same object.
func TestSourceSCIMFromAPI_ComputedURLAndToken(t *testing.T) {
	ctx := context.Background()
	r := &sourceSCIMResource{}

	data := &sourceSCIMModel{
		PropertyMappings:      types.ListNull(types.StringType),
		PropertyMappingsGroup: types.ListNull(types.StringType),
	}

	res := &api.SCIMSource{
		Pk:       "uuid-1",
		Name:     "scim",
		Slug:     "scim",
		RootUrl:  "https://authentik.invalid/source/scim/scim/v2",
		TokenObj: api.Token{Identifier: "token-identifier"},
	}

	diags := r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	assert.Equal(t, "https://authentik.invalid/source/scim/scim/v2", data.SCIMURL.ValueString())
	assert.Equal(t, "token-identifier", data.Token.ValueString())
	assert.Equal(t, "scim", data.ID.ValueString(), "id is the slug (H3)")
	assert.Equal(t, "uuid-1", data.UUID.ValueString())
}
