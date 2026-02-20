package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
				Default:  100,
			},
			"distance_tolerance_km": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  50,
			},
			"history_login_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"check_impossible_travel": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"impossible_tolerance_km": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"countries": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: helpers.EnumToDescription(api.AllowedCountryCodeEnumEnumValues),
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: helpers.StringInEnum(api.AllowedCountryCodeEnumEnumValues),
				},
			},
		},
	}
}

func resourcePolicyGeoIPSchemaToProvider(d *schema.ResourceData) *api.GeoIPPolicyRequest {
	r := api.GeoIPPolicyRequest{
		Name:                  d.Get("name").(string),
		ExecutionLogging:      new(d.Get("execution_logging").(bool)),
		CheckHistoryDistance:  helpers.GetP[bool](d, "check_history_distance"),
		HistoryMaxDistanceKm:  helpers.GetInt64P(d, "history_max_distance_km"),
		DistanceToleranceKm:   helpers.GetIntP(d, "distance_tolerance_km"),
		HistoryLoginCount:     helpers.GetIntP(d, "history_login_count"),
		CheckImpossibleTravel: helpers.GetP[bool](d, "check_impossible_travel"),
		ImpossibleToleranceKm: helpers.GetIntP(d, "impossible_tolerance_km"),
		Asns:                  helpers.CastSliceInt32(d.Get("asns").([]any)),
	}

	r.Countries = make([]api.CountryCodeEnum, 0)
	for _, c := range helpers.CastSlice[string](d, "countries") {
		r.Countries = append(r.Countries, api.CountryCodeEnum(c))
	}
	return &r
}

func resourcePolicyGeoIPCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyGeoIPSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesGeoipCreate(ctx).GeoIPPolicyRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyGeoIPRead(ctx, d, m)
}

func resourcePolicyGeoIPRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesGeoipRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "execution_logging", res.ExecutionLogging)
	helpers.SetWrapper(d, "check_history_distance", res.CheckHistoryDistance)
	helpers.SetWrapper(d, "history_max_distance_km", res.HistoryMaxDistanceKm)
	helpers.SetWrapper(d, "distance_tolerance_km", res.DistanceToleranceKm)
	helpers.SetWrapper(d, "history_login_count", res.HistoryLoginCount)
	helpers.SetWrapper(d, "check_impossible_travel", res.CheckImpossibleTravel)
	helpers.SetWrapper(d, "impossible_tolerance_km", res.ImpossibleToleranceKm)
	helpers.SetWrapper(d, "asns", helpers.ListConsistentMerge(
		helpers.CastSlice[int](d, "asns"),
		helpers.Slice32ToInt(res.Asns),
	))
	helpers.SetWrapper(d, "countries", helpers.ListConsistentMerge(
		helpers.CastSliceString[api.CountryCodeEnum](helpers.CastSlice[string](d, "countries")),
		res.Countries,
	))
	return diags
}

func resourcePolicyGeoIPUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyGeoIPSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesGeoipUpdate(ctx, d.Id()).GeoIPPolicyRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyGeoIPRead(ctx, d, m)
}

func resourcePolicyGeoIPDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesGeoipDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
