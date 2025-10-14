package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/helpers"
)

func resourceStageIdentification() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
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
					Type:             schema.TypeString,
					Description:      helpers.EnumToDescription(api.AllowedUserFieldsEnumEnumValues),
					ValidateDiagFunc: helpers.StringInEnum(api.AllowedUserFieldsEnumEnumValues),
				},
			},
			"password_stage": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"captcha_stage": {
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
			"pretend_user_exists": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_source_labels": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_remember_me": {
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
		PretendUserExists:       api.PtrBool(d.Get("pretend_user_exists").(bool)),
		ShowMatchedUser:         api.PtrBool(d.Get("show_matched_user").(bool)),
		EnableRememberMe:        api.PtrBool(d.Get("enable_remember_me").(bool)),
		ShowSourceLabels:        api.PtrBool(d.Get("show_source_labels").(bool)),
		CaseInsensitiveMatching: api.PtrBool(d.Get("case_insensitive_matching").(bool)),
		Sources:                 helpers.CastSlice_New[string](d, "sources"),
		PasswordStage:           *api.NewNullableString(api.PtrString(d.Get("password_stage").(string))),
		CaptchaStage:            *api.NewNullableString(api.PtrString(d.Get("captcha_stage").(string))),
		EnrollmentFlow:          *api.NewNullableString(helpers.GetP[string](d, "enrollment_flow")),
		RecoveryFlow:            *api.NewNullableString(helpers.GetP[string](d, "recovery_flow")),
		PasswordlessFlow:        *api.NewNullableString(helpers.GetP[string](d, "passwordless_flow")),
	}

	userFields := make([]api.UserFieldsEnum, 0)
	for _, userFieldsS := range d.Get("user_fields").([]interface{}) {
		userFields = append(userFields, api.UserFieldsEnum(userFieldsS.(string)))
	}
	r.UserFields = userFields
	return &r
}

func resourceStageIdentificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageIdentificationSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesIdentificationCreate(ctx).IdentificationStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageIdentificationRead(ctx, d, m)
}

func resourceStageIdentificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesIdentificationRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "user_fields", res.UserFields)
	if res.PasswordStage.IsSet() {
		helpers.SetWrapper(d, "password_stage", res.PasswordStage.Get())
	}
	if res.CaptchaStage.IsSet() {
		helpers.SetWrapper(d, "captcha_stage", res.CaptchaStage.Get())
	}
	helpers.SetWrapper(d, "case_insensitive_matching", res.CaseInsensitiveMatching)
	helpers.SetWrapper(d, "show_matched_user", res.ShowMatchedUser)
	helpers.SetWrapper(d, "enable_remember_me", res.EnableRememberMe)
	helpers.SetWrapper(d, "show_source_labels", res.ShowSourceLabels)
	helpers.SetWrapper(d, "pretend_user_exists", res.PretendUserExists)
	if res.EnrollmentFlow.IsSet() {
		helpers.SetWrapper(d, "enrollment_flow", res.EnrollmentFlow.Get())
	}
	if res.RecoveryFlow.IsSet() {
		helpers.SetWrapper(d, "recovery_flow", res.RecoveryFlow.Get())
	}
	if res.PasswordlessFlow.IsSet() {
		helpers.SetWrapper(d, "passwordless_flow", res.PasswordlessFlow.Get())
	}
	localSources := helpers.CastSlice_New[string](d, "sources")
	helpers.SetWrapper(d, "sources", helpers.ListConsistentMerge(localSources, res.Sources))
	return diags
}

func resourceStageIdentificationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageIdentificationSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesIdentificationUpdate(ctx, d.Id()).IdentificationStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageIdentificationRead(ctx, d, m)
}

func resourceStageIdentificationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesIdentificationDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
