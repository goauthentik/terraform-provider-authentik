package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

func TestDataSourceBase_Configure(t *testing.T) {
	t.Run("nil ProviderData is a no-op, not an error", func(t *testing.T) {
		d := &dataSourceBase{}
		resp := &datasource.ConfigureResponse{}
		d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: nil}, resp)
		assert.False(t, resp.Diagnostics.HasError())
		assert.Nil(t, d.client)
	})

	t.Run("wrong ProviderData type is an error", func(t *testing.T) {
		d := &dataSourceBase{}
		resp := &datasource.ConfigureResponse{}
		d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: 42}, resp)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("correct ProviderData type is stored", func(t *testing.T) {
		d := &dataSourceBase{}
		client := &api.APIClient{}
		resp := &datasource.ConfigureResponse{}
		d.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: client}, resp)
		require.False(t, resp.Diagnostics.HasError())
		assert.Same(t, client, d.client)
	})
}

func TestDataSourceBase_Span(t *testing.T) {
	d := &dataSourceBase{}
	finish := d.span(context.Background(), "read")
	require.NotNil(t, finish)
	finish()
}
