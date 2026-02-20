package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceProviderLDAP() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderLDAPCreate,
		ReadContext:   resourceProviderLDAPRead,
		UpdateContext: resourceProviderLDAPUpdate,
		DeleteContext: resourceProviderLDAPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"unbind_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_server_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uid_start_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2000,
			},
			"gid_start_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4000,
			},
			"search_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.LDAPAPIACCESSMODE_DIRECT,
			},
			"bind_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.LDAPAPIACCESSMODE_DIRECT,
			},
			"mfa_support": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceProviderLDAPSchemaToProvider(d *schema.ResourceData) *api.LDAPProviderRequest {
	r := api.LDAPProviderRequest{
		Name:              d.Get("name").(string),
		AuthorizationFlow: d.Get("bind_flow").(string),
		InvalidationFlow:  d.Get("unbind_flow").(string),
		BaseDn:            new(d.Get("base_dn").(string)),
		UidStartNumber:    new(int32(d.Get("uid_start_number").(int))),
		GidStartNumber:    new(int32(d.Get("gid_start_number").(int))),
		SearchMode:        api.LDAPAPIAccessMode(d.Get("search_mode").(string)).Ptr(),
		BindMode:          api.LDAPAPIAccessMode(d.Get("bind_mode").(string)).Ptr(),
		MfaSupport:        new(d.Get("mfa_support").(bool)),
		Certificate:       *api.NewNullableString(helpers.GetP[string](d, "certificate")),
		TlsServerName:     helpers.GetP[string](d, "tls_server_name"),
	}
	return &r
}

func resourceProviderLDAPCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderLDAPSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersLdapCreate(ctx).LDAPProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderLDAPRead(ctx, d, m)
}

func resourceProviderLDAPRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersLdapRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "bind_flow", res.AuthorizationFlow)
	helpers.SetWrapper(d, "unbind_flow", res.InvalidationFlow)
	helpers.SetWrapper(d, "base_dn", res.BaseDn)
	helpers.SetWrapper(d, "certificate", res.Certificate.Get())
	helpers.SetWrapper(d, "tls_server_name", res.TlsServerName)
	helpers.SetWrapper(d, "uid_start_number", res.UidStartNumber)
	helpers.SetWrapper(d, "gid_start_number", res.GidStartNumber)
	helpers.SetWrapper(d, "bind_mode", res.BindMode)
	helpers.SetWrapper(d, "search_mode", res.SearchMode)
	helpers.SetWrapper(d, "mfa_support", res.MfaSupport)
	return diags
}

func resourceProviderLDAPUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderLDAPSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersLdapUpdate(ctx, int32(id)).LDAPProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderLDAPRead(ctx, d, m)
}

func resourceProviderLDAPDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersLdapDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
