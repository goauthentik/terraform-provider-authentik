package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceUsers() *schema.Resource {
	userSchema := map[string]*schema.Schema{}
	for k, v := range dataSourceUser().Schema {
		userSchema[k] = &schema.Schema{}
		*userSchema[k] = *v
		userSchema[k].Computed = true
		userSchema[k].Optional = false
		userSchema[k].Required = false
		userSchema[k].ExactlyOneOf = []string{}
	}
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Description: "Directory --- Get users list",
		Schema: map[string]*schema.Schema{
			"attributes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups_by_name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"groups_by_pk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_superuser": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ordering": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path_startswith": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: userSchema,
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	req := c.client.CoreApi.CoreUsersList(ctx)

	for key := range dataSourceUsers().Schema {
		if v, ok := d.GetOk(key); ok {
			switch key {
			case "attributes":
				req = req.Attributes(v.(string))
			case "email":
				req = req.Email(v.(string))
			case "groups_by_name":
				req = req.GroupsByName(v.([]string))
			case "groups_by_pk":
				req = req.GroupsByPk(v.([]string))
			case "is_active":
				req = req.IsActive(v.(bool))
			case "is_superuser":
				req = req.IsSuperuser(v.(bool))
			case "name":
				req = req.Name(v.(string))
			case "ordering":
				req = req.Ordering(v.(string))
			case "path":
				req = req.Path(v.(string))
			case "path_startswith":
				req = req.PathStartswith(v.(string))
			case "search":
				req = req.Search(v.(string))
			case "username":
				req = req.Username(v.(string))
			case "uuid":
				req = req.Uuid(v.(string))
			}
		}
	}

	users := make([]map[string]any, 0)
	for page := int32(1); true; page++ {
		req = req.Page(page)
		res, hr, err := req.Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}

		for _, userRes := range res.Results {
			u, err := mapFromUser(userRes)
			if err != nil {
				return diag.FromErr(err)
			}
			users = append(users, u)
		}

		if res.Pagination.Next == 0 {
			break
		}
	}

	d.SetId("0")
	helpers.SetWrapper(d, "users", users)
	return diag.Diagnostics{}
}
