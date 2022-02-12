package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
)

func resourceStageIdentification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStageIdentificationCreate,
		ReadContext:   resourceStageIdentificationRead,
		UpdateContext: resourceStageIdentificationUpdate,
		DeleteContext: resourceStageIdentificationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"password_stage": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"case_insensitive_matching": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"show_matched_user": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_source_labels": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enrollment_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"recovery_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"passwordless_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceStageIdentificationSchemaToProvider(d *schema.ResourceData) *api.IdentificationStageRequest {
	r := api.IdentificationStageRequest{
		Name:                    d.Get("name").(string),
		ShowMatchedUser:         boolToPointer(d.Get("show_matched_user").(bool)),
		ShowSourceLabels:        boolToPointer(d.Get("show_source_labels").(bool)),
		CaseInsensitiveMatching: boolToPointer(d.Get("case_insensitive_matching").(bool)),
	}

	if h, hSet := d.GetOk("password_stage"); hSet {
		r.PasswordStage.Set(stringToPointer(h.(string)))
	}
	if h, hSet := d.GetOk("enrollment_flow"); hSet {
		r.EnrollmentFlow.Set(stringToPointer(h.(string)))
	}
	if h, hSet := d.GetOk("recovery_flow"); hSet {
		r.RecoveryFlow.Set(stringToPointer(h.(string)))
	}
	if h, hSet := d.GetOk("passwordless_flow"); hSet {
		r.PasswordlessFlow.Set(stringToPointer(h.(string)))
	}

	userFields := make([]api.UserFieldsEnum, 0)
	for _, userFieldsS := range d.Get("user_fields").([]interface{}) {
		userFields = append(userFields, api.UserFieldsEnum(userFieldsS.(string)))
	}
	r.UserFields = &userFields

	sources := make([]string, 0)
	for _, sourcesS := range d.Get("sources").([]interface{}) {
		sources = append(sources, sourcesS.(string))
	}
	if len(sources) > 0 {
		r.Sources = &sources
	}

	return &r
}

func resourceStageIdentificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageIdentificationSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesIdentificationCreate(ctx).IdentificationStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageIdentificationRead(ctx, d, m)
}

func resourceStageIdentificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesIdentificationRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.Set("name", res.Name)
	d.Set("user_fields", res.UserFields)
	if res.PasswordStage.IsSet() {
		d.Set("password_stage", res.PasswordStage.Get())
	}
	d.Set("case_insensitive_matching", res.CaseInsensitiveMatching)
	d.Set("show_matched_user", res.ShowMatchedUser)
	d.Set("show_source_labels", res.ShowSourceLabels)
	if res.EnrollmentFlow.IsSet() {
		d.Set("enrollment_flow", res.EnrollmentFlow.Get())
	}
	if res.RecoveryFlow.IsSet() {
		d.Set("recovery_flow", res.RecoveryFlow.Get())
	}
	if res.PasswordlessFlow.IsSet() {
		d.Set("passwordless_flow", res.PasswordlessFlow.Get())
	}
	d.Set("sources", res.Sources)
	return diags
}

func resourceStageIdentificationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageIdentificationSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesIdentificationUpdate(ctx, d.Id()).IdentificationStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageIdentificationRead(ctx, d, m)
}

func resourceStageIdentificationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesIdentificationDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
