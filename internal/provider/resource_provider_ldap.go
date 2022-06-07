package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
			"bind_flow": {
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
		},
	}
}

func resourceProviderLDAPSchemaToProvider(d *schema.ResourceData) *api.LDAPProviderRequest {
	r := api.LDAPProviderRequest{
		Name:              d.Get("name").(string),
		AuthorizationFlow: d.Get("bind_flow").(string),
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
	r.SetSearchMode(api.LDAPAPIAccessMode(d.Get("search_mode").(string)))
	r.SetBindMode(api.LDAPAPIAccessMode(d.Get("bind_mode").(string)))
	return &r
}

func resourceProviderLDAPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderLDAPSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersLdapCreate(ctx).LDAPProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
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
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "bind_flow", res.AuthorizationFlow)
	setWrapper(d, "base_dn", res.BaseDn)
	if res.SearchGroup.IsSet() {
		setWrapper(d, "search_group", res.SearchGroup.Get())
	}
	if res.Certificate.IsSet() {
		setWrapper(d, "certificate", res.Certificate.Get())
	}
	setWrapper(d, "tls_server_name", res.TlsServerName)
	setWrapper(d, "uid_start_number", res.UidStartNumber)
	setWrapper(d, "gid_start_number", res.GidStartNumber)
	setWrapper(d, "bind_mode", res.BindMode)
	setWrapper(d, "search_mode", res.SearchMode)
	return diags
}

func resourceProviderLDAPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderLDAPSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersLdapUpdate(ctx, int32(id)).LDAPProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
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
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
