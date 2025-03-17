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
			"asns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"check_history_distance": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"history_max_distance_km": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"distance_tolerance_km": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"history_login_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"check_impossible_travel": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"impossible_tolerance_km": {
				Type:     schema.TypeInt,
				Optional: true,
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
	if x, xok := d.GetOk("check_history_distance"); xok {
		r.CheckHistoryDistance = api.PtrBool(x.(bool))
	}
	if x, xok := d.GetOk("history_max_distance_km"); xok {
		r.HistoryMaxDistanceKm = api.PtrInt64(int64(x.(int)))
	}
	if x, xok := d.GetOk("distance_tolerance_km"); xok {
		r.DistanceToleranceKm = api.PtrInt32(int32(x.(int)))
	}
	if x, xok := d.GetOk("history_login_count"); xok {
		r.HistoryLoginCount = api.PtrInt32(int32(x.(int)))
	}
	if x, xok := d.GetOk("check_impossible_travel"); xok {
		r.CheckImpossibleTravel = api.PtrBool(x.(bool))
	}
	if x, xok := d.GetOk("impossible_tolerance_km"); xok {
		r.ImpossibleToleranceKm = api.PtrInt32(int32(x.(int)))
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
	setWrapper(d, "check_history_distance", res.CheckHistoryDistance)
	setWrapper(d, "history_max_distance_km", res.HistoryMaxDistanceKm)
	setWrapper(d, "distance_tolerance_km", res.DistanceToleranceKm)
	setWrapper(d, "history_login_count", res.HistoryLoginCount)
	setWrapper(d, "check_impossible_travel", res.CheckImpossibleTravel)
	setWrapper(d, "impossible_tolerance_km", res.ImpossibleToleranceKm)
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
