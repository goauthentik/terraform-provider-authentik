package provider

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Description: "Directory --- Get users by pk or username",
		Schema: map[string]*schema.Schema{
			"pk": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"pk", "username"},
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"pk", "username"},
			},
			"type": {
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
			"date_joined": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_superuser": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"avatar": {
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
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func mapFromUser(user api.User) (map[string]any, error) {
	m := map[string]any{
		"pk":           int(user.Pk),
		"username":     user.GetUsername(),
		"name":         user.GetName(),
		"is_active":    user.GetIsActive(),
		"last_login":   "",
		"is_superuser": user.GetIsSuperuser(),
		"groups":       user.GetGroups(),
		"email":        user.GetEmail(),
		"avatar":       user.GetAvatar(),
		"attributes":   "",
		"uid":          user.GetUid(),
		"uuid":         user.GetUuid(),
		"path":         user.GetPath(),
		"type":         string(user.GetType()),
	}

	if t, ok := user.GetLastLoginOk(); ok && t != nil {
		if last_login, err := t.MarshalText(); err == nil && last_login != nil {
			m["last_login"] = string(last_login)
		}
	}
	if t, ok := user.GetDateJoinedOk(); ok && t != nil {
		if date_joined, err := t.MarshalText(); err == nil && date_joined != nil {
			m["date_joined"] = string(date_joined)
		}
	}

	b, err := json.Marshal(user.GetAttributes())
	if err != nil {
		return nil, err
	}
	m["attributes"] = string(b)
	return m, nil
}

func setUser(data *schema.ResourceData, user api.User) diag.Diagnostics {
	m, err := mapFromUser(user)
	if err != nil {
		return diag.FromErr(err)
	}
	for key, value := range m {
		switch key {
		case "pk":
			data.SetId(strconv.Itoa(value.(int)))
			helpers.SetWrapper(data, key, value.(int))
		case "is_active", "is_superuser":
			helpers.SetWrapper(data, key, value.(bool))
		case "groups":
			helpers.SetWrapper(data, key, value.([]string))
		default:
			helpers.SetWrapper(data, key, value.(string))
		}
	}
	return diag.Diagnostics{}
}

func dataSourceUserReadByPk(ctx context.Context, d *schema.ResourceData, c *APIClient, pk int) diag.Diagnostics {
	req := c.client.CoreApi.CoreUsersRetrieve(ctx, int32(pk))

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	return setUser(d, *res)
}

func dataSourceUserReadByUsername(ctx context.Context, d *schema.ResourceData, c *APIClient, username string) diag.Diagnostics {
	req := c.client.CoreApi.CoreUsersList(ctx)
	req = req.Username(username)

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching users found")
	}

	if len(res.Results) > 1 {
		return diag.Errorf("Multiple users found")
	}

	return setUser(d, res.Results[0])
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	if n, ok := d.GetOk("pk"); ok {
		return dataSourceUserReadByPk(ctx, d, c, n.(int))
	}

	if n, ok := d.GetOk("username"); ok {
		return dataSourceUserReadByUsername(ctx, d, c, n.(string))
	}

	return diag.Errorf("Neither pk nor username were provided")
}
