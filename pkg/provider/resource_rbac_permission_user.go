package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceRBACUserObjectPermission() *schema.Resource {
	return &schema.Resource{
		Description:   "RBAC --- ",
		CreateContext: resourceRBACUserObjectPermissionCreate,
		ReadContext:   resourceRBACUserObjectPermissionRead,
		DeleteContext: resourceRBACUserObjectPermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"model": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Description:      helpers.EnumToDescription(api.AllowedModelEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedModelEnumEnumValues),
			},
			"permission": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRBACUserObjectPermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceRBACUserObjectPermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceRBACUserObjectPermissionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceRBACUserObjectPermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}
