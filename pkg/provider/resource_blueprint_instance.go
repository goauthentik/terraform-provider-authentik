package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				Description:      JSONDescription,
				DiffSuppressFunc: diffSuppressJSON,
				ValidateDiagFunc: ValidateJSON,
			},
		},
	}
}

func resourceBlueprintInstanceSchemaToModel(d *schema.ResourceData) (*api.BlueprintInstanceRequest, diag.Diagnostics) {
	m := api.BlueprintInstanceRequest{
		Name:    d.Get("name").(string),
		Enabled: api.PtrBool(d.Get("enabled").(bool)),
		Path:    getP[string](d.Get("path")),
		Content: getP[string](d.Get("content")),
	}

	context, err := getJSON[map[string]interface{}](d.Get("context"))
	m.Context = context
	return &m, err
}

func resourceBlueprintInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceBlueprintInstanceSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ManagedApi.ManagedBlueprintsCreate(ctx).BlueprintInstanceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)

	return resourceBlueprintInstanceRead(ctx, d, m)
}

func resourceBlueprintInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.ManagedApi.ManagedBlueprintsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "path", res.Path)
	setWrapper(d, "content", res.Content)
	setWrapper(d, "enabled", res.Enabled)
	b, err := json.Marshal(res.Context)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "context", string(b))
	return diags
}

func resourceBlueprintInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceBlueprintInstanceSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ManagedApi.ManagedBlueprintsUpdate(ctx, d.Id()).BlueprintInstanceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceBlueprintInstanceRead(ctx, d, m)
}

func resourceBlueprintInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.ManagedApi.ManagedBlueprintsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
