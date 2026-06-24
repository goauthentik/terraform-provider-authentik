package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccResourceGroup(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName),
					resource.TestCheckResourceAttr("authentik_group.group", "name", rName),
				),
			},
			{
				Config: testAccResourceGroup(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName+"test"),
					resource.TestCheckResourceAttr("authentik_group.group", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceGroup(name string) string {
	return fmt.Sprintf(`
resource "authentik_user" "name" {
  username = "%[1]s"
  name = "%[1]s"
}
resource "authentik_rbac_role" "role" {
  name = "%[1]s"
}
resource "authentik_group" "group" {
  name = "%[1]s"
  users = [authentik_user.name.id]
  is_superuser = true
  roles = [authentik_rbac_role.role.id]
}
`, name)
}

func TestResourceGroupReadRolesPreserveConfiguredOrder(t *testing.T) {
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v3/core/groups/group-1/", r.URL.Path)
		assert.Equal(t, "false", r.URL.Query().Get("include_users"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"pk": "group-1",
			"num_pk": 1,
			"name": "infrastructure",
			"is_superuser": false,
			"parents": [],
			"parents_obj": [],
			"users": [],
			"users_obj": [],
			"attributes": {},
			"roles": ["role-c", "role-a", "role-b", "role-d"],
			"roles_obj": [],
			"inherited_roles_obj": [],
			"children": [],
			"children_obj": []
		}`))
	}))
	t.Cleanup(server.Close)

	config := api.NewConfiguration()
	config.Servers = api.ServerConfigurations{{
		URL: server.URL + "/api/v3",
	}}
	client := &APIClient{
		client: api.NewAPIClient(config),
	}

	d := schema.TestResourceDataRaw(t, resourceGroup().Schema, map[string]any{
		"name":       "infrastructure",
		"attributes": "{}",
		"roles":      []any{"role-a", "role-b", "role-c"},
	})
	d.SetId("group-1")

	diags := resourceGroupRead(ctx, d, client)

	require.False(t, diags.HasError(), diags)
	assert.Equal(t, []string{
		"role-a",
		"role-b",
		"role-c",
		"role-d",
	}, helpers.CastSlice[string](d, "roles"))
}
