package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceRBACUserObjectPermission() *schema.Resource {
	return &schema.Resource{
		Description:   "RBAC --- ",
		CreateContext: schema.NoopContext,
		ReadContext:   schema.NoopContext,
		Delete:        schema.RemoveFromState,
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
