package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceToken() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceTokenCreate,
		ReadContext:   resourceTokenRead,
		UpdateContext: resourceTokenUpdate,
		DeleteContext: resourceTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Computed
			"key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"expires_in": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Meta
			"retrieve_key": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Actual
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"intent": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.INTENTENUM_API,
				Description:      EnumToDescription(api.AllowedIntentEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedIntentEnumEnumValues),
			},
			"expires": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTokenSchemaToModel(d *schema.ResourceData) (*api.TokenRequest, diag.Diagnostics) {
	m := api.TokenRequest{
		Identifier: d.Get("identifier").(string),
		User:       api.PtrInt32(int32(d.Get("user").(int))),
		Expiring:   api.PtrBool(d.Get("expiring").(bool)),
	}

	if l, ok := d.Get("description").(string); ok {
		m.Description = &l
	}

	if l, ok := d.Get("expires").(string); ok && l != "" {
		t, err := time.Parse(time.RFC3339, l)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		m.Expires.Set(&t)
	}
	int := api.IntentEnum(d.Get("intent").(string))
	m.Intent = &int
	return &m, nil
}

func resourceTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceTokenSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreTokensCreate(ctx).TokenRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Identifier)
	return resourceTokenRead(ctx, d, m)
}

func resourceTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreTokensRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "identifier", res.Identifier)
	setWrapper(d, "user", res.User)
	setWrapper(d, "description", res.Description)
	setWrapper(d, "intent", res.Intent)
	if res.Expires.IsSet() && res.Expires.Get() != nil {
		setWrapper(d, "expires_in", time.Until(*res.Expires.Get()).Seconds())
	}
	if rt, ok := d.Get("retrieve_key").(bool); ok && rt {
		res, hr, err := c.client.CoreApi.CoreTokensViewKeyRetrieve(ctx, d.Id()).Execute()
		if err != nil {
			return httpToDiag(d, hr, err)
		}
		setWrapper(d, "key", res.Key)
	}
	return diags
}

func resourceTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceTokenSchemaToModel(d)
	if di != nil {
		return di
	}
	res, hr, err := c.client.CoreApi.CoreTokensUpdate(ctx, d.Id()).TokenRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Identifier)
	return resourceTokenRead(ctx, d, m)
}

func resourceTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreTokensDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
