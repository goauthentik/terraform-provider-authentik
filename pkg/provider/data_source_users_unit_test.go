package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// Filtering `authentik_users` by group used to panic with
// "interface {} is []interface {}, not []string", because list attributes
// come back from ResourceData as []any.
func TestDataSourceUsersGroupFilters(t *testing.T) {
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"pagination":{"next":0,"previous":0,"count":0,"current":1,"total_pages":1,"start_index":0,"end_index":0},"results":[],"autocomplete":{}}`))
	}))
	defer srv.Close()

	cfg := api.NewConfiguration()
	cfg.Servers = api.ServerConfigurations{{URL: srv.URL + "/api/v3"}}
	cfg.HTTPClient = srv.Client()
	client := &APIClient{client: api.NewAPIClient(cfg)}

	d := schema.TestResourceDataRaw(t, dataSourceUsers().Schema, map[string]any{
		"groups_by_name": []any{"admins", "ops"},
		"groups_by_pk":   []any{"7d1d7bb0-4f2b-4b1a-9f6e-0c8a5b6d1e2f"},
	})

	diags := dataSourceUsersRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.ElementsMatch(t, []string{"admins", "ops"}, query["groups_by_name"])
	assert.Equal(t, []string{"7d1d7bb0-4f2b-4b1a-9f6e-0c8a5b6d1e2f"}, query["groups_by_pk"])
}
