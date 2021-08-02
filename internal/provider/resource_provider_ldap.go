package provider

import (
	"context"
	"strconv"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProviderLDAP() *schema.Resource {
	return &schema.Resource{
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
			"authorization_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_group": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceProviderLDAPSchemaToProvider(d *schema.ResourceData) (*api.LDAPProviderRequest, diag.Diagnostics) {
	r := api.LDAPProviderRequest{
		Name:              d.Get("name").(string),
		AuthorizationFlow: d.Get("authorization_flow").(string),
		BaseDn:            stringToPointer(d.Get("base_dn").(string)),
		UidStartNumber:    intToPointer(d.Get("uid_start_number").(int)),
		GidStartNumber:    intToPointer(d.Get("gid_start_number").(int)),
	}

	if s, sok := d.GetOk("search_group"); sok && s.(string) != "" {
		r.SearchGroup.Set(stringToPointer(s.(string)))
	}
	if s, sok := d.GetOk("certificate"); sok && s.(string) != "" {
		r.Certificate.Set(stringToPointer(s.(string)))
	}
	if s, sok := d.GetOk("tls_server_name"); sok && s.(string) != "" {
		r.TlsServerName = stringToPointer(s.(string))
	}
	return &r, nil
}

func resourceProviderLDAPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceProviderLDAPSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersLdapCreate(ctx).LDAPProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderLDAPRead(ctx, d, m)
}

func resourceProviderLDAPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersLdapRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("authorization_flow", res.AuthorizationFlow)
	d.Set("base_dn", res.BaseDn)
	if res.SearchGroup.IsSet() {
		d.Set("search_group", res.SearchGroup.Get())
	}
	if res.Certificate.IsSet() {
		d.Set("certificate", res.Certificate.Get())
	}
	d.Set("tls_server_name", res.TlsServerName)
	d.Set("uid_start_number", res.UidStartNumber)
	d.Set("gid_start_number", res.GidStartNumber)
	return diags
}

func resourceProviderLDAPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	app, di := resourceProviderLDAPSchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.ProvidersApi.ProvidersLdapUpdate(ctx, int32(id)).LDAPProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderLDAPRead(ctx, d, m)
}

func resourceProviderLDAPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersLdapDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
