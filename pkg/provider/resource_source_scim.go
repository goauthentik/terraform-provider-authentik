package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceSourceSCIM() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceSourceSCIMCreate,
		ReadContext:   resourceSourceSCIMRead,
		UpdateContext: resourceSourceSCIMUpdate,
		DeleteContext: resourceSourceSCIMDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_path_template": {
				Type:     schema.TypeString,
				Default:  "goauthentik.io/sources/%(slug)s",
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"property_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"property_mappings_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"scim_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCIM URL",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCIM URL",
			},
		},
	}
}

func resourceSourceSCIMSchemaToSource(d *schema.ResourceData) *api.SCIMSourceRequest {
	r := api.SCIMSourceRequest{
		Name:             d.Get("name").(string),
		Slug:             d.Get("slug").(string),
		Enabled:          new(d.Get("enabled").(bool)),
		UserPathTemplate: new(d.Get("user_path_template").(string)),

		UserPropertyMappings:  helpers.CastSlice[string](d, "property_mappings"),
		GroupPropertyMappings: helpers.CastSlice[string](d, "property_mappings_group"),
	}
	return &r
}

func resourceSourceSCIMCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceSCIMSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesScimCreate(ctx).SCIMSourceRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSCIMRead(ctx, d, m)
}

func resourceSourceSCIMRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesScimRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "slug", res.Slug)
	helpers.SetWrapper(d, "uuid", res.Pk)
	helpers.SetWrapper(d, "user_path_template", res.UserPathTemplate)
	helpers.SetWrapper(d, "enabled", res.Enabled)
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.UserPropertyMappings,
	))
	helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings_group"),
		res.GroupPropertyMappings,
	))

	meta, hr, err := c.client.SourcesApi.SourcesScimRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	helpers.SetWrapper(d, "scim_url", meta.RootUrl)
	helpers.SetWrapper(d, "token", meta.TokenObj.Identifier)
	return diags
}

func resourceSourceSCIMUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceSCIMSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesScimUpdate(ctx, d.Id()).SCIMSourceRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSCIMRead(ctx, d, m)
}

func resourceSourceSCIMDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesScimDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
