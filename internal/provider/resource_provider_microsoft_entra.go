package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				Description: EnumToDescription([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
				ValidateDiagFunc: StringInEnum([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
			},
			"group_delete_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.OUTGOINGSYNCDELETEACTION_DELETE,
				Description: EnumToDescription([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
				ValidateDiagFunc: StringInEnum([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
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
		PropertyMappings:           castSlice[string](d.Get("property_mappings").([]interface{})),
		PropertyMappingsGroup:      castSlice[string](d.Get("property_mappings_group").([]interface{})),
		ExcludeUsersServiceAccount: api.PtrBool(d.Get("exclude_users_service_account").(bool)),
		UserDeleteAction:           api.OutgoingSyncDeleteAction(d.Get("user_delete_action").(string)).Ptr(),
		GroupDeleteAction:          api.OutgoingSyncDeleteAction(d.Get("group_delete_action").(string)).Ptr(),
	}

	if l, ok := d.Get("filter_group").(string); ok {
		r.FilterGroup = *api.NewNullableString(&l)
	}
	return &r, nil
}

func resourceProviderMicrosoftEntraCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceProviderMicrosoftEntraSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersMicrosoftEntraCreate(ctx).MicrosoftEntraProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderMicrosoftEntraRead(ctx, d, m)
}

func resourceProviderMicrosoftEntraRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersMicrosoftEntraRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "client_id", res.ClientId)
	setWrapper(d, "client_secret", res.ClientSecret)
	setWrapper(d, "tenant_id", res.TenantId)
	setWrapper(d, "exclude_users_service_account", res.ExcludeUsersServiceAccount)
	setWrapper(d, "user_delete_action", res.UserDeleteAction)
	setWrapper(d, "group_delete_action", res.GroupDeleteAction)
	setWrapper(d, "filter_group", res.FilterGroup)
	localMappings := castSlice[string](d.Get("property_mappings").([]interface{}))
	if len(localMappings) > 0 {
		setWrapper(d, "property_mappings", listConsistentMerge(localMappings, res.PropertyMappings))
	}
	localGroupMappings := castSlice[string](d.Get("property_mappings_group").([]interface{}))
	if len(localGroupMappings) > 0 {
		setWrapper(d, "property_mappings_group", listConsistentMerge(localGroupMappings, res.PropertyMappingsGroup))
	}
	return diags
}

func resourceProviderMicrosoftEntraUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderMicrosoftEntraRead(ctx, d, m)
}

func resourceProviderMicrosoftEntraDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersMicrosoftEntraDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
