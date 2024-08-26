package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePolicyGeoIP() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourcePolicyGeoIPCreate,
		ReadContext:   resourcePolicyGeoIPRead,
		UpdateContext: resourcePolicyGeoIPUpdate,
		DeleteContext: resourcePolicyGeoIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"countries": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: EnumToDescription(api.AllowedCountryCodeEnumEnumValues),
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: StringInEnum(api.AllowedCountryCodeEnumEnumValues),
				},
			},
		},
	}
}

func resourcePolicyGeoIPSchemaToProvider(d *schema.ResourceData) *api.GeoIPPolicyRequest {
	r := api.GeoIPPolicyRequest{
		Name:             d.Get("name").(string),
		ExecutionLogging: api.PtrBool(d.Get("execution_logging").(bool)),
	}

	asns := d.Get("asns").([]interface{})
	r.Asns = make([]int32, len(asns))
	for i, prov := range asns {
		r.Asns[i] = int32(prov.(int))
	}
	if a, ok := d.Get("countries").([]interface{}); ok {
		r.Countries = make([]api.CountryCodeEnum, 0)
		for _, c := range castSlice[string](a) {
			r.Countries = append(r.Countries, api.CountryCodeEnum(c))
		}
	}
	return &r
}

func resourcePolicyGeoIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyGeoIPSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesGeoipCreate(ctx).GeoIPPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyGeoIPRead(ctx, d, m)
}

func resourcePolicyGeoIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesGeoipRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "execution_logging", res.ExecutionLogging)
	if res.HasAsns() {
		localAsns := castSlice[int](d.Get("asns").([]interface{}))
		setWrapper(d, "asns", listConsistentMerge(localAsns, slice32ToInt(res.Asns)))
	}
	if res.Countries != nil {
		localCountries := make([]api.CountryCodeEnum, 0)
		for _, c := range castSlice[string](d.Get("countries").([]interface{})) {
			localCountries = append(localCountries, api.CountryCodeEnum(c))
		}
		setWrapper(d, "countries", listConsistentMerge(localCountries, res.Countries))
	}
	return diags
}

func resourcePolicyGeoIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyGeoIPSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesGeoipUpdate(ctx, d.Id()).GeoIPPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyGeoIPRead(ctx, d, m)
}

func resourcePolicyGeoIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesGeoipDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
