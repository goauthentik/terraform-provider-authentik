package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceObjectAttribute() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourceObjectAttributeCreate,
		ReadContext:   resourceObjectAttributeRead,
		UpdateContext: resourceObjectAttributeUpdate,
		DeleteContext: resourceObjectAttributeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"object_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Content-type of the objects this attribute definition applies to, e.g. `authentik_core.user`.",
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      helpers.EnumToDescription(api.AllowedObjectAttributeTypeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedObjectAttributeTypeEnumEnumValues),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Objects that are managed by authentik are created/updated automatically and may be overwritten by later migrations.",
			},
			"is_unique": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// Computed
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceObjectAttributeSchemaToModel(d *schema.ResourceData) *api.ObjectAttributeRequest {
	m := api.ObjectAttributeRequest{
		ObjectType: d.Get("object_type").(string),
		Enabled:    new(d.Get("enabled").(bool)),
		Key:        d.Get("key").(string),
		Label:      d.Get("label").(string),
		Type:       api.ObjectAttributeTypeEnum(d.Get("type").(string)),
		Regex:      helpers.GetP[string](d, "regex"),
		Group:      helpers.GetP[string](d, "group"),
		Managed:    *api.NewNullableString(helpers.GetP[string](d, "managed")),
		IsUnique:   helpers.GetP[bool](d, "is_unique"),
		IsRequired: helpers.GetP[bool](d, "is_required"),
	}
	return &m
}

func resourceObjectAttributeCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceObjectAttributeSchemaToModel(d)

	res, hr, err := c.client.CoreAPI.CoreObjectAttributesCreate(ctx).ObjectAttributeRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceObjectAttributeRead(ctx, d, m)
}

func resourceObjectAttributeRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CoreAPI.CoreObjectAttributesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	// res.ObjectType is a human-readable display string (e.g. "authentik Core | User");
	// the machine-readable "app_label.model" form accepted on write round-trips via ObjectTypeObj.
	helpers.SetWrapper(d, "object_type", res.ObjectTypeObj.FullyQualifiedModel)
	helpers.SetWrapper(d, "key", res.Key)
	helpers.SetWrapper(d, "label", res.Label)
	helpers.SetWrapper(d, "type", res.Type)
	helpers.SetWrapper(d, "enabled", res.Enabled)
	helpers.SetWrapper(d, "regex", res.Regex)
	helpers.SetWrapper(d, "group", res.Group)
	helpers.SetWrapper(d, "managed", res.Managed.Get())
	helpers.SetWrapper(d, "is_unique", res.IsUnique)
	helpers.SetWrapper(d, "is_required", res.IsRequired)
	helpers.SetWrapper(d, "created", res.Created.Format(time.RFC3339))
	helpers.SetWrapper(d, "last_updated", res.LastUpdated.Format(time.RFC3339))
	return diags
}

func resourceObjectAttributeUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceObjectAttributeSchemaToModel(d)

	res, hr, err := c.client.CoreAPI.CoreObjectAttributesUpdate(ctx, d.Id()).ObjectAttributeRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceObjectAttributeRead(ctx, d, m)
}

func resourceObjectAttributeDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreAPI.CoreObjectAttributesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
