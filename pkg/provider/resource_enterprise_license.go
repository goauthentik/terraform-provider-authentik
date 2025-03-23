package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceEnterpriseLicense() *schema.Resource {
	return &schema.Resource{
		Description:   "Enterprise --- ",
		CreateContext: resourceEnterpriseLicenseCreate,
		ReadContext:   resourceEnterpriseLicenseRead,
		UpdateContext: resourceEnterpriseLicenseUpdate,
		DeleteContext: resourceEnterpriseLicenseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiry": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"internal_users": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"external_users": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceEnterpriseLicenseSchemaToProvider(d *schema.ResourceData) *api.LicenseRequest {
	r := api.LicenseRequest{
		Key: d.Get("key").(string),
	}
	return &r
}

func resourceEnterpriseLicenseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceEnterpriseLicenseSchemaToProvider(d)

	res, hr, err := c.client.EnterpriseApi.EnterpriseLicenseCreate(ctx).LicenseRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.LicenseUuid)
	return resourceEnterpriseLicenseRead(ctx, d, m)
}

func resourceEnterpriseLicenseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.EnterpriseApi.EnterpriseLicenseRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expiry", res.Expiry.Unix())
	setWrapper(d, "key", res.Key)
	setWrapper(d, "internal_users", res.InternalUsers)
	setWrapper(d, "external_users", res.ExternalUsers)
	return diags
}

func resourceEnterpriseLicenseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceEnterpriseLicenseSchemaToProvider(d)

	res, hr, err := c.client.EnterpriseApi.EnterpriseLicenseUpdate(ctx, d.Id()).LicenseRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.LicenseUuid)
	return resourceEnterpriseLicenseRead(ctx, d, m)
}

func resourceEnterpriseLicenseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.EnterpriseApi.EnterpriseLicenseDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
