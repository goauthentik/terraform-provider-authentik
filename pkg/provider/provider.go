// Package provider implements the terraform-plugin-framework half of the authentik
// provider. It is muxed together with pkg/sdkprovider (the original SDKv2 provider) via
// terraform-plugin-mux; resources and data sources move from pkg/sdkprovider to here in
// batches as they are migrated. See migration-plan.md at the repo root.
package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var _ provider.Provider = &authentikProvider{}

type authentikProvider struct {
	version string
	testing bool
}

// New returns the framework half of the muxed authentik provider. testing mirrors
// pkg/sdkprovider.Provider's testing flag: when true, all HTTP requests fail with a
// mocked 400 response instead of reaching a real server.
func New(version string, testing bool) provider.Provider {
	return &authentikProvider{
		version: version,
		testing: testing,
	}
}

type authentikProviderModel struct {
	URL      types.String `tfsdk:"url"`
	Insecure types.Bool   `tfsdk:"insecure"`
	Token    types.String `tfsdk:"token"`
	Headers  types.Map    `tfsdk:"headers"`
}

func (p *authentikProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "authentik"
	resp.Version = p.version
}

// Schema must stay byte-for-byte identical to pkg/sdkprovider's, because
// terraform-plugin-mux compares both servers' tfprotov6.Schema for the provider block
// and errors on any mismatch. MarkdownDescription (not Description) is used throughout
// because the SDKv2 provider sets the global schema.DescriptionKind = StringMarkdown, so
// its descriptions are emitted as MARKDOWN; Description would emit PLAIN and mismatch.
// The provider block's own top-level description is intentionally left empty to match
// the SDKv2 side, whose provider.Description is never set.
func (p *authentikProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The authentik API endpoint, can optionally be passed as `AUTHENTIK_URL` environmental variable",
			},
			"insecure": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to skip TLS verification, can optionally be passed as `AUTHENTIK_INSECURE` environmental variable",
			},
			"token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The authentik API token, can optionally be passed as `AUTHENTIK_TOKEN` environmental variable",
			},
			"headers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Optional HTTP headers sent with every request",
			},
		},
	}
}

// Configure mirrors pkg/sdkprovider's providerConfigure exactly, including the env
// fallback being keyed on the resolved value being empty rather than on IsNull(), so
// both providers agree on precedence when muxed together. It builds the API client
// through the same memoised helpers.NewAPIClient constructor pkg/sdkprovider uses:
// terraform-plugin-mux calls ConfigureProvider on every underlying server for a single
// plan/apply, so without sharing that constructor this provider would double the
// RootConfigRetrieve round trip and re-run the global sentry.Init.
func (p *authentikProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data authentikProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiURL := data.URL.ValueString()
	if apiURL == "" {
		apiURL = os.Getenv("AUTHENTIK_URL")
	}
	token := data.Token.ValueString()
	if token == "" {
		token = os.Getenv("AUTHENTIK_TOKEN")
	}
	insecure := data.Insecure.ValueBool()
	if !insecure {
		insecure, _ = strconv.ParseBool(os.Getenv("AUTHENTIK_INSECURE"))
	}

	if apiURL == "" {
		resp.Diagnostics.AddError("Missing authentik URL", "no authentik URL configured, set `url` or the AUTHENTIK_URL environment variable")
		return
	}
	if token == "" {
		resp.Diagnostics.AddError("Missing authentik token", "no authentik token configured, set `token` or the AUTHENTIK_TOKEN environment variable")
		return
	}

	headers := map[string]string{}
	if !data.Headers.IsNull() {
		resp.Diagnostics.Append(data.Headers.ElementsAs(ctx, &headers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	apiClient, err := helpers.NewAPIClient(ctx, helpers.ClientOptions{
		URL:      apiURL,
		Token:    token,
		Insecure: insecure,
		Headers:  headers,
		Version:  p.version,
		Testing:  p.testing,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure authentik client", err.Error())
		return
	}

	resp.ResourceData = apiClient
	resp.DataSourceData = apiClient
}

func (p *authentikProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newGroupResource,
		newApplicationResource,
		newStageDummyResource,
		newStageUserDeleteResource,
		newStageUserLogoutResource,
		newStageDenyResource,
		newStageInvitationResource,
		newStageEndpointsResource,
		newStageSourceResource,
		newStagePromptResource,
		newStageAuthenticatorEndpointGDTCResource,
		newStageConsentResource,
		newStageAuthenticatorTOTPResource,
		newStageAuthenticatorStaticResource,
		newStageRedirectResource,
		newStageMutualTLSResource,
		newStagePasswordResource,
		newStageAccountLockdownResource,
		newStageIdentificationResource,
		newStageUserWriteResource,
		newStageAuthenticatorWebAuthnResource,
		newStageAuthenticatorSmsResource,
		newStageCaptchaResource,
		newStageEmailResource,
		newStageAuthenticatorDuoResource,
		newStageAuthenticatorEmailResource,
		newStageAuthenticatorValidateResource,
		newStageUserLoginResource,
		newStagePromptFieldResource,
		newPropertyMappingNotificationResource,
		newPropertyMappingProviderGoogleWorkspaceResource,
		newPropertyMappingProviderMicrosoftEntraResource,
		newPropertyMappingProviderRACResource,
		newPropertyMappingProviderRadiusResource,
		newPropertyMappingProviderSAMLResource,
		newPropertyMappingProviderSCIMResource,
		newPropertyMappingProviderScopeResource,
		newPropertyMappingSourceKerberosResource,
		newPropertyMappingSourceLDAPResource,
		newPropertyMappingSourceOAuthResource,
		newPropertyMappingSourcePlexResource,
		newPropertyMappingSourceSAMLResource,
		newPropertyMappingSourceSCIMResource,
	}
}

func (p *authentikProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newGroupDataSource,
	}
}
