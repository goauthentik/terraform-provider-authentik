package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"goauthentik.io/api/v3"
)

func dataSourceGroupMember() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_login": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Description: "Directory --- Get groups by pk or name",
		Schema: map[string]*schema.Schema{
			"pk": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"pk", "name"},
			},
			"num_pk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"pk", "name"},
			},
			"include_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to include group members. Note that depending on group size, this can make the Terraform state a lot larger.",
			},
			"is_superuser": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"parent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"attributes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users_obj": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceGroupMember(),
			},
		},
	}
}

func mapFromGroupMember(member api.GroupMember) (map[string]interface{}, error) {
	m := map[string]interface{}{
		"pk":         int(member.GetPk()),
		"username":   member.GetUsername(),
		"name":       member.GetName(),
		"is_active":  member.GetIsActive(),
		"last_login": "",
		"email":      member.GetEmail(),
		"attributes": "",
		"uid":        member.GetUid(),
	}

	if t, ok := member.GetLastLoginOk(); ok && t != nil {
		if last_login, err := t.MarshalText(); err == nil && last_login != nil {
			m["last_login"] = string(last_login)
		}
	}

	b, err := json.Marshal(member.GetAttributes())
	if err != nil {
		return nil, err
	}
	m["attributes"] = string(b)

	return m, nil
}

func mapFromGroup(group api.Group) (map[string]interface{}, error) {
	m := map[string]interface{}{
		"pk":           group.GetPk(),
		"num_pk":       int(group.GetNumPk()),
		"name":         group.GetName(),
		"is_superuser": group.GetIsSuperuser(),
		"parent":       group.GetParent(),
		"parent_name":  group.GetParentName(),
		"users":        []int{},
		"attributes":   "",
		"users_obj":    []map[string]interface{}{},
	}

	b, err := json.Marshal(group.GetAttributes())
	if err != nil {
		return nil, err
	}
	m["attributes"] = string(b)

	users := make([]int, len(group.GetUsers()))
	for i, user_pk := range group.GetUsers() {
		users[i] = int(user_pk)
	}
	m["users"] = users

	users_obj := make([]map[string]interface{}, len(group.GetUsersObj()))
	for i, member := range group.GetUsersObj() {
		memberMap, err := mapFromGroupMember(member)
		if err != nil {
			return nil, err
		}
		users_obj[i] = memberMap
	}
	m["users_obj"] = users_obj

	return m, nil
}

func setGroup(data *schema.ResourceData, group api.Group) diag.Diagnostics {
	m, err := mapFromGroup(group)
	if err != nil {
		return diag.FromErr(err)
	}
	for key, value := range m {
		switch key {
		case "pk":
			data.SetId(value.(string))
			setWrapper(data, key, value.(string))
		case "num_pk":
			setWrapper(data, key, value.(int))
		case "is_superuser":
			setWrapper(data, key, value.(bool))
		case "users":
			setWrapper(data, key, value.([]int))
		case "users_obj":
			setWrapper(data, key, value.([]map[string]interface{}))
		default:
			setWrapper(data, key, value.(string))
		}
	}
	return diag.Diagnostics{}
}

func dataSourceGroupReadByPk(ctx context.Context, d *schema.ResourceData, c *APIClient, pk string, includeUsers bool) diag.Diagnostics {
	req := c.client.CoreApi.CoreGroupsRetrieve(ctx, pk)
	req = req.IncludeUsers(includeUsers)

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	return setGroup(d, *res)
}

func dataSourceGroupReadByName(ctx context.Context, d *schema.ResourceData, c *APIClient, name string, includeUsers bool) diag.Diagnostics {
	req := c.client.CoreApi.CoreGroupsList(ctx)
	req = req.IncludeUsers(includeUsers)
	req = req.Name(name)

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching groups found")
	}

	if len(res.Results) > 1 {
		return diag.Errorf("Multiple groups found")
	}

	return setGroup(d, res.Results[0])
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	includeUsers := true
	if i := d.Get("include_users"); i != nil {
		includeUsers = i.(bool)
	}

	if n, ok := d.GetOk("pk"); ok {
		return dataSourceGroupReadByPk(ctx, d, c, n.(string), includeUsers)
	}

	if n, ok := d.GetOk("name"); ok {
		return dataSourceGroupReadByName(ctx, d, c, n.(string), includeUsers)
	}

	return diag.Errorf("Neither pk nor name were provided")
}
