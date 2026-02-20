package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceProviderMicrosoftEntra() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderMicrosoftEntraCreate,
		ReadContext:   resourceProviderMicrosoftEntraRead,
		UpdateContext: resourceProviderMicrosoftEntraUpdate,
		DeleteContext: resourceProviderMicrosoftEntraDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"property_mappings": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"property_mappings_group": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"exclude_users_service_account": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"filter_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_delete_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.OUTGOINGSYNCDELETEACTION_DELETE,
				Description: helpers.EnumToDescription([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
				ValidateDiagFunc: helpers.StringInEnum([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
			},
			"group_delete_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.OUTGOINGSYNCDELETEACTION_DELETE,
				Description: helpers.EnumToDescription([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
				ValidateDiagFunc: helpers.StringInEnum([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
			},
			"sync_page_timeout": {
				Type:             schema.TypeString,
				Default:          "minutes=30",
				Optional:         true,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
				Description:      helpers.RelativeDurationDescription,
			},
			"sync_page_size": {
				Type:     schema.TypeInt,
				Default:  100,
				Optional: true,
			},
		},
	}
}

func resourceProviderMicrosoftEntraSchemaToProvider(d *schema.ResourceData) (*api.MicrosoftEntraProviderRequest, diag.Diagnostics) {
	r := api.MicrosoftEntraProviderRequest{
		Name:                       d.Get("name").(string),
		ClientId:                   d.Get("client_id").(string),
		ClientSecret:               d.Get("client_secret").(string),
		TenantId:                   d.Get("tenant_id").(string),
		PropertyMappings:           helpers.CastSlice[string](d, "property_mappings"),
		PropertyMappingsGroup:      helpers.CastSlice[string](d, "property_mappings_group"),
		ExcludeUsersServiceAccount: new(d.Get("exclude_users_service_account").(bool)),
		UserDeleteAction:           api.OutgoingSyncDeleteAction(d.Get("user_delete_action").(string)).Ptr(),
		GroupDeleteAction:          api.OutgoingSyncDeleteAction(d.Get("group_delete_action").(string)).Ptr(),
		FilterGroup:                *api.NewNullableString(helpers.GetP[string](d, "filter_group")),
		DryRun:                     new(d.Get("dry_run").(bool)),
		SyncPageTimeout:            helpers.GetP[string](d, "sync_page_timeout"),
		SyncPageSize:               helpers.GetIntP(d, "sync_page_size"),
	}
	return &r, nil
}

func resourceProviderMicrosoftEntraCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceProviderMicrosoftEntraSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersMicrosoftEntraCreate(ctx).MicrosoftEntraProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderMicrosoftEntraRead(ctx, d, m)
}

func resourceProviderMicrosoftEntraRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersMicrosoftEntraRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "client_id", res.ClientId)
	helpers.SetWrapper(d, "client_secret", res.ClientSecret)
	helpers.SetWrapper(d, "tenant_id", res.TenantId)
	helpers.SetWrapper(d, "exclude_users_service_account", res.ExcludeUsersServiceAccount)
	helpers.SetWrapper(d, "user_delete_action", res.UserDeleteAction)
	helpers.SetWrapper(d, "group_delete_action", res.GroupDeleteAction)
	helpers.SetWrapper(d, "filter_group", res.FilterGroup)
	helpers.SetWrapper(d, "dry_run", res.DryRun)
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.PropertyMappings,
	))
	helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings_group"),
		res.PropertyMappingsGroup,
	))
	helpers.SetWrapper(d, "sync_page_timeout", res.SyncPageTimeout)
	helpers.SetWrapper(d, "sync_page_size", res.SyncPageSize)
	return diags
}

func resourceProviderMicrosoftEntraUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app, diags := resourceProviderMicrosoftEntraSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersMicrosoftEntraUpdate(ctx, int32(id)).MicrosoftEntraProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderMicrosoftEntraRead(ctx, d, m)
}

func resourceProviderMicrosoftEntraDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersMicrosoftEntraDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
