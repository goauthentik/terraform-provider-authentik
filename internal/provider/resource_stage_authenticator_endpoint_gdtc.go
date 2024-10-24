package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageAuthenticatorEndpointGDTC() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorEndpointGDTCCreate,
		ReadContext:   resourceStageAuthenticatorEndpointGDTCRead,
		UpdateContext: resourceStageAuthenticatorEndpointGDTCUpdate,
		DeleteContext: resourceStageAuthenticatorEndpointGDTCDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"credentials": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceStageAuthenticatorEndpointGDTCSchemaToProvider(d *schema.ResourceData) (*api.AuthenticatorEndpointGDTCStageRequest, error) {
	r := api.AuthenticatorEndpointGDTCStageRequest{
		Name: d.Get("name").(string),
	}
	var creds interface{}
	err := json.NewDecoder(strings.NewReader(d.Get("credentials").(string))).Decode(&creds)
	if err != nil {
		return nil, err
	}
	r.Credentials = creds
	return &r, nil
}

func resourceStageAuthenticatorEndpointGDTCCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, err := resourceStageAuthenticatorEndpointGDTCSchemaToProvider(d)
	if err != nil {
		return diag.FromErr(err)
	}

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEndpointGdtcCreate(ctx).AuthenticatorEndpointGDTCStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorEndpointGDTCRead(ctx, d, m)
}

func resourceStageAuthenticatorEndpointGDTCRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEndpointGdtcRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	b, err := json.Marshal(res.Credentials)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "credentials", string(b))
	return diags
}

func resourceStageAuthenticatorEndpointGDTCUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, err := resourceStageAuthenticatorEndpointGDTCSchemaToProvider(d)
	if err != nil {
		return diag.FromErr(err)
	}

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEndpointGdtcUpdate(ctx, d.Id()).AuthenticatorEndpointGDTCStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorEndpointGDTCRead(ctx, d, m)
}

func resourceStageAuthenticatorEndpointGDTCDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorEndpointGdtcDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
