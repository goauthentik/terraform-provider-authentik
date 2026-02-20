package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceTaskSchedule() *schema.Resource {
	return &schema.Resource{
		Description:   "Tasks --- ",
		CreateContext: resourceTaskScheduleCreate,
		ReadContext:   resourceTaskScheduleRead,
		UpdateContext: resourceTaskScheduleUpdate,
		DeleteContext: resourceTaskScheduleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"app_model": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      helpers.EnumToDescription(api.AllowedModelEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedModelEnumEnumValues),
			},
			"model_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"actor_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"crontab": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Crontab expression at which this task will run.",
			},
			"paused": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourceTaskScheduleSchemaToProvider(d *schema.ResourceData) api.PatchedScheduleRequest {
	return api.PatchedScheduleRequest{
		RelObjId: *api.NewNullableString(helpers.GetP[string](d, "model_id")),
		Crontab:  new(d.Get("crontab").(string)),
		Paused:   new(d.Get("paused").(bool)),
	}
}

func resourceTaskScheduleImport(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	parts := strings.Split(d.Get("app_model").(string), ".")

	req := c.client.TasksApi.
		TasksSchedulesList(ctx).
		RelObjContentTypeAppLabel(parts[0]).
		RelObjContentTypeModel(parts[1]).
		RelObjId(d.Get("model_id").(string))
	if act, ok := d.GetOk("actor_name"); ok {
		req = req.ActorName(act.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	if len(res.Results) < 1 {
		return diag.Errorf("Failed to find task schedule")
	}

	d.SetId(res.Results[0].Id)
	return nil
}

func resourceTaskScheduleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	diag := resourceTaskScheduleImport(ctx, d, m)
	if diag != nil {
		return diag
	}
	diag = resourceTaskScheduleUpdate(ctx, d, m)
	if diag != nil {
		return diag
	}

	return resourceTaskScheduleRead(ctx, d, m)
}

func resourceTaskScheduleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.TasksApi.TasksSchedulesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "model_id", res.RelObjId.Get())
	helpers.SetWrapper(d, "crontab", res.Crontab)
	helpers.SetWrapper(d, "paused", res.Paused)
	helpers.SetWrapper(d, "actor_name", res.ActorName)
	return diags
}

func resourceTaskScheduleUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceTaskScheduleSchemaToProvider(d)

	res, hr, err := c.client.TasksApi.TasksSchedulesPartialUpdate(ctx, d.Id()).
		PatchedScheduleRequest(app).
		Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Id)
	return resourceTaskScheduleRead(ctx, d, m)
}

func resourceTaskScheduleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return diag.Diagnostics{}
}
