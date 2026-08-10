package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ resource.Resource                = &policyGeoIPResource{}
	_ resource.ResourceWithConfigure   = &policyGeoIPResource{}
	_ resource.ResourceWithImportState = &policyGeoIPResource{}
)

func newPolicyGeoIPResource() resource.Resource {
	return &policyGeoIPResource{}
}

type policyGeoIPResource struct {
	resourceBase
}

type policyGeoIPModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	ExecutionLogging      types.Bool   `tfsdk:"execution_logging"`
	Asns                  types.List   `tfsdk:"asns"`
	CheckHistoryDistance  types.Bool   `tfsdk:"check_history_distance"`
	HistoryMaxDistanceKm  types.Int64  `tfsdk:"history_max_distance_km"`
	DistanceToleranceKm   types.Int32  `tfsdk:"distance_tolerance_km"`
	HistoryLoginCount     types.Int32  `tfsdk:"history_login_count"`
	CheckImpossibleTravel types.Bool   `tfsdk:"check_impossible_travel"`
	ImpossibleToleranceKm types.Int32  `tfsdk:"impossible_tolerance_km"`
	Countries             types.List   `tfsdk:"countries"`
}

func (r *policyGeoIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_geoip"
}

func (r *policyGeoIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	executionLoggingDefault := helpers.BoolDefault(false)
	// history_max_distance_km is the provider's only int64 API field, hence Int64Attribute
	// and Int64Default rather than the usual Int32 pair. Both serialise to protocol
	// `number`, so state and docs are unchanged.
	historyMaxDistanceKmDefault := helpers.Int64Default(100)
	distanceToleranceKmDefault := helpers.Int32Default(50)
	historyLoginCountDefault := helpers.Int32Default(5)
	impossibleToleranceKmDefault := helpers.Int32Default(100)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Customization --- ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"execution_logging": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             executionLoggingDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(executionLoggingDefault.Value())),
			},
			"asns": schema.ListAttribute{
				ElementType: types.Int32Type,
				Optional:    true,
			},
			// check_history_distance and check_impossible_travel are the two attributes
			// this whole file is about: Optional with no Default, so they stay plain
			// Optional booleans and are the only two using helpers.BoolPtr on the write
			// side. See toRequest for why that matters.
			"check_history_distance": schema.BoolAttribute{
				Optional: true,
			},
			"history_max_distance_km": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             historyMaxDistanceKmDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(historyMaxDistanceKmDefault.Value())),
			},
			"distance_tolerance_km": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             distanceToleranceKmDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(distanceToleranceKmDefault.Value())),
			},
			"history_login_count": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             historyLoginCountDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(historyLoginCountDefault.Value())),
			},
			"check_impossible_travel": schema.BoolAttribute{
				Optional: true,
			},
			"impossible_tolerance_km": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             impossibleToleranceKmDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(impossibleToleranceKmDefault.Value())),
			},
			// Unlike the enum lists in the stages batch, SDKv2 declared this enum prose on
			// the *list* rather than on its Elem, so it does reach the generated docs and
			// must be kept here.
			"countries": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: helpers.EnumToDescription(api.AllowedCountryCodeEnumEnumValues),
				Validators: []validator.List{
					listvalidator.ValueStringsAre(helpers.OneOf(api.AllowedCountryCodeEnumEnumValues)),
				},
			},
		},
	}
}

func (r *policyGeoIPResource) toRequest(ctx context.Context, data *policyGeoIPModel) (*api.GeoIPPolicyRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	asns, d := helpers.SliceOrEmpty[int32](ctx, data.Asns)
	diags.Append(d...)

	countries, d := helpers.SliceOrEmpty[string](ctx, data.Countries)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.GeoIPPolicyRequest{
		Name:             data.Name.ValueString(),
		ExecutionLogging: new(data.ExecutionLogging.ValueBool()),
		Asns:             asns,
		Countries:        helpers.CastSliceString[api.CountryCodeEnum](countries),
		// These two are the fix migration-plan.md flags as "where the null-semantics change
		// is user-visible". SDKv2 built them with GetP[bool], which is d.GetOk underneath
		// and so returned nil for *any* false value - it could not tell "unset" from
		// "explicitly false". With omitempty on the request field and the API leaving
		// absent fields unchanged on PUT, that meant either check could be switched on but
		// never back off: the request omitted it, the API kept true, Read returned true,
		// and every subsequent plan showed the same diff. BoolPtr sends an explicit false
		// and omits only a genuinely null config.
		CheckHistoryDistance:  helpers.BoolPtr(data.CheckHistoryDistance),
		CheckImpossibleTravel: helpers.BoolPtr(data.CheckImpossibleTravel),
		HistoryMaxDistanceKm:  helpers.Int64Ptr(data.HistoryMaxDistanceKm),
		DistanceToleranceKm:   helpers.Int32Ptr(data.DistanceToleranceKm),
		HistoryLoginCount:     helpers.Int32Ptr(data.HistoryLoginCount),
		ImpossibleToleranceKm: helpers.Int32Ptr(data.ImpossibleToleranceKm),
	}, diags
}

func (r *policyGeoIPResource) fromAPI(ctx context.Context, data *policyGeoIPModel, res *api.GeoIPPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)

	asns, d := helpers.MergeInt32List(ctx, data.Asns, res.Asns)
	diags.Append(d...)
	data.Asns = asns

	countries := make([]string, len(res.Countries))
	for i, c := range res.Countries {
		countries[i] = string(c)
	}
	mergedCountries, d := helpers.MergeStringList(ctx, data.Countries, countries)
	diags.Append(d...)
	data.Countries = mergedCountries

	// The two no-Default booleans keep the prior-aware helper, so a never-configured check
	// stays null rather than becoming false.
	data.CheckHistoryDistance = helpers.BoolOrNull(data.CheckHistoryDistance, res.GetCheckHistoryDistance())
	data.CheckImpossibleTravel = helpers.BoolOrNull(data.CheckImpossibleTravel, res.GetCheckImpossibleTravel())
	// Everything below has a Default and takes the API value verbatim.
	data.ExecutionLogging = types.BoolValue(res.GetExecutionLogging())
	data.HistoryMaxDistanceKm = types.Int64Value(res.GetHistoryMaxDistanceKm())
	data.DistanceToleranceKm = types.Int32Value(res.GetDistanceToleranceKm())
	data.HistoryLoginCount = types.Int32Value(res.GetHistoryLoginCount())
	data.ImpossibleToleranceKm = types.Int32Value(res.GetImpossibleToleranceKm())

	return diags
}

func (r *policyGeoIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data policyGeoIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesGeoipCreate(ctx).GeoIPPolicyRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyGeoIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data policyGeoIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesGeoipRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyGeoIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data policyGeoIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesGeoipUpdate(ctx, data.ID.ValueString()).GeoIPPolicyRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyGeoIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data policyGeoIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PoliciesAPI.PoliciesGeoipDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
