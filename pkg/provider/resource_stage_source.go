package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageSource() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageSourceCreate,
		ReadContext:   resourceStageSourceRead,
		UpdateContext: resourceStageSourceUpdate,
		DeleteContext: resourceStageSourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resume_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=10",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
		},
	}
}

func resourceStageSourceSchemaToProvider(d *schema.ResourceData) *api.SourceStageRequest {
	r := api.SourceStageRequest{
		Name:          d.Get("name").(string),
		Source:        d.Get("source").(string),
		ResumeTimeout: new(d.Get("resume_timeout").(string)),
	}
	return &r
}

func resourceStageSourceCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageSourceSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesSourceCreate(ctx).SourceStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageSourceRead(ctx, d, m)
}

func resourceStageSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesSourceRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "source", res.Source)
	helpers.SetWrapper(d, "resume_timeout", res.ResumeTimeout)
	return diags
}

func resourceStageSourceUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageSourceSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesSourceUpdate(ctx, d.Id()).SourceStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageSourceRead(ctx, d, m)
}

func resourceStageSourceDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesSourceDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
