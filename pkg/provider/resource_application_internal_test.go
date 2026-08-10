package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// TestApplicationToRequest_ClearingSendsEmptyString is the executable spec for #865:
// clearing group/meta_icon/meta_launch_url in config must send "" to the API to clear
// the field server-side, not omit it (which would leave the field unchanged). This is
// why these three use StringPtrEmpty rather than StringPtr in toRequest.
func TestApplicationToRequest_ClearingSendsEmptyString(t *testing.T) {
	ctx := context.Background()
	r := &applicationResource{}

	data := &applicationModel{
		Name:                 types.StringValue("app"),
		Slug:                 types.StringValue("app"),
		Group:                types.StringNull(),
		MetaIcon:             types.StringNull(),
		MetaLaunchURL:        types.StringNull(),
		MetaDescription:      types.StringNull(),
		MetaPublisher:        types.StringNull(),
		ProtocolProvider:     types.Int32Null(),
		BackchannelProviders: types.ListNull(types.Int32Type),
		PolicyEngineMode:     types.StringValue(string(api.POLICYENGINEMODE_ANY)),
		OpenInNewTab:         types.BoolValue(false),
		MetaHide:             types.BoolValue(false),
	}

	req, diags := r.toRequest(ctx, data)
	require.False(t, diags.HasError(), diags)

	require.NotNil(t, req.Group, "group must be sent as a pointer to clear it server-side, not omitted")
	assert.Equal(t, "", *req.Group)
	require.NotNil(t, req.MetaIcon)
	assert.Equal(t, "", *req.MetaIcon)
	require.NotNil(t, req.MetaLaunchUrl)
	assert.Equal(t, "", *req.MetaLaunchUrl)

	// meta_description/meta_publisher use StringPtr (null -> omit), the opposite choice
	assert.Nil(t, req.MetaDescription)
	assert.Nil(t, req.MetaPublisher)
}

// TestApplicationFromAPI_PreservesNullVsEmpty is H1's core guarantee for this resource:
// a never-configured optional string stays null in state even though the API returns
// "" for it, so Terraform doesn't see a state/plan mismatch on every Read.
func TestApplicationFromAPI_PreservesNullVsEmpty(t *testing.T) {
	ctx := context.Background()
	r := &applicationResource{}

	data := &applicationModel{
		Group:                types.StringNull(),
		MetaIcon:             types.StringNull(),
		BackchannelProviders: types.ListNull(types.Int32Type),
		OpenInNewTab:         types.BoolValue(false),
		MetaHide:             types.BoolValue(false),
	}

	empty := ""
	res := &api.Application{
		Pk:           "uuid-1",
		Name:         "app",
		Slug:         "app",
		Group:        &empty,
		MetaIcon:     &empty,
		OpenInNewTab: new(false),
		MetaHide:     new(false),
		Provider:     *api.NewNullableInt32(nil),
	}

	diags := r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	assert.True(t, data.Group.IsNull(), "group should stay null: prior was null and API returned \"\"")
	assert.True(t, data.MetaIcon.IsNull(), "meta_icon should stay null: prior was null and API returned \"\"")
	assert.True(t, data.ProtocolProvider.IsNull())
}
