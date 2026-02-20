package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceProviderOAuth2Config() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProviderOAuth2ConfigRead,
		Description: "Applications --- Get OAuth2 provider config",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"provider_id"},
				Description:   "Find provider by name",
			},
			"provider_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "Find provider by ID",
			},

			"issuer_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorize_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_info_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_info_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logout_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"jwks_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceProviderOAuth2ConfigRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	id, ok := d.GetOk("provider_id")
	if !ok {
		req := c.client.ProvidersApi.ProvidersOauth2List(ctx)
		if m, ok := d.Get("name").(string); ok {
			req = req.Name(m)
		}
		res, hr, err := req.Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		if len(res.Results) < 1 {
			return diag.Errorf("no matching providers found")
		}
		id = int(res.Results[0].Pk)
	}
	finalId := int32(id.(int))
	d.SetId(strconv.FormatInt(int64(finalId), 10))

	meta, hr, err := c.client.ProvidersApi.ProvidersOauth2SetupUrlsRetrieve(ctx, finalId).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	helpers.SetWrapper(d, "issuer_url", meta.Issuer)
	helpers.SetWrapper(d, "authorize_url", meta.Authorize)
	helpers.SetWrapper(d, "token_url", meta.Token)
	helpers.SetWrapper(d, "user_info_url", meta.UserInfo)
	helpers.SetWrapper(d, "provider_info_url", meta.ProviderInfo)
	helpers.SetWrapper(d, "logout_url", meta.Logout)
	helpers.SetWrapper(d, "jwks_url", meta.Jwks)
	return diags
}
