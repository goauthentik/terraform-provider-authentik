package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceGroups() *schema.Resource {
	groupSchema := map[string]*schema.Schema{}
	for k, v := range dataSourceGroup().Schema {
		if v.Default != nil {
			continue
		}
		groupSchema[k] = &schema.Schema{}
		*groupSchema[k] = *v
		groupSchema[k].Computed = true
		groupSchema[k].Optional = false
		groupSchema[k].Required = false
		groupSchema[k].ExactlyOneOf = []string{}
	}
	return &schema.Resource{
		ReadContext: dataSourceGroupsRead,
		Description: "Directory --- Get groups list",
		Schema: map[string]*schema.Schema{
			"attributes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_superuser": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"include_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to include group members. Note that depending on group size, this can make the Terraform state a lot larger.",
			},
			"members_by_pk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"members_by_username": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ordering": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: groupSchema,
				},
			},
		},
	}
}

func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	req := c.client.CoreApi.CoreGroupsList(ctx)

	for key, _schema := range dataSourceGroups().Schema {
		var v any
		if _schema.Type == schema.TypeBool && _schema.Default != nil {
			v = d.Get(key)
			if v == nil {
				continue
			}
		} else {
			vv, ok := d.GetOk(key)
			if !ok {
				continue
			}
			v = vv
		}
		switch key {
		case "attributes":
			req = req.Attributes(v.(string))
		case "is_superuser":
			req = req.IsSuperuser(v.(bool))
		case "members_by_pk":
			members := make([]int32, len(v.([]int)))
			for i, pk := range v.([]int) {
				members[i] = int32(pk)
			}
			req = req.MembersByPk(members)
		case "members_by_username":
			req = req.MembersByUsername(v.([]string))
		case "name":
			req = req.Name(v.(string))
		case "ordering":
			req = req.Ordering(v.(string))
		case "search":
			req = req.Search(v.(string))
		case "include_users":
			req = req.IncludeUsers(v.(bool))
		}
	}

	groups := make([]map[string]any, 0)
	for page := int32(1); true; page++ {
		req = req.Page(page)
		res, hr, err := req.Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}

		for _, groupRes := range res.Results {
			u, err := mapFromGroup(groupRes)
			if err != nil {
				return diag.FromErr(err)
			}
			groups = append(groups, u)
		}

		if res.Pagination.Next == 0 {
			break
		}
	}

	d.SetId("0")
	helpers.SetWrapper(d, "groups", groups)
	return diag.Diagnostics{}
}
