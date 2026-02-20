package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStagePassword() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStagePasswordCreate,
		ReadContext:   resourceStagePasswordRead,
		UpdateContext: resourceStagePasswordUpdate,
		DeleteContext: resourceStagePasswordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backends": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					Description:      helpers.EnumToDescription(api.AllowedBackendsEnumEnumValues),
					ValidateDiagFunc: helpers.StringInEnum(api.AllowedBackendsEnumEnumValues),
				},
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"failed_attempts_before_cancel": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"allow_show_password": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourceStagePasswordSchemaToProvider(d *schema.ResourceData) *api.PasswordStageRequest {
	r := api.PasswordStageRequest{
		Name:                       d.Get("name").(string),
		AllowShowPassword:          new(d.Get("allow_show_password").(bool)),
		ConfigureFlow:              *api.NewNullableString(helpers.GetP[string](d, "configure_flow")),
		FailedAttemptsBeforeCancel: helpers.GetIntP(d, "failed_attempts_before_cancel"),
	}

	backend := make([]api.BackendsEnum, 0)
	for _, backendS := range d.Get("backends").([]any) {
		backend = append(backend, api.BackendsEnum(backendS.(string)))
	}
	r.Backends = backend
	return &r
}

func resourceStagePasswordCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStagePasswordSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPasswordCreate(ctx).PasswordStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePasswordRead(ctx, d, m)
}

func resourceStagePasswordRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesPasswordRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "backends", res.Backends)
	helpers.SetWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	helpers.SetWrapper(d, "failed_attempts_before_cancel", res.FailedAttemptsBeforeCancel)
	helpers.SetWrapper(d, "allow_show_password", res.AllowShowPassword)
	return diags
}

func resourceStagePasswordUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStagePasswordSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPasswordUpdate(ctx, d.Id()).PasswordStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePasswordRead(ctx, d, m)
}

func resourceStagePasswordDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesPasswordDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
