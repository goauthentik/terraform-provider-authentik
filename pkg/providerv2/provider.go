package providerv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type authentikProvider struct {
}

func New() provider.Provider {
	return &authentikProvider{}
}

func (p *authentikProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "authentik"
}

func (p *authentikProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Required: true,
				// DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_URL", nil),
				MarkdownDescription: "The authentik API endpoint, can optionally be passed as `AUTHENTIK_URL` environmental variable",
			},
			"insecure": schema.BoolAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_INSECURE", false),
				MarkdownDescription: "Whether to skip TLS verification, can optionally be passed as `AUTHENTIK_INSECURE` environmental variable",
			},
			"token": schema.StringAttribute{
				Required: true,
				// Default:  stringdefault.StaticString(""),
				// DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_TOKEN", nil),
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

func (p *authentikProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *authentikProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// func() resource.Resource {
		// 	return resourceExample{}
		// },
	}
}

func (p *authentikProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// func() datasource.DataSource {
		// 	return dataSourceExample{}
		// },
	}
}
