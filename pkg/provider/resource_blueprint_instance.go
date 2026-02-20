package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceBlueprintInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Blueprints --- ",
		CreateContext: resourceBlueprintInstanceCreate,
		ReadContext:   resourceBlueprintInstanceRead,
		UpdateContext: resourceBlueprintInstanceUpdate,
		DeleteContext: resourceBlueprintInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"context": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
		},
	}
}

func resourceBlueprintInstanceSchemaToModel(d *schema.ResourceData) (*api.BlueprintInstanceRequest, diag.Diagnostics) {
	m := api.BlueprintInstanceRequest{
		Name:    d.Get("name").(string),
		Enabled: new(d.Get("enabled").(bool)),
		Path:    helpers.GetP[string](d, "path"),
		Content: helpers.GetP[string](d, "content"),
	}

	context, err := helpers.GetJSON[map[string]any](d, ("context"))
	m.Context = context
	return &m, err
}

func resourceBlueprintInstanceCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceBlueprintInstanceSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ManagedApi.ManagedBlueprintsCreate(ctx).BlueprintInstanceRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)

	return resourceBlueprintInstanceRead(ctx, d, m)
}

func resourceBlueprintInstanceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.ManagedApi.ManagedBlueprintsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "path", res.Path)
	helpers.SetWrapper(d, "content", res.Content)
	helpers.SetWrapper(d, "enabled", res.Enabled)
	return helpers.SetJSON(d, "context", res.Context)
}

func resourceBlueprintInstanceUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceBlueprintInstanceSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ManagedApi.ManagedBlueprintsUpdate(ctx, d.Id()).BlueprintInstanceRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceBlueprintInstanceRead(ctx, d, m)
}

func resourceBlueprintInstanceDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.ManagedApi.ManagedBlueprintsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
