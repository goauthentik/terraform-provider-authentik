package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

func TestResourceBase_Configure(t *testing.T) {
	t.Run("nil ProviderData is a no-op, not an error", func(t *testing.T) {
		r := &resourceBase{}
		resp := &resource.ConfigureResponse{}
		r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: nil}, resp)
		assert.False(t, resp.Diagnostics.HasError())
		assert.Nil(t, r.client)
	})

	t.Run("wrong ProviderData type is an error", func(t *testing.T) {
		r := &resourceBase{}
		resp := &resource.ConfigureResponse{}
		r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: "not-a-client"}, resp)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("correct ProviderData type is stored", func(t *testing.T) {
		r := &resourceBase{}
		client := &api.APIClient{}
		resp := &resource.ConfigureResponse{}
		r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: client}, resp)
		require.False(t, resp.Diagnostics.HasError())
		assert.Same(t, client, r.client)
	})
}

func TestResourceBase_Span(t *testing.T) {
	r := &resourceBase{}
	finish := r.span(context.Background(), "create")
	require.NotNil(t, finish)
	finish()
}
